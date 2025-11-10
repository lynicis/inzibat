package router

import (
	"sync"

	"github.com/gofiber/fiber/v2"

	"inzibat/config"
	"inzibat/handler"
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

// TODO:
// make it open/closed
// probably command pattern
func (mainRouter *MainRouter) processRoute(
	routeChannel *handler.RouteChannel,
	waitGroup *sync.WaitGroup,
) {
	defer waitGroup.Done()

	var routeFunction fiber.Handler
	if routeChannel.Route.RequestTo.Method != "" {
		routeFunction = mainRouter.ClientHandler.CreateHandler(routeChannel.RouteIndex)
	}

	if routeChannel.Route.FakeResponse.StatusCode > 0 {
		routeFunction = mainRouter.EndpointHandler.CreateHandler(routeChannel.RouteIndex)
	}

	if routeFunction != nil {
		mainRouter.FiberApp.Add(routeChannel.Route.Method, routeChannel.Route.Path, routeFunction)
	}
}
