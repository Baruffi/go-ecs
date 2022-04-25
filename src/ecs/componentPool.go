package ecs

import "math"

type ComponentIndex uint32

const (
	INVALID_COMPONENT ComponentIndex = ComponentIndex(math.MaxUint32)
	PAGE_SIZE         EntityIndex    = 32
)

type ComponentPool[T anyComponent] struct {
	pages         [][]ComponentIndex // sparse slice of ComponentIndex, which has indices matching the respective entityIndex (split into pages of PAGE_SIZE)
	entityIndexes []EntityIndex      // packed slice of EntityIndex, which has indices matching the respective componentIndex

	components []*T
}

func NewComponentPool[T anyComponent]() *ComponentPool[T] {
	return &ComponentPool[T]{
		pages:         make([][]ComponentIndex, 0),
		entityIndexes: make([]EntityIndex, 0),
		components:    make([]*T, 0),
	}
}

func (pool *ComponentPool[T]) Assign(entityIndex EntityIndex, component anyComponent) {
	entityPage := entityIndex / PAGE_SIZE
	pool.entityIndexes = append(pool.entityIndexes, entityIndex)
	pool.components = append(pool.components, component.(*T))
	for EntityIndex(len(pool.pages)) <= entityPage {
		pool.pages = append(pool.pages, nil)
	}
	if pool.pages[entityPage] == nil {
		newPage := make([]ComponentIndex, PAGE_SIZE)
		for i := range newPage {
			newPage[i] = INVALID_COMPONENT
		}
		pool.pages[entityPage] = newPage
	}
	pool.pages[entityPage][entityIndex%PAGE_SIZE] = ComponentIndex(len(pool.components) - 1)
}

func (pool *ComponentPool[T]) GetComponent(entityIndex EntityIndex) anyComponent {
	return pool.components[pool.pages[entityIndex/PAGE_SIZE][entityIndex%PAGE_SIZE]]
}

func (pool *ComponentPool[T]) HasEntity(entityIndex EntityIndex) bool {
	return EntityIndex(len(pool.pages)) > entityIndex/PAGE_SIZE && pool.pages[entityIndex/PAGE_SIZE] != nil && pool.pages[entityIndex/PAGE_SIZE][entityIndex%PAGE_SIZE] != INVALID_COMPONENT
}

func (pool *ComponentPool[T]) RemoveEntity(entityIndex EntityIndex) {
	targetComponentIndex := pool.pages[entityIndex/PAGE_SIZE][entityIndex%PAGE_SIZE]
	lastEntityIndex := pool.entityIndexes[len(pool.entityIndexes)-1]
	lastComponent := pool.components[len(pool.components)-1]
	pool.pages[lastEntityIndex/PAGE_SIZE][lastEntityIndex%PAGE_SIZE] = targetComponentIndex
	pool.pages[entityIndex/PAGE_SIZE][entityIndex%PAGE_SIZE] = INVALID_COMPONENT
	pool.entityIndexes[targetComponentIndex] = lastEntityIndex
	pool.components[targetComponentIndex] = lastComponent
	pool.entityIndexes = pool.entityIndexes[:len(pool.entityIndexes)-1]
	pool.components = pool.components[:len(pool.components)-1]
}
