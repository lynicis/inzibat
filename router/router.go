package router

import (
	"errors"
	"net/http"
	"net/url"
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
			err                     error
			ConfigRequestToHostName = routeConfig.RequestTo.Host
			ConfigRequestToPath     = routeConfig.RequestTo.Path
			ConfigRequestToMethod   = routeConfig.RequestTo.Method
		)

		var requestUrlBuilder *url.URL
		requestUrlBuilder, err = url.Parse(ConfigRequestToHostName)
		if err != nil {
			return errors.New(ErrorUrlParse)
		}
		requestUrlBuilder.Path = ConfigRequestToPath
		ConfigRequestToURL := requestUrlBuilder.String()

		cloneOfClientStruct := r.client.GetCloneOfStruct()
		methodName := cases.Title(language.Und).String(strings.ToLower(ConfigRequestToMethod))
		request := reflect.ValueOf(cloneOfClientStruct).MethodByName(methodName)

		var (
			RequestWithHeaders = routeConfig.RequestTo.Header
			RequestWithBody    = routeConfig.RequestTo.Body
			params             []reflect.Value
		)

		if ConfigRequestToMethod == http.MethodGet {
			params = []reflect.Value{
				reflect.ValueOf(ConfigRequestToURL),
				reflect.ValueOf(RequestWithHeaders),
			}
		} else {
			params = []reflect.Value{
				reflect.ValueOf(ConfigRequestToURL),
				reflect.ValueOf(RequestWithHeaders),
				reflect.ValueOf(RequestWithBody),
			}
		}
		returnValues := request.Call(params)

		var isOk bool
		var returnedHttpResponse *client.HttpResponse
		returnedHttpResponse, isOk = returnValues[0].Interface().(*client.HttpResponse)
		if !isOk {
			return errors.New(ErrorTypeCasting)
		}

		returnedErr := returnValues[1].Interface()
		if returnedErr != nil {
			return errors.New(ErrorTypeCasting)
		}

		return ctx.Status(returnedHttpResponse.Status).Send(returnedHttpResponse.Body)
	}
}
