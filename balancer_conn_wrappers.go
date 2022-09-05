package wecross

import (
	"fmt"
	"strings"
	"sync"

	"github.com/WeBankBlockchain/WeCross-Go-SDK/balancer"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/connectivity"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/internal/balancer/gracefulswitch"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/internal/buffer"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/internal/wecrosssync"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/resolver"
)

type ccBalancerWrapper struct {
	cc *ClientConn

	// Since these fields are accessed only from handleXxx() methods which are
	// synchronized by the watcher goroutine, we do not need a mutex to protect
	// these fields.
	balancer        *gracefulswitch.Balancer
	curBalancerName string

	updateCh *buffer.Unbounded  // Updates written on this channel are processed by watcher().
	resultCh *buffer.Unbounded  // Results of calls to UpdateClientConnState() are pushed here.
	closed   *wecrosssync.Event // Indicates if close has been called.
	done     *wecrosssync.Event // Indicates if close has completed its work.
}

func newCCBalancerWrapper(cc *ClientConn, bopts balancer.BuildOptions) *ccBalancerWrapper {
	ccb := &ccBalancerWrapper{
		cc:       cc,
		updateCh: buffer.NewUnbounded(),
		resultCh: buffer.NewUnbounded(),
		closed:   wecrosssync.NewEvent(),
		done:     wecrosssync.NewEvent(),
	}
	go ccb.watcher()
	ccb.balancer = gracefulswitch.NewBalancer(ccb, bopts)
	return ccb
}

// The following xxxUpdate structs wrap the arguments received as part of the
// corresponding update. The watcher goroutine uses the 'type' of the update to
// invoke the appropriate handler routine to handle the update.

type ccStateUpdate struct {
	ccs *balancer.ClientConnState
}

type scStateUpdate struct {
	sc    balancer.SubConn
	state connectivity.State
	err   error
}

type exitIdleUpdate struct{}

type resolverErrorUpdate struct {
	err error
}

type switchToUpdate struct {
	name string
}

type subConnUpdate struct {
	acbw *acBalancerWrapper
}

func (ccb *ccBalancerWrapper) watcher() {
	for {
		select {
		case u := <-ccb.updateCh.Get():
			ccb.updateCh.Load()
			if ccb.closed.HasFired() {
				break
			}
			switch update := u.(type) {
			case *ccStateUpdate:
				ccb.handleClientConnStateChange(update.ccs)
			case *scStateUpdate:
				ccb.handlerSubConnStateChange(update)
			case *exitIdleUpdate:
				ccb.handleExitIdle()
			case *resolverErrorUpdate:
				ccb.handleResolverError(update.err)
			case *switchToUpdate:
				ccb.handleSwitchTo(update.name)
			case *subConnUpdate:
				ccb.handleRemoveSubConn(update.acbw)
			default:
				logger.Errorf("ccBalancerWrapper.watcher: unknown update %+v, type %T", update, update)
			}
		case <-ccb.closed.Done():
		}

		if ccb.closed.HasFired() {
			ccb.handleClose()
			return
		}
	}
}

// handleClientConnStateChange handles a ClientConnState update from the update
// channel and invokes the appropriate method on the underlying balancer.
//
// If the addresses specified in the update contain addresses of type "wecrosslb"
// and the selected LB policy is not "wecrosslb", these addresses will be filtered
// out and ccs will be modified with the updated address list.
func (ccb *ccBalancerWrapper) handleClientConnStateChange(ccs *balancer.ClientConnState) {
	if ccb.curBalancerName != wecrosslbName {
		// Filter any wecrosslb addresses since we don't have the wecrosslb balancer.
		var addrs []resolver.Address
		for _, addr := range ccs.ResolverState.Addresses {
			if addr.BalancerAttributes == nil {
				continue
			}
			addrs = append(addrs, addr)
		}
		ccs.ResolverState.Addresses = addrs
	}
	ccb.resultCh.Put(ccb.balancer.UpdateClientConnState(*ccs))
}

// updateSubConnState is invoked by WeCross to push a subConn state update to the
// underlying balancer.
func (ccb *ccBalancerWrapper) updateSubConnState(sc balancer.SubConn, s connectivity.State, err error) {
	// When updating addresses for a SubConn, if the address in use is not in
	// the new addresses, the old ac will be tearDown() and a new ac will be
	// created. tearDown() generates a state change with Shutdown state, we
	// don't want the balancer to receive this state change. So before
	// tearDown() on the old ac, ac.acbw (acWrapper) will be set to nil, and
	// this function will be called with (nil, Shutdown, err). We don't need to call
	// balancer method in this case.
	if sc == nil {
		return
	}
	ccb.updateCh.Put(&scStateUpdate{
		sc:    sc,
		state: s,
		err:   err,
	})
}

// handleSubConnStateChange handles a SubConnState update from the update
// channel and invokes the appropriate method on the underlying balancer.
func (ccb *ccBalancerWrapper) handlerSubConnStateChange(update *scStateUpdate) {
	ccb.balancer.UpdateSubConnState(update.sc, balancer.SubConnState{ConnectivityState: update.state, ConnectionError: update.err})
}

func (ccb *ccBalancerWrapper) handleExitIdle() {
	if ccb.cc.GetState() != connectivity.Idle {
		return
	}
	ccb.balancer.ExitIdle()
}

func (ccb *ccBalancerWrapper) handleResolverError(err error) {
	ccb.balancer.ResolverError(err)
}

// switchTo is invoked by WeCross to instruct the balancer wrapper to switch to the
// LB policy identified by name.
//
// ClientConn calls newCCBalancerWrapper() at creation time. Upon receipt of the
// first good update from the name resolver, it determines the LB policy to use
// and invokes the switchTo() method. Upon receipt of every subsequent update
// from the name resolver, it invokes this method.
//
// the ccBalancerWrapper keeps track of the current LB policy name, and skips
// the graceful balancer switching process if the name does not change.
func (ccb *ccBalancerWrapper) switchTo(name string) {
	ccb.updateCh.Put(&switchToUpdate{name: name})
}

// handleSwitchTo handles a balancer switch update from the update channel. It
// calls the SwitchTo() method on the gracefulswitch.Balancer with a
// balancer.Builder corresponding to name. If no balancer.Builder is registered
// for the given name, it uses the default LB policy which is "pick_first".
func (ccb *ccBalancerWrapper) handleSwitchTo(name string) {
	if strings.EqualFold(ccb.curBalancerName, name) {
		return
	}

	builder := balancer.Get(name)
	if builder == nil {
		logger.Warningf("Channel switches to new LB policy %q, since the specified LB policy %q was not registered", PickFirstBalancerName, name)
		builder = newPickfirstBuilder()
	} else {
		logger.Infof("Channel switches to new LB policy %q", name)
	}

	if err := ccb.balancer.SwitchTo(builder); err != nil {
		logger.Errorf("Channel failed to build new LB policy %q: %v", name, err)
		return
	}
	ccb.curBalancerName = builder.Name()
}

// handleRemoveSucConn handles a request from the underlying balancer to remove
// a subConn.
//
// See comments in RemoveSubConn() for more details.
func (ccb *ccBalancerWrapper) handleRemoveSubConn(acbw *acBalancerWrapper) {
	ccb.cc.removeAddrConn(acbw.getAddrConn(), errConnDrain)
}

func (ccb *ccBalancerWrapper) handleClose() {
	ccb.balancer.Close()
	ccb.done.Done()
}

func (ccb *ccBalancerWrapper) NewSubConn(addrs []resolver.Address) (balancer.SubConn, error) {
	if len(addrs) <= 0 {
		return nil, fmt.Errorf("grpc: cannot create SubConn with empty address list")
	}
	ac, err := ccb.cc.newAddrConn(addrs)
	if err != nil {
		logger.Warningf("acBalancerWrapper: NewSubConn: failed to newAddrConn: %v", err)
		return nil, err
	}
	acbw := &acBalancerWrapper{ac: ac}
	acbw.ac.mu.Lock()
	ac.acbw = acbw
	acbw.ac.mu.Unlock()
	return acbw, nil
}

func (ccb *ccBalancerWrapper) RemoveSubConn(sc balancer.SubConn) {
	// Before we switched the ccBalancerWrapper to use gracefulswitch.Balancer, it
	// was required to handle the RemoveSubConn() method asynchronously by pushing
	// the update onto the update channel. This was done to avoid a deadlock as
	// switchBalancer() was holding cc.mu when calling Close() on the old
	// balancer, which would in turn call RemoveSubConn().
	//
	// With the use of gracefulswitch.Balancer in ccBalancerWrapper, handling this
	// asynchronously is probably not required anymore since the switchTo() method
	// handles the balancer switch by pushing the update onto the channel.
	// TODO(easwars): Handle this inline.
	acbw, ok := sc.(*acBalancerWrapper)
	if !ok {
		return
	}
	ccb.updateCh.Put(&subConnUpdate{acbw: acbw})
}

func (ccb *ccBalancerWrapper) UpdateAddresses(sc balancer.SubConn, addrs []resolver.Address) {
	acbw, ok := sc.(*acBalancerWrapper)
	if !ok {
		return
	}
	acbw.UpdateAddresses(addrs)
}

func (ccb *ccBalancerWrapper) UpdateState(s balancer.State) {
	// Update picker before updating state.  Even though the ordering here does
	// not matter, it can lead to multiple calls of Pick in the common start-up
	// case where we wait for ready and then perform an RPC.  If the picker is
	// updated later, we could call the "connecting" picker when the state is
	// updated, and then call the "ready" picker after the picker gets updated.
	ccb.cc.blockingpicker.updatePicker(s.Picker)
	ccb.cc.csMgr.updateState(s.ConnectivityState)
}

// acBalancerWrapper is a wrapper on top of ac for balancers.
// It implements balancer.SubConn interface.
type acBalancerWrapper struct {
	mu sync.Mutex
	ac *addrConn
}

func (acbw *acBalancerWrapper) UpdateAddresses(addrs []resolver.Address) {
	acbw.mu.Lock()
	defer acbw.mu.Unlock()
	if len(addrs) <= 0 {
		acbw.ac.cc.removeAddrConn(acbw.ac, errConnDrain)
		return
	}
	if !acbw.ac.tryUpdateAddrs(addrs) {

	}
}

func (acbw *acBalancerWrapper) Connect() {
	acbw.mu.Lock()
	defer acbw.mu.Unlock()
	go acbw.ac.connect()
}

func (acbw *acBalancerWrapper) getAddrConn() *addrConn {
	acbw.mu.Lock()
	defer acbw.mu.Unlock()
	return acbw.ac
}
