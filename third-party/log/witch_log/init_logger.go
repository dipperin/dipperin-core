package witch_log

import (
	"github.com/dipperin/dipperin-core/third-party/log"
)

var Log log.Logger

var LogConf = log.LoggerConfig{
	Type:      log.LoggerConsole,
	LogLevel:  log.LvlCrit,
	FilePath:  "",
	DirName:   "witch",
	RemoveOld: true,
}

func InitWitchLogger(conf log.LoggerConfig, nodeName string) {
	Log = log.SetInitLogger(conf, nodeName)
}

func init() {
	InitWitchLogger(LogConf, "")
}
