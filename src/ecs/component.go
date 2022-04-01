package ecs

type ComponentId uint64

type ComponentData interface {
}

type Component[C ComponentData] struct {
	id   ComponentId
	data map[EntityId]C
}

// NewComponent - Creates a new component filling in required initialization parameters
func NewComponent[C ComponentData]() Component[C] {
	return Component[C]{
		id:   ComponentId(generateId()),
		data: make(map[EntityId]C),
	}
}
