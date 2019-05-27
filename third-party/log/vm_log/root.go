package vm_log

import (
"github.com/dipperin/dipperin-core/third-party/log"
"github.com/mattn/go-colorable"
"github.com/mattn/go-isatty"
"os"
"path/filepath"
"os/user"
)

// Predefined handlers
var (
	root          log.Logger
	StdoutHandler = log.StreamHandler(os.Stdout, log.LogfmtFormat())
	StderrHandler = log.StreamHandler(os.Stderr, log.LogfmtFormat())
)

func init() {
	if isatty.IsTerminal(os.Stdout.Fd()) {
		StdoutHandler = log.StreamHandler(colorable.NewColorableStdout(), log.TerminalFormat())
	}

	if isatty.IsTerminal(os.Stderr.Fd()) {
		StderrHandler = log.StreamHandler(colorable.NewColorableStderr(), log.TerminalFormat())
	}

	root = log.New()

	// default output nothing
	root.SetHandler(log.LvlFilterHandler(log.LvlInfo, StdoutHandler))
}

// Root returns the root logger
func Root() log.Logger {
	return root
}

// The following functions bypass the exported logger methods (logger.Debug,
// etc.) to keep the call depth the same for all paths to logger.write so
// runtime.Caller(2) always refers to the call site in client code.

// Debug is a convenient alias for Root().Debug
func Debug(msg string, ctx ...interface{}) {
	root.Debug(msg, ctx...)
}

// Info is a convenient alias for Root().Info
func Info(msg string, ctx ...interface{}) {
	root.Info(msg, ctx...)
}

// Warn is a convenient alias for Root().Warn
func Warn(msg string, ctx ...interface{}) {
	root.Warn(msg, ctx...)
}

// Error is a convenient alias for Root().Error
func Error(msg string, ctx ...interface{}) {
	root.Error(msg, ctx...)
}

// Crit is a convenient alias for Root().Crit
func Crit(msg string, ctx ...interface{}) {
	root.Crit(msg, ctx...)
}

func homeDir() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	if usr, err := user.Current(); err == nil {
		return usr.HomeDir
	}
	return ""
}

// 外边可能要文件输出，可能要控制台输出，可能两个都需要
func InitHaltLogger(logLevel log.Lvl, nodeName string, removeOld bool) {
	targetDir := filepath.Join(homeDir(), "tmp", "cs_debug", "halt")

	if !log.PathExists(targetDir) {
		_ = os.MkdirAll(targetDir, os.ModePerm)
	}

	var handlers []log.Handler
	logFilePath := filepath.Join(targetDir, nodeName + ".log")

	if removeOld {
		_ = os.RemoveAll(logFilePath)
	}

	fileHandler, err := log.FileHandler(logFilePath, log.LogfmtFormat())
	if err != nil {
		panic(err)
	}
	log.Debug("write pbft debug log to file", "path", logFilePath)
	handlers = append(handlers, log.LvlFilterHandler(logLevel, fileHandler))

	Root().SetHandler(log.MultiHandler(handlers...))
}

