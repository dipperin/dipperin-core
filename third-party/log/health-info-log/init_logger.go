package health_info_log

import (
	"github.com/dipperin/dipperin-core/third-party/log"
	"os"
)

var Log log.Logger

var LogConf = log.LoggerConfig{
	Type:      log.LoggerFile,
	LogLevel:  log.LvlDebug,
	FilePath:  "",
	DirName:   "health_info",
	RemoveOld: true,
}

func InitHealthInfoLogger(conf log.LoggerConfig, nodeName string) {
	Log = log.SetInitLogger(conf, nodeName)
}

func init() {
	if os.Getenv("boots_env") == "mercury" {
		LogConf.LogLevel = log.LvlWarn
	}
	InitHealthInfoLogger(LogConf, "")
}

func OutputHealthLog() bool {
	if os.Getenv("boots_env") == "mercury" {
		return false
	}
	return true
}
