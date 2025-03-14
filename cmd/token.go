package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/zaker/oauth2local/ipc"
)

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Gets access token from the local server instance",
	Long:  `Gets access token from the local server instance"`,
	Run: func(cmd *cobra.Command, args []string) {

		cli, err := ipc.NewClient()
		if err != nil {
			jww.ERROR.Println(err)
			os.Exit(1)
		}

		a, err := cli.GetAccessToken()
		if err != nil {
			jww.ERROR.Println(err)

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
