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

func HasComponent[C ComponentData](e Entity, c C) bool {
	return Has(e.scene.registry, e.id, c)
}

func HasComponentGroup[C ComponentData](e Entity, c C) bool {
	return HasGroup(e.scene.registry, e.id, c)
}

func GetComponent[C ComponentData](e Entity) (C, bool) {
	return Get[C](e.scene.registry, e.id)
}

func GetComponentGroup[C ComponentData](e Entity) []C {
	return GetGroup[C](e.scene.registry, e.id)
}

func AddComponent[C ComponentData](e Entity, c C) Component[C] {
	return Link(e.scene.registry, e.id, c)
}

func AddComponentGroup[C ComponentData](e Entity, c C) ComponentGroup[C] {
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
