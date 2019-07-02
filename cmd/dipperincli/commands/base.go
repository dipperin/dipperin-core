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
		Name:    "miner",
		Aliases: []string{"m"},
		Usage:   "miner func",
		Flags:   commonFlags,
		Action: func(c *cli.Context) error {
			RpcCall(c)
			return nil
		},
	},
	{
		Name:  "personal",
		Usage: "personal func",
		Flags: commonFlags,
		Action: func(c *cli.Context) error {
			RpcCall(c)
			return nil
		},
	},
	{
		Name:    "tx",
		Aliases: []string{"t"},
		Usage:   "tx func",
		Flags:   txFlags,
		Action: func(c *cli.Context) error {
			RpcCall(c)
			return nil
		},
	},
	{
		Name:    "verifier",
		Aliases: []string{"v"},
		Usage:   "verifier func",
		Flags:   commonFlags,
		Action: func(c *cli.Context) error {
			RpcCall(c)
			return nil
		},
	},
	{
		Name:    "chain",
		Aliases: []string{"c"},
		Usage:   "chain func",
		Flags:   commonFlags,
		Action: func(c *cli.Context) error {
			RpcCall(c)
			return nil
		},
	},
}

var commonFlags = []cli.Flag{
	cli.StringFlag{Name: "p", Usage: "parameters"},
}

var txFlags = []cli.Flag{
	cli.StringFlag{Name: "p", Usage: "parameters"},
	cli.StringFlag{Name: "abi", Usage: "abi path"},
	cli.StringFlag{Name: "wasm", Usage: "wasm path"},
	cli.StringFlag{Name: "input", Usage: "contract params"},
	cli.BoolFlag{Name: "is-create", Usage: "create contract or not"},
	cli.StringFlag{Name: "func-name", Usage: "call function name"},
}
