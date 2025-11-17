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

	t.Run("returns no proxy command when requestTo method is empty", func(t *testing.T) {
		route := config.Route{
			Method: fiber.MethodGet,
			Path:   "/no-method",
			RequestTo: &config.RequestTo{
				Method: "",
				Host:   "http://example.com",
			},
		}

		commands := CreateRouteCommands(route, 0, noopHandler{}, noopHandler{})

		assert.Empty(t, commands)
	})

	t.Run("returns no mock command when status code is zero", func(t *testing.T) {
		route := config.Route{
			Method:       fiber.MethodGet,
			Path:         "/no-status",
			FakeResponse: &config.FakeResponse{StatusCode: 0},
		}

		commands := CreateRouteCommands(route, 0, noopHandler{}, noopHandler{})

		assert.Empty(t, commands)
	})
}

func TestMockRouteCommand(t *testing.T) {
	t.Run("happy path - ShouldExecute returns true", func(t *testing.T) {
		command := NewMockRouteCommand(0, noopHandler{})

		result := command.ShouldExecute()

		assert.True(t, result)
	})

	t.Run("happy path - Execute returns handler", func(t *testing.T) {
		command := NewMockRouteCommand(0, noopHandler{})

		handler, err := command.Execute()

		assert.NoError(t, err)
		assert.NotNil(t, handler)
	})

	t.Run("happy path - Execute uses correct route index", func(t *testing.T) {
		command := NewMockRouteCommand(5, noopHandler{})

		handler, err := command.Execute()

		assert.NoError(t, err)
		assert.NotNil(t, handler)
	})
}

func TestProxyRouteCommand(t *testing.T) {
	t.Run("happy path - ShouldExecute returns true", func(t *testing.T) {
		command := NewProxyRouteCommand(0, noopHandler{})

		result := command.ShouldExecute()

		assert.True(t, result)
	})

	t.Run("happy path - Execute returns handler", func(t *testing.T) {
		command := NewProxyRouteCommand(0, noopHandler{})

		handler, err := command.Execute()

		assert.NoError(t, err)
		assert.NotNil(t, handler)
	})

	t.Run("happy path - Execute uses correct route index", func(t *testing.T) {
		command := NewProxyRouteCommand(3, noopHandler{})

		handler, err := command.Execute()

		assert.NoError(t, err)
		assert.NotNil(t, handler)
	})
}
