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


package chain_state

import (
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/dipperin/dipperin-core/core/cs-chain/chain-writer"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/tests"
	"github.com/dipperin/dipperin-core/third-party/crypto"
	"github.com/stretchr/testify/assert"
	"gopkg.in/check.v1"
	"math/big"
)

type BaseChainSuite struct {
	chainState   *ChainState
	env          *tests.GenesisEnv
	txBuilder    *tests.TxBuilder
	blockBuilder *tests.BlockBuilder
	bootNode 	 []tests.Account
}

func (suite *BaseChainSuite) SetUpTest(c *check.C) {
	account, _ := tests.ChangeVerBootNodeAddress()
	suite.bootNode = account
	f := chain_writer.NewChainWriterFactory()

	conf := chain_config.GetChainConfig()
	conf.SlotSize = 3
	conf.VerifierNumber = 3
	suite.chainState = NewChainState(&ChainStateConfig{
		DataDir:       "",
		WriterFactory: f,
		ChainConfig:   conf,
	})
	f.SetChain(suite.chainState)

	suite.env = tests.NewGenesisEnv(suite.chainState.GetChainDB(), suite.chainState.GetStateStorage(), nil)

	suite.txBuilder = &tests.TxBuilder{
		Nonce:  0,
		To:     common.HexToAddress(fmt.Sprintf("0x123a%v", 1)),
		Amount: big.NewInt(1),
		Pk:     suite.env.DefaultVerifiers()[0].Pk,
		Fee:	testFee,
	}
	suite.blockBuilder = &tests.BlockBuilder{
		ChainState: suite.chainState,
		PreBlock:   suite.chainState.CurrentBlock(),
		MinerPk:    suite.env.Miner().Pk,
	}
}

func (suite *BaseChainSuite) TearDownTest(c *check.C) {
	//c.Check(os.RemoveAll("/tmp/base_chain_suite1_test"), check.IsNil)
}

func (suite *BaseChainSuite) InsertBlock(t *check.C, num int) {
	curBlock := suite.chainState.CurrentBlock()
	curNum := int(curBlock.Number())
	config := chain_config.GetChainConfig()
	suite.blockBuilder.SetPreBlock(curBlock)
	for i := curNum; i < curNum+num; i++ {

		// curBlock on chain
		curBlock = suite.chainState.CurrentBlock()
		var block model.AbstractBlock
		if curBlock.Number() == 0 {
			block = suite.blockBuilder.Build()
		} else {

			// votes for curBlock on chain
			var curBlockVotes []model.AbstractVerification
			if curBlock.IsSpecial() {
				vote := suite.CreateVerBootVote(curBlock)
				curBlockVotes = append(curBlockVotes, vote)
			} else {
				curBlockVotes = suite.env.VoteBlock(config.VerifierNumber, 0, curBlock)
			}
			suite.blockBuilder.SetVerifivations(curBlockVotes)
			block = suite.blockBuilder.Build()
		}

		// votes for build block
		votes := suite.env.VoteBlock(config.VerifierNumber, 0, block)
		err := suite.chainState.SaveBftBlock(block, votes)
		suite.blockBuilder.SetPreBlock(block)
		assert.NoError(t, err)
		assert.Equal(t, uint64(i+1), suite.chainState.CurrentBlock().Number())
		assert.Equal(t, false, suite.chainState.CurrentBlock().IsSpecial())
	}
}

func (suite *BaseChainSuite) InsertSpecialBlock(t *check.C, num int) {
	curBlock := suite.chainState.CurrentBlock()
	curNum := int(curBlock.Number())
	config := chain_config.GetChainConfig()
	suite.blockBuilder.SetPreBlock(curBlock)
	suite.blockBuilder.SetMinerPk(suite.bootNode[0].Pk)
	for i := curNum; i < curNum+num; i++ {

		// curBlock on chain
		curBlock = suite.chainState.CurrentBlock()
		var block model.AbstractBlock
		if curBlock.Number() == 0 {
			block = suite.blockBuilder.BuildSpecialBlock()
		} else {

			// votes for curBlock on chain
			var curBlockVotes []model.AbstractVerification
			if curBlock.IsSpecial() {
				vote := suite.CreateVerBootVote(curBlock)
				curBlockVotes = append(curBlockVotes, vote)
			} else {
				curBlockVotes = suite.env.VoteBlock(config.VerifierNumber, 0, curBlock)
			}
			suite.blockBuilder.SetVerifivations(curBlockVotes)
			block = suite.blockBuilder.BuildSpecialBlock()
		}

		// votes for build block
		vote := suite.CreateVerBootVote(block)
		err := suite.chainState.SaveBftBlock(block, []model.AbstractVerification{vote})
		suite.blockBuilder.SetPreBlock(block)
		assert.NoError(t, err)
		assert.Equal(t, uint64(i+1), suite.chainState.CurrentBlock().Number())
		assert.Equal(t, true, suite.chainState.CurrentBlock().IsSpecial())
	}
}

func (suite *BaseChainSuite) CreateVerBootVote(block model.AbstractBlock) model.AbstractVerification {
	voteA := model.NewVoteMsg(block.Number(), uint64(0), block.Hash(), model.VerBootNodeVoteMessage)
	sign, _ := crypto.Sign(voteA.Hash().Bytes(), suite.bootNode[0].Pk)
	voteA.Witness.Address = suite.bootNode[0].Address()
	voteA.Witness.Sign = sign
	return voteA
}