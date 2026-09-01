package main

type LoadBalancer interface {
	Pick() string
	Rest()
}
