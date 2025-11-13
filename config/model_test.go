package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCfg_GetServerAddr(t *testing.T) {
	cfg := &Cfg{
		ServerPort: 8080,
	}
	result := cfg.GetServerAddr()

	assert.Equal(t, ":8080", result)
}

func TestRequestTo_GetParsedUrl(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		requestTo := RequestTo{
			Host: "http://localhost:8080",
			Path: "/test",
		}
		parsedUrl, err := requestTo.GetParsedUrl()

		assert.NoError(t, err)
		assert.Equal(t, "http://localhost:8080/test", parsedUrl.String())
	})

	t.Run("invalid host", func(t *testing.T) {
		requestTo := RequestTo{
			Host: "http://host%zz",
			Path: "/path",
		}
		_, err := requestTo.GetParsedUrl()

		assert.Error(t, err)
	})
}
