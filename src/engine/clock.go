package engine

import "time"

type Clock struct {
	dt   float64
	last time.Time
}

func (t *Clock) Init() {
	t.last = time.Now()
}

func (t *Clock) Tick() {
	t.dt = time.Since(t.last).Seconds()
	t.last = time.Now()
}
