package types

import (
	"github.com/google/uuid"
)

// TuiState TUI 상태
type TuiState int

const (
	TuiStateTalk TuiState = iota
	TuiStateSelect
	TuiStateModel
)

type ToolStatus int 

const (
	ToolPending ToolStatus = iota
	ToolSuccess
	ToolError
	ToolDefault
)

type StartTool struct {
	UUID uuid.UUID
	Name string
	Path string
}

type UpdateToolStatus struct{
	UUID uuid.UUID
	Status ToolStatus 
	Info string
}

type SubmitInput struct {
	Text string
}

type ResetCommand struct {}