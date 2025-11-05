package router

import (
  "testing"

  "github.com/gofiber/fiber/v2"
  "github.com/stretchr/testify/assert"

  "inzibat/client/http"
  "inzibat/config"
  "inzibat/handler"
)

func TestRouter_CreateRoute(t *testing.T) {
  t.Run("happy path", func(t *testing.T) {
    t.Run("GET Method", func(t *testing.T) {
      fiberApp := fiber.New()
      router := &MainRouter{
        Config: &config.Cfg{
          ServerPort: 3000,
          Routes: []config.Route{
            {
              Method: fiber.MethodGet,
              Path:   "/",
              RequestTo: config.RequestTo{
                Method: fiber.MethodGet,
                Host:   "http://127.0.0.1:3001",
                Path:   "/",
              },
            },
          },
          Concurrency: config.Concurrency{
            RouteCreatorLimit: 1,
          },
        },
        FiberApp: fiberApp,
        ClientHandler: &handler.ClientHandler{
          Client:      http.NewHttpClient(),
          RouteConfig: nil,
        },
      }
      router.CreateRoutes()

      assert.Len(t, fiberApp.GetRoutes(), 1)
    })

    t.Run("POST Method", func(t *testing.T) {
      fiberApp := fiber.New()
      router := &MainRouter{
        Config: &config.Cfg{
          ServerPort: 3000,
          Routes: []config.Route{
            {
              Method: fiber.MethodPost,
              Path:   "/",
              RequestTo: config.RequestTo{
                Method: fiber.MethodGet,
                Host:   "http://127.0.0.1:3001",
                Path:   "/",
              },
            },
          },
          Concurrency: config.Concurrency{
            RouteCreatorLimit: 1,
          },
        },
        FiberApp:        nil,
        EndpointHandler: &handler.EndpointHandler{},
        ClientHandler:   &handler.ClientHandler{},
      }
      router.CreateRoutes()

      assert.Len(t, fiberApp.GetRoutes(), 1)
    })

    t.Run("PUT Method", func(t *testing.T) {
      fiberApp := fiber.New()
      router := &MainRouter{
        Config: &config.Cfg{
          ServerPort: 3000,
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
        },
        FiberApp:        fiberApp,
        EndpointHandler: &handler.EndpointHandler{},
        ClientHandler:   &handler.ClientHandler{},
      }
      router.CreateRoutes()

      assert.Len(t, fiberApp.GetRoutes(), 1)
    })

    t.Run("POST Method", func(t *testing.T) {
      fiberApp := fiber.New()
      router := &MainRouter{
        Config: &config.Cfg{
          ServerPort: 3000,
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
        },
        FiberApp:        fiberApp,
        EndpointHandler: &handler.EndpointHandler{},
        ClientHandler:   &handler.ClientHandler{},
      }
      router.CreateRoutes()

      assert.Len(t, fiberApp.GetRoutes(), 1)
    })

    t.Run("PATCH Method", func(t *testing.T) {
      fiberApp := fiber.New()
      r := &MainRouter{
        Config: &config.Cfg{
          ServerPort: 3000,
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
        },
        FiberApp:        fiberApp,
        EndpointHandler: &handler.EndpointHandler{},
        ClientHandler:   &handler.ClientHandler{},
      }
      r.CreateRoutes()

      assert.Len(t, fiberApp.GetRoutes(), 1)
    })
  })

  t.Run("when client get malicious http method", func(t *testing.T) {
    assert.Panics(t, func() {
      fiberApp := fiber.New()
      router := &MainRouter{
        Config: &config.Cfg{
          ServerPort: 3000,
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
        },
        FiberApp:        fiberApp,
        EndpointHandler: &handler.EndpointHandler{},
        ClientHandler:   &handler.ClientHandler{},
      }
      router.CreateRoutes()
    })
  })
}
