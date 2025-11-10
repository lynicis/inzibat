package config

import (
	"errors"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestReader_Read(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

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
			Concurrency: 5,
		}

		mockReader := NewMockReaderStrategy(ctrl)
		mockReader.EXPECT().
			Read(gomock.Any()).
			Return(expectedCfg, nil).
			Times(1)

		cfgLoader := &Reader{
			ConfigReader: mockReader,
		}
		cfg, err := cfgLoader.Read("test-file-name")

		assert.NoError(t, err)
		assert.Equal(t, expectedCfg, cfg)
	})

	t.Run("when reader return error should return it", func(t *testing.T) {
		mockReader := NewMockReaderStrategy(ctrl)
		mockReader.EXPECT().
			Read(gomock.Any()).
			Return(nil, errors.New("something went wrong")).
			Times(1)

		cfgLoader := &Reader{
			ConfigReader: mockReader,
		}
		cfg, err := cfgLoader.Read("test-file-name")

		assert.Nil(t, cfg)
		assert.Errorf(t, err, "something went wrong")
	})

	t.Run("against healthcheck route", func(t *testing.T) {
		mockReader := NewMockReaderStrategy(ctrl)
		mockReader.EXPECT().
			Read(gomock.Any()).
			Return(&Cfg{HealthCheckRoute: true}, nil).
			Times(1)

		configReader := &Reader{
			ConfigReader: mockReader,
		}

		cfg, err := configReader.Read("")

		assert.NoError(t, err)
		assert.NotNil(t, cfg)
	})

	t.Run("against concurrency route creator limit", func(t *testing.T) {
		mockReader := NewMockReaderStrategy(ctrl)
		mockReader.EXPECT().
			Read(gomock.Any()).
			Return(&Cfg{
				Concurrency: 0,
			}, nil).
			Times(1)

		configReader := &Reader{
			ConfigReader: mockReader,
		}
		cfg, err := configReader.Read("")

		assert.NoError(t, err)
		assert.NotNil(t, cfg)
	})
}
