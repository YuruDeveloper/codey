package anthropicAuth

import (
	"encoding/json"

	"github.com/YuruDeveloper/codey/internal/auth"
	"github.com/YuruDeveloper/codey/internal/config"
	appError "github.com/YuruDeveloper/codey/internal/error"
)

var _ auth.Auth = (*ApiKeyAuth)(nil)

func NewApiKeyAuth(config config.AppConfig) (*ApiKeyAuth , error) {
	data := config.GetProviderAuth(name)
	if data == nil {
		return nil , appError.NewValidError(appError.FailLoadJsonData,"Fail Load anthropic api key data")
	}
	var auth AuthData
	err := json.Unmarshal(data, &auth)
	if err != nil {
		return nil , appError.NewError(appError.JsonUnMarshalError,err)
	}
	return &ApiKeyAuth{
		key: auth.Key,
	} , nil
}

type ApiKeyAuth struct {
	key string
}

func (instance *ApiKeyAuth) SetApiKey(key string) {
	instance.key = key
}

func (instance *ApiKeyAuth) Key() string {
	return instance.key
}

func (instance *ApiKeyAuth) Save(config config.AppConfig) error {
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
