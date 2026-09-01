package load_balancer

type LoadBalancer interface {
	Pick() string
	Rest()
}
