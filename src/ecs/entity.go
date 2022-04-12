package ecs

type Entity struct {
	id    int
	scene *Scene
}

func (entity *Entity) Die() {
	entity.scene.RemoveEntity(entity.id)
	entity.scene = nil
}

func Has[T any](entity Entity) bool {
	return HasComponent[T](entity.scene, entity.id)
}

func Add[T any](entity Entity, component T) T {
	return Assign(entity.scene, entity.id, component)
}

func Get[T any](entity Entity) (component T, ok bool) {
	return GetComponent[T](entity.scene, entity.id)
}
