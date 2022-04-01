package ecs

import (
	"math/rand"
)

// generateId - Generate a random uint64 id
func generateId() uint64 {
	return rand.Uint64()
}
