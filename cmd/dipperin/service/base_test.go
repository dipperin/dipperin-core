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
	"fmt"
	"github.com/dipperin/dipperin-core/cmd/dipperin/config"
	"github.com/dipperin/dipperin-core/cmd/utils"
	"github.com/dipperin/dipperin-core/cmd/utils/debug"
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
	"os"
	"testing"
)

func TestStartNode(t *testing.T) {
	app := cli.NewApp()
	app.Flags = append(config.Flags, debug.Flags...)
	app.Action = func(c *cli.Context) {
		assert.NoError(t, os.Setenv("cslog", "enable"))
		assert.NoError(t, c.Set(config.DataDirFlagName, "/tmp/test_start_node"))
		assert.NoError(t, c.Set(config.LogLevelFlagName, "haha"))
		assert.NoError(t, c.Set(config.P2PListenerFlagName, "16888"))
		assert.NoError(t, c.Set(config.HttpPortFlagName, "16889"))
		assert.NoError(t, c.Set(config.WsPortFlagName, "16887"))
		dataDir := c.String(config.DataDirFlagName)
		assert.Equal(t, "/tmp/test_start_node", dataDir)
		defer os.RemoveAll(dataDir)
		utils.SetupGenesis(dataDir, chain_config.GetChainConfig())

		n, err := StartNode(c, true, true, false)
		assert.NoError(t, err)
		assert.NotNil(t, n)
		n.Stop()

		assert.NoError(t, c.Set(config.IsStartMine, "1"))
		assert.NoError(t, c.Set(config.NodeTypeFlagName, fmt.Sprintf("%v", chain_config.NodeTypeOfMineMaster)))
		n, err = StartNode(c, true, true, false)
		assert.Error(t, err)
		assert.Nil(t, n)
	}
	assert.NoError(t, app.Run([]string{os.Args[0]}))
}
