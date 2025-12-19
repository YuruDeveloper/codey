package ports

import (
	"context"

	"github.com/YuruDeveloper/codey/internal/types"
)

type AppRouter interface {
	Provider() Provider
	Engine() AppNodeEngine
	Config() AppConfig
	SendMessage(ctx context.Context,parts []types.Part) types.Message
	RegisterTool(tool types.Tool)
	ExecuteTool(ctx context.Context,toolUse types.ToolUsePart) types.ToolResultPart
}