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

package spec_util

import (
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/dipperin/dipperin-core/core/chain/cachedb"
	"github.com/dipperin/dipperin-core/core/cs-chain"
	"github.com/dipperin/dipperin-core/core/cs-chain/chain-state"
	"github.com/dipperin/dipperin-core/core/cs-chain/chain-writer"
	"github.com/dipperin/dipperin-core/core/tx-pool"
	"github.com/dipperin/dipperin-core/tests"
	"github.com/hashicorp/golang-lru"
	"math/big"
	"time"
)

type BaseChainServiceSuite struct {
	TxPool         *tx_pool.TxPool
	CacheDB        *cachedb.CacheDB
	Env            *tests.GenesisEnv
	TxBuilder      *tests.TxBuilder
	BlockBuilder   *tests.BlockBuilder
	CsChainService *cs_chain.CsChainService
}

func (suite *BaseChainServiceSuite) SetupTest() {

	config := &cs_chain.CsChainServiceConfig{
		ChainStateConfig: &chain_state.ChainStateConfig{
			DataDir:     "",
			ChainConfig: chain_config.GetChainConfig(),
		}}

	config.ChainStateConfig.WriterFactory = chain_writer.NewChainWriterFactory()
	chainState := chain_state.NewChainState(config.ChainStateConfig)

	futureBlocks, _ := lru.New(100)

	//cache chain state
	ccs, err := cs_chain.NewCacheChainState(chainState)
	if err != nil {
		panic(err)
	}

	config.ChainStateConfig.WriterFactory.SetChain(ccs)
	service := &cs_chain.CsChainService{
		CsChainServiceConfig: config,
		CacheChainState:      ccs,
		FutureBlocks:         futureBlocks,
		Quit:                 make(chan struct{}),
	}

	suite.Env = tests.NewGenesisEnv(chainState.GetChainDB(), chainState.GetStateStorage(), nil)

	if err := service.InitService(); err != nil {
		panic(err)
	}

	suite.CsChainService = service

	// tx pool
	suite.TxPool = tx_pool.NewTxPool(tx_pool.DefaultTxPoolConfig, *suite.CsChainService.GetChainConfig(),
		suite.CsChainService)

	// start tx pool
	go func() {
		if err := suite.TxPool.Start(); err != nil {
			panic(err)
		}
	}()

	// cache db
	suite.CacheDB = cachedb.NewCacheDB(suite.CsChainService.GetDB())
	cachedb.SetCacheDataDecoder(&cachedb.BFTCacheDataDecoder{})

	suite.CsChainService.TxPool = suite.TxPool
	suite.CsChainService.CacheDB = suite.CacheDB

	suite.TxBuilder = &tests.TxBuilder{
		Nonce:  1,
		To:     common.HexToAddress(fmt.Sprintf("0x123a%v", 1)),
		Amount: big.NewInt(1),
		Pk:     suite.Env.DefaultVerifiers()[0].Pk,
		Fee:    testFee,
	}
	suite.BlockBuilder = &tests.BlockBuilder{
		ChainState: suite.CsChainService.ChainState,
		PreBlock:   suite.CsChainService.ChainState.CurrentBlock(),
		MinerPk:    suite.Env.Miner().Pk,
	}
}

func (suite *BaseChainServiceSuite) TearDownTest() {
	suite.TxPool.Stop()
	time.Sleep(100 * time.Microsecond)
	suite.TxPool = nil
	suite.CacheDB = nil
	suite.Env = nil
	suite.TxBuilder = nil
	suite.BlockBuilder = nil
	suite.CsChainService = nil
}
