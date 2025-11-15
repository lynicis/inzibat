package main

import (
	"log"
	"os"

	"go.uber.org/zap"

	"inzibat/cmd"
	"inzibat/config"
	_ "inzibat/log"
	"inzibat/server"
)

func main() {
	if len(os.Args) > 1 {
		if err := config.InitGlobalConfig(); err != nil {
			log.Fatal("failed to initialize global config")
		}
		cmd.Execute()
		return
	}

	if err := server.StartServer(""); err != nil {
		zap.L().Fatal("failed to start server", zap.Error(err))
	}
}
