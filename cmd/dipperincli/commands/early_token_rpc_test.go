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

package commands

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"

	"github.com/urfave/cli"
)

func TestMain(m *testing.M) {
	//m.Run()
}

func wrapRpcArgs(c *cli.Context, m string, p string) {
	err := c.Set("m", m)
	if err != nil {
		panic(err)
	}
	err = c.Set("p", p)
	if err != nil {
		panic(err)
	}
}

func addRpcFlags(app *cli.App) {
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "m"},
		cli.StringFlag{Name: "p"},
	}
}

func getRpcTestApp() *cli.App {
	app := cli.NewApp()
	addRpcFlags(app)
	return app
}

func Test_rpcCaller_TransferEDIPToDIP(t *testing.T) {
	app := getRpcTestApp()
	app.Action = func(context *cli.Context) {
		c := &rpcCaller{}
		c.TransferEDIPToDIP(context)

		wrapRpcArgs(context, "TransferEDIPToDIP", "")
		c.TransferEDIPToDIP(context)

		wrapRpcArgs(context, "TransferEDIPToDIP", "0x00005033874289F4F823A896700D94274683535cF0,e,t")
		c.TransferEDIPToDIP(context)

		wrapRpcArgs(context, "TransferEDIPToDIP", "0x00005033874289F4F823A896700D94274683535cF0E1,e,t")
		c.TransferEDIPToDIP(context)

		wrapRpcArgs(context, "TransferEDIPToDIP", "0x00005033874289F4F823A896700D94274683535cF0E1,12,t")
		c.TransferEDIPToDIP(context)

		assert.Panics(t, func() {
			wrapRpcArgs(context, "TransferEDIPToDIP", "0x00005033874289F4F823A896700D94274683535cF0E1,12,2")
			c.TransferEDIPToDIP(context)
		})
	}
	assert.NoError(t, app.Run([]string{os.Args[0]}))
}

func Test_rpcCaller_SetExchangeRate(t *testing.T) {
	app := getRpcTestApp()

	app.Action = func(context *cli.Context) {
		c := &rpcCaller{}
		c.SetExchangeRate(context)

		wrapRpcArgs(context, "SetExchangeRate", "")
		c.SetExchangeRate(context)

		wrapRpcArgs(context, "SetExchangeRate", "0x00005033874289F4F823A896700D94274683535cF0,p,y")
		c.SetExchangeRate(context)

		wrapRpcArgs(context, "SetExchangeRate", "0x00005033874289F4F823A896700D94274683535cF0E1,p,y")
		c.SetExchangeRate(context)

		assert.Panics(t, func() {
			wrapRpcArgs(context, "SetExchangeRate", "0x00005033874289F4F823A896700D94274683535cF0E1,p,1")
			c.SetExchangeRate(context)
		})
	}
	assert.NoError(t, app.Run([]string{os.Args[0]}))
}
