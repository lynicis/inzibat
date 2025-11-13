package config

import (
	"net/http"
	"runtime"

	"github.com/go-playground/validator/v10"
)

type Reader struct {
	ConfigReader ReaderStrategy
	Validator    *validator.Validate
}

func (reader *Reader) Read(filename string) (*Cfg, error) {
	config, err := reader.ConfigReader.Read(filename)
	if err != nil {
		return nil, err
	}

	if reader.Validator != nil {
		if err = reader.Validator.Struct(config); err != nil {
			return nil, err
		}
	}

	for routeIndex, route := range config.Routes {
		var (
			RequestToMethod = route.RequestTo.Method
			RequestToBody   = route.RequestTo.Body
		)

		if RequestToMethod == "" {
			route.RequestTo.Method = http.MethodGet
		}

		if RequestToMethod == http.MethodGet && RequestToBody != nil {
			return nil, ErrorGetSendBody
		}

		config.Routes[routeIndex].Method = route.Method
		config.Routes[routeIndex].RequestTo.Method = route.RequestTo.Method
	}

	if config.HealthCheckRoute {
		config.Routes = append(
			config.Routes,
			Route{
				Method: "GET",
				Path:   "/health",
				FakeResponse: FakeResponse{
					StatusCode: http.StatusOK,
				},
			},
		)
	}

	if config.Concurrency == 0 {
		config.Concurrency = runtime.GOMAXPROCS(3)
	}

	return config, nil
}
