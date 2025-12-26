package types

// TuiState TUI 상태
type TuiState int

const (
	TuiStateMain TuiState = iota
	TuiStateTalk
	TuiStateSelect
	TuiStateModel
)