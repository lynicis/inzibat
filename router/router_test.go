package router

import (
	"bytes"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Lynicis/inzibat/client"
	"github.com/Lynicis/inzibat/config"
	"github.com/Lynicis/inzibat/test-utils"
)

func TestNewRouter(t *testing.T) {
	r := NewRouter(nil, nil, nil)
	assert.Implements(t, (*Router)(nil), r)
}

func TestRouter_CreateRoutes(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		t.Run("GET Method", func(t *testing.T) {
			configInstance := &config.Config{
				ServerPort: "3000",
				Routes: []config.Route{
					{
						Method: fiber.MethodGet,
						Path:   "/",
						RequestTo: config.RequestTo{
							Host: "http://127.0.0.1:3001",
							Path: "/",
						},
					},
				},
				Concurrency: config.Concurrency{
					RouteCreatorLimit: 1,
				},
			}
			app := fiber.New()
			r := NewRouter(configInstance, app, nil)
			r.CreateRoutes()

			assert.Len(t, app.GetRoutes(), 1)
		})

		t.Run("POST Method", func(t *testing.T) {
			configInstance := &config.Config{
				ServerPort: "3000",
				Routes: []config.Route{
					{
						Method: fiber.MethodPost,
						Path:   "/",
						RequestTo: config.RequestTo{
							Host: "http://127.0.0.1:3001",
							Path: "/",
						},
					},
				},
				Concurrency: config.Concurrency{
					RouteCreatorLimit: 1,
				},
			}
			app := fiber.New()
			r := NewRouter(configInstance, app, nil)
			r.CreateRoutes()

			assert.Len(t, app.GetRoutes(), 1)
		})

		t.Run("PUT Method", func(t *testing.T) {
			configInstance := &config.Config{
				ServerPort: "3000",
				Routes: []config.Route{
					{
						Method: fiber.MethodPut,
						Path:   "/",
						RequestTo: config.RequestTo{
							Host: "http://127.0.0.1:3001",
							Path: "/",
						},
					},
				},
				Concurrency: config.Concurrency{
					RouteCreatorLimit: 1,
				},
			}
			app := fiber.New()
			r := NewRouter(configInstance, app, nil)
			r.CreateRoutes()

			assert.Len(t, app.GetRoutes(), 1)
		})

		t.Run("POST Method", func(t *testing.T) {
			configInstance := &config.Config{
				ServerPort: "3000",
				Routes: []config.Route{
					{
						Method: fiber.MethodDelete,
						Path:   "/",
						RequestTo: config.RequestTo{
							Host: "http://127.0.0.1:3001",
							Path: "/",
						},
					},
				},
				Concurrency: config.Concurrency{
					RouteCreatorLimit: 1,
				},
			}
			app := fiber.New()
			r := NewRouter(configInstance, app, nil)
			r.CreateRoutes()

			assert.Len(t, app.GetRoutes(), 1)
		})

		t.Run("PATCH Method", func(t *testing.T) {
			configInstance := &config.Config{
				ServerPort: "3000",
				Routes: []config.Route{
					{
						Method: fiber.MethodPatch,
						Path:   "/",
						RequestTo: config.RequestTo{
							Host: "http://127.0.0.1:3001",
							Path: "/",
						},
					},
				},
				Concurrency: config.Concurrency{
					RouteCreatorLimit: 1,
				},
			}
			app := fiber.New()
			r := NewRouter(configInstance, app, nil)
			r.CreateRoutes()

			assert.Len(t, app.GetRoutes(), 1)
		})
	})

	t.Run("when client get malicious http method", func(t *testing.T) {
		assert.Panics(t, func() {
			configInstance := &config.Config{
				ServerPort: "3000",
				Routes: []config.Route{
					{
						Method: "MALICIOUS",
						Path:   "/",
						RequestTo: config.RequestTo{
							Host: "http://127.0.0.1:3001",
							Path: "/",
						},
					},
				},
				Concurrency: config.Concurrency{
					RouteCreatorLimit: 1,
				},
			}
			app := fiber.New()
			r := NewRouter(configInstance, app, nil)
			r.CreateRoutes()
		})
	})
}

