package log

import (
	"os"
	"path/filepath"
)

const (
	LoggerConsole LoggerType = iota
	LoggerFile
	LoggerConSoleAndFile
)

type LoggerType uint8

type LoggerConfig struct {
	Type      LoggerType
	LogLevel  Lvl
	FilePath  string
	DirName   string
	RemoveOld bool
}

type dpLogger struct {
	Logger
	conf LoggerConfig
}

func DefaultDpLogger(dirName string) *dpLogger {
	conf := DefaultLogConf
	conf.DirName = dirName
	return &dpLogger{
		Logger: SetInitLogger(conf, ""),
		conf:   conf,
	}
}

var (
	DefaultLogConf = LoggerConfig{
		Type:      LoggerConsole,
		LogLevel:  LvlInfo,
		FilePath:  "",
		DirName:   "",
		RemoveOld: true,
	}
)

//dipperIn logger
var (
	Mpt       *dpLogger
	Halt      *dpLogger
	Health    *dpLogger
	PBft      *dpLogger
	Witch     *dpLogger
	Vm        *dpLogger
	VmMem     *dpLogger
	Pm        *dpLogger
	dpLoggers map[string]*dpLogger
)

func init() {
	Mpt = DefaultDpLogger("mpt")
	Halt = DefaultDpLogger("ver_halt")
	Health = DefaultDpLogger("health_info")
	PBft = DefaultDpLogger("PBft")
	Witch = DefaultDpLogger("witch")
	Vm = DefaultDpLogger("vm")
	VmMem = DefaultDpLogger("vm_memory")
	Pm = DefaultDpLogger("pm")

	dpLoggers = map[string]*dpLogger{
		"mpt":         Mpt,
		"ver_halt":    Halt,
		"health_info": Health,
		"PBft":        PBft,
		"witch":       Witch,
		"vm":          Vm,
		"vm_memory":   VmMem,
		"pm":          Pm,
	}
}

func SetInitLogger(conf LoggerConfig, nodeName string) Logger {
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
	if nodeName == "" {
		nodeName = conf.DirName
	} else {
		//remove init creat log file
		file := filepath.Join(targetDir, conf.DirName+".log")
		os.Remove(file)
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

func InitDPLogger(nodeName string) {
	for k, v := range dpLoggers {
		v.conf.Type = LoggerFile
		if os.Getenv("boots_env") == "venus" {
			switch k {
			case "mpt", "health_info", "vm_memory", "witch":
				v.conf.LogLevel = LvlWarn
			}
		}
		v.Logger = SetInitLogger(v.conf, nodeName)
	}
}

func OutputHealthLog() bool {
	if os.Getenv("boots_env") == "venus" {
		return false
	}
	return true
}
