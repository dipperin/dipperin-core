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

package dipperin

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/chain/state-processor"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third-party/rpc"
	"github.com/ethereum/go-ethereum/rlp"
)

type NodeInfo struct {
	InProcHandler *rpc.Server // In-process RPC request handler to process the API requests
}

// manage dipperin app lifetime
type Node interface {
	// add extra service before node start
	AddService(service NodeService)
	Start() error
	Stop()
	Wait()
	GetNodeInfo() NodeInfo
}

type VerifiersReader interface {
	CurrentVerifiers() []common.Address
	NextVerifiers() []common.Address
	PrimaryNode() common.Address
	GetPBFTPrimaryNode() common.Address
	VerifiersTotalCount() int
	ShouldChangeVerifier() bool
}

// fixme use cs_chain
type Chain interface {
	InsertBlocks(blocks []model.AbstractBlock) error

	Genesis() model.AbstractBlock

	CurrentBlock() model.AbstractBlock
	CurrentHeader() model.AbstractHeader

	CurrentSeed() (common.Hash, uint64)
	IsChangePoint(block model.AbstractBlock, isProcessPackageBlock bool) bool
	GetLastChangePoint(block model.AbstractBlock) *uint64
	GetSlotByNum(num uint64) *uint64
	GetSlot(block model.AbstractBlock) *uint64

	GetCurrVerifiers() []common.Address
	GetVerifiers(round uint64) []common.Address
	GetNextVerifiers() []common.Address

	GetBlock(hash common.Hash, number uint64) model.AbstractBlock
	GetBlockByHash(hash common.Hash) model.AbstractBlock
	GetBlockByNumber(number uint64) model.AbstractBlock
	GetLatestNormalBlock() model.AbstractBlock

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

	GetStateStorage() state_processor.StateStorage
	CurrentState() (*state_processor.AccountStateDB, error)
	StateAtByBlockNumber(num uint64) (*state_processor.AccountStateDB, error)
	StateAtByStateRoot(root common.Hash) (*state_processor.AccountStateDB, error)

	ValidTx(tx model.AbstractTransaction) error

	GetSeenCommit(height uint64) []model.AbstractVerification
	SaveBlock(block model.AbstractBlock, seenCommits []model.AbstractVerification) error
}
