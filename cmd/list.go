package cmd

import (
	"inzibat/config"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
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
			log.Fatalf("failed to list routes")
		}

		log.Print(cfg.Routes)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
