package commands

import (
	"errors"
	"github.com/dipperin/dipperin-core/common"
	"github.com/urfave/cli"
)

var (
	from         = "0x00005586B883Ec6dd4f8c26063E18eb4Bd228e59c3E9"
	to           = "0x0000970e8128aB834E8EAC17aB8E3812f010678CF791"
	contractAddr = "0x0014B5Df12F50295469Fe33951403b8f4E63231Ef488"
	txHash       = "0xeb7dd095c6339b9f64e6dc5677a371adf3629f0261e49c79a3d4dd7a5c1bfc1a"
	fromAddr     = common.HexToAddress(from)
	testErr      = errors.New("test error")
)

func addRpcFlags(app *cli.App) {
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "p", Usage: "parameters"},
		cli.StringFlag{Name: "abi", Usage: "abi path"},
		cli.StringFlag{Name: "wasm", Usage: "wasm path"},
		cli.StringFlag{Name: "input", Usage: "contract params"},
		cli.BoolFlag{Name: "is-create", Usage: "create contract or not"},
		cli.StringFlag{Name: "func-name", Usage: "call function name"},
	}
}

func getRpcTestApp() *cli.App {
	app := cli.NewApp()
	addRpcFlags(app)
	return app
}
