package ecs

import (
	"github.com/google/uuid"
)

// generateId - Generate a random uint64 id
func generateId() string {
	return uuid.NewString()
}
