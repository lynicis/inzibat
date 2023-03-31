package main

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestServer(t *testing.T) {
	t.Run("should create server instance and return server instance", func(t *testing.T) {
		cfg := &Config{
			ServerPort: "8080",
		}

		testServer := NewServer(cfg)

		assert.IsType(t, &server{}, testServer)
	})

	t.Run("should server start and stop", func(t *testing.T) {
		cfg := &Config{
			ServerPort: "8080",
		}

		testServer := NewServer(cfg)

		go func() {
			err := testServer.Start()
			assert.NoError(t, err)
		}()

		err := testServer.Shutdown()
		assert.NoError(t, err)
	})
}

func TestServer_GetFiberInstance(t *testing.T) {
	testServer := &server{
		fiber: fiber.New(),
	}
	fiberInstance := testServer.GetFiberInstance()

	assert.IsType(t, fiberInstance, testServer.fiber)
}
