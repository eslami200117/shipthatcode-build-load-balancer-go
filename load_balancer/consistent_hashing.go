package load_balancer

import (
	"fmt"
	"sort"
	"strconv"
)

type ConsistentHashing struct {
	ring  map[int]string
	keys  []int
	dirty bool
}

func NewConsistentHashing(backends []string) *ConsistentHashing {
	ch := &ConsistentHashing{
		ring: make(map[int]string),
		keys: make([]int, 0),
	}
	for _, b := range backends {
		ch.Add(b)
	}
	return ch
}

func (ch *ConsistentHashing) Add(b string) {
	for j := 0; j < 5; j++ {
		p := hashString(b + "#" + strconv.Itoa(j))
		if _, exists := ch.ring[p]; !exists {
			ch.ring[p] = b
			ch.keys = append(ch.keys, p)
			ch.dirty = true
		}
	}
}

func (ch *ConsistentHashing) LookUp(p string) string {
	if len(ch.keys) == 0 {
		return ""
	}

	if ch.dirty {
		sort.Ints(ch.keys)
		ch.dirty = false
	}

	key := hashString(p)
	idx := sort.Search(len(ch.keys), func(i int) bool {
		return ch.keys[i] >= key
	})

	if idx == len(ch.keys) {
		idx = 0
	}

	return ch.ring[ch.keys[idx]]
}

func (ch *ConsistentHashing) Ring() {
	if ch.dirty {
		sort.Ints(ch.keys)
		ch.dirty = false
	}

	for _, k := range ch.keys {
		fmt.Printf("%d:%s\n", k, string(ch.ring[k][0]))
	}
}

func hashString(s string) int {
	h := 0
	for _, c := range s {
		h = (h*31 + int(c)) % 1000
	}
	return h
}
