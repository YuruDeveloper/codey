package types

import "encoding/json"

type PartType int

const (
	PartTypeText PartType = PartType(iota)
	PartTypeImage
	PartTypeTool
	PartTypeToolResult
	PartTypeThink
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
	Input json.RawMessage
}

func (ToolUsePart) GetType() PartType { return PartTypeTool }

type ToolResultPart struct {
	ToolUseID string
	Content   string
	IsError   bool
}

func (ToolResultPart) GetType() PartType { return PartTypeToolResult }

type ThinkPart struct {
	Thinking string
}

func (ThinkPart) GetType() PartType { return PartTypeThink }