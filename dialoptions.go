package wecross

import (
	"time"

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
