package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "inzibat",
	Short: "A lightweight HTTP mock server for microservices testing and development",
	Long: `Inzibat (from Turkish, meaning "Military Police") is a small, fully-customizable 
mock service intended for use as a lightweight HTTP mock server.

It reads simple configuration files (JSON/TOML/YAML) and serves mock responses, 
allowing teams to simulate downstream services during development and integration testing.

Key features:
  - Config-driven (JSON, TOML, YAML) for easy scenario definition
  - Fast — built on top of Fiber (which uses fasthttp)
  - Simple, declarative API for defining routes and responses
  - No-code scenarios — implement complex mock behavior without writing server code`,
}

var exitFunc = os.Exit

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		exitFunc(1)
	}
}
