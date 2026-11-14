package router

import (
	"github.com/gofiber/fiber/v2"

	"github.com/lynicis/inzibat/config"
	_ "github.com/lynicis/inzibat/log"
)

type Handler interface {
	CreateHandler(handlerIndex int) func(ctx *fiber.Ctx) error
}

type MainRouter struct {
	Config          *config.Cfg
	FiberApp        *fiber.App
	EndpointHandler Handler
	ClientHandler   Handler
}

func (mainRouter *MainRouter) CreateRoutes() {
	for routeIndex, route := range mainRouter.Config.Routes {
		mainRouter.processRoute(route, routeIndex)
	}
}

func (mainRouter *MainRouter) processRoute(route config.Route, routeIndex int) {
	if route.RequestTo != nil && route.RequestTo.Method != "" {
		routeFunction := mainRouter.ClientHandler.CreateHandler(routeIndex)
		mainRouter.FiberApp.Add(route.Method, route.Path, routeFunction)

		return
	}

	if route.FakeResponse != nil && route.FakeResponse.StatusCode > 0 {
		routeFunction := mainRouter.EndpointHandler.CreateHandler(routeIndex)
		mainRouter.FiberApp.Add(route.Method, route.Path, routeFunction)
	}
}
