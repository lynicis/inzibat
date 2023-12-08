package config

import (
	"errors"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestReadConfig(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		cfg, err := ReadConfig("../test-data/test.json")

		assert.NoError(t, err)
		assert.Equal(t, &Config{
			ServerPort: "8080",
			Routes: []Route{
				{
					Method: fiber.MethodGet,
					Path:   "/route-one",
					RequestTo: RequestTo{
						Method: http.MethodPost,
						Headers: map[string]string{
							"xtestheader": "TestHeaderValue",
						},
						Body: map[string]interface{}{
							"testkey": "testValue",
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
						Headers: map[string]string{
							"xtestheader": "TestHeaderValue",
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
		}, cfg)
	})

	t.Run("config file doesn't exist", func(t *testing.T) {
		t.Run("with file extension", func(t *testing.T) {
			cfg, err := ReadConfig("not-found.json")

			assert.Empty(t, cfg)
			if assert.Error(t, err) {
				expectedError := errors.New(ErrorFileNotFound)
				assert.Equal(t, expectedError, err)
			}
		})

		t.Run("without file extension", func(t *testing.T) {
			cfg, err := ReadConfig("not-found")

			assert.Empty(t, cfg)
			if assert.Error(t, err) {
				expectedError := errors.New(ErrorFileNotFound)
				assert.Equal(t, expectedError, err)
			}
		})
	})
}
