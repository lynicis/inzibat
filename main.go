package main

import (
	"os"

	"github.com/Lynicis/inzibat/client"
	"github.com/Lynicis/inzibat/config"
	"github.com/Lynicis/inzibat/router"
	"github.com/Lynicis/inzibat/server"
)

func main() {
	workingDirectory, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	configFileName := os.Getenv(config.EnvironmentVariableConfigFileName)
	if configFileName == "" {
		configFileName = config.DefaultConfigFileName
	}

	configInstance, err := config.ReadConfig(workingDirectory, configFileName)
	if err != nil {
		panic(err)
	}

	serverInstance := server.NewServer(configInstance)
	app := serverInstance.GetFiberInstance()

	clientInstance := client.NewClient()
	routerInstance := router.NewRouter(configInstance, app, clientInstance)
	routerInstance.CreateRoutes()

	configInstance.Print()
	err = serverInstance.Start()
	if err != nil {
		panic(err)
	}
}
