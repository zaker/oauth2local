package config

import (
	"fmt"
	"net/url"

	"github.com/spf13/viper"
)

type Config struct {
	authServer *url.URL
}

var cfg *Config

func SetDefaults() {
	viper.SetDefault("CustomScheme", "loc-auth")
}

func Load() error {
	cfg = new(Config)

	if !viper.GetBool("NO_AUTH") {
		a, err := parseURL(viper.GetString("AUTHSERVER"))
		if err != nil {
			return err
		}
		cfg.authServer = a
	}

	return nil
}

func parseURL(s string) (*url.URL, error) {

	if len(s) == 0 {
		return nil, fmt.Errorf("Url value empty")
	}
	u, err := url.Parse(s)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func AuthServer() string {
	return viper.GetString("Authserver")
}
func TenantID() string {
	return viper.GetString("TenantID")
}

func ClientID() string {
	return viper.GetString("ClientID")
}
func ClientSecret() string {
	return viper.GetString("ClientSecret")
}
func ResourceID() string {
	return viper.GetString("ResourceID")
}
func CustomScheme() string {
	return viper.GetString("CustomScheme")
}
