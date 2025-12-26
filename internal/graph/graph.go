package graph

import (
	"context"
	"github.com/YuruDeveloper/codey/internal/ports"
	"github.com/YuruDeveloper/codey/internal/types"
	"github.com/google/uuid"
)

var _ ports.AppRunner = (*Runner)(nil)

type Runner struct {

}

func (r *Runner) Invoke(router ports.AppRouter,appChecker ports.AppCheckpointer, ctx context.Context, threadID uuid.UUID, input types.GraphState)  {
	
}
