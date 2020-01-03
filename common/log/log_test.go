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
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"testing"
)

func TestDefaultCsPrint(t *testing.T) {
	DLogger.Debug("test Debug", zap.String("a", "b"))
	DLogger.Info("test Info", zap.String("a", "b"))
	DLogger.Warn("test Warn", zap.String("a", "b"))
	DLogger.Error("test Error", zap.String("a", "b"))
}

func TestDefaultCsPrint_DPainc(t *testing.T) {
	defer func() {
		e := recover()
		fmt.Println(e)
	}()
	DLogger.DPanic("test DPanic", zap.String("a", "b"))
}

func TestDefaultCsPrint_Panic(t *testing.T) {
	defer func() {
		e := recover()
		fmt.Println(e)
	}()
	DLogger.Panic("test Panic", zap.String("a", "b"))
}

func TestDefaultCsPrint_Fatal(t *testing.T) {
	t.Skip()
	defer func() {
		e := recover()
		fmt.Println(e)
	}()
	DLogger.Fatal("test Fatal", zap.String("a", "b"))
}

func TestInitLogger(t *testing.T) {
	cnf := LoggerConfig{
		Lvl:         zapcore.InfoLevel,
		FilePath:    "",
		Filename:    "",
		WithConsole: false,
		WithFile:    false,
	}
	InitLogger(cnf)
	DLogger.Debug("test Debug", zap.String("a", "b"))
	DLogger.Info("test Info", zap.String("a", "b"))
	DLogger.Warn("test Warn", zap.String("a", "b"))
	DLogger.Error("test Error", zap.String("a", "b"))
}

func TestInitLogger_WithConsole_AddCaller(t *testing.T) {
	cnf := LoggerConfig{
		Lvl:           zapcore.InfoLevel,
		FilePath:      "",
		Filename:      "",
		WithConsole:   true,
		WithFile:      false,
		DisableCaller: false,
	}
	InitLogger(cnf)
	DLogger.Debug("test Debug", zap.String("a", "b"))
	DLogger.Info("test Info", zap.String("a", "b"))
	DLogger.Warn("test Warn", zap.String("a", "b"))
	DLogger.Error("test Error", zap.String("a", "b"))
}

func TestInitLogger_WithConsole_DisableCaller(t *testing.T) {
	cnf := LoggerConfig{
		Lvl:           zapcore.InfoLevel,
		FilePath:      "",
		Filename:      "",
		WithConsole:   true,
		WithFile:      false,
		DisableCaller: true,
	}
	InitLogger(cnf)
	DLogger.Debug("test Debug", zap.String("a", "b"))
	DLogger.Info("test Info", zap.String("a", "b"))
	DLogger.Warn("test Warn", zap.String("a", "b"))
	DLogger.Error("test Error", zap.String("a", "b"))
}

func TestInitLogger_WithFile(t *testing.T) {
	cnf := LoggerConfig{
		Lvl:         zapcore.InfoLevel,
		FilePath:    "",
		Filename:    "",
		WithConsole: false,
		WithFile:    true,
	}
	InitLogger(cnf)
	DLogger.Debug("test Debug", zap.String("a", "b"))
	DLogger.Info("test Info", zap.String("a", "b"))
	DLogger.Warn("test Warn", zap.String("a", "b"))
	DLogger.Error("test Error", zap.String("a", "b"))
}

func TestInitLogger_WithConsole_WithFile(t *testing.T) {
	cnf := LoggerConfig{
		Lvl:         zapcore.InfoLevel,
		FilePath:    "",
		Filename:    "",
		WithConsole: true,
		WithFile:    true,
	}
	InitLogger(cnf)
	DLogger.Debug("test Debug", zap.String("a", "b"))
	DLogger.Info("test Info", zap.String("a", "b"))
	DLogger.Warn("test Warn", zap.String("a", "b"))
	DLogger.Error("test Error", zap.String("a", "b"))
}

func TestLvlFromString(t *testing.T) {
	for v, r := range map[string]zapcore.Level{
		"error": zapcore.ErrorLevel,
		"eror":  zapcore.ErrorLevel,
		"warn":  zapcore.WarnLevel,
		"debug": zapcore.DebugLevel,
		"dbug":  zapcore.DebugLevel,
		"info":  zapcore.InfoLevel,
	} {
		lv, _ := LvlFromString(v)
		assert.Equal(t, r, lv)
	}
}
