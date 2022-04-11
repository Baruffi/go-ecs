package queue

import (
	"fmt"
	"sort"
)

const (
	NO_PRIORITY PriorityLevel = -1
)

type PriorityLevel int

type PriorityQueue[T any] struct {
	queueHead    int
	bucketHead   int
	bucketOrder  []int
	buckets      map[PriorityLevel][]T
	enqueueLevel PriorityLevel
}

func NewPriorityQueue[T any]() *PriorityQueue[T] {
	return &PriorityQueue[T]{
		bucketOrder:  make([]int, 0),
		buckets:      make(map[PriorityLevel][]T, 0),
		enqueueLevel: NO_PRIORITY,
	}
}

func (q *PriorityQueue[T]) SetEnqueueLevel(level PriorityLevel) {
	q.enqueueLevel = level
}

func (q *PriorityQueue[T]) Enqueue(item T) error {
	if q.enqueueLevel == NO_PRIORITY {
		return fmt.Errorf("no priority level set")
	}

	if items, ok := q.buckets[q.enqueueLevel]; ok {
		q.buckets[q.enqueueLevel] = append(items, item)
	} else {
		q.buckets[q.enqueueLevel] = []T{item}
		q.bucketOrder = append(q.bucketOrder, int(q.enqueueLevel))
		sort.Ints(q.bucketOrder)
	}

	return nil
}

func (q *PriorityQueue[T]) Peek() (item T, err error) {
	if len(q.buckets) > 0 {
		item = q.buckets[PriorityLevel(q.bucketOrder[q.queueHead])][q.bucketHead]
	} else {
		err = fmt.Errorf("no items in queue")
	}

	return item, err
}

func (q *PriorityQueue[T]) Rotate(amount int) {
	if len(q.buckets) > 0 {
		var totalLen, prevQueueHead int
		bucketLengths := make(map[PriorityLevel]int, len(q.buckets))
		for level, items := range q.buckets {
			totalLen += len(items)
			bucketLengths[level] = len(items)
			if PriorityLevel(q.bucketOrder[q.queueHead]) > level {
				prevQueueHead += len(items)
			}
		}
		head := (prevQueueHead + q.bucketHead + amount + totalLen) % totalLen
		for order, level := range q.bucketOrder {
			curLength := bucketLengths[PriorityLevel(level)]
			if head < curLength {
				q.queueHead = order
				q.bucketHead = head
				break
			}
			head -= curLength
		}
	}
}

func (q *PriorityQueue[T]) Dequeue() (item T, err error) {
	if len(q.buckets) > 0 {
		head := PriorityLevel(q.bucketOrder[q.queueHead])
		item = q.buckets[head][q.bucketHead]
		if len(q.buckets[head]) == 1 {
			delete(q.buckets, head)
			q.bucketOrder = append(q.bucketOrder[:q.queueHead], q.bucketOrder[q.queueHead+1:]...)
			if q.queueHead == len(q.buckets) {
				q.queueHead = 0
			}
		} else {
			q.buckets[head] = append(q.buckets[head][:q.bucketHead], q.buckets[head][q.bucketHead+1:]...)
			if q.bucketHead == len(q.buckets[head]) {
				q.bucketHead = 0
			}
		}
	} else {
		err = fmt.Errorf("no items in queue")
	}

	return item, err
}

func (q *PriorityQueue[T]) HasItems() bool {
	return len(q.buckets) > 0
}

func (q *PriorityQueue[T]) GetLen() int {
	var totalLen int
	for _, items := range q.buckets {
		totalLen += len(items)
	}
	return totalLen
}
