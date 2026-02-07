package clock

import (
	"time"
)

type Clock struct {
	remaining map[uint8]time.Duration
	lastStart time.Time
	running   bool
	current   uint8
}

func NewClock(mainTime time.Duration) *Clock {
	return &Clock{
		remaining: map[uint8]time.Duration{
			1: mainTime,
			2: mainTime,
		},
	}
}

func (c *Clock) Start(player uint8) {
	c.current = player
	c.lastStart = time.Now()
	c.running = true
}

func (c *Clock) Stop() {
	if !c.running {
		return
	}
	elapsed := time.Since(c.lastStart)
	c.remaining[c.current] -= elapsed
	c.running = false
}

func (c *Clock) Switch(next uint8) {
	c.Stop()
	c.Start(next)
}

func (c *Clock) OutOfTime(p uint8) bool {
	if c.running && p == c.current {
		elapsed := time.Since(c.lastStart)
		return c.remaining[p]-elapsed <= 0
	}
	return c.remaining[p] <= 0
}
