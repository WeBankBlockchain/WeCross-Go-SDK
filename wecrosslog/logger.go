package wecrosslog

import (
	"encoding/json"
	"fmt"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/internal/wecrosslog"
	"io"
	"log"
	"os"
)

const (
	// infoLog indicates Info severity.
	infoLog int = iota
	// waringLog indicates Waring severity.
	warningLog
	// errorLog indicates Error severity.
	errorLog
	// fatalLog indicates Fatal severity.
	fatalLog
)

// severityName contains the string representation of each severity.
var severityName = []string{
	infoLog:    "INFO",
	warningLog: "WARNING",
	errorLog:   "ERROR",
	fatalLog:   "FATAL",
}

// loggerT is the default logger used by wecrosslog.
type loggerT struct {
	m          []*log.Logger
	v          int
	jsonFormat bool
}

// SetLoggerV1 sets logger that is used in wecross to a V1 logger.
// Not mutex-protected, should be called before any wecross functions.
func SetLoggerV1(l wecrosslog.LoggerV1) {
	wecrosslog.Logger = l
}

// NewLogger creates a logger with the provided writers.
// Fatal logs will be written to errorW, warningW, infoW, followed by exit(1).
// Error logs will be written to errorW, warningW and infoW.
// Warning logs will be written to warningW and infoW.
// Info logs will be written to infoW.
func NewLogger(infoW, waringW, errorW io.Writer) wecrosslog.LoggerV1 {
	return newLoggerWithConfig(infoW, waringW, errorW, loggerConfig{})
}

// NewLoggerWithVerbosity creates a logger with the provided writers and
// verbosity level.
func NewLoggerWithVerbosity(infoW, warningW, errorW io.Writer, v int) wecrosslog.LoggerV1 {
	return newLoggerWithConfig(infoW, warningW, errorW, loggerConfig{verbose: v})
}

type loggerConfig struct {
	verbose    int
	jsonFormat bool
}

func newLoggerWithConfig(infoW, waringW, errorW io.Writer, c loggerConfig) wecrosslog.LoggerV1 {
	var m []*log.Logger
	flag := log.LstdFlags
	if c.jsonFormat {
		flag = 0
	}
	m = append(m, log.New(infoW, "", flag))
	m = append(m, log.New(io.MultiWriter(infoW, waringW), "", flag))
	ew := io.MultiWriter(infoW, waringW, errorW) // errorW will be used for error and fatal.
	m = append(m, log.New(ew, "", flag))
	m = append(m, log.New(ew, "", flag))
	return &loggerT{m: m, v: c.verbose, jsonFormat: c.jsonFormat}
}

func (g *loggerT) output(severity int, s string) {
	sevStr := severityName[severity]
	if !g.jsonFormat {
		_ = g.m[severity].Output(2, fmt.Sprintf("%v: %v", sevStr, s)).Error()
		return
	}
	b, _ := json.Marshal(map[string]string{
		"severity": sevStr,
		"message":  s,
	})
	_ = g.m[severity].Output(2, string(b)).Error()
}

func (g *loggerT) Info(args ...interface{}) {
	g.output(infoLog, fmt.Sprint(args...))
}

func (g *loggerT) Infoln(args ...interface{}) {
	g.output(infoLog, fmt.Sprintln(args...))
}

func (g *loggerT) Infof(format string, args ...interface{}) {
	g.output(infoLog, fmt.Sprintf(format, args...))
}

func (g *loggerT) Warning(args ...interface{}) {
	g.output(warningLog, fmt.Sprint(args...))
}

func (g *loggerT) Warningln(args ...interface{}) {
	g.output(warningLog, fmt.Sprintln(args...))
}

func (g *loggerT) Warningf(format string, args ...interface{}) {
	g.output(warningLog, fmt.Sprintf(format, args...))
}

func (g *loggerT) Error(args ...interface{}) {
	g.output(errorLog, fmt.Sprint(args...))
}

func (g *loggerT) Errorln(args ...interface{}) {
	g.output(errorLog, fmt.Sprintln(args...))
}

func (g *loggerT) Errorf(format string, args ...interface{}) {
	g.output(errorLog, fmt.Sprintf(format, args...))
}

func (g *loggerT) Fatal(args ...interface{}) {
	g.output(fatalLog, fmt.Sprint(args...))
	os.Exit(1)
}

func (g *loggerT) Fatalln(args ...interface{}) {
	g.output(fatalLog, fmt.Sprintln(args...))
	os.Exit(1)
}

func (g *loggerT) Fatalf(format string, args ...interface{}) {
	g.output(fatalLog, fmt.Sprintf(format, args...))
	os.Exit(1)
}

func (g *loggerT) V(l int) bool {
	return l <= g.v
}
