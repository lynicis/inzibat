package router

import (
	"errors"
	"strings"

	"github.com/goccy/go-reflect"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/Lynicis/inzibat/client"
	"github.com/Lynicis/inzibat/config"
)

type Router interface {
	CreateRoutes()
	HandleClientMethod(routeConfig *config.Route) func(ctx *fiber.Ctx) error
}

type router struct {
	config *config.Config
	app    *fiber.App
	client client.Client
}

func NewRouter(config *config.Config, app *fiber.App, client client.Client) Router {
	return &router{
		config: config,
		app:    app,
		client: client,
	}
}

func (r *router) CreateRoutes() {
	workerCount := r.config.Concurrency.RouteCreatorLimit
	routeChannel := make(chan config.Route, workerCount)
	defer close(routeChannel)

	for _, route := range r.config.Routes {
		routeChannel <- route
		r.routeCreatorWorker(routeChannel)
	}
}

func (r *router) routeCreatorWorker(routeChannel chan config.Route) {
	route := <-routeChannel
	r.app.Add(route.Method, route.Path, r.HandleClientMethod(&route))
}

func (r *router) HandleClientMethod(routeConfig *config.Route) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		var ConfigRequestToMethod = routeConfig.RequestTo.Method
		cloneOfClientStruct := r.client

		methodName := cases.Title(language.Und).String(strings.ToLower(ConfigRequestToMethod))
		requestMethod := reflect.ValueOf(cloneOfClientStruct).MethodByName(methodName)

		var (
			ConfigRequestToURL = routeConfig.RequestTo.Host + routeConfig.RequestTo.Path
			RequestWithHeaders map[string]string
			RequestWithBody    = ctx.Body()
		)

		if routeConfig.RequestTo.Headers != nil {
			RequestWithHeaders = routeConfig.RequestTo.Headers
		}

		methodArgumentsForClient := []reflect.Value{
			reflect.ValueOf(ConfigRequestToURL),
			reflect.ValueOf(RequestWithHeaders),
			reflect.ValueOf(RequestWithBody),
		}

		if ConfigRequestToMethod == fiber.MethodGet {
			methodArgumentsForClient = methodArgumentsForClient[:len(methodArgumentsForClient)-1]
		}

		returnedArguments := requestMethod.Call(methodArgumentsForClient)

		var isSafeToGetReturnArguments bool
		var returnedHttpResponse *client.HttpResponse
		returnedHttpResponse, isSafeToGetReturnArguments = returnedArguments[0].Interface().(*client.HttpResponse)
		if !isSafeToGetReturnArguments {
			return errors.New(ErrorTypeCasting)
		}

		returnedError := returnedArguments[1].Interface()
		if returnedError != nil {
			return returnedError.(error)
		}

		return ctx.Status(returnedHttpResponse.Status).Send(returnedHttpResponse.Body)
	}
}
