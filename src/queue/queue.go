package queue

type Queue[T any] interface {
	Enqueue(T) error
	Dequeue() (T, error)
	Peek() (T, error)
	Rotate(int)
	HasItems() bool
	GetLen() int
}

type QueueConsumer[Q any] interface {
	Consume(Queue[Q]) error
}
