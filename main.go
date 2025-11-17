package main

import (
	"os"

	"go.uber.org/zap"

	"inzibat/cmd"
	_ "inzibat/log"
	"inzibat/server"
)

var (
	executeCmd    = cmd.Execute
	startServerFn = server.StartServer
)

func run(args []string) error {
	if len(args) > 1 {
		executeCmd()
		return nil
	}

	return startServerFn("")
}

func main() {
	if err := run(os.Args); err != nil {
		zap.L().Fatal("failed to start server", zap.Error(err))
	}
}
