package load_balancer

type CS struct {
	rr            RoundRobin
	state         map[string]string
	stickedClient map[string]string
}

func NewCS(backends []string) *CS {
	state := make(map[string]string)
	stickedClient := make(map[string]string)
	for _, b := range backends {
		state[b] = "UP"
	}
	return &CS{
		rr:            *NewRoundRobin(backends),
		state:         state,
		stickedClient: stickedClient,
	}
}

func (c *CS) ChangeState(backends []string, state string) {
	for _, b := range backends {
		c.state[b] = state
	}
}

func (c *CS) Rest() {
	c.rr.Rest()
}

func (c *CS) Request(args []string) string {
	if len(args) < 2 {
		b := c.rr.Pick()
		if b == "" {
			return "NONE"
		}
		return b + " new"
	} else {
		status, ok := c.state[args[1]]
		if !ok || status == "DOWN" {
			b := c.rr.Pick()
			if b == "" {
				return "NONE"
			}
			return b + " new"
		} else {
			return args[1] + " sticky"
		}
	}
}
