package router

import (
	"github.com/gofiber/fiber/v2"

	"github.com/Lynicis/inzibat/config"
)

type MockHandler struct {
	RouteConfig []config.Route
}

func NewMockHandler(routeConfig []config.Route) Handler {
	return &MockHandler{
		RouteConfig: routeConfig,
	}
}

func (mockRoute *MockHandler) CreateRoute(indexOfRoute int) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		var (
			headers    = mockRoute.RouteConfig[indexOfRoute].Mock.Headers
			body       = mockRoute.RouteConfig[indexOfRoute].Mock.Body
			bodyString = mockRoute.RouteConfig[indexOfRoute].Mock.BodyString
			statusCode = mockRoute.RouteConfig[indexOfRoute].Mock.StatusCode
		)

		if len(headers) > 0 {
			for headerKey, headerValue := range headers {
				ctx.Set(headerKey, headerValue)
			}
		}
		ctx.Status(statusCode)

		if len(bodyString) > 0 {
			return ctx.SendString(bodyString)
		}

		if len(body) > 0 {
			return ctx.JSON(body)
		}

		return nil
	}
}
