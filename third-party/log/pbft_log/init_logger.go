package pbft_log

import (
	"github.com/dipperin/dipperin-core/third-party/log"
	"os"
)

var Log log.Logger

var LogConf = log.LoggerConfig{
	Type:      log.LoggerFile,
	LogLevel:  log.LvlDebug,
	FilePath:  "",
	DirName:   "pbft",
	RemoveOld: true,
}

func InitPBFTLogger(conf log.LoggerConfig, nodeName string) {
	if os.Getenv("boots_env") == "mercury"{
		LogConf.LogLevel = log.LvlWarn
	}
	Log = log.SetInitLogger(conf, nodeName)
}

func init() {
	InitPBFTLogger(LogConf, "")
}
