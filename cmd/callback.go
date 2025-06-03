package cmd

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zaker/oauth2local/ipc"
)

var breakB bool

// callbackCmd represents the callback command
var callbackCmd = &cobra.Command{
	Use:   "callback",
	Short: "Send callback url to sovereign",
	Long:  `Send callback url to sovereign`,
	Run: func(cmd *cobra.Command, args []string) {
		err := sendCallback(args)
		if err != nil {

			slog.Error("send callback", "error", err.Error())
			_, _ = bufio.NewReader(os.Stdin).ReadBytes('\n')
			return
		}
		if breakB {
			slog.Info("sent calback %v", "Args", strings.Join(args, ","))
			_, _ = bufio.NewReader(os.Stdin).ReadBytes('\n')
		}
	},
}

func sendCallback(args []string) error {
	cli, err := ipc.NewClient()
	if err != nil {
		return err
	}

	if len(args) != 1 {
		return fmt.Errorf("only one arg supported")
	}

	err = cli.SendCallback(args[0])
	if err != nil {
		return err
	}

	return nil
}

func init() {
	rootCmd.AddCommand(callbackCmd)

	callbackCmd.Flags().BoolVarP(&breakB, "break", "b", false, "Break before exit")
}
