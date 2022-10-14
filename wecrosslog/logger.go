package wecrosslog

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/WeBankBlockchain/WeCross-Go-SDK/internal/wecrosslog"
)

// LoggerV1 does underlying logging work for wecrosslog.
type LoggerV1 interface {
	// Info logs to INFO log. Arguments are handled in the manner of fmt.Print.
	Info(args ...interface{})
	// Infoln logs to INFO log. Arguments are handled in the manner of fmt.Println.
	Infoln(args ...interface{})
	// Infof logs to INFO log. Arguments are handled in the manner of fmt.Printf.
	Infof(format string, args ...interface{})
	// Warning logs to WARNING log. Arguments are handled in the manner of fmt.Print.
	Warning(args ...interface{})
	// Warningln logs to WARNING log. Arguments are handled in the manner of fmt.Println.
	Warningln(args ...interface{})
	// Warningf logs to WARNING log. Arguments are handled in the manner of fmt.Printf.
	Warningf(format string, args ...interface{})
	// Error logs to ERROR log. Arguments are handled in the manner of fmt.Print.
	Error(args ...interface{})
	// Errorln logs to ERROR log. Arguments are handled in the manner of fmt.Println.
	Errorln(args ...interface{})
	// Errorf logs to ERROR log. Arguments are handled in the manner of fmt.Printlf.
	Errorf(format string, args ...interface{})
	// Fatal logs to ERROR log. Arguments are handled in the manner of fmt.Print.
	// WeCross ensures that all Fatal logs will exit with os.Exit(1).
	// Implementations may also call os.Exit() with a non-zero exit code.
	Fatal(args ...interface{})
	// Fatalln logs to ERROR log. Arguments are handled in the manner of fmt.Println.
	// WeCross ensures that all Fatal logs will exit with os.Exit(1).
	// Implementations may also call os.Exit() with a non-zero exit code.
	Fatalln(args ...interface{})
	// Fatalf logs to ERROR log. Arguments are handled in the manner of fmt.Printf.
	// WeCross ensures that all Fatal logs will exit with os.Exit(1).
	// Implementations may also call os.Exit() with a non-zero exit code.
	Fatalf(format string, args ...interface{})
	// V reports whether verbosity level l is at least the requested verbose level.
	V(l int) bool
}

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
	if _, ok := l.(*componentData); ok {
		panic("cannot use component logger as wecross logger")
	}
	wecrosslog.Logger = l
	wecrosslog.DepthLogger, _ = l.(wecrosslog.DepthLoggerV1)
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

// newLoggerV1 creates a loggerV1 to be used as default logger.
// All logs are written to stderr.
func newLoggerV1() LoggerV1 {
	errorW := ioutil.Discard
	warningW := ioutil.Discard
	infoW := ioutil.Discard

	logLevel := os.Getenv("WECORSS_GO_LOG_SEVERITY_LEVEL")
	switch logLevel {
	case "", "ERROR", "error": // If env is unset, set level to ERROR.
		errorW = os.Stderr
	case "WARNING", "warning":
		warningW = os.Stderr
	case "INFO", "info":
		infoW = os.Stderr
	}

	var v int
	vLevel := os.Getenv("WECORSS_GO_LOG_VERBOSITY_LEVEL")
	if vl, err := strconv.Atoi(vLevel); err == nil {
		v = vl
	}

	jsonFormat := strings.EqualFold(os.Getenv("WECORSS_GO_LOG_FORMATTER"), "json")

	return newLoggerWithConfig(infoW, warningW, errorW, loggerConfig{
		verbose:    v,
		jsonFormat: jsonFormat,
	})
}

func (g *loggerT) output(severity int, s string) {
	sevStr := severityName[severity]
	if !g.jsonFormat {
		g.m[severity].Output(2, fmt.Sprintf("%v: %v", sevStr, s))
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
