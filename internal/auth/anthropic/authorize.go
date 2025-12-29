package anthropicAuth

import (
	"github.com/YuruDeveloper/codey/internal/browser"
	"github.com/grokify/go-pkce"
	"golang.org/x/oauth2"
)

type AuthMode int 

const (
	console = AuthMode(iota) 
	Max 
)

type ShowUrl func (url string)

func AuthorizeURL(mode AuthMode) (string, string){
	verifier , err := pkce.NewCodeVerifier(-1)
	if err != nil {
		return "" , ""
	}

	baseURL := "https://console.anthropic.com"
	if mode == Max {
		baseURL = "https://claude.ai"
	}
	config := &oauth2.Config {
		ClientID: ClientID,
		ClientSecret: "",
		Endpoint: oauth2.Endpoint{
			AuthURL: baseURL + "/oauth/authorize",
			TokenURL: "https://console.anthropic.com/v1/oauth/token",
		},
		RedirectURL: "https://console.anthropic.com/oauth/code/callback",
		Scopes: []string{
          "org:create_api_key",
          "user:profile",
          "user:inference",
      	},
	}
	url := config.AuthCodeURL(verifier,oauth2.S256ChallengeOption(verifier))
	return url , verifier 
}

func Authorize(mode AuthMode,show ShowUrl)  string {
	authURL , verifier := AuthorizeURL(mode)
	if err := browser.Browser(authURL); err != nil {
		show(authURL)
	}

	return verifier
}

