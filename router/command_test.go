package router

import (
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"

	"inzibat/config"
)

type noopHandler struct{}

func (n noopHandler) CreateHandler(handlerIndex int) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		return nil
	}
}

func TestCreateRouteCommands(t *testing.T) {
	t.Run("returns proxy command when requestTo is configured", func(t *testing.T) {
		route := config.Route{
			Method: fiber.MethodGet,
			Path:   "/proxy-only",
			RequestTo: &config.RequestTo{
				Method: http.MethodGet,
				Host:   "http://example.com",
				Path:   "/proxy-only",
			},
		}

		commands := CreateRouteCommands(route, 0, noopHandler{}, noopHandler{})

		assert.Len(t, commands, 1)
		_, isProxy := commands[0].(*ProxyRouteCommand)
		assert.True(t, isProxy)
	})

	t.Run("returns mock command when fake response is configured", func(t *testing.T) {
		route := config.Route{
			Method:       fiber.MethodGet,
			Path:         "/mock-only",
			FakeResponse: &config.FakeResponse{StatusCode: http.StatusOK},
		}

		commands := CreateRouteCommands(route, 0, noopHandler{}, noopHandler{})

		assert.Len(t, commands, 1)
		_, isMock := commands[0].(*MockRouteCommand)
		assert.True(t, isMock)
	})

	t.Run("returns both commands when request and fake response configured", func(t *testing.T) {
		route := config.Route{
			Method: fiber.MethodGet,
			Path:   "/both",
			RequestTo: &config.RequestTo{
				Method: http.MethodPost,
				Host:   "http://example.com",
				Path:   "/both",
			},
			FakeResponse: &config.FakeResponse{StatusCode: http.StatusCreated},
		}

		commands := CreateRouteCommands(route, 0, noopHandler{}, noopHandler{})

		assert.Len(t, commands, 2)
		assert.IsType(t, &ProxyRouteCommand{}, commands[0])
		assert.IsType(t, &MockRouteCommand{}, commands[1])
	})

	t.Run("returns no commands when neither request nor fake response configured", func(t *testing.T) {
		route := config.Route{
			Method: fiber.MethodGet,
			Path:   "/missing",
		}

		commands := CreateRouteCommands(route, 0, noopHandler{}, noopHandler{})

		assert.Empty(t, commands)
	})
}
