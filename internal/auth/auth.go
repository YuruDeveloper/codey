package auth

import (
	"context"

	"github.com/YuruDeveloper/codey/internal/config"
)

type Auth interface {
	Update(ctx context.Context)
	Key() string
	Save(config *config.Config)
}
