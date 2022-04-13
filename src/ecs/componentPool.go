package ecs

import "math"

type ComponentIndex uint32

const (
	INVALID_COMPONENT ComponentIndex = ComponentIndex(math.MaxUint32)
)

type ComponentPool[T any] struct {
	componentIds []ComponentIndex // sparse slice of ComponentIndex, which has indices matching the respective entityIndex
	entityIds    []EntityIndex    // packed slice of EntityIndex, which has indices matching the respective componentIndex

	components []T
}

func NewComponentPool[T any]() *ComponentPool[T] {
	return &ComponentPool[T]{
		componentIds: make([]ComponentIndex, 0),
		entityIds:    make([]EntityIndex, 0),
		components:   make([]T, 0),
	}
}

func (pool *ComponentPool[T]) HasEntity(entityId EntityIndex) bool {
	return len(pool.componentIds) > int(entityId) && pool.componentIds[entityId] != INVALID_COMPONENT
}

func (pool *ComponentPool[T]) RemoveEntity(entityId EntityIndex) {
	targetComponentId := pool.componentIds[entityId]
	lastEntityId := pool.entityIds[len(pool.entityIds)-1]
	lastComponent := pool.components[len(pool.components)-1]
	pool.componentIds[entityId] = INVALID_COMPONENT
	pool.componentIds[lastEntityId] = targetComponentId
	pool.entityIds[targetComponentId] = lastEntityId
	pool.components[targetComponentId] = lastComponent
	pool.entityIds = pool.entityIds[:len(pool.entityIds)-1]
	pool.components = pool.components[:len(pool.components)-1]
}
