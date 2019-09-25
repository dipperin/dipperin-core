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
	"github.com/dipperin/dipperin-core/cmd/base"
	"github.com/dipperin/dipperin-core/cmd/dipperin/config"
	"github.com/dipperin/dipperin-core/core/dipperin"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/urfave/cli"
	"os"
	"os/signal"
	"syscall"
)

const (
	masterFlagName     = "master"
	minerCountFlagName = "m_count"
	coinbaseFlagName   = "coinbase"
)

func main() {
	app := base.NewApp("dipperin miner", "miner for dipperin")
	app.Action = action
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
			Usage: "p2p port",
			Value: ":62060",
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Error("miner run failed", "err", err)
	}
}

func action(c *cli.Context) {
	n, err := dipperin.NewMinerNode(c.String(masterFlagName), c.String(coinbaseFlagName), c.Int(minerCountFlagName), c.String(config.P2PListenerFlagName))
	if err != nil {
		panic("make mine node failed:" + err.Error())
	}

	// listen kill signal
	go signalListen(n)

	// start miner
	n.Start()
	n.Wait()
}

func signalListen(n dipperin.Node) {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	s := <-c
	log.Info("got system signal", "signal", s)
	n.Stop()
}
