package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

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


func main() {
	// file, err := os.Open("tests/04-power-of-two/4.in")
	// if err != nil {
	// 	panic(err)
	// }
	// defer file.Close()
	// sc := bufio.NewScanner(file)
	sc := bufio.NewScanner(os.Stdin)
	sc.Buffer(make([]byte, 1024*1024), 1024*1024)
	var lb *P2C
	for sc.Scan() {
		if sc.Err() != nil{
			panic("error in scaning")
		}
		line := sc.Text()
		if line == "" {
			continue
		}
		args := strings.Split(line, " ")
		switch args[0] {
		case "POOL":
			lb = NewP2C(args[1:])
			fmt.Println("OK")
		case "PICK":
			ans := lb.Pick(args[1], args[2])
			fmt.Println(ans)
		case "STATUS":
			lb.Status()
		default:
			fmt.Println("wrong input:", args[0])
		}
	}
}
