package load_balancer

import "fmt"

type P2C struct {
	m        map[string]int
	backends []string
}

func (p *P2C) Pick(a, b string) string {
	valA, okA := p.m[a]
	valB, okB := p.m[b]

	if okA && okB {
		if valA != valB {
			if valA < valB {
				p.m[a] += 1
				return a
			}
			p.m[b] += 1
			return b
		}
	} else {
		panic("empty pool")
	}
	if a < b {
		p.m[a] += 1
		return a
	}
	p.m[b] += 1
	return b
}

func NewP2C(backends []string) *P2C {
	t := make(map[string]int)
	for _, b := range backends {
		t[b] = 0

	}
	return &P2C{
		backends: backends,
		m:        t,
	}
}

func (p *P2C) Status() {
	for _, key := range p.backends {
		fmt.Printf("%s:%d\n", key, p.m[key])
	}
}
