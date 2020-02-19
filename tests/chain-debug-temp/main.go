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
	"fmt"
	"github.com/dipperin/dipperin-core/core/chain"
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/dipperin/dipperin-core/core/chain/state-processor"
	"github.com/dipperin/dipperin-core/core/contract"
	"github.com/dipperin/dipperin-core/core/cs-chain"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/urfave/cli"
	"os"
	"reflect"
)

const (
	DataDirFName = "data_dir"
)

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: DataDirFName, Value: "/home/qydev/csdebug/full_chain_data"},
	}
	app.Action = run
	fmt.Println(os.Args)
	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}

/**/
func run(c *cli.Context) {
	nContext := initContext(c)
	curBlock := nContext.fullChain.CurrentBlock()
	log.Info("cur block num", "num", curBlock.Number())

	state, err := chain.NewBlockProcessor(nContext.fullChain, curBlock.StateRoot(), nContext.stateStorage)
	if err != nil {
		panic(err)
	}

	cb, err := state.GetContract(contract.EarlyContractAddress, reflect.TypeOf(contract.EarlyRewardContract{}))
	if err != nil {
		panic(err)
	}
	fmt.Println(cb.Interface())
}

func initContext(c *cli.Context) *nodeContext {
	dataDir := c.String(DataDirFName)
	db, err := ethdb.NewLDBDatabase(dataDir, 0, 0)
	if err != nil {
		panic(err)
	}

	model.SetBlockRlpHandler(&model.PBFTBlockRlpHandler{})
	model.SetBlockJsonHandler(&model.PBFTBlockJsonHandler{})

	context := &nodeContext{
		db: db, blockDecoder: model.MakeDefaultBlockDecoder(),
		chainConfig: *chain_config.GetChainConfig(),
	}

	context.stateStorage = state_processor.NewStateStorageWithCache(context.db)

	//context.fullChain = chain.MakeFullChain(context)

	return context
}

type nodeContext struct {
	chainConfig  chain_config.ChainConfig
	blockDecoder model.BlockDecoder
	db           ethdb.Database
	stateStorage state_processor.StateStorage

	fullChain cs_chain.Chain
}
