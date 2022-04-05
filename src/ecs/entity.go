package ecs

type EntityId string

type Entity struct {
	id    EntityId
	scene *Scene
}

// NewEntity - Creates a new entity filling in required initialization parameters
func NewEntity(scene *Scene) Entity {
	return Entity{
		id:    EntityId(GenerateId()),
		scene: scene,
	}
}

func (e *Entity) JoinScene(scene *Scene) {
	e.scene = scene
}

func AddComponent[C ComponentData](e Entity, c C) TypedComponentId[C] {
	return Link(e.scene.registry, e.id, c)
}

func AddComponentGroup[C ComponentData](e Entity, c C) TypedComponentGroupId[C] {
	return Group[C](e.scene.registry, e.id, c)
}

func HasComponent[C ComponentData](e Entity, is ...TypedComponentId[C]) bool {
	switch {
	case len(is) > 1:
		panic("More than 1 component id for component type is not allowed")
	case len(is) == 1:
		return HasById(e.scene.registry, is[0], e.id)
	default:
		return Has[C](e.scene.registry, e.id)
	}
}

func HasComponentGroup[C ComponentData](e Entity, is ...TypedComponentGroupId[C]) bool {
	switch {
	case len(is) > 1:
		panic("More than 1 component id for component type is not allowed")
	case len(is) == 1:
		return HasGroupById(e.scene.registry, is[0], e.id)
	default:
		return HasGroup[C](e.scene.registry, e.id)
	}
}

func GetComponent[C ComponentData](e Entity, is ...TypedComponentId[C]) (C, bool) {
	switch {
	case len(is) > 1:
		panic("More than 1 component id for component type is not allowed")
	case len(is) == 1:
		return GetById(e.scene.registry, is[0], e.id)
	default:
		return Get[C](e.scene.registry, e.id)
	}
}

func GetComponentGroup[C ComponentData](e Entity, is ...TypedComponentGroupId[C]) []C {
	switch {
	case len(is) > 1:
		panic("More than 1 component id for component type is not allowed")
	case len(is) == 1:
		return GetGroupById(e.scene.registry, is[0], e.id)
	default:
		return GetGroup[C](e.scene.registry, e.id)
	}
}

func RemoveComponent[C ComponentData](e Entity, is ...TypedComponentId[C]) {
	switch {
	case len(is) > 1:
		panic("More than 1 component id for component type is not allowed")
	case len(is) == 1:
		UnlinkById(e.scene.registry, is[0], e.id)
	default:
		Unlink[C](e.scene.registry, e.id)
	}
}

func RemoveComponentGroup[C ComponentData](e Entity, is ...TypedComponentGroupId[C]) {
	switch {
	case len(is) > 1:
		panic("More than 1 component id for component type is not allowed")
	case len(is) == 1:
		UngroupById(e.scene.registry, is[0], e.id)
	default:
		Ungroup[C](e.scene.registry, e.id)
	}
}
