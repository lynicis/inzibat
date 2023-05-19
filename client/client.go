package client

import (
	"errors"
	"net/http"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

const (
	ResponseFailed = "response failed"
)

type Client interface {
	Get(uri string, requestHeader HttpHeader) (*HttpResponse, error)
	Post(uri string, requestHeader HttpHeader, requestBody HttpBody) (*HttpResponse, error)
	Put(uri string, requestHeader HttpHeader, requestBody HttpBody) (*HttpResponse, error)
	Patch(uri string, requestHeader HttpHeader, requestBody HttpBody) (*HttpResponse, error)
	Delete(uri string, requestHeader HttpHeader, requestBody HttpBody) (*HttpResponse, error)
	GetCloneOfStruct() *client
}

type client struct {
	fasthttp *fasthttp.Client
}

func NewClient() Client {
	fasthttpInstance := &fasthttp.Client{
		ReadTimeout:              10 * time.Second,
		WriteTimeout:             10 * time.Second,
		MaxIdleConnDuration:      10 * time.Second,
		NoDefaultUserAgentHeader: true,
	}

	return &client{
		fasthttp: fasthttpInstance,
	}
}

func (c *client) Get(uri string, requestHeader HttpHeader) (*HttpResponse, error) {
	response, err := makeRequest(uri, http.MethodGet, requestHeader)
	if err != nil {
		return nil, err
	}

	return &HttpResponse{
		Status: response.StatusCode(),
		Body:   response.Body(),
	}, nil
}

func (c *client) Post(uri string, requestHeader HttpHeader, requestBody HttpBody) (*HttpResponse, error) {
	response, err := makeRequest(uri, http.MethodPost, requestHeader, requestBody)
	if err != nil {
		return nil, err
	}

	return &HttpResponse{
		Status: response.StatusCode(),
		Body:   response.Body(),
	}, nil
}

func (c *client) Put(uri string, requestHeader HttpHeader, requestBody HttpBody) (*HttpResponse, error) {
	response, err := makeRequest(uri, http.MethodPut, requestHeader, requestBody)
	if err != nil {
		return nil, err
	}

	return &HttpResponse{
		Status: response.StatusCode(),
		Body:   response.Body(),
	}, nil
}

func (c *client) Patch(uri string, requestHeader HttpHeader, requestBody HttpBody) (*HttpResponse, error) {
	response, err := makeRequest(uri, http.MethodPatch, requestHeader, requestBody)
	if err != nil {
		return nil, err
	}

	return &HttpResponse{
		Status: response.StatusCode(),
		Body:   response.Body(),
	}, nil
}

func (c *client) Delete(uri string, requestHeader HttpHeader, requestBody HttpBody) (*HttpResponse, error) {
	response, err := makeRequest(uri, http.MethodDelete, requestHeader, requestBody)
	if err != nil {
		return nil, err
	}

	return &HttpResponse{
		Status: response.StatusCode(),
		Body:   response.Body(),
	}, nil
}

func makeRequest(uri string, method string, requestHeader HttpHeader, requestBody ...HttpBody) (*fiber.Response, error) {
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
		reqBody, err := json.Marshal(requestBody)
		if err != nil {
			return nil, err
		}
		request.SetBody(reqBody)
	}

	response := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(response)

	err := fasthttp.Do(request, response)
	if err != nil {
		return nil, err
	}

	if response.StatusCode() >= fasthttp.StatusMultipleChoices {
		return nil, errors.New(ResponseFailed)
	}

	var responseBody map[string]interface{}
	responseBodyBytes := response.Body()
	if len(responseBodyBytes) > 0 {
		err = json.Unmarshal(responseBodyBytes, &responseBody)
		if err != nil {
			return nil, err
		}
	}

	return response, nil
}

func (c *client) GetCloneOfStruct() *client {
	copOfClientStruct := c
	return copOfClientStruct
}
