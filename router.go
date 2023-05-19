package main

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Router interface {
	CreateRoutes()
	HandleClientMethod(routeConfig *Route) func(ctx *fiber.Ctx) error
}

type router struct {
	config *Config
	app    *fiber.App
	client Client
}

func NewRouter(config *Config, app *fiber.App, client Client) Router {
	return &router{
		config: config,
		app:    app,
		client: client,
	}
}

func (r *router) CreateRoutes() {
	for _, route := range r.config.Routes {
		r.app.Add(route.Method, route.Path, r.HandleClientMethod(&route))
	}
}

func (r *router) HandleClientMethod(routeConfig *Route) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		var (
			err                error
			HostName           = routeConfig.RequestTo.Host
			Path               = routeConfig.RequestTo.Path
			Method             = routeConfig.RequestTo.Method
			RequestToUrl       = fmt.Sprintf("%s%s", HostName, Path)
			RequestWithHeaders = routeConfig.RequestTo.Header
			RequestWithBody    = routeConfig.RequestTo.Body
		)

		cloneOfClientStruct := r.client.GetCloneOfStruct()
		methodName := cases.Title(language.Und).String(strings.ToLower(Method))
		method := reflect.ValueOf(&client{fasthttp: cloneOfClientStruct.fasthttp}).MethodByName(methodName)

		var params []reflect.Value
		if Method == http.MethodGet {
			params = []reflect.Value{
				reflect.ValueOf(RequestToUrl),
				reflect.ValueOf(RequestWithHeaders),
			}
		} else {
			params = []reflect.Value{
				reflect.ValueOf(RequestToUrl),
				reflect.ValueOf(RequestWithHeaders),
				reflect.ValueOf(RequestWithBody),
			}
		}
		returnValues := method.Call(params)

		var ok bool
		var response *HttpResponse
		response, ok = returnValues[0].Interface().(*HttpResponse)
		if !ok {
			return errors.New("type casting error")
		}

		err, ok = returnValues[1].Interface().(error)
		if !ok {
			return errors.New("type casting error")
		}

		if err != nil {
			return err
		}

		return ctx.Status(response.Status).Send(response.Body)
	}
}
