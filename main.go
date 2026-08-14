package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
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

func (r *RoundRobin) PICK() string {
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

type SmoothWeightedRR struct {
	backends    []string
	weights     []int
	current     []int
	totalWeight int
}

func (s *SmoothWeightedRR) PICKN(n string) []string {
	N, err := strconv.Atoi(n)
	if err != nil {
		panic("Invalid argument for PICKN. Expected an integer.")
	}
	ans := make([]string, 0, N)
	for range N {
		p := slices.Max(s.current)

	InerLoop:
		for j, key := range s.backends {
			if s.current[j] == p {
				ans = append(ans, key)
				s.current[j] -= s.totalWeight
				break InerLoop
			}
		}
		for j := range len(s.current) {
			s.current[j] += s.weights[j]
		}

	}
	return ans
}

func NewSmoothWeightedRR(backends []string) *SmoothWeightedRR {
	backendList := make([]string, len(backends))
	current := make([]int, len(backends))
	weights := make([]int, len(backends))
	totalWeight := 0
	for i, backend := range backends {
		parts := strings.Split(backend, ":")
		if len(parts) == 2 {
			weight, err := strconv.Atoi(parts[1])
			if err != nil {
				panic("Invalid backend weight. Expected an integer.")
			}
			backendList[i] = parts[0]
			current[i] = weight
			weights[i] = weight
			totalWeight += weight
		} else {
			panic("Invalid backend format. Expected 'backend:weight'")
		}
	}
	return &SmoothWeightedRR{
		backends:    backendList,
		current:     current,
		weights:     weights,
		totalWeight: totalWeight,
	}
}

func main() {
	// file, err := os.Open("tests/02-weighted-rr/4.in")
	// if err != nil {
	// 	panic(err)
	// }
	// defer file.Close()
	// sc := bufio.NewScanner(file)
	sc := bufio.NewScanner(os.Stdin)
	sc.Buffer(make([]byte, 1024*1024), 1024*1024)
	var swrr *SmoothWeightedRR
	for sc.Scan() {
		line := sc.Text()
		if line == "" {
			continue
		}
		args := strings.Split(line, " ")
		switch args[0] {
		case "POOL":
			swrr = NewSmoothWeightedRR(args[1:])
			fmt.Println("OK")
		case "PICKN":
			ans := swrr.PICKN(args[1])
			fmt.Println(strings.Join(ans, ","))

		default:
			fmt.Println("wrong input:", args[0])
		}
	}
}
