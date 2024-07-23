package router

import (
	"sync"

	"github.com/gofiber/fiber/v2"

	"github.com/Lynicis/inzibat/config"
)

type Router interface {
	CreateRoutes()
}

type Handler interface {
	CreateRoute(indexOfRoute int) func(ctx *fiber.Ctx) error
}

type MainRouter struct {
	Config        *config.Cfg
	FiberApp      *fiber.App
	MockHandler   Handler
	ClientHandler Handler
}

func NewMainRouter(
	cfg *config.Cfg,
	fiberApp *fiber.App,
	mockHandler Handler,
	clientHandler Handler,
) Router {
	return &MainRouter{
		Config:        cfg,
		FiberApp:      fiberApp,
		MockHandler:   mockHandler,
		ClientHandler: clientHandler,
	}
}

func (mainRouter *MainRouter) CreateRoutes() {
	var (
		workerCount = mainRouter.Config.Concurrency.RouteCreatorLimit
		routes      = mainRouter.Config.Routes
	)

	routeChannel := make(chan *RouteChannel, workerCount)
	defer close(routeChannel)

	var waitGroup sync.WaitGroup
	for indexOfRoute, route := range routes {
		waitGroup.Add(1)
		routeChannel <- &RouteChannel{
			Route:        route,
			IndexOfRoute: indexOfRoute,
		}
		mainRouter.routeCreatorWorker(routeChannel, &waitGroup)
	}
	waitGroup.Wait()
}

func (mainRouter *MainRouter) routeCreatorWorker(routeChannel chan *RouteChannel, waitGroup *sync.WaitGroup) {
	var (
		resultOfChannel = <-routeChannel
		routeFunction   func(ctx *fiber.Ctx) error
	)

	if resultOfChannel.Route.RequestTo.Method != "" {
		routeFunction = mainRouter.ClientHandler.CreateRoute(resultOfChannel.IndexOfRoute)
	}

	if resultOfChannel.Route.Mock.StatusCode > 0 {
		routeFunction = mainRouter.MockHandler.CreateRoute(resultOfChannel.IndexOfRoute)
	}

	if routeFunction != nil {
		mainRouter.FiberApp.Add(resultOfChannel.Route.Method, resultOfChannel.Route.Path, routeFunction)
	}

	waitGroup.Done()
}
