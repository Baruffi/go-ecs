package managers

import (
	"sync/atomic"
	"time"

	"example.com/v0/src/ecs"
)

type EventCall int
type Event func()

type EventManager struct {
	count *int64
	ecs.BasicQueueManager[EventCall, Event]
}

func NewEventManager(maxAllowed int64, limit int, unlimitRate time.Duration) EventManager {
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
	return EventManager{
		count: &count,
		BasicQueueManager: ecs.NewBasicQueueManager[EventCall, Event](
			&eventScheduler{
				count:      &count,
				maxAllowed: maxAllowed,
				schedule:   schedule,
			},
		),
	}
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

func (s *eventScheduler) Handle(e Event) {
	if atomic.LoadInt64(s.count) < s.maxAllowed {
		atomic.AddInt64(s.count, int64(1))
		go func() {
			<-s.schedule
			e()
			atomic.AddInt64(s.count, -1)
		}()
	} else {
		e()
	}
}
