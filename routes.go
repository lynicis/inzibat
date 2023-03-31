package main

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type Router interface {
	CreateRoutes() []fiber.Router
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

func (r *router) CreateRoutes() []fiber.Router {
	var routers []fiber.Router

	for _, route := range r.config.Routes {
		r := r.app.Add(
			route.Method,
			route.Path,
			func(ctx *fiber.Ctx) error {
				return ctx.SendStatus(http.StatusOK)
			},
		)
		routers = append(routers, r)
	}

	return routers
}
