package provider

import (
	"context"

	"github.com/YuruDeveloper/codey/internal/auth"
	"github.com/YuruDeveloper/codey/internal/types"
)

type Provider interface {
	Send(ctx context.Context, params SendParams) <-chan types.StreamEvent
	Models() []string
	Model() string
	SetModel(index int)
}

type ClientProvider interface {
	Reconnect(auth auth.Auth) error
}

type SendParams struct {
	Messages     []types.Message
	Tool         types.Tool
	SystemPrompt string
	MaxTokens    int
}
