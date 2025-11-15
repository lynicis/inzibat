package http

import (
	"errors"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

type RetryConfig struct {
	MaxRetries      int
	InitialBackoff time.Duration
	MaxBackoff      time.Duration
	BackoffMultiplier float64
}

func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries:      3,
		InitialBackoff:  100 * time.Millisecond,
		MaxBackoff:      2 * time.Second,
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

func (httpClient *Client) Get(
	uri string,
	requestHeader http.Header,
) (*Response, error) {
	return httpClient.makeRequest(uri, http.MethodGet, requestHeader, nil)
}

func (httpClient *Client) Post(
	uri string,
	requestHeader http.Header,
	requestBody []byte,
) (*Response, error) {
	return httpClient.makeRequest(uri, http.MethodPost, requestHeader, requestBody)
}

func (httpClient *Client) Put(
	uri string,
	requestHeader http.Header,
	requestBody []byte,
) (*Response, error) {
	return httpClient.makeRequest(uri, http.MethodPut, requestHeader, requestBody)
}

func (httpClient *Client) Patch(
	uri string,
	requestHeader http.Header,
	requestBody []byte,
) (*Response, error) {
	return httpClient.makeRequest(uri, http.MethodPatch, requestHeader, requestBody)
}

func (httpClient *Client) Delete(
	uri string,
	requestHeader http.Header,
	requestBody []byte,
) (*Response, error) {
	return httpClient.makeRequest(uri, http.MethodDelete, requestHeader, requestBody)
}

func isRetryableError(err error, statusCode int) bool {
	if err != nil {
		return true
	}
	return statusCode >= http.StatusInternalServerError && statusCode < 600
}

func (httpClient *Client) calculateBackoff(attempt int) time.Duration {
	backoff := time.Duration(float64(httpClient.retryConfig.InitialBackoff) * 
		pow(httpClient.retryConfig.BackoffMultiplier, float64(attempt)))
	if backoff > httpClient.retryConfig.MaxBackoff {
		backoff = httpClient.retryConfig.MaxBackoff
	}
	return backoff
}

func pow(base, exp float64) float64 {
	result := 1.0
	for i := 0; i < int(exp); i++ {
		result *= base
	}
	return result
}

func (httpClient *Client) makeRequest(
	uri string,
	method string,
	requestHeader http.Header,
	requestBody []byte,
) (*Response, error) {
	var lastErr error
	var lastResponse *Response

	for attempt := 0; attempt <= httpClient.retryConfig.MaxRetries; attempt++ {
		if attempt > 0 {
			backoff := httpClient.calculateBackoff(attempt - 1)
			time.Sleep(backoff)
		}

		req := fasthttp.AcquireRequest()
		req.SetRequestURI(uri)
		req.Header.SetMethod(method)
		req.Header.SetContentType(fiber.MIMEApplicationJSON)

		if len(requestHeader) > 0 {
			for headerKey, headerValue := range requestHeader {
				req.Header.Set(headerKey, strings.Join(headerValue, ""))
			}
		}

		if requestBody != nil {
			req.SetBody(requestBody)
		}

		resp := fasthttp.AcquireResponse()
		err := httpClient.client.Do(req, resp)
		
		if err != nil {
			lastErr = err
			fasthttp.ReleaseRequest(req)
			fasthttp.ReleaseResponse(resp)
			
			if isRetryableError(err, 0) && attempt < httpClient.retryConfig.MaxRetries {
				continue
			}
			return nil, err
		}

		statusCode := resp.StatusCode()
		
		if statusCode >= http.StatusMultipleChoices {
			lastErr = errors.New("response failed")
			fasthttp.ReleaseRequest(req)
			fasthttp.ReleaseResponse(resp)
			
			if isRetryableError(nil, statusCode) && attempt < httpClient.retryConfig.MaxRetries {
				continue
			}
			return nil, lastErr
		}

		body := make([]byte, len(resp.Body()))
		copy(body, resp.Body())
		
		lastResponse = &Response{
			Status: statusCode,
			Body:   body,
		}

		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(resp)
		break
	}

	if lastResponse != nil {
		return lastResponse, nil
	}

	if lastErr != nil {
		return nil, lastErr
	}

	return nil, errors.New("request failed after all retries")
}

func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	tcpListener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer tcpListener.Close()

	return tcpListener.Addr().(*net.TCPAddr).Port, nil
}
