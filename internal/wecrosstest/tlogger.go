package wecrosstest

import (
	"os"
	"regexp"
	"strconv"
	"sync"
	"testing"
	"time"
)

// Tlogger serves as the wecrosslog logger and is the interface through which
// expected errors are declared in tests.
var TLogger *tLogger

const callingFrame = 4

type logType int

func (l logType) String() string {
	switch l {
	case infoLog:
		return "INFO"
	case warningLog:
		return "WARNING"
	case errorLog:
		return "ERROR"
	case fatalLog:
		return "FATAL"
	}
	return "UNKNOWN"
}

const (
	infoLog logType = iota
	warningLog
	errorLog
	fatalLog
)

type tLogger struct {
	v           int
	initialized bool

	mu     sync.Mutex // guards t, start, and errors
	t      *testing.T
	start  time.Time
	errors map[*regexp.Regexp]int
}

func init() {
	TLogger = &tLogger{errors: map[*regexp.Regexp]int{}}
	vLevel := os.Getenv("WECROSS_GO_LOG_VERBOSITY_LEVEL")
	if vl, err := strconv.Atoi(vLevel); err == nil {
		TLogger.v = vl
	}
}
