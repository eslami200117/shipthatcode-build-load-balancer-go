package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

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
		for range c.state {
			b := c.rr.Pick()
			if b == "" {
				return "NONE"
			}
			if c.state[b] == "UP" {
				return b + " new"
			}
		}

		return "NONE"
	} else {
		status, ok := c.state[args[1]]
		if !ok || status == "DOWN" {
			for range c.state {
				b := c.rr.Pick()
				if b == "" {
					return "NONE"
				}
				if c.state[b] == "UP" {
					return b + " new"
				}
			}
			return "NONE"
		} else {
			return args[1] + " sticky"
		}
	}
}

func main() {
	// var test int
	// var err error
	// if len(os.Args) < 2 {
	// 	test = 1
	// } else {
	// 	test, err = strconv.Atoi(os.Args[1])
	// 	if err != nil {
	// 		panic("wrong arguman")
	// 	}
	// }

	// file, err := os.Open(fmt.Sprintf("tests/09-sticky-sessions/%d.in", test))
	// if err != nil {
	// 	panic(err)
	// }
	// defer file.Close()
	// sc := bufio.NewScanner(file)
	sc := bufio.NewScanner(os.Stdin)
	sc.Buffer(make([]byte, 1024*1024), 1024*1024)
	var lb *CS
	for sc.Scan() {
		if sc.Err() != nil {
			panic("error in scaning")
		}
		line := sc.Text()
		if line == "" {
			continue
		}
		args := strings.Split(line, " ")
		switch args[0] {
		case "POOL":
			lb = NewCS(args[1:])
			fmt.Println("OK")
		case "REQUEST":
			ans := lb.Request(args[1:])
			fmt.Println(ans)
		case "UP":
			lb.ChangeState(args[1:], "UP")
			fmt.Println("OK")
		case "DOWN":
			lb.ChangeState(args[1:], "DOWN")
			fmt.Println("OK")
		default:
			fmt.Println("wrong input:", args[0])
		}
	}
}
