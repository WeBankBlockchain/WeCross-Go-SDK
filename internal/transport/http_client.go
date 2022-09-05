package transport

import (
	"context"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/WeBankBlockchain/WeCross-Go-SDK/credentials"
	icredentials "github.com/WeBankBlockchain/WeCross-Go-SDK/internal/credentials"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/internal/syscall"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/internal/transport/networktype"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/keepalive"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/resolver"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/stats"
)

// httpClient implements the ClientTransport interface with HTTP.
type httpClient struct {
	lastRead   int64 // Keep this field 64-bit aligned. Accessed atomically.
	ctx        context.Context
	cancel     context.CancelFunc
	ctxDone    <-chan struct{} // Cache the ctx.Done() chan.
	conn       net.Conn
	remoteAddr net.Addr
	localAddr  net.Addr
	authInfo   credentials.AuthInfo // auth info about the connection

	readerDone chan struct{} // sync point to enable testing.
	writerDone chan struct{} // sync point to enable testing.

	// controlBuf delivers all the control related tasks (e.g., window
	// updates, reset streams, and various settings) to the controller.
	controlBuf *controlBuffer
	// The scheme used: https if TLS is on, http otherwise.
	scheme string

	isSecure bool

	perRPCCreds []credentials.PerRPCCredentials

	kp               keepalive.ClientParameters
	keepaliveEnabled bool

	statsHandlers []stats.Handler

	mu            sync.Mutex // guard the following variables
	activeStreams map[uint32]*Stream
	state         transportState
	// A condition variable used to signal when the keepalive goroutine should
	// go dormant. The condition for dormancy is based on the number of active
	// streams and the `PermitWithoutStream` keepalive client parameter. And
	// since the number of active streams is guarded by the above mutex, we use
	// the same for this condition variable as well.
	kpDormancyCond *sync.Cond
	// A boolean to track whether the keepalive goroutine is dormant or not.
	// This is checked before attempting to signal the above condition
	// variable.
	kpDormant bool

	onClose func()
}

func dial(ctx context.Context, fn func(context.Context, string) (net.Conn, error), addr resolver.Address) (net.Conn, error) {
	address := addr.Addr
	networkType, _ := networktype.Get(addr)
	if fn != nil {
		return fn(ctx, address)
	}
	return (&net.Dialer{}).DialContext(ctx, networkType, address)
}

func isTemporary(err error) bool {
	switch err := err.(type) {
	case interface {
		Temporary() bool
	}:
		return err.Temporary()
	case interface {
		Timeout() bool
	}:
		// Timeouts may be resolved upon retry, and are thus treated as
		// temporary.
		return err.Timeout()
	}
	return true
}

// newHTTPClient constructs a connected ClientTransport to addr based on HTTP
// and starts to receive messages on it. Non-nil error returns if construction
// fails.
func newHTTPClient(connectCtx, ctx context.Context, addr resolver.Address, opts ConnectOptions, onPrefaceReceipt func(), onClose func()) (_ *httpClient, err error) {
	scheme := "http"
	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		if err != nil {
			cancel()
		}
	}()

	// resolver, balancer etc. can specify arbitrary data in the
	// Attributes field of resolver. Address, which is shoved into connectCtx
	// and passed to the dialer and credential handshaker. This makes it possible for
	// address specific arbitrary data to reach custom dialers and credential hanshakers.
	connectCtx = icredentials.NewClientHandshakeInfoContext(connectCtx, credentials.ClientHandshakeInfo{Attributes: addr.Attributes})

	conn, err := dial(connectCtx, opts.Dialer, addr)
	if err != nil {
		if opts.FailOnNonTempDialError {
			return nil, connectionErrorf(isTemporary(err), err, "transport: error while dialing: %v", err)
		}
		return nil, connectionErrorf(true, err, "transport: error while dialing: %v", err)
	}
	// Any further errors will close the underlying connection
	defer func(conn net.Conn) {
		if err != nil {
			conn.Close()
		}
	}(conn)
	kp := opts.KeepaliveParams
	// Validate keepalive parameters.
	if kp.Time == 0 {
		kp.Time = defaultClientKeepaliveTime
	}
	if kp.Timeout == 0 {
		kp.Timeout = defaultClientKeppaliveTimeout
	}
	keepaliveEnabled := false
	if kp.Time != infinity {
		if err = syscall.SetTCPUserTimeout(conn, kp.Timeout); err != nil {
			return nil, connectionErrorf(false, err, "transport: failed to set TCP_USER_TIMEOUT: %v", err)
		}
		keepaliveEnabled = true
	}
	var (
		isSecure bool
		authInfo credentials.AuthInfo
	)
	transportCreds := opts.TransportCredentials
	perRPCCreds := opts.PerRPCCredentials

	if b := opts.CredsBundle; b != nil {
		if t := b.TransportCredentials(); t != nil {
			transportCreds = t
		}
		if t := b.PerRPCCredentials(); t != nil {
			perRPCCreds = append(perRPCCreds, t)
		}
	}
	if transportCreds != nil {
		rawConn := conn
		// Pull the deadline from the connectCtx, which will be used for
		// timeouts in the authentication protocol handshake. Can ignore the
		// boolean as the deadline will return the zero value, which will make
		// the conn not timeout on I/O operations.
		deadline, _ := connectCtx.Deadline()
		rawConn.SetDeadline(deadline)
		conn, authInfo, err = transportCreds.ClientHandshake(connectCtx, addr.Addr, rawConn)
		rawConn.SetDeadline(time.Time{})
		if err != nil {
			return nil, connectionErrorf(isTemporary(err), err, "transport: authentication handshake failed: %v", err)
		}
		for _, cd := range perRPCCreds {
			if cd.RequireTransportSecurity() {
				if ci, ok := authInfo.(interface {
					GetCommonAuthInfo() credentials.CommonAuthInfo
				}); ok {
					secLevel := ci.GetCommonAuthInfo().SecurityLevel
					if secLevel < credentials.SslOff {
						return nil, connectionErrorf(true, nil, "transport: cannot send secure credentials on an insecure connection")
					}
				}
			}
		}
		isSecure = true
		if transportCreds.Info().SecurityProtocol == "tls" {
			scheme = "https"
		}
	}
	t := &httpClient{
		ctx:              ctx,
		ctxDone:          ctx.Done(),
		cancel:           cancel,
		conn:             conn,
		remoteAddr:       conn.RemoteAddr(),
		localAddr:        conn.LocalAddr(),
		authInfo:         authInfo,
		readerDone:       make(chan struct{}),
		writerDone:       make(chan struct{}),
		scheme:           scheme,
		isSecure:         isSecure,
		perRPCCreds:      perRPCCreds,
		kp:               kp,
		statsHandlers:    opts.StatsHandlers,
		keepaliveEnabled: keepaliveEnabled,
	}
	for _, sh := range t.statsHandlers {
		t.ctx = sh.TagConn(t.ctx, &stats.ConnTagInfo{
			RemoteAddr: t.remoteAddr,
			LocalAddr:  t.localAddr,
		})
		connBegin := &stats.ConnBegin{
			Client: true,
		}
		sh.HandleConn(t.ctx, connBegin)
	}
	if t.keepaliveEnabled {
		t.kpDormancyCond = sync.NewCond(&t.mu)
		go t.keepalive()
	}

}

