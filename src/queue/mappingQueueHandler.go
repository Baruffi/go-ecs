package queue

import "fmt"

type MappingQueueHandler[M comparable, H any] struct {
	mappings map[M]H
	handler  Handler[H]
}

func NewMappingQueueHandler[M comparable, H any](handler Handler[H]) *MappingQueueHandler[M, H] {
	return &MappingQueueHandler[M, H]{
		mappings: make(map[M]H),
		handler:  handler,
	}
}

func (h *MappingQueueHandler[M, H]) SetMapping(mapping M, handled H) {
	h.mappings[mapping] = handled
}

func (h *MappingQueueHandler[M, H]) UnsetMapping(mapping M) {
	delete(h.mappings, mapping)
}

func (h *MappingQueueHandler[M, H]) Consume(q Queue[M]) error {
	len := q.GetLen()
	for i := 0; i < len; i++ {
		if mapping, err := q.Peek(); err == nil {
			if handled, ok := h.mappings[mapping]; ok {
				result := h.handler.Handle(handled)

				if result == DONE {
					q.Dequeue()
				} else {
					q.Rotate(1)
				}
			} else {
				return fmt.Errorf("invalid mapping - no match found for %v", mapping)
			}
		} else {
			return err
		}
	}

	return nil
}
