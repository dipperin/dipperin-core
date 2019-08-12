package mem_manage

import "github.com/dipperin/dipperin-core/third-party/log"

var l log.Logger
var MemoryLogConf = log.LoggerConfig{
	Type:      log.LoggerFile,
	LogLevel:  log.LvlCrit,
	FilePath:  "",
	DirName:   "vm_memory",
	RemoveOld: true,
}

func InitVmMemoryLogger(conf log.LoggerConfig, nodeName string) {
	l = log.SetInitLogger(conf, nodeName)
}

func init() {
	InitVmMemoryLogger(MemoryLogConf, "")
}