package types

type Message struct {
	Role  Role
	Parts []Part
}


type SendParams struct {
	Messages     []Message
	Tool         Tool
	SystemPrompt string
	MaxTokens    int
}