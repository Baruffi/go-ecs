package ecs

type ComponentId string

type TypedComponentId[D ComponentData] ComponentId

type ComponentData interface {
}

type Component[D ComponentData] struct {
	id   ComponentId
	data map[EntityId]D
}

// NewComponent - Creates a new component filling in required initialization parameters
func NewComponent[D ComponentData]() Component[D] {
	return Component[D]{
		id:   ComponentId(GenerateId()),
		data: make(map[EntityId]D),
	}
}

func (c Component[D]) GetId() ComponentId {
	return c.id
}

func (c Component[D]) Has(e EntityId) bool {
	_, ok := c.data[e]
	return ok
}

func (c Component[D]) Get(e EntityId) (ComponentData, bool) {
	d, ok := c.data[e]
	return d, ok
}

func (c Component[D]) Set(e EntityId, d ComponentData) {
	c.data[e] = d.(D)
}

func (c Component[D]) Unset(e EntityId) {
	delete(c.data, e)
}
