package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	validatorPkg "github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"inzibat/client/http"
	"inzibat/cmd"
	"inzibat/config"
	"inzibat/handler"
	_ "inzibat/log"
	"inzibat/router"
)

func main() {
	if len(os.Args) > 1 {
		if err := config.InitGlobalConfig(); err != nil {
			log.Fatal("failed to initialize global config")
		}
		cmd.Execute()
		return
	}

	var err error

	validator := validatorPkg.New()
	configLoader := config.NewLoader(validator, false)
	cfg, err := configLoader.Read()
	if err != nil {
		zap.L().Fatal("failed to read config", zap.Error(err))
	}

	endpointHandler := &handler.EndpointHandler{
		RouteConfig: &cfg.Routes,
	}
	httpClient := http.NewHttpClient()
	clientHandler := &handler.ClientHandler{
		Client:      httpClient,
		RouteConfig: &cfg.Routes,
	}

	fiberApp := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		JSONDecoder:           json.Unmarshal,
		JSONEncoder:           json.Marshal,
	})
	mainRouter := &router.MainRouter{
		Config:          cfg,
		FiberApp:        fiberApp,
		EndpointHandler: endpointHandler,
		ClientHandler:   clientHandler,
	}
	mainRouter.CreateRoutes()

	fmt.Print(
		"🫡 INZIBAT 🪖\n",
		fmt.Sprintf(
			"Open Routes: %d\n", len(cfg.Routes),
		),
		fmt.Sprintf(
			"Server Port: %d", cfg.ServerPort,
		),
	)

	go func() {
		if err = fiberApp.Listen(cfg.GetServerAddr()); err != nil {
			zap.L().Fatal("failed to start http server", zap.Error(err))
		}
	}()

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	<-signalChannel

	if err = fiberApp.ShutdownWithTimeout(5 * time.Second); err != nil {
		zap.L().Fatal("failed to shutdown gracefully", zap.Error(err))
	}
}
