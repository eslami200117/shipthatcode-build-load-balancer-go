package main

import (
	"strconv"
	"strings"
)


type SmoothWeightedRR struct {
	backends    []string
	weights     []int
	current     []int
	totalWeight int
}

func (s *SmoothWeightedRR) Pick() string {
	ans := s.PickN("1")

	return ans[0]
}

func (s *SmoothWeightedRR) PickN(n string) []string {
	N, err := strconv.Atoi(n)
	if err != nil {
		panic("Invalid argument for PickN. Expected an integer.")
	}
	ans := make([]string, 0, N)
	for i := 0; i < N; i++ {
		p := Max(s.current)

	InerLoop:
		for j, key := range s.backends {
			if s.current[j] == p {
				ans = append(ans, key)
				s.current[j] -= s.totalWeight
				break InerLoop
			}
		}
		for j := 0; j < len(s.current); j++ {
			s.current[j] += s.weights[j]
		}

	}
	return ans
}

func Max(in []int) int {
	ans := -1000000000000
	for i := 0; i < len(in); i++ {
		if in[i] > ans {
			ans = in[i]
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


func (s *SmoothWeightedRR) Rest() {
	copy(s.current, s.weights)
}