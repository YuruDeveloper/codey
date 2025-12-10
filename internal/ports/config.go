package ports

import "encoding/json"

type AppConfig interface {
	Save() error
	Load() error
	GetProviderAuth(name string) json.RawMessage
	SetProviderAuth(name string, data json.RawMessage)
}