package anthropic

import (
	"github.com/YuruDeveloper/codey/internal/types"
	
	"github.com/anthropics/anthropic-sdk-go"
)

func errorHandler(eventChan chan types.StreamEvent, err error) {
	eventChan <- types.StreamEvent{
		Type:  types.EventTypeError,
		Error: err,
	}
}

func blockStartHandler(eventChan chan types.StreamEvent, index int64, contentBlock anthropic.ContentBlockStartEventContentBlockUnion, toolList map[int64]string) {
	if contentBlock.Type == ToolStart {
		toolList[index] = contentBlock.ID
		eventChan <- types.StreamEvent{
			Type:      types.EventTypeToolUseStart,
			ToolUseID: contentBlock.ID,
		}
	}
}

func blockDetlaHandler(eventChan chan types.StreamEvent, index int64, delta anthropic.MessageStreamEventUnionDelta, toolList map[int64]string) {
	switch delta.Type {
	case MessageDelta:
		eventChan <- types.StreamEvent{
			Type: types.EventTypeTextDelta,
			Text: delta.Text,
		}
	case ToolInputDelta:
		eventChan <- types.StreamEvent{
			Type:      types.EventTypeToolUseInput,
			ToolUseID: toolList[index],
			ToolInput: delta.PartialJSON,
		}
	}
}

func blockStopHandler(eventChan chan types.StreamEvent, index int64, toolList map[int64]string) {
	if id, exists := toolList[index]; exists {
		eventChan <- types.StreamEvent{
			Type:      types.EventTypeToolUseEnd,
			ToolUseID: id,
		}
		delete(toolList, index)
	}
}

func messageEndHandler(eventChan chan types.StreamEvent, stopReason string) {
	eventChan <- types.StreamEvent{
		Type:       types.EventTypeMessageEnd,
		StopReason: stopReason,
	}
}
