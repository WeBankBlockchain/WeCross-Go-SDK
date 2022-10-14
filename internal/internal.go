package internal

var (
	// BalancerUnregister is exported by package balancer to unregister a balancer.
	BalancerUnregister func(name string)
)
