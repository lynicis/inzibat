package handler

import (
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"

	httpPkg "inzibat/client/http"
	"inzibat/config"
)

func TestClientHandler_CreateEndpoint(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		httpClient := httpPkg.NewHttpClient()
		c := &ClientHandler{
			Client: httpClient,
			RouteConfig: &[]config.Route{
				{
					Method: http.MethodGet,
					Path:   "/test",
					RequestTo: config.RequestTo{
						Method:                 http.MethodGet,
						Headers:                nil,
						Body:                   nil,
						Host:                   "http://localhost:8081",
						Path:                   "",
						PassWithRequestBody:    false,
						PassWithRequestHeaders: false,
						InErrorReturn500:       false,
					},
				},
			},
		}

		handler := c.CreateHandler(0)
		err := handler(&fiber.Ctx{})

		assert.NoError(t, err)
	})
}
