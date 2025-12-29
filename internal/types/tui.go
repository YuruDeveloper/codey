package types

// TuiState TUI 상태
type TuiState int

const (
	TuiStateMain TuiState = iota
	TuiStateTalk
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

type UpdateToolStatus struct{
	Status ToolStatus 
	Info string
}