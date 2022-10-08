package types

import (
	"WeCross-Go-SDK/common"
	"sync"
	"time"
)

var CALLBACK_TIMEOUT = 30000 * time.Millisecond

type CallBack struct {
	finishMux sync.Mutex
	isFinish  bool // should be changed in atomic manner

	onSuccess func(response *Response)
	onFailed  func(err common.WeCrossSDKError)

	quit chan struct{}
}

// atomically set isFinish and back the old value
func (cb *CallBack) GetAndSetIsFinish(newValue bool) bool {
	cb.finishMux.Lock()
	defer cb.finishMux.Unlock()
	old := cb.isFinish
	cb.isFinish = newValue
	return old
}

func (cb *CallBack) CallOnSuccess(response *Response) {
	if !cb.GetAndSetIsFinish(true) {
		close(cb.quit)
		cb.onSuccess(response)
	}
}

func (cb *CallBack) CallOnFailed(err common.WeCrossSDKError) {
	if !cb.GetAndSetIsFinish(true) {
		close(cb.quit)
		cb.onFailed(err)
	}
}

func NewCallBack(onSuccess func(response *Response), onFailed func(err common.WeCrossSDKError)) *CallBack {
	newCB := &CallBack{
		isFinish:  false,
		onSuccess: onSuccess,
		onFailed:  onFailed,
		quit:      make(chan struct{}),
	}
	go newCB.timeOut()
	return newCB
}

func (cb *CallBack) timeOut() {
	timer := time.NewTimer(CALLBACK_TIMEOUT)
	select {
	case <-timer.C:
		if !cb.GetAndSetIsFinish(true) {
			err := common.NewWeCrossSDKFromString(common.REMOTECALL_ERROR, "Timeout")
			cb.onFailed(err)
			return
		}
	case <-cb.quit:
		return
	}
}
