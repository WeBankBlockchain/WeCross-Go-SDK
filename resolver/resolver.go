package resolver

import "github.com/WeBankBlockchain/WeCross-Go-SDK/attributes"

// Address represents a server the client connects to.
type Address struct {
	// Addr is the server address on which a connection will be established.
	Addr string

	// Attributes contains arbitrary data about this address intended for
	// consumption by the SubConn.
	Attributes *attributes.Attributes

	// BalancerAttributes contains arbitrary data about this address intended
	// for consumption by the LB policy. These attribes do not affect SubConn
	// creation, connection establishment, handshaking, etc.
	BalancerAttributes *attributes.Attributes
}

// Equal returns whether a and o are identical.
func (a Address) Equal(o Address) bool {
	return a.Addr == o.Addr &&
		a.Attributes.Equal(o.Attributes) &&
		a.BalancerAttributes.Equal(o.BalancerAttributes)
}

// State contains the current Resolver state relevant to the ClientConn.
type State struct {
	// Addresses is the latest set of resolved addresses for the target.
	Addresses []Address

	// Attributes contains arbitrary data about the resolver intended for
	// consumption by the load balancing policy.
	Attributes *attributes.Attributes
}
