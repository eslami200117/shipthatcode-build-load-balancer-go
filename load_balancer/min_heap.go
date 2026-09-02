package load_balancer

import (
	"container/heap"
)

// Item represents an element in the heap
type Item struct {
	key   string
	value int
	index int // index in the heap (for updates)
	order int // insertion order number for stability
}

// PriorityQueue implements heap with stability support
type PriorityQueue []*Item

// Len returns the number of elements
func (pq PriorityQueue) Len() int { return len(pq) }

// Less compares two items to determine priority
func (pq PriorityQueue) Less(i, j int) bool {
	// First compare by value
	if pq[i].value != pq[j].value {
		return pq[i].value < pq[j].value
	}
	// If values are equal, compare by insertion order (stability)
	return pq[i].order < pq[j].order
}

// Swap swaps two elements
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

// Push adds an element to the heap
func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

// Pop removes and returns the highest priority element
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update updates an item's value
func (pq *PriorityQueue) update(item *Item, value int) {
	if value < 0 {
		value = 0
	}
	item.value = value
	heap.Fix(pq, item.index)
}

// MinHeap combines dictionary and heap functionality
type MinHeap struct {
	pq           PriorityQueue
	items        map[string]*Item
	order        int      // insertion order counter
	insertionSeq []string // keeps track of insertion order
}

// NewMinHeap creates a new MinHeap
func NewMinHeap() *MinHeap {
	mh := &MinHeap{
		pq:    PriorityQueue{},
		items: make(map[string]*Item),
		order: 0,
	}
	heap.Init(&mh.pq)
	return mh
}

// InsertOrUpdate adds a new key or updates an existing key's value
func (mh *MinHeap) InsertOrUpdate(key string, value int) {
	if item, exists := mh.items[key]; exists {
		// Update existing item
		mh.pq.update(item, value)
	} else {
		// Add new item
		mh.order++
		item := &Item{
			key:   key,
			value: value,
			order: mh.order,
		}
		mh.items[key] = item
		mh.insertionSeq = append(mh.insertionSeq, key) // track insertion order
		heap.Push(&mh.pq, item)
	}
}

// Add increments a key's value by delta
func (mh *MinHeap) Add(key string, delta int) bool {
	if item, exists := mh.items[key]; exists {
		mh.pq.update(item, item.value+delta)
		return true
	}
	return false
}

// GetMin returns the key with the minimum value (without removing)
func (mh *MinHeap) GetMin() (string, int, bool) {
	if len(mh.pq) == 0 {
		return "", 0, false
	}
	item := mh.pq[0]
	return item.key, item.value, true
}

// PopMin removes and returns the key with the minimum value
func (mh *MinHeap) PopMin() (string, int, bool) {
	if len(mh.pq) == 0 {
		return "", 0, false
	}
	item := heap.Pop(&mh.pq).(*Item)
	delete(mh.items, item.key)
	return item.key, item.value, true
}

// GetValue returns the value of a specific key
func (mh *MinHeap) GetValue(key string) (int, bool) {
	if item, exists := mh.items[key]; exists {
		return item.value, true
	}
	return 0, false
}

// Size returns the number of elements
func (mh *MinHeap) Size() int {
	return len(mh.pq)
}

// GetAllOrdered returns all items in original insertion order (more efficient)
func (mh *MinHeap) GetAllOrdered() []struct {
	Key   string
	Value int
} {
	result := []struct {
		Key   string
		Value int
	}{}

	// Iterate over insertion sequence
	for _, key := range mh.insertionSeq {
		if item, exists := mh.items[key]; exists {
			result = append(result, struct {
				Key   string
				Value int
			}{
				Key:   item.key,
				Value: item.value,
			})
		}
	}
	return result
}
