package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zaker/oauth2local/ipc"
)

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Gets access token from the local server instance",
	Long:  `Gets access token from the local server instance"`,
	Run: func(cmd *cobra.Command, args []string) {

		cli, err := ipc.NewClient()
		if err != nil {
			slog.Error("error", "inner", err)
			os.Exit(1)
		}

		a, err := cli.GetAccessToken()
		if err != nil {
			slog.Error("error", "inner", err)

			s := err.Error()

			if strings.Contains(s, "code = Unavailable") {
				os.Exit(8)
			} else {
				os.Exit(1)
			}
		}

		fmt.Println(a)
	},
}

func init() {
	rootCmd.AddCommand(tokenCmd)
}
