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

package consensus_spec

import (
	"github.com/dipperin/dipperin-core/common/g-error"
	"github.com/dipperin/dipperin-core/core/cs-chain/chain-writer/middleware"
	"github.com/dipperin/dipperin-core/tests/chain-test-spec/spec-util"
	"github.com/stretchr/testify/suite"
	"testing"
)

type consensusTestSuite struct {
	suite.Suite
	spec_util.BaseChainSuite
}

func (suite *consensusTestSuite) TearDownTest() {
	suite.BaseChainSuite.TearDownTest()
}

func (suite *consensusTestSuite) SetupTest() {
	suite.BaseChainSuite.SetupTest()

	if suite.BaseChainSuite.ChainState == nil {
		suite.Fail("chain state is nil")
	}
}

func TestConsensusTestSuite(t *testing.T) {
	suite.Run(t, new(consensusTestSuite))
}

func fakeChainSuite() *spec_util.BaseChainSuite {
	chainTool := &spec_util.BaseChainSuite{}
	chainTool.SetupTest()
	return chainTool
}

func (suite *consensusTestSuite) Test_ValidateBlockNumber() {
	//chainHelper used for making future block
	chainHelper := fakeChainSuite()
	block := suite.BlockBuilder.Build()

	// seen commit
	seenCommit := suite.Env.VoteBlock(suite.ChainState.GetChainConfig().VerifierNumber, 1, block)

	// create middleware context
	context := middleware.NewBftBlockContext(block, seenCommit, suite.ChainState)
	context.Use(middleware.ValidateBlockNumber(&context.BlockContext))
	suite.NoError(context.Process())
	suite.NoError(suite.ChainState.SaveBftBlock(block, seenCommit))
	suite.NoError(chainHelper.ChainState.SaveBftBlock(block, seenCommit))
	suite.Equal(block.Number(), suite.ChainState.CurrentBlock().Number())

	// create middleware context
	context2 := middleware.NewBftBlockContext(block, seenCommit, suite.ChainState)
	context2.Use(middleware.ValidateBlockNumber(&context2.BlockContext))
	suite.EqualError(context2.Process(), g_error.ErrBlockHeightIsCurrentAndIsNotSpecial.Error())

	//test special block
	blockSpec := suite.BlockBuilder.BuildSpecialBlock()
	context3 := middleware.NewBftBlockContext(blockSpec, seenCommit, suite.ChainState)
	context3.Use(middleware.ValidateBlockNumber(&context3.BlockContext))
	suite.NoError(context3.Process())

	suite.BlockBuilder.SetPreBlock(block)
	suite.BlockBuilder.SetVerifivations(seenCommit)
	block1 := suite.BlockBuilder.Build()
	context4 := middleware.NewBftBlockContext(block1, seenCommit, suite.ChainState)
	context4.Use(middleware.ValidateBlockNumber(&context4.BlockContext))
	suite.NoError(context4.Process())

	//test future block
	seenCommit1 := suite.Env.VoteBlock(suite.ChainState.GetChainConfig().VerifierNumber, 1, block1)
	suite.NoError(chainHelper.ChainState.SaveBftBlock(block1, seenCommit1))
	chainHelper.BlockBuilder.SetPreBlock(block1)
	chainHelper.BlockBuilder.SetVerifivations(seenCommit1)
	block2 := chainHelper.BlockBuilder.Build()
	context5 := middleware.NewBftBlockContext(block2, seenCommit1, suite.ChainState)
	context5.Use(middleware.ValidateBlockNumber(&context5.BlockContext))
	suite.EqualError(context5.Process(), g_error.ErrFutureBlock.Error())

	//test far-away future block
	block3 := chainHelper.BlockBuilder.BuildFuture()
	context6 := middleware.NewBftBlockContext(block3, seenCommit1, suite.ChainState)
	context6.Use(middleware.ValidateBlockNumber(&context6.BlockContext))
	suite.EqualError(context6.Process(), g_error.ErrFutureBlockTooFarAway.Error())
}

