package managers

import "example.com/v0/src/ecs"

type EventCall int
type Event func()

type EventManager ecs.BasicQueueManager[EventCall, Event]

func NewEventManager() EventManager {
	return EventManager(ecs.NewBasicQueueManager[EventCall, Event](ecs.HandlerFunc[Event](handle)))
}

func handle(e Event) {
	e()
}
