package handler

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/goccy/go-reflect"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	httpPkg "inzibat/client/http"
	"inzibat/config"
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
			return err
		}

		var bodyBytes []byte
		if len(requestTo.Body) > 0 {
			bodyBytes, err = json.Marshal(requestTo.Body)
			if err != nil {
				return err
			}
		}

		methodArgumentsForClient := []reflect.Value{
			reflect.ValueOf(parsedUrl.String()),
			reflect.ValueOf(requestTo.Headers),
			reflect.ValueOf(bodyBytes),
		}

		if requestTo.Method == fiber.MethodGet {
			methodArgumentsForClient = methodArgumentsForClient[:len(methodArgumentsForClient)-1]
		}

		methodName := cases.Title(language.Und).String(strings.ToLower(requestTo.Method))
		requestMethod := reflect.ValueOf(clientRoute.Client).MethodByName(methodName)
		returnedArguments := requestMethod.Call(methodArgumentsForClient)

		returnedHttpResponse, ok := returnedArguments[0].Interface().(*httpPkg.Response)
		if !ok {
			return errors.New("failed to cast response to http.Response")
		}

		var returnedError error
		if !returnedArguments[1].IsNil() {
			returnedError = returnedArguments[1].Interface().(error)
		}

		if returnedError != nil {
			if requestTo.InErrorReturn500 {
				ctx.Status(fiber.StatusInternalServerError)
				return ctx.Send(nil)
			}

			return ctx.
				Status(fiber.StatusInternalServerError).
				SendString(returnedError.Error())
		}

		return ctx.
			Status(returnedHttpResponse.Status).
			Send(returnedHttpResponse.Body)
	}
}
