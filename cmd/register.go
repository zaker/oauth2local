package cmd

import (
	"os"

	"github.com/equinor/oauth2local/config"
	"github.com/equinor/oauth2local/register"
	"github.com/spf13/cobra"
)

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register app as url handler for custom url",
	Long:  `Register app as url handler for custom url`,
	Run: func(cmd *cobra.Command, args []string) {
		register.RegMe(config.CustomScheme(), os.Args[0])
	},
}

func init() {
	rootCmd.AddCommand(registerCmd)
}
