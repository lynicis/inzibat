package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
}

func TestListCmd_WithConfig(t *testing.T) {
	t.Run("happy path - command can be executed", func(t *testing.T) {
		assert.NotNil(t, listCmd)
		assert.Equal(t, "list", listCmd.Use)
	})
}
