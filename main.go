package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
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
		info.lastUpdate = nowRL
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

func main() {
	// var test int
	// var err error
	// if len(os.Args) < 2 {
	// 	test = 2
	// } else {
	// 	test, err = strconv.Atoi(os.Args[1])
	// 	if err != nil {
	// 		panic("wrong arguman")
	// 	}
	// }

	// file, err := os.Open(fmt.Sprintf("tests/10-rate-limiting/%d.in", test))
	// if err != nil {
	// 	panic(err)
	// }
	// defer file.Close()
	// sc := bufio.NewScanner(file)
	sc := bufio.NewScanner(os.Stdin)
	sc.Buffer(make([]byte, 1024*1024), 1024*1024)
	lb := NewRateLimit()
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
		case "NOW":
			lb.Now(args[1])
			fmt.Println("OK")
		case "REQUEST":
			ans := lb.Request(args[1])
			fmt.Println(ans)
		case "STATUS":
			lb.Status(args[1])
		}
	}
}
