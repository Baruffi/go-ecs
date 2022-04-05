package ecs

type ComponentGroupId string

type TypedComponentGroupId[C ComponentData] ComponentGroupId

type ComponentGroup[C ComponentData] struct {
	id      ComponentGroupId
	members map[ComponentId]Component[C]
}

// NewComponentGroup - Creates a new ComponentGroup filling in required initialization parameters
func NewComponentGroup[C ComponentData]() ComponentGroup[C] {
	return ComponentGroup[C]{
		id:      ComponentGroupId(GenerateId()),
		members: make(map[ComponentId]Component[C]),
	}
}

func (g ComponentGroup[C]) GetId() ComponentGroupId {
	return g.id
}

func (g ComponentGroup[C]) Has(e EntityId) bool {
	for _, c := range g.members {
		if c.Has(e) {
			return true
		}
	}
	return false
}

func (g ComponentGroup[C]) Get(e EntityId) []ComponentData {
	ds := make([]ComponentData, 0)
	for _, c := range g.members {
		if d, ok := c.Get(e); ok {
			ds = append(ds, d)
		}
	}
	return ds
}

func (g ComponentGroup[C]) Set(c AnyComponent) {
	g.members[c.GetId()] = c.(Component[C])
}

func (g ComponentGroup[C]) Unset(c ComponentId) {
	delete(g.members, c)
}

func (g ComponentGroup[C]) UnsetEntity(e EntityId) {
	for _, c := range g.members {
		c.Unset(e)
	}
}
