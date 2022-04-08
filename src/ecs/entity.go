package ecs

type EntityId string

type Entity struct {
	id    EntityId
	scene *Scene
}

// NewEntity - Creates a new entity filling in required initialization parameters
func NewEntity(scene *Scene, manualIdInput ...string) Entity {
	var id string
	if len(manualIdInput) == 1 {
		id = manualIdInput[0]
	} else {
		id = GenerateId()
	}
	return Entity{
		id:    EntityId(id),
		scene: scene,
	}
}

func (e Entity) JoinScene(scene *Scene) {
	*e.scene = *scene
}

func AddComponent[D ComponentData](e Entity, d D) ComponentId {
	return Link(e.scene.Registry, e.id, d)
}

func AddComponentGroup[D ComponentData](e Entity, d D) (ComponentId, ComponentGroupId) {
	return Group[D](e.scene.Registry, e.id, d)
}

func HasComponent[D ComponentData](e Entity) bool {
	return Has[D](e.scene.Registry, e.id)
}

func HasComponentGroup[D ComponentData](e Entity) bool {
	return HasGroup[D](e.scene.Registry, e.id)
}

func GetComponent[D ComponentData](e Entity) (D, bool) {
	return Get[D](e.scene.Registry, e.id)
}

func GetComponentGroup[D ComponentData](e Entity) map[ComponentId]D {
	return GetGroup[D](e.scene.Registry, e.id)
}

func RemoveComponent[D ComponentData](e Entity) {
	Unlink[D](e.scene.Registry, e.id)
}

func RemoveComponentGroup[D ComponentData](e Entity) {
	Ungroup[D](e.scene.Registry, e.id)
}
