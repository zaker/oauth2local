package cmd

import (
	"log"
	"os"

	"github.com/pkg/browser"

	"github.com/equinor/oauth2local/ipc"
	"github.com/equinor/oauth2local/oauth2"
	"github.com/equinor/oauth2local/storage"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
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
		jww.ERROR.Println("No config file loaded")
		os.Exit(1)
	}

	jww.INFO.Println("Using config file:", viper.ConfigFileUsed())
	if ipc.HasSovereign() {
		jww.INFO.Println("A server is already running")
		return
	}

	oauthHandler, err := oauth2.NewAdalHandler(
		oauth2.Oauth2Settings{
			AuthServer:   viper.GetString("Authserver"),
			TenantID:     viper.GetString("TenantID"),
			ClientID:     viper.GetString("ClientID"),
			ClientSecret: viper.GetString("ClientSecret"),
		},
		storage.Memory(),
		viper.GetString("CustomScheme"))
	if err != nil {
		log.Printf("Error with oauth client: %v", err)
		return
	}

	jww.INFO.Println("starting browser...")

	lpu, err := oauthHandler.LoginProviderURL()
	if err != nil {
		log.Printf("Login provider url isn't an url: %v", err)
		return
	}

	browser.OpenURL(lpu)
	s := ipc.NewServer(oauthHandler)

	jww.ERROR.Println("Cannot serve:", s.Serve())
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
