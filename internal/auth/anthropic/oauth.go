package anthropicAuth

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/YuruDeveloper/codey/internal/auth"
	"github.com/YuruDeveloper/codey/internal/config"
	"golang.org/x/oauth2"
)

var _ auth.Auth = (*OAuthAuth)(nil)

const ClientID = "9d1c250a-e61b-44d9-88ed-5944d1962f5e"

type OAuthAuth struct {
	token        *oauth2.Token
	oauth2Config *oauth2.Config
}

func NewOAuthAuth(config *config.Config) *OAuthAuth {
	data := config.GetProviderAuth(name)
	var auth AuthData
	json.Unmarshal(data, &auth)

	token := &oauth2.Token{
		AccessToken:  auth.Access,
		RefreshToken: auth.Refresh,
		Expiry:       time.Unix(auth.Expires, 0),
	}

	oauth2Config := &oauth2.Config{
		ClientID: ClientID,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://console.anthropic.com/oauth/authorize",
			TokenURL: "https://console.anthropic.com/v1/oauth/token",
		},
		RedirectURL: "https://console.anthropic.com/oauth/code/callback",
		Scopes: []string{
			"org:create_api_key",
			"user:profile",
			"user:inference",
		},
	}

	return &OAuthAuth{
		token:        token,
		oauth2Config: oauth2Config,
	}
}

func (instance *OAuthAuth) Key() string {
	return instance.token.AccessToken
}

func (instance *OAuthAuth) Update(ctx context.Context) {
	if instance.token.Valid() {
		return
	}

	tokenSource := instance.oauth2Config.TokenSource(ctx, instance.token)
	newToken, err := tokenSource.Token()
	if err != nil {
		return
	}
	instance.token = newToken
}

func (instance *OAuthAuth) Save(config *config.Config) {
	data, _ := json.Marshal(AuthData{
		Type:    OAuth,
		Refresh: instance.token.RefreshToken,
		Access:  instance.token.AccessToken,
		Expires: instance.token.Expiry.Unix(),
	})
	config.SetProviderAuth(name, data)
	config.Save()
}


func (instance *OAuthAuth) ExchangeToken(ctx context.Context, code string, verifier string) {
	parts := strings.Split(code,"#")
	if len(parts) != 2 {
		return 
	}

	codeValue := parts[0]
	stateValue := parts[1]

	payload := map[string]string {
		"grant_type" : "authorization_code",
		"code" : codeValue,
		"state" : stateValue,
		"client_id" : ClientID,
		"redirect_uri":  "https://console.anthropic.com/oauth/code/callback",
		"code_verifier": verifier,
	}

	payloadBytes, err := json.Marshal(payload)
    if err != nil {
        return 
    }

	request , err := http.NewRequestWithContext(
		ctx,
		"POST",
		instance.oauth2Config.Endpoint.TokenURL,
        bytes.NewReader(payloadBytes),
	)
	if err != nil {
		return
	}
	 client := http.DefaultClient
	if c, ok := ctx.Value(oauth2.HTTPClient).(*http.Client); ok {
		client = c
	} 
	responce , err := client.Do(request)
	if err != nil {
		return
	}
	defer responce.Body.Close()

	if responce.StatusCode != http.StatusOK {
		return
	}

	var result struct {
		AccessToken  string `json:"access_token"`  
        RefreshToken string `json:"refresh_token"`  
        ExpiresIn    int64  `json:"expires_in"`    
        TokenType    string `json:"token_type"`     
	}

	if err := json.NewDecoder(responce.Body).Decode(&result) ; err != nil {
		return
	}

	instance.token.AccessToken = result.AccessToken
	instance.token.RefreshToken = result.RefreshToken
	instance.token.Expiry = time.Unix(result.ExpiresIn,0)
}