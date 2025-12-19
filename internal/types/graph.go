package types

import (
	"context"

	"github.com/google/uuid"
)

type GraphState struct {
	Provider uuid.UUID
	Messages []Message
	CurrentStep uuid.UUID
}

type GraphUpdate struct {
	Message Message
}

type NodeFunc func(ctx context.Context,state *GraphState) GraphUpdate 

type NodeEdge func(state *GraphState) uuid.UUID 

