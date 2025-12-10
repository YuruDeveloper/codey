package types

type PartType int

const (
	PartTypeText PartType = PartType(iota)
	PartTypeImage
	PartTypeTool
	PartTypeToolResult
)

type Part interface {
	GetType() PartType
}

type TextPart struct {
	Text string
}

func (TextPart) GetType() PartType { return PartTypeText }

type ImagePart struct {
	MediaType string
	Data      string
}

func (ImagePart) GetType() PartType { return PartTypeImage }

type ToolUsePart struct {
	ID    string
	Input map[string]string
}

func (ToolUsePart) GetType() PartType { return PartTypeTool }

type ToolResultPart struct {
	ToolUseID string
	Content   string
	IsError   bool
}

func (ToolResultPart) GetType() PartType { return PartTypeToolResult }