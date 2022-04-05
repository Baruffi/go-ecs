package ecs

type ComponentGroupId string

type TypedComponentGroupId[D ComponentData] ComponentGroupId

type ComponentGroup[D ComponentData] struct {
	id      ComponentGroupId
	members map[ComponentId]Component[D]
}

// NewComponentGroup - Creates a new ComponentGroup filling in required initialization parameters
func NewComponentGroup[D ComponentData]() ComponentGroup[D] {
	return ComponentGroup[D]{
		id:      ComponentGroupId(GenerateId()),
		members: make(map[ComponentId]Component[D]),
	}
}

func (g ComponentGroup[D]) GetId() ComponentGroupId {
	return g.id
}

func (g ComponentGroup[D]) Has(e EntityId) bool {
	for _, c := range g.members {
		if c.Has(e) {
			return true
		}
	}
	return false
}

func (g ComponentGroup[D]) Get(e EntityId) []ComponentData {
	ds := make([]ComponentData, 0)
	for _, c := range g.members {
		if d, ok := c.Get(e); ok {
			ds = append(ds, d)
		}
	}
	return ds
}

func (g ComponentGroup[D]) Set(c AnyComponent) {
	g.members[c.GetId()] = c.(Component[D])
}

func (g ComponentGroup[D]) Unset(c ComponentId) {
	delete(g.members, c)
}

func (g ComponentGroup[D]) UnsetEntity(e EntityId) {
	for _, c := range g.members {
		c.Unset(e)
	}
}
