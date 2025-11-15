package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStartServerCmd(t *testing.T) {
	t.Run("command structure is correct", func(t *testing.T) {
		assert.NotNil(t, startServerCmd)
		assert.Equal(t, "start", startServerCmd.Use)
		assert.Contains(t, startServerCmd.Aliases, "start-server")
		assert.Contains(t, startServerCmd.Aliases, "server")
		assert.Contains(t, startServerCmd.Aliases, "s")
		assert.Equal(t, "Start the Inzibat mock server", startServerCmd.Short)
		assert.NotEmpty(t, startServerCmd.Long)
	})
}

