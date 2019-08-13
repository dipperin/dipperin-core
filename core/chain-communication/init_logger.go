package chain_communication

import "github.com/dipperin/dipperin-core/third-party/log"

var pmLog log.Logger

var PMLogConf = log.LoggerConfig{
	Type:      log.LoggerFile,
	LogLevel:  log.LvlInfo,
	FilePath:  "",
	DirName:   "pm",
	RemoveOld: true,
}

func InitPMLogger(conf log.LoggerConfig, nodeName string) {
	pmLog = log.SetInitLogger(conf, nodeName)
}

func init() {
	InitPMLogger(PMLogConf, "")
}
