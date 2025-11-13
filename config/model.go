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
)

type Cfg struct {
	ServerPort       int     `koanf:"serverPort" validate:"required"`
	Routes           []Route `koanf:"routes" validate:"required,gt=0,dive,required"`
	Concurrency      int     `koanf:"concurrency"`
	HealthCheckRoute bool    `koanf:"isHealthCheckRouteEnabled"`
}

func (cfg *Cfg) GetServerAddr() string {
	return fmt.Sprintf(":%d", cfg.ServerPort)
}

type Route struct {
	Method       string       `koanf:"method" validate:"oneof=GET PUT PATCH DELETE"`
	Path         string       `koanf:"path" validate:"required,startswith=/"`
	RequestTo    RequestTo    `koanf:"requestTo"`
	FakeResponse FakeResponse `koanf:"fakeResponse"`
}

type RequestTo struct {
	Method                 string      `koanf:"method" validate:"oneof=GET PUT PATCH DELETE"`
	Headers                http.Header `koanf:"headers"`
	Body                   HttpBody    `koanf:"body"`
	Host                   string      `koanf:"host" validate:"url"`
	Path                   string      `koanf:"path" validate:"required,startswith=/"`
	PassWithRequestBody    bool        `koanf:"passWithRequestBody"`
	PassWithRequestHeaders bool        `koanf:"passWithRequestHeaders"`
	InErrorReturn500       bool        `koanf:"inErrorReturn500"`
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
	Headers    http.Header `koanf:"headers"`
	Body       HttpBody    `koanf:"body" validate:"required_without=BodyString"`
	BodyString string      `koanf:"bodyString" validate:"required_without=Body"`
	StatusCode int         `koanf:"statusCode" validate:"required"`
}
