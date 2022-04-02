package ecs

type ComponentId string

type ComponentData interface {
}

type Component[C ComponentData] struct {
	id       ComponentId
	entities map[EntityId]Entity
	data     map[EntityId]C
}

// NewComponent - Creates a new component filling in required initialization parameters
func NewComponent[C ComponentData]() Component[C] {
	return Component[C]{
		id:       ComponentId(generateId()),
		entities: make(map[EntityId]Entity),
		data:     make(map[EntityId]C),
	}
}
