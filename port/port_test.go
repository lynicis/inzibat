package port

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFreePort(t *testing.T) {
	port, err := GetFreePort()

	assert.NoError(t, err)
	assert.Greater(t, "0", port)
}
