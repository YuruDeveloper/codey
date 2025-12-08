package types

type Role int  

const (
	UserRole = Role(iota) 
	AssistantRole 
	ToolRole 
	AlarmRole
)

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
	Data string
}

func (ImagePart) GetType() PartType { return PartTypeImage }

type ToolUsePart struct {
	ID string
	Input map[string]string
}

func (ToolUsePart) GetType() PartType { return PartTypeTool }

type ToolResultPart struct {
	ToolUseID string
	Content string
	IsError bool
}

func (ToolResultPart) GetType() PartType { return PartTypeToolResult }

type Message struct {
	Role Role
	Parts []Part
}

type Tool struct {
	Name string
	Description string
	InputSchema map[string]string
}

type EventType int 

const (
	EventTypeTextDelta = EventType(iota)
	EventTypeToolUseStart
	EventTypeToolUseInput
	EventTypeToolUseEnd
	EventTypeMessageEnd
	EventTypeError
)

type StreamEvent struct {
	Type EventType

	Text string

	ToolUseID string
	ToolInput string

	Error error
	StopReason string
}