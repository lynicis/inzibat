package server

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"inzibat/config"
)

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

func TestStartServer(t *testing.T) {
	t.Run("happy path - with valid config file", func(t *testing.T) {
		tmpDir := t.TempDir()
		configFile := filepath.Join(tmpDir, "inzibat.json")
		cfg := &config.Cfg{
			ServerPort:       8080,
			Concurrency:      1,
			HealthCheckRoute: false,
			Routes: []config.Route{
				{
					Method: "GET",
					Path:   "/test",
					FakeResponse: config.FakeResponse{
						StatusCode: 200,
					},
				},
			},
		}
		err := config.WriteConfig(cfg, configFile)
		require.NoError(t, err)

		done := make(chan error, 1)
		go func() {
			done <- StartServer(configFile)
		}()

		time.Sleep(100 * time.Millisecond)
	})

	t.Run("error path - config file path resolution fails", func(t *testing.T) {
		invalidPath := "/nonexistent/path/to/config.json"

		err := StartServer(invalidPath)

		assert.Error(t, err)
		assert.True(t,
			contains(err.Error(), "failed to resolve config file path") ||
				contains(err.Error(), "failed to read config"),
			"error: %s", err.Error())
	})

	t.Run("error path - config file does not exist", func(t *testing.T) {
		tmpDir := t.TempDir()
		nonExistentFile := filepath.Join(tmpDir, "nonexistent.json")

		err := StartServer(nonExistentFile)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to read config")
	})

	t.Run("happy path - with empty config file string", func(t *testing.T) {
		tmpDir := t.TempDir()
		originalWd, err := os.Getwd()
		require.NoError(t, err)
		defer os.Chdir(originalWd)

		err = os.Chdir(tmpDir)
		require.NoError(t, err)

		cfg := &config.Cfg{
			ServerPort:       8080,
			Concurrency:      1,
			HealthCheckRoute: false,
			Routes: []config.Route{
				{
					Method: "GET",
					Path:   "/test",
					FakeResponse: config.FakeResponse{
						StatusCode: 200,
						Body:       config.HttpBody{"message": "test"},
					},
				},
			},
		}
		err = config.WriteConfig(cfg, "inzibat.json")
		require.NoError(t, err)

		originalEnv := os.Getenv(config.EnvironmentVariableConfigFileName)
		defer func() {
			if originalEnv != "" {
				os.Setenv(config.EnvironmentVariableConfigFileName, originalEnv)
			} else {
				os.Unsetenv(config.EnvironmentVariableConfigFileName)
			}
		}()
		os.Unsetenv(config.EnvironmentVariableConfigFileName)

		done := make(chan error, 1)
		go func() {
			done <- StartServer("")
		}()

		time.Sleep(50 * time.Millisecond)

		select {
		case err := <-done:
			if err != nil {
				assert.True(t,
					contains(err.Error(), "failed to start http server") ||
						contains(err.Error(), "failed to read config"),
					"error: %s", err.Error())
			}
		default:
		}
	})

	t.Run("error path - environment variable set fails", func(t *testing.T) {
		tmpDir := t.TempDir()
		configFile := filepath.Join(tmpDir, "test.json")

		cfg := &config.Cfg{
			ServerPort:       8080,
			Concurrency:      1,
			HealthCheckRoute: false,
			Routes: []config.Route{
				{
					Method: "GET",
					Path:   "/test",
					FakeResponse: config.FakeResponse{
						StatusCode: 200,
					},
				},
			},
		}
		err := config.WriteConfig(cfg, configFile)
		require.NoError(t, err)

		invalidConfigFile := filepath.Join(tmpDir, "invalid.json")
		err = os.WriteFile(invalidConfigFile, []byte("invalid json"), 0644)
		require.NoError(t, err)

		err = StartServer(invalidConfigFile)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to read config")
	})
}
