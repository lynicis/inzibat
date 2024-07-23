package config

import (
	"errors"
	"net/http"
	"runtime"
)

type Config interface {
	LoadConfig(filename string) (*Cfg, error)
}

type Loader struct {
	ConfigReader Reader
}

func (loader *Loader) LoadConfig(filename string) (*Cfg, error) {
	config, err := loader.ConfigReader.ReadConfig(filename)
	if err != nil {
		return nil, err
	}

	for indexOfRoute, route := range config.Routes {
		var (
			RequestToMethod = route.RequestTo.Method
			RequestToBody   = route.RequestTo.Body
		)

		if RequestToMethod == "" {
			route.RequestTo.Method = http.MethodGet
		}

		if RequestToMethod == http.MethodGet && RequestToBody != nil {
			return nil, errors.New(ErrorGetSendBody)
		}

		config.Routes[indexOfRoute].Method = route.Method
		config.Routes[indexOfRoute].RequestTo.Method = route.RequestTo.Method
	}

	if config.HealthCheckRoute {
		config.Routes = append(
			config.Routes,
			Route{
				Method: "GET",
				Path:   "/health",
				Mock: Mock{
					StatusCode: http.StatusOK,
				},
			},
		)
	}

	if config.Concurrency.RouteCreatorLimit == 0 {
		config.Concurrency.RouteCreatorLimit = runtime.GOMAXPROCS(3)
	}

	return config, nil
}
