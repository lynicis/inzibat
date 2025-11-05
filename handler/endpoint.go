package handler

import (
	"inzibat/config"

	"github.com/gofiber/fiber/v2"
)

type EndpointHandler struct {
	RouteConfig *[]config.Route
}

func (mockRoute *EndpointHandler) CreateHandler(routeIndex int) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		resp := (*mockRoute.RouteConfig)[routeIndex].FakeResponse
		ctx = ctx.Status(resp.StatusCode)

		if len(resp.Headers) > 0 {
			for headerKey, headerValue := range resp.Headers {
				ctx.Set(headerKey, headerValue)
			}
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
