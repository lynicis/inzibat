package main

import (
	"os"

	"github.com/Lynicis/inzibat/client"
	"github.com/Lynicis/inzibat/config"
	"github.com/Lynicis/inzibat/router"
	"github.com/Lynicis/inzibat/server"
)

func main() {
	var err error

	var workingDirectory string
	workingDirectory, err = os.Getwd()
	if err != nil {
		panic(err)
	}

	configFileName := os.Getenv(config.EnvironmentVariableConfigFileName)
	if configFileName == "" {
		configFileName = config.DefaultConfigFileName
	}

	var cfg *config.Config
	cfg, err = config.ReadConfig(workingDirectory, configFileName)
	if err != nil {
		panic(err)
	}

	serverInstance := server.NewServer(cfg)
	app := serverInstance.GetFiberInstance()

	clientInstance := client.NewClient()
	routerInstance := router.NewRouter(cfg, app, clientInstance)
	routerInstance.CreateRoutes()

	cfg.Print()
	err = serverInstance.Start()
	if err != nil {
		panic(err)
	}
}
