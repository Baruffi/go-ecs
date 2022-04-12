package queue

import (
	"sync"
)

type ThreadSafeQueue[T any, Q Queue[T]] struct {
	mutex       sync.RWMutex
	unsafeQueue Q
}

func NewThreadSafeQueue[T any, Q Queue[T]](queue Q) *ThreadSafeQueue[T, Q] {
	return &ThreadSafeQueue[T, Q]{
		unsafeQueue: queue,
	}
}

func (q *ThreadSafeQueue[T, Q]) SafeRead(read func(queue Q)) {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	read(q.unsafeQueue)
}

func (q *ThreadSafeQueue[T, Q]) SafeWrite(write func(queue Q)) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	write(q.unsafeQueue)
}

func (q *ThreadSafeQueue[T, Q]) Enqueue(item T) error {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	return q.unsafeQueue.Enqueue(item)
}

func (q *ThreadSafeQueue[T, Q]) Peek() (item T, err error) {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	return q.unsafeQueue.Peek()
}

func (q *ThreadSafeQueue[T, Q]) Rotate(amount int) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.unsafeQueue.Rotate(amount)
}

func (q *ThreadSafeQueue[T, Q]) Dequeue() (item T, err error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	return q.unsafeQueue.Dequeue()
}

func (q *ThreadSafeQueue[T, Q]) HasItems() bool {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	return q.unsafeQueue.HasItems()
}

func (q *ThreadSafeQueue[T, Q]) GetLen() int {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	return q.unsafeQueue.GetLen()
}
