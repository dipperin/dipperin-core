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
	"github.com/dipperin/dipperin-core/core/chain/registerdb"
	"github.com/dipperin/dipperin-core/core/chain/stateprocessor"
	"github.com/dipperin/dipperin-core/core/chainconfig"
	"github.com/dipperin/dipperin-core/core/cschain/chainstate"
	"github.com/dipperin/dipperin-core/core/cschain/chainwriter"
)

func SetupGenesis(dataDir string, cConfig *chainconfig.ChainConfig) {
	cs := chainstate.NewChainState(&chainstate.ChainStateConfig{
		ChainConfig:   cConfig,
		DataDir:       dataDir,
		WriterFactory: chainwriter.NewChainWriterFactory(),
	})
	setupGenesis(cs)
	cs.ChainDB.DB().Close()
}

func setupGenesis(cs *chainstate.ChainState) {
	genesisAccountStateProcessor, err := stateprocessor.MakeGenesisAccountStateProcessor(cs.StateStorage)
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
