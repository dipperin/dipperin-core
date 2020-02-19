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

package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/dipperin/dipperin-core/common/util"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

func Test_main(t *testing.T) {
	main()
	time.Sleep(1 * time.Millisecond)
}

func Test_newApp(t *testing.T) {
	assert.NotNil(t, newApp())
	os.RemoveAll(filepath.Join(util.HomeDir(), ".dipperin", "start_conf.json"))
}

func Test_initStartFlag(t *testing.T) {
	assert.NotNil(t, initStartFlag())
	os.RemoveAll(filepath.Join(util.HomeDir(), ".dipperin", "start_conf.json"))
}

func Test_doPrompts(t *testing.T) {
	startConfPath := filepath.Join(util.HomeDir(), ".dipperin", "start_conf.json")
	defer os.RemoveAll(startConfPath)

	var conf startConf
	doPrompts(&conf, startConfPath)
}

func Test_appAction(t *testing.T) {
	app := newApp()

	assert.Panics(t, func() {
		appAction(cli.NewContext(app, nil, nil))
	})
}

func Test_haveCmd(t *testing.T) {
	assert.Equal(t, true, haveCmd("quit"))
	assert.Equal(t, false, haveCmd("ssss"))
}

func TestExecutor(t *testing.T) {
	assert.NotNil(t, Executor(nil))
}

func Test_startNode(t *testing.T) {
	assert.Panics(t, func() {
		startNode(nil)
	})
}

func testTempFile(t *testing.T) (string, func()) {
	t.Helper()
	filePath := filepath.Join(util.HomeDir(), ".dipperin", "registration")
	tf, err := os.Create(filePath)
	if err != nil {
		t.Fatalf(":err: %s", err)
	}

	tf.Close()
	return tf.Name(), func() { os.Remove(tf.Name()) }
}

func Test_getNodeType(t *testing.T) {

	assert.Equal(t, getNodeType(2), "verifier")
	assert.Equal(t, getNodeType(1), "mine master")
	assert.Equal(t, getNodeType(0), "normal")
	_, tfClean := testTempFile(t)
	defer tfClean()
	assert.Equal(t, getNodeType(1), "")
	assert.Equal(t, getNodeType(0), "")
}
