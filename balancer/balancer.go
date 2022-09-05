package balancer

import (
	"context"
	"errors"
	"strings"

	"github.com/WeBankBlockchain/WeCross-Go-SDK/connectivity"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/internal"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/resolver"
)

// Builder creates a balancer.
type Builder interface {
	// Build creates a new balancer with the ClientConn.
	Build(cc ClientConn, opts BuildOptions) Balancer
	// Name returns the name of balancers built by this builder.
	// It will be used to pick balancers.
	Name() string
}

var (
	// m is a map from name to balancer builder.
	m = make(map[string]Builder)
)

// Register registers the balancer builder to the balancer map.b.Name
// (lowercased) will be used as the name registered with this builder. If the
// Builder implements ConfigParser, ParseConfig will be called when new service
// configs are received by the resolver, and the result will be provided to the
// Balancer in UpdateClientConnState.
//
// NOTE: this function must only be called during initialization time (i.e. in
// an init() function), and is not thread-safe. If multiple Balancers are
// registered with the same name, the one registered last will take effect.
func Register(b Builder) {
	m[strings.ToLower(b.Name())] = b
}

// unregisterForTesting deletes the balancer with the given name from the
// balancer map.
//
// This function is not thread-safe.
func unregisterForTesting(name string) {
	delete(m, name)
}

func init() {
	internal.BalancerUnregister = unregisterForTesting
}

// Get returns the resolver builder registered with the given name.
// Note that the compare is done in a case-insensitive fashion.
// If no builder is register with the name, nil will be returned.
func Get(name string) Builder {
	if b, ok := m[strings.ToLower(name)]; ok {
		return b
	}
	return nil
}

// Balancer takes input from WeCross, manages SubConns, and collects and aggregates
// the connectivity states.
//
// It also generates and updates the Picker used by WeCross to pick SubConns for RPCs.
//
// UpdateClientConnState, ResolverError, UpdateSubConnState, and Close are
// guaranteed to be called synchronously from the same goroutine. There's no
// guarantee on picker.Pick, it may be called anytime.
type Balancer interface {
	// UpdateClientConnState is called by WeCross when the state of the ClientConn
	// changes. If the error returned is ErrBadResolverState, the ClientConn
	// will begin calling ResolveNow on the active name resolver with
	// exponential backoff until a subsequent call to UpdateClientConnState
	// returns a nil error. Any other errors are currently ignored.
	UpdateClientConnState(ClientConnState) error
	// ResolverError is called by WeCross when the name resolver reports an error.
	ResolverError(error)
	// UpdateSubConnState is called by WeCross when the state of a SubConn
	// changes.
	UpdateSubConnState(SubConn, SubConnState)
	// Close closes the balancer. The balancer is not required to call
	// ClientConn.RemoveSubConn for its existing SubConns.
	Close()
}

// SubConnState describes the state of a SubConn.
type SubConnState struct {
	// ConnectivityState is the connectivity state of the SubConn.
	ConnectivityState connectivity.State
	// ConnectionError is set if the ConnectivityState is TransientFailure,
	// describing the reason the SubConn failed. Otherwise, it is nil.
	ConnectionError error
}

// A SubConn represents a single connection to a backend service.
//
// All SubConns start in IDLE, and will not try to connect. To trigger the
// connecting, Balancers must call Connect. If a connection re-enters IDLE,
// Balancers must call Connect again to trigger a new connection attempt.
//
// WeCross will try to connect to the addresses in sequence, and stop trying the
// remainder once the first connection is successful. If an attempt to connect
// to all addresses encounters an error, the SubConn will enter
// TRANSIENT_FAILURE for a backoff period, and then transition to IDLE.
//
// Once established, if a connection is lost, the SubConn will transition
// directly to IDLE.
//
// This interface is to be implemented by WeCross. Users should not need their own
// implementation of this interface. For situations like testing, any
// implementations should embed this interface. This allows WeCross to add new
// methods to this interface.
type SubConn interface {
	Connect()
}

// State contains the balancer's state relevant to the WeCross ClientConn.
type State struct {
	// State contains the connectivity state of the balancer, which is used to
	// determine the state of the ClientConn.
	ConnectivityState connectivity.State
	// Picker is used to choose connections (SubConns) for RPCs.
	Picker Picker
}

// ClientConn represents a WeCross ClientConn.
//
// This interface is to be implemented by WeCross. Users should not need a
// brand new implementation of this interface. For the situations like
// testing, the new implementation should embed this interface. This allows
// WeCross to add new methods to this interface.
type ClientConn interface {
	// NewSubConn is called by balancer to create a new SubConn.
	// It doesn't block and wait for the connections to be established.
	// Behaviors of the SubConn can be controlled by options.
	NewSubConn([]resolver.Address) (SubConn, error)
	// RemoveSubConn removes the SubConn from ClientConn.
	// The SubConn will be shutdown.
	RemoveSubConn(SubConn)
	// UpdateAddresses updates the addresses used in the passed in SubConn.
	// WeCross checks if the currently connected address is still in the new list.
	// If so, the connection will be kept. Else, the connection will be
	// gracefully closed, and a new connection will be created.
	//
	// This will trigger a state transition for the SubConn.
	UpdateAddresses(SubConn, []resolver.Address)
	// UpdateState notifies WeCross that the balancer's internal state has
	// changed.
	//
	// WeCross will update the connectivity state of the ClientConn, and will call
	// Pick on the new Picker to pick new SubConns.
	UpdateState(State)
}

// BuildOptions contains additional information for Build.
type BuildOptions struct{}

type PickInfo struct {
	// FullMethodName is the method name that NewClientStream() is called
	// with. The canonical format is /service/Method.
	FullMethodName string
	// Ctx is the RPC's context, and may contain relevant RPC-level information
	// like the outgoing header metadata.
	Ctx context.Context
}

// DoneInfo contains additional information for done.
type DoneInfo struct {
	// Err is the rpc error the RPC finished with. It could be nil.
	Err error
	// BytesSent indicates if any bytes have been sent to the server.
	BytesSent bool
	// BytesReceived indicates if any byte has been received from the server.
	BytesReceived bool
}

var (
	// ErrNoSubConnAvailable indicates no SubConn is available for pick().
	// WeCross will block the WeCross until a new picker is available via UpdateState().
	ErrNoSubConnAvailable = errors.New("no SubConn is available")
)

// PickResult contains information related to a connection chosen for an RPC.
type PickResult struct {
	// SubConn is the connection to use for this pick, if its state is Ready.
	// If the state is not Ready, WeCross will block the RPC until a new Picker is
	// provided by the balancer (using ClientConn.UpdateState). The SubConn
	// must be one returned by ClientConn.NewSubConn.
	SubConn SubConn

	// Done is called when the RPC is completed.  If the SubConn is not ready,
	// this will be called with a nil parameter.  If the SubConn is not a valid
	// type, Done may not be called.  May be nil if the balancer does not wish
	// to be notified when the RPC completes.
	Done func(DoneInfo)
}

// Picker is used by WeCross to pick a SubConn to send an RPC.
// Balancer is expected to generate a new picker from its snapshot every time its
// internal state has changed.
//
// The pickers used by WeCross can be updated by ClientConn.UpdateState().
type Picker interface {
	// Pick returns the connection to use for this RPC and related information.
	//
	// Pick should not block. If the balancer needs to do I/O or any blocking
	// or time-consuming work to service this call, it should return
	// ErrNoSubConnAvailable, and Pick call will  be repeated by WeCross when
	// the Picker is updated (using ClientConn.UpdateState.)
	//
	// If an error is returned:
	//
	// - If the error is a ErrNoSubConnAvailable, WeCross will block until a new
	//   Picker is provided by the balancer (using ClientConn.UpdateState).
	//
	// - If the error is a status error (implemented by the WeCross/status
	//   package), WeCross will terminate the RPC with the code and message
	//   provided.
	//
	// - For all other errors, wait for ready RPCs will wait, but non-wait for
	//   ready RPCs will be terminated with this error's Error() string and
	//   status code Unavailable.
	Pick(info PickInfo) (PickResult, error)
}

// ExitIdler is an optional interface for balancers to implement. If
// implemented, ExitIdle will be called when ClientConn.Connect is called, if
// the ClientConn is idle. If unimplemented, ClientConn.Connect will cause
// all SubConns to connect.
type ExitIdler interface {
	// ExitIdle instructs the LB policy to reconnect to backends / exit the
	// IDLE state, if appropriate and possible. Note that SubConns that enter
	// the IDLE state will not reconnect until SubConn.Connect is called.
	ExitIdle()
}

// ClientConnState describes the state of a ClientConn relevant to the
// balancer.
type ClientConnState struct {
	ResolverState resolver.State
}

// ErrBadResolverState may be returned by UpdateClientConnState to indicate a
// problem with the provided name resolver data.
var ErrBadResolverState = errors.New("bad resolver state")
