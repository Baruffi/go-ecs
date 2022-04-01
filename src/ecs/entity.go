package ecs

type EntityId uint64

type Entity struct {
	id EntityId
}

// NewEntity - Creates a new entity filling in required initialization parameters
func NewEntity() Entity {
	return Entity{
		id: EntityId(generateId()),
	}
}
