package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/pkg/browser"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zaker/oauth2local/config"
	"github.com/zaker/oauth2local/ipc"
	"github.com/zaker/oauth2local/oauth2"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "serve a local auth provider",
	Long:  `serve a local auth provider.`,
	Run:   runServe,
}

func runServe(cmd *cobra.Command, args []string) {

	if viper.ConfigFileUsed() == "" {
		slog.Error("No config file loaded")
		os.Exit(1)
	}

	if ipc.HasSovereign() {
		slog.Info("A server is already running")
		os.Exit(9)
	}

	opts := []oauth2.Option{oauth2.WithOauth2Settings(
		oauth2.Oauth2Settings{
			AuthServer:   config.AuthServer(),
			TenantID:     config.TenantID(),
			ClientID:     config.ClientID(),
			ClientSecret: config.ClientSecret(),
			ResourceID:   config.ResourceID(),
		})}
	if config.IgnoreStateCheck() {
		opts = append(opts, oauth2.WithState("none"))
	}

	var port = config.LocalHttpServerPort()
	if port > 0 && port < 65536 {
		slog.Debug("Using local http server to receive callback")
		opts = append(opts, oauth2.WithLocalhostHttpServer(port))

	}
	var oauthHandler oauth2.Handler
	var err error
	switch config.AuthServerType() {
	case "adal":
		oauthHandler, err = oauth2.NewMsal(opts...)
	case "msal":
		oauthHandler, err = oauth2.NewAdal(opts...)
	default:
		err = fmt.Errorf("no authserver selected")
	}

	if err != nil {
		slog.Error("error with oauth client: %v", "inner", err)
		os.Exit(1)
	}

	slog.Info("starting browser...")

	lpu, err := oauthHandler.LoginProviderURL()
	if err != nil {
		slog.Error("login provider url isn't an url", "error", err)
		os.Exit(1)
	}

	if !testing {
		err = browser.OpenURL(lpu)
		if err != nil {
			slog.Error("error", "inner", err)
			os.Exit(1)
		}
	}
	s := ipc.NewServer(oauthHandler)

	slog.Error("Cannot serve", "serverUrl", s.Serve())
	os.Exit(1)
}

func init() {
	rootCmd.AddCommand(serveCmd)

}
