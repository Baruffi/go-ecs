package ecs

import "time"

type Clock struct {
	last time.Time
}

func (c *Clock) Init() {
	c.last = time.Now()
}

func (c *Clock) Dt() float64 {
	c.last = time.Now()

	return time.Since(c.last).Seconds()
}
