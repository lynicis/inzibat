package main

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"

	"github.com/Lynicis/inzibat/client"
	"github.com/Lynicis/inzibat/config"
	"github.com/Lynicis/inzibat/router"
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
		configFileName = path.Clean(
			fmt.Sprintf("%s.json", configFileName),
		)
	}
	configFilePath := path.Join(workingDirectory, configFileName)

	var configReader config.Reader
	configReader, err = config.NewReader(extensionOfFilePath)
	if err != nil {
		panic(err)
	}

	configLoader := &config.Loader{
		ConfigReader: configReader,
	}

	var cfg *config.Cfg
	cfg, err = configLoader.LoadConfig(configFilePath)
	if err != nil {
		panic(err)
	}

	fiberApp := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		JSONDecoder:           json.Unmarshal,
		JSONEncoder:           json.Marshal,
	})

	httpClient := &client.HttpClient{
		FasthttpClient: &fasthttp.Client{
			ReadTimeout:                   10 * time.Second,
			WriteTimeout:                  10 * time.Second,
			MaxIdleConnDuration:           10 * time.Second,
			NoDefaultUserAgentHeader:      true,
			DisableHeaderNamesNormalizing: true,
			DisablePathNormalizing:        true,
		},
	}

	mockHandler := router.NewMockHandler(cfg.Routes)
	clientHandler := router.NewClientHandler(httpClient, cfg.Routes)
	mainRouter := router.NewMainRouter(
		cfg,
		fiberApp,
		mockHandler,
		clientHandler,
	)
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

	serverPort := fmt.Sprintf(":%d", cfg.ServerPort)
	err = fiberApp.Listen(serverPort)
	if err != nil {
		panic(err)
	}
}
