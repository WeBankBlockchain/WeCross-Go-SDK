package stats

// RPCStats contains stats information about RPCs.
type RPCStats interface {
	isRPCStats()
	// IsClient returns true if this RPCStats is from client side.
	IsClient() bool
}

// ConnStats contains stats information about connections.
type ConnStats interface {
	isConnStats()
	// IsClient returns true if this ConnStats is from client side.
	IsClient() bool
}

// ConnBegin contains the stats of a connection when it is established.
type ConnBegin struct {
	// Client is true if this ConnBegin is from client side.
	Client bool
}

// IsClient indicates if this is from client side.
func (s *ConnBegin) IsClient() bool { return s.Client }

func (s *ConnBegin) isConnStats() {}

// ConnEnd contains the stats of a connection when it ends.
type ConnEnd struct {
	// Client is true if this ConnEnd is from client side.
	Client bool
}

// IsClient indicates if this is from client side.
func (s *ConnEnd) IsClient() bool { return s.Client }

func (s *ConnEnd) isConnStats() {}
