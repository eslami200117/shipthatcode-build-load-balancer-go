package load_balancer

import (
	"fmt"
	"sort"
	"strconv"
)

type ConsistentHashing struct {
	keys []int
	vals map[int]string
}

func NewConsistentHashing(backends []string) *ConsistentHashing {
	ch := &ConsistentHashing{
		vals: make(map[int]string),
		keys: make([]int, 0),
	}
	ch.Add(backends)
	return ch
}

// add with O(n) insertion maintaining sorted order
func (ch *ConsistentHashing) add(key int, value string) {
	if _, exists := ch.vals[key]; exists {
		return
	}

	// Find insertion point using binary search
	idx := sort.Search(len(ch.keys), func(i int) bool {
		return ch.keys[i] >= key
	})

	// Insert at idx
	ch.keys = append(ch.keys, 0)
	copy(ch.keys[idx+1:], ch.keys[idx:])
	ch.keys[idx] = key
	ch.vals[key] = value
}

func (ch *ConsistentHashing) Add(bachends []string) {
	for _, b := range bachends {
		for j := 0; j < 5; j++ {
			p := hashString(b + "#" + strconv.Itoa(j))
			ch.add(p, b)
		}
	}
}

// LookUp finds the smallest key >= input key (successor with wrap-around)
func (ch *ConsistentHashing) LookUp(p string) string {
	if len(ch.keys) == 0 {
		return ""
	}

	key := hashString(p)

	// Find first key >= input
	idx := sort.Search(len(ch.keys), func(i int) bool {
		return ch.keys[i] >= key
	})

	// If no key found, wrap around to first (smallest)
	if idx == len(ch.keys) {
		idx = 0
	}

	return ch.vals[ch.keys[idx]]
}

// Remove a backend (optional but good to have)
func (ch *ConsistentHashing) Remove(b string) {
	for j := 0; j < 5; j++ {
		p := hashString(b + "#" + strconv.Itoa(j))
		if _, exists := ch.vals[p]; exists {
			delete(ch.vals, p)
			// Remove from keys slice
			idx := sort.Search(len(ch.keys), func(i int) bool {
				return ch.keys[i] >= p
			})
			if idx < len(ch.keys) && ch.keys[idx] == p {
				ch.keys = append(ch.keys[:idx], ch.keys[idx+1:]...)
			}
		}
	}
}

// Get all backends (useful for debugging)
func (ch *ConsistentHashing) GetBackends() []string {
	seen := make(map[string]bool)
	backends := make([]string, 0)
	for _, key := range ch.keys {
		if val := ch.vals[key]; !seen[val] {
			seen[val] = true
			backends = append(backends, val)
		}
	}
	return backends
}

func hashString(s string) int {
	h := 0
	for _, c := range s {
		h = (h*31 + int(c)) % 1000
	}
	return h
}

func (ch *ConsistentHashing) Ring() {
	for _, k := range ch.keys {
		fmt.Printf("%d:%s\n", k, string(ch.vals[k][0]))
	}
}
