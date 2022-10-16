package transactionContext

import (
	"WeCross-Go-SDK/logger"
	"context"
)

const (
	TransactionContextTag = "TxContex"
)

type TxCtxKey string

const (
	txCtxKey TxCtxKey = "TxCtx"
)

var TransactionContextLogger = logger.NewLogger(TransactionContextTag)

type TransactionContext struct {
	baseCtx    context.Context
	currentCtx context.Context
	cancel     context.CancelFunc
}

type TxCtx struct {
	TxID              string
	Seq               int
	PathInTransaction []string
}

func NewTransactionContex() *TransactionContext {

	baseCtx, cancel := context.WithCancel(context.Background())
	newTxCtx := &TransactionContext{
		baseCtx: baseCtx,
		cancel:  cancel,
	}

	newTxCtx.ClearAll()

	return newTxCtx
}

func (tc *TransactionContext) SetContex(txCtx *TxCtx) {
	tc.currentCtx = context.WithValue(tc.baseCtx, txCtxKey, txCtx)
}

func (tc *TransactionContext) GetContex() *TxCtx {
	txCtx := tc.currentCtx.Value(txCtxKey).(*TxCtx)
	return txCtx.Copy()
}

func (tc *TransactionContext) ClearAll() {
	emptyTxCtx := &TxCtx{
		TxID:              "",
		Seq:               -1,
		PathInTransaction: make([]string, 0),
	}
	tc.SetContex(emptyTxCtx)
}

func (tc *TransactionContext) Close() {
	tc.cancel()
}

func (tctx *TxCtx) IsPathInTransaction(path string) bool {
	if len(tctx.PathInTransaction) == 0 {
		return false
	}
	paths := tctx.PathInTransaction
	if paths == nil { // this should not happen as the pathInTransaction could not be nil
		TransactionContextLogger.Errorf("IsPathInTransaction: pathList is nil.")
		return false
	}
	if len(paths) == 0 {
		TransactionContextLogger.Errorf("IsPathInTransaction: TransactionID exist, but pathList doesn't.")
	}
	for i := 0; i < len(paths); i++ {
		if paths[i] == path {
			return true
		}
	}
	return false
}

func (tctx *TxCtx) Copy() *TxCtx {
	pathInTx := make([]string, 0)
	pathInTx = append(pathInTx, tctx.PathInTransaction...)
	return &TxCtx{
		TxID:              tctx.TxID,
		Seq:               tctx.Seq,
		PathInTransaction: pathInTx,
	}
}
