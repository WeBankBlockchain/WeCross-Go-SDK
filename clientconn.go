package wecross

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"

	"github.com/WeBankBlockchain/WeCross-Go-SDK/codes"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/connectivity"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/internal/backoff"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/internal/status"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/wecrosslog"
)

// minium time to give a connection to complete
const minConnectTimeout = 20 * time.Second

var (
	ErrClientConnClosing = status.Err(codes.Canceled, "grpc: the client connection is closing")
	// errConnDrain indicates that the connection starts to be drained and does not accept any new RPCs.
	errConnDrain = errors.New("wecross: the connection is drained")
	// errConnClosing indicates that the connection is closing.
	errConnClosing = errors.New("wecross: the connection is closing")
	// invalidDefaultServiceConfigErrPrefix is used to prefix the json parsing error for the default
	// service config.
	invalidDefaultServiceConfigErrPrefix = "wecross: the provided default service config is invalid"
)

// The following errors are returned from Dial and DialContext
var (
	// errNoTransportSecurity indicates that there is no transport security
	// being set for ClientConn. Users should either set one or explicitly
	// call WithInsecure DialOption to disable security.
	errNoTransportSecurity = errors.New("wecross: no transport security set (use wecross.WithTransportCredentials(insecure.NewCredentials()) explicitly or set credentials)")
	// errTransportCredsAndBundle indicates that creds bundle is used together
	// with other individual Transport Credentials.
	errTransportCredsAndBundle = errors.New("wecorss: credentials.Bundle may not be used with individual TransportCredentials")
	// errNoTransportCredsInBundle indicates that the configured creds bundle
	// returned a transport credentials which was nil.
	errNoTransportCredsInBundle = errors.New("wecross: credentials.Bundle must return non-nil transport credentials")
	// errTransportCredentialsMissing indicates that users want to transmit
	// security information (e.g., OAuth2 token) which requires secure
	// connection on an insecure connection.
	errTransportCredentialsMissing = errors.New("wecorss: the credentials require transport level security(use wecross.WithTransportCredentials() to set)")
)

const (
	defaultClientMaxReceiveMessageSize = 1024 * 1024 * 4
	defaultClientMaxSendMessageSize    = math.MaxInt32
	// http2IOBufSize specifies the buffer size for sending frames.
	defaultWriteBufSize = 32 * 1024
	defaultReadBufSize  = 32 * 1024
)

var logger = wecrosslog.Component("core")

// Dial creates a clinet connection to the given target.
func Dial(target string, opts ...DialOption) (*ClientConn, error) {
	return DialContext(context.Background(), target, opts...)
}

// DialContext creates a client connection to the given target. By default, it's
// a non-blocking dial (the function won't wait for connecting to be
// established, and connecting happens in the background). To make it a blocking
// dial, use WithBlock() dial option.
//
// In the non-blocking case, the ctx does not act against the connection. It
// only controls the setup steps.
//
// In the blocking case, ctx can be used to cancel or expire the pending
// connection. Once this function returns, the cancellation and expiration of
// ctx will be noop. Users should call ClientConn.Close to terminate all the
// pending operations after this function returns.
func DialContext(ctx context.Context, target string, opts ...DialOption) (conn *ClientConn, err error) {
	cc := &ClientConn{
		target: target,
		csMgr:  &connectivityStateManager{},
		conns:  make(map[*addrConn]struct{}),
		dopts:  defaultDialOptions(),
	}
	cc.retryThrottler.Store((*retryThrottler)(nil))
	cc.ctx, cc.cancel = context.WithCancel(context.Background())

	for _, opt := range opts {
		opt.apply(&cc.dopts)
	}

	defer func() {
		if err != nil {
			cc.Close()
		}
	}()

	if cc.dopts.copts.TransportCredentials == nil && cc.dopts.copts.CredsBundle == nil {
		return nil, errNoTransportSecurity
	}
	if cc.dopts.copts.TransportCredentials != nil && cc.dopts.copts.CredsBundle != nil {
		return nil, errTransportCredsAndBundle
	}
	if cc.dopts.copts.CredsBundle != nil && cc.dopts.copts.CredsBundle.TransportCredentials() == nil {
		return nil, errNoTransportCredsInBundle
	}
	transportCreds := cc.dopts.copts.TransportCredentials
	if transportCreds == nil {
		transportCreds = cc.dopts.copts.CredsBundle.TransportCredentials()
	}
	if transportCreds.Info().SecurityProtocol == "insecure" {
		for _, cd := range cc.dopts.copts.PerRPCCredentials {
			if cd.RequireTransportSecurity() {
				return nil, errTransportCredentialsMissing
			}
		}
	}

	if cc.dopts.timout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, cc.dopts.timout)
		defer cancel()
	}
	defer func() {
		select {
		case <-ctx.Done():
			switch {
			case ctx.Err() == err:
				conn = nil
			case err == nil || !cc.dopts.returnLastError:
				conn, err = nil, ctx.Err()
			default:
				conn, err = nil, fmt.Errorf("%v: %v", ctx.Err(), err)
			}
		default:
		}
	}()

	if cc.dopts.bs == nil {
		cc.dopts.bs = backoff.DefaultExponential
	}

	// A blocking dial blocks until the clientConn is ready.
	if cc.dopts.block {
		for {
			s := cc.GetState()
			if s == connectivity.Ready {
				break
			} else if cc.dopts.copts.FailOnNonTempDialError && s == connectivity.TransientFailure {
				if err = cc.connectionError(); err != nil {
					terr, ok := err.(interface {
						Temporary() bool
					})
					if ok && !terr.Temporary() {
						return nil, err
					}
					return nil, ctx.Err()
				}
			}
		}
	}

	return cc, nil
}

