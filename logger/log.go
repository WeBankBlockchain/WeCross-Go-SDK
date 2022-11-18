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

// AddStdOutLogSystem quickly add an stdOut log system with the given log level
func AddStdOutLogSystem(logLevel int) LogSystem {
	var sys LogSystem
	sys = NewStdLogSystem(os.Stdout, log.LstdFlags, LogLevel(logLevel))
	AddLogSystem(sys)
	return sys
}

// AddNewLogSystem can customize a log system with the given args and add it.
// If logFile is empty, will use stdOut as the log output.
func AddNewLogSystem(datadir string, logFile string, flags int, logLevel int) LogSystem {
	var writer io.Writer
	if logFile == "" {
		writer = os.Stdout
	} else {
		writer = openLogFile(datadir, logFile)
	}

	var sys LogSystem
	sys = NewStdLogSystem(writer, flags, LogLevel(logLevel))
	AddLogSystem(sys)

	return sys
}
