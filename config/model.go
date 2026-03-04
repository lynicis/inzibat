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
	ServerPort       int                   `json:"serverPort" koanf:"serverPort" validate:"required"`
	Routes           []Route               `json:"routes" koanf:"routes" validate:"required,gt=0,dive,required"`
	Concurrency      int                   `json:"concurrency" koanf:"concurrency"`
	HealthCheckRoute bool                  `json:"healthCheckRoute" koanf:"isHealthCheckRouteEnabled"`
	CircuitBreaker   *CircuitBreakerConfig `json:"circuitBreaker,omitempty" koanf:"circuitBreaker"`
}

func (cfg *Cfg) GetServerAddr() string {
	return fmt.Sprintf(":%d", cfg.ServerPort)
}

type Route struct {
	Method       string        `json:"method" koanf:"method" validate:"oneof=GET POST PUT PATCH DELETE"`
	Path         string        `json:"path" koanf:"path" validate:"required,startswith=/"`
	RequestTo    *RequestTo    `json:"requestTo,omitempty" koanf:"requestTo" validate:"required_without=FakeResponse"`
	FakeResponse *FakeResponse `json:"fakeResponse,omitempty" koanf:"fakeResponse" validate:"required_without=RequestTo"`
}

func (cfg *Cfg) ConvertRoutesTuiTable() [][]string {
	var rows [][]string
	for _, route := range cfg.Routes {
		routeType := "UNKNOWN"

		if route.FakeResponse != nil {
			routeType = "MOCK"
		}
		if route.RequestTo != nil {
			routeType = "PROXY"
		}

		rows = append(rows, []string{
			route.Method,
			route.Path,
			routeType,
		})
	}

	return rows
}

type RequestTo struct {
	Method                 string                `json:"method" koanf:"method" validate:"oneof=GET POST PUT PATCH DELETE"`
	Headers                http.Header           `json:"headers" koanf:"headers"`
	Body                   HttpBody              `json:"body,omitempty" koanf:"body"`
	Host                   string                `json:"host" koanf:"host" validate:"url"`
	Path                   string                `json:"path" koanf:"path" validate:"required,startswith=/"`
	PassWithRequestBody    bool                  `json:"passWithRequestBody,omitempty" koanf:"passWithRequestBody"`
	PassWithRequestHeaders bool                  `json:"passWithRequestHeaders,omitempty" koanf:"passWithRequestHeaders"`
	InErrorReturn500       bool                  `json:"inErrorReturn500,omitempty" koanf:"inErrorReturn500"`
	CircuitBreaker         *CircuitBreakerConfig `json:"circuitBreaker,omitempty" koanf:"circuitBreaker"`
}

type CircuitBreakerConfig struct {
	Enabled             *bool `json:"enabled,omitempty" koanf:"enabled"`
	FailureThreshold    int   `json:"failureThreshold,omitempty" koanf:"failureThreshold" validate:"omitempty,gt=0"`
	MinimumRequests     int   `json:"minimumRequests,omitempty" koanf:"minimumRequests" validate:"omitempty,gt=0"`
	OpenTimeoutMs       int   `json:"openTimeoutMs,omitempty" koanf:"openTimeoutMs" validate:"omitempty,gt=0"`
	HalfOpenMaxRequests int   `json:"halfOpenMaxRequests,omitempty" koanf:"halfOpenMaxRequests" validate:"omitempty,gt=0"`
	SuccessThreshold    int   `json:"successThreshold,omitempty" koanf:"successThreshold" validate:"omitempty,gt=0"`
}

func BoolPointer(value bool) *bool {
	boolValue := value
	return &boolValue
}

func DefaultCircuitBreakerConfig() CircuitBreakerConfig {
	return CircuitBreakerConfig{
		Enabled:             BoolPointer(false),
		FailureThreshold:    5,
		MinimumRequests:     10,
		OpenTimeoutMs:       30000,
		HalfOpenMaxRequests: 2,
		SuccessThreshold:    2,
	}
}

func MergeCircuitBreakerConfig(base *CircuitBreakerConfig, override *CircuitBreakerConfig) *CircuitBreakerConfig {
	if base == nil && override == nil {
		return nil
	}

	mergedConfig := DefaultCircuitBreakerConfig()
	if base != nil {
		mergedConfig = applyCircuitBreakerConfig(mergedConfig, *base)
	}
	if override != nil {
		mergedConfig = applyCircuitBreakerConfig(mergedConfig, *override)
	}

	return &mergedConfig
}

func applyCircuitBreakerConfig(destination CircuitBreakerConfig, source CircuitBreakerConfig) CircuitBreakerConfig {
	if source.Enabled != nil {
		destination.Enabled = BoolPointer(*source.Enabled)
	}

	if source.FailureThreshold > 0 {
		destination.FailureThreshold = source.FailureThreshold
	}
	if source.MinimumRequests > 0 {
		destination.MinimumRequests = source.MinimumRequests
	}
	if source.OpenTimeoutMs > 0 {
		destination.OpenTimeoutMs = source.OpenTimeoutMs
	}
	if source.HalfOpenMaxRequests > 0 {
		destination.HalfOpenMaxRequests = source.HalfOpenMaxRequests
	}
	if source.SuccessThreshold > 0 {
		destination.SuccessThreshold = source.SuccessThreshold
	}

	return destination
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
