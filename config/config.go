package config

import (
	"net/http"
	"runtime"
)

type Reader struct {
	ConfigReader ReaderStrategy
}

func (reader *Reader) Read(filename string) (*Cfg, error) {
	config, err := reader.ConfigReader.Read(filename)
	if err != nil {
		return nil, err
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
