package client

import (
	"errors"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

type Client interface {
	Get(uri string, requestHeader HttpHeader) (*HttpResponse, error)
	Post(uri string, requestHeader HttpHeader, requestBody []byte) (*HttpResponse, error)
	Put(uri string, requestHeader HttpHeader, requestBody []byte) (*HttpResponse, error)
	Patch(uri string, requestHeader HttpHeader, requestBody []byte) (*HttpResponse, error)
	Delete(uri string, requestHeader HttpHeader, requestBody []byte) (*HttpResponse, error)
	GetCloneOfStruct() *client
}

type client struct {
	httpClient *fasthttp.Client
}

func NewClient() Client {
	httpClient := &fasthttp.Client{
		ReadTimeout:                   10 * time.Second,
		WriteTimeout:                  10 * time.Second,
		MaxIdleConnDuration:           10 * time.Second,
		NoDefaultUserAgentHeader:      true,
		DisableHeaderNamesNormalizing: true,
		DisablePathNormalizing:        true,
	}

	return &client{
		httpClient: httpClient,
	}
}

func (c client) Get(uri string, requestHeader HttpHeader) (*HttpResponse, error) {
	response, err := c.makeRequest(uri, http.MethodGet, requestHeader, nil)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c client) Post(uri string, requestHeader HttpHeader, requestBody []byte) (*HttpResponse, error) {
	response, err := c.makeRequest(uri, http.MethodPost, requestHeader, requestBody)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c client) Put(uri string, requestHeader HttpHeader, requestBody []byte) (*HttpResponse, error) {
	response, err := c.makeRequest(uri, http.MethodPut, requestHeader, requestBody)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c client) Patch(uri string, requestHeader HttpHeader, requestBody []byte) (*HttpResponse, error) {
	response, err := c.makeRequest(uri, http.MethodPatch, requestHeader, requestBody)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c client) Delete(uri string, requestHeader HttpHeader, requestBody []byte) (*HttpResponse, error) {
	response, err := c.makeRequest(uri, http.MethodDelete, requestHeader, requestBody)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c client) GetCloneOfStruct() *client {
	copOfClientStruct := c
	return &copOfClientStruct
}

func (c client) makeRequest(uri string, method string, requestHeader HttpHeader, requestBody []byte) (*HttpResponse, error) {
	request := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(request)

	request.SetRequestURI(uri)
	request.Header.SetMethod(method)
	request.Header.SetContentType(fiber.MIMEApplicationJSON)

	if len(requestHeader) > 0 {
		for headerKey, headerValue := range requestHeader {
			request.Header.Set(headerKey, headerValue)
		}
	}

	if requestBody != nil {
		request.SetBody(requestBody)
	}

	response := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(response)

	err := fasthttp.Do(request, response)
	if err != nil {
		return nil, err
	}

	if response.StatusCode() >= http.StatusMultipleChoices {
		return nil, errors.New(ResponseFailed)
	}

	return &HttpResponse{
		Status: response.StatusCode(),
		Body:   response.Body(),
	}, nil
}
