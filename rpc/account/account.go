package account

import (
	"fmt"

	internalwecrosslog "github.com/WeBankBlockchain/WeCross-Go-SDK/internal/wecrosslog"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/wecrosslog"
)

var logger = wecrosslog.Component("account")

type accountB struct{}

func (ab *accountB) Build() {
	aB := &account{}
	aB.logger = internalwecrosslog.NewPrefixLogger(logger, fmt.Sprintf("[rpc-account %p]", aB))
}

type account struct {
	logger *internalwecrosslog.PrefixLogger
}
