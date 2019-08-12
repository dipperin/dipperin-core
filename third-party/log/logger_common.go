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

type LoggerConfig struct {
	Type     LoggerType
	LogLevel Lvl
	FilePath string
	DirName  string
	RemoveOld bool
}

func SetInitLogger(conf LoggerConfig,nodeName string) Logger {
	l := &logger{[]interface{}{}, new(swapHandler)}
	var targetDir string
	if conf.FilePath == "" {
		targetDir = filepath.Join(homeDir(), "tmp", "dp_debug", conf.DirName)
	} else {
		targetDir = conf.FilePath
	}

	if !PathExists(targetDir) {
		os.MkdirAll(targetDir, os.ModePerm)
	}
	if nodeName == ""{
		nodeName = conf.DirName
	}

	logFilePath := filepath.Join(targetDir, nodeName+".log")
	if conf.RemoveOld {
		_ = os.RemoveAll(logFilePath)
	}
	fileHandler, err := FileHandler(logFilePath, TerminalFormat())
	if err != nil {
		panic(err.Error())
	}

	var handlers []Handler
	switch conf.Type {
	case LoggerConsole:
		handlers = append(handlers, LvlFilterHandler(conf.LogLevel, StdoutHandler))
	case LoggerFile:
		handlers = append(handlers, LvlFilterHandler(conf.LogLevel, fileHandler))
	case LoggerConSoleAndFile:
		handlers = append(handlers, LvlFilterHandler(conf.LogLevel, StdoutHandler), LvlFilterHandler(conf.LogLevel, fileHandler))
	}
	l.SetHandler(MultiHandler(handlers...))
	return l
}
