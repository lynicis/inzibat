package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"inzibat/client/http"
	"inzibat/config"
	"inzibat/handler"
	_ "inzibat/log"
	"inzibat/router"
)

func main() {
	var err error

	configFileName := os.Getenv(config.EnvironmentVariableConfigFileName)
	if configFileName == "" {
		configFileName = config.DefaultConfigFileName
	}

	var workingDirectory string
	workingDirectory, err = os.Getwd()
	if err != nil {
		zap.L().Fatal("failed to get current working directory path", zap.Error(err))
	}

	extensionOfFilePath := filepath.Ext(configFileName)
	if extensionOfFilePath == "" {
		configFileName = filepath.Clean(fmt.Sprintf("%s.json", configFileName))
	}
	configFilePath := filepath.Join(workingDirectory, configFileName)

	var configReader config.ReaderStrategy
	configReader, err = config.NewReaderStrategy(extensionOfFilePath)
	if err != nil {
		zap.L().Fatal("failed to create config reader strategy", zap.Error(err))
	}

	configLoader := &config.Reader{
		ConfigReader: configReader,
	}

	var cfg *config.Cfg
	cfg, err = configLoader.Read(configFilePath)
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
