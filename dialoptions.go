package wecross

import (
	"time"

	"github.com/WeBankBlockchain/WeCross-Go-SDK/credentials"

	"github.com/WeBankBlockchain/WeCross-Go-SDK/internal/transport"

	internalbackoff "github.com/WeBankBlockchain/WeCross-Go-SDK/internal/backoff"
)

// dialOptions configure a Dial call. dialOptions are set by the DialOption
// values passed to Dial.
type dialOptions struct {
	bs              internalbackoff.Stratgy
	block           bool
	returnLastError bool
	timout          time.Duration
	copts           transport.ConnectOptions
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

// WithTransportCredentials returns a DialOption which configures a connection
// level security credentials (e.g., TLS/SSL). This should not be used together
// with WithCredentialsBundle.
func WithTransportCredentials(creds credentials.TransportCredentials) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.copts.TransportCredentials = creds
	})
}

// WithBlock returns a DialOption which makes callers of Dial block until the
// underlying connection is up. Without this, Dial returns immediately and
// connecting the server happens in background.
func WithBlock() DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.block = true
	})
}
