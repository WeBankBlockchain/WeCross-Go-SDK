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
	txID              TxCtxKey = "currentXATransactionID"
	seq               TxCtxKey = "currentXATransactionSeq"
	pathInTransaction TxCtxKey = "pathInTransaction"
)

type TransactionContext struct {
	logger               *logger.Logger
	localCtx             context.Context
	cancel               context.CancelFunc
	txIDCtx              context.Context
	seqCtx               context.Context
	pathInTransactionCtx context.Context
}

func NewTransactionContex() *TransactionContext {
	logger := logger.NewLogger(TransactionContextTag)
	localCtx, cancel := context.WithCancel(context.Background())

	return &TransactionContext{
		logger:               logger,
		localCtx:             localCtx,
		cancel:               cancel,
		txIDCtx:              context.WithValue(localCtx, txID, ""),
		seqCtx:               context.WithValue(localCtx, seq, int(0)),
		pathInTransactionCtx: context.WithValue(localCtx, pathInTransaction, []string{}),
	}
}

func (tc *TransactionContext) ClearAll() {
	tc.txIDCtx = context.WithValue(tc.localCtx, txID, "")
	tc.seqCtx = context.WithValue(tc.localCtx, seq, int(0))
	tc.pathInTransactionCtx = context.WithValue(tc.localCtx, pathInTransaction, []string{})
}

func (tc *TransactionContext) Close() {
	tc.cancel()
}

func (tc *TransactionContext) SetTransactionID(inputTxID string) {
	tc.txIDCtx = context.WithValue(tc.localCtx, txID, inputTxID)
}

func (tc *TransactionContext) GetTransactionID() string {
	return tc.txIDCtx.Value(txID).(string)
}

func (tc *TransactionContext) SetSeqNumber(inputNumber int) {
	tc.seqCtx = context.WithValue(tc.localCtx, seq, inputNumber)
}

func (tc *TransactionContext) GetSeqNumber() int {
	return tc.seqCtx.Value(seq).(int)
}

func (tc *TransactionContext) SetPathInTransaction(paths []string) {
	tc.pathInTransactionCtx = context.WithValue(tc.localCtx, pathInTransaction, paths)
}

func (tc *TransactionContext) GetPathInTransaction() []string {
	return tc.pathInTransactionCtx.Value(pathInTransaction).([]string)
}

func (tc *TransactionContext) IsPathInTransaction(path string) bool {
	if len(tc.txIDCtx.Value(txID).(string)) == 0 {
		return false
	}
	paths := tc.pathInTransactionCtx.Value(pathInTransaction).([]string)
	if paths == nil { // this should not happen as the pathInTransaction could not be nil
		tc.logger.Errorf("IsPathInTransaction: pathList is nil.")
		return false
	}
	if len(paths) == 0 {
		tc.logger.Errorf("IsPathInTransaction: TransactionID exist, but pathList doesn't.")
	}
	for i := 0; i < len(paths); i++ {
		if paths[i] == path {
			return true
		}
	}
	return false
}
