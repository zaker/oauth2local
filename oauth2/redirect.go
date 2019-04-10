package oauth2

import "net/url"

type RedirectParams struct {
	scheme string
	code   string
	state  string
}

func DecodeRedirect(u *url.URL) *RedirectParams {

	params := u.Query()
	return &RedirectParams{
		scheme: u.Scheme,
		code:   params.Get("code"),
		state:  params.Get("state"),
	}

}
