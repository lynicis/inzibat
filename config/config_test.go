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

	t.Run("zero value config", func(t *testing.T) {
		mockReader := NewMockReaderStrategy(ctrl)
		mockReader.EXPECT().
			Read(gomock.Any()).
			Return(&Cfg{}, nil).
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

	t.Run("should merge global and route circuit breaker config", func(t *testing.T) {
		expectedCfg := &Cfg{
			ServerPort: 8080,
			CircuitBreaker: &CircuitBreakerConfig{
				Enabled:          BoolPointer(true),
				FailureThreshold: 9,
			},
			Routes: []Route{
				{
					Method: fiber.MethodGet,
					Path:   "/proxy",
					RequestTo: &RequestTo{
						Path: "/target",
						CircuitBreaker: &CircuitBreakerConfig{
							FailureThreshold: 3,
						},
					},
				},
			},
		}

		mockReader := NewMockReaderStrategy(ctrl)
		mockReader.EXPECT().Read(gomock.Any()).Return(expectedCfg, nil).Times(1)

		cfgLoader := &Reader{ConfigReader: mockReader}
		cfg, err := cfgLoader.Read()

		assert.NoError(t, err)
		assert.NotNil(t, cfg)
		assert.NotNil(t, cfg.CircuitBreaker)
		assert.NotNil(t, cfg.CircuitBreaker.Enabled)
		assert.True(t, *cfg.CircuitBreaker.Enabled)
		assert.Equal(t, 9, cfg.CircuitBreaker.FailureThreshold)
		assert.NotNil(t, cfg.Routes[0].RequestTo.CircuitBreaker)
		assert.NotNil(t, cfg.Routes[0].RequestTo.CircuitBreaker.Enabled)
		assert.True(t, *cfg.Routes[0].RequestTo.CircuitBreaker.Enabled)
		assert.Equal(t, 3, cfg.Routes[0].RequestTo.CircuitBreaker.FailureThreshold)
		assert.Equal(t, 10, cfg.Routes[0].RequestTo.CircuitBreaker.MinimumRequests)
	})
}

func TestReadOrCreateConfig(t *testing.T) {
	t.Run("happy path - create new config when file does not exist", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "inzibat.json")

		cfg, err := ReadOrCreateConfig(configPath)

		assert.NoError(t, err)
		assert.NotNil(t, cfg)
		assert.Equal(t, DefaultSchemaFileName, cfg.Schema)
		assert.Equal(t, 8080, cfg.ServerPort)
		assert.False(t, cfg.HealthCheckRoute)
		assert.Empty(t, cfg.Routes)
	})

	t.Run("happy path - read existing JSON config without schema defaults to schema file", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "inzibat.json")
		configJSON := `{
			"serverPort": 9090,
			"concurrency": 10,
			"healthCheckRoute": true,
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

		assert.NoError(t, err)
		assert.NotNil(t, cfg)
		assert.Equal(t, DefaultSchemaFileName, cfg.Schema)
		assert.Equal(t, 9090, cfg.ServerPort)
		assert.True(t, cfg.HealthCheckRoute)
		assert.Len(t, cfg.Routes, 1)
		assert.Equal(t, "/test", cfg.Routes[0].Path)
	})

	t.Run("happy path - read existing JSON config with custom schema", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "inzibat.json")
		configJSON := `{
			"$schema": "custom.schema.json",
			"serverPort": 9090,
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

		assert.NoError(t, err)
		assert.NotNil(t, cfg)
		assert.Equal(t, "custom.schema.json", cfg.Schema)
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

	t.Run("error path - config path is a directory", func(t *testing.T) {
		cfg, err := ReadOrCreateConfig(t.TempDir())

		assert.Nil(t, cfg)
		assert.Error(t, err)
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

func TestResolveConfigPath(t *testing.T) {
	t.Run("explicit path keeps its extension", func(t *testing.T) {
		path, err := ResolveConfigPath(false, "/tmp/custom.yaml")

		assert.NoError(t, err)
		assert.Equal(t, "/tmp/custom.yaml", path)
	})

	t.Run("explicit path without extension defaults to json", func(t *testing.T) {
		path, extension, err := resolveConfigPath(true, "/tmp/custom")

		assert.NoError(t, err)
		assert.Equal(t, "/tmp/custom", path)
		assert.Equal(t, DefaultConfigExtension, extension)
	})

	t.Run("global config ignores the working directory", func(t *testing.T) {
		homeDir := t.TempDir()
		t.Setenv("HOME", homeDir)

		path, extension, err := resolveConfigPath(true, "")

		assert.NoError(t, err)
		assert.Equal(t, filepath.Join(homeDir, GlobalConfigFileName), path)
		assert.Equal(t, DefaultConfigExtension, extension)
	})

	t.Run("default filename in the working directory", func(t *testing.T) {
		workingDirectory, err := os.Getwd()
		require.NoError(t, err)

		t.Setenv(EnvironmentVariableConfigFileName, "")

		path, extension, err := resolveConfigPath(false, "")

		assert.NoError(t, err)
		assert.Equal(t, filepath.Join(workingDirectory, DefaultConfigFileName), path)
		assert.Equal(t, DefaultConfigExtension, extension)
	})

	t.Run("environment variable overrides the filename", func(t *testing.T) {
		workingDirectory, err := os.Getwd()
		require.NoError(t, err)

		t.Setenv("INZIBAT_CONFIG_FILE", "custom.yaml")

		path, extension, err := resolveConfigPath(false, "")

		assert.NoError(t, err)
		assert.Equal(t, filepath.Join(workingDirectory, "custom.yaml"), path)
		assert.Equal(t, ".yaml", extension)
	})

	t.Run("environment variable without extension defaults to json", func(t *testing.T) {
		workingDirectory, err := os.Getwd()
		require.NoError(t, err)

		t.Setenv(EnvironmentVariableConfigFileName, "custom")

		path, extension, err := resolveConfigPath(false, "")

		assert.NoError(t, err)
		assert.Equal(t, filepath.Join(workingDirectory, "custom.json"), path)
		assert.Equal(t, DefaultConfigExtension, extension)
	})
}

func TestHealthCheckRouteRoundTrip(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "inzibat.json")
	cfg := &Cfg{
		ServerPort:       8080,
		HealthCheckRoute: true,
		Routes: []Route{
			{
				Method:       "GET",
				Path:         "/test",
				FakeResponse: &FakeResponse{StatusCode: 200},
			},
		},
	}

	require.NoError(t, WriteConfig(cfg, configPath))

	loader, err := NewLoader(nil, false, configPath)
	require.NoError(t, err)

	readCfg, err := loader.Read()
	require.NoError(t, err)

	assert.True(t, readCfg.HealthCheckRoute)
	assert.Len(t, readCfg.Routes, 2)
	assert.Equal(t, "/health", readCfg.Routes[1].Path)
}

func TestWriteConfig(t *testing.T) {
	t.Run("happy path - write config to file with empty schema defaults to standard schema", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "inzibat.json")
		cfg := &Cfg{
			ServerPort:       8080,
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
		assert.Equal(t, DefaultSchemaFileName, readCfg.Schema)
		assert.Equal(t, cfg.ServerPort, readCfg.ServerPort)
		assert.Len(t, readCfg.Routes, 1)
	})

	t.Run("happy path - write config to file preserves custom schema", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "inzibat.json")
		cfg := &Cfg{
			Schema:           "https://example.com/custom.schema.json",
			ServerPort:       8080,
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
		assert.Equal(t, "https://example.com/custom.schema.json", readCfg.Schema)
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

		loader, err := NewLoader(validator, false, "")
		require.NoError(t, err)

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

		loader, err := NewLoader(validator, false, "")
		require.NoError(t, err)

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

		loader, err := NewLoader(validator, false, "")
		require.NoError(t, err)

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

		loader, err := NewLoader(validator, true, "")
		require.NoError(t, err)

		assert.NotNil(t, loader)
		assert.NotNil(t, loader.ConfigReader)
		assert.Contains(t, loader.Filepath, GlobalConfigFileName)
	})

	t.Run("happy path - explicit path overrides env", func(t *testing.T) {
		tmpDir := t.TempDir()
		explicitPath := filepath.Join(tmpDir, "custom.yaml")

		loader, err := NewLoader(validator, false, explicitPath)
		require.NoError(t, err)

		assert.Equal(t, explicitPath, loader.Filepath)
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

		loader, err := NewLoader(validator, false, "")

		assert.Nil(t, loader)
		assert.ErrorContains(t, err, "ambiguous file extension")
	})
}
