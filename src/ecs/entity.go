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

func (e Entity) JoinScene(scene *Scene) {
	*e.scene = *scene
}

func AddComponent[D ComponentData](e Entity, d D) TypedComponentId[D] {
	return Link(e.scene.registry, e.id, d)
}

func AddComponentGroup[D ComponentData](e Entity, d D) (TypedComponentId[D], TypedComponentGroupId[D]) {
	return Group[D](e.scene.registry, e.id, d)
}

func HasComponent[D ComponentData](e Entity, is ...TypedComponentId[D]) bool {
	switch {
	case len(is) > 1:
		panic("More than 1 component id for component type is not allowed")
	case len(is) == 1:
		return HasById(e.scene.registry, is[0], e.id)
	default:
		return Has[D](e.scene.registry, e.id)
	}
}

func HasComponentGroup[D ComponentData](e Entity, is ...TypedComponentGroupId[D]) bool {
	switch {
	case len(is) > 1:
		panic("More than 1 component id for component type is not allowed")
	case len(is) == 1:
		return HasGroupById(e.scene.registry, is[0], e.id)
	default:
		return HasGroup[D](e.scene.registry, e.id)
	}
}

func GetComponent[D ComponentData](e Entity, is ...TypedComponentId[D]) (D, bool) {
	switch {
	case len(is) > 1:
		panic("More than 1 component id for component type is not allowed")
	case len(is) == 1:
		return GetById(e.scene.registry, is[0], e.id)
	default:
		return Get[D](e.scene.registry, e.id)
	}
}

func GetComponentGroup[D ComponentData](e Entity, is ...TypedComponentGroupId[D]) []D {
	switch {
	case len(is) > 1:
		panic("More than 1 component id for component type is not allowed")
	case len(is) == 1:
		return GetGroupById(e.scene.registry, is[0], e.id)
	default:
		return GetGroup[D](e.scene.registry, e.id)
	}
}

// Ids are not worth it right now
func GetComponentFromGroup[D ComponentData](e Entity, ci TypedComponentId[D], gi TypedComponentGroupId[D]) (D, bool) {
	return GetFromGroup(e.scene.registry, gi, ci, e.id)
}

func RemoveComponent[D ComponentData](e Entity, is ...TypedComponentId[D]) {
	switch {
	case len(is) > 1:
		panic("More than 1 component id for component type is not allowed")
	case len(is) == 1:
		UnlinkById(e.scene.registry, is[0], e.id)
	default:
		Unlink[D](e.scene.registry, e.id)
	}
}

func RemoveComponentGroup[D ComponentData](e Entity, is ...TypedComponentGroupId[D]) {
	switch {
	case len(is) > 1:
		panic("More than 1 component id for component type is not allowed")
	case len(is) == 1:
		UngroupById(e.scene.registry, is[0], e.id)
	default:
		Ungroup[D](e.scene.registry, e.id)
	}
}
