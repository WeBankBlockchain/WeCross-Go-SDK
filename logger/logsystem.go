package logger

import (
	"io"
	"log"
	"sync/atomic"
)

// LogSystem is implemented by log output devices.
// All methods can be called concurrently from multiple goroutines.
type LogSystem interface {
	LogPrint(LogMsg)
}

// NewStdLogSystem creates a LogSystem that prints to the given writer.
// The flag values are defined package log.
func NewStdLogSystem(writer io.Writer, flags int, level LogLevel) *StdLogSystem {
	logger := log.New(writer, "", flags)
	return &StdLogSystem{logger, uint32(level)}
}

type StdLogSystem struct {
	logger *log.Logger
	level  uint32
}

func (t *StdLogSystem) LogPrint(msg LogMsg) {
	stdmsg, ok := msg.(stdMsg)
	if ok {
		if t.GetLogLevel() >= stdmsg.Level() {
			t.logger.Print(stdmsg.String())
		}
	}
}

func (t *StdLogSystem) SetLogLevel(i LogLevel) {
	atomic.StoreUint32(&t.level, uint32(i))
}

func (t *StdLogSystem) GetLogLevel() LogLevel {
	return LogLevel(atomic.LoadUint32(&t.level))
}
