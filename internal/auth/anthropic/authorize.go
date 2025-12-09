package anthropicAuth

import (
	"net/url"

	"github.com/YuruDeveloper/codey/internal/browser"
)

type AuthMode int 

const (
	console = AuthMode(iota) 
	Max 
)

type ShowUrl func (url string)

func AuthorizeURL(mode AuthMode) (string ,*PKCE){
	pkce := GeneratePKCE()
	if pkce == nil {
		return "" , pkce
	}

	baseURL := "https://console.anthropic.com"
	if mode == Max {
		baseURL = "https://claude.ai"
	}

	authURL, _ := url.Parse(baseURL + "/oauth/authorize")
	query := authURL.Query()

	query.Set("client_id", ClientID)
	query.Set("response_type", "code")
	query.Set("redirect_uri", "https://console.anthropic.com/oauth/code/callback")
	query.Set("scope", "org:create_api_key user:profile user:inference")
	query.Set("code_challenge", pkce.Challenge)
	query.Set("code_challenge_method", "S256")
	query.Set("state", pkce.Verifier)
	query.Set("code", "true")
	authURL.RawQuery = query.Encode()

	return authURL.String(), pkce
}

func Authorize(mode AuthMode,show ShowUrl) *PKCE {
	authURL , PKCE := AuthorizeURL(mode)
	if PKCE == nil {
		return nil
	}
	if err := browser.Browser(authURL); err != nil {
		show(authURL)
	}

	return  PKCE
}

