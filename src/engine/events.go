package engine

type EventData interface {
}

type EventCall interface {
	~uint
}

type Events[E EventCall] map[E][]EventData

type EventEmitter[E EventCall] interface {
	Emit() Events[E]
}

type EventConsumer[E EventCall] interface {
	Consume(Events[E])
}

func EmitEvents[E EventCall](emitters []EventEmitter[E]) Events[E] {
	events := make(Events[E])
	for _, e := range emitters {
		for k, v := range e.Emit() {
			events[k] = append(events[k], v...)
		}
	}
	return events
}

func ConsumeEvents[E EventCall](consumers []EventConsumer[E], events Events[E]) {
	for _, c := range consumers {
		c.Consume(events)
	}
}
