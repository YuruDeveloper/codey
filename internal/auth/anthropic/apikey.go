package anthropicAuth

import (
	"context"
	"encoding/json"

	"github.com/YuruDeveloper/codey/internal/auth"
	"github.com/YuruDeveloper/codey/internal/config"
)

var _ auth.Auth = (*ApiKeyAuth)(nil)

func NewApiKeyAuth(config *config.Config) *ApiKeyAuth {
	data := config.GetProviderAuth(name)
	var auth AuthData
	json.Unmarshal(data, &auth)
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

func (*ApiKeyAuth) Update(ctx context.Context) {}

func (instance *ApiKeyAuth) Key() string {
	return instance.key
}

func (instance *ApiKeyAuth) Save(config *config.Config) {
	data, _ := json.Marshal(AuthData{
		Type: ApiKey,
		Key:  instance.key,
	})
	config.SetProviderAuth(name, data)
	config.Save()
}
