package ports

import "context"

type Auth interface {
	Key() string
	Save(config AppConfig) error
}

type DynamicAuth interface {
	Update(ctx context.Context) error
}

type AuthContext interface {
	ShowMessage(message string)
	GetUserInput(prompt string) string
	GetConfig() AppConfig
}

type AuthManager interface {
	LoadAuth(config AppConfig) (Auth , error)
	Authenticate(index int,authContext AuthContext,ctx context.Context) (Auth , error)
	SupportedAuths() []string
}