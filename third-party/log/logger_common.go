package log

import (
	"os"
	"path/filepath"
)

type LoggerType uint8

const (
	LoggerConsole LoggerType = iota
	LoggerFile
	LoggerConSoleAndFile
)

func GetInitLogger(loggerType LoggerType, logLevel Lvl, filePath, dirName, logName string, initLog Logger) Logger {
	var l Logger
	if initLog == nil {
		l = &logger{[]interface{}{}, new(swapHandler)}
	} else {
		l = initLog
	}

	var targetDir string
	if filePath == "" {
		targetDir = filepath.Join(homeDir(), "tmp", "log", dirName)
	} else {
		targetDir = filePath
	}

	if !PathExists(targetDir) {
		os.MkdirAll(targetDir, os.ModePerm)
	}

	logFilePath := filepath.Join(targetDir, logName)
	fileHandler, err := FileHandler(logFilePath, TerminalFormat())
	if err != nil {
		panic(err.Error())
	}

	var handlers []Handler
	switch loggerType {
	case LoggerConsole:
		handlers = append(handlers, LvlFilterHandler(logLevel, StdoutHandler))
	case LoggerFile:
		handlers = append(handlers, LvlFilterHandler(logLevel, fileHandler))
	case LoggerConSoleAndFile:
		handlers = append(handlers, LvlFilterHandler(logLevel, StdoutHandler), LvlFilterHandler(logLevel, fileHandler))
	}

	l.SetHandler(MultiHandler(handlers...))
	return l
}
