package load_balancer

import "strings"

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

func (l7 *L7) AddRoute(host, prefix, upStream string) {
	route := Route{
		host:     host,
		prefix:   prefix,
		upStream: upStream,
	}

	if host == "*" {
		l7.wildcardRoutes = append(l7.wildcardRoutes, route)
	} else {
		l7.routesByHost[host] = append(l7.routesByHost[host], route)
	}
}

func (l7 *L7) Default(def string) {
	l7.def = def
}

func (l7 *L7) L7_(host, path string) string {
	// Get routes for this host
	routes := l7.routesByHost[host]

	// First try to find a match in host-specific routes
	bestUpstream := ""
	bestLen := 0

	if len(routes) > 0 {
		for _, route := range routes {
			prefix := route.prefix

			// Check if path starts with this prefix
			if strings.HasPrefix(path, prefix) || prefix == "" {
				length := len(prefix)

				// Exact match gets priority
				if prefix == path {
					length += 1000
				}

				if length > bestLen {
					bestLen = length
					bestUpstream = route.upStream
				}
			}
		}

		// If we found a match in host-specific routes, return it
		if bestUpstream != "" {
			return bestUpstream
		}
	}

	// If no host-specific match found, try wildcard routes
	for _, route := range l7.wildcardRoutes {
		prefix := route.prefix

		// Check if path starts with this prefix
		if strings.HasPrefix(path, prefix) || prefix == "" {
			length := len(prefix)

			// Exact match gets priority
			if prefix == path {
				length += 1000
			}

			if length > bestLen {
				bestLen = length
				bestUpstream = route.upStream
			}
		}
	}

	if bestUpstream == "" {
		return l7.def
	}

	return bestUpstream
}
