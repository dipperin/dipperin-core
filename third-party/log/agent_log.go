package log

import (
	"os"
    "sync"
	"flag"
)

var agentLog Logger
var mutex sync.Mutex

func InitAgentLog(dataDir string) {
	agentLog = GetAgentLog(LvlDebug, dataDir)
}

func Agent(msg string, ctx ...interface{}) {
    if agentLog != nil {
        agentLog.Info(msg, ctx...)
        return
    }

    mutex.Lock()
    defer mutex.Unlock()

    if IsTestEnv() && agentLog == nil {
    	InitAgentLog(defaultDir)
	}
    agentLog.Info(msg, ctx...)
}

// 判断是否是测试环境
func IsTestEnv() bool {
	return flag.Lookup("test.v") != nil
}

const (
	// 用默认目录，如果一台电脑有多个节点就会炸
	defaultDir  = "/tmp/dipperin/log_agent"
	DefaultPath = defaultDir + "/collect_agent.log"
)

// 这里需要返回一个log 给Dipperin-core 用, 需要在指定目录的生成log文件
func GetAgentLog(logLevel Lvl, path string) Logger {

	var targetDir string

	if path == "" {
		targetDir = defaultDir
	} else {
		targetDir = path
	}

	if !PathExists(targetDir) {
		os.MkdirAll(targetDir, os.ModePerm)
	}
	//log.Info("set agent log to: " + targetDir)
	fileHandler, err := FileHandler(targetDir+"/collect_agent.log", JsonFormat())
	if err != nil {
		panic(err.Error())
	}

	logger := New()
	logger.SetHandler(MultiHandler(LvlFilterHandler(logLevel, fileHandler)))

	return logger
}
