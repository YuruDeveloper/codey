package anthropicAuth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	appError "github.com/YuruDeveloper/codey/internal/error"
	"github.com/YuruDeveloper/codey/internal/ports"
	"github.com/YuruDeveloper/codey/internal/types"
	"golang.org/x/oauth2"
)

var _ ports.Auth = (*OAuthAuth)(nil)
var _ ports.DynamicAuth = (*OAuthAuth)(nil)

const ClientID = "9d1c250a-e61b-44d9-88ed-5944d1962f5e"

type OAuthAuth struct {
	token        *oauth2.Token
	oauth2Config *oauth2.Config
}

func NewOAuthAuth(auth AuthData) *OAuthAuth {
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

func (instance *OAuthAuth) Key() (string , types.AuthType){
	return instance.token.AccessToken , types.OAuth
}

func (instance *OAuthAuth) Update(ctx context.Context) error {
	if instance.token.Valid() {
		return nil
	}

	tokenSource := instance.oauth2Config.TokenSource(ctx, instance.token)
	newToken, err := tokenSource.Token()
	if err != nil {
		return appError.NewError(appError.FailRefreshToken,err)
	}
	instance.token = newToken
	return nil
}

func (instance *OAuthAuth) Save(config ports.AppConfig) error{
	data, err := json.Marshal(AuthData{
		Type:    OAuth,
		Refresh: instance.token.RefreshToken,
		Access:  instance.token.AccessToken,
		Expires: instance.token.Expiry.Unix(),
	})
	if err != nil {
		return appError.NewError(appError.JsonMarshalError,err)
	}
	config.SetProviderAuth(name, data)
	return config.Save()
}


func (instance *OAuthAuth) ExchangeToken(ctx context.Context, code string,verifier string) error {
	parts := strings.Split(code,"#")
	if len(parts) != 2 {
		return appError.NewValidError(appError.UnexpectedCode,"Fail Exchange anthropic token")
	}

	codeValue := parts[0]
	stateValue := parts[1]
	// Anthropic OAuth는 JSON 형식 사용 (OpenCode 참조)
	payload := map[string]string{
		"code":          codeValue,
		"state":         stateValue,
		"grant_type":    "authorization_code",
		"client_id":     ClientID,
		"redirect_uri":  "https://console.anthropic.com/oauth/code/callback",
		"code_verifier": verifier,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return appError.NewError(appError.JsonMarshalError, err)
	}

	request, err := http.NewRequestWithContext(
		ctx,
		"POST",
		instance.oauth2Config.Endpoint.TokenURL,
		strings.NewReader(string(payloadBytes)),
	)
	if err != nil {
		return appError.NewError(appError.FailMakeRequest, err)
	}

	request.Header.Set("Content-Type", "application/json")

	 client := http.DefaultClient
	if newClient, ok := ctx.Value(oauth2.HTTPClient).(*http.Client); ok {
		client = newClient
	} 
	response , err := client.Do(request)
	if err != nil {
		return appError.NewError(appError.FailGetResponse,err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		return appError.NewValidError(appError.HttpNotOK,fmt.Sprintf("Fail http anthropic token, status code: %s , body: %s",response.Status,string(body)))
	}

	var result struct {
		AccessToken  string `json:"access_token"`  
        RefreshToken string `json:"refresh_token"`  
        ExpiresIn    int64  `json:"expires_in"`    
        TokenType    string `json:"token_type"`     
	}

	if err := json.NewDecoder(response.Body).Decode(&result) ; err != nil {
		return appError.NewError(appError.FailDecodeHttpBody,err)
	}

	instance.token.AccessToken = result.AccessToken
	instance.token.RefreshToken = result.RefreshToken
	instance.token.Expiry = time.Now().Add(time.Duration(result.ExpiresIn) * time.Second)
	return nil
}