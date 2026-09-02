package load_balancer

import (
	"fmt"
	"strconv"
)

type CBStatus string

const (
	OPEN      CBStatus = "OPEN"
	HALF_OPEN CBStatus = "HALF_OPEN"
	CLOSE     CBStatus = "CLOSE"
)

type CBInfo struct {
	lastUpdate  int
	status      CBStatus
	succ_streak int
	fail_streak int
}
type CicuitBreaker struct {
	backends []string
	states   map[string]*CBInfo
}

var now int

func (c *CicuitBreaker) Now(t string) {
	now, _ = strconv.Atoi(t)
	for _, state := range c.states {
		if now-state.lastUpdate >= 30 && state.status == OPEN {
			state.status = HALF_OPEN
		}
	}
}
func (c *CicuitBreaker) Call(b string, ok bool) string {
	state := c.states[b]
	state.lastUpdate = now

	if !ok {
		switch state.status {
		case CLOSE:
			state.fail_streak += 1
			if state.fail_streak >= 5 {
				state.status = OPEN
			}
			return "FAIL"
		case HALF_OPEN:
			state.status = OPEN
			state.succ_streak = 0
			return "FAIL"
		case OPEN:
			return "SHORT"
		default:
			return "SHORT"
		}
	} else {
		switch state.status {
		case CLOSE:
			state.fail_streak = 0
			return "OK"
		case HALF_OPEN:
			state.succ_streak += 1
			if state.succ_streak >= 3 {
				state.status = CLOSE
			}
			return "OK"
		case OPEN:
			return "OK"
		default:
			return "OK"
		}
	}
}

func (c *CicuitBreaker) Status() {
	for _, b := range c.backends {
		fmt.Printf("%s %s\n", b, c.states[b].status)
	}
}

func NewCB(backends []string) *CicuitBreaker {
	states := make(map[string]*CBInfo)
	for _, b := range backends {
		states[b] = &CBInfo{
			status: CLOSE,
		}
	}
	return &CicuitBreaker{
		backends: backends,
		states:   states,
	}
}
