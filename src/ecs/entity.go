package ecs

import "math"

type EntityIndex uint32
type EntityVersion uint32
type EntityId uint64

const (
	INVALID_ENTITY EntityIndex = EntityIndex(math.MaxUint32)
)

func CreateEntityId(index EntityIndex, version EntityVersion) EntityId {
	return EntityId(index)<<32 | EntityId(version)
}

func GetEntityIndex(id EntityId) EntityIndex {
	return EntityIndex(id >> 32)
}

func GetEntityVersion(id EntityId) EntityVersion {
	return EntityVersion(id)
}

func IsEntityValid(id EntityId) bool {
	return EntityIndex(id>>32) != INVALID_ENTITY
}

type Entity struct {
	id    EntityId
	scene *Scene
}

func (entity *Entity) Die() {
	entity.scene.RemoveEntity(entity.id)
	entity.scene = nil
}

func Has[T any](entity Entity) bool {
	return HasComponent[T](entity.scene, entity.id)
}

func Add[T any](entity Entity) *T {
	return Assign[T](entity.scene, entity.id)
}

func Get[T any](entity Entity) (component *T, ok bool) {
	return GetComponent[T](entity.scene, entity.id)
}
