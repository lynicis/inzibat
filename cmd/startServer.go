package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/lynicis/inzibat/server"
)

var (
	configFile      string
	isGlobalConfig  bool
	recordEnabled   bool
	startServerFunc = server.StartServer
)

func addServerFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(
		&configFile,
		"config",
		"c",
		"",
		"Path to the configuration file",
	)
	cmd.Flags().BoolVarP(
		&isGlobalConfig,
		"global",
		"g",
		false,
		"Use the global config file (~/.inzibat.config.json)",
	)
	cmd.Flags().BoolVarP(
		&recordEnabled,
		"record",
		"r",
		false,
		"Enable request recording to capture incoming HTTP traffic",
	)
}

func runStartServer(cmd *cobra.Command, args []string) error {
	if err := startServerFunc(configFile, isGlobalConfig, recordEnabled); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

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
	RunE: runStartServer,
}

func init() {
	addServerFlags(startServerCmd)
	addServerFlags(rootCmd)

	rootCmd.AddCommand(startServerCmd)
}
