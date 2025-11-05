package config

import (
	"errors"
	"net/url"
	"path/filepath"
)

const (
	EnvironmentVariableConfigFileName = "CONFIG_FN"
	DefaultConfigFileName             = "inzibat.config.json"
)

type Cfg struct {
	ServerPort       int
	Routes           []Route
	Concurrency      Concurrency
	HealthCheckRoute bool
}

type Route struct {
	Method       string
	Path         string
	RequestTo    RequestTo
	FakeResponse FakeResponse
}

type Concurrency struct {
	RouteCreatorLimit int
}

type RequestTo struct {
	Method                 string
	Headers                map[string][]string
	Body                   map[string]interface{}
	Host                   string
	Path                   string
	PassWithRequestBody    bool
	PassWithRequestHeaders bool
	InErrorReturn500       bool
}

func (requestTo *RequestTo) GetParsedUrl() (*url.URL, error) {
	parsedUrl, err := url.Parse(filepath.Join(requestTo.Host, requestTo.Path))
	if err != nil {
		return nil, errors.New("failed to parse url")
	}

	return parsedUrl, nil
}

type FakeResponse struct {
	Headers    map[string]string
	Body       map[string]interface{}
	BodyString string
	StatusCode int
}
