package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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

	t.Run("error path - fails when config file is invalid (subprocess)", func(t *testing.T) {
		_, testFile, _, ok := runtime.Caller(0)
		require.True(t, ok, "failed to get test file path")

		testDir := filepath.Dir(testFile)
		projectRoot := filepath.Dir(testDir)
		projectRoot, err := filepath.Abs(projectRoot)
		require.NoError(t, err)

		binaryPath := filepath.Join(projectRoot, "inzibat")

		if _, statErr := os.Stat(binaryPath); os.IsNotExist(statErr) {
			buildCmd := exec.Command("go", "build", "-o", binaryPath, ".")
			buildCmd.Dir = projectRoot
			buildErr := buildCmd.Run()
			require.NoError(t, buildErr, "failed to build binary")
		}

		tmpHomeDir := t.TempDir()
		globalConfigPath := filepath.Join(tmpHomeDir, config.GlobalConfigFileName)

		// Write invalid JSON to cause config read to fail
		invalidJSON := `{"invalid": json}`
		err = os.WriteFile(globalConfigPath, []byte(invalidJSON), 0644)
		require.NoError(t, err)

		cmd := exec.Command(binaryPath, "list")
		cmd.Env = append(os.Environ(), "HOME="+tmpHomeDir)
		err = cmd.Run()

		if exitError, ok := err.(*exec.ExitError); ok {
			assert.Equal(t, 1, exitError.ExitCode(), "list command should exit with code 1 on config read error")
		} else if err != nil {
			t.Fatalf("unexpected error type: %v (expected ExitError)", err)
		} else {
			t.Fatalf("expected command to fail with exit code 1, but it succeeded")
		}
	})

}
