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
	"github.com/dipperin/dipperin-core/cmd/dipperin/service"
	"github.com/dipperin/dipperin-core/cmd/utils/debug"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/ethereum/go-ethereum/metrics"
	"github.com/urfave/cli"
	"os"
	"time"
)

func main() {
	log.Info("~~~~~~~~~start app ~~~~~~~~~~~~")
	app := base.NewApp("dipperin", "dipperin node and console")
	app.Flags = append(config.Flags, debug.Flags...)
	app.Action = func(c *cli.Context) error {
		//use pprof
		debug.Setup(c)

		// Start system runtime metrics collection
		go metrics.CollectProcessMetrics(3 * time.Second)

		_, err := service.StartNode(c, false, true, false)
		return err
	}
	if err := app.Run(os.Args); err != nil {
		panic("run dipperin failed: " + err.Error())
	}
}
