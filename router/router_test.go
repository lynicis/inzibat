package router

import (
	"net/http"
	"testing"

	"github.com/Lynicis/inzibat/config"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
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
						Method: http.MethodGet,
						Path:   "/",
						RequestTo: config.RequestTo{
							Host: "http://127.0.0.1:3001",
							Path: "/",
						},
					},
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
						Method: http.MethodPost,
						Path:   "/",
						RequestTo: config.RequestTo{
							Host: "http://127.0.0.1:3001",
							Path: "/",
						},
					},
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
						Method: http.MethodPut,
						Path:   "/",
						RequestTo: config.RequestTo{
							Host: "http://127.0.0.1:3001",
							Path: "/",
						},
					},
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
						Method: http.MethodDelete,
						Path:   "/",
						RequestTo: config.RequestTo{
							Host: "http://127.0.0.1:3001",
							Path: "/",
						},
					},
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
						Method: http.MethodPatch,
						Path:   "/",
						RequestTo: config.RequestTo{
							Host: "http://127.0.0.1:3001",
							Path: "/",
						},
					},
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
			}
			app := fiber.New()
			r := NewRouter(configInstance, app, nil)
			r.CreateRoutes()
		})
	})
}
