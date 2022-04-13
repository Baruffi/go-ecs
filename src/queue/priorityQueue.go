package queue

import (
	"fmt"
)

const (
	ZERO PriorityLevel = 9 - iota
	ONE
	TWO
	THREE
	FOUR
	FIVE
	SIX
	SEVEN
	EIGHT
	NINE
	NO_PRIORITY
)

type PriorityLevel int

type PriorityQueue[T any] struct {
	queueHead    int
	bucketHead   int
	totalLen     int
	buckets      [10][]T
	enqueueLevel PriorityLevel
}

func NewPriorityQueue[T any]() *PriorityQueue[T] {
	pq := &PriorityQueue[T]{
		enqueueLevel: NO_PRIORITY,
	}
	for i := 0; i < 10; i++ {
		pq.buckets[i] = make([]T, 0)
	}
	return pq
}

func (q *PriorityQueue[T]) SetEnqueueLevel(level PriorityLevel) {
	q.enqueueLevel = level
}

func (q *PriorityQueue[T]) Enqueue(item T) error {
	if q.enqueueLevel == NO_PRIORITY {
		return fmt.Errorf("no priority level set")
	}

	q.buckets[q.enqueueLevel] = append(q.buckets[q.enqueueLevel], item)
	q.totalLen++

	return nil
}

func (q *PriorityQueue[T]) Peek() (item T, err error) {
	if q.totalLen > 0 {
		for len(q.buckets[q.queueHead]) == 0 {
			q.queueHead = (q.queueHead + 1) % 10
		}
		item = q.buckets[q.queueHead][q.bucketHead]
	} else {
		err = fmt.Errorf("no items in queue")
	}

	return item, err
}

func (q *PriorityQueue[T]) Rotate(amount int) {
	if q.totalLen > 0 {
		head := (q.bucketHead + amount) % q.totalLen
		if head >= 0 {
			for level, bucket := range q.buckets[q.queueHead:] {
				bucketLength := len(bucket)
				if head < bucketLength {
					q.queueHead += level
					q.bucketHead = head
					return
				}
				head -= bucketLength
			}
			for level, bucket := range q.buckets[:q.queueHead] {
				bucketLength := len(bucket)
				if head < bucketLength {
					q.queueHead = level
					q.bucketHead = head
					return
				}
				head -= bucketLength
			}
		}
		for i := q.queueHead - 1; i >= 0; i-- {
			head += len(q.buckets[i])
			if head > 0 {
				q.queueHead = i
				q.bucketHead = head
				return
			}
		}
		for i := q.totalLen; i >= q.queueHead; i-- {
			head += len(q.buckets[i])
			if head > 0 {
				q.queueHead = i
				q.bucketHead = head
				return
			}
		}
	}
}

func (q *PriorityQueue[T]) Dequeue() (item T, err error) {
	if q.totalLen > 0 {
		for len(q.buckets[q.queueHead]) == 0 {
			q.queueHead = (q.queueHead + 1) % 10
		}
		item = q.buckets[q.queueHead][q.bucketHead]
		q.buckets[q.queueHead] = append(q.buckets[q.queueHead][:q.bucketHead], q.buckets[q.queueHead][q.bucketHead+1:]...)
		if q.bucketHead >= len(q.buckets[q.queueHead]) {
			q.bucketHead = 0
		}
		q.totalLen--
	} else {
		err = fmt.Errorf("no items in queue")
	}

	return item, err
}

func (q *PriorityQueue[T]) HasItems() bool {
	return q.totalLen > 0
}

func (q *PriorityQueue[T]) GetLen() int {
	return q.totalLen
}
