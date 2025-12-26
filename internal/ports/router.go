package ports

import (
	"context"

	"github.com/YuruDeveloper/codey/internal/types"
	"github.com/google/uuid"
)

type AppRouter interface {
	GetProvider(uuid.UUID) Provider
	Config() AppConfig
	SendMessage(ctx context.Context,parts []types.Part) types.Message
	RegisterTool(tool types.Tool)
	ExecuteTool(ctx context.Context,toolUse types.ToolUsePart) types.ToolResultPart
}