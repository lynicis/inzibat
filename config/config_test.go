package config

import (
	"errors"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestLoader_LoadConfig(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		expectedCfg := &Cfg{
			ServerPort: 8080,
			Routes: []Route{
				{
					Method: fiber.MethodGet,
					Path:   "/route-one",
					RequestTo: RequestTo{
						Method: http.MethodPost,
						Headers: map[string][]string{
							"X-Test-Header": {"Test-Header-Value"},
						},
						Body: map[string]interface{}{
							"testKey": "testValue",
						},
						Host:                   "http://localhost:8081",
						Path:                   "/route-one",
						PassWithRequestBody:    true,
						PassWithRequestHeaders: true,
					},
				},
				{
					Method: fiber.MethodGet,
					Path:   "/route-two",
					RequestTo: RequestTo{
						Method: http.MethodGet,
						Headers: map[string][]string{
							"X-Test-Header": {"Test-Header-Value"},
						},
						Host:                   "http://localhost:8081",
						Path:                   "/route-two",
						PassWithRequestBody:    true,
						PassWithRequestHeaders: true,
					},
				},
			},
			Concurrency: Concurrency{
				RouteCreatorLimit: 5,
			},
		}

		mockReader := &MockReader{
			OutputCfg:   expectedCfg,
			OutputError: nil,
		}

		cfgLoader := &Reader{
			ConfigReader: mockReader,
		}
		cfg, err := cfgLoader.Read("test-file-name")

		assert.NoError(t, err)
		assert.Equal(t, expectedCfg, cfg)
	})

	t.Run("when reader return error should return it", func(t *testing.T) {
		expectedError := errors.New("something went wrong")
		mockReader := &MockReader{
			OutputCfg:   nil,
			OutputError: expectedError,
		}

		cfgLoader := &Reader{
			ConfigReader: mockReader,
		}
		cfg, err := cfgLoader.Read("test-file-name")

		assert.Nil(t, cfg)
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})

	t.Run("against healthcheck route", func(t *testing.T) {
		configLoader := &Reader{
			ConfigReader: &MockReader{
				OutputCfg: &Cfg{
					HealthCheckRoute: true,
				},
				OutputError: nil,
			},
		}
		cfg, err := configLoader.Read("")

		assert.NoError(t, err)
		assert.NotNil(t, cfg)
	})

	t.Run("against concurrency route creator limit", func(t *testing.T) {
		configLoader := &Reader{
			ConfigReader: &MockReader{
				OutputCfg: &Cfg{
					Concurrency: Concurrency{
						RouteCreatorLimit: 0,
					},
				},
				OutputError: nil,
			},
		}
		cfg, err := configLoader.Read("")

		assert.NoError(t, err)
		assert.NotNil(t, cfg)
	})
}
