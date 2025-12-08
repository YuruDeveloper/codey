package anthropic

import (
	"github.com/YuruDeveloper/codey/internal/types"
	"github.com/anthropics/anthropic-sdk-go"
)

func roleAdapter(role types.Role) anthropic.MessageParamRole {
	switch role {
		case types.UserRole:
			return anthropic.MessageParamRoleUser
		case types.AssistantRole:
			return anthropic.MessageParamRoleAssistant
		case types.ToolRole:
			return anthropic.MessageParamRole("tool")
		case types.AlarmRole:
			return anthropic.MessageParamRole("alarm")
		default:
			return anthropic.MessageParamRole("unknown")
	}
}

func partsAdapter(parts []types.Part) []anthropic.ContentBlockParamUnion {
	blocks := make([]anthropic.ContentBlockParamUnion,0,len(parts))

	for _ , part := range parts {
		switch part := part.(type) {
			case types.TextPart:
				blocks = append(blocks, anthropic.NewTextBlock(part.Text))
			case types.ImagePart:
				blocks = append(blocks, anthropic.NewImageBlockBase64(part.MediaType,part.Data))
			case types.ToolUsePart:
				blocks = append(blocks, anthropic.NewToolUseBlock(part.ID,part.Input,""))
			case types.ToolResultPart:
				blocks = append(blocks, anthropic.NewToolResultBlock(
					part.ToolUseID,
					part.Content,
					part.IsError,
				))		
		}
	}
	return blocks
}

func toolAdapter(tool types.Tool) []anthropic.ToolUnionParam {
	result := make([]anthropic.ToolUnionParam,1)

	inputSchema := anthropic.ToolInputSchemaParam {
		Properties: tool.InputSchema,
	}

	toolParm := anthropic.ToolUnionParamOfTool(inputSchema,tool.Name)
	
	if toolParm.OfTool != nil {
		toolParm.OfTool.Description = anthropic.String(tool.Description)
	}

	result[0] = toolParm
	
	return result
}

func messageAdapter(messages []types.Message) []anthropic.MessageParam {
	result := make([]anthropic.MessageParam,0,len(messages))

	for _ , message := range messages {
		blocks := partsAdapter(message.Parts)
		result = append(result, anthropic.MessageParam{
			Role: roleAdapter(message.Role),
			Content: blocks,
		})
	}

	return result
}