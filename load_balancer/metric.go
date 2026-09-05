package load_balancer

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type Metric struct {
	pool      []string
	requests  map[string]map[string]int
	bytesIn   map[string]int
	bytesOut  map[string]int
	latencies map[string][]int
	health    map[string]int
}

func NewMetric() *Metric {
	return &Metric{
		pool:      []string{},
		requests:  make(map[string]map[string]int),
		bytesIn:   make(map[string]int),
		bytesOut:  make(map[string]int),
		latencies: make(map[string][]int),
		health:    make(map[string]int),
	}
}

func (m *Metric) Pool(backends ...string) {
	m.pool = backends
	for _, b := range m.pool {
		m.health[b] = 1
		m.bytesIn[b] = 0
		m.bytesOut[b] = 0
		m.requests[b] = make(map[string]int)
		m.latencies[b] = []int{}
	}
}

func (m *Metric) Request(backend, status string, durationMs string) {
	if _, ok := m.requests[backend]; !ok {
		m.requests[backend] = make(map[string]int)
	}
	m.requests[backend][status]++
	d, err := strconv.Atoi(durationMs)
	if err != nil {
		panic("invalid input")
	}
	m.latencies[backend] = append(m.latencies[backend], d)
}

func (m *Metric) Bytes(backend, direction string, b string) {
	n, err := strconv.Atoi(b)
	if err != nil {
		panic("invalid input")
	}
	if direction == "in" {
		m.bytesIn[backend] += n
	} else {
		m.bytesOut[backend] += n
	}
}

func (m *Metric) Health(backend, status string) {
	if status == "UP" {
		m.health[backend] = 1
	} else {
		m.health[backend] = 0
	}
}

func (m *Metric) Metrics() string {
	var output strings.Builder

	sort.Strings(m.pool)

	// Requests - only for pairs that have occurred
	var reqKeys []string
	for b, statuses := range m.requests {
		for s := range statuses {
			reqKeys = append(reqKeys, b+"|"+s)
		}
	}
	sort.Strings(reqKeys)
	for _, key := range reqKeys {
		parts := strings.Split(key, "|")
		b := parts[0]
		s := parts[1]
		output.WriteString(fmt.Sprintf("lb_requests_total{backend=\"%s\",status=\"%s\"} %d\n", b, s, m.requests[b][s]))
	}

	// Bytes In - all backends sorted
	for _, b := range m.pool {
		output.WriteString(fmt.Sprintf("lb_bytes_in_total{backend=\"%s\"} %d\n", b, m.bytesIn[b]))
	}

	// Bytes Out - all backends sorted
	for _, b := range m.pool {
		output.WriteString(fmt.Sprintf("lb_bytes_out_total{backend=\"%s\"} %d\n", b, m.bytesOut[b]))
	}

	// Latency - all backends sorted
	for _, b := range m.pool {
		avg := 0
		if len(m.latencies[b]) > 0 {
			sum := 0
			for _, d := range m.latencies[b] {
				sum += d
			}
			avg = sum / len(m.latencies[b]) // Truncates toward zero
		}
		output.WriteString(fmt.Sprintf("lb_latency_avg_ms{backend=\"%s\"} %d\n", b, avg))
	}

	// Health - all backends sorted
	for _, b := range m.pool {
		output.WriteString(fmt.Sprintf("lb_health{backend=\"%s\"} %d\n", b, m.health[b]))
	}

	return output.String()
}
