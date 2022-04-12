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
	HasEntity(entityId int) bool
	RemoveEntity(entityId int)
}

type Scene struct {
	Updater
	entityCounter        int
	componentPoolCounter int
	componentPoolMap     map[reflect.Type]int
	componentPools       []AnyComponentPool
}

func NewScene(updater Updater) *Scene {
	return &Scene{
		Updater:          updater,
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
	id := scene.entityCounter
	scene.entityCounter++
	return Entity{
		id:    id,
		scene: scene,
	}
}

func (scene *Scene) RemoveEntity(entityId int) {
	for _, pool := range scene.componentPools {
		if pool.HasEntity(entityId) {
			pool.RemoveEntity(entityId)
		}
	}
}

func HasComponent[T any](scene *Scene, entityId int) bool {
	var base T
	poolId := scene.getPoolId(base)
	if poolId < len(scene.componentPools) {
		pool := scene.componentPools[poolId]
		return pool.HasEntity(entityId)
	}
	return false
}

func Assign[T any](scene *Scene, entityId int, component T) T {
	poolId := scene.getPoolId(component)
	if len(scene.componentPools) <= poolId {
		reflectType := reflect.TypeOf(component)
		scene.componentPoolMap[reflectType] = poolId
		scene.componentPools = append(scene.componentPools, NewComponentPool[T]())
		scene.componentPoolCounter++
	}
	pool := scene.componentPools[poolId].(*ComponentPool[T])
	pool.entityIds = append(pool.entityIds, entityId)
	pool.components = append(pool.components, component)
	for len(pool.componentIds) <= entityId {
		pool.componentIds = append(pool.componentIds, -1)
	}
	pool.componentIds[entityId] = len(pool.components) - 1
	return component
}

func GetComponent[T any](scene *Scene, entityId int) (component T, ok bool) {
	poolId := scene.getPoolId(component)
	if poolId < len(scene.componentPools) {
		pool, ok := scene.componentPools[poolId].(*ComponentPool[T])
		if ok {
			componentId := pool.componentIds[entityId]
			if componentId != -1 {
				component = pool.components[componentId]
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
