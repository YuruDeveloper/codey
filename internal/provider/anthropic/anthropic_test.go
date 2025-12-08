package anthropic

import (
	"reflect"
	"testing"

	"github.com/YuruDeveloper/codey/internal/provider"
	"github.com/YuruDeveloper/codey/internal/types"
	"github.com/anthropics/anthropic-sdk-go"
)

func TestInitStream_SimpleText(t *testing.T) {
	sendParams := provider.SendParams{
		Messages: []types.Message{
			{
				Role: types.UserRole,
				Parts: []types.Part{
					types.TextPart{Text: "Hello, world!"},
				},
			},
		},
		MaxTokens:    100,
		SystemPrompt: "You are a helpful assistant.",
		Tool:         types.Tool{}, // No tool for this test
	}

	expectedMessages := []anthropic.MessageParam{
		{
			Role: anthropic.MessageParamRoleUser,
			Content: []anthropic.ContentBlockParamUnion{
				anthropic.NewTextBlock("Hello, world!"),
			},
		},
	}

	expectedSystem := []anthropic.TextBlockParam{{Text: "You are a helpful assistant."}}

	expectedStreamInput := StreamInput{
		messages: expectedMessages,
		tools:    []anthropic.ToolUnionParam{}, // No tools expected
		maxToken: 100,
		system:   expectedSystem,
	}

	actualStreamInput := initStream(sendParams)

	if !reflect.DeepEqual(actualStreamInput.messages, expectedStreamInput.messages) {
		t.Errorf("initStream messages mismatch:\nGot: %+v\nWant: %+v", actualStreamInput.messages, expectedStreamInput.messages)
	}
	if !reflect.DeepEqual(actualStreamInput.tools, expectedStreamInput.tools) {
		t.Errorf("initStream tools mismatch:\nGot: %+v\nWant: %+v", actualStreamInput.tools, expectedStreamInput.tools)
	}
	if actualStreamInput.maxToken != expectedStreamInput.maxToken {
		t.Errorf("initStream maxToken mismatch:\nGot: %d\nWant: %d", actualStreamInput.maxToken, expectedStreamInput.maxToken)
	}
	if !reflect.DeepEqual(actualStreamInput.system, expectedStreamInput.system) {
		t.Errorf("initStream system mismatch:\nGot: %+v\nWant: %+v", actualStreamInput.system, expectedStreamInput.system)
	}
}

func TestInitStream_ImageMessage(t *testing.T) {
	sendParams := provider.SendParams{
		Messages: []types.Message{
			{
				Role: types.UserRole,
				Parts: []types.Part{
					types.ImagePart{
						MediaType: "image/jpeg",
						Data:      "fake-image-data",
					},
				},
			},
		},
	}

	expectedMessages := []anthropic.MessageParam{
		{
			Role: anthropic.MessageParamRoleUser,
			Content: []anthropic.ContentBlockParamUnion{
				anthropic.NewImageBlockBase64("image/jpeg", "fake-image-data"),
			},
		},
	}

	expectedStreamInput := StreamInput{
		messages: expectedMessages,
		tools:    []anthropic.ToolUnionParam{},
		maxToken: 4096, // default value
		system:   []anthropic.TextBlockParam{{Text: ""}},
	}

	actualStreamInput := initStream(sendParams)

	if !reflect.DeepEqual(actualStreamInput.messages, expectedStreamInput.messages) {
		t.Errorf("initStream messages mismatch:\nGot: %+v\nWant: %+v", actualStreamInput.messages, expectedStreamInput.messages)
	}
}

func TestInitStream_ToolUse(t *testing.T) {
	sendParams := provider.SendParams{
		Tool: types.Tool{
			Name:        "get_weather",
			Description: "Get the current weather for a location",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"location": map[string]any{
						"type":        "string",
						"description": "The city and state, e.g. San Francisco, CA",
					},
				},
				"required": []string{"location"},
			},
		},
	}

	expectedTool := anthropic.ToolUnionParamOfTool(
		anthropic.ToolInputSchemaParam{
			Type:       "object",
			Properties: sendParams.Tool.InputSchema, // Directly use the InputSchema from sendParams
		},
		"get_weather",
	)
	expectedTool.OfTool.Description = anthropic.String("Get the current weather for a location")

	expectedTools := []anthropic.ToolUnionParam{expectedTool}

	actualStreamInput := initStream(sendParams)

	if len(actualStreamInput.tools) != len(expectedTools) {
		t.Fatalf("Tool count mismatch: got %d, want %d", len(actualStreamInput.tools), len(expectedTools))
	}
	if len(expectedTools) > 0 {
		actualOfTool := actualStreamInput.tools[0].OfTool
		expectedOfTool := expectedTools[0].OfTool

		if actualOfTool.Name != expectedOfTool.Name {
			t.Errorf("Tool Name mismatch: got %s, want %s", actualOfTool.Name, expectedOfTool.Name)
		}
		// Correctly access the string value from param.Opt[string]
		if actualOfTool.Description.Value != expectedOfTool.Description.Value {
			t.Errorf("Tool Description mismatch: got %s, want %s", actualOfTool.Description.Value, expectedOfTool.Description.Value)
		}

		if actualOfTool.InputSchema.Type != expectedOfTool.InputSchema.Type {
			t.Errorf("Tool InputSchema.Type mismatch: got %s, want %s", actualOfTool.InputSchema.Type, expectedOfTool.InputSchema.Type)
		}
		if !reflect.DeepEqual(actualOfTool.InputSchema.Properties, expectedOfTool.InputSchema.Properties) {
			t.Errorf("Tool InputSchema.Properties mismatch:\nGot: %+v\nWant: %+v", actualOfTool.InputSchema.Properties, expectedOfTool.InputSchema.Properties)
		}
	}
}
