package ports

import (
	"context"

	"github.com/YuruDeveloper/codey/internal/types"
	"github.com/google/uuid"
)

type Provider interface {
	Send(ctx context.Context, params types.SendParams) (types.Message, error)
	Models() []string
	Model() string
	SetModel(index int)
	GetUUID() uuid.UUID
}

type ClientProvider interface {
	Reconnect(key Auth) error
}