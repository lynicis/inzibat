package router

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/Lynicis/inzibat/client"
	"github.com/Lynicis/inzibat/config"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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
	for index := range r.config.Routes {
		route := r.config.Routes[index]
		r.app.Add(route.Method, route.Path, r.HandleClientMethod(&route))
	}
}

func (r *router) HandleClientMethod(routeConfig *config.Route) func(ctx *fiber.Ctx) error {
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
		method := reflect.ValueOf(cloneOfClientStruct).MethodByName(methodName)

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

		var isOk bool
		var response *client.HttpResponse
		response, isOk = returnValues[0].Interface().(*client.HttpResponse)
		if !isOk {
			return errors.New(ErrorTypeCasting)
		}

		err, isOk = returnValues[1].Interface().(error)
		if !isOk {
			return errors.New(ErrorTypeCasting)
		}

		if err != nil {
			return err
		}

		return ctx.Status(response.Status).JSON(response.Body)
	}
}
