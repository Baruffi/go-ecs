package ecs

import (
	"github.com/google/uuid"
)

// GenerateId - Generate a random uint64 id
func GenerateId() string {
	return uuid.NewString()
}
