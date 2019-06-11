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
	"github.com/dipperin/dipperin-core/core/cs-chain/chain-state"
	"github.com/dipperin/dipperin-core/core/cs-chain/chain-writer"
	"github.com/dipperin/dipperin-core/core/economy-model"
	"github.com/dipperin/dipperin-core/tests"
	"github.com/dipperin/dipperin-core/third-party/log"
	"math/big"
	"github.com/dipperin/dipperin-core/core/model"
)

var testFee = economy_model.GetMinimumTxFee(50)

type BaseChainSuite struct {
	ChainState   *chain_state.ChainState
	Env          *tests.GenesisEnv
	TxBuilder    *tests.TxBuilder
	BlockBuilder *tests.BlockBuilder
}

func (suite *BaseChainSuite) SetupTest() {
	model.IgnoreDifficultyValidation = true
	f := chain_writer.NewChainWriterFactory()

	suite.ChainState = chain_state.NewChainState(&chain_state.ChainStateConfig{
		DataDir:       "",
		WriterFactory: f,
		ChainConfig:   chain_config.GetChainConfig(),
	})
	f.SetChain(suite.ChainState)

	suite.Env = tests.NewGenesisEnv(suite.ChainState.GetChainDB(), suite.ChainState.GetStateStorage(), nil)
	suite.TxBuilder = &tests.TxBuilder{
		Nonce:  1,
		To:     common.HexToAddress(fmt.Sprintf("0x123a%v", 1)),
		Amount: big.NewInt(1),
		Fee:    testFee,
		Pk:     suite.Env.DefaultVerifiers()[0].Pk,
	}
	suite.BlockBuilder = &tests.BlockBuilder{
		ChainState: suite.ChainState,
		PreBlock:   suite.ChainState.CurrentBlock(),
		MinerPk:    suite.Env.Miner().Pk,
	}

	log.Info("the suit preBlock is~~~~:","preBlock",suite.ChainState.CurrentBlock())
}

func (suite *BaseChainSuite) TearDownTest() {
	//c.Check(os.RemoveAll("/tmp/base_chain_suite1_test"), check.IsNil)
}
