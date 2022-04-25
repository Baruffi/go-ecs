package ecs

import (
	"reflect"
)

type anyComponent interface{}

type anyComponentPool interface {
	Assign(EntityIndex, anyComponent)
	GetComponent(EntityIndex) anyComponent
	HasEntity(EntityIndex) bool
	RemoveEntity(EntityIndex)
}

type Scene struct {
	destroyed            EntityId
	entities             []EntityId
	componentPoolCounter int
	componentPoolMap     map[reflect.Type]int
	componentPools       []anyComponentPool
}

func NewScene() *Scene {
	return &Scene{
		destroyed:        CreateEntityId(INVALID_ENTITY, 0),
		entities:         make([]EntityId, 0),
		componentPoolMap: make(map[reflect.Type]int),
		componentPools:   make([]anyComponentPool, 0),
	}
}

func (scene *Scene) getPoolId(component anyComponent) int {
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
		destroyedIndex := GetEntityIndex(scene.destroyed)
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

func HasComponent[T anyComponent](scene *Scene, entityId EntityId) bool {
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

func Assign[T anyComponent](scene *Scene, entityId EntityId) *T {
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
	pool := scene.componentPools[poolId]
	entityIndex := GetEntityIndex(entityId)
	pool.Assign(entityIndex, &component)
	return &component
}

func GetComponent[T anyComponent](scene *Scene, entityId EntityId) (*T, bool) {
	if scene.HasEntity(entityId) {
		var example *T
		poolId := scene.getPoolId(example)
		if poolId < len(scene.componentPools) {
			pool := scene.componentPools[poolId]
			entityIndex := GetEntityIndex(entityId)
			if pool.HasEntity(entityIndex) {
				component := pool.GetComponent(entityIndex).(*T)
				return component, true
			}
		}
	}
	return nil, false
}

func RemoveComponent[T anyComponent](scene *Scene, entityId EntityId) {
	if !scene.HasEntity(entityId) {
		panic("cannot remove a component from an entity that is not in the scene")
	}
	var example *T
	poolId := scene.getPoolId(example)
	if poolId < len(scene.componentPools) {
		pool := scene.componentPools[poolId]
		entityIndex := GetEntityIndex(entityId)
		if pool.HasEntity(entityIndex) {
			pool.RemoveEntity(entityIndex)
		}
	}
}

type EntityView struct {
	currEntityIndex EntityIndex
	examples        []anyComponent
	scene           *Scene
}

func (view *EntityView) GetEntity() Entity {
	return Entity{
		id:    view.scene.entities[view.currEntityIndex],
		scene: view.scene,
	}
}

func (view *EntityView) GetComponents() []anyComponent {
	components := make([]anyComponent, 0)
	for _, example := range view.examples {
		poolId := view.scene.getPoolId(example)
		pool := view.scene.componentPools[poolId]
		if pool.HasEntity(view.currEntityIndex) {
			component := pool.GetComponent(view.currEntityIndex)
			components = append(components, component)
		}
	}
	return components
}

func (view *EntityView) Next() bool {
	temp := view.currEntityIndex + 1
	for !(temp > EntityIndex(len(view.scene.entities))) {
		for _, example := range view.examples {
			poolId := view.scene.getPoolId(example)
			if poolId < len(view.scene.componentPools) {
				pool := view.scene.componentPools[poolId]
				if pool.HasEntity(temp) {
					view.currEntityIndex = temp
					return true
				}
			}
		}
		temp++
	}
	view.currEntityIndex = INVALID_ENTITY
	return false
}

func View(scene *Scene, examples ...anyComponent) *EntityView {
	view := &EntityView{
		currEntityIndex: INVALID_ENTITY,
		scene:           scene,
		examples:        examples,
	}

	return view
}
