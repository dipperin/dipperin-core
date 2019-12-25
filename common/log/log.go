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
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"time"
)

var (
	DLogger *zap.Logger // common

	Mpt        *zap.Logger
	Halt       *zap.Logger
	Health     *zap.Logger
	PBft       *zap.Logger
	Witch      *zap.Logger
	Vm         *zap.Logger
	VmMem      *zap.Logger
	Pm         *zap.Logger
	Middleware *zap.Logger
	P2P        *zap.Logger
	Stack      *zap.Logger
	Rpc        *zap.Logger
)

func init() {
	DLogger = defaultLogger()
	Mpt = defaultLogger()
	Halt = defaultLogger()
	Health = defaultLogger()
	PBft = defaultLogger()
	Witch = defaultLogger()
	Vm = defaultLogger()
	VmMem = defaultLogger()
	Pm = defaultLogger()
	Middleware = defaultLogger()
	P2P = defaultLogger()
	Stack = defaultLogger()
	Rpc = defaultLogger()
}

func defaultLogger() *zap.Logger {
	cnf := zap.NewDevelopmentConfig()
	cnf.DisableCaller = true
	cnf.OutputPaths = []string{"stdout"}
	cnf.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	logger, _ := cnf.Build()
	return logger
}

// init logger
func InitLogger(debug bool, targetDir, filename string) {
	DLogger = newLogger(debug, path.Join(targetDir, "dipperin"), filename)
	Mpt = newLogger(debug, path.Join(targetDir, "mpt"), filename)
	Halt = newLogger(debug, path.Join(targetDir, "halt"), filename)
	Health = newLogger(debug, path.Join(targetDir, "health"), filename)
	PBft = newLogger(debug, path.Join(targetDir, "pbft"), filename)
	Witch = newLogger(debug, path.Join(targetDir, "witch"), filename)
	Vm = newLogger(debug, path.Join(targetDir, "vm"), filename)
	VmMem = newLogger(debug, path.Join(targetDir, "vmmem"), filename)
	Pm = newLogger(debug, path.Join(targetDir, "pm"), filename)
	Middleware = newLogger(debug, path.Join(targetDir, "middleware"), filename)
	P2P = newLogger(debug, path.Join(targetDir, "p2p"), filename)
	Stack = newLogger(debug, path.Join(targetDir, "stack"), filename)
	Rpc = newLogger(debug, path.Join(targetDir, "rpc"), filename)
}

func newLogger(debug bool, targetDir, filename string) *zap.Logger {
	return zap.New(
		newLogCore(debug, targetDir, filename),
		newLogOptions(debug)...,
	)
}

func newLogOptions(debug bool) []zap.Option {
	if debug {
		return []zap.Option{
			// print stack messages
			zap.AddStacktrace(zapcore.ErrorLevel),
		}
	}
	return nil
}

func newLogCore(debug bool, targetDir, filename string) zapcore.Core {
	encoderConfig := newLogEncoderConfig()
	jsonEncoder := newJSONEncoder(encoderConfig)
	consoleEncoder := newConsoleEncoder(encoderConfig)
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel && lvl < zapcore.ErrorLevel
	})
	if debug {
		lowPriority = func(lvl zapcore.Level) bool {
			return lvl < zapcore.ErrorLevel
		}
	}

	stdoutW := stdoutWriteSyncer()
	stderrW := stderrWriteSyncer()

	out, errOut := getLogFilePath(targetDir, filename)
	logFileOutW := backupsLogWriteSyncer(out)
	logFileErrW := backupsLogWriteSyncer(errOut)

	return zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, stderrW, highPriority),
		zapcore.NewCore(consoleEncoder, stdoutW, lowPriority),
		zapcore.NewCore(jsonEncoder, logFileErrW, highPriority),
		zapcore.NewCore(jsonEncoder, logFileOutW, lowPriority))
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

// encoder config
func newLogEncoderConfig() zapcore.EncoderConfig {
	cnf := zap.NewProductionEncoderConfig()
	cnf.EncodeDuration = zapcore.SecondsDurationEncoder
	cnf.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
	}
	cnf.EncodeCaller = func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
		var invokeName string
		if pc, _, lineNo, ok := runtime.Caller(6); ok {
			invokeName = runtime.FuncForPC(pc).Name() + ":" + strconv.FormatInt(int64(lineNo), 10)
		}
		enc.AppendString(os.Getenv("HOSTNAME") + "	" + invokeName)
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
