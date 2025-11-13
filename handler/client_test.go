package handler

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	httpPkg "inzibat/client/http"
	"inzibat/config"
)

func TestClientHandler_CreateEndpoint(t *testing.T) {
	port, err := httpPkg.GetFreePort()
	require.NoError(t, err)

	host := fmt.Sprintf("http://127.0.0.1:%d", port)

	t.Run("happy path", func(t *testing.T) {
		httpClient := httpPkg.NewHttpClient()
		c := &ClientHandler{
			Client: httpClient,
			RouteConfig: &[]config.Route{
				{
					Method: http.MethodGet,
					Path:   "/test",
					RequestTo: config.RequestTo{
						Method: http.MethodGet,
						Headers: http.Header{
							"Test-Header": {"Test"},
						},
						Body: config.HttpBody{
							"Test": "Test",
						},
						Host:                   host,
						Path:                   "/test",
						PassWithRequestBody:    false,
						PassWithRequestHeaders: false,
						InErrorReturn500:       false,
					},
				},
			},
		}

		handler := c.CreateHandler(0)
		srv := fiber.New()
		srv.Get("/test", func(ctx *fiber.Ctx) error {
			return ctx.Status(200).SendString("OK")
		})

		go func() {
			err = srv.Listen(fmt.Sprintf(":%d", port))
			require.NoError(t, err)
		}()
		srv.ShutdownWithContext(t.Context())

		srv.Get("/mock", func(ctx *fiber.Ctx) error {
			err = handler(ctx)
			assert.NoError(t, err)

			return nil
		})
	})
}
