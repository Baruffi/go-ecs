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

func (entity *Entity) GetId() EntityId {
	return entity.id
}

func (entity *Entity) IsAlive() bool {
	if IsEntityValid(entity.id) {
		return entity.scene.HasEntity(entity.id)
	}
	return false
}

func (entity *Entity) Die() {
	if IsEntityValid(entity.id) {
		entity.scene.RemoveEntity(entity.id)
		entity.id = CreateEntityId(INVALID_ENTITY, 0)
		entity.scene = nil
	}
}

func Has[T anyComponent](entity Entity) bool {
	if IsEntityValid(entity.id) {
		return HasComponent[T](entity.scene, entity.id)
	}
	return false
}

func Add[T anyComponent](entity Entity) *T {
	if IsEntityValid(entity.id) {
		return Assign[T](entity.scene, entity.id)
	}
	panic("cannot assign a component to a dead entity")
}

func Get[T anyComponent](entity Entity) (*T, bool) {
	if IsEntityValid(entity.id) {
		return GetComponent[T](entity.scene, entity.id)
	}
	return nil, false
}

func Remove[T anyComponent](entity Entity) {
	if IsEntityValid(entity.id) {
		RemoveComponent[T](entity.scene, entity.id)
	}
	panic("cannot remove a component from a dead entity")
}
