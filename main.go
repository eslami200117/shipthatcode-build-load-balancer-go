package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type RoundRobin struct {
	backands []string
	index	int
}

func NewRoundRobin(backends []string) *RoundRobin {
	return &RoundRobin{
		backands: backends,
		index: 0,
	}
}

func (r *RoundRobin) PICK() string{
	if len(r.backands) == 0 {
		return ""
	}
	r.index = r.index % len(r.backands)
	backend := r.backands[r.index]
	r.index += 1
	return backend
}

func (r *RoundRobin) REST() {
	r.index = 0
}




func main() {
	sc := bufio.NewScanner(os.Stdin)
	sc.Buffer(make([]byte, 1024*1024), 1024*1024)
	var rr *RoundRobin
	for sc.Scan() {
		line := sc.Text()
		if line == "" { continue }
		args := strings.Split(line, " ")
		switch args[0] {
		case "POOL":
			rr = NewRoundRobin(args[1:])
			fmt.Println("OK")
		case "PICK":
			ans := rr.PICK()
			fmt.Println(ans)
		case "RESET":
			rr.REST()
			fmt.Println("OK")
		default:
			fmt.Println("wrong input:", args[0])
		}
	}
}
