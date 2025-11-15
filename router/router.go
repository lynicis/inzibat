package router

import (
	"sync"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"inzibat/config"
	"inzibat/handler"
	_ "inzibat/log"
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

	commands := CreateRouteCommands(
		routeChannel.Route,
		routeChannel.RouteIndex,
		mainRouter.EndpointHandler,
		mainRouter.ClientHandler,
	)

	for _, cmd := range commands {
		if !cmd.ShouldExecute() {
			continue
		}

		routeFunction, err := cmd.Execute()
		if err != nil {
			zap.L().Warn("failed to execute route command",
				zap.String("method", routeChannel.Route.Method),
				zap.String("path", routeChannel.Route.Path),
				zap.Int("route_index", routeChannel.RouteIndex),
				zap.Error(err),
			)
			continue
		}

		if routeFunction != nil {
			mainRouter.FiberApp.Add(routeChannel.Route.Method, routeChannel.Route.Path, routeFunction)
		}
	}
}
