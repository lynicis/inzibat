package http

import (
	"errors"
	"math"
	"math/rand/v2"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

type RetryConfig struct {
	MaxRetries        int
	InitialBackoff    time.Duration
	MaxBackoff        time.Duration
	BackoffMultiplier float64
}

func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries:        3,
		InitialBackoff:    100 * time.Millisecond,
		MaxBackoff:        2 * time.Second,
		BackoffMultiplier: 2.0,
	}
}

type Client struct {
	client      *fasthttp.Client
	retryConfig RetryConfig
}

func NewHttpClient() *Client {
	return &Client{
		client: &fasthttp.Client{
			ReadTimeout:                   10 * time.Second,
			WriteTimeout:                  10 * time.Second,
			MaxConnsPerHost:               1000,
			MaxIdleConnDuration:           time.Minute,
			MaxConnDuration:               time.Minute * 5,
			NoDefaultUserAgentHeader:      true,
			DisableHeaderNamesNormalizing: true,
			DisablePathNormalizing:        true,
			Dial: (&fasthttp.TCPDialer{
				Concurrency:      4096,
				DNSCacheDuration: time.Hour,
			}).Dial,
		},
		retryConfig: DefaultRetryConfig(),
	}
}

func (httpClient *Client) SetRetryConfig(config RetryConfig) {
	httpClient.retryConfig = config
}

func (httpClient *Client) Do(
	method, uri string,
	requestHeader http.Header,
	requestBody []byte,
) (*Response, error) {
	return httpClient.makeRequest(uri, method, requestHeader, requestBody)
}

func isRetryableError(err error, statusCode int) bool {
	if err != nil {
		return true
	}
	return statusCode >= http.StatusInternalServerError && statusCode < 600
}

func (httpClient *Client) calculateBackoff(attempt int) time.Duration {
	backoff := time.Duration(float64(httpClient.retryConfig.InitialBackoff) *
		math.Pow(httpClient.retryConfig.BackoffMultiplier, float64(attempt)))

	// ±20% jitter keeps retries from synchronizing across clients.
	backoff = time.Duration(float64(backoff) * (0.8 + 0.4*rand.Float64()))

	if backoff > httpClient.retryConfig.MaxBackoff {
		backoff = httpClient.retryConfig.MaxBackoff
	}

	return backoff
}

func (httpClient *Client) buildRequest(
	uri, method string,
	requestHeader http.Header,
	requestBody []byte,
) *fasthttp.Request {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(uri)
	req.Header.SetMethod(method)
	req.Header.SetContentType(fiber.MIMEApplicationJSON)

	if len(requestHeader) > 0 {
		for headerKey, headerValues := range requestHeader {
			if IsUnforwardableHeader(headerKey) {
				continue
			}

			for _, headerValue := range headerValues {
				req.Header.Add(headerKey, headerValue)
			}
		}
	}

	if requestBody != nil {
		req.SetBody(requestBody)
	}

	return req
}

// IsUnforwardableHeader reports headers the transport owns or that only apply to one hop.
func IsUnforwardableHeader(name string) bool {
	switch http.CanonicalHeaderKey(name) {
	case "Host", "Content-Length", "Transfer-Encoding", "Connection", "Keep-Alive",
		"Upgrade", "Proxy-Connection", "Te", "Trailer", "Proxy-Authorization":
		return true
	default:
		return false
	}
}

func (httpClient *Client) executeRequest(req *fasthttp.Request) (*fasthttp.Response, error) {
	resp := fasthttp.AcquireResponse()
	err := httpClient.client.Do(req, resp)
	if err != nil {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(resp)
		return nil, err
	}
	return resp, nil
}

func (httpClient *Client) handleResponse(
	resp *fasthttp.Response,
	req *fasthttp.Request,
) (*Response, int, error) {
	statusCode := resp.StatusCode()

	if statusCode >= http.StatusInternalServerError {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(resp)
		return nil, statusCode, errors.New("response failed")
	}

	body := make([]byte, len(resp.Body()))
	copy(body, resp.Body())

	fasthttp.ReleaseRequest(req)
	fasthttp.ReleaseResponse(resp)

	return &Response{
		Status: statusCode,
		Body:   body,
	}, statusCode, nil
}

func (httpClient *Client) shouldRetry(method string, err error, statusCode, attempt int) bool {
	if !isIdempotent(method) {
		return false
	}

	return isRetryableError(err, statusCode) && attempt < httpClient.retryConfig.MaxRetries
}

func isIdempotent(method string) bool {
	switch method {
	case http.MethodGet, http.MethodPut, http.MethodDelete:
		return true
	default:
		return false
	}
}

func (httpClient *Client) makeRequest(
	uri string,
	method string,
	requestHeader http.Header,
	requestBody []byte,
) (*Response, error) {
	var lastErr error

	for attempt := 0; attempt <= httpClient.retryConfig.MaxRetries; attempt++ {
		if attempt > 0 {
			backoff := httpClient.calculateBackoff(attempt - 1)
			time.Sleep(backoff)
		}

		req := httpClient.buildRequest(uri, method, requestHeader, requestBody)
		resp, err := httpClient.executeRequest(req)

		if err != nil {
			lastErr = err
			if httpClient.shouldRetry(method, err, 0, attempt) {
				continue
			}
			return nil, err
		}

		response, statusCode, err := httpClient.handleResponse(resp, req)
		if err != nil {
			lastErr = err
			if httpClient.shouldRetry(method, nil, statusCode, attempt) {
				continue
			}
			return nil, lastErr
		}

		return response, nil
	}

	if lastErr != nil {
		return nil, lastErr
	}

	return nil, errors.New("request failed after all retries")
}
