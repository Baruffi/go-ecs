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

func HasComponents[C ComponentData](g EntityGroup) (hasComponent bool) {
	for _, e := range g.members {
		hasComponent = hasComponent && Has[C](e.scene.registry, e.id)
	}
	return hasComponent
}

func HasComponentGroups[C ComponentData](g EntityGroup) (hasGroup bool) {
	for _, e := range g.members {
		hasGroup = hasGroup && HasGroup[C](e.scene.registry, e.id)
	}
	return hasGroup
}

func GetComponents[C ComponentData](g EntityGroup, is ...TypedComponentId[C]) []C {
	components := make([]C, 0)
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
			if c, ok := Get[C](e.scene.registry, e.id); ok {
				components = append(components, c)
			}
		}
	}
	return components
}

func GetComponentGroups[C ComponentData](g EntityGroup, is ...TypedComponentGroupId[C]) []C {
	components := make([]C, 0)
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
			components = append(components, GetGroup[C](e.scene.registry, e.id)...)
		}
	}
	return components
}

func AddComponents[C ComponentData](g EntityGroup, c C) (typedComponentId TypedComponentId[C]) {
	for _, e := range g.members {
		typedComponentId = Link(e.scene.registry, e.id, c)
	}
	return typedComponentId
}

func AddComponentGroups[C ComponentData](g EntityGroup, c C) (typedGroupId TypedComponentGroupId[C]) {
	for _, e := range g.members {
		typedGroupId = Group[C](e.scene.registry, e.id, c)
	}
	return typedGroupId
}

func RemoveComponents[C ComponentData](g EntityGroup) {
	for _, e := range g.members {
		Unlink[C](e.scene.registry, e.id)
	}
}

func RemoveComponentGroups[C ComponentData](g EntityGroup) {
	for _, e := range g.members {
		Ungroup[C](e.scene.registry, e.id)
	}
}

func ClearComponentTypes[C ComponentData](g EntityGroup) {
	for _, e := range g.members {
		ClearType[C](e.scene.registry)
		// Clear operations only need to happen once
		break
	}
}

func ClearComponentGroups[C ComponentData](g EntityGroup) {
	for _, e := range g.members {
		ClearGroup[C](e.scene.registry)
		// Clear operations only need to happen once
		break
	}
}
