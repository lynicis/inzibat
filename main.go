package main

import (
	"os"

	"go.uber.org/zap"

	"github.com/lynicis/inzibat/cmd"
	_ "github.com/lynicis/inzibat/log"
	"github.com/lynicis/inzibat/server"
)

func main() {
	if len(os.Args) > 1 {
		cmd.Execute()
		return
	}

	if err := server.StartServer(""); err != nil {
		zap.L().Fatal("failed to start server", zap.Error(err))
	}
}
