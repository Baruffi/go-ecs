package queue

import (
	"sync"
)

type ThreadSafeQueue[T any, Q Queue[T]] struct {
	mutex     sync.RWMutex
	UnsafePtr Q
}

func NewThreadSafeQueue[T any, Q Queue[T]](queue Q) *ThreadSafeQueue[T, Q] {
	return &ThreadSafeQueue[T, Q]{
		UnsafePtr: queue,
	}
}

func (q *ThreadSafeQueue[T, Q]) Enqueue(item T) error {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	return q.UnsafePtr.Enqueue(item)
}

func (q *ThreadSafeQueue[T, Q]) Peek() (item T, err error) {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	return q.UnsafePtr.Peek()
}

func (q *ThreadSafeQueue[T, Q]) Rotate(amount int) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.UnsafePtr.Rotate(amount)
}

func (q *ThreadSafeQueue[T, Q]) Dequeue() (item T, err error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	return q.UnsafePtr.Dequeue()
}

func (q *ThreadSafeQueue[T, Q]) HasItems() bool {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	return q.UnsafePtr.HasItems()
}

func (q *ThreadSafeQueue[T, Q]) GetLen() int {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	return q.UnsafePtr.GetLen()
}
