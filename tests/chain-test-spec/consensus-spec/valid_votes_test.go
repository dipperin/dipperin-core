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
	"github.com/stretchr/testify/suite"
	"github.com/dipperin/dipperin-core/tests/chain-test-spec/spec-util"
	"testing"
	"github.com/dipperin/dipperin-core/core/cs-chain/chain-writer/middleware"
	"github.com/dipperin/dipperin-core/common/g-error"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/tests"
	"github.com/dipperin/dipperin-core/third-party/crypto"
)

var (
	alicePriv            = "289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232031"
	aliceAddr            = common.HexToAddress("0x00005586B883Ec6dd4f8c26063E18eb4Bd228e59c3E9")
)


type votesTestSuite struct {
	suite.Suite
	spec_util.BaseChainSuite
}

func (suite *votesTestSuite) TearDownTest() {
	suite.BaseChainSuite.TearDownTest()
}

func (suite *votesTestSuite) SetupTest() {
	suite.BaseChainSuite.SetupTest()

	if suite.BaseChainSuite.ChainState == nil {
		suite.Fail("chain state is nil")
	}
}

func TestVotesTestSuite(t *testing.T) {
	suite.Run(t, new(votesTestSuite))
}

func (suite *votesTestSuite) Test_ValidateVotesForBFT_ErrFirstBlockShouldNotHaveVerifications() {
	config := suite.ChainState.GetChainConfig()
	b0 := suite.ChainState.CurrentBlock()
	votes := suite.Env.VoteBlock(config.VerifierNumber, 1, b0)

	suite.BlockBuilder.SetVerifivations(votes)
	block := suite.BlockBuilder.Build()

	// create middleware context
	context := middleware.NewBftBlockContext(block, nil, suite.ChainState)
	context.Use(middleware.ValidateVotesForBFT(&context.BlockContext))
	suite.EqualError(context.Process(), g_error.ErrFirstBlockShouldNotHaveVerifications.Error())
}

func (suite *votesTestSuite) Test_ValidateVotesForBFT_ErrSameVoteSingerInVotes() {
	suite.insertFirstBlock()

	// create same vote
	config := suite.ChainState.GetChainConfig()
	b0 := suite.ChainState.CurrentBlock()
	votes := suite.Env.VoteBlock(1, 1, b0)
	for i := 0; i < config.VerifierNumber; i++ {
		votes = append(votes, votes[0])
	}
	suite.BlockBuilder.SetVerifivations(votes)
	block := suite.BlockBuilder.Build()

	// create middleware context
	context := middleware.NewBftBlockContext(block, nil, suite.ChainState)
	context.Use(middleware.ValidateVotesForBFT(&context.BlockContext))
	suite.EqualError(context.Process(), g_error.ErrSameVoteSingerInVotes.Error())
}

func (suite *votesTestSuite)Test_ValidateVotesForBFT_ErrBlockVotesNotEnough() {
	suite.insertFirstBlock()

	// create only one vote
	votes := suite.Env.VoteBlock(1, 1, suite.ChainState.CurrentBlock())
	suite.BlockBuilder.SetVerifivations(votes)
	block := suite.BlockBuilder.Build()

	// create middleware context
	context := middleware.NewBftBlockContext(block, nil, suite.ChainState)
	context.Use(middleware.ValidateVotesForBFT(&context.BlockContext))
	suite.EqualError(context.Process(), g_error.ErrBlockVotesNotEnough.Error())
}

func (suite *votesTestSuite)Test_ValidateVotesForBFT_ErrNotCurrentVerifier() {
	suite.insertFirstBlock()
	config := suite.ChainState.GetChainConfig()
	b0 := suite.ChainState.CurrentBlock()
	votes := suite.Env.VoteBlock(config.VerifierNumber, 1, b0)

	// create alice vote
	voteA := createVote(b0)
	votes = append(votes, voteA)

	suite.BlockBuilder.SetVerifivations(votes)
	block := suite.BlockBuilder.Build()

	// create middleware context
	context := middleware.NewBftBlockContext(block, nil, suite.ChainState)
	context.Use(middleware.ValidateVotesForBFT(&context.BlockContext))
	suite.EqualError(context.Process(), g_error.ErrNotCurrentVerifier.Error())
}

func (suite *votesTestSuite)Test_ValidateVotes_ErrInvalidFirstVoteInSpecialBlock() {
	bootNode, _ := tests.ChangeVerBootNodeAddress()
	suite.insertFirstBlock()
	b0 := suite.ChainState.CurrentBlock()

	b0Votes := suite.Env.VoteBlock(suite.ChainState.GetChainConfig().VerifierNumber, 1, b0)
	suite.BlockBuilder.SetVerifivations(b0Votes)
	suite.BlockBuilder.SetMinerPk(bootNode[0].Pk)
	specialBlock := suite.BlockBuilder.BuildSpecialBlock()

	// save special block
	voteA := createVote(b0)
	err := suite.ChainState.SaveBftBlock(specialBlock, []model.AbstractVerification{voteA})
	suite.EqualError(err, g_error.ErrInvalidFirstVoteInSpecialBlock.Error())
}

func (suite *votesTestSuite) Test_ValidateVotesForBFT_ErrInvalidBlockHashInVotes() {
	config := suite.ChainState.GetChainConfig()
	b0 := suite.ChainState.CurrentBlock()
	suite.insertFirstBlock()

	votes := suite.Env.VoteBlock(config.VerifierNumber, 1, b0)
	suite.BlockBuilder.SetVerifivations(votes)
	block := suite.BlockBuilder.Build()

	// create middleware context
	context := middleware.NewBftBlockContext(block, nil, suite.ChainState)
	context.Use(middleware.ValidateVotesForBFT(&context.BlockContext))
	suite.EqualError(context.Process(), g_error.ErrInvalidBlockHashInVotes.Error())
}

func (suite *votesTestSuite) insertFirstBlock() {
	block := suite.BlockBuilder.Build()
	seenCommit := suite.Env.VoteBlock(suite.ChainState.GetChainConfig().VerifierNumber, 1, block)
	err := suite.ChainState.SaveBftBlock(block, seenCommit)
	suite.NoError(err)

	suite.BlockBuilder.SetPreBlock(block)
}

func createVote(block model.AbstractBlock) model.AbstractVerification {
	voteA := model.NewVoteMsg(block.Number(), uint64(0), block.Hash(), model.VoteMessage)
	key, _ := crypto.HexToECDSA(alicePriv)
	sign, _ := crypto.Sign(voteA.Hash().Bytes(), key)
	voteA.Witness.Address = aliceAddr
	voteA.Witness.Sign = sign
	return voteA
}