package cmd

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"inzibat/config"
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
					FakeResponse: config.FakeResponse{
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
