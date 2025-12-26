package ports

import (
	"context"

	"github.com/YuruDeveloper/codey/internal/types"
	"github.com/google/uuid"
)

type AppCheckpointer interface {
	Save(ctx context.Context,threadID uuid.UUID,state *types.GraphState) error
	Load(ctx context.Context,threadID uuid.UUID) (*types.GraphState,error)
	List(ctx context.Context,threadID uuid.UUID) ([]uuid.UUID, error)
	LoadVersion(ctx context.Context,threadID uuid.UUID,version int) (*types.GraphState, error)
}

type AppRunner interface {
	Invoke(router AppRouter,appChecker AppCheckpointer,ctx context.Context,threadID uuid.UUID,input types.GraphState)
}