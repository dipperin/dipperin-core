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
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/tests"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/dipperin/dipperin-core/third-party/log/mpt_log"
	"github.com/stretchr/testify/assert"
	"gopkg.in/check.v1"
	"math/big"
)

func (suite *chainWriterSuite) TestChainState_BuildStateProcessor(t *check.C) {
	mpt_log.InitMptLogger(mpt_log.LogConf, "state_writer")

	testAccount1 := tests.AccFactory.GetAccount(0)
	suite.txBuilder.Amount = big.NewInt(100)
	suite.txBuilder.Pk = suite.env.DefaultVerifiers()[0].Pk
	suite.txBuilder.To = testAccount1.Address()

	suite.blockBuilder.Txs = []*model.Transaction{suite.txBuilder.Build()}

	block := suite.blockBuilder.Build()
	log.Info("block txs", "len", block.TxCount())
	err := suite.chainState.SaveBlock(block)
	t.Check(err, check.IsNil)

	// no this account
	s0, err := suite.chainState.StateAtByBlockNumber(0)
	assert.NoError(t, err)
	b0, err := s0.GetBalance(testAccount1.Address())
	assert.Error(t, err)
	assert.Nil(t, b0)

	// this account exists after the transaction
	s1, err := suite.chainState.CurrentState()
	assert.NoError(t, err)
	b1, err := s1.GetBalance(testAccount1.Address())
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(100), b1)
}

func (suite *chainWriterSuite) TestChainState_GetStateStorage(t *check.C) {
	assert.NotNil(t, suite.chainState.GetStateStorage())
}

func (suite *chainWriterSuite) TestChainState_CurrentState(t *check.C) {
	s0, err := suite.chainState.CurrentState()
	assert.NoError(t, err)
	b0, err := s0.GetBalance(suite.env.DefaultVerifiers()[0].Address())
	assert.NoError(t, err)
	assert.NotNil(t, b0)
	assert.Equal(t, 1, b0.Cmp(big.NewInt(0)))
}

func (suite *chainWriterSuite) TestChainState_StateAtByBlockNumber(t *check.C) {
	s0, err := suite.chainState.StateAtByBlockNumber(0)
	assert.NoError(t, err)
	b0, err := s0.GetBalance(suite.env.DefaultVerifiers()[0].Address())
	assert.NoError(t, err)
	assert.NotNil(t, b0)
	assert.Equal(t, 1, b0.Cmp(big.NewInt(0)))
}

func (suite *chainWriterSuite) TestChainState_StateAtByStateRoot(t *check.C) {
	s0, err := suite.chainState.StateAtByStateRoot(suite.chainState.CurrentBlock().StateRoot())
	assert.NoError(t, err)
	b0, err := s0.GetBalance(suite.env.DefaultVerifiers()[0].Address())
	assert.NoError(t, err)
	assert.NotNil(t, b0)
	assert.Equal(t, 1, b0.Cmp(big.NewInt(0)))
}
