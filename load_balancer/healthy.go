package load_balancer

import (
	"fmt"
	"strings"
)

type Healthy struct {
	backends    []string
	succ_streak map[string]int
	fail_streak map[string]int
	pass_streak map[string]int
	state       map[string]string
}

func (h *Healthy) Report(backend string, ok bool) {
	if ok {
		h.fail_streak[backend] = 0
		h.succ_streak[backend] += 1
		if h.state[backend] == "DOWN" && h.succ_streak[backend] >= 2 {
			h.state[backend] = "UP"
		}
	} else {
		h.succ_streak[backend] = 0
		h.fail_streak[backend] += 1
		if h.state[backend] == "UP" && h.fail_streak[backend] >= 3 {
			h.state[backend] = "DOWN"
		}
	}
}

func (h *Healthy) Status() {
	for _, b := range h.backends {
		fmt.Println(b, h.state[b])
	}
}

func (h *Healthy) Healthy() {
	ups := make([]string, 0)
	for _, b := range h.backends {
		if h.state[b] == "UP" {
			ups = append(ups, b)
		}
	}
	if len(ups) != 0 {
		fmt.Println(strings.Join(ups, ","))
	} else {
		fmt.Println("none")
	}
}

func (h *Healthy) Observe(backend, status string) {
	ok := true
	if status[0] == '5' || status[0] == '0' {
		ok = false
	}
	if ok {
		h.pass_streak[backend] = 0
		if h.state[backend] == "DOWN" {
			h.state[backend] = "UP"
		}
	} else {
		h.pass_streak[backend] += 1
		if h.state[backend] == "UP" && h.pass_streak[backend] >= 3 {
			h.state[backend] = "DOWN"
		}
	}
}

func NewHealthy(backends []string) *Healthy {
	s := make(map[string]string)
	for _, b := range backends {
		s[b] = "UP"
	}
	return &Healthy{
		backends:    backends,
		succ_streak: make(map[string]int),
		fail_streak: make(map[string]int),
		pass_streak: make(map[string]int),
		state:       s,
	}
}
