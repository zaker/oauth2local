package oauth2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/equinor/oauth2local/storage"
	"github.com/google/uuid"
	jww "github.com/spf13/jwalterweatherman"
)

type AdalHandler struct {
	net           *http.Client
	o2o           Oauth2Settings
	appRedirect   string
	handleScheme  string
	codeChallenge string
	store         storage.Storage
	jwtParser     *jwt.Parser
	ticker        *time.Ticker
	mut           *sync.Mutex
}

type Oauth2Settings struct {
	TenantID     string
	AuthServer   string
	ClientID     string
	ClientSecret string
}

const (
	authGrant    = "authorization_code"
	refreshGrant = "refresh_token"
)

func generateCodeChallenge() string {
	return uuid.New().String() + "-" + uuid.New().String()
}

func (o2o Oauth2Settings) Valid() bool {

	if o2o.AuthServer == "" {
		return false
	}
	strings.TrimRight(o2o.AuthServer, "/")

	if o2o.ClientID == "" {
		return false
	}

	if o2o.ClientSecret == "" {
		return false
	}
	return true
}
func NewAdalHandler(o2o Oauth2Settings, store storage.Storage, scheme string) (*AdalHandler, error) {

	if !o2o.Valid() {
		return nil, fmt.Errorf("Oauth2 Settings is not valid")
	}

	h := &AdalHandler{
		net:           new(http.Client),
		o2o:           o2o,
		appRedirect:   scheme + "://callback",
		handleScheme:  scheme,
		codeChallenge: generateCodeChallenge(),
		store:         store,
		jwtParser:     new(jwt.Parser),
		ticker:        time.NewTicker(1 * time.Minute),
		mut:           new(sync.Mutex)}

	go func() {
		for range h.ticker.C {
			err := h.renewTokens()
			if err != nil {
				jww.INFO.Println("Couldn't renew token", err)
			}
		}
	}()
	return h, nil
}

func (h *AdalHandler) renewTokens() error {

	a, err := h.store.GetToken(storage.AccessToken)
	if err != nil {
		return err
	}

	token, _, err := h.jwtParser.ParseUnverified(a, &jwt.StandardClaims{})
	//Reissue to authorize if old
	if claims, ok := token.Claims.(*jwt.StandardClaims); ok {

		tokenPeriod := claims.ExpiresAt - claims.IssuedAt
		currentPeriod := claims.ExpiresAt - time.Now().Unix()

		if currentPeriod > tokenPeriod/5 {
			jww.INFO.Println("Token still in grace period")
			return nil
		}

	}

	r, err := h.store.GetToken(storage.RefreshToken)
	if err != nil {
		return err
	}
	err = h.updateTokens(r, refreshGrant)
	if err != nil {
		return err
	}
	return nil
}

func (h *AdalHandler) tokenURL() string {

	return fmt.Sprintf("%s/oauth2/token", h.getAuthEndpoint())
}

func (h *AdalHandler) getAuthEndpoint() string {
	if h.o2o.TenantID == "" {
		return h.o2o.AuthServer
	}
	return fmt.Sprintf("%s/%s", h.o2o.AuthServer, h.o2o.TenantID)
}

func (h *AdalHandler) LoginProviderURL() (string, error) {
	u, err := url.Parse(fmt.Sprintf("%s/oauth2/authorize", h.getAuthEndpoint()))
	if err != nil {
		return "", err
	}
	jww.DEBUG.Println("LoginProvider at:", u)
	params := u.Query()

	params.Set("redirect_uri", h.appRedirect)
	params.Set("client_id", h.o2o.ClientID)
	params.Set("response_type", "code")
	params.Set("state", "none")
	params.Set("code_challenge", h.codeChallenge)

	u.RawQuery = params.Encode()
	return u.String(), nil

}

func CodeFromURL(callbackURL, scheme string) (string, error) {
	u, err := url.Parse(callbackURL)
	if err != nil {
		return "", err
	}

	if u.Scheme != scheme {
		return "", fmt.Errorf("App doesn't handle scheme: %s", u.Scheme)

	}
	params := u.Query()
	code := params.Get("code")

	return code, nil
}

func (h *AdalHandler) CodeFromURL(callbackURL string) (string, error) {
	return CodeFromURL(callbackURL, h.handleScheme)
}

func (h *AdalHandler) updateTokens(code, grant string) error {

	params := url.Values{}
	params.Set("client_id", h.o2o.ClientID)
	params.Set("client_secret", h.o2o.ClientSecret)
	params.Set("grant_type", grant)

	if grant == authGrant {
		params.Set("code_verifier", h.codeChallenge)
		params.Set("code", code)
		params.Set("redirect_uri", h.appRedirect)
	} else if grant == refreshGrant {
		params.Set("refresh_token", code)
	}
	params.Set("resource", h.o2o.ClientID)
	body := bytes.NewBufferString(params.Encode())

	tokenURL := h.tokenURL()
	resp, err := h.net.Post(tokenURL, "application/x-www-form-urlencoded", body)
	if err != nil {
		return fmt.Errorf("Error posting to token url %s: %s ", tokenURL, err)
	}
	if resp.StatusCode != 200 {

		return fmt.Errorf("Did not receive token: %v", resp.Status)

	}

	decoder := json.NewDecoder(resp.Body)
	var dat map[string]interface{}
	err = decoder.Decode(&dat)
	if err != nil {
		return err
	}

	if t, ok := dat["access_token"]; ok {

		err = h.store.SetToken(storage.AccessToken, t.(string))
		if err != nil {
			return err
		}
	}
	if t, ok := dat["id_token"]; ok {
		err = h.store.SetToken(storage.IDToken, t.(string))
		if err != nil {
			return err
		}
	}
	if t, ok := dat["refresh_token"]; ok {
		err = h.store.SetToken(storage.RefreshToken, t.(string))
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *AdalHandler) getValidAccessToken() (string, error) {
	a, err := h.store.GetToken(storage.AccessToken)
	if err != nil {
		return "", err
	}

	token, _, err := h.jwtParser.ParseUnverified(a, &jwt.StandardClaims{})
	//Reissue to authorize if old
	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return a, nil
	}

	return "", err
}

func (h *AdalHandler) GetAccessToken() (string, error) {

	a, err := h.store.GetToken(storage.AccessToken)
	if err != nil {
		return "", err
	}

	return a, nil
}
func (h *AdalHandler) UpdateFromRedirect(redirect *url.URL) error {

	// TODO: Validate state/nonce
	// Decode to authorize code
	h.mut.Lock()
	defer h.mut.Unlock()

	c, err := h.CodeFromURL(redirect.String())
	if err != nil {
		return err
	}

	err = h.updateTokens(c, authGrant)
	if err != nil {
		return err
	}

	return nil
}

func (h *AdalHandler) UpdateFromCode(code string) error {
	h.mut.Lock()
	defer h.mut.Unlock()
	return nil
}
