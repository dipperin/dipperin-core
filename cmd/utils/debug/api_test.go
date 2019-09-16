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

package debug

import (
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/stretchr/testify/assert"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"testing"

	"github.com/urfave/cli"
)

func TestSetup(t *testing.T) {
	app := cli.NewApp()
	app.Flags = Flags
	app.Action = func(c *cli.Context) {
		defer os.RemoveAll(filepath.Join(util.HomeDir(), "tmp", "TestSetup1"))
		defer os.RemoveAll(filepath.Join(util.HomeDir(), "tmp", "TestSetup2"))
		assert.Nil(t, Setup(c))
		assert.Nil(t, c.GlobalSet(traceFlag.Name, "~/tmp/TestSetup1"))
		assert.Nil(t, c.GlobalSet(cpuprofileFlag.Name, "~/tmp/TestSetup2"))
		assert.Nil(t, c.GlobalSet(pprofFlag.Name, "1"))
		Setup(c)
	}
	app.Run([]string{"setup_test"})
}

func TestExit(t *testing.T) {
	Exit()
}
