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

package dipperin_console

import (
	"github.com/dipperin/dipperin-core/common/util"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/dipperin/dipperin-core/cmd/dipperincli/config"
	"github.com/stretchr/testify/assert"
	"runtime"
)

func TestGetConfigDir(t *testing.T) {
	assert.NotEqual(t, "", GetConfigDir())

	oldConfigDir := util.ExecutablePath()
	cFile := filepath.Join(oldConfigDir, ConfigName)
	defer os.RemoveAll(cFile)
	assert.NoError(t, ioutil.WriteFile(cFile, []byte("test"), 0644))

	GetConfigDir()
}

func TestGetConfigDir1(t *testing.T) {
	assert.NoError(t, os.Setenv(EnvConfigDir, "/tmp"))

	assert.NotEqual(t, "", GetConfigDir())
}

func TestGetConfigDir2(t *testing.T) {
	assert.NoError(t, os.Setenv(EnvConfigDir, "./config"))

	assert.NotEqual(t, "", GetConfigDir())
}

func TestNewConsole(t *testing.T) {
	//c := NewConsole(func(command string) {}, config.DipperinCliCompleter)
	//assert.NotNil(t, c)
	//c.WrapExecutor(func(s string) {})("123")
}

func TestNewConsole1(t *testing.T) {
	if runtime.GOOS != "linux" {
		return
	}

	historyFilePath = "/tmp/aaa/cs_command_history.txt"
	os.RemoveAll(filepath.Dir(historyFilePath))
	assert.Panics(t, func() {
		NewConsole(func(command string) {}, config.DipperinCliCompleter)
	})
}

func TestGetWinConfigDir(t *testing.T) {
	getWinConfigDir("")
	assert.NoError(t, os.Setenv("APPDATA", "/home/x"))
	getWinConfigDir("")
}
