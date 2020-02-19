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

package utils

import (
	"github.com/dipperin/dipperin-core/core/chain"
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/dipperin/dipperin-core/core/chain/registerdb"
	"github.com/dipperin/dipperin-core/core/chain/state-processor"
	"github.com/dipperin/dipperin-core/core/cs-chain/chain-state"
	"github.com/dipperin/dipperin-core/core/cs-chain/chain-writer"
)

func SetupGenesis(dataDir string, cConfig *chain_config.ChainConfig) {
	cs := chain_state.NewChainState(&chain_state.ChainStateConfig{
		ChainConfig:   cConfig,
		DataDir:       dataDir,
		WriterFactory: chain_writer.NewChainWriterFactory(),
	})
	setupGenesis(cs)
	cs.ChainDB.DB().Close()
}

func setupGenesis(cs *chain_state.ChainState) {
	genesisAccountStateProcessor, err := state_processor.MakeGenesisAccountStateProcessor(cs.StateStorage)
	if err != nil {
		panic("open account state processor for genesis failed: " + err.Error())
	}

	genesisRegisterProcessor, err := registerdb.MakeGenesisRegisterProcessor(cs.StateStorage)
	if err != nil {
		panic("make registerDB processor for genesis failed: " + err.Error())
	}
	// setup genesis block
	defaultGenesis := chain.DefaultGenesisBlock(cs.ChainDB, genesisAccountStateProcessor, genesisRegisterProcessor,
		cs.ChainConfig)

	if _, _, err = chain.SetupGenesisBlock(defaultGenesis); err != nil {
		panic("setup genesis block failed: " + err.Error())
	}
}
