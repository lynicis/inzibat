package router

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lynicis/inzibat/config"
)

type stubHandler struct {
	body  string
	calls []int
}

func (stub *stubHandler) CreateHandler(handlerIndex int) func(ctx *fiber.Ctx) error {
	stub.calls = append(stub.calls, handlerIndex)

	return func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusOK).SendString(stub.body)
	}
}

func request(t *testing.T, fiberApp *fiber.App, method, path string) string {
	t.Helper()

	req := httptest.NewRequest(method, path, nil)

	response, err := fiberApp.Test(req)
	require.NoError(t, err)

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)
	require.NoError(t, response.Body.Close())

	return string(body)
}

func TestRouter_CreateRoutes(t *testing.T) {
	t.Run("mock route is registered", func(t *testing.T) {
		fiberApp := fiber.New()
		endpointHandler := &stubHandler{body: "mock"}
		clientHandler := &stubHandler{body: "proxy"}

		mainRouter := &MainRouter{
			Config: &config.Cfg{
				Routes: []config.Route{
					{
						Method:       fiber.MethodGet,
						Path:         "/mock",
						FakeResponse: &config.FakeResponse{StatusCode: fiber.StatusOK},
					},
				},
			},
			FiberApp:        fiberApp,
			EndpointHandler: endpointHandler,
			ClientHandler:   clientHandler,
		}

		mainRouter.CreateRoutes()

		assert.Equal(t, []int{0}, endpointHandler.calls)
		assert.Empty(t, clientHandler.calls)
		assert.Equal(t, "mock", request(t, fiberApp, fiber.MethodGet, "/mock"))
	})

	t.Run("proxy route is registered", func(t *testing.T) {
		fiberApp := fiber.New()
		endpointHandler := &stubHandler{body: "mock"}
		clientHandler := &stubHandler{body: "proxy"}

		mainRouter := &MainRouter{
			Config: &config.Cfg{
				Routes: []config.Route{
					{
						Method: fiber.MethodPost,
						Path:   "/proxy",
						RequestTo: &config.RequestTo{
							Method: fiber.MethodPost,
							Host:   "http://127.0.0.1:3001",
							Path:   "/proxy",
						},
					},
				},
			},
			FiberApp:        fiberApp,
			EndpointHandler: endpointHandler,
			ClientHandler:   clientHandler,
		}

		mainRouter.CreateRoutes()

		assert.Equal(t, []int{0}, clientHandler.calls)
		assert.Empty(t, endpointHandler.calls)
		assert.Equal(t, "proxy", request(t, fiberApp, fiber.MethodPost, "/proxy"))
	})

	t.Run("route with neither handler is skipped", func(t *testing.T) {
		fiberApp := fiber.New()
		endpointHandler := &stubHandler{}
		clientHandler := &stubHandler{}

		mainRouter := &MainRouter{
			Config: &config.Cfg{
				Routes: []config.Route{
					{Method: fiber.MethodGet, Path: "/unknown"},
					{
						Method:       fiber.MethodGet,
						Path:         "/no-method",
						RequestTo:    &config.RequestTo{},
						FakeResponse: &config.FakeResponse{},
					},
				},
			},
			FiberApp:        fiberApp,
			EndpointHandler: endpointHandler,
			ClientHandler:   clientHandler,
		}

		mainRouter.CreateRoutes()

		assert.Empty(t, endpointHandler.calls)
		assert.Empty(t, clientHandler.calls)
	})
}
