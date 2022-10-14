package wecross

import (
	"errors"
	"fmt"

	"github.com/WeBankBlockchain/WeCross-Go-SDK/balancer"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/connectivity"
)

// PickFirstBalancerName is the name of the pick_first balancer.
const PickFirstBalancerName = "pick_first"

func newPickfirstBuilder() balancer.Builder {
	return &pickfirstBuilder{}
}

type pickfirstBuilder struct{}

func (*pickfirstBuilder) Build(cc balancer.ClientConn, opt balancer.BuildOptions) balancer.Balancer {
	return &pickfirstBalancer{cc: cc}
}

func (*pickfirstBuilder) Name() string {
	return PickFirstBalancerName
}

type pickfirstBalancer struct {
	state   connectivity.State
	cc      balancer.ClientConn
	subConn balancer.SubConn
}

func (b *pickfirstBalancer) ResolverError(err error) {
	if logger.V(2) {
		logger.Infof("pickfirstBalancer: ResolverError called with error %v", err)
	}
	if b.subConn == nil {
		b.state = connectivity.TransientFailure
	}

	if b.state != connectivity.TransientFailure {
		// The picker will not change since the balancer does not currently
		// report an error.
		return
	}
	b.cc.UpdateState(balancer.State{
		ConnectivityState: connectivity.TransientFailure,
		Picker:            &picker{err: fmt.Errorf("name resolver error: %v", err)},
	})
}

func (b *pickfirstBalancer) UpdateClientConnState(state balancer.ClientConnState) error {
	if len(state.ResolverState.Addresses) == 0 {
		// The resolver reported an empty address list. Treat it like an error by
		// calling b.ResolverError.
		if b.subConn != nil {
			// Remove the old subConn. All addresses were removed, so it is no longer
			// valid.
			b.cc.RemoveSubConn(b.subConn)
			b.subConn = nil
		}
		b.ResolverError(errors.New("produced zero addresses"))
		return balancer.ErrBadResolverState
	}

	if b.subConn != nil {
		b.cc.UpdateAddresses(b.subConn, state.ResolverState.Addresses)
		return nil
	}

	subConn, err := b.cc.NewSubConn(state.ResolverState.Addresses)
	if err != nil {
		if logger.V(2) {
			logger.Errorf("pickfirstBalancer: failed to NewSubConn: %v", err)
		}
		b.state = connectivity.TransientFailure
		b.cc.UpdateState(balancer.State{
			ConnectivityState: connectivity.TransientFailure,
			Picker:            &picker{err: fmt.Errorf("error creating connection: %v", err)},
		})
		return balancer.ErrBadResolverState
	}
	b.subConn = subConn
	b.state = connectivity.Idle
	b.cc.UpdateState(balancer.State{
		ConnectivityState: connectivity.Idle,
		Picker:            &picker{result: balancer.PickResult{SubConn: b.subConn}},
	})
	b.subConn.Connect()
	return nil
}

func (b *pickfirstBalancer) UpdateSubConnState(subConn balancer.SubConn, state balancer.SubConnState) {
	if logger.V(2) {
		logger.Infof("pickfirstBalancer: UpdateSubConnState: %p, %v", subConn, state)
	}
	if b.subConn != subConn {
		if logger.V(2) {
			logger.Infof("pickfirstBalancer: ignored state change because subConn is not recognized")
		}
		return
	}
	b.state = state.ConnectivityState
	if state.ConnectivityState == connectivity.Shutdown {
		b.subConn = nil
		return
	}

	switch state.ConnectivityState {
	case connectivity.Ready:
		b.cc.UpdateState(balancer.State{
			ConnectivityState: state.ConnectivityState,
			Picker:            &picker{result: balancer.PickResult{SubConn: subConn}},
		})
	case connectivity.Connecting:
		b.cc.UpdateState(balancer.State{
			ConnectivityState: state.ConnectivityState,
			Picker:            &picker{err: balancer.ErrNoSubConnAvailable},
		})
	case connectivity.Idle:
		b.cc.UpdateState(balancer.State{
			ConnectivityState: state.ConnectivityState,
			Picker:            &idlePicker{subConn},
		})
	case connectivity.TransientFailure:
		b.cc.UpdateState(balancer.State{
			ConnectivityState: state.ConnectivityState,
			Picker:            &picker{err: state.ConnectionError},
		})
	}
}

func (b *pickfirstBalancer) Close() {
}

func (b *pickfirstBalancer) ExitIdle() {
	if b.subConn != nil && b.state == connectivity.Idle {
		b.subConn.Connect()
	}
}

type picker struct {
	result balancer.PickResult
	err    error
}

func (p *picker) Pick(balancer.PickInfo) (balancer.PickResult, error) {
	return p.result, p.err
}

// idlePicker is used when the SubConn is IDLE and kicks the SubConn into
// CONNECTING when Pick is called.
type idlePicker struct {
	subConn balancer.SubConn
}

func (i *idlePicker) Pick(balancer.PickInfo) (balancer.PickResult, error) {
	i.subConn.Connect()
	return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
}

func init() {
	balancer.Register(newPickfirstBuilder())
}
