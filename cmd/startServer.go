package cmd

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/lynicis/inzibat/server"
)

var (
	configFile      string
	isGlobalConfig  = true
	startServerFunc = server.StartServer
)

var startServerCmd = &cobra.Command{
	Use:     "start",
	Aliases: []string{"start-server", "server", "s"},
	Short:   "Start the Inzibat mock server",
	Long: `Start the Inzibat mock server using the configuration file.

The server will read the configuration from (in order of precedence):
  1. The file specified by the --config flag
  2. The file specified by the INZIBAT_CONFIG_FILE environment variable
  3. inzibat.json in the current working directory

The server will start listening on the port specified in the configuration
and serve the routes defined in the config file.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := startServerFunc(configFile, isGlobalConfig); err != nil {
			zap.L().Fatal("failed to start server", zap.Error(err))
		}
	},
}

func init() {
	startServerCmd.Flags().StringVarP(
		&configFile,
		"config",
		"c",
		"",
		"Path to the configuration file",
	)
	startServerCmd.Flags().BoolVarP(
		&isGlobalConfig,
		"global",
		"g",
		false,
		"Use the global config file (~/.inzibat.config.json)",
	)
	rootCmd.AddCommand(startServerCmd)
}
