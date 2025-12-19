package ports

import (
	"context"

	"github.com/YuruDeveloper/codey/internal/types"
	"github.com/google/uuid"
)

type AppNodeEngine interface {
	AddNode(name uuid.UUID,fn types.NodeFunc)
	AddEdge(from uuid.UUID,edge types.NodeEdge)
	Invoke(ctx context.Context,threadID uuid.UUID,input types.GraphState) types.GraphUpdate
}

type AppCheckpointer interface {
	Save(ctx context.Context,threadID uuid.UUID,state *types.GraphState) error
	Load(ctx context.Context,threadID uuid.UUID) (*types.GraphState,error)
	List(ctx context.Context,threadID uuid.UUID) ([]uuid.UUID, error)
	LoadVersion(ctx context.Context,threadID uuid.UUID,version int) (*types.GraphState, error)
}

type AppRunner interface {
	
}