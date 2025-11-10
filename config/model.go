package config

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
)

const (
	EnvironmentVariableConfigFileName = "CONFIG_FN"
	DefaultConfigFileName             = "inzibat.json"
)

type Cfg struct {
	ServerPort       int     `koanf:"serverPort"`
	Routes           []Route `koanf:"routes"`
	Concurrency      int     `koanf:"concurrency"`
	HealthCheckRoute bool    `koanf:"isHealthCheckRouteEnabled"`
}

func (cfg *Cfg) GetServerAddr() string {
	return fmt.Sprintf(":%d", cfg.ServerPort)
}

type Route struct {
	Method       string       `koanf:"method" validate:"oneof=GET,PUT,PATCH,DELETE"`
	Path         string       `koanf:"path"`
	RequestTo    RequestTo    `koanf:"requestTo"`
	FakeResponse FakeResponse `koanf:"fakeResponse"`
}

type RequestTo struct {
	Method                 string                 `koanf:"method" validate:"oneof=GET,PUT,PATCH,DELETE"`
	Headers                http.Header            `koanf:"headers"`
	Body                   map[string]interface{} `koanf:"body"`
	Host                   string                 `koanf:"host" validate:"url"`
	Path                   string                 `koanf:"path"`
	PassWithRequestBody    bool                   `koanf:"passWithRequestBody"`
	PassWithRequestHeaders bool                   `koanf:"passWithRequestHeaders"`
	InErrorReturn500       bool                   `koanf:"inErrorReturn500"`
}

func (requestTo *RequestTo) GetParsedUrl() (*url.URL, error) {
	parsedUrl, err := url.Parse(filepath.Join(requestTo.Host, requestTo.Path))
	if err != nil {
		return nil, errors.New("failed to parse url")
	}

	return parsedUrl, nil
}

type FakeResponse struct {
	Headers    map[string]string      `koanf:"headers"`
	Body       map[string]interface{} `koanf:"body"`
	BodyString string                 `koanf:"bodyString"`
	StatusCode int                    `koanf:"statusCode"`
}
