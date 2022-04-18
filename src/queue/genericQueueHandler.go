package queue

type GenericQueueHandler[Q any] struct {
	handler Handler[Q]
}

func NewGenericQueueHandler[Q any](handler Handler[Q]) *GenericQueueHandler[Q] {
	return &GenericQueueHandler[Q]{
		handler: handler,
	}
}

func (h *GenericQueueHandler[Q]) Consume(q Queue[Q]) error {
	q.Reset()
	len := q.GetLen()
	for i := 0; i < len; i++ {
		if item, err := q.Peek(); err == nil {
			result := h.handler.Handle(item)
			if result == DONE {
				q.Dequeue()
			} else {
				q.Rotate(1)
			}
		} else {
			return err
		}
	}
	return nil
}