func TestRouter_HandleClientMethod(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		t.Run("GET Method", func(t *testing.T) {
			mockServer := fiber.New(fiber.Config{
				DisableStartupMessage: true,
			})
			mockServer.Get("/user", func(ctx *fiber.Ctx) error {
				return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
					"nick":     "lynicis",
					"password": 12345,
				})
			})

			mockServerPort, err := testUtils.GetFreePort()
			require.NoError(t, err)

			go mockServer.Listen(fmt.Sprintf(":%s", mockServerPort))
			defer mockServer.Shutdown()
			time.Sleep(1 * time.Second)

			port, err := testUtils.GetFreePort()
			require.NoError(t, err)

			app := fiber.New()
			cfg := &config.Config{
				ServerPort: port,
				Routes: []config.Route{
					{
						Method: fiber.MethodGet,
						Path:   "/user",
						RequestTo: config.RequestTo{
							Method: fiber.MethodGet,
							Host:   fmt.Sprintf("http://localhost:%s", mockServerPort),
							Path:   "/user",
						},
					},
				},
				Concurrency: config.Concurrency{
					RouteCreatorLimit: 1,
				},
			}

			router := NewRouter(
				cfg,
				app,
				client.NewClient(),
			)
			firstRoute := cfg.Routes[0]
			handler := router.HandleClientMethod(&firstRoute)
			app.Get("/user", handler)

			request := httptest.NewRequest(fiber.MethodGet, "/user", nil)
			response, err := app.Test(request)
			require.NoError(t, err)

			responseBody, err := io.ReadAll(response.Body)
			require.NoError(t, err)

			var body map[string]interface{}
			err = json.Unmarshal(responseBody, &body)
			require.NoError(t, err)

			assert.Equal(t, fiber.StatusOK, response.StatusCode)
			assert.Equal(t, "lynicis", body["nick"])
			assert.Equal(t, float64(12345), body["password"])
		})

		t.Run("Other HTTP Method", func(t *testing.T) {
			mockServer := fiber.New(fiber.Config{
				DisableStartupMessage: true,
			})
			mockServer.Post("/user", func(ctx *fiber.Ctx) error {
				return ctx.SendStatus(fiber.StatusCreated)
			})

			mockServerPort, err := testUtils.GetFreePort()
			require.NoError(t, err)

			go mockServer.Listen(fmt.Sprintf(":%s", mockServerPort))
			defer mockServer.Shutdown()
			time.Sleep(1 * time.Second)

			port, err := testUtils.GetFreePort()
			require.NoError(t, err)

			app := fiber.New()
			cfg := &config.Config{
				ServerPort: port,
				Routes: []config.Route{
					{
						Method: fiber.MethodPost,
						Path:   "/register",
						RequestTo: config.RequestTo{
							Method: fiber.MethodPost,
							Host:   fmt.Sprintf("http://localhost:%s", mockServerPort),
							Path:   "/user",
						},
					},
				},
				Concurrency: config.Concurrency{
					RouteCreatorLimit: 1,
				},
			}

			router := NewRouter(
				cfg,
				app,
				client.NewClient(),
			)
			firstRoute := cfg.Routes[0]
			handler := router.HandleClientMethod(&firstRoute)
			app.Post("/register", handler)

			requestBody := bytes.NewBufferString(`{"nick":"lynicis","password":"1234"}"`)
			request := httptest.NewRequest(fiber.MethodPost, "/register", requestBody)
			response, err := app.Test(request)
			require.NoError(t, err)

			assert.Equal(t, fiber.StatusCreated, response.StatusCode)
		})
	})

	t.Run("when client return error and inErrorReturn500 is true should return ctx with 500 status code", func(t *testing.T) {
		port, err := testUtils.GetFreePort()
		require.NoError(t, err)

		app := fiber.New()
		cfg := &config.Config{
			ServerPort: port,
			Routes: []config.Route{
				{
					Method: fiber.MethodGet,
					Path:   "/user",
					RequestTo: config.RequestTo{
						Method:           fiber.MethodGet,
						Host:             "http://localhost:12345",
						Path:             "/user",
						InErrorReturn500: true,
					},
				},
			},
			Concurrency: config.Concurrency{
				RouteCreatorLimit: 1,
			},
		}

		router := NewRouter(
			cfg,
			app,
			client.NewClient(),
		)
		firstRoute := cfg.Routes[0]
		handler := router.HandleClientMethod(&firstRoute)
		app.Get("/user", handler)

		request := httptest.NewRequest(fiber.MethodGet, "/user", nil)
		response, err := app.Test(request)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusInternalServerError, response.StatusCode)
	})
}

func TestRouter_HandleMockMethod(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		t.Run("GET Method", func(t *testing.T) {
			router := NewRouter(nil, nil, nil)
			route := config.Route{
				Method: fiber.MethodGet,
				Path:   "/user",
				Mock: config.Mock{
					Headers: map[string]string{
						"X-Test-Header": "test-header-value",
					},
					Body: map[string]interface{}{
						"nick":     "lynicis",
						"password": 12345,
					},
					Status: fiber.StatusOK,
				},
			}
			handler := router.HandleMockMethod(&route)

			app := fiber.New()
			app.Get("/user", handler)

			request := httptest.NewRequest(fiber.MethodGet, "/user", nil)
			response, err := app.Test(request)
			require.NoError(t, err)

			responseBody, err := io.ReadAll(response.Body)
			require.NoError(t, err)

			var body map[string]interface{}
			err = json.Unmarshal(responseBody, &body)
			require.NoError(t, err)

			assert.Equal(t, fiber.StatusOK, response.StatusCode)
			assert.Equal(t, "lynicis", body["nick"])
			assert.Equal(t, float64(12345), body["password"])
			assert.Equal(t, "test-header-value", response.Header["X-Test-Header"][0])
		})

		t.Run("Other HTTP Method", func(t *testing.T) {
			router := NewRouter(nil, nil, nil)
			firstRoute := config.Route{
				Method: fiber.MethodPost,
				Path:   "/user",
				Mock: config.Mock{
					Headers: map[string]string{
						"X-Test-Header": "test-header-value",
					},
					Body: map[string]interface{}{
						"token": "abcd.abcd.abcd",
					},
					Status: fiber.StatusCreated,
				},
			}
			handler := router.HandleMockMethod(&firstRoute)

			app := fiber.New()
			app.Post("/user", handler)

			requestBody := bytes.NewBufferString(`{"nick":"lynicis","password":"1234"}"`)
			request := httptest.NewRequest(fiber.MethodPost, "/user", requestBody)
			response, err := app.Test(request)
			require.NoError(t, err)

			responseBody, err := io.ReadAll(response.Body)
			require.NoError(t, err)

			var body map[string]interface{}
			err = json.Unmarshal(responseBody, &body)
			require.NoError(t, err)

			assert.Equal(t, fiber.StatusCreated, response.StatusCode)
			assert.Equal(t, "abcd.abcd.abcd", body["token"])
			assert.Equal(t, "test-header-value", response.Header["X-Test-Header"][0])
		})
	})
}
