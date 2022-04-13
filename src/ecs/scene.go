package ecs

import (
	"reflect"
)

type Updater interface {
	Update(float64)
}

type UpdaterFunc func(float64)

func (f UpdaterFunc) Update(dt float64) {
	f(dt)
}

type AnyComponentPool interface {
	HasEntity(entityId EntityIndex) bool
	RemoveEntity(entityId EntityIndex)
}

type Scene struct {
	Updater
	destroyed            EntityId
	entities             []EntityId
	componentPoolCounter int
	componentPoolMap     map[reflect.Type]int
	componentPools       []AnyComponentPool
}

func NewScene(updater Updater) *Scene {
	return &Scene{
		Updater:          updater,
		destroyed:        CreateEntityId(INVALID_ENTITY, 0),
		entities:         make([]EntityId, 0),
		componentPoolMap: make(map[reflect.Type]int),
		componentPools:   make([]AnyComponentPool, 0),
	}
}

func (scene *Scene) getPoolId(component any) (id int) {
	reflectType := reflect.TypeOf(component)
	if reflectType == nil {
		panic("a nil interface is not a valid component type")
	}
	id, ok := scene.componentPoolMap[reflectType]
	if !ok {
		id = scene.componentPoolCounter
	}
	return id
}

func (scene *Scene) CreateEntity() Entity {
	if IsEntityValid(scene.destroyed) {
		destroyed := scene.destroyed
		if IsEntityValid(scene.entities[GetEntityIndex(destroyed)]) {
			scene.destroyed = scene.entities[GetEntityIndex(destroyed)]
		} else {
			scene.destroyed = CreateEntityId(INVALID_ENTITY, 0)
		}
		scene.entities[destroyed] = destroyed
		return Entity{
			id:    scene.entities[destroyed],
			scene: scene,
		}
	} else {
		id := CreateEntityId(EntityIndex(len(scene.entities)), 0)
		scene.entities = append(scene.entities, id)
		return Entity{
			id:    id,
			scene: scene,
		}
	}
}

func (scene *Scene) RemoveEntity(entityId EntityId) {
	if IsEntityValid(entityId) && IsEntityValid(scene.entities[GetEntityIndex(entityId)]) {
		entityVersion := GetEntityVersion(scene.entities[GetEntityIndex(entityId)])
		scene.entities[GetEntityIndex(entityId)] = scene.destroyed
		scene.destroyed = CreateEntityId(GetEntityIndex(entityId), entityVersion+1)
		for _, pool := range scene.componentPools {
			if pool.HasEntity(GetEntityIndex(entityId)) {
				pool.RemoveEntity(GetEntityIndex(entityId))
			}
		}
	}
}

func HasComponent[T any](scene *Scene, entityId EntityId) bool {
	var base T
	poolId := scene.getPoolId(base)
	if poolId < len(scene.componentPools) {
		pool := scene.componentPools[poolId]
		return pool.HasEntity(GetEntityIndex(entityId))
	}
	return false
}

func Assign[T any](scene *Scene, entityId EntityId, component T) T {
	poolId := scene.getPoolId(component)
	if len(scene.componentPools) <= poolId {
		reflectType := reflect.TypeOf(component)
		scene.componentPoolMap[reflectType] = poolId
		scene.componentPools = append(scene.componentPools, NewComponentPool[T]())
		scene.componentPoolCounter++
	}
	pool := scene.componentPools[poolId].(*ComponentPool[T])
	pool.entityIndexes = append(pool.entityIndexes, GetEntityIndex(entityId))
	pool.components = append(pool.components, component)
	for len(pool.componentIndexes) <= int(GetEntityIndex(entityId)) {
		pool.componentIndexes = append(pool.componentIndexes, INVALID_COMPONENT)
	}
	pool.componentIndexes[GetEntityIndex(entityId)] = ComponentIndex(len(pool.components) - 1)
	return component
}

func GetComponent[T any](scene *Scene, entityId EntityId) (component T, ok bool) {
	poolId := scene.getPoolId(component)
	if poolId < len(scene.componentPools) {
		pool, ok := scene.componentPools[poolId].(*ComponentPool[T])
		if ok {
			componentIndex := pool.componentIndexes[GetEntityIndex(entityId)]
			if componentIndex != INVALID_COMPONENT {
				component = pool.components[componentIndex]
				return component, true
			}
		}
	}
	return component, false
}

func View[T any](scene *Scene) (components []T, ok bool) {
	var base T
	poolId := scene.getPoolId(base)
	pool, ok := scene.componentPools[poolId].(*ComponentPool[T])
	if ok {
		components = pool.components
	}

	return components, ok
}
