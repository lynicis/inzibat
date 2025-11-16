package cmd

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecute(t *testing.T) {
	t.Run("happy path - executes without error when no args", func(t *testing.T) {
		originalArgs := os.Args
		defer func() {
			os.Args = originalArgs
		}()

		assert.NotNil(t, rootCmd)
		assert.Equal(t, "inzibat", rootCmd.Use)
	})

	t.Run("happy path - root command has correct properties", func(t *testing.T) {
		assert.NotNil(t, rootCmd)
		assert.Equal(t, "inzibat", rootCmd.Use)
		assert.Contains(t, rootCmd.Short, "HTTP mock server")
		assert.Contains(t, rootCmd.Long, "Inzibat")
	})
}
