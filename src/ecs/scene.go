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

type anyComponentPool interface {
	HasEntity(entityId EntityIndex) bool
	RemoveEntity(entityId EntityIndex)
}

type Scene struct {
	Updater
	destroyed            EntityId
	entities             []EntityId
	componentPoolCounter int
	componentPoolMap     map[reflect.Type]int
	componentPools       []anyComponentPool
}

func NewScene[U Updater]() *Scene {
	var u U
	return &Scene{
		Updater:          u,
		destroyed:        CreateEntityId(INVALID_ENTITY, 0),
		entities:         make([]EntityId, 0),
		componentPoolMap: make(map[reflect.Type]int),
		componentPools:   make([]anyComponentPool, 0),
	}
}

func (scene *Scene) getPoolId(component any) int {
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

func (scene *Scene) HasEntity(entityId EntityId) bool {
	if IsEntityValid(entityId) && GetEntityIndex(entityId) < EntityIndex(len(scene.entities)) && IsEntityValid(scene.entities[GetEntityIndex(entityId)]) && entityId == scene.entities[GetEntityIndex(entityId)] {
		return true
	}
	return false
}

func (scene *Scene) CreateEntity() Entity {
	if IsEntityValid(scene.destroyed) {
		destroyed := scene.destroyed
		destroyedIndex := GetEntityIndex(destroyed)
		if IsEntityValid(scene.entities[destroyedIndex]) {
			scene.destroyed = scene.entities[destroyedIndex]
		} else {
			scene.destroyed = CreateEntityId(INVALID_ENTITY, 0)
		}
		scene.entities[destroyedIndex] = destroyed
		return Entity{
			id:    destroyed,
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
	if scene.HasEntity(entityId) {
		entityIndex := GetEntityIndex(entityId)
		scene.entities[entityIndex] = scene.destroyed
		scene.destroyed = CreateEntityId(entityIndex, GetEntityVersion(entityId)+1)
		for _, pool := range scene.componentPools {
			if pool.HasEntity(entityIndex) {
				pool.RemoveEntity(entityIndex)
			}
		}
	}
}

func HasComponent[T any](scene *Scene, entityId EntityId) bool {
	if scene.HasEntity(entityId) {
		var example *T
		poolId := scene.getPoolId(example)
		if poolId < len(scene.componentPools) {
			pool := scene.componentPools[poolId]
			return pool.HasEntity(GetEntityIndex(entityId))
		}
	}
	return false
}

func Assign[T any](scene *Scene, entityId EntityId) *T {
	if !scene.HasEntity(entityId) {
		panic("cannot assign a component to an entity that is not in the scene")
	}
	var component T
	poolId := scene.getPoolId(&component)
	if len(scene.componentPools) <= poolId {
		reflectType := reflect.TypeOf(&component)
		scene.componentPoolMap[reflectType] = poolId
		scene.componentPools = append(scene.componentPools, NewComponentPool[T]())
		scene.componentPoolCounter++
	}
	pool := scene.componentPools[poolId].(*ComponentPool[T])
	entityIndex := GetEntityIndex(entityId)
	pool.Assign(entityIndex, &component)
	return &component
}

func GetComponent[T any](scene *Scene, entityId EntityId) (*T, bool) {
	if scene.HasEntity(entityId) {
		var example *T
		poolId := scene.getPoolId(example)
		if poolId < len(scene.componentPools) {
			pool, ok := scene.componentPools[poolId].(*ComponentPool[T])
			if ok {
				entityIndex := GetEntityIndex(entityId)
				if pool.HasEntity(entityIndex) {
					component := pool.GetComponent(entityIndex)
					return component, true
				}
			}
		}
	}
	return nil, false
}
