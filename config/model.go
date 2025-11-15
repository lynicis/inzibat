package config

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

const (
	EnvironmentVariableConfigFileName = "CONFIG_FN"
	DefaultConfigFileName             = "inzibat.json"
	GlobalConfigFileName              = ".inzibat.config.json"
)

type Cfg struct {
	ServerPort       int     `json:"serverPort" koanf:"serverPort" validate:"required"`
	Routes           []Route `json:"routes" koanf:"routes" validate:"required,gt=0,dive,required"`
	Concurrency      int     `json:"concurrency" koanf:"concurrency"`
	HealthCheckRoute bool    `json:"healthCheckRoute" koanf:"isHealthCheckRouteEnabled"`
}

func (cfg *Cfg) GetServerAddr() string {
	return fmt.Sprintf(":%d", cfg.ServerPort)
}

type Route struct {
	Method       string       `json:"method" koanf:"method" validate:"oneof=GET PUT PATCH DELETE"`
	Path         string       `json:"path" koanf:"path" validate:"required,startswith=/"`
	RequestTo    RequestTo    `json:"requestTo,omitempty" koanf:"requestTo"`
	FakeResponse FakeResponse `json:"fakeResponse,omitempty" koanf:"fakeResponse"`
}

type RequestTo struct {
	Method                 string      `json:"method" koanf:"method" validate:"oneof=GET PUT PATCH DELETE"`
	Headers                http.Header `json:"headers" koanf:"headers"`
	Body                   HttpBody    `json:"body,omitempty" koanf:"body"`
	Host                   string      `json:"host" koanf:"host" validate:"url"`
	Path                   string      `json:"path" koanf:"path" validate:"required,startswith=/"`
	PassWithRequestBody    bool        `json:"passWithRequestBody,omitempty" koanf:"passWithRequestBody"`
	PassWithRequestHeaders bool        `json:"passWithRequestHeaders,omitempty" koanf:"passWithRequestHeaders"`
	InErrorReturn500       bool        `json:"inErrorReturn500,omitempty" koanf:"inErrorReturn500"`
}

func (requestTo *RequestTo) GetParsedUrl() (*url.URL, error) {
	parsedUrl, err := url.Parse(requestTo.Host + requestTo.Path)
	if err != nil {
		return nil, errors.New("failed to parse url")
	}

	return parsedUrl, nil
}

type HttpBody map[string]any

type FakeResponse struct {
	Headers    http.Header `json:"headers" koanf:"headers"`
	Body       HttpBody    `json:"body,omitempty" koanf:"body" validate:"required_without=BodyString"`
	BodyString string      `json:"bodyString,omitempty" koanf:"bodyString" validate:"required_without=Body"`
	StatusCode int         `json:"statusCode" koanf:"statusCode" validate:"required"`
}
