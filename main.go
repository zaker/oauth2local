package main
import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/zaker/oauth2local/register"

	"github.com/joho/godotenv"
	"github.com/pkg/browser"
)

const handleURLScheme = "loc-auth"

var redirectCallback = flag.String("r", "", "Handles redirect from azure ad")

func redirectURL() string {
	return handleURLScheme + "://callback"
}

func authorizeURL(tenant, clientID, state string) string {
	params := url.Values{}

	params.Set("redirect_uri", redirectURL())
	params.Set("client_id", clientID)
	params.Set("response_type", "code")
	params.Set("state", state)
	return fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/authorize?%s", tenant, params.Encode())
}

func tokenURL(tenant string) string {

	return fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/token", tenant)
}

func getToken(tokenURL, clientID, clientSecret, code string) (string, error) {

	params := url.Values{}

	params.Set("redirect_uri", redirectURL())
	params.Set("client_id", clientID)
	params.Set("client_secret", clientSecret)
	params.Set("grant_type", "authorization_code")
	params.Set("code", code)
	params.Set("resource", clientID)
	body := bytes.NewBufferString(params.Encode())
	cli := new(http.Client)
	resp, err := cli.Post(tokenURL, "application/x-www-form-urlencoded", body)
	if err != nil {
		return "", err
	}
	decoder := json.NewDecoder(resp.Body)
	var dat map[string]interface{}
	err = decoder.Decode(&dat)
	if err != nil {
		return "", err
	}

	accessToken := dat["access_token"].(string)
	return accessToken, nil
}

func main() {
	flag.Parse()

	if len(*redirectCallback) > 0 {
		// fmt.Print(*redirectCallback)
		u, err := url.Parse(*redirectCallback)
		if err != nil {
			log.Println("Error parsing url", err, *redirectCallback)
			time.Sleep(time.Second * 3)
			return
		}

		if u.Scheme != handleURLScheme {
			log.Println("We dont handle", u.Scheme)
			time.Sleep(time.Second * 3)
			return
		}
		params := u.Query()
		code := params.Get("code")
		accessToken, err := getToken(
			tokenURL(os.Getenv("AAD_TENANT_ID")),
			os.Getenv("AAD_CLIENT_ID"),
			os.Getenv("AAD_CLIENT_SECRET"),
			code)
		if err != nil {
			log.Println("Error parsing url", err)
			time.Sleep(time.Second * 3)
			return
		}
		fmt.Println(accessToken)
		time.Sleep(time.Second * 3)
		return
	}
	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		log.Fatal("Error loading .env file", err)
	}

	register.RegMe(handleURLScheme, os.Args[0])
	browser.OpenURL(authorizeURL(os.Getenv("AAD_TENANT_ID"), os.Getenv("AAD_CLIENT_ID"), "none"))
}
