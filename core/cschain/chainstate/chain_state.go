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

package chainstate

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/gerror"
	"github.com/dipperin/dipperin-core/common/gevent"
	"github.com/dipperin/dipperin-core/common/log"
	"github.com/dipperin/dipperin-core/core/chain"
	"github.com/dipperin/dipperin-core/core/chain/chaindb"
	"github.com/dipperin/dipperin-core/core/chain/stateprocessor"
	"github.com/dipperin/dipperin-core/core/chainconfig"
	"github.com/dipperin/dipperin-core/core/cschain/chainwriter"
	"github.com/dipperin/dipperin-core/core/economymodel"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/ethereum/go-ethereum/ethdb"
	"go.uber.org/zap"
	"path/filepath"
)

// configuration of ChainState
type ChainStateConfig struct {
	ChainConfig   *chainconfig.ChainConfig
	DataDir       string
	WriterFactory chainwriter.AbstractChainWriterFactory
}

// the struct of ChainState
type ChainState struct {
	*ChainStateConfig

	ethDB        ethdb.Database
	ChainDB      chaindb.Database
	StateStorage stateprocessor.StateStorage
	EconomyModel economymodel.EconomyModel
}

//get AccountStateDB
func (cs *ChainState) AccountStateDB(root common.Hash) (*stateprocessor.AccountStateDB, error) {
	aDB, err := stateprocessor.NewAccountStateDB(root, cs.StateStorage)
	if err != nil {
		return nil, err
	}

	return aDB, nil
}

// create a new BlockProcessor according to the root hash
func (cs *ChainState) BlockProcessor(root common.Hash) (*chain.BlockProcessor, error) {
	return chain.NewBlockProcessor(cs, root, cs.StateStorage)
}

// create a new BlockProcessor according to the number
func (cs *ChainState) BlockProcessorByNumber(num uint64) (*chain.BlockProcessor, error) {
	block := cs.GetBlockByNumber(num)
	if block == nil {
		return nil, gerror.ErrBlockNotFound
	}
	return chain.NewBlockProcessor(cs, block.StateRoot(), cs.StateStorage)
}

// create a new ChainState
func NewChainState(conf *ChainStateConfig) *ChainState {
	gevent.Add(gevent.NewBlockInsertEvent)
	cs := &ChainState{ChainStateConfig: conf}
	cs.initConfigAndDB(conf.DataDir)
	cs.WriterFactory = conf.WriterFactory
	cs.WriterFactory.SetChain(cs)
	return cs
}

// get the database of the ChainState
func (cs *ChainState) GetDB() ethdb.Database {
	return cs.ethDB
}

func (cs *ChainState) initConfigAndDB(dataDir string) {
	// init ethdb
	ethDB := initEthDB(dataDir)
	cs.ethDB = ethDB

	// init block decoder
	blockDecoder := model.MakeDefaultBlockDecoder()

	// init chain config
	cs.ChainConfig = chainconfig.GetChainConfig()

	// init chainDB
	cs.ChainDB = chaindb.NewChainDB(ethDB, blockDecoder)

	cs.StateStorage = stateprocessor.NewStateStorageWithCache(ethDB)

	// init economy model
	cs.EconomyModel = economymodel.MakeDipperinEconomyModel(cs, economymodel.DIPProportion)
}

// init database
func initEthDB(dataDir string) ethdb.Database {
	var db ethdb.Database

	switch dataDir {
	case "mem", "test", "":
		db = ethdb.NewMemDatabase()

	default:
		dataDir = filepath.Join(dataDir, "full_chain_data")
		log.DLogger.Info("open chain data", zap.String("dir", dataDir))
		tmpDB, err := ethdb.NewLDBDatabase(dataDir, 0, 0)

		if err != nil {
			panic(err)
		}

		db = tmpDB
	}

	return db
}
