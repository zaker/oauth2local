package cmd

import (
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Shows the config settings",
	Long:  `Shows the current congfig settings for the sovereign`,
	Run: func(cmd *cobra.Command, args []string) {
		jww.INFO.Println("config called")
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

}
