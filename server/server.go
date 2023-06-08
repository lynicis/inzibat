package server

import (
	"fmt"
	"time"

	"github.com/Lynicis/inzibat/config"
	"github.com/gofiber/fiber/v2"
)

type Server interface {
	Start() error
	Shutdown() error
	GetFiberInstance() *fiber.App
}

type server struct {
	port  string
	fiber *fiber.App
}

func NewServer(config *config.Config) Server {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	return &server{
		port:  config.ServerPort,
		fiber: app,
	}
}

func (s *server) Start() error {
	serverAddress := fmt.Sprintf(":%s", s.port)
	return s.fiber.Listen(serverAddress)
}

func (s *server) Shutdown() error {
	return s.fiber.ShutdownWithTimeout(10 * time.Second)
}

func (s *server) GetFiberInstance() *fiber.App {
	return s.fiber
}
