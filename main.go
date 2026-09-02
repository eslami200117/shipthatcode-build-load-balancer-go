package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Healthy struct {
	backends    []string
	succ_streak map[string]int
	fail_streak map[string]int
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

func NewHealthy(backends []string) *Healthy {
	s := make(map[string]string)
	for _, b := range backends {
		s[b] = "UP"
	}
	return &Healthy{
		backends:    backends,
		succ_streak: make(map[string]int),
		fail_streak: make(map[string]int),
		state:       s,
	}
}

func main() {
	// file, err := os.Open("tests/05-health-checks/3.in")
	// if err != nil {
	// 	panic(err)
	// }
	// defer file.Close()
	// sc := bufio.NewScanner(file)
	sc := bufio.NewScanner(os.Stdin)
	sc.Buffer(make([]byte, 1024*1024), 1024*1024)
	var lb *Healthy
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
			lb = NewHealthy(args[1:])
			fmt.Println("OK")
		case "REPORT":
			ok := args[2] == "OK"
			lb.Report(args[1], ok)
			fmt.Println("OK")
		case "HEALTHY":
			lb.Healthy()
		case "STATUS":
			lb.Status()
		default:
			fmt.Println("wrong input:", args[0])
		}
	}
}
