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

func HasComponent[C ComponentData](e Entity) bool {
	return Has[C](e.scene.registry, e.id)
}

func HasComponentGroup[C ComponentData](e Entity) bool {
	return HasGroup[C](e.scene.registry, e.id)
}

func GetComponent[C ComponentData](e Entity, is ...TypedComponentId[C]) (C, bool) {
	// Since only 1 id per type will exist in the registry, there should be no use case with multiple ids as args. Using ... as an optional notation
	for _, i := range is {
		return GetById(e.scene.registry, i, e.id)
	}
	return Get[C](e.scene.registry, e.id)
}

func GetComponentGroup[C ComponentData](e Entity, is ...TypedComponentId[C]) []C {
	// Since only 1 id per type will exist in the registry, there should be no use case with multiple ids as args. Using ... as an optional notation
	for _, i := range is {
		return GetGroupById(e.scene.registry, i, e.id)
	}
	return GetGroup[C](e.scene.registry, e.id)
}

func AddComponent[C ComponentData](e Entity, c C) TypedComponentId[C] {
	return Link(e.scene.registry, e.id, c)
}

func AddComponentGroup[C ComponentData](e Entity, c C) TypedComponentId[C] {
	return Group[C](e.scene.registry, e.id, c)
}

func RemoveComponent[C ComponentData](e Entity) {
	Unlink[C](e.scene.registry, e.id)
}

func RemoveComponentGroup[C ComponentData](e Entity) {
	Ungroup[C](e.scene.registry, e.id)
}

func ClearComponentType[C ComponentData](e Entity) {
	ClearType[C](e.scene.registry)
}

func ClearComponentGroup[C ComponentData](e Entity) {
	ClearGroup[C](e.scene.registry)
}

func ClearComponents(e Entity) {
	Clear(e.scene.registry)
}
