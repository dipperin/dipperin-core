package mpt_log

import (
	"github.com/dipperin/dipperin-core/third-party/log"
	"os"
)

var Log log.Logger

var LogConf = log.LoggerConfig{
	Type:      log.LoggerFile,
	LogLevel:  log.LvlDebug,
	FilePath:  "",
	DirName:   "mpt",
	RemoveOld: true,
}

func InitMptLogger(conf log.LoggerConfig, nodeName string) {
	if os.Getenv("boots_env") == "venus" {
		LogConf.LogLevel = log.LvlWarn
	}
	Log = log.SetInitLogger(conf, nodeName)
}

func init() {
	InitMptLogger(LogConf, "")
}
