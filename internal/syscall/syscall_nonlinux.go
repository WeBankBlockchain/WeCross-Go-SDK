package syscall

import (
	"net"
	"sync"
	"time"

	"github.com/WeBankBlockchain/WeCross-Go-SDK/wecrosslog"
)

var once sync.Once
var logger = wecrosslog.Component("core")

func log() {
	once.Do(func() {
		logger.Info("CPU time info is unavailable on non-linux environments.")
	})
}

// SetTCPUserTimeout is a no-op function under non-linux environments.
func SetTCPUserTimeout(conn net.Conn, timeout time.Duration) error {
	log()
	return nil
}
