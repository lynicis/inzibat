package router

import (
	"github.com/gofiber/fiber/v2"

	"inzibat/config"
)

type RouteCommand interface {
	Execute() (fiber.Handler, error)
	ShouldExecute() bool
}

type MockRouteCommand struct {
	routeIndex      int
	endpointHandler Handler
}

func NewMockRouteCommand(routeIndex int, endpointHandler Handler) *MockRouteCommand {
	return &MockRouteCommand{
		routeIndex:      routeIndex,
		endpointHandler: endpointHandler,
	}
}

func (c *MockRouteCommand) ShouldExecute() bool {
	return true
}

func (c *MockRouteCommand) Execute() (fiber.Handler, error) {
	return c.endpointHandler.CreateHandler(c.routeIndex), nil
}

type ProxyRouteCommand struct {
	routeIndex    int
	clientHandler Handler
}

func NewProxyRouteCommand(routeIndex int, clientHandler Handler) *ProxyRouteCommand {
	return &ProxyRouteCommand{
		routeIndex:    routeIndex,
		clientHandler: clientHandler,
	}
}

func (c *ProxyRouteCommand) ShouldExecute() bool {
	return true
}

func (c *ProxyRouteCommand) Execute() (fiber.Handler, error) {
	return c.clientHandler.CreateHandler(c.routeIndex), nil
}

func CreateRouteCommands(
	route config.Route,
	routeIndex int,
	endpointHandler Handler,
	clientHandler Handler,
) []RouteCommand {
	var commands []RouteCommand

	if route.RequestTo != nil && route.RequestTo.Method != "" {
		commands = append(commands, NewProxyRouteCommand(routeIndex, clientHandler))
	}

	if route.FakeResponse != nil && route.FakeResponse.StatusCode > 0 {
		commands = append(commands, NewMockRouteCommand(routeIndex, endpointHandler))
	}

	return commands
}
