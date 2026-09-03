package load_balancer

import (
	"fmt"
	"math"
)

type RateLimit struct {
	capacity float64
	status   map[string]*RLInfo
}

type RLInfo struct {
	lastUpdate float64
	token      float64
}

var nowRL float64

func NewRateLimit() *RateLimit {
	return &RateLimit{
		capacity: 10,
		status:   make(map[string]*RLInfo),
	}
}

func (r *RateLimit) Now(str string) {
	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		panic("invalid input")
	}
	nowRL = f
	for _, info := range r.status {
		elapse := nowRL - info.lastUpdate
		info.token = math.Min(r.capacity, info.token+elapse)
	}

}

func (r *RateLimit) Request(b string) string {
	info, ok := r.status[b]
	if !ok {
		r.status[b] = &RLInfo{
			lastUpdate: nowRL,
			token:      10,
		}
		info = r.status[b]
	}
	elapse := nowRL - info.lastUpdate
	info.token = math.Min(r.capacity, info.token+elapse)
	info.lastUpdate = nowRL
	if info.token >= 1 {
		info.token -= 1
		return "OK"
	}
	return "LIMITED"
}

func (r *RateLimit) Status(b string) {
	fmt.Printf("%.2f\n", r.status[b].token)
}