// Close kicks off the shutdown process of the transport. This should be called
// only once on a transport. Once it is called, the transport should not be
// accessed any more.
//
// This method blocks until the addrConn that initiated this transport is
// re-connected. This happens because t.onClose() begins reconnect logic at the
// addrConn level and blocks until the addrConn is successfully connected.
func (t *httpClient) Close(err error) {
	t.mu.Lock()
	// Make sure we only Close once.
	if t.state == closing {
		t.mu.Unlock()
		return
	}
	// Call t.onClose before setting the state to closing to prevent the client
	// from attempting to create new streams ASAP.
	t.onClose()
	t.state = closing
	streams := t.activeStreams
	t.activeStreams = nil
	if t.kpDormant {
		// If the keepalive goroutine is blocked on this condition variable, we
		// should unblock it so that the goroutine eventually exits.
		t.kpDormancyCond.Signal()
	}
	t.mu.Unlock()
	t.controlBuf.finish()
	t.cancel()
	t.conn.Close()

	// Notify all active streams.
	for _, s := range streams {
		t.closeStream()
	}
	for _, sh := range t.statsHandlers {
		connEnd := &stats.ConnEnd{
			Client: true,
		}
		sh.HandleConn(t.ctx, connEnd)
	}
}

func minTime(a, b time.Duration) time.Duration {
	if a < b {
		return a
	}
	return b
}

// keepalive running in a separate goroutine makes sure the connection is alive by sending pings.
func (t *httpClient) keepalive() {
	p := &ping{data: [8]byte{}}
	// Ture if a ping has been sent, and no data has been received since then.
	outstandingPing := false
	// Amount of time remaining before which we should receive an ACK for the
	// last sent ping.
	timeoutLeft := time.Duration(0)
	// Records the last value of t.lastRead before we go block on the timer.
	// This is required to check for read activity since then.
	prevNano := time.Now().UnixNano()
	timer := time.NewTimer(t.kp.Time)
	for {
		select {
		case <-timer.C:
			lastRead := atomic.LoadInt64(&t.lastRead)
			if lastRead > prevNano {
				// There has been read activity since the last time we were here.
				outstandingPing = false
				// Next timer should fire at kp.Time seconds from lastRead time.
				timer.Reset(time.Duration(lastRead) + t.kp.Time - time.Duration(time.Now().UnixNano()))
				prevNano = lastRead
				continue
			}
			if outstandingPing && timeoutLeft <= 0 {
				t.Close(connectionErrorf(true, nil, "keepalive ping failed to receive ACK within timeout"))
				return
			}
			t.mu.Lock()
			if t.state == closing {
				// If the transport is closing, we should exit from the
				// keepalive goroutine here. If not, we could have a race
				// between the call to Signal() from Close() and the call to
				// Wait() here, whereby the keepalive goroutine ends up
				// blocking on the condition variable which will never be
				// signalled again.
				t.mu.Unlock()
				return
			}
			t.kpDormant = false
			t.mu.Unlock()

			// We get here because the keepalive timer expired.
			// We need to send a ping.
			if !outstandingPing {
				t.controlBuf.put(p)
				timeoutLeft = t.kp.Timeout
				outstandingPing = true
			}
			// The amount of time to sleep here is the minimum of kp.Time and
			// timeoutLeft. This will ensure that we wait only for kp.Time
			// before sending out the next ping (for cases where the ping is
			// acked).
			sleepDuration := minTime(t.kp.Time, timeoutLeft)
			timeoutLeft -= sleepDuration
			timer.Reset(sleepDuration)
		case <-t.ctx.Done():
			if !timer.Stop() {
				<-timer.C
			}
			return
		}
	}
}
