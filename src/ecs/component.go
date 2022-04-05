package ecs

type ComponentId string

type TypedComponentId[C ComponentData] ComponentId

type ComponentData interface {
}

type Component[C ComponentData] struct {
	id   ComponentId
	data map[EntityId]C
}

// NewComponent - Creates a new component filling in required initialization parameters
func NewComponent[C ComponentData]() Component[C] {
	return Component[C]{
		id:   ComponentId(GenerateId()),
		data: make(map[EntityId]C),
	}
}

func (c Component[C]) GetId() ComponentId {
	return c.id
}

func (c Component[C]) Has(e EntityId) bool {
	_, ok := c.data[e]
	return ok
}

func (c Component[C]) Get(e EntityId) (ComponentData, bool) {
	d, ok := c.data[e]
	return d, ok
}

func (c Component[C]) Set(e EntityId, d ComponentData) {
	c.data[e] = d.(C)
}

func (c Component[C]) Unset(e EntityId) {
	delete(c.data, e)
}
