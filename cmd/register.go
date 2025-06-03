package cmd

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/zaker/oauth2local/config"
	"github.com/zaker/oauth2local/register"
)

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register app as url handler for custom url",
	Long:  `Register app as url handler for custom url`,
	Run: func(cmd *cobra.Command, args []string) {
		err := register.RegMe(config.CustomScheme(), os.Args[0])
		if err != nil {
			slog.Error("error", "inner", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(registerCmd)
}
