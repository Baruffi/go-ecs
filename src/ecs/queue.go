package ecs

type Queue[T any] interface {
	Feed(T)
	View() []T
	Consume() []T
	Len() int
}

type BasicQueue[T any] struct {
	items []T
}

func (q *BasicQueue[T]) Feed(item T) {
	if q.items == nil {
		q.items = make([]T, 0)
	}
	q.items = append(q.items, item)
}

func (q *BasicQueue[T]) View() []T {
	return q.items
}

func (q *BasicQueue[T]) Consume() []T {
	items := q.items
	q.items = []T{}
	return items
}

func (q *BasicQueue[T]) Len() int {
	return len(q.items)
}
