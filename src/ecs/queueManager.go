package ecs

type Handler[T any] interface {
	Handle(T)
}

type HandlerFunc[T any] func(T)

func (f HandlerFunc[T]) Handle(t T) {
	f(t)
}

type QueueManager[T1 comparable, T2 any, Q1 Queue[T1], Q2 Queue[T2]] struct {
	Handler[T2]
	defaults map[T1]T2
	queue1   Q1
	queue2   Q2
}

func NewQueueManager[T1 comparable, T2 any, Q1 Queue[T1], Q2 Queue[T2]](handler Handler[T2], queue1 Q1, queue2 Q2) QueueManager[T1, T2, Q1, Q2] {
	return QueueManager[T1, T2, Q1, Q2]{
		Handler:  handler,
		defaults: make(map[T1]T2),
		queue1:   queue1,
		queue2:   queue2,
	}
}

func (m *QueueManager[T1, T2, Q1, Q2]) SetDefault(t1 T1, t2 T2) {
	m.defaults[t1] = t2
}

func (m *QueueManager[T1, T2, Q1, Q2]) AddT1(c T1) {
	m.queue1.Feed(c)
}

func (m *QueueManager[T1, T2, Q1, Q2]) AddT2(e T2) {
	m.queue2.Feed(e)
}

func (m *QueueManager[T1, T2, Q1, Q2]) Execute() {
	for _, item := range m.queue1.Consume() {
		if defaultItem, ok := m.defaults[item]; ok {
			m.Handle(defaultItem)
		}
	}
	for _, item := range m.queue2.Consume() {
		m.Handle(item)
	}
}

type BasicQueueManager[T1 comparable, T2 any] struct {
	QueueManager[T1, T2, *BasicQueue[T1], *BasicQueue[T2]]
}

func NewBasicQueueManager[T1 comparable, T2 any](handler Handler[T2]) BasicQueueManager[T1, T2] {
	return BasicQueueManager[T1, T2]{NewQueueManager[T1](handler, &BasicQueue[T1]{}, &BasicQueue[T2]{})}
}
