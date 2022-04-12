package ecs

type ComponentPool[T any] struct {
	componentIds []int // sparse slice of ComponentId
	entityIds    []int // packed slice of EntityId

	components []T
}

func NewComponentPool[T any]() *ComponentPool[T] {
	return &ComponentPool[T]{
		componentIds: make([]int, 0),
		entityIds:    make([]int, 0),
		components:   make([]T, 0),
	}
}

func (pool *ComponentPool[T]) HasEntity(entityId int) bool {
	return len(pool.componentIds) > entityId && pool.componentIds[entityId] > -1
}

func (pool *ComponentPool[T]) RemoveEntity(entityId int) {
	targetComponentId := pool.componentIds[entityId]
	lastEntityId := pool.entityIds[len(pool.entityIds)-1]
	lastComponent := pool.components[len(pool.components)-1]
	pool.componentIds[entityId] = -1
	pool.componentIds[lastEntityId] = targetComponentId
	pool.entityIds[targetComponentId] = lastEntityId
	pool.components[targetComponentId] = lastComponent
	pool.entityIds = pool.entityIds[:len(pool.entityIds)-1]
	pool.components = pool.components[:len(pool.components)-1]
}
