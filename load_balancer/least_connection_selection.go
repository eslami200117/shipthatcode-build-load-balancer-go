package load_balancer

import "fmt"

type LeastConnSel struct {
	mh *MinHeap
}

func NewLeastConnSel(bacnends []string) *LeastConnSel {
	minHeap := NewMinHeap()
	for _, b := range bacnends {
		minHeap.InsertOrUpdate(b, 0)
	}
	return &LeastConnSel{
		mh: minHeap,
	}
}

func (l *LeastConnSel) Pick() string {
	key, _, ok := l.mh.GetMin()
	if !ok {
		panic("pool is empty")
	}
	l.mh.Add(key, 1)

	return key
}

func (l *LeastConnSel) Done(backend string) {
	l.mh.Add(backend, -1)
}


func (l *LeastConnSel) Status() {
	items := l.mh.GetAllOrdered()
	for _, item := range items {
		fmt.Printf("%s:%d\n", item.Key, item.Value)
	}
}
