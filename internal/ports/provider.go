package ports

import (
	"context"

	"github.com/YuruDeveloper/codey/internal/types"
)

type Provider interface {
	Send(ctx context.Context, params types.SendParams) <-chan types.StreamEvent
	Models() []string
	Model() string
	SetModel(index int)
}

type ClientProvider interface {
	Reconnect(key Auth) error
}