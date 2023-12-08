package main

import (
	"fmt"
	"os"

	"github.com/Lynicis/inzibat/client"
	"github.com/Lynicis/inzibat/config"
	"github.com/Lynicis/inzibat/router"
	"github.com/Lynicis/inzibat/server"
)

func main() {
	var err error

	configFileName := os.Getenv(config.EnvironmentVariableConfigFileName)
	if configFileName == "" {
		configFileName = config.DefaultConfigFileName
	}

	var cfg *config.Config
	cfg, err = config.ReadConfig(configFileName)
	if err != nil {
		panic(err)
	}

	serverInstance := server.NewServer(cfg)
	fiberInstance := serverInstance.GetFiberInstance()

	clientInstance := client.NewClient()
	routerInstance := router.NewRouter(cfg, fiberInstance, clientInstance)
	routerInstance.CreateRoutes()

	fmt.Println("🫡 INZIBAT 🪖")
	fmt.Printf(
		"Open Routes: %d\n", len(cfg.Routes),
	)
	fmt.Printf(
		"Server Port: %s", cfg.ServerPort,
	)

	err = serverInstance.Start()
	if err != nil {
		panic(err)
	}
}
