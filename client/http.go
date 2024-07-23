package client

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

type Client interface {
	Get(uri string, requestHeader map[string]string) (*HttpResponse, error)
	Post(uri string, requestHeader map[string]string, requestBody []byte) (*HttpResponse, error)
	Put(uri string, requestHeader map[string]string, requestBody []byte) (*HttpResponse, error)
	Patch(uri string, requestHeader map[string]string, requestBody []byte) (*HttpResponse, error)
	Delete(uri string, requestHeader map[string]string, requestBody []byte) (*HttpResponse, error)
}

type HttpClient struct {
	FasthttpClient *fasthttp.Client
}

func (httpClient *HttpClient) Get(
	uri string,
	requestHeader map[string]string,
) (*HttpResponse, error) {
	return httpClient.makeRequest(uri, http.MethodGet, requestHeader, nil)
}

func (httpClient *HttpClient) Post(
	uri string,
	requestHeader map[string]string,
	requestBody []byte,
) (*HttpResponse, error) {
	return httpClient.makeRequest(uri, http.MethodPost, requestHeader, requestBody)
}

func (httpClient *HttpClient) Put(
	uri string,
	requestHeader map[string]string,
	requestBody []byte,
) (*HttpResponse, error) {
	return httpClient.makeRequest(uri, http.MethodPut, requestHeader, requestBody)
}

func (httpClient *HttpClient) Patch(
	uri string,
	requestHeader map[string]string,
	requestBody []byte,
) (*HttpResponse, error) {
	return httpClient.makeRequest(uri, http.MethodPatch, requestHeader, requestBody)
}

func (httpClient *HttpClient) Delete(
	uri string,
	requestHeader map[string]string,
	requestBody []byte,
) (*HttpResponse, error) {
	return httpClient.makeRequest(uri, http.MethodDelete, requestHeader, requestBody)
}

func (httpClient *HttpClient) makeRequest(
	uri string,
	method string,
	requestHeader map[string]string,
	requestBody []byte,
) (*HttpResponse, error) {
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
		return nil, errors.New(ErrorResponseFailed)
	}

	return &HttpResponse{
		Status: response.StatusCode(),
		Body:   response.Body(),
	}, nil
}
