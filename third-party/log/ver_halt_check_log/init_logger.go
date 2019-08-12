package ver_halt_check_log

import "github.com/dipperin/dipperin-core/third-party/log"

var Log log.Logger

var VertHaltLogConf = log.LoggerConfig{
	Type:      log.LoggerFile,
	LogLevel:  log.LvlInfo,
	FilePath:  "",
	DirName:   "ver_halt",
	RemoveOld: true,
}

func InitVerHaltLogger(conf log.LoggerConfig, nodeName string) {
	Log = log.SetInitLogger(conf, nodeName)
}

func init() {
	InitVerHaltLogger(VertHaltLogConf, "")
}
