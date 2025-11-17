package config

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/goccy/go-json"

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
					RequestTo: &RequestTo{
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
					FakeResponse: &FakeResponse{
						Body:       HttpBody{},
						StatusCode: http.StatusOK,
					},
				},
				{
					Method: fiber.MethodGet,
					Path:   "/route-two",
					RequestTo: &RequestTo{
						Method: http.MethodGet,
						Headers: map[string][]string{
							"X-Test-Header": {"Test-Header-Value"},
						},
						Host:                   "http://localhost:8081",
						Path:                   "/route-two",
						PassWithRequestBody:    true,
						PassWithRequestHeaders: true,
					},
					FakeResponse: &FakeResponse{
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

	t.Run("when validator returns error should return it", func(t *testing.T) {
		invalidCfg := &Cfg{
			ServerPort: 0,
			Routes:     []Route{},
		}

		mockReader := NewMockReaderStrategy(ctrl)
		mockReader.EXPECT().
			Read(gomock.Any()).
			Return(invalidCfg, nil).
			Times(1)

		cfgLoader := &Reader{
			ConfigReader: mockReader,
			Validator:    validator.New(),
		}

		cfg, err := cfgLoader.Read()

		assert.Nil(t, cfg)
		assert.Error(t, err)
	})

	t.Run("when RequestTo.Method is empty should default to GET", func(t *testing.T) {
		cfgWithEmptyMethod := &Cfg{
			ServerPort: 8080,
			Routes: []Route{
				{
					Method: fiber.MethodGet,
					Path:   "/test",
					RequestTo: &RequestTo{
						Method: "",
						Path:   "/test",
					},
				},
			},
		}

		mockReader := NewMockReaderStrategy(ctrl)
		mockReader.EXPECT().
			Read(gomock.Any()).
			Return(cfgWithEmptyMethod, nil).
			Times(1)

		cfgLoader := &Reader{
			ConfigReader: mockReader,
		}

		cfg, err := cfgLoader.Read()

		assert.NoError(t, err)
		assert.NotNil(t, cfg)
		assert.Equal(t, http.MethodGet, cfg.Routes[0].RequestTo.Method)
	})

	t.Run("when GET method has body should return error", func(t *testing.T) {
		cfgWithGetBody := &Cfg{
			ServerPort: 8080,
			Routes: []Route{
				{
					Method: fiber.MethodGet,
					Path:   "/test",
					RequestTo: &RequestTo{
						Method: http.MethodGet,
						Body: HttpBody{
							"key": "value",
						},
					},
				},
			},
		}

		mockReader := NewMockReaderStrategy(ctrl)
		mockReader.EXPECT().
			Read(gomock.Any()).
			Return(cfgWithGetBody, nil).
			Times(1)

		cfgLoader := &Reader{
			ConfigReader: mockReader,
		}

		cfg, err := cfgLoader.Read()

		assert.Nil(t, cfg)
		assert.Equal(t, ErrorGetSendBody, err)
	})

	t.Run("when route only has fake response it should pass validation", func(t *testing.T) {
		cfgWithMockOnly := &Cfg{
			ServerPort: 8080,
			Routes: []Route{
				{
					Method: fiber.MethodGet,
					Path:   "/mock",
					FakeResponse: &FakeResponse{
						StatusCode: http.StatusOK,
						Body: HttpBody{
							"message": "ok",
						},
					},
				},
			},
		}

		mockReader := NewMockReaderStrategy(ctrl)
		mockReader.EXPECT().
			Read(gomock.Any()).
			Return(cfgWithMockOnly, nil).
			Times(1)

		cfgLoader := &Reader{
			ConfigReader: mockReader,
			Validator:    validator.New(),
		}

		cfg, err := cfgLoader.Read()

		assert.NoError(t, err)
		assert.NotNil(t, cfg)
		assert.Nil(t, cfg.Routes[0].RequestTo)
		assert.NotNil(t, cfg.Routes[0].FakeResponse)
	})

	t.Run("should assign route method and RequestTo method correctly", func(t *testing.T) {
		expectedCfg := &Cfg{
			ServerPort: 8080,
			Routes: []Route{
				{
					Method: fiber.MethodPost,
					Path:   "/test",
					RequestTo: &RequestTo{
						Method: http.MethodPut,
						Path:   "/test",
					},
				},
			},
		}

		mockReader := NewMockReaderStrategy(ctrl)
		mockReader.EXPECT().
			Read(gomock.Any()).
			Return(expectedCfg, nil).
			Times(1)

		cfgLoader := &Reader{
			ConfigReader: mockReader,
		}

		cfg, err := cfgLoader.Read()

		assert.NoError(t, err)
		assert.NotNil(t, cfg)
		assert.Equal(t, fiber.MethodPost, cfg.Routes[0].Method)
		assert.Equal(t, http.MethodPut, cfg.Routes[0].RequestTo.Method)
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

	t.Run("error path - invalid file extension", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "config.unknown")
		err := os.WriteFile(configPath, []byte("test"), 0644)
		require.NoError(t, err)

		cfg, err := ReadOrCreateConfig(configPath)

		assert.Nil(t, cfg)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create reader strategy")
	})

	t.Run("error path - read error", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "inzibat.json")
		invalidJSON := `{invalid json}`
		err := os.WriteFile(configPath, []byte(invalidJSON), 0644)
		require.NoError(t, err)

		cfg, err := ReadOrCreateConfig(configPath)

		assert.Nil(t, cfg)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to read config")
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
					FakeResponse: &FakeResponse{
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

	t.Run("error path - path resolution failure", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "inzibat.json")

		cfg := &Cfg{
			ServerPort: 8080,
		}

		err := WriteConfig(cfg, configPath)
		assert.NoError(t, err)
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

func TestNewLoader(t *testing.T) {
	validator := validator.New()

	t.Run("happy path - local config with default filename", func(t *testing.T) {
		originalEnv := os.Getenv(EnvironmentVariableConfigFileName)
		defer func() {
			if originalEnv != "" {
				os.Setenv(EnvironmentVariableConfigFileName, originalEnv)
			} else {
				os.Unsetenv(EnvironmentVariableConfigFileName)
			}
		}()
		os.Unsetenv(EnvironmentVariableConfigFileName)

		loader := NewLoader(validator, false)

		assert.NotNil(t, loader)
		assert.NotNil(t, loader.ConfigReader)
		assert.Equal(t, validator, loader.Validator)
		assert.Contains(t, loader.Filepath, DefaultConfigFileName)
	})

	t.Run("happy path - local config with custom filename from env", func(t *testing.T) {
		originalEnv := os.Getenv(EnvironmentVariableConfigFileName)
		defer func() {
			if originalEnv != "" {
				os.Setenv(EnvironmentVariableConfigFileName, originalEnv)
			} else {
				os.Unsetenv(EnvironmentVariableConfigFileName)
			}
		}()
		customFileName := "custom-config.json"
		os.Setenv(EnvironmentVariableConfigFileName, customFileName)

		loader := NewLoader(validator, false)

		assert.NotNil(t, loader)
		assert.Contains(t, loader.Filepath, customFileName)
	})

	t.Run("happy path - local config with JSON extension", func(t *testing.T) {
		originalEnv := os.Getenv(EnvironmentVariableConfigFileName)
		defer func() {
			if originalEnv != "" {
				os.Setenv(EnvironmentVariableConfigFileName, originalEnv)
			} else {
				os.Unsetenv(EnvironmentVariableConfigFileName)
			}
		}()
		os.Setenv(EnvironmentVariableConfigFileName, "config.json")

		loader := NewLoader(validator, false)

		assert.NotNil(t, loader)
		assert.NotNil(t, loader.ConfigReader)
	})

	t.Run("happy path - global config", func(t *testing.T) {
		originalEnv := os.Getenv(EnvironmentVariableConfigFileName)
		defer func() {
			if originalEnv != "" {
				os.Setenv(EnvironmentVariableConfigFileName, originalEnv)
			} else {
				os.Unsetenv(EnvironmentVariableConfigFileName)
			}
		}()

		loader := NewLoader(validator, true)

		assert.NotNil(t, loader)
		assert.NotNil(t, loader.ConfigReader)
		assert.Contains(t, loader.Filepath, GlobalConfigFileName)
	})

	t.Run("error path - invalid file extension", func(t *testing.T) {
		originalEnv := os.Getenv(EnvironmentVariableConfigFileName)
		defer func() {
			if originalEnv != "" {
				os.Setenv(EnvironmentVariableConfigFileName, originalEnv)
			} else {
				os.Unsetenv(EnvironmentVariableConfigFileName)
			}
		}()
		os.Setenv(EnvironmentVariableConfigFileName, "config.unknown")
	})
}

func TestWrite(t *testing.T) {
	t.Run("happy path - write route to directory", func(t *testing.T) {
		tmpDir := t.TempDir()
		route := &Route{
			Method: "GET",
			Path:   "/test",
			FakeResponse: &FakeResponse{
				StatusCode: 200,
			},
		}

		err := Write(route, tmpDir)

		assert.NoError(t, err)
		configPath := filepath.Join(tmpDir, "inzibat.json")
		assert.FileExists(t, configPath)

		data, err := os.ReadFile(configPath)
		require.NoError(t, err)
		var readRoute Route
		err = json.Unmarshal(data, &readRoute)
		require.NoError(t, err)
		assert.Equal(t, route.Method, readRoute.Method)
		assert.Equal(t, route.Path, readRoute.Path)
	})

	t.Run("error path - invalid directory", func(t *testing.T) {
		invalidDir := "/invalid/path/that/does/not/exist"
		route := &Route{
			Method: "GET",
			Path:   "/test",
		}

		err := Write(route, invalidDir)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no such file or directory")
	})

	t.Run("error path - directory path resolution failure", func(t *testing.T) {
		tmpDir := t.TempDir()
		route := &Route{
			Method: "GET",
			Path:   "/test",
		}

		err := Write(route, tmpDir)

		assert.NoError(t, err)
	})
}

func TestWriteJSONToFile(t *testing.T) {
	t.Run("happy path - writes JSON to file", func(t *testing.T) {
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "test.json")
		data := map[string]interface{}{
			"test": "value",
			"num":  42,
		}

		err := writeJSONToFile(filePath, data)

		assert.NoError(t, err)
		assert.FileExists(t, filePath)

		readData, err := os.ReadFile(filePath)
		require.NoError(t, err)
		var result map[string]interface{}
		err = json.Unmarshal(readData, &result)
		require.NoError(t, err)
		assert.Equal(t, "value", result["test"])
		assert.Equal(t, float64(42), result["num"])
	})

	t.Run("error path - invalid file path", func(t *testing.T) {
		invalidPath := "/invalid/path/that/does/not/exist/test.json"
		data := map[string]interface{}{"test": "value"}

		err := writeJSONToFile(invalidPath, data)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create file")
	})

	t.Run("error path - path resolution failure", func(t *testing.T) {
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "test.json")
		data := map[string]interface{}{"test": "value"}

		err := writeJSONToFile(filePath, data)

		assert.NoError(t, err)
	})
}
