package router

import (
	"github.com/Lynicis/inzibat/client"
	"github.com/Lynicis/inzibat/config"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"net/http"
	"testing"
)

func TestNewClientHandler(t *testing.T) {
	c := NewClientHandler(nil, nil)
	assert.Implements(t, (*Handler)(nil), c)
}

func TestClientHandler_CreateRoute(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		httpClient := &client.HttpClient{
			FasthttpClient: &fasthttp.Client{},
		}
		cfg := []config.Route{
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
		}

		c := NewClientHandler(httpClient, cfg)

		handler := c.CreateRoute(0)
		err := handler(&fiber.Ctx{})

		assert.NoError(t, err)
	})
}
