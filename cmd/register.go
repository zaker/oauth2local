package cmd
import (
	"os"

	"github.com/spf13/viper"

	"github.com/equinor/oauth2local/register"

	"github.com/spf13/cobra"
)

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register app ass url handler for custom url",
	Long:  `Register app ass url handler for custom url`,
	Run: func(cmd *cobra.Command, args []string) {
		register.RegMe(viper.GetString("CustomScheme"), os.Args[0])
	},
}

func init() {
	rootCmd.AddCommand(registerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// registerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// registerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
