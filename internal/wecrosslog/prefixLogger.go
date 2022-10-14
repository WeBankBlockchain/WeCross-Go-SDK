package wecrosslog

import "fmt"

// PrefixLogger does logging with a prefix.
// Logging method on a nil logs without any prefix.
type PrefixLogger struct {
	logger DepthLoggerV1
	prefix string
}

// Infof does info logging.
func (pl *PrefixLogger) Infof(format string, args ...interface{}) {
	if pl != nil {
		// Handle nil, so the tests can pass in a nil logger.
		format = pl.prefix + format
		pl.logger.InfoDepth(1, fmt.Sprintf(format, args...))
		return
	}
	InfoDepth(1, fmt.Sprintf(format, args...))
}

// Warningf does warning logging
func (pl *PrefixLogger) Warningf(format string, args ...interface{}) {
	if pl != nil {
		// Handle nil, so the tests can pass in a nil logger.
		format = pl.prefix + format
		pl.logger.WarningDepth(1, fmt.Sprintf(format, args...))
		return
	}
	WarningDepth(1, fmt.Sprintf(format, args...))
}

// Errorf does error logging
func (pl *PrefixLogger) Errorf(format string, args ...interface{}) {
	if pl != nil {
		// Handle nil, so the tests can pass in a nil logger.
		format = pl.prefix + format
		pl.logger.ErrorDepth(1, fmt.Sprintf(format, args...))
		return
	}
	ErrorDepth(1, fmt.Sprintf(format, args...))
}

// Debugf does error logging at verbose level 2.
func (pl *PrefixLogger) Debugf(format string, args ...interface{}) {
	if !Logger.V(2) {
		return
	}
	if pl != nil {
		// Handle nil, so the tests can pass in a nil logger.
		format = pl.prefix + format
		pl.logger.InfoDepth(1, fmt.Sprintf(format, args...))
		return
	}
	InfoDepth(1, fmt.Sprintf(format, args...))
}

// NewPrefixLogger creates a prefix logger with the given prefix.
func NewPrefixLogger(logger DepthLoggerV1, prefix string) *PrefixLogger {
	return &PrefixLogger{logger: logger, prefix: prefix}
}
