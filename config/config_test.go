package config

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadConfig(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		wd, err := os.Getwd()
		require.NoError(t, err)

		folderPath := filepath.Join(wd, "../testdata")
		cfg, err := ReadConfig(folderPath, "test.json")

		assert.NoError(t, err)
		assert.Equal(t, &Config{
			ServerPort: "8080",
			Routes: []Route{
				{
					Method: fiber.MethodGet,
					Path:   "/",
					RequestTo: RequestTo{
						Method: http.MethodGet,
						Host:   "http://localhost:8081",
						Path:   "/health",
					},
				},
			},
		}, cfg)
	})

	t.Run("config file doesn't exist", func(t *testing.T) {
		folderPath := filepath.Base("./testdata")
		cfg, err := ReadConfig(folderPath, "notfound")

		assert.Empty(t, cfg)
		assert.Error(t, err)
	})

	t.Run("config file unmarshalling error", func(t *testing.T) {
		folderPath := filepath.Base("./testdata")
		cfg, err := ReadConfig(folderPath, "unmarshalling-error")

		assert.Empty(t, cfg)
		assert.Error(t, err)
	})
}

func TestConfig_Print(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)

	path := filepath.Join(wd, "../testdata")
	cfg, err := ReadConfig(path, "test")
	cfg.Print()

	assert.NoError(t, err)
}
