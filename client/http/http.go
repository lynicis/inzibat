package http

import (
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

type Client interface {
	Get(uri string, requestHeader map[string]string) (*Response, error)
	Post(uri string, requestHeader map[string]string, requestBody []byte) (*Response, error)
	Put(uri string, requestHeader map[string]string, requestBody []byte) (*Response, error)
	Patch(uri string, requestHeader map[string]string, requestBody []byte) (*Response, error)
	Delete(uri string, requestHeader map[string]string, requestBody []byte) (*Response, error)
}

type HttpClient struct {
	client       *fasthttp.Client
	requestPool  sync.Pool
	responsePool sync.Pool
}

func NewHttpClient() *HttpClient {
	return &HttpClient{
		client: &fasthttp.Client{
			ReadTimeout:                   10 * time.Second,
			WriteTimeout:                  10 * time.Second,
			MaxConnsPerHost:               1000,
			MaxIdleConnDuration:           time.Minute,
			MaxConnDuration:               time.Minute * 5,
			NoDefaultUserAgentHeader:      true,
			DisableHeaderNamesNormalizing: true,
			DisablePathNormalizing:        true,
		},
		requestPool: sync.Pool{
			New: func() interface{} {
				return &fasthttp.Request{}
			},
		},
		responsePool: sync.Pool{
			New: func() interface{} {
				return &fasthttp.Response{}
			},
		},
	}
}

func (httpClient *HttpClient) Get(
	uri string,
	requestHeader map[string]string,
) (*Response, error) {
	return httpClient.makeRequest(uri, http.MethodGet, requestHeader, nil)
}

func (httpClient *HttpClient) Post(
	uri string,
	requestHeader map[string]string,
	requestBody []byte,
) (*Response, error) {
	return httpClient.makeRequest(uri, http.MethodPost, requestHeader, requestBody)
}

func (httpClient *HttpClient) Put(
	uri string,
	requestHeader map[string]string,
	requestBody []byte,
) (*Response, error) {
	return httpClient.makeRequest(uri, http.MethodPut, requestHeader, requestBody)
}

func (httpClient *HttpClient) Patch(
	uri string,
	requestHeader map[string]string,
	requestBody []byte,
) (*Response, error) {
	return httpClient.makeRequest(uri, http.MethodPatch, requestHeader, requestBody)
}

func (httpClient *HttpClient) Delete(
	uri string,
	requestHeader map[string]string,
	requestBody []byte,
) (*Response, error) {
	return httpClient.makeRequest(uri, http.MethodDelete, requestHeader, requestBody)
}

func (httpClient *HttpClient) makeRequest(
	uri string,
	method string,
	requestHeader map[string]string,
	requestBody []byte,
) (*Response, error) {
	req := httpClient.requestPool.Get().(*fasthttp.Request)
	defer httpClient.requestPool.Put(req)
	req.Reset()

	req.SetRequestURI(uri)
	req.Header.SetMethod(method)
	req.Header.SetContentType(fiber.MIMEApplicationJSON)

	if len(requestHeader) > 0 {
		for headerKey, headerValue := range requestHeader {
			req.Header.Set(headerKey, headerValue)
		}
	}

	if requestBody != nil {
		req.SetBody(requestBody)
	}

	resp := httpClient.responsePool.Get().(*fasthttp.Response)
	defer httpClient.responsePool.Put(resp)
	resp.Reset()

	err := httpClient.client.Do(req, resp)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() >= http.StatusMultipleChoices {
		return nil, errors.New("response failed")
	}

	return &Response{
		Status: resp.StatusCode(),
		Body:   append([]byte(nil), resp.Body()...),
	}, nil
}
