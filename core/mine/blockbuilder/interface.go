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

package blockbuilder

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/chain"
	"github.com/dipperin/dipperin-core/core/chain/registerdb"
	"github.com/dipperin/dipperin-core/core/chaincommunication"
	"github.com/dipperin/dipperin-core/core/chainconfig"
	"github.com/dipperin/dipperin-core/core/economymodel"
	"github.com/dipperin/dipperin-core/core/model"
)

type AbstractBlockBuilder interface {
	BuildWaitPackBlock(coinbaseAddr common.Address, gasFloor, gasCeil uint64) model.AbstractBlock
}

type Chain interface {
	CurrentBlock() model.AbstractBlock
	GetBlockByNumber(number uint64) model.AbstractBlock
	GetVerifiers(round uint64) []common.Address

	IsChangePoint(block model.AbstractBlock, isProcessPackageBlock bool) bool
	GetLastChangePoint(block model.AbstractBlock) *uint64
	GetSlot(block model.AbstractBlock) *uint64
	GetSeenCommit(height uint64) []model.AbstractVerification
	GetLatestNormalBlock() model.AbstractBlock

	BlockProcessor(root common.Hash) (*chain.BlockProcessor, error)
	BuildRegisterProcessor(preRoot common.Hash) (*registerdb.RegisterDB, error)
	GetEconomyModel() economymodel.EconomyModel
}

//go:generate mockgen -destination=./tx_pool_test.go -package=builder github.com/dipperin/dipperin-core/core/model/builder TxPool
type TxPool interface {
	RemoveTxs(newBlock model.AbstractBlock)
	Pending() (map[common.Address][]model.AbstractTransaction, error)
	AddRemotes(txs []model.AbstractTransaction) []error
	Reset(oldHead, newHead *model.Header)
}

type ModelConfig struct {
	ChainReader        Chain
	TxPool             TxPool
	PriorityCalculator model.PriofityCalculator
	TxSigner           model.Signer
	MsgSigner          chaincommunication.PbftSigner
	ChainConfig        chainconfig.ChainConfig
}

//go:generate mockgen -destination=./signer_mock_test.go -package=builder github.com/dipperin/dipperin-core/core/model Signer

//go:generate mockgen -destination=./pbft_signer_mock_test.go -package=builder github.com/dipperin/dipperin-core/core/chaincommunication PbftSigner
