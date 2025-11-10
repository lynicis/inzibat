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

type Client struct {
	client *fasthttp.Client
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
	}
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

// TODO: implement retry mechanism
func (httpClient *Client) makeRequest(
	uri string,
	method string,
	requestHeader http.Header,
	requestBody []byte,
) (*Response, error) {
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
	if err := httpClient.client.Do(req, resp); err != nil {
		return nil, err
	}

	if resp.StatusCode() >= http.StatusMultipleChoices {
		return nil, errors.New("response failed")
	}

	return &Response{
		Status: resp.StatusCode(),
		Body:   resp.Body(),
	}, nil
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
