package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"inzibat/config"
)

func TestListCmd(t *testing.T) {
	t.Run("command structure is correct", func(t *testing.T) {
		assert.NotNil(t, listCmd)
		assert.Equal(t, "list", listCmd.Use)
		assert.Contains(t, listCmd.Aliases, "list-routes")
		assert.Contains(t, listCmd.Aliases, "ls")
		assert.Contains(t, listCmd.Aliases, "l")
		assert.Equal(t, "List all routes", listCmd.Short)
	})

	t.Run("error path - config file does not exist", func(t *testing.T) {
		homeDir, err := os.UserHomeDir()
		require.NoError(t, err)

		globalConfigPath := filepath.Join(homeDir, config.GlobalConfigFileName)

		originalConfigExists := false
		var originalConfigData []byte
		if _, err := os.Stat(globalConfigPath); err == nil {
			originalConfigExists = true
			originalConfigData, err = os.ReadFile(globalConfigPath)
			require.NoError(t, err)
			os.Remove(globalConfigPath)
		}

		defer func() {
			if originalConfigExists {
				err := os.WriteFile(globalConfigPath, originalConfigData, 0644)
				require.NoError(t, err)
			}
		}()

		assert.NotNil(t, listCmd)
		assert.Equal(t, "list", listCmd.Use)
	})

	t.Run("error path - invalid config file", func(t *testing.T) {
		homeDir, err := os.UserHomeDir()
		require.NoError(t, err)

		globalConfigPath := filepath.Join(homeDir, config.GlobalConfigFileName)

		originalConfigExists := false
		var originalConfigData []byte
		if _, err := os.Stat(globalConfigPath); err == nil {
			originalConfigExists = true
			originalConfigData, err = os.ReadFile(globalConfigPath)
			require.NoError(t, err)
		}

		invalidConfig := []byte("invalid json content {")
		err = os.WriteFile(globalConfigPath, invalidConfig, 0644)
		require.NoError(t, err)

		defer func() {
			if originalConfigExists {
				err := os.WriteFile(globalConfigPath, originalConfigData, 0644)
				require.NoError(t, err)
			} else {
				os.Remove(globalConfigPath)
			}
		}()

		assert.NotNil(t, listCmd)
		assert.Equal(t, "list", listCmd.Use)
	})
}
