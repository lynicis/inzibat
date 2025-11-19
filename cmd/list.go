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

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"list-routes", "ls", "l"},
	Short:   "List all routes",
	Run: func(cmd *cobra.Command, args []string) {
		cfgLoader := config.NewLoader(nil, true)

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
	rootCmd.AddCommand(listCmd)
}
