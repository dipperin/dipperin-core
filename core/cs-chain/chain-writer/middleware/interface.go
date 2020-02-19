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

package middleware

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/chain"
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/dipperin/dipperin-core/core/chain/chaindb"
	"github.com/dipperin/dipperin-core/core/chain/registerdb"
	"github.com/dipperin/dipperin-core/core/chain/state-processor"
	"github.com/dipperin/dipperin-core/core/economy-model"
	"github.com/dipperin/dipperin-core/core/model"
	model2 "github.com/dipperin/dipperin-core/core/vm/model"
	"github.com/ethereum/go-ethereum/rlp"
)

//go:generate mockgen -destination=./chain_interface_mock_test.go -package=middleware github.com/dipperin/dipperin-core/core/cs-chain/chain-writer/middleware ChainInterface
//go:generate mockgen -destination=../../../../tests/g-mockFile/chainInterface_mock_test.go -package=g_mockFile github.com/dipperin/dipperin-core/core/cs-chain/chain-writer/middleware ChainInterface
type ChainInterface interface {
	StateReader
	StateWriter
	VerifierHelper
	StateHelper
	ChainHelper
}

type StateWriter interface {
	//SaveBlock(block model.AbstractBlock) error
	Rollback(target uint64) error
}

type StateReader interface {
	Genesis() model.AbstractBlock
	CurrentBlock() model.AbstractBlock
	CurrentHeader() model.AbstractHeader
	GetBlock(hash common.Hash, number uint64) model.AbstractBlock
	GetBlockByHash(hash common.Hash) model.AbstractBlock
	GetBlockByNumber(number uint64) model.AbstractBlock
	HasBlock(hash common.Hash, number uint64) bool
	GetBody(hash common.Hash) model.AbstractBody
	GetBodyRLP(hash common.Hash) rlp.RawValue
	GetHeader(hash common.Hash, number uint64) model.AbstractHeader
	GetHeaderByHash(hash common.Hash) model.AbstractHeader
	GetHeaderByNumber(number uint64) model.AbstractHeader
	GetHeaderRLP(hash common.Hash) rlp.RawValue
	HasHeader(hash common.Hash, number uint64) bool
	GetBlockNumber(hash common.Hash) *uint64
	GetTransaction(txHash common.Hash) (model.AbstractTransaction, common.Hash, uint64, uint64)
	GetReceipts(hash common.Hash, number uint64) model2.Receipts
	GetBloomLog(hash common.Hash, number uint64) model2.Bloom
	GetBloomBits(head common.Hash, bit uint, section uint64) []byte
	GetLatestNormalBlock() model.AbstractBlock

	BlockProcessor(root common.Hash) (*chain.BlockProcessor, error)
	BlockProcessorByNumber(num uint64) (*chain.BlockProcessor, error)

	AccountStateDB(root common.Hash) (*state_processor.AccountStateDB, error)
}

type VerifierHelper interface {
	CurrentSeed() (common.Hash, uint64)
	IsChangePoint(block model.AbstractBlock, isProcessPackageBlock bool) bool
	GetLastChangePoint(block model.AbstractBlock) *uint64
	GetSlotByNum(num uint64) *uint64
	GetSlot(block model.AbstractBlock) *uint64
	GetCurrVerifiers() []common.Address
	GetVerifiers(round uint64) []common.Address
	GetNextVerifiers() []common.Address
	NumBeforeLastBySlot(slot uint64) *uint64
	BuildRegisterProcessor(preRoot common.Hash) (*registerdb.RegisterDB, error)
}

type StateHelper interface {
	GetStateStorage() state_processor.StateStorage
	CurrentState() (*state_processor.AccountStateDB, error)
	StateAtByBlockNumber(num uint64) (*state_processor.AccountStateDB, error)
	StateAtByStateRoot(root common.Hash) (*state_processor.AccountStateDB, error)
	BuildStateProcessor(preAccountStateRoot common.Hash) (*state_processor.AccountStateDB, error)
}

type ChainHelper interface {
	GetChainConfig() *chain_config.ChainConfig
	GetEconomyModel() economy_model.EconomyModel
	GetChainDB() chaindb.Database
}
