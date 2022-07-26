package wecross

import (
	"testing"

	"github.com/WeBankBlockchain/WeCross-Go-SDK/internal/wecrosstest"
)

type s struct {
	wecrosstest.Tester
}

func Test(t *testing.T) {
	wecrosstest.RunSubTests(t, s{})
}
