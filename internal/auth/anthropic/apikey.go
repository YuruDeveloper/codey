package anthropicAuth

import (
	"encoding/json"

	appError "github.com/YuruDeveloper/codey/internal/error"
	"github.com/YuruDeveloper/codey/internal/ports"
	"github.com/YuruDeveloper/codey/internal/types"
)

var _ ports.Auth = (*ApiKeyAuth)(nil)

func NewApiKeyAuth(auth AuthData) *ApiKeyAuth {
	return &ApiKeyAuth{
		key: auth.Key,
	} 
}

type ApiKeyAuth struct {
	key string
}

func (instance *ApiKeyAuth) SetApiKey(key string) {
	instance.key = key
}

func (instance *ApiKeyAuth) Key() (string , types.AuthType){
	return instance.key , types.ApiKey
}

func (instance *ApiKeyAuth) Save(config ports.AppConfig) error {
	data, err := json.Marshal(AuthData{
		Type: ApiKey,
		Key:  instance.key,
	})
	if err != nil {
		return appError.NewError(appError.JsonMarshalError,err)
	}
	config.SetProviderAuth(name, data)
	return config.Save()
}
