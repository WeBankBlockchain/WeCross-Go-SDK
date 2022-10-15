package transactionContext

import (
	"reflect"
	"testing"
	"time"
)

func TestTransactionContex(t *testing.T) {
	testTxCtx := NewTransactionContex()
	testPaths := []string{"20", "13", "13"}

	go func() {
		time.Sleep(1 * time.Second)
		testTxCtx.SetTransactionID("201313")
		testTxCtx.SetSeqNumber(201313)
		testTxCtx.SetPathInTransaction(testPaths)
	}()

	go func() {
		t.Log(testTxCtx.GetTransactionID())
		t.Log(testTxCtx.GetSeqNumber())
		t.Log(testTxCtx.GetPathInTransaction())
	}()

	time.Sleep(2 * time.Second)
	if testTxCtx.GetTransactionID() != "201313" {
		t.Fatalf("fain in set txID")
	}
	if testTxCtx.GetSeqNumber() != 201313 {
		t.Fatalf("fail in set seq")
	}
	if !reflect.DeepEqual(testTxCtx.GetPathInTransaction(), testPaths) {
		t.Fatalf("fail in set pathInTransaction")
	}
	t.Log(testTxCtx.GetTransactionID())
	t.Log(testTxCtx.GetSeqNumber())
	t.Log(testTxCtx.GetPathInTransaction())

}
