package main

import (
	"os"
)

func main() {
	workingDirectory, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	configFileName := os.Getenv(EnvironmentVariableConfigFileName)
	if configFileName == "" {
		configFileName = DefaultConfigFileName
	}

	config, err := ReadConfig(workingDirectory, configFileName)
	if err != nil {
		panic(err)
	}

	server := NewServer(config)
	app := server.GetFiberInstance()

	client := NewClient()
	router := NewRouter(config, app, client)
	router.CreateRoutes()

	config.Print()
	err = server.Start()
	if err != nil {
		panic(err)
	}
}
