package tools

import "time"

type Clock struct {
	last      time.Time
	frameTick *time.Ticker
}

func (c *Clock) Init(fps int) {
	c.last = time.Now()

	if fps <= 0 {
		c.frameTick = nil
	} else {
		c.frameTick = time.NewTicker(time.Second / time.Duration(fps))
	}
}

func (c *Clock) Dt() float64 {
	c.last = time.Now()

	return time.Since(c.last).Seconds()
}

func (c *Clock) Tick() {
	if c.frameTick != nil {
		<-c.frameTick.C
	}
}
