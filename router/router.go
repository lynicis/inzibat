package router

import (
	"sync"

	"github.com/gofiber/fiber/v2"

	"github.com/lynicis/inzibat/config"
	"github.com/lynicis/inzibat/handler"
	_ "github.com/lynicis/inzibat/log"
)

type Router interface {
	CreateRoutes()
}

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
	routeCount := len(mainRouter.Config.Routes)
	routeChannel := make(chan *handler.RouteChannel, routeCount)
	defer close(routeChannel)

	var waitGroup sync.WaitGroup
	waitGroup.Add(routeCount)

	for workerCount := 0; workerCount < mainRouter.Config.Concurrency; workerCount++ {
		go func() {
			for route := range routeChannel {
				mainRouter.processRoute(route, &waitGroup)
			}
		}()
	}

	for routeIndex, route := range mainRouter.Config.Routes {
		routeChannel <- &handler.RouteChannel{
			Route:      route,
			RouteIndex: routeIndex,
		}
	}

	waitGroup.Wait()
}

func (mainRouter *MainRouter) processRoute(
	routeChannel *handler.RouteChannel,
	waitGroup *sync.WaitGroup,
) {
	defer waitGroup.Done()

	route := routeChannel.Route
	routeIndex := routeChannel.RouteIndex

	if route.RequestTo != nil && route.RequestTo.Method != "" {
		routeFunction := mainRouter.ClientHandler.CreateHandler(routeIndex)
		mainRouter.FiberApp.Add(route.Method, route.Path, routeFunction)
	}

	if route.FakeResponse != nil && route.FakeResponse.StatusCode > 0 {
		routeFunction := mainRouter.EndpointHandler.CreateHandler(routeIndex)
		mainRouter.FiberApp.Add(route.Method, route.Path, routeFunction)
	}
}
