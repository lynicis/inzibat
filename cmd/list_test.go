package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"inzibat/config"
)

func TestListCmd(t *testing.T) {
	t.Run("happy path - command is registered", func(t *testing.T) {
		assert.NotNil(t, listCmd)
		assert.Equal(t, "list", listCmd.Use)
		assert.Contains(t, listCmd.Aliases, "list-routes")
		assert.Contains(t, listCmd.Aliases, "ls")
		assert.Contains(t, listCmd.Aliases, "l")
	})

	t.Run("happy path - command has correct short description", func(t *testing.T) {
		assert.Contains(t, listCmd.Short, "List")
	})

	t.Run("happy path - command has Run function", func(t *testing.T) {
		assert.NotNil(t, listCmd.Run)
	})
}

func TestListCmd_Run(t *testing.T) {
	t.Run("happy path - successfully reads and logs routes from global config", func(t *testing.T) {
		tmpHomeDir := t.TempDir()
		originalHomeDir := os.Getenv("HOME")
		defer func() {
			if originalHomeDir != "" {
				os.Setenv("HOME", originalHomeDir)
			} else {
				os.Unsetenv("HOME")
			}
		}()
		os.Setenv("HOME", tmpHomeDir)

		globalConfigPath := filepath.Join(tmpHomeDir, config.GlobalConfigFileName)
		cfg := &config.Cfg{
			ServerPort:       8080,
			Concurrency:      5,
			HealthCheckRoute: false,
			Routes: []config.Route{
				{
					Method: "GET",
					Path:   "/test-route",
					FakeResponse: config.FakeResponse{
						StatusCode: 200,
						Body:       config.HttpBody{"message": "test"},
					},
					RequestTo: config.RequestTo{
						Method: "GET",
						Host:   "http://localhost:8081",
						Path:   "/test-route",
					},
				},
			},
		}

		configData, err := json.Marshal(cfg)
		require.NoError(t, err)
		err = os.WriteFile(globalConfigPath, configData, 0644)
		require.NoError(t, err)

		assert.NotPanics(t, func() {
			listCmd.Run(listCmd, []string{})
		})

		_, err = os.Stat(globalConfigPath)
		assert.NoError(t, err)
	})

	t.Run("happy path - handles single route", func(t *testing.T) {
		tmpHomeDir := t.TempDir()
		originalHomeDir := os.Getenv("HOME")
		defer func() {
			if originalHomeDir != "" {
				os.Setenv("HOME", originalHomeDir)
			} else {
				os.Unsetenv("HOME")
			}
		}()
		os.Setenv("HOME", tmpHomeDir)

		globalConfigPath := filepath.Join(tmpHomeDir, config.GlobalConfigFileName)
		cfg := &config.Cfg{
			ServerPort:       8080,
			Concurrency:      5,
			HealthCheckRoute: false,
			Routes: []config.Route{
				{
					Method: "GET",
					Path:   "/single-route-test",
					FakeResponse: config.FakeResponse{
						StatusCode: 200,
						Body:       config.HttpBody{"message": "single route"},
					},
					RequestTo: config.RequestTo{
						Method: "GET",
						Host:   "http://localhost:8081",
						Path:   "/single-route-test",
					},
				},
			},
		}

		configData, err := json.Marshal(cfg)
		require.NoError(t, err)
		err = os.WriteFile(globalConfigPath, configData, 0644)
		require.NoError(t, err)

		assert.NotPanics(t, func() {
			listCmd.Run(listCmd, []string{})
		})

		_, err = os.Stat(globalConfigPath)
		assert.NoError(t, err)
	})

	t.Run("happy path - handles multiple routes", func(t *testing.T) {
		tmpHomeDir := t.TempDir()
		originalHomeDir := os.Getenv("HOME")
		defer func() {
			if originalHomeDir != "" {
				os.Setenv("HOME", originalHomeDir)
			} else {
				os.Unsetenv("HOME")
			}
		}()
		os.Setenv("HOME", tmpHomeDir)

		globalConfigPath := filepath.Join(tmpHomeDir, config.GlobalConfigFileName)
		cfg := &config.Cfg{
			ServerPort:       8080,
			Concurrency:      5,
			HealthCheckRoute: false,
			Routes: []config.Route{
				{
					Method: "GET",
					Path:   "/route-one",
					FakeResponse: config.FakeResponse{
						StatusCode: 200,
						Body:       config.HttpBody{"message": "route one"},
					},
					RequestTo: config.RequestTo{
						Method: "GET",
						Host:   "http://localhost:8081",
						Path:   "/route-one",
					},
				},
				{
					Method: "PUT",
					Path:   "/route-two",
					FakeResponse: config.FakeResponse{
						StatusCode: 201,
						Body:       config.HttpBody{"message": "route two"},
					},
					RequestTo: config.RequestTo{
						Method: "PUT",
						Host:   "http://localhost:8081",
						Path:   "/route-two",
					},
				},
				{
					Method: "PUT",
					Path:   "/route-three",
					FakeResponse: config.FakeResponse{
						StatusCode: 200,
						Body:       config.HttpBody{"message": "route three"},
					},
					RequestTo: config.RequestTo{
						Method: "PUT",
						Host:   "http://localhost:8081",
						Path:   "/route-three",
					},
				},
			},
		}

		configData, err := json.Marshal(cfg)
		require.NoError(t, err)
		err = os.WriteFile(globalConfigPath, configData, 0644)
		require.NoError(t, err)

		assert.NotPanics(t, func() {
			listCmd.Run(listCmd, []string{})
		})

		_, err = os.Stat(globalConfigPath)
		assert.NoError(t, err)
	})
}
