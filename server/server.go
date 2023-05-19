package server

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

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
	shutdownChannel := make(chan os.Signal, 1)
	signal.Notify(shutdownChannel, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-shutdownChannel
		_ = s.fiber.Shutdown()
	}()

	serverAddress := fmt.Sprintf(":%s", s.port)
	return s.fiber.Listen(serverAddress)
}

func (s *server) Shutdown() error {
	return s.fiber.Shutdown()
}

func (s *server) GetFiberInstance() *fiber.App {
	return s.fiber
}
