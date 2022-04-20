package resource

import (
	"fmt"
	internalwecrosslog "github.com/WeBankBlockchain/WeCross-Go-SDK/internal/wecrosslog"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/wecrosslog"
)

var logger = wecrosslog.Component("resource")

type resourceB struct{}

func (resourceB) Build() {
	rB := &resource{}
	rB.logger = internalwecrosslog.NewPrefixLogger(logger, fmt.Sprintf("[resource-resource %p]", rB))
}

type resource struct {
	logger *internalwecrosslog.PrefixLogger
}
