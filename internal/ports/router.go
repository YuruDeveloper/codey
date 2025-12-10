package ports

import "github.com/YuruDeveloper/codey/internal/types"

type AppRouter interface {
	//provider manage
	SetProvider(provider Provider,manager AuthManager)
	GetProivder() Provider
	// message send
	SendMessage(parts []types.Part) error
	// session manage
	SetSession(session AppSession)
	GetSession() AppSession  
	// tool manage
	RegisterTool(tool types.Tool)
	ExecuteTool(toolUse types.ToolUsePart) types.ToolResultPart
	// command manage
	HandleCommand(cmd string) error
	RegisterCommand(name string)
	// config manage
	GetConfig() AppConfig
	SaveConfig() error
}