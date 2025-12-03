package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/lynicis/inzibat/config"
	_ "github.com/lynicis/inzibat/log"
)

var (
	listConfigFile     string
	listIsGlobalConfig bool
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"list-routes", "ls", "l"},
	Short:   "List all routes",
	Long: `List all routes from the configuration file.

The routes will be read from (in order of precedence):
  1. The file specified by the --config flag
  2. The file specified by the INZIBAT_CONFIG_FILE environment variable
  3. inzibat.json in the current working directory
  4. ~/.inzibat.config.json if --global flag is used`,
	Run: func(cmd *cobra.Command, args []string) {
		cfgLoader := config.NewLoader(nil, listIsGlobalConfig, listConfigFile)

		cfg, err := cfgLoader.Read()
		if err != nil {
			zap.L().Fatal("failed to list routes", zap.Error(err))
		}

		t := table.New().
			Border(lipgloss.NormalBorder()).
			Headers("METHOD", "PATH", "TYPE").
			Rows(cfg.ConvertRoutesTuiTable()...)

		fmt.Println(t)
	},
}

func init() {
	listCmd.Flags().StringVarP(
		&listConfigFile,
		"config",
		"c",
		"",
		"Path to the configuration file",
	)
	listCmd.Flags().BoolVarP(
		&listIsGlobalConfig,
		"global",
		"g",
		false,
		"Use the global config file (~/.inzibat.config.json)",
	)
	rootCmd.AddCommand(listCmd)
}
