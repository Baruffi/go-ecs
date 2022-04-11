package managers

import (
	"sync/atomic"
	"time"

	"example.com/v0/src/queue"
)

type EventCall int
type Event func()

type EventManager struct {
	count        *int64
	callQueue    *queue.ThreadSafeQueue[EventCall, *queue.BasicQueue[EventCall]]
	eventQueue   *queue.ThreadSafeQueue[Event, *queue.BasicQueue[Event]]
	callHandler  *queue.MappingQueueHandler[EventCall, Event]
	eventHandler *queue.GenericQueueHandler[Event]
}

func NewEventManager(maxAllowed int64, limit int, unlimitRate time.Duration) *EventManager {
	var count int64
	schedule := make(chan time.Time, limit)
	for i := 0; i < limit; i++ {
		schedule <- time.Now()
	}
	go func() {
		for t := range time.NewTicker(unlimitRate).C {
			schedule <- t
		}
	}()
	eventHandler := &eventScheduler{
		count:      &count,
		maxAllowed: maxAllowed,
		schedule:   schedule,
	}
	return &EventManager{
		count:        &count,
		callQueue:    queue.NewThreadSafeQueue[EventCall](queue.NewBasicQueue[EventCall]()),
		eventQueue:   queue.NewThreadSafeQueue[Event](queue.NewBasicQueue[Event]()),
		callHandler:  queue.NewMappingQueueHandler[EventCall, Event](eventHandler),
		eventHandler: queue.NewGenericQueueHandler[Event](eventHandler),
	}
}

func (m *EventManager) SetMapping(call EventCall, event Event) {
	m.callHandler.SetMapping(call, event)
}

func (m *EventManager) EnqueueCall(c EventCall) error {
	return m.callQueue.Enqueue(c)
}

func (m *EventManager) EnqueueEvent(e Event) error {
	return m.eventQueue.Enqueue(e)
}

func (m *EventManager) Execute() (callHandlerErr error, eventHandlerErr error) {
	callHandlerErr = m.callHandler.Consume(m.callQueue)
	eventHandlerErr = m.eventHandler.Consume(m.eventQueue)

	return callHandlerErr, eventHandlerErr
}

func (m *EventManager) Executing() bool {
	return atomic.LoadInt64(m.count) > 0
}

func (m *EventManager) GetTaskCount() int64 {
	return atomic.LoadInt64(m.count)
}

type eventScheduler struct {
	maxAllowed int64
	count      *int64
	schedule   chan time.Time
}

func (s *eventScheduler) Handle(e Event) queue.HandlerResult {
	if atomic.LoadInt64(s.count) < s.maxAllowed {
		atomic.AddInt64(s.count, int64(1))
		go func() {
			<-s.schedule
			e()
			atomic.AddInt64(s.count, -1)
		}()
		return queue.DONE
	}
	return queue.NOT_DONE
}
