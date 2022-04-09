package ecs

type PriorityQueueManager[T2 any] struct {
	QueueManager[PriorityLevel, PriorityPackage[T2], *BasicQueue[PriorityLevel], *PriorityQueue[T2]]
}

func NewPriorityQueueManager[T2 any](handler Handler[PriorityPackage[T2]]) PriorityQueueManager[T2] {
	priorityManager := PriorityQueueManager[T2]{NewQueueManager[PriorityLevel](handler, &BasicQueue[PriorityLevel]{}, &PriorityQueue[T2]{})}

	// Priority manager always runs through all defaults on execute, in reverse level order (priority level 0 is therefore highest)
	for i := 9; i >= 0; i-- {
		priorityManager.InsertLevel(PriorityLevel(i))
	}

	return priorityManager
}

func (m *PriorityQueueManager[T2]) SetDefault(level PriorityLevel, elements ...T2) {
	m.defaults[level] = PriorityPackage[T2]{level, elements}
}

func (m *PriorityQueueManager[T2]) Execute() {
	// Priority manager always runs through all defaults on execute, in reverse level order (priority level 0 is therefore highest)
	for _, item := range m.queue1.View() {
		if defaultItem, ok := m.defaults[item]; ok {
			m.Handle(defaultItem)
		}
	}
	for _, item := range m.queue2.Consume() {
		m.Handle(item)
	}
}

func (m *PriorityQueueManager[T2]) InsertLevel(level PriorityLevel) {
	m.queue1.Feed(level)
}

func (m *PriorityQueueManager[T2]) InsertElement(level PriorityLevel, elements ...T2) {
	m.queue2.Insert(level, elements...)
}

func (m *PriorityQueueManager[T2]) AddDefault(level PriorityLevel, elements ...T2) {
	if curDefault, ok := m.defaults[level]; ok {
		curDefault.Element = append(curDefault.Element, elements...)
		m.defaults[level] = curDefault
	} else {
		m.defaults[level] = PriorityPackage[T2]{level, elements}
	}
}
