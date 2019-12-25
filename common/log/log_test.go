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
	"go.uber.org/zap"
	"os"
	"path"
	"sync"
	"testing"
	"time"
)

func TestDefaultPrint(t *testing.T) {
	DLogger.Debug("test Debug", zap.String("a", "b"))
	DLogger.Info("test Info", zap.String("a", "b"))
	DLogger.Warn("test Warn", zap.String("a", "b"))
	DLogger.Error("test Error", zap.String("a", "b"))
}

func TestDefaultPrint_DPainc(t *testing.T) {
	defer func() {
		e := recover()
		fmt.Println(e)
	}()
	DLogger.DPanic("test DPanic", zap.String("a", "b"))
}

func TestDefaultPrint_Panic(t *testing.T) {
	defer func() {
		e := recover()
		fmt.Println(e)
	}()
	DLogger.Panic("test Panic", zap.String("a", "b"))
}

func TestDefaultPrint_Fatal(t *testing.T) {
	t.Skip()
	defer func() {
		e := recover()
		fmt.Println(e)
	}()
	DLogger.Fatal("test Fatal", zap.String("a", "b"))
}

func TestInitLogger(t *testing.T) {
	InitLogger(false, "", "")
	DLogger.Debug("test Debug", zap.String("a", "b"))
	DLogger.Info("test Info", zap.String("a", "b"))
	DLogger.Warn("test Warn", zap.String("a", "b"))
	DLogger.Error("test Error", zap.String("a", "b"))
}

func TestInitLogger_debug(t *testing.T) {
	InitLogger(true, "", "")
	DLogger.Debug("test Debug", zap.String("a", "b"))
	DLogger.Info("test Info", zap.String("a", "b"))
	DLogger.Warn("test Warn", zap.String("a", "b"))
	DLogger.Error("test Error", zap.String("a", "b"))
}

func TestInitLogger_LogFile(t *testing.T) {
	testDir, err := os.UserHomeDir()
	if err != nil {
		t.Skip()
	}
	testDir = path.Join(testDir, "tmp")
	InitLogger(false, testDir, fmt.Sprintf("%d.log", time.Now().Unix()))
	DLogger.Debug("test Debug", zap.String("a", "b"))
	DLogger.Info("test Info", zap.String("a", "b"))
	DLogger.Warn("test Warn", zap.String("a", "b"))
	DLogger.Error("test Error", zap.String("a", "b"))
}

func TestInitLogger_LogFile_debug(t *testing.T) {
	testDir, err := os.UserHomeDir()
	if err != nil {
		t.Skip()
	}
	testDir = path.Join(testDir, "tmp")
	InitLogger(true, testDir, fmt.Sprintf("%d.log", time.Now().Unix()))
	DLogger.Debug("test Debug", zap.String("a", "b"))
	DLogger.Info("test Info", zap.String("a", "b"))
	DLogger.Warn("test Warn", zap.String("a", "b"))
	DLogger.Error("test Error", zap.String("a", "b"))
}

func TestInitLogger_backups(t *testing.T) {
	t.Skip("Backups: only local")
	testDir, err := os.UserHomeDir()
	if err != nil {
		t.Skip()
	}
	testDir = path.Join(testDir, "tmp", "logs")
	InitLogger(true, testDir, fmt.Sprintf("%d.log", time.Now().Unix()))
	sw := sync.WaitGroup{}
	bt := time.Now()
	for i := 0; i < (1 << 20); i++ {
		i := i
		sw.Add(1)
		go func() {
			defer sw.Done()
			DLogger.Debug("test Debug", zap.String("a", "b"), zap.Int("count", i))
			DLogger.Info("test Info", zap.String("a", "b"), zap.Int("count", i))
			DLogger.Warn("test Warn", zap.String("a", "b"), zap.Int("count", i))
			DLogger.Error("test Error", zap.String("a", "b"), zap.Int("count", i))
		}()
	}
	sw.Wait()
	fmt.Println(time.Since(bt))
}
