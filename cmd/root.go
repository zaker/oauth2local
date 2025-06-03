package cmd

import (
	"log/slog"
	"os"

	"github.com/zaker/oauth2local/config"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var cfgFile string
var verbose bool
var testing bool
var rootCmd = &cobra.Command{
	Use:   "oauth2local",
	Short: "oauth2local is providing oauth2 authenticated tokens to local processes",
	Long:  "oauth2local is providing oauth2 authenticated tokens to local processes",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		slog.Error("error", "inner", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.oauth2local.yaml)")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "log to console to console")
	rootCmd.PersistentFlags().BoolVar(&testing, "testing", false, "run in testing mode")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			slog.Error("error", "inner", err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".oauth2local")
	}

	if verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)

	}
	slog.Info("using config file", "filename", viper.ConfigFileUsed())

	config.SetDefaults()
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		slog.Error("error", "inner", err)
	}

}
