package main

import (
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
			zap.L().Fatal("failed to initialize global config", zap.Error(err))
		}
		cmd.Execute()
		return
	}

	if err := server.StartServer(""); err != nil {
		zap.L().Fatal("failed to start server", zap.Error(err))
	}
}
