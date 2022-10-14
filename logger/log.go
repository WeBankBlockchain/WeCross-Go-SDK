package logger

import (
	"WeCross-Go-SDK/utils"
	"fmt"
	"io"
	"log"
	"os"
)

func openLogFile(datadir string, filename string) *os.File {
	path := utils.AbsolutePath(datadir, filename)
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Sprintf("error opening log file '%s': %v", filename, err))
	}
	return file
}

func New(datadir string, logFile string, logLevel int) LogSystem {
	var writer io.Writer
	if logFile == "" {
		writer = os.Stdout
	} else {
		writer = openLogFile(datadir, logFile)
	}

	var sys LogSystem
	sys = NewStdLogSystem(writer, log.LstdFlags, LogLevel(logLevel))
	AddLogSystem(sys)

	return sys
}
