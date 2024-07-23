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

type ClientHandler struct {
	Client      client.Client
	RouteConfig []config.Route
}

func NewClientHandler(httpClient client.Client, routeConfig []config.Route) Handler {
	return &ClientHandler{
		Client:      httpClient,
		RouteConfig: routeConfig,
	}
}

func (clientRoute *ClientHandler) CreateRoute(indexOfRoute int) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		var (
			requestTo        = clientRoute.RouteConfig[indexOfRoute].RequestTo
			method           = requestTo.Method
			url              = requestTo.Host + requestTo.Path
			headers          = requestTo.Headers
			body             = ctx.Body()
			inErrorReturn500 = requestTo.InErrorReturn500
		)

		methodName := cases.Title(language.Und).String(strings.ToLower(method))
		requestMethod := reflect.ValueOf(clientRoute.Client).MethodByName(methodName)

		methodArgumentsForClient := []reflect.Value{
			reflect.ValueOf(url),
			reflect.ValueOf(headers),
			reflect.ValueOf(body),
		}

		if method == fiber.MethodGet {
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
			if inErrorReturn500 {
				return ctx.SendStatus(fiber.StatusInternalServerError)
			}

			return returnedError.(error)
		}

		return ctx.Status(returnedHttpResponse.Status).Send(returnedHttpResponse.Body)
	}
}
