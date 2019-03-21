package main
import (
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/joho/godotenv"
	"github.com/pkg/browser"
)

const handleURLScheme = "loc-auth://"

func redirectUri() string {
	return handleURLScheme + "callback"
}

func authorizeURL(tenant, clientID, state string) string {
	params := url.Values{}

	params.Set("redirect_uri", redirectUri())
	params.Set("client_id", clientID)
	params.Set("state", state)
	return fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/authorize?%s", tenant, params.Encode())
}

func tokenURL(tenant string) string {
	return fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/token", tenant)
}

func main() {
	log.Println(os.Args)
	log.Println(os.Environ())

	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		log.Fatal("Error loading .env file", err)
	}

	log.Println(os.Environ())

	browser.OpenURL(authorizeURL(os.Getenv("AAD_TENANT_ID"), os.Getenv("AAD_CLIENT_ID"), "none"))
}