// ClientConn represents a virtual connection to a conceptual endpoint, to
// perform RPCs.
//
// A ClientConn encapsulates a range of functionality including name
// resolution, TCP connection establishment (with retries and backoff) and TLS
// handshakes. It also handles errors on established connections by
// re-resolving the name and reconnecting.
type ClientConn struct {
	ctx    context.Context    // Initialized using the background context at dial time.
	cancel context.CancelFunc // Cancelled on close.

	// THe following are initialized at dial time, and are read-only after that.
	target string      // User's dial target.
	dopts  dialOptions // Default and user specified dial options.

	// The following provide their own synchronization, and therefore don't
	// require cc.mu to be held to access them.
	csMgr          *connectivityStateManager
	retryThrottler atomic.Value // Update from config.

	// mu protects the following fields.
	mu    sync.RWMutex
	conns map[*addrConn]struct{} // Set to nil on close.

	lceMu               sync.Mutex // protects lastConnectionError
	lastConnectionError error
}

// GetState returns the connectivity.State of ClientConn.
func (cc *ClientConn) GetState() connectivity.State {
	return cc.csMgr.getState()
}

func (cc *ClientConn) connectionError() error {
	cc.lceMu.Lock()
	defer cc.lceMu.Unlock()
	return cc.lastConnectionError
}

// connectivityStateManager keeps the connectivity.State of ClientConn.
type connectivityStateManager struct {
	mu         sync.Mutex
	state      connectivity.State
	notifyChan chan struct{}
}

// updateState updates the connectivity. State of ClientConn.
// If there's a change it notifies goroutines waiting on state change to
// happen.
func (csm *connectivityStateManager) updateState(state connectivity.State) {
	csm.mu.Lock()
	defer csm.mu.Unlock()
	if csm.state == connectivity.Shutdown {
		return
	}
	if csm.state == state {
		return
	}
	csm.state = state
	logger.Infof("Channel Connectivity change to %v", state)
	if csm.notifyChan != nil {
		// There are other goroutines waiting on this channel.
		close(csm.notifyChan)
		csm.notifyChan = nil
	}
}

func (csm *connectivityStateManager) getState() connectivity.State {
	csm.mu.Lock()
	defer csm.mu.Unlock()
	return csm.state
}

// addrConn is a network connection to a given address.
type addrConn struct {
	ctx    context.Context
	cancel context.CancelFunc

	cc    *ClientConn
	dopts dialOptions

	mu sync.Mutex

	// Use updateConnectivityState for updating addrConn's connectivity state.
	state connectivity.State

	backoffIdx   int // Needs to be stateful for resetConnectBackoff.
	resetBackoff chan struct{}
}

// Note: this requires a lock on ac.mu.
func (ac *addrConn) updateConnectivityState(s connectivity.State) {
	if ac.state == s {
		return
	}
	ac.state = s
	logger.Infof("Subchannel Connectivity change to %v", s)
}

// tearDown starts to tear down the addrConn.
//
// Note that tearDown doesn't remove ac from ac.cc.conns, so the addrConn struct
// will leak. In most cases, call cc.removeAddrConn() instead.
func (ac *addrConn) tearDown(err error) {
	ac.mu.Lock()
	if ac.state == connectivity.Shutdown {
		ac.mu.Unlock()
		return
	}
	// We have to set the state to Shutdown before anything else to prevent races
	// between setting the state and logic that waits on context cancellatio / etc.
	ac.updateConnectivityState(connectivity.Shutdown)
	ac.cancel()
	ac.mu.Unlock()
}

// removeAddrConn removes the addrConn in the subConn from clientConn.
// It also tears down the ac with the given error.
func (cc *ClientConn) removeAddrConn(ac *addrConn, err error) {
	cc.mu.Lock()
	if cc.conns == nil {
		cc.mu.Unlock()
		return
	}
	delete(cc.conns, ac)
	cc.mu.Unlock()
	ac.tearDown(err)
}

func (cc *ClientConn) Close() error {
	defer cc.cancel()

	cc.mu.Lock()
	if cc.conns == nil {
		cc.mu.Unlock()
		return ErrClientConnClosing
	}
	conns := cc.conns
	cc.conns = nil
	cc.csMgr.updateState(connectivity.Shutdown)

	cc.mu.Unlock()

	for ac := range conns {
		ac.tearDown(ErrClientConnClosing)
	}

	return nil
}

type retryThrottler struct {
	max    float64
	thresh float64
	ratio  float64

	mu     sync.Mutex
	tokens float64
}
