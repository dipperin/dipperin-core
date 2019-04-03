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
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/chain"
	"github.com/dipperin/dipperin-core/core/chain/registerdb"
	"github.com/dipperin/dipperin-core/core/chain/state-processor"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third-party/log"
	"gopkg.in/check.v1"
	"testing"
)

// Hook up gocheck into the "go test" runner
func TestChainState(t *testing.T) { check.TestingT(t) }

type ChainStateSuite struct {
	BaseChainSuite
	//chainState *ChainState
	block      model.AbstractBlock
}

var _ = check.Suite(&ChainStateSuite{})

func (suite *ChainStateSuite) SetUpTest(c *check.C) {
	log.InitLogger(log.LvlError)
	suite.BaseChainSuite.SetUpTest(c)

	// insert block
	// mock block
	block := suite.blockBuilder.Build()
	config := suite.chainState.ChainConfig
	votes := suite.env.VoteBlock(config.VerifierNumber, 0, block)
	err := suite.chainState.SaveBftBlock(block, votes)
	c.Check(err, check.IsNil)
	suite.block = block
	suite.blockBuilder.SetPreBlock(block)
}

func (suite *ChainStateSuite) createGenesisBlock() {
	if suite.chainState.StateStorage == nil || suite.chainState.ChainDB == nil {
		panic("you need new chain state first")
	}

	genesisAccountStateProcessor, err := state_processor.MakeGenesisAccountStateProcessor(suite.chainState.StateStorage)
	if err != nil {
		panic("open account state processor for genesis failed: " + err.Error())
	}

	genesisRegisterProcessor, err := registerdb.MakeGenesisRegisterProcessor(suite.chainState.StateStorage)
	if err != nil {
		panic("make registerDB processor for genesis failed: " + err.Error())
	}

	// setup genesis block
	defaultGenesis := chain.DefaultGenesisBlock(suite.chainState.ChainDB, genesisAccountStateProcessor, genesisRegisterProcessor,
		suite.chainState.ChainConfig)

	if _, _, err = chain.SetupGenesisBlock(defaultGenesis); err != nil {
		panic("setup genesis block failed: " + err.Error())
	}
}

func (suite *ChainStateSuite) TearDownTest(c *check.C) {
	suite.chainState = nil
}

func (suite *ChainStateSuite) TestChainState_Genesis(c *check.C) {
	genesis := suite.chainState.Genesis()

	c.Check(genesis, check.NotNil)

	c.Check(genesis.Number(),  check.Equals, uint64(0))

	c.Check(genesis.Hash().IsEmpty(), check.Equals, false)
}

func (suite *ChainStateSuite) TestChainState_CurrentBlock(c *check.C) {
	curBlock := suite.chainState.CurrentBlock()

	c.Check(curBlock.Number(),  check.Equals, suite.block.Number())

	c.Check(curBlock.Hash().IsEmpty(), check.Equals, false)
}

func (suite *ChainStateSuite) TestChainState_CurrentHeader(c *check.C) {
	curHeader := suite.chainState.CurrentHeader()

	c.Check(curHeader.GetNumber(),  check.Equals, suite.block.Number())

	c.Check(curHeader.Hash().IsEmpty(), check.Equals, false)
}

func (suite *ChainStateSuite) TestChainState_GetBlock(c *check.C) {
	block := suite.chainState.GetBlock(suite.block.Hash(), suite.block.Number())

	c.Check(block, check.NotNil)

	c.Check(block.Number(),  check.Equals, suite.block.Number())

	c.Check(block.Hash().IsEqual(suite.block.Hash()), check.Equals, true)

	suite.chainState.GetHeaderRLP(common.Hash{})
	suite.chainState.GetHeaderRLP(suite.block.Hash())
	suite.chainState.GetTransaction(common.Hash{})
}

func (suite *ChainStateSuite) TestChainState_GetBlockByHash(c *check.C) {
	block := suite.chainState.GetBlockByHash(suite.block.Hash())

	c.Check(block, check.NotNil)

	c.Check(block.Number(),  check.Equals, suite.block.Number())

	c.Check(block.Hash().IsEqual(suite.block.Hash()), check.Equals, true)
}

func (suite *ChainStateSuite) TestChainState_GetBlockByNumber(c *check.C) {
	block := suite.chainState.GetBlockByNumber(suite.block.Number())

	c.Check(block, check.NotNil)

	c.Check(block.Number(),  check.Equals, suite.block.Number())

	c.Check(block.Hash().IsEqual(suite.block.Hash()), check.Equals, true)
}

func (suite *ChainStateSuite) TestChainState_HasBlock(c *check.C) {
	result := suite.chainState.HasBlock(suite.block.Hash(), suite.block.Number())

	c.Check(result, check.Equals, true)
}

func (suite *ChainStateSuite) TestChainState_GetBody(c *check.C) {
	body := suite.chainState.GetBody(suite.block.Hash())

	c.Check(body, check.NotNil)
	c.Check(body.GetTxsSize(), check.Equals, suite.block.TxCount())
}

func (suite *ChainStateSuite) TestChainState_GetBodyRLP(c *check.C) {
	body := suite.chainState.GetBodyRLP(suite.block.Hash())

	c.Check(len(body) > 0, check.Equals, true)
}

func (suite *ChainStateSuite) TestChainState_GetHeader(c *check.C) {
	header := suite.chainState.GetHeader(suite.block.Hash(), suite.block.Number())

	c.Check(header, check.NotNil)
}

func (suite *ChainStateSuite) TestChainState_GetHeaderByHash(c *check.C) {
	header := suite.chainState.GetHeaderByHash(suite.block.Hash())

	c.Check(header, check.NotNil)
}

func (suite *ChainStateSuite) TestChainState_GetHeaderByNumber(c *check.C) {
	header := suite.chainState.GetHeaderByNumber(suite.block.Number())

	c.Check(header, check.NotNil)
}

func (suite *ChainStateSuite) TestChainState_HasHeader(c *check.C) {
	result := suite.chainState.HasHeader(suite.block.Hash(), suite.block.Number())

	c.Check(result, check.Equals, true)
}

func (suite *ChainStateSuite) TestChainState_GetBlockNumber(c *check.C) {
	result := suite.chainState.GetBlockNumber(suite.block.Hash())

	c.Check(*result, check.Equals, suite.block.Number())
}

func (suite *ChainStateSuite) TestChainState_GetLatestNormalBlock(c *check.C) {
	// insert special blocks
	config := suite.chainState.ChainConfig
	suite.InsertBlock(c, int(config.SlotSize))
	curBlock := suite.chainState.CurrentBlock()
	suite.InsertSpecialBlock(c, int(config.SlotSize))
	block := suite.chainState.GetLatestNormalBlock()
	c.Check(block.Number(), check.Equals, curBlock.Number())
}