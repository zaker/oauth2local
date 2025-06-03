package cmd

import (
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(defaultsCmd)
}

var defaultsCmd = &cobra.Command{
	Use:   "defaults",
	Short: "Writes default config values to config file",
	Long:  `Writes default config values to config file, specified with --config`,
	Run: func(cmd *cobra.Command, args []string) {
		err := viper.WriteConfig()
		if err != nil {
			slog.Error("write config", "err", err)
		}
	},
}
