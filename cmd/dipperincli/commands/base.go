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
	"github.com/urfave/cli"
	"os"
)

var CliCommands = []cli.Command{
	{
		Name:    "quit",
		Aliases: []string{"exit"},
		Usage:   "quit",
		Action: func(c *cli.Context) error {
			os.Exit(0)
			return cli.NewExitError("", 0)
		},
		Hidden:   false,
		HideHelp: false,
	},
	{
		Name:    "rpc",
		Aliases: []string{"r"},
		Usage:   "control node",
		Flags:   rpcFlags,
		Action: func(c *cli.Context) error {
			RpcCall(c)
			return nil
		},
	},
}

var rpcFlags = []cli.Flag{
	cli.StringFlag{Name: "m", Usage: "operation"},
	cli.StringFlag{Name: "p", Usage: "parameters"},
	cli.StringFlag{Name: "abi", Usage:"abi path"},
	cli.StringFlag{Name: "wasm", Usage:"wasm path"},
	cli.StringFlag{Name: "input", Usage: "contract params"},
	cli.BoolFlag{Name:   "isCreate", Usage: "create contract or not"},
	cli.StringFlag{Name: "funcName", Usage: "call function name"},
}
