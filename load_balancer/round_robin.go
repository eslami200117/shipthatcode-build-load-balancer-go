package load_balancer


type RoundRobin struct {
	backands []string
	index    int
}

func NewRoundRobin(backends []string) *RoundRobin {
	return &RoundRobin{
		backands: backends,
		index:    0,
	}
}

func (r *RoundRobin) Pick() string {
	if len(r.backands) == 0 {
		return ""
	}
	r.index = r.index % len(r.backands)
	backend := r.backands[r.index]
	r.index += 1
	return backend
}

func (r *RoundRobin) Rest() {
	r.index = 0
}