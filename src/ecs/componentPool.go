package ecs

import "math"

type ComponentIndex uint32

const (
	INVALID_COMPONENT ComponentIndex = ComponentIndex(math.MaxUint32)
	PAGE_SIZE         EntityIndex    = 32
)

type ComponentPool[T any] struct {
	pages         [][]ComponentIndex // sparse slice of ComponentIndex, which has indices matching the respective entityIndex (split into pages of PAGE_SIZE)
	entityIndexes []EntityIndex      // packed slice of EntityIndex, which has indices matching the respective componentIndex

	components []*T
}

func NewComponentPool[T any]() *ComponentPool[T] {
	return &ComponentPool[T]{
		pages:         make([][]ComponentIndex, 0),
		entityIndexes: make([]EntityIndex, 0),
		components:    make([]*T, 0),
	}
}

func (pool *ComponentPool[T]) HasEntity(entityIndex EntityIndex) bool {
	return len(pool.pages) > int(entityIndex/PAGE_SIZE) && pool.pages[entityIndex/PAGE_SIZE] != nil && pool.pages[entityIndex/PAGE_SIZE][entityIndex%PAGE_SIZE] != INVALID_COMPONENT
}

func (pool *ComponentPool[T]) RemoveEntity(entityIndex EntityIndex) {
	targetComponentIndex := pool.pages[entityIndex/PAGE_SIZE][entityIndex%PAGE_SIZE]
	lastEntityIndex := pool.entityIndexes[len(pool.entityIndexes)-1]
	lastComponent := pool.components[len(pool.components)-1]
	pool.pages[entityIndex/PAGE_SIZE][entityIndex%PAGE_SIZE] = INVALID_COMPONENT
	pool.pages[lastEntityIndex/PAGE_SIZE][lastEntityIndex%PAGE_SIZE] = targetComponentIndex
	pool.entityIndexes[targetComponentIndex] = lastEntityIndex
	pool.components[targetComponentIndex] = lastComponent
	pool.entityIndexes = pool.entityIndexes[:len(pool.entityIndexes)-1]
	pool.components = pool.components[:len(pool.components)-1]
}
