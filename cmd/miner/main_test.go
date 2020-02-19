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
	"github.com/dipperin/dipperin-core/cmd/dipperin/config"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	"github.com/urfave/cli"
)

func Test_main(t *testing.T) {
	go main()
	time.Sleep(10 * time.Millisecond)
}

func Test_action(t *testing.T) {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  masterFlagName,
			Usage: "master info",
		},
		cli.StringFlag{
			Name:  coinbaseFlagName,
			Usage: "coinbase",
		},
		cli.IntFlag{
			Name:  minerCountFlagName,
			Usage: "number of miners",
			Value: 1,
		},
		cli.StringFlag{
			Name:  config.P2PListenerFlagName,
			Usage: "p2p port",
			Value: ":62060",
		},
		cli.BoolFlag{
			Name:  config.NoWalletStartFlagName,
			Usage: "not need to set SoftWalletPasswordFlag SoftWalletPassPhraseFlag SoftWalletPathFlag when this flag is true",
		},
	}
	app.Action = func(c *cli.Context) {
		assert.Panics(t, func() {
			action(c)
		})

		c.Set(masterFlagName, "enode://b832f4f2fe19dbc5604766bbb268a6d0f7ce9ce381b034b262a92f0ad8283a1b5fa058dea5269b66fbb2014a24fa7198c6dc2d8c9cbac7a348258fc20702561f@127.0.0.1:10003")
		c.Set(coinbaseFlagName, "123")
		c.Set(minerCountFlagName, "1")
		c.Set(config.P2PListenerFlagName, "123")
		c.Set(config.NoWalletStartFlagName, "true")

		//fmt.Println(123)
		go action(c)
		time.Sleep(10 * time.Millisecond)
	}
	assert.NoError(t, app.Run([]string{"x"}))
}
