package main

import (
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestNewRouter(t *testing.T) {
	r := NewRouter(nil, nil, nil)
	assert.Implements(t, (*Router)(nil), r)
}

func TestRouter_CreateRoutes(t *testing.T) {
	config := &Config{
		ServerPort: "3000",
		Routes: []Routes{
			{
				Method: http.MethodGet,
				Path:   "/",
				RequestTo: RequestTo{
					Host: "http://127.0.0.1:3001",
					Path: "/",
				},
			},
		},
	}
	app := fiber.New()
	r := NewRouter(config, app, nil)

	routers := r.CreateRoutes()

	assert.NotEmpty(t, routers)
}
