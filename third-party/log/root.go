package log

import (
	"os"

	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
	osUser "os/user"
	"path/filepath"
)

// Predefined handlers
var (
	root          *logger
	StdoutHandler = StreamHandler(os.Stdout, LogfmtFormat())
	StderrHandler = StreamHandler(os.Stderr, LogfmtFormat())

	CliOutHandler = StreamHandler(colorable.NewColorableStdout(), TerminalFormat())
)

func init() {
	if isatty.IsTerminal(os.Stdout.Fd()) {
		StdoutHandler = StreamHandler(colorable.NewColorableStdout(), TerminalFormat())
	}

	if isatty.IsTerminal(os.Stderr.Fd()) {
		StderrHandler = StreamHandler(colorable.NewColorableStderr(), TerminalFormat())
	}

	root = &logger{[]interface{}{}, new(swapHandler)}
	root.SetHandler(StdoutHandler)
}

// New returns a new logger with the given context.
// New is a convenient alias for Root().New
// 该new是在root的基础上做扩展，因此必须传偶数个参数
func New(ctx ...interface{}) Logger {
	return root.New(ctx...)
}

// Root returns the root logger
func Root() Logger {
	return root
}

// The following functions bypass the exported logger methods (logger.Debug,
// etc.) to keep the call depth the same for all paths to logger.write so
// runtime.Caller(2) always refers to the call site in client code.

// Debug is a convenient alias for Root().Debug
func Debug(msg string, ctx ...interface{}) {
	root.write(msg, LvlDebug, ctx)
}

// Info is a convenient alias for Root().Info
func Info(msg string, ctx ...interface{}) {
	root.write(msg, LvlInfo, ctx)
}

// Warn is a convenient alias for Root().Warn
func Warn(msg string, ctx ...interface{}) {
	root.write(msg, LvlWarn, ctx)
}

// Error is a convenient alias for Root().Error
func Error(msg string, ctx ...interface{}) {
	root.write(msg, LvlError, ctx)
}

// Crit is a convenient alias for Root().Crit
func Crit(msg string, ctx ...interface{}) {
	root.write(msg, LvlCrit, ctx)
}

func InitLogger(logLevel Lvl) {
	Root().SetHandler(LvlFilterHandler(logLevel, StdoutHandler))
}

func homeDir() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	if usr, err := osUser.Current(); err == nil {
		return usr.HomeDir
	}
	return ""
}

// 外边可能要文件输出，可能要控制台输出，可能两个都需要
func InitCsLogger(logLevel Lvl, filePath string, withConsole bool, withFile bool) {

	var targetDir string
	if filePath == "" {
		targetDir = filepath.Join(homeDir(), "tmp", "log", "dipperin")
	} else {
		targetDir = filePath
	}

	if !PathExists(targetDir) {
		os.MkdirAll(targetDir, os.ModePerm)
	}

	var handlers []Handler
	if withFile {
		logFilePath := filepath.Join(targetDir, "dipperin.log")
		fileHandler, err := FileHandler(logFilePath, TerminalFormat())
		if err != nil {
			panic(err.Error())
		}
		Info("write log to file", "path", logFilePath)
		handlers = append(handlers, LvlFilterHandler(logLevel, fileHandler))
	}
	if withConsole {
		handlers = append(handlers, LvlFilterHandler(logLevel, StdoutHandler))
	}

	Root().SetHandler(MultiHandler(handlers...))
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
