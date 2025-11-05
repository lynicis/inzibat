package main

import (
	"fmt"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	json "github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"

	"inzibat/client/http"
	"inzibat/config"
	"inzibat/handler"
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
		panic(err)
	}

	extensionOfFilePath := path.Ext(configFileName)
	if extensionOfFilePath == "" {
		configFileName = path.Clean(fmt.Sprintf("%s.json", configFileName))
	}
	configFilePath := path.Join(workingDirectory, configFileName)

	var configReader config.ReaderStrategy
	configReader, err = config.NewReaderStrategy(extensionOfFilePath)
	if err != nil {
		panic(err)
	}

	configLoader := &config.Reader{
		ConfigReader: configReader,
	}

	var cfg *config.Cfg
	cfg, err = configLoader.Read(configFilePath)
	if err != nil {
		panic(err)
	}

	httpClient := http.NewHttpClient()

	endpointHandler := &handler.EndpointHandler{
		RouteConfig: &cfg.Routes,
	}
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
		if err = fiberApp.Listen(fmt.Sprintf(":%d", cfg.ServerPort)); err != nil {
			panic(err)
		}
	}()

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	<-signalChannel

	if err = fiberApp.ShutdownWithTimeout(5 * time.Second); err != nil {
		panic(err)
	}
}
