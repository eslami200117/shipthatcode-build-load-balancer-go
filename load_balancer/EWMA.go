package load_balancer

import (
	"fmt"
	"math"
	"sort"
	"strconv"
)

type EWMA struct {
	backends []string
	state    map[string]float64
}

func (e *EWMA) Len() int {
	return len(e.backends)
}

func (e *EWMA) Less(i, j int) bool {
	if e.state[e.backends[i]] == e.state[e.backends[j]] {
		return e.backends[i] < e.backends[j]
	}
	return e.state[e.backends[i]] < e.state[e.backends[j]]
}

func (e *EWMA) Swap(i, j int) {
	e.backends[i], e.backends[j] = e.backends[j], e.backends[i]
}

func NewEWMA(backends []string) *EWMA {
	state := make(map[string]float64)
	for _, b := range backends {
		state[b] = 0
	}
	return &EWMA{
		backends: backends,
		state:    state,
	}
}

func (e *EWMA) Record(b string, latency string) {
	l, err := strconv.ParseFloat(latency, 64)
	if err != nil {
		panic("invalid argument")
	}
	e.state[b] = 0.3*l + 0.7*e.state[b]
}

func (e *EWMA) Pick() string {
	sort.Sort(e)
	return e.backends[0]
}

func (e *EWMA) Status() {
	sort.Strings(e.backends)
	for _, b := range e.backends {
		fmt.Printf("%s:%d\n", b, int(math.Round(e.state[b])))
	}
}
