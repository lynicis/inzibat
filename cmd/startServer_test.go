package cmd

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lynicis/inzibat/config"
)

func TestStartServerCmd(t *testing.T) {
	t.Run("happy path - command is registered", func(t *testing.T) {
		assert.NotNil(t, startServerCmd)
		assert.Equal(t, "start", startServerCmd.Use)
		assert.Contains(t, startServerCmd.Aliases, "start-server")
		assert.Contains(t, startServerCmd.Aliases, "server")
		assert.Contains(t, startServerCmd.Aliases, "s")
	})

	t.Run("happy path - command has config flag", func(t *testing.T) {
		flag := startServerCmd.Flag("config")
		require.NotNil(t, flag)

		assert.Equal(t, "c", flag.Shorthand)
	})
}

func TestStartServerCmd_ConfigFile(t *testing.T) {
	t.Run("happy path - config file flag can be set", func(t *testing.T) {
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
					FakeResponse: &config.FakeResponse{
						StatusCode: 200,
					},
				},
			},
		}
		err := config.WriteConfig(cfg, configFile)
		require.NoError(t, err)

		err = startServerCmd.Flag("config").Value.Set(configFile)
		require.NoError(t, err)

		assert.Equal(t, configFile, configFile)
	})
}

func TestStartServerCmd_GlobalFlag(t *testing.T) {
	t.Run("happy path - command has global flag", func(t *testing.T) {
		flag := startServerCmd.Flag("global")
		require.NotNil(t, flag)

		assert.Equal(t, "g", flag.Shorthand)
	})

	t.Run("happy path - start server invoked with global flag", func(t *testing.T) {
		originalStartServerFunc := startServerFunc
		defer func() {
			startServerFunc = originalStartServerFunc
			_ = startServerCmd.Flags().Set("global", "false")
		}()

		var calledWithGlobal bool
		startServerFunc = func(_ string, isGlobal bool) error {
			calledWithGlobal = isGlobal
			return nil
		}

		err := startServerCmd.Flags().Set("global", "true")
		require.NoError(t, err)

		startServerCmd.Run(startServerCmd, []string{})

		assert.True(t, calledWithGlobal)
	})
}
