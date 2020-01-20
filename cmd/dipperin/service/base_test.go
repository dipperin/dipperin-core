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

package service

import (
	"github.com/dipperin/dipperin-core/cmd/dipperin/config"
	"github.com/dipperin/dipperin-core/cmd/utils"
	"github.com/dipperin/dipperin-core/cmd/utils/debug"
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/core/chainconfig"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
	"os"
	"path/filepath"
	"testing"
)

func TestStartNode(t *testing.T) {
	app := cli.NewApp()
	app.Flags = append(config.Flags, debug.Flags...)
	app.Action = func(c *cli.Context) {
		//assert.NoError(t, os.Setenv("cslog", "enable"))
		path := filepath.Join(util.HomeDir(), "tmp/test_start_node")
		assert.NoError(t, c.Set(config.DataDirFlagName, path))
		assert.NoError(t, c.Set(config.LogLevelFlagName, "haha"))
		assert.NoError(t, c.Set(config.P2PListenerFlagName, "16888"))
		assert.NoError(t, c.Set(config.HttpPortFlagName, "16889"))
		assert.NoError(t, c.Set(config.WsPortFlagName, "16887"))
		assert.NoError(t, c.Set(config.NoWalletStartFlagName, "true"))
		dataDir := c.String(config.DataDirFlagName)
		assert.Equal(t, path, dataDir)
		defer os.RemoveAll(dataDir)
		utils.SetupGenesis(dataDir, chainconfig.GetChainConfig())

		n, err := StartNode(c, true, true, false)
		assert.NoError(t, err)
		assert.NotNil(t, n)
		n.Stop()
	}
	assert.NoError(t, app.Run([]string{os.Args[0]}))
}
