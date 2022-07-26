package wecrosslog

import (
	"os"

	"github.com/WeBankBlockchain/WeCross-Go-SDK/internal/wecrosslog"
)

func init() {
	SetLoggerV1(newLoggerV1())
}

// V reports whether verbosity level l is at least the requested verbose level.
func V(l int) bool {
	return wecrosslog.Logger.V(l)
}

// Info logs to the INFO log.
func Info(args ...interface{}) {
	wecrosslog.Logger.Info(args...)
}

// Infof logs to the INFO log. Arguments are handled in the manner of fmt.Printf.
func Infof(format string, args ...interface{}) {
	wecrosslog.Logger.Infof(format, args...)
}

// Infoln logs to the INFO log. Arguments are handled in the manner of fmt.Println.
func Infoln(args ...interface{}) {
	wecrosslog.Logger.Infoln(args...)
}

// Warning logs to the WARNING log.
func Warning(args ...interface{}) {
	wecrosslog.Logger.Warning(args...)
}

// Warningf logs to the WARNING log. Arguments are handled in the manner of fmt.Printf.
func Warningf(format string, args ...interface{}) {
	wecrosslog.Logger.Warningf(format, args...)
}

// Warningln logs to the WARNING log. Arguments are handled in the manner of fmt.Println.
func Warningln(args ...interface{}) {
	wecrosslog.Logger.Warningln(args...)
}

// Error logs to the ERROR log.
func Error(args ...interface{}) {
	wecrosslog.Logger.Error(args...)
}

// Errorf logs to the ERROR log. Arguments are handled in the manner of fmt.Printf.
func Errorf(format string, args ...interface{}) {
	wecrosslog.Logger.Errorf(format, args...)
}

// Errorln logs to the ERROR log. Arguments are handled in the manner of fmt.Println.
func Errorln(args ...interface{}) {
	wecrosslog.Logger.Errorln(args...)
}

// Fatal logs to the FATAL log. Arguments are handled in the manner of fmt.Print.
// It calls os.Exit() with exit code 1.
func Fatal(args ...interface{}) {
	wecrosslog.Logger.Fatal(args...)
	// Make sure fatal logs will exit.
	os.Exit(1)
}

// Fatalf logs to the FATAL log. Arguments are handled in the manner of fmt.Printf.
// It calls os.Exit() with exit code 1.
func Fatalf(format string, args ...interface{}) {
	wecrosslog.Logger.Fatalf(format, args...)
	// Make sure fatal logs will exit.
	os.Exit(1)
}

// Fatalln logs to the FATAL log. Arguments are handled in the manner of fmt.Println.
// It calle os.Exit()) with exit code 1.
func Fatalln(args ...interface{}) {
	wecrosslog.Logger.Fatalln(args...)
	// Make sure fatal logs will exit.
	os.Exit(1)
}

// Print prints to the logger. Arguments are handled in the manner of fmt.Print.
//
// Deprecated: use Info.
func Print(args ...interface{}) {
	wecrosslog.Logger.Info(args...)
}

// Printf prints to the logger. Arguments are handled in the manner of fmt.Printf.
//
// Deprecated: use Infof.
func Printf(format string, args ...interface{}) {
	wecrosslog.Logger.Infof(format, args...)
}

// Println prints to the logger. Arguments are handled in the manner of fmt.Println.
//
// Deprecated: use Infoln.
func Println(args ...interface{}) {
	wecrosslog.Logger.Infoln(args...)
}
