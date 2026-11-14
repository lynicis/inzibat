package handler

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/lynicis/inzibat/config"
)

type EndpointHandler struct {
	RouteConfig *[]config.Route
}

func (mockRoute *EndpointHandler) CreateHandler(routeIndex int) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		resp := (*mockRoute.RouteConfig)[routeIndex].FakeResponse
		ctx.Status(resp.StatusCode)

		for headerKey, headerValue := range resp.Headers {
			ctx.Set(headerKey, strings.Join(headerValue, ","))
		}

		if len(resp.BodyString) > 0 {
			return ctx.SendString(resp.BodyString)
		}

		if len(resp.Body) > 0 {
			return ctx.JSON(resp.Body)
		}

		return nil
	}
}
