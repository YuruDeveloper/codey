package auth

import (
	"context"

	"github.com/YuruDeveloper/codey/internal/config"
)

type Auth interface {
	Key() string
	Save(config config.AppConfig) error
}

type DynamicAuth interface {
	Update(ctx context.Context) error
}
