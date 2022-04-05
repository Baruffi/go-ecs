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

type EntityGroupId string

type EntityGroup struct {
	id      EntityGroupId
	members map[EntityId]Entity
}

// NewEntityGroup - Creates a new EntityGroup filling in required initialization parameters
func NewEntityGroup() EntityGroup {
	return EntityGroup{
		id:      EntityGroupId(GenerateId()),
		members: make(map[EntityId]Entity),
	}
}

func (g *EntityGroup) AddMember(e Entity) {
	g.members[e.id] = e
}

func (g *EntityGroup) RemoveMember(e Entity) {
	delete(g.members, e.id)
}

func (g *EntityGroup) Clear() {
	g.members = make(map[EntityId]Entity)
}
