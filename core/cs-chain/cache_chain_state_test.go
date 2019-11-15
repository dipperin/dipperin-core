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

package cs_chain

import (
	"fmt"
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/dipperin/dipperin-core/core/cs-chain/chain-state"
	"github.com/dipperin/dipperin-core/core/cs-chain/chain-writer"
	"github.com/dipperin/dipperin-core/core/economy-model"
	"github.com/dipperin/dipperin-core/tests"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"

	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/model"
)

var testFee = economy_model.GetMinimumTxFee(50)

func getTestCacheEnv() (*CacheChainState, *tests.GenesisEnv, *tests.TxBuilder, *tests.BlockBuilder) {
	model.IgnoreDifficultyValidation = true
	f := chain_writer.NewChainWriterFactory()

	cConf := chain_config.GetChainConfig()
	cConf.SlotSize = 3
	cConf.VerifierNumber = 3

	chainState := chain_state.NewChainState(&chain_state.ChainStateConfig{
		DataDir:       "",
		WriterFactory: f,
		ChainConfig:   cConf,
	})

	attackEnv := tests.NewGenesisEnv(chainState.GetChainDB(), chainState.GetStateStorage(), nil)

	txB := &tests.TxBuilder{
		Nonce:  1,
		To:     common.HexToAddress(fmt.Sprintf("0x123a%v", 1)),
		Amount: big.NewInt(1),
		Pk:     attackEnv.DefaultVerifiers()[0].Pk,
		Fee:    testFee,
	}
	bb := &tests.BlockBuilder{
		ChainState: chainState,
		PreBlock:   chainState.CurrentBlock(),
		MinerPk:    attackEnv.Miner().Pk,
	}

	cs, err := NewCacheChainState(chainState)
	if err != nil {
		panic(err)
	}
	f.SetChain(cs)

	return cs, attackEnv, txB, bb
}

func TestCacheChainState_CurrentBlock(t *testing.T) {
	cs, err := NewCacheChainState(chain_state.NewChainState(&chain_state.ChainStateConfig{ChainConfig: chain_config.GetChainConfig(), WriterFactory: chain_writer.NewChainWriterFactory()}))
	assert.NoError(t, err)
	assert.Nil(t, cs.CurrentBlock())
	assert.Nil(t, cs.GetBody(common.Hash{}))
	assert.Nil(t, cs.GetBodyRLP(common.Hash{}))
	assert.Nil(t, cs.GetBlock(common.Hash{}, 1))
	assert.Nil(t, cs.GetHeader(common.Hash{}, 1))
	assert.Nil(t, cs.GetBlockByHash(common.Hash{}))

	cs, _, _, bb := getTestCacheEnv()
	curB := cs.CurrentBlock()
	//votes := env.VoteBlock(len(env.DefaultVerifiers()),0,curB)
	//bb.SetVerifications(votes)
	//bb.PreBlock = curB
	block := bb.BuildFuture()

	//block := model.CreateBlock(curB.Number() + 1, curB.Hash(),0)
	//err = cs.SaveBftBlock(block, votes)
	cs.ChainDB.InsertBlock(block)
	assert.NoError(t, err)
	assert.NotNil(t, cs.GetBody(curB.Hash()))
	assert.NotNil(t, cs.GetBody(curB.Hash()))
	assert.NotNil(t, cs.GetBodyRLP(curB.Hash()))
	assert.NotNil(t, cs.GetBodyRLP(curB.Hash()))
	assert.True(t, cs.HasBlock(curB.Hash(), curB.Number()))
	assert.NotNil(t, cs.GetBlock(curB.Hash(), curB.Number()))
	assert.True(t, cs.HasBlock(curB.Hash(), curB.Number()))
	assert.NotNil(t, cs.GetHeader(curB.Hash(), curB.Number()))
	assert.True(t, cs.HasHeader(curB.Hash(), curB.Number()))
	assert.NotNil(t, cs.GetHeader(curB.Hash(), curB.Number()))

	assert.NotNil(t, cs.GetSlotByNum(curB.Number()))
	assert.NotNil(t, cs.GetSlotByNum(curB.Number()))
	assert.Nil(t, cs.GetSlotByNum(22))

	cs.slotCache.Remove(curB.Number())
	assert.NotNil(t, cs.GetSlot(curB))
	assert.NotNil(t, cs.GetSlot(curB))
	assert.Nil(t, cs.GetSlot(model.NewBlock(&model.Header{Number: 22}, nil, nil)))

	assert.Error(t, cs.Rollback(curB.Number()+3))
	assert.NoError(t, cs.Rollback(curB.Number()+1))
}
