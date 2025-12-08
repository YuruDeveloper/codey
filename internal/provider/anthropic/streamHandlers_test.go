package anthropic

import (
	"errors"
	"testing"
	"time"

	"github.com/YuruDeveloper/codey/internal/types"
	"github.com/anthropics/anthropic-sdk-go"
)

// Helper function to assert no events are received on a channel
func assertNoEvent(t *testing.T, eventChan chan types.StreamEvent) {
	t.Helper()
	select {
	case unexpectedEvent := <-eventChan:
		t.Fatalf("Received unexpected event: %+v", unexpectedEvent)
	case <-time.After(10 * time.Millisecond):
		// Success, no event received
	}
}

func Test_errorHandler(t *testing.T) {
	eventChan := make(chan types.StreamEvent, 1)
	testError := errors.New("test error")

	errorHandler(eventChan, testError)

	event := <-eventChan
	if event.Type != types.EventTypeError {
		t.Errorf("Expected event type %v, got %v", types.EventTypeError, event.Type)
	}
	if event.Error != testError {
		t.Errorf("Expected error %v, got %v", testError, event.Error)
	}
}

func Test_blockStartHandler(t *testing.T) {
	t.Run("ToolStart event", func(t *testing.T) {
		eventChan := make(chan types.StreamEvent, 1)
		toolList := make(map[int64]string)
		contentBlock := anthropic.ContentBlockStartEventContentBlockUnion{
			Type: "tool_use", // This is the value of the 'ToolStart' constant
			ID:   "tool-123",
		}

		blockStartHandler(eventChan, 0, contentBlock, toolList)

		event := <-eventChan
		if event.Type != types.EventTypeToolUseStart {
			t.Errorf("Expected event type %v, got %v", types.EventTypeToolUseStart, event.Type)
		}
		if event.ToolUseID != "tool-123" {
			t.Errorf("Expected ToolUseID 'tool-123', got '%s'", event.ToolUseID)
		}
		if toolList[0] != "tool-123" {
			t.Errorf("Expected toolList to contain the tool ID, but it doesn't")
		}
	})

	t.Run("Non-ToolStart event", func(t *testing.T) {
		eventChan := make(chan types.StreamEvent, 1)
		toolList := make(map[int64]string)
		contentBlock := anthropic.ContentBlockStartEventContentBlockUnion{
			Type: "text",
		}

		blockStartHandler(eventChan, 0, contentBlock, toolList)
		assertNoEvent(t, eventChan)
	})
}

func Test_blockDeltaHandler(t *testing.T) {
	t.Run("MessageDelta event", func(t *testing.T) {
		eventChan := make(chan types.StreamEvent, 1)
		toolList := make(map[int64]string)
		delta := anthropic.MessageStreamEventUnionDelta{
			Type: "text_delta", // This is the value of the 'MessageDelta' constant
			Text: "hello",
		}

		blockDetlaHandler(eventChan, 0, delta, toolList)

		event := <-eventChan
		if event.Type != types.EventTypeTextDelta {
			t.Errorf("Expected event type %v, got %v", types.EventTypeTextDelta, event.Type)
		}
		if event.Text != "hello" {
			t.Errorf("Expected text 'hello', got '%s'", event.Text)
		}
	})

	t.Run("ToolInputDelta event", func(t *testing.T) {
		eventChan := make(chan types.StreamEvent, 1)
		toolList := map[int64]string{0: "tool-123"}
		delta := anthropic.MessageStreamEventUnionDelta{
			Type:        "input_json_delta", // This is the value of the 'ToolInputDelta' constant
			PartialJSON: `{"arg":`,
		}

		blockDetlaHandler(eventChan, 0, delta, toolList)

		event := <-eventChan
		if event.Type != types.EventTypeToolUseInput {
			t.Errorf("Expected event type %v, got %v", types.EventTypeToolUseInput, event.Type)
		}
		if event.ToolUseID != "tool-123" {
			t.Errorf("Expected ToolUseID 'tool-123', got '%s'", event.ToolUseID)
		}
		if event.ToolInput != `{"arg":` {
			t.Errorf("Expected ToolInput to be `{\"arg\":`, got `%s`", event.ToolInput)
		}
	})
}

func Test_blockStopHandler(t *testing.T) {
	t.Run("Tool exists", func(t *testing.T) {
		eventChan := make(chan types.StreamEvent, 1)
		toolList := map[int64]string{0: "tool-123"}

		blockStopHandler(eventChan, 0, toolList)

		event := <-eventChan
		if event.Type != types.EventTypeToolUseEnd {
			t.Errorf("Expected event type %v, got %v", types.EventTypeToolUseEnd, event.Type)
		}
		if event.ToolUseID != "tool-123" {
			t.Errorf("Expected ToolUseID 'tool-123', got '%s'", event.ToolUseID)
		}
		if _, exists := toolList[0]; exists {
			t.Error("Expected tool ID to be deleted from toolList")
		}
	})

	t.Run("Tool does not exist", func(t *testing.T) {
		eventChan := make(chan types.StreamEvent, 1)
		toolList := map[int64]string{1: "tool-456"}

		blockStopHandler(eventChan, 0, toolList)

		assertNoEvent(t, eventChan)
	})
}

func Test_messageEndHandler(t *testing.T) {
	eventChan := make(chan types.StreamEvent, 1)
	stopReason := "end_turn"

	messageEndHandler(eventChan, stopReason)

	event := <-eventChan
	if event.Type != types.EventTypeMessageEnd {
		t.Errorf("Expected event type %v, got %v", types.EventTypeMessageEnd, event.Type)
	}
	if event.StopReason != "end_turn" {
		t.Errorf("Expected StopReason 'end_turn', got '%s'", event.StopReason)
	}
}
