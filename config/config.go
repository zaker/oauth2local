package config
// import (
// 	"fmt"
// 	"os"

// 	"github.com/joho/godotenv"
// )

// // Config for oauth local service
// type Config struct {
// 	ClientID     string
// 	ClientSecret string
// 	TenantID     string
// 	// AppRedirect   string
// 	HandleScheme  string
// 	startAsServer bool
// 	RedirectURL   string
// }
// type ClientHandle int

// // Client handle scenarios
// const (
// 	None ClientHandle = iota
// 	Redirect
// 	AccessToken
// )

// // const handleURLScheme = "loc-auth"

// // func AredirectURL() string {
// // 	return handleURLScheme + "://callback"
// // }

// func Init(isClient bool) *Config {

// 	c := &Config{
// 		HandleScheme: handleURLScheme,
// 		// AppRedirect:   redirectURL(),
// 		startAsServer: !isClient}
// 	return c
// }

// func (cfg *Config) AsClient() bool {
// 	return !cfg.startAsServer
// }

// // LoadEnv loads environment variables into config
// func (cfg *Config) LoadEnv() error {
// 	err := godotenv.Load("C:\\bin\\.env")
// 	if err != nil {
// 		return err
// 	}
// 	tenant := os.Getenv("AAD_TENANT_ID")
// 	if len(tenant) == 0 {
// 		return fmt.Errorf("No tenant id => AAD_TENANT_ID")
// 	}
// 	cfg.TenantID = tenant

// 	clientID := os.Getenv("AAD_CLIENT_ID")
// 	if len(clientID) == 0 {
// 		return fmt.Errorf("No client id => AAD_CLIENT_ID")
// 	}
// 	cfg.ClientID = clientID

// 	clientSecret := os.Getenv("AAD_CLIENT_SECRET")
// 	if len(clientSecret) == 0 {
// 		return fmt.Errorf("No client secret => AAD_CLIENT_SECRET")
// 	}
// 	cfg.ClientSecret = clientSecret
// 	return nil
// }

// func (cfg *Config) ClientType() ClientHandle {

// 	if cfg.startAsServer {
// 		return None
// 	}
// 	if len(cfg.RedirectURL) > 0 {
// 		return Redirect
// 	}
// 	return None
// }
