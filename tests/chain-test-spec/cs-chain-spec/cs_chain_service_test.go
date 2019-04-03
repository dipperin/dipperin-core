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


package cs_chain_spec

import (
	"github.com/dipperin/dipperin-core/core/cs-chain/chain-writer/middleware"
	"github.com/dipperin/dipperin-core/tests/chain-test-spec/spec-util"
	"github.com/stretchr/testify/suite"
	"testing"
	"github.com/dipperin/dipperin-core/core/model"
)

type chainServiceTestSuite struct {
	suite.Suite
	spec_util.BaseChainServiceSuite
}

func TestChainServiceSuite(t *testing.T) {
	suite.Run(t, new(chainServiceTestSuite))
}

func (suite *chainServiceTestSuite) TearDownTest() {
	suite.BaseChainServiceSuite.TearDownTest()
}

func (suite *chainServiceTestSuite) SetupTest() {
	suite.BaseChainServiceSuite.SetupTest()
}

func (suite *chainServiceTestSuite) TestCsChainService_handleFutureBlock() {
	block := suite.BlockBuilder.Build()

	// seen commit
	seenCommit := suite.Env.VoteBlock(suite.CsChainService.GetChainConfig().VerifierNumber, 1, block)

	// create middleware context
	context := middleware.NewBftBlockContext(block, seenCommit, suite.CsChainService.ChainState)

	context.Use(middleware.ValidateBlockNumber(&context.BlockContext))

	suite.NoError(context.Process())

	model.IgnoreDifficultyValidation = true

	suite.NoError(suite.CsChainService.SaveBftBlock(block, seenCommit))

	suite.Equal(block.Number(), suite.CsChainService.CurrentBlock().Number())

}