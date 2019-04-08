package cmd
import (
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
)

var cfgFile string
var verbose bool

var rootCmd = &cobra.Command{
	Use:   "oauth2local",
	Short: "oauth2local is providing oauth2 authenticated tokens to local processes",
	Long:  "oauth2local is providing oauth2 authenticated tokens to local processes",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		jww.ERROR.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.oauth2local.yaml)")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "log to console to console")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initDefaultConfig() {
	viper.SetDefault("Authserver", "https://login.microsoftonline.com")
	viper.SetDefault("TenantID", "common")
	viper.SetDefault("ClientID", "")
	viper.SetDefault("ClientSecret", "")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			jww.ERROR.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".oauth2local" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".oauth2local")
	}

	if verbose {
		jww.SetLogThreshold(jww.LevelTrace)
		jww.SetStdoutThreshold(jww.LevelTrace)
	}
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		jww.ERROR.Println(err)
	}

}
