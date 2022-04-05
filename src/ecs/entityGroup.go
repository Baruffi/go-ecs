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

func (g *EntityGroup) HasMember(e Entity) bool {
	_, ok := g.members[e.id]
	return ok
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

func (g *EntityGroup) JoinScene(scene *Scene) {
	for _, e := range g.members {
		e.JoinScene(scene)
	}
}

func HasComponents[D ComponentData](g EntityGroup) (hasComponent bool) {
	for _, e := range g.members {
		hasComponent = hasComponent && Has[D](e.scene.registry, e.id)
	}
	return hasComponent
}

func HasComponentGroups[D ComponentData](g EntityGroup) (hasGroup bool) {
	for _, e := range g.members {
		hasGroup = hasGroup && HasGroup[D](e.scene.registry, e.id)
	}
	return hasGroup
}

func GetComponents[D ComponentData](g EntityGroup, is ...TypedComponentId[D]) []D {
	components := make([]D, 0)
	switch {
	case len(is) > 1:
		panic("More than 1 component id for component type is not allowed")
	case len(is) == 1:
		for _, e := range g.members {
			if c, ok := GetById(e.scene.registry, is[0], e.id); ok {
				components = append(components, c)
			}
		}
	default:
		for _, e := range g.members {
			if c, ok := Get[D](e.scene.registry, e.id); ok {
				components = append(components, c)
			}
		}
	}
	return components
}

func GetComponentGroups[D ComponentData](g EntityGroup, is ...TypedComponentGroupId[D]) []D {
	components := make([]D, 0)
	switch {
	case len(is) > 1:
		panic("More than 1 component id for component type is not allowed")
	case len(is) == 1:
		for _, e := range g.members {
			components = append(components, GetGroupById(e.scene.registry, is[0], e.id)...)
		}
		return components
	default:
		for _, e := range g.members {
			components = append(components, GetGroup[D](e.scene.registry, e.id)...)
		}
	}
	return components
}

func AddComponents[D ComponentData](g EntityGroup, c D) (typedComponentId TypedComponentId[D]) {
	for _, e := range g.members {
		typedComponentId = Link(e.scene.registry, e.id, c)
	}
	return typedComponentId
}

func AddComponentGroups[D ComponentData](g EntityGroup, c D) (typedGroupId TypedComponentGroupId[D]) {
	for _, e := range g.members {
		typedGroupId = Group[D](e.scene.registry, e.id, c)
	}
	return typedGroupId
}

func RemoveComponents[D ComponentData](g EntityGroup) {
	for _, e := range g.members {
		Unlink[D](e.scene.registry, e.id)
	}
}

func RemoveComponentGroups[D ComponentData](g EntityGroup) {
	for _, e := range g.members {
		Ungroup[D](e.scene.registry, e.id)
	}
}

func ClearComponentTypes[D ComponentData](g EntityGroup) {
	for _, e := range g.members {
		ClearType[D](e.scene.registry)
		// Clear operations only need to happen once
		break
	}
}

func ClearComponentGroups[D ComponentData](g EntityGroup) {
	for _, e := range g.members {
		ClearGroup[D](e.scene.registry)
		// Clear operations only need to happen once
		break
	}
}
