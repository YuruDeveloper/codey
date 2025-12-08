package anthropicAuth

import (
	"context"
	"encoding/json"
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
