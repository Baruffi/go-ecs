package queue

import "fmt"

type BasicQueue[T any] struct {
	head  int
	items []T
}

func NewBasicQueue[T any]() *BasicQueue[T] {
	return &BasicQueue[T]{
		items: make([]T, 0),
	}
}

func (q *BasicQueue[T]) Enqueue(item T) error {
	q.items = append(q.items, item)

	return nil
}

func (q *BasicQueue[T]) Peek() (item T, err error) {
	if len(q.items) > 0 {
		item = q.items[0]
	} else {
		err = fmt.Errorf("no items in queue")
	}

	return item, err
}

func (q *BasicQueue[T]) Rotate(amount int) {
	if len(q.items) > 0 {
		q.head = (q.head + amount + len(q.items)) % len(q.items)
	}
}

func (q *BasicQueue[T]) Dequeue() (item T, err error) {
	if len(q.items) > 0 {
		item = q.items[0]
		q.items = append(q.items[:q.head], q.items[q.head+1:]...)
	} else {
		err = fmt.Errorf("no items in queue")
	}

	return item, err
}

func (q *BasicQueue[T]) HasItems() bool {
	return len(q.items) > 0
}

func (q *BasicQueue[T]) GetLen() int {
	return len(q.items)
}
