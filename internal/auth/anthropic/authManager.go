package anthropicAuth

import (
	"context"
	"encoding/json"

	appError "github.com/YuruDeveloper/codey/internal/error"
	"github.com/YuruDeveloper/codey/internal/ports"
)

var _ ports.AuthManager = (*AuthManager)(nil)

var auths []string = []string { "Apikey", "OAuth" }

type AuthManager struct {}

func (AuthManager) LoadAuth(config ports.AppConfig) (ports.Auth, error) {
	data := config.GetProviderAuth(name)
	var auth AuthData
	if data == nil {
		return nil , appError.NewValidError(appError.FailLoadJsonData,"Fail Make anthropic auth")
	}
	err :=  json.Unmarshal(data, &auth)
	if err != nil {
		return nil , appError.NewError(appError.JsonUnMarshalError,err)
	}
	switch auth.Type {
		case ApiKey:
			return NewApiKeyAuth(auth) , nil
		case OAuth:
			return NewOAuthAuth(auth) , nil
		default:
			return nil , appError.NewValidError(0,"UnknownAuthType")
	} 
}

func (AuthManager) SupportedAuths() []string {
	return auths
}

func (AuthManager) Authenticate(index int,authContext ports.AuthContext,ctx context.Context) (ports.Auth , error) {
	switch index {
		case 0:
			key := authContext.GetUserInput("ApiKey:")
			auth := &ApiKeyAuth { key: key }
			err := auth.Save(authContext.GetConfig())
			return auth , err
		case 1:
			verifier := Authorize(Max,authContext.ShowMessage)
			auth := &OAuthAuth{}
			err := auth.ExchangeToken(context.Background(),authContext.GetUserInput("Code: "),verifier)
			if err != nil {
				return nil , err
			}
			err = auth.Save(authContext.GetConfig())
			return auth , err
		default:
			return nil , appError.NewValidError(appError.UnexpectedAuthIndex,"anthropic auth index error")
	}
}

