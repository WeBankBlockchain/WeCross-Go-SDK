package wecrosstest

import (
	"testing"

	wecrosslogi "github.com/WeBankBlockchain/WeCross-Go-SDK/internal/wecrosslog"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/wecrosslog"
)

type s struct {
	Tester
}

func Test(t *testing.T) {
	RunSubTests(t, s{})
}

func (s) TestInfo(t *testing.T) {
	wecrosslog.Info("Info", "message.")
}

func (s) TestInfoln(t *testing.T) {
	wecrosslog.Infoln("Info", "message.")
}

func (s) TestInfof(t *testing.T) {
	wecrosslog.Infof("%v %v.", "Info", "message")
}

func (s) TestInfoDepth(t *testing.T) {
	wecrosslogi.InfoDepth(0, "Info", "depth", "message.")
}

func (s) TestWarning(t *testing.T) {
	wecrosslog.Warning("Warning", "message.")
}

func (s) TestWarningln(t *testing.T) {
	wecrosslog.Warningln("Warning", "message.")
}

func (s) TestWarningf(t *testing.T) {
	wecrosslog.Warningf("%v %v.", "Warning", "message")
}

func (s) TestWarningDepth(t *testing.T) {
	wecrosslogi.WarningDepth(0, "Warning", "depth", "message.")
}

func (s) TestError(t *testing.T) {
	const numErrors = 10
	TLogger.ExpectError("Expected error")
	TLogger.ExpectError("Expected ln error")
	TLogger.ExpectError("Expected formatted error")
	TLogger.ExpectErrorN("Expected repeated error", numErrors)
	wecrosslog.Error("Expected", "error")
	wecrosslog.Errorln("Expected", "ln", "error")
	wecrosslog.Errorf("%v %v %v", "Expected", "formatted", "error")
	for i := 0; i < numErrors; i++ {
		wecrosslog.Error("Expected repeated error")
	}
}
