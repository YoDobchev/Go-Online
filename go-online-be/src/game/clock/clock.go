package clock

import (
	"sync"
	"time"
)

type TimeFormat struct {
	MainTime       time.Duration
	ByoYomi        time.Duration
	ByoYomiPeriods int
}

type PlayerTime struct {
	MainRemaining  time.Duration
	ByoRemaining   time.Duration
	ByoPeriodsLeft int
	InByoYomi      bool
}

type Clock struct {
	mu        sync.Mutex
	format    TimeFormat
	players   map[uint8]*PlayerTime
	lastStart time.Time
	running   bool
	current   uint8
}

func NewClock(format TimeFormat) *Clock {
	return &Clock{
		format: format,
		players: map[uint8]*PlayerTime{
			1: {
				MainRemaining:  format.MainTime,
				ByoRemaining:   format.ByoYomi,
				ByoPeriodsLeft: format.ByoYomiPeriods,
			},
			2: {
				MainRemaining:  format.MainTime,
				ByoRemaining:   format.ByoYomi,
				ByoPeriodsLeft: format.ByoYomiPeriods,
			},
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
	p := c.players[c.current]

	if !p.InByoYomi {
		if elapsed < p.MainRemaining {
			p.MainRemaining -= elapsed
		} else {
			elapsed -= p.MainRemaining
			p.MainRemaining = 0
			p.InByoYomi = true
		}
	}

	if p.InByoYomi {
		for elapsed > 0 {
			if p.ByoRemaining > elapsed {
				p.ByoRemaining -= elapsed
				elapsed = 0
			} else {
				elapsed -= p.ByoRemaining
				p.ByoPeriodsLeft--
				if p.ByoPeriodsLeft <= 0 {
					p.ByoRemaining = 0
					c.running = false
					return
				}
				p.ByoRemaining = c.format.ByoYomi
			}
		}
	}

	c.running = false
}

func (c *Clock) Switch(next uint8) {
	c.Stop()
	c.Start(next)
}

func (c *Clock) OutOfTime(p uint8) bool {
	pt := c.players[p]

	if pt.MainRemaining > 0 {
		if c.running && p == c.current {
			return pt.MainRemaining-time.Since(c.lastStart) <= 0 &&
				pt.ByoPeriodsLeft == 0
		}
		return false
	}

	if pt.ByoPeriodsLeft <= 0 {
		return true
	}

	if c.running && p == c.current {
		return pt.ByoRemaining-time.Since(c.lastStart) <= 0 &&
			pt.ByoPeriodsLeft == 1
	}

	return false
}

func (c *Clock) GetClockUpdate(player uint8) map[string]any {
	c.mu.Lock()
	defer c.mu.Unlock()
	playerTime := c.players[player]

	var main time.Duration
	var byoYomi time.Duration

	if playerTime.MainRemaining <= time.Since(c.lastStart) {
		difference := time.Since(c.lastStart) - playerTime.MainRemaining
		main = 0
		byoYomi = playerTime.ByoRemaining - difference
		if byoYomi < 0 {
			byoYomi = 0
		}
	} else {
		main = playerTime.MainRemaining - time.Since(c.lastStart)
		byoYomi = playerTime.ByoRemaining
	}

	return map[string]any{
		"type": "clock_update",
		"data": map[string]any{
			"player":            c.current,
			"main_remaining_ms": main.Milliseconds(),
			"byo_remaining_ms":  byoYomi.Milliseconds(),
			"byo_periods_left":  playerTime.ByoPeriodsLeft,
			"in_byo_yomi":       playerTime.InByoYomi,
			"server_time":       time.Now().UnixMilli(),
		},
	}
}
