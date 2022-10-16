package transactionContext

import (
	"fmt"
	"testing"
	"time"
)

func TestTransactionContex(t *testing.T) {
	testTxCtx := NewTransactionContex()
	tempCtx := &TxCtx{
		TxID:              "201313",
		Seq:               1,
		PathInTransaction: []string{"20", "13", "13"},
	}
	testTxCtx.SetContex(tempCtx)

	errChan := make(chan error)
	go func() {
		tempCtx := testTxCtx.GetContex()
		t.Log(tempCtx)
		time.Sleep(1 * time.Second)
		tempCtx.TxID = "201314"
		tempCtx.Seq = 2
		tempCtx.PathInTransaction = []string{"20", "13", "14"}
		testTxCtx.SetContex(tempCtx)
	}()

	go func() {
		tempCtx1 := testTxCtx.GetContex()
		if tempCtx1.TxID != "201313" {
			errChan <- fmt.Errorf("txID is changed")
			return
		}
		time.Sleep(2 * time.Second)
		tempCtx2 := testTxCtx.GetContex()

		if tempCtx1.TxID != "201313" {
			errChan <- fmt.Errorf("txID is changed but should not")
			return
		}

		if tempCtx2.TxID != "201314" {
			errChan <- fmt.Errorf("txID is not right")
			return
		}
		t.Log(tempCtx2)
		errChan <- nil
	}()
	select {
	case err := <-errChan:
		if err != nil {
			t.Fatalf(err.Error())
		}
	}

}