func (suite *consensusTestSuite) Test_ValidateBlockHash() {
	//chainHelper used for making future block
	chainHelper := fakeChainSuite()
	//fork chain
	chainHelperA := fakeChainSuite()
	block := suite.BlockBuilder.Build()

	// seen commit
	seenCommit := suite.Env.VoteBlock(suite.ChainState.GetChainConfig().VerifierNumber, 1, block)
	suite.NoError(suite.ChainState.SaveBftBlock(block, seenCommit))
	context := middleware.NewBftBlockContext(block, seenCommit, suite.ChainState)
	context.Use(middleware.ValidateBlockHash(&context.BlockContext))
	suite.NoError(context.Process())
	suite.NoError(chainHelper.ChainState.SaveBftBlock(block, seenCommit))

	blockA := suite.BlockBuilder.Build()
	seenCommitA := suite.Env.VoteBlock(suite.ChainState.GetChainConfig().VerifierNumber, 1, blockA)
	suite.NoError(chainHelperA.ChainState.SaveBftBlock(blockA, seenCommitA))

	suite.BlockBuilder.SetPreBlock(block)
	suite.BlockBuilder.SetVerifivations(seenCommit)
	block1 := suite.BlockBuilder.Build()
	seenCommit1 := suite.Env.VoteBlock(suite.ChainState.GetChainConfig().VerifierNumber, 1, block1)
	context1 := middleware.NewBftBlockContext(block1, seenCommit1, suite.ChainState)
	context1.Use(middleware.ValidateBlockHash(&context1.BlockContext))
	suite.NoError(context1.Process())

	//test not match prehash
	chainHelperA.BlockBuilder.SetPreBlock(blockA)
	chainHelperA.BlockBuilder.SetVerifivations(seenCommitA)
	block1A := chainHelperA.BlockBuilder.Build()
	seenCommit1A := chainHelperA.Env.VoteBlock(suite.ChainState.GetChainConfig().VerifierNumber, 1, block1A)
	context2 := middleware.NewBftBlockContext(block1A, seenCommit1A, suite.ChainState)
	context2.Use(middleware.ValidateBlockHash(&context2.BlockContext))
	suite.EqualError(context2.Process(), g_error.ErrPreBlockHashNotMatch.Error())

	//test preblock not existed
	suite.NoError(chainHelper.ChainState.SaveBftBlock(block1, seenCommit1))
	chainHelper.BlockBuilder.SetPreBlock(block1)
	chainHelper.BlockBuilder.SetVerifivations(seenCommit1)
	block2 := chainHelper.BlockBuilder.Build()
	context3 := middleware.NewBftBlockContext(block2, seenCommit1, suite.ChainState)
	context3.Use(middleware.ValidateBlockHash(&context3.BlockContext))
	suite.EqualError(context3.Process(), g_error.ErrPreBlockIsNil.Error())

}

func (suite *consensusTestSuite) Test_ValidateBlockCoinBase() {
	block := suite.BlockBuilder.Build()

	// seen commit
	seenCommit := suite.Env.VoteBlock(suite.ChainState.GetChainConfig().VerifierNumber, 1, block)
	context := middleware.NewBftBlockContext(block, seenCommit, suite.ChainState)
	context.Use(middleware.ValidateBlockCoinBase(&context.BlockContext))
	suite.NoError(context.Process())

	blockSpec := suite.BlockBuilder.BuildSpecialBlock()
	seenCommitSpec := suite.Env.VoteBlock(suite.ChainState.GetChainConfig().VerifierNumber, 1, blockSpec)
	context1 := middleware.NewBftBlockContext(blockSpec, seenCommitSpec, suite.ChainState)
	context1.Use(middleware.ValidateBlockCoinBase(&context1.BlockContext))
	suite.EqualError(context1.Process(), g_error.ErrSpecialInvalidCoinBase.Error())
}

func (suite *consensusTestSuite) Test_ValidateBlockDifficulty() {
	block := suite.BlockBuilder.Build()
	// seen commit
	seenCommit := suite.Env.VoteBlock(suite.ChainState.GetChainConfig().VerifierNumber, 1, block)
	suite.NoError(suite.ChainState.SaveBftBlock(block, seenCommit))
	context := middleware.NewBftBlockContext(block, seenCommit, suite.ChainState)
	context.Use(middleware.ValidateBlockDifficulty(&context.BlockContext))
	suite.NoError(context.Process())
}

func (suite *consensusTestSuite) Test_ValidateSeed() {
	block := suite.BlockBuilder.Build()
	// seen commit
	seenCommit := suite.Env.VoteBlock(suite.ChainState.GetChainConfig().VerifierNumber, 1, block)
	context := middleware.NewBftBlockContext(block, seenCommit, suite.ChainState)
	context.Use(middleware.ValidateSeed(&context.BlockContext))
	suite.NoError(context.Process())
}

func (suite *consensusTestSuite) Test_ValidateBlockVersion() {
	block := suite.BlockBuilder.Build()
	// seen commit
	seenCommit := suite.Env.VoteBlock(suite.ChainState.GetChainConfig().VerifierNumber, 1, block)
	context := middleware.NewBftBlockContext(block, seenCommit, suite.ChainState)
	context.Use(middleware.ValidateBlockVersion(&context.BlockContext))
	suite.NoError(context.Process())
}
