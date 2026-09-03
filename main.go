package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type L4 struct {
	backends []string
}

func NewL4(backends []string) *L4 {
	return &L4{
		backends: backends,
	}
}

func (l *L4) L4(ip, port string) string {
	index := hashString(ip+":"+port) % len(l.backends)
	return l.backends[index]
}

func hashString(s string) int {
	h := 0
	for _, c := range s {
		h = (h*31 + int(c)) % 1000
	}
	return h
}

type Route struct {
	host     string
	prefix   string
	upStream string
}

type L7 struct {
	routesByHost   map[string][]Route
	wildcardRoutes []Route
	def            string
}

func NewL7() *L7 {
	return &L7{
		routesByHost:   make(map[string][]Route),
		wildcardRoutes: []Route{},
	}
}


func (l7 *L7) Default(def string) {
	l7.def = def
}

func (l7 *L7) AddRoute(host, prefix, upStream string) {
	route := Route{
		host:     host,
		prefix:   prefix,
		upStream: upStream,
	}

	if host == "*" {
		// Check if wildcard route with same prefix exists, replace it
		for i, r := range l7.wildcardRoutes {
			if r.prefix == prefix {
				l7.wildcardRoutes[i] = route
				return
			}
		}
		l7.wildcardRoutes = append(l7.wildcardRoutes, route)
	} else {
		// Check if route with same host and prefix exists, replace it
		for i, r := range l7.routesByHost[host] {
			if r.prefix == prefix {
				l7.routesByHost[host][i] = route
				return
			}
		}
		l7.routesByHost[host] = append(l7.routesByHost[host], route)
	}
}

func (l7 *L7) L7_(host, path string) string {
	// Get routes for this host
	hostRoutes := l7.routesByHost[host]
	
	// First, find the best match among host-specific routes
	bestHostUpstream := ""
	bestHostLen := 0
	
	for i := len(hostRoutes) - 1; i >= 0; i-- {
		route := hostRoutes[i]
		prefix := route.prefix
		
		if strings.HasPrefix(path, prefix) || prefix == "" {
			length := len(prefix)
			if prefix == path {
				length += 1000
			}
			if length >= bestHostLen {
				bestHostLen = length
				bestHostUpstream = route.upStream
			}
		}
	}
	
	// Find the best match among wildcard routes
	bestWildcardUpstream := ""
	bestWildcardLen := 0
	
	for i := len(l7.wildcardRoutes) - 1; i >= 0; i-- {
		route := l7.wildcardRoutes[i]
		prefix := route.prefix
		
		if strings.HasPrefix(path, prefix) || prefix == "" {
			length := len(prefix)
			if prefix == path {
				length += 1000
			}
			if length >= bestWildcardLen {
				bestWildcardLen = length
				bestWildcardUpstream = route.upStream
			}
		}
	}
	
	// Choose the best match: host-specific wins only if it's more specific
	if bestHostUpstream != "" && bestHostLen >= bestWildcardLen {
		return bestHostUpstream
	}
	
	if bestWildcardUpstream != "" {
		return bestWildcardUpstream
	}
	
	return l7.def
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

	// file, err := os.Open(fmt.Sprintf("tests/12-l4-vs-l7/%d.in", test))
	// if err != nil {
	// 	panic(err)
	// }
	// defer file.Close()
	// sc := bufio.NewScanner(file)
	sc := bufio.NewScanner(os.Stdin)
	sc.Buffer(make([]byte, 1024*1024), 1024*1024)
	var l4 *L4
	l7 := NewL7()
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
		case "POOL4":
			l4 = NewL4(args[1:])
			fmt.Println("OK")
		case "L4":
			ans := l4.L4(args[1], args[2])
			fmt.Println(ans)
		case "ROUTE":
			l7.AddRoute(args[1], args[2], args[3])
			fmt.Println("OK")
		case "DEFAULT":
			l7.Default(args[1])
			fmt.Println("OK")
		case "L7":
			ans := l7.L7_(args[1], args[2])
			fmt.Println(ans)
		default:
			fmt.Println("Wrong argument:", args[0])
		}
	}
}
