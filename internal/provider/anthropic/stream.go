package anthropic

import (
	"context"

	"github.com/YuruDeveloper/codey/internal/types"

	"github.com/anthropics/anthropic-sdk-go"
)

type StreamInput struct {
	messages []anthropic.MessageParam
	tools    []anthropic.ToolUnionParam
	maxToken int64
	system   []anthropic.TextBlockParam
}

const (
	BlockStart     = "content_block_start"
	ToolStart      = "tool_use"
	BlockDelta     = "content_block_delta"
	MessageDelta   = "text_delta"
	ToolInputDelta = "input_json_delta"
	BlockStop      = "content_block_stop"
	MessageEnd     = "message_delta"
)

func initStream(parms types.SendParams) StreamInput {
	maxTokens := parms.MaxTokens

	if maxTokens == 0 {
		maxTokens = 4096
	}

	return StreamInput{
		messages: messageAdapter(parms.Messages),
		tools:    toolAdapter(parms.Tool),
		maxToken: int64(maxTokens),
		system: []anthropic.TextBlockParam{
			{Text: parms.SystemPrompt},
		},
	}
}

func stream(ctx context.Context, client *anthropic.Client, model anthropic.Model, eventChan chan types.StreamEvent, input StreamInput) {
	go func() {
		defer close(eventChan)

		stream := client.Messages.NewStreaming(ctx, anthropic.MessageNewParams{
			Model:     model,
			Messages:  input.messages,
			MaxTokens: input.maxToken,
			Tools:     input.tools,
			System:    input.system,
		})

		toolList := make(map[int64]string, 5)

		for stream.Next() {
			streamEvent := stream.Current()

			switch streamEvent.Type {
			case BlockStart:
				blockStartHandler(eventChan, streamEvent.Index, streamEvent.ContentBlock, toolList)
			case BlockDelta:
				blockDetlaHandler(eventChan, streamEvent.Index, streamEvent.Delta, toolList)
			case BlockStop:
				blockStopHandler(eventChan, streamEvent.Index, toolList)
			case MessageEnd:
				messageEndHandler(eventChan, string(streamEvent.Delta.StopReason))
			}
		}

		if err := stream.Err(); err != nil {
			errorHandler(eventChan, err)
		}
	}()
}
