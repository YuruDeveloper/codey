package types

type StreamEvent struct {
	Type EventType

	Text string

	ToolUseID string
	ToolInput string

	Error      error
	StopReason string
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