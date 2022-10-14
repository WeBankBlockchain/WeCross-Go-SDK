package wecrosstest

import (
	"reflect"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/WeBankBlockchain/WeCross-Go-SDK/internal/leakcheck"
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

// Setup updates the tlogger.
func (Tester) Setup(t *testing.T) {
	TLogger.Update(t)
}

// Teardown performs a leak check.
func (Tester) Teardown(t *testing.T) {
	if atomic.LoadUint32(&lcFailed) == 1 {
		return
	}
	leakcheck.Check(errorer{t: t})
	if atomic.LoadUint32(&lcFailed) == 1 {
		t.Log("Leak check disabled for future tests")
	}
	TLogger.EndTest(t)
}

func getTestFunc(t *testing.T, xv reflect.Value, name string) func(t2 *testing.T) {
	if m := xv.MethodByName(name); m.IsValid() {
		if f, ok := m.Interface().(func(*testing.T)); ok {
			return f
		}
		// Method exists but has the wrong type signature.
		t.Fatalf("wecrosstest: function %v has unexpected signature (%T)", name, m.Interface())
	}
	return func(*testing.T) {}
}

// RunSubTests runs all "Test___" functions that are methods of x as subtests
// of the current test. If x contains methods "Setup(*testing.T)" or
// "Teardown(*testing.T)", those are run before or after each of the test
// function, respectively.
//
// For example usage, see example_test.go. Run it using:
//  	$ go test -v -run TestExample
// To run a specific test/subtest:
// 		$ go test -v -run 'TestExample/^Something$'.
func RunSubTests(t *testing.T, x interface{}) {
	xt := reflect.TypeOf(x)
	xv := reflect.ValueOf(x)

	setup := getTestFunc(t, xv, "Setup")
	teardown := getTestFunc(t, xv, "Teardown")

	for i := 0; i < xv.NumMethod(); i++ {
		methodName := xt.Method(i).Name
		if !strings.HasPrefix(methodName, "Test") {
			continue
		}
		tFunc := getTestFunc(t, xv, methodName)
		t.Run(strings.TrimPrefix(methodName, "Test"), func(t *testing.T) {
			// Run leakcheck in t.Cleanup() to guarantee it is run even if tfunc
			// or setup uses t.Fatal().
			//
			// Note that a defer would run before t.Cleanup, so if a goroutine
			// is closed by a test's t.Cleanup, a deferred leakcheck would fail.
			t.Cleanup(func() { teardown(t) })
			setup(t)
			tFunc(t)
		})
	}
}
