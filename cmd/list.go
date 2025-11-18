package cmd

import (
	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/lynicis/inzibat/config"
	_ "github.com/lynicis/inzibat/log"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"list-routes", "ls", "l"},
	Short:   "List all routes",
	Run: func(cmd *cobra.Command, args []string) {
		v := validator.New()
		cfgLoader := config.NewLoader(v, true)

		cfg, err := cfgLoader.Read()
		if err != nil {
			zap.L().Fatal("failed to list routes", zap.Error(err))
		}

		zap.L().Info("Routes", zap.Any("routes", cfg.Routes))
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
