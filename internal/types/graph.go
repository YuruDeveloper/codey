package types

import (
	"github.com/google/uuid"
)

type GraphState struct {
	Provider uuid.UUID
	Messages []Message
}
