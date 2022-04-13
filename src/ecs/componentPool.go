package ecs

import "math"

type ComponentIndex uint32

const (
	INVALID_COMPONENT ComponentIndex = ComponentIndex(math.MaxUint32)
)

type ComponentPool[T any] struct {
	componentIndexes []ComponentIndex // sparse slice of ComponentIndex, which has indices matching the respective entityIndex
	entityIndexes    []EntityIndex    // packed slice of EntityIndex, which has indices matching the respective componentIndex

	components []*T
}

func NewComponentPool[T any]() *ComponentPool[T] {
	return &ComponentPool[T]{
		componentIndexes: make([]ComponentIndex, 0),
		entityIndexes:    make([]EntityIndex, 0),
		components:       make([]*T, 0),
	}
}

func (pool *ComponentPool[T]) HasEntity(entityIndex EntityIndex) bool {
	return len(pool.componentIndexes) > int(entityIndex) && pool.componentIndexes[entityIndex] != INVALID_COMPONENT
}

func (pool *ComponentPool[T]) RemoveEntity(entityIndex EntityIndex) {
	targetComponentIndex := pool.componentIndexes[entityIndex]
	lastEntityIndex := pool.entityIndexes[len(pool.entityIndexes)-1]
	lastComponent := pool.components[len(pool.components)-1]
	pool.componentIndexes[entityIndex] = INVALID_COMPONENT
	pool.componentIndexes[lastEntityIndex] = targetComponentIndex
	pool.entityIndexes[targetComponentIndex] = lastEntityIndex
	pool.components[targetComponentIndex] = lastComponent
	pool.entityIndexes = pool.entityIndexes[:len(pool.entityIndexes)-1]
	pool.components = pool.components[:len(pool.components)-1]
}
