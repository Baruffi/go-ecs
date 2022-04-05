package ecs

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

func (g EntityGroup) HasMember(e Entity) bool {
	_, ok := g.members[e.id]
	return ok
}

func (g EntityGroup) AddMember(e Entity) {
	g.members[e.id] = e
}

func (g EntityGroup) RemoveMember(e Entity) {
	delete(g.members, e.id)
}

func (g EntityGroup) JoinScene(scene *Scene) {
	for _, e := range g.members {
		e.JoinScene(scene)
	}
}

func HasComponents[D ComponentData](g EntityGroup, is ...TypedComponentId[D]) (hasComponent bool) {
	for _, e := range g.members {
		hasComponent = hasComponent && HasComponent(e, is...)
	}
	return hasComponent
}

func HasComponentGroups[D ComponentData](g EntityGroup, is ...TypedComponentGroupId[D]) (hasGroup bool) {
	for _, e := range g.members {
		hasGroup = hasGroup && HasComponentGroup(e, is...)
	}
	return hasGroup
}

func GetComponents[D ComponentData](g EntityGroup, is ...TypedComponentId[D]) []D {
	components := make([]D, 0)
	for _, e := range g.members {
		if d, ok := GetComponent(e, is...); ok {
			components = append(components, d)
		}
	}
	return components
}

func GetComponentGroups[D ComponentData](g EntityGroup, is ...TypedComponentGroupId[D]) []D {
	components := make([]D, 0)
	for _, e := range g.members {
		components = append(components, GetComponentGroup(e, is...)...)
	}
	return components
}

func AddComponents[D ComponentData](g EntityGroup, d D) (typedComponentId TypedComponentId[D]) {
	for _, e := range g.members {
		typedComponentId = AddComponent(e, d)
	}
	return typedComponentId
}

func AddComponentGroups[D ComponentData](g EntityGroup, d D) (typedGroupId TypedComponentGroupId[D]) {
	for _, e := range g.members {
		typedGroupId = AddComponentGroup(e, d)
	}
	return typedGroupId
}

func RemoveComponents[D ComponentData](g EntityGroup, is ...TypedComponentId[D]) {
	for _, e := range g.members {
		RemoveComponent(e, is...)
	}
}

func RemoveComponentGroups[D ComponentData](g EntityGroup, is ...TypedComponentGroupId[D]) {
	for _, e := range g.members {
		RemoveComponentGroup(e, is...)
	}
}
