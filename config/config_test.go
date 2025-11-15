package config

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
						Method: http.MethodPut,
						Headers: map[string][]string{
							"X-Test-Header": {"Test-Header-Value"},
						},
						Body: HttpBody{
							"testKey": "testValue",
						},
						Host:                   "http://localhost:8081",
						Path:                   "/route-one",
						PassWithRequestBody:    true,
						PassWithRequestHeaders: true,
					},
					FakeResponse: FakeResponse{
						Body:       HttpBody{},
						StatusCode: http.StatusOK,
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
					FakeResponse: FakeResponse{
						Body:       HttpBody{},
						StatusCode: http.StatusOK,
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
			Validator:    validator.New(),
		}
		cfg, err := cfgLoader.Read()

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
		cfg, err := cfgLoader.Read()

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

		cfg, err := configReader.Read()

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
		cfg, err := configReader.Read()

		assert.NoError(t, err)
		assert.NotNil(t, cfg)
	})
}

func TestReadOrCreateConfig(t *testing.T) {
	t.Run("happy path - create new config when file does not exist", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "inzibat.json")

		cfg, err := ReadOrCreateConfig(configPath)

		assert.NoError(t, err)
		assert.NotNil(t, cfg)
		assert.Equal(t, 8080, cfg.ServerPort)
		assert.Equal(t, 5, cfg.Concurrency)
		assert.False(t, cfg.HealthCheckRoute)
		assert.Empty(t, cfg.Routes)
	})

	t.Run("happy path - read existing JSON config", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "inzibat.json")
		// Use koanf field names for proper unmarshaling
		configJSON := `{
			"serverPort": 9090,
			"concurrency": 10,
			"isHealthCheckRouteEnabled": true,
			"routes": [
				{
					"method": "GET",
					"path": "/test"
				}
			]
		}`
		err := os.WriteFile(configPath, []byte(configJSON), 0644)
		require.NoError(t, err)

		cfg, err := ReadOrCreateConfig(configPath)

		assert.NotNil(t, cfg)
		assert.Equal(t, 9090, cfg.ServerPort)
		assert.Equal(t, 10, cfg.Concurrency)
		assert.True(t, cfg.HealthCheckRoute)
		assert.Len(t, cfg.Routes, 1)
		assert.Equal(t, "/test", cfg.Routes[0].Path)
	})

	t.Run("happy path - file without extension defaults to JSON", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "config")

		cfg, err := ReadOrCreateConfig(configPath)

		assert.NoError(t, err)
		assert.NotNil(t, cfg)
		assert.Equal(t, 8080, cfg.ServerPort)
	})
}

func TestWriteConfig(t *testing.T) {
	t.Run("happy path - write config to file", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "inzibat.json")
		cfg := &Cfg{
			ServerPort:       8080,
			Concurrency:      5,
			HealthCheckRoute: false,
			Routes: []Route{
				{
					Method: "GET",
					Path:   "/test",
					FakeResponse: FakeResponse{
						StatusCode: 200,
					},
				},
			},
		}

		err := WriteConfig(cfg, configPath)

		assert.NoError(t, err)
		assert.FileExists(t, configPath)

		data, err := os.ReadFile(configPath)
		require.NoError(t, err)
		var readCfg Cfg
		err = json.Unmarshal(data, &readCfg)
		require.NoError(t, err)
		assert.Equal(t, cfg.ServerPort, readCfg.ServerPort)
		assert.Equal(t, cfg.Concurrency, readCfg.Concurrency)
		assert.Len(t, readCfg.Routes, 1)
	})

	t.Run("error path - invalid directory", func(t *testing.T) {
		invalidPath := "/invalid/path/that/does/not/exist/inzibat.json"
		cfg := &Cfg{
			ServerPort: 8080,
		}

		err := WriteConfig(cfg, invalidPath)

		assert.Contains(t, err.Error(), "failed to create file")
	})
}

func TestInitGlobalConfig(t *testing.T) {
	t.Run("happy path - file does not exist", func(t *testing.T) {
		homeDir, err := os.UserHomeDir()
		require.NoError(t, err, "should be able to get home directory in test environment")

		globalConfigPath := filepath.Join(homeDir, DefaultConfigFileName)

		var originalFileExists bool
		var originalFileData []byte
		if _, err := os.Stat(globalConfigPath); err == nil {
			originalFileExists = true
			originalFileData, err = os.ReadFile(globalConfigPath)
			require.NoError(t, err)
			err = os.Remove(globalConfigPath)
			require.NoError(t, err)
		}
		defer func() {
			if originalFileExists {
				err := os.WriteFile(globalConfigPath, originalFileData, 0644)
				require.NoError(t, err)
			} else {
				os.Remove(globalConfigPath)
			}
		}()

		err = InitGlobalConfig()

		assert.NoError(t, err)
		_, statErr := os.Stat(globalConfigPath)
		assert.True(t, os.IsNotExist(statErr), "config file should not exist")
	})

	t.Run("happy path - file exists and gets rewritten", func(t *testing.T) {
		homeDir, err := os.UserHomeDir()
		require.NoError(t, err, "should be able to get home directory in test environment")

		globalConfigPath := filepath.Join(homeDir, DefaultConfigFileName)

		var originalFileExists bool
		var originalFileData []byte
		if _, err := os.Stat(globalConfigPath); err == nil {
			originalFileExists = true
			originalFileData, err = os.ReadFile(globalConfigPath)
			require.NoError(t, err)
		}
		defer func() {
			if originalFileExists {
				err := os.WriteFile(globalConfigPath, originalFileData, 0644)
				require.NoError(t, err)
			} else {
				os.Remove(globalConfigPath)
			}
		}()

		configJSON := `{
			"serverPort": 9090,
			"concurrency": 10,
			"isHealthCheckRouteEnabled": true,
			"routes": [
				{
					"method": "GET",
					"path": "/existing",
					"fakeResponse": {
						"statusCode": 200
					}
				}
			]
		}`
		err = os.WriteFile(globalConfigPath, []byte(configJSON), 0644)
		require.NoError(t, err)
		require.FileExists(t, globalConfigPath)

		err = InitGlobalConfig()

		assert.NoError(t, err)
		assert.FileExists(t, globalConfigPath)

		data, err := os.ReadFile(globalConfigPath)
		require.NoError(t, err)
		var readCfg Cfg
		err = json.Unmarshal(data, &readCfg)
		require.NoError(t, err)
		assert.Equal(t, 9090, readCfg.ServerPort)
		assert.Equal(t, 10, readCfg.Concurrency)
		assert.True(t, readCfg.HealthCheckRoute)
		assert.Len(t, readCfg.Routes, 1)
		assert.Equal(t, "/existing", readCfg.Routes[0].Path)
	})
}
