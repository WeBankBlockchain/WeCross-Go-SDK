package wecrosstest

import (
	"sync/atomic"
	"testing"
)

// leak check failed
var lcFailed uint32

type errorer struct {
	t *testing.T
}

func (e errorer) Errorf(format string, args ...interface{}) {
	atomic.StoreUint32(&lcFailed, 1)
	e.t.Errorf(format, args...)
}

// Tester is an implementation of x interface paramter to
// wecrosstest. RunSubTests with default Setup and Teardown behavior. Setup updates
// the tlogger and Teardown performs a leak check. Embed in a struct with tests
// defined to use.
type Tester struct{}

func (Tester) Setup(t *testing.T) {

}
