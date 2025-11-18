package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/goccy/go-json"

	"github.com/goccy/go-reflect"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	httpPkg "github.com/lynicis/inzibat/client/http"
	"github.com/lynicis/inzibat/config"
)

type ClientHandler struct {
	Client      *httpPkg.Client
	RouteConfig *[]config.Route
}

func (clientRoute *ClientHandler) CreateHandler(routeIndex int) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		requestTo := (*clientRoute.RouteConfig)[routeIndex].RequestTo
		parsedUrl, err := requestTo.GetParsedUrl()
		if err != nil {
			return fmt.Errorf("failed to parse request URL: %w", err)
		}

		bodyBytes, err := json.Marshal(&requestTo.Body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}

		methodArguments := clientRoute.prepareMethodArguments(
			parsedUrl.String(),
			requestTo.Headers,
			bodyBytes,
			requestTo.Method,
		)
		response, err := clientRoute.executeHttpMethod(requestTo.Method, methodArguments)
		if err != nil {
			if requestTo.InErrorReturn500 {
				ctx.Status(fiber.StatusInternalServerError)
				return ctx.Send(nil)
			}

			return ctx.
				Status(fiber.StatusInternalServerError).
				SendString(err.Error())
		}

		return ctx.
			Status(response.Status).
			Send(response.Body)
	}
}

func (clientRoute *ClientHandler) prepareMethodArguments(
	url string,
	headers http.Header,
	bodyBytes []byte,
	method string,
) []reflect.Value {
	arguments := []reflect.Value{
		reflect.ValueOf(url),
		reflect.ValueOf(headers),
		reflect.ValueOf(bodyBytes),
	}

	if method == fiber.MethodGet {
		arguments = arguments[:len(arguments)-1]
	}

	return arguments
}

func (clientRoute *ClientHandler) executeHttpMethod(
	method string,
	arguments []reflect.Value,
) (*httpPkg.Response, error) {
	methodName := cases.Title(language.Und).String(strings.ToLower(method))
	requestMethod := reflect.ValueOf(clientRoute.Client).MethodByName(methodName)
	returnedArguments := requestMethod.Call(arguments)

	response, ok := returnedArguments[0].Interface().(*httpPkg.Response)
	if !ok {
		return nil, fmt.Errorf("failed to cast response to http.Response for method %s", method)
	}

	var err error
	if !returnedArguments[1].IsNil() {
		err = returnedArguments[1].Interface().(error)
	}

	return response, err
}
