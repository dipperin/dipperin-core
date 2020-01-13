// Copyright 2019, Keychain Foundation Ltd.
// This file is part of the dipperin-core library.
//
// The dipperin-core library is free software: you can redistribute
// it and/or modify it under the terms of the GNU Lesser General Public License
// as published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// The dipperin-core library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package log

import (
	"fmt"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"time"
)

func defaultLogger() *zap.Logger {
	cnf := zap.NewDevelopmentConfig()
	cnf.DisableCaller = true
	cnf.OutputPaths = []string{"stdout"}
	lvl := zapcore.DebugLevel
	cnf.Level = zap.NewAtomicLevelAt(lvl)
	cnf.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	cnf.DisableStacktrace = true
	logger, _ := cnf.Build()
	return logger
}

func NewLogger() *zap.Logger {
	cnf := zap.NewProductionConfig()
	cnf.DisableCaller = true
	cnf.OutputPaths = []string{"stdout"}
	cnf.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	cnf.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, _ := cnf.Build()
	return logger
}

var DLogger *zap.Logger // root

func init() {
	DLogger = defaultLogger()
}

type LoggerConfig struct {
	Lvl           zapcore.Level
	FilePath      string
	Filename      string
	WithConsole   bool
	WithFile      bool
	DisableCaller bool
}

func InitLogger(cnf LoggerConfig) {
	var cores []zapcore.Core

	if cnf.WithFile {
		jsonEncoder := newJSONEncoder(newFileEncoderConfig())

		highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.ErrorLevel
		})
		if cnf.Lvl >= zapcore.ErrorLevel {
			cnf.Lvl = zapcore.WarnLevel
		}
		lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= cnf.Lvl && lvl < zapcore.ErrorLevel
		})

		if cnf.FilePath == "" {
			homeDir, _ := os.UserHomeDir()
			cnf.FilePath = path.Join(homeDir, "tmp", "logs", "dipperin")
		}

		if cnf.Filename == "" {
			cnf.Filename = "dipperin.log"
		}

		out, errOut := getLogFilePath(cnf.FilePath, cnf.Filename)
		logFileOutW := backupsLogWriteSyncer(out)
		logFileErrW := backupsLogWriteSyncer(errOut)

		cores = append(cores,
			zapcore.NewCore(jsonEncoder, logFileErrW, highPriority),
			zapcore.NewCore(jsonEncoder, logFileOutW, lowPriority))
	}

	if cnf.WithConsole {
		priority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= cnf.Lvl && lvl <= zapcore.FatalLevel
		})

		consoleEncoder := newConsoleEncoder(newConsoleEncoderConfig())
		stdoutW := stdoutWriteSyncer()
		cores = append(cores,
			zapcore.NewCore(consoleEncoder, stdoutW, priority))
	}

	DLogger = zap.New(zapcore.NewTee(cores...))
	if !cnf.DisableCaller {
		DLogger = DLogger.WithOptions(zap.AddCaller())
	}
}

func LvlFromString(lv string) (zapcore.Level, error) {
	switch lv {
	case "debug", "dbug":
		return zapcore.DebugLevel, nil
	case "info":
		return zapcore.InfoLevel, nil
	case "warn":
		return zapcore.WarnLevel, nil
	case "error", "eror":
		return zapcore.ErrorLevel, nil
	default:
		return zapcore.DebugLevel, fmt.Errorf("unknown level: %s", lv)
	}
}

/*
#################
#	log file path
#################
*/

func getLogFilePath(targetDir, filename string) (out, errOut string) {
	if filename == "" {
		return
	}

	if !pathExists(targetDir) {
		if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
			panic(err.Error() + "; dir=" + targetDir)
		}
	}

	out = filepath.Join(targetDir, filename)
	errOut = filepath.Join(targetDir, filename+"-err")
	return
}

// Determine if the path file exists
func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

/*
###############
#	log encoder
###############
*/

// json encoder
func newJSONEncoder(cnf zapcore.EncoderConfig) zapcore.Encoder {
	return zapcore.NewJSONEncoder(cnf)
}

// console encoder
func newConsoleEncoder(cnf zapcore.EncoderConfig) zapcore.Encoder {
	return zapcore.NewConsoleEncoder(cnf)
}

// console encoder config
func newConsoleEncoderConfig() zapcore.EncoderConfig {
	cnf := zap.NewProductionEncoderConfig()
	cnf.EncodeDuration = zapcore.SecondsDurationEncoder
	cnf.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
	}
	cnf.EncodeCaller = func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
		var invokeName string
		if _, file, lineNo, ok := runtime.Caller(5); ok {
			invokeName = fmt.Sprintf("%s:%d", file, lineNo)
		}
		enc.AppendString(invokeName)
	}
	cnf.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return cnf
}

// file encoder config
func newFileEncoderConfig() zapcore.EncoderConfig {
	cnf := zap.NewProductionEncoderConfig()
	cnf.EncodeDuration = zapcore.SecondsDurationEncoder
	cnf.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
	}
	cnf.EncodeCaller = func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
		var invokeName string
		if _, file, lineNo, ok := runtime.Caller(5); ok {
			invokeName = fmt.Sprintf("%s:%d", file, lineNo)
		}
		enc.AppendString(invokeName)
	}
	return cnf
}

/*
#####################
#	log writer syncer
#####################
*/

// log files. include Backups
func backupsLogWriteSyncer(filename string) zapcore.WriteSyncer {
	return zapcore.AddSync(&lumberjack.Logger{
		Filename:   filename,
		MaxSize:    500, // megabytes
		MaxBackups: 10,
		MaxAge:     7, // days
		Compress:   true,
	})
}

func stdoutWriteSyncer() zapcore.WriteSyncer {
	return zapcore.Lock(os.Stdout)
}

func stderrWriteSyncer() zapcore.WriteSyncer {
	return zapcore.Lock(os.Stderr)
}
