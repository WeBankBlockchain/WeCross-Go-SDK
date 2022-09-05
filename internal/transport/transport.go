package transport

import (
	"context"
	"fmt"
	"net"

	"github.com/WeBankBlockchain/WeCross-Go-SDK/credentials"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/keepalive"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/resolver"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/stats"
)

// state of transport
type transportState int

const (
	reachable transportState = iota
	closing
	draining
)

// Stream represents an RPC in the transport layer.
type Stream struct {
	id       uint32
	ct       *httpClient
	ctx      context.Context // the associated context of the stream
	done     chan struct{}   // closed at the end of stream to unblock writers. On the client side.
	doneFunc func()          // invoked at the end of stream on client side.
	method   string          // the associated RPC method of the stream
}

// ConnectOptions covers all relevant options for coummunicating with the server.
type ConnectOptions struct {
	// Dialer specifies how to dial a network address.
	Dialer func(context.Context, string) (net.Conn, error)
	// FailOnNonTempDialError specifies if WeCross fails on non-temporary dial errors.
	FailOnNonTempDialError bool
	// PerRPCCredentials stores the PerRPCCredentials required to issue RPCs.
	// Note that PerPRCCredentials has not been implemented.
	PerRPCCredentials []credentials.PerRPCCredentials
	// TransportCredentials stores the Authenticator required to setup a client
	// connection.
	TransportCredentials credentials.TransportCredentials
	CredsBundle          credentials.Bundle
	// KeepaliveParams stores the keepalive parameters.
	KeepaliveParams keepalive.ClientParameters
	// StatsHandlers stores the handler for stats.
	StatsHandlers []stats.Handler
	// WriteBufferSize sets the size of write buffer which in turn determines how much data can be batched before it's written on the wire.
	WriteBufferSize int
	// ReadBufferSize sets the size of read buffer, which in turn determines how much data can be read at most for one read syscall.
	ReadBufferSize int
}

func NewClientTransport(connectCtx, ctx context.Context, addr resolver.Address, opts ConnectOptions, onPrefaceReceipt func(), onClose func()) (ClientTransport, error) {
	return newHTTP2Client(connectCtx, ctx, addr, opts, onPrefaceReceipt, onClose)
}

// ClientTransport is the common interface for all WeCross SDK transport
// implementations.
type ClientTransport interface {
}

// connectionErrorf creates an ConnectionError with the specified error description.
func connectionErrorf(temp bool, e error, format string, a ...interface{}) ConnectionError {
	return ConnectionError{
		Desc: fmt.Sprintf(format, a...),
		temp: temp,
		err:  e,
	}
}

// ConnectionError is an error that results in the termination of the
// entire connection and the retry of all the active streams.
type ConnectionError struct {
	Desc string
	temp bool
	err  error
}

func (e ConnectionError) Error() string {
	return fmt.Sprintf("connection error: desc = %q", e.Desc)
}

// Temporary indicates if this connection error is temporary or fatal.
func (e ConnectionError) Temporary() bool {
	return e.temp
}

// Origin returns the original error of this connection error.
func (e ConnectionError) Origin() error {
	// Never return nil error here.
	// If the original error is nil, return itself.
	if e.err == nil {
		return e
	}
	return e.err
}

// Unwrap returns the original error of this connection error or nil when the
// origin is nil.
func (e ConnectionError) Unwrap() error {
	return e.err
}

var (
	// ErrConnClosing indicates that the transport is closing.
	ErrConnClosing = connectionErrorf(true, nil, "transport is closing")
)
