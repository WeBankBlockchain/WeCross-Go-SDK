package transaction

import (
	"fmt"

	internalwecrosslog "github.com/WeBankBlockchain/WeCross-Go-SDK/internal/wecrosslog"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/wecrosslog"
)

var logger = wecrosslog.Component("transaction")

type transactionB struct{}

func (tb *transactionB) Build() {
	tB := &transaction{}
	tB.logger = internalwecrosslog.NewPrefixLogger(logger, fmt.Sprintf("[rpc-transaction %p]", tB))
}

type transaction struct {
	logger *internalwecrosslog.PrefixLogger
}
