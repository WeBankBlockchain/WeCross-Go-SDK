package wecross

import (
	"context"
	"net"
	"time"

	"github.com/WeBankBlockchain/WeCross-Go-SDK/credentials"

	"github.com/WeBankBlockchain/WeCross-Go-SDK/internal/transport"

	internalbackoff "github.com/WeBankBlockchain/WeCross-Go-SDK/internal/backoff"
)

// dialOptions configure a DialContext call. dialOptions are set by the DialOption
// values passed to DialContext.
type dialOptions struct {
	bs                internalbackoff.Stratgy
	block             bool
	returnLastError   bool
	timeout           time.Duration
	copts             transport.ConnectOptions
	minConnectTimeout func() time.Duration
}

// DialOption configures how we set up the connection.
type DialOption interface {
	apply(*dialOptions)
}

func defaultDialOptions() dialOptions {
	return dialOptions{
		copts: transport.ConnectOptions{
			WriteBufferSize: defaultWriteBufSize,
			ReadBufferSize:  defaultReadBufSize,
		},
	}
}

// funcDialOption wraps a function that modifies dialOptions into an
// implementation of the DialOption interface.
type funcDialOption struct {
	f func(*dialOptions)
}

func (fdo *funcDialOption) apply(do *dialOptions) {
	fdo.f(do)
}

func newFuncDialOption(f func(*dialOptions)) *funcDialOption {
	return &funcDialOption{
		f: f,
	}
}

// WithContextDialer returns a DialOption that sets a dialer to create
// connections. If FailOnNonTempDialError() is set to true, and an error is
// returned by f, WeCross checks the error's Temporary() method to decide if it
// should try to reconnect to the network address.
func WithContextDialer(f func(context.Context, string) (net.Conn, error)) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.copts.Dialer = f
	})
}

// FailOnNonTempDialError returns a DialOption that specifies if WeCross fails on
// non-temporary dial errors. If f is true, and dialer returns a non-temporary
// error, WeCross will fail the connection to the network address and won't try to
// reconnect. The default value of FailOnNonTempDialError is false.
//
// FailOnNonTempDialError only affects the initial dial, and does not do
// anything useful unless you are also using WithBlock().
func FailOnNonTempDialError(f bool) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.copts.FailOnNonTempDialError = f
	})
}

// WithTransportCredentials returns a DialOption which configures a connection
// level security credentials (e.g., TLS/SSL). This should not be used together
// with WithCredentialsBundle.
func WithTransportCredentials(creds credentials.TransportCredentials) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.copts.TransportCredentials = creds
	})
}

// WithBlock returns a DialOption which makes callers of DialContext block until the
// underlying connection is up. Without this, DialContext returns immediately and
// connecting the server happens in background.
func WithBlock() DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.block = true
	})
}
