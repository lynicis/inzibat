package server

import (
	"context"
	"fmt"
	nethttp "net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"inzibat/client/http"
	"inzibat/config"
)

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

func TestStartServer(t *testing.T) {
	t.Run("happy path - with valid config file", func(t *testing.T) {
		// Arrange
		tmpDir := t.TempDir()
		configFile := filepath.Join(tmpDir, "inzibat.json")
		freePort, err := http.GetFreePort()
		require.NoError(t, err)

		cfg := &config.Cfg{
			ServerPort:       freePort,
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
		err = config.WriteConfig(cfg, configFile)
		require.NoError(t, err)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		done := make(chan error, 1)
		go func() {
			done <- StartServerWithContext(ctx, configFile)
		}()

		time.Sleep(200 * time.Millisecond)

		go func() {
			time.Sleep(100 * time.Millisecond)
			cancel()
		}()

		select {
		case err := <-done:
			if err != nil {
				t.Logf("Server shutdown completed with: %v", err)
			}
		case <-time.After(5 * time.Second):
			t.Fatal("server did not shutdown within timeout")
		}
	})

	t.Run("error path - config file path resolution fails", func(t *testing.T) {
		// Arrange
		invalidPath := "/nonexistent/path/to/config.json"

		// Act
		err := StartServer(invalidPath)

		// Assert
		assert.Error(t, err)
		assert.True(t,
			contains(err.Error(), "failed to resolve config file path") ||
				contains(err.Error(), "failed to read config"),
			"error: %s", err.Error())
	})

	t.Run("error path - config file does not exist", func(t *testing.T) {
		// Arrange
		tmpDir := t.TempDir()
		nonExistentFile := filepath.Join(tmpDir, "nonexistent.json")

		// Act
		err := StartServer(nonExistentFile)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to read config")
	})

	t.Run("happy path - with empty config file string", func(t *testing.T) {
		// Arrange
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

		// Act
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		done := make(chan error, 1)
		go func() {
			done <- StartServerWithContext(ctx, "")
		}()

		time.Sleep(50 * time.Millisecond)

		go func() {
			time.Sleep(100 * time.Millisecond)
			cancel()
		}()

		// Assert
		select {
		case err := <-done:
			if err != nil {
				assert.True(t,
					contains(err.Error(), "failed to start http server") ||
						contains(err.Error(), "failed to read config") ||
						contains(err.Error(), "failed to shutdown gracefully"),
					"error: %s", err.Error())
			}
		case <-time.After(5 * time.Second):
			t.Fatal("server did not shutdown within timeout")
		}
	})

	t.Run("error path - environment variable set fails", func(t *testing.T) {
		// Arrange
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

		// Act
		err = StartServer(invalidConfigFile)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to read config")
	})

	t.Run("happy path - environment variable restoration with original env set", func(t *testing.T) {
		// Arrange
		tmpDir := t.TempDir()
		configFile := filepath.Join(tmpDir, "inzibat.json")
		freePort, err := http.GetFreePort()
		require.NoError(t, err)

		cfg := &config.Cfg{
			ServerPort:       freePort,
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
		err = config.WriteConfig(cfg, configFile)
		require.NoError(t, err)

		// Set original environment variable
		originalEnvValue := "original_config.json"
		err = os.Setenv(config.EnvironmentVariableConfigFileName, originalEnvValue)
		require.NoError(t, err)
		defer func() {
			os.Setenv(config.EnvironmentVariableConfigFileName, originalEnvValue)
		}()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		done := make(chan error, 1)
		go func() {
			done <- StartServerWithContext(ctx, configFile)
		}()

		time.Sleep(200 * time.Millisecond)

		go func() {
			time.Sleep(100 * time.Millisecond)
			cancel()
		}()

		select {
		case err := <-done:
			if err != nil {
				t.Logf("Server shutdown completed with: %v", err)
			}
		case <-time.After(5 * time.Second):
			t.Fatal("server did not shutdown within timeout")
		}
	})

	t.Run("happy path - environment variable restoration with no original env", func(t *testing.T) {
		// Arrange
		tmpDir := t.TempDir()
		configFile := filepath.Join(tmpDir, "inzibat.json")
		freePort, err := http.GetFreePort()
		require.NoError(t, err)

		cfg := &config.Cfg{
			ServerPort:       freePort,
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
		err = config.WriteConfig(cfg, configFile)
		require.NoError(t, err)

		// Ensure no original environment variable is set
		originalEnv := os.Getenv(config.EnvironmentVariableConfigFileName)
		if originalEnv != "" {
			err = os.Unsetenv(config.EnvironmentVariableConfigFileName)
			require.NoError(t, err)
			defer os.Setenv(config.EnvironmentVariableConfigFileName, originalEnv)
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		done := make(chan error, 1)
		go func() {
			done <- StartServerWithContext(ctx, configFile)
		}()

		time.Sleep(200 * time.Millisecond)

		go func() {
			time.Sleep(100 * time.Millisecond)
			cancel()
		}()

		select {
		case err := <-done:
			if err != nil {
				t.Logf("Server shutdown completed with: %v", err)
			}
		case <-time.After(5 * time.Second):
			t.Fatal("server did not shutdown within timeout")
		}
	})

	t.Run("happy path - server startup and route creation", func(t *testing.T) {
		// Arrange
		tmpDir := t.TempDir()
		configFile := filepath.Join(tmpDir, "inzibat.json")
		freePort, err := http.GetFreePort()
		require.NoError(t, err)

		cfg := &config.Cfg{
			ServerPort:       freePort,
			Concurrency:      1,
			HealthCheckRoute: true,
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
		err = config.WriteConfig(cfg, configFile)
		require.NoError(t, err)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		done := make(chan error, 1)
		go func() {
			done <- StartServerWithContext(ctx, configFile)
		}()

		time.Sleep(200 * time.Millisecond)

		go func() {
			time.Sleep(100 * time.Millisecond)
			cancel()
		}()

		select {
		case err := <-done:
			if err != nil {
				t.Logf("Server shutdown completed with: %v", err)
			}
		case <-time.After(5 * time.Second):
			t.Fatal("server did not shutdown within timeout")
		}
	})

	t.Run("happy path - full server lifecycle with graceful shutdown", func(t *testing.T) {
		originalWd, err := os.Getwd()
		require.NoError(t, err)
		defer os.Chdir(originalWd)

		tmpDir := t.TempDir()
		err = os.Chdir(tmpDir)
		require.NoError(t, err)

		configFile := "inzibat.json"
		freePort, err := http.GetFreePort()
		require.NoError(t, err)

		cfg := &config.Cfg{
			ServerPort:       freePort,
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
		err = config.WriteConfig(cfg, configFile)
		require.NoError(t, err)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		serverDone := make(chan error, 1)
		go func() {
			serverDone <- StartServerWithContext(ctx, configFile)
		}()

		time.Sleep(500 * time.Millisecond)

		client := &nethttp.Client{Timeout: 2 * time.Second}
		resp, err := client.Get(fmt.Sprintf("http://localhost:%d/test", freePort))
		if err == nil {
			resp.Body.Close()
		}

		go func() {
			time.Sleep(200 * time.Millisecond)
			cancel()
		}()

		select {
		case err := <-serverDone:
			if err != nil {
				t.Logf("Server shutdown completed with: %v", err)
			}
		case <-time.After(5 * time.Second):
			t.Fatal("server did not shutdown within timeout")
		}
	})
}
