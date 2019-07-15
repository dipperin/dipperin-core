// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package chain_state

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/model"
	model2 "github.com/dipperin/dipperin-core/core/vm/model"
)

// NewTxsEvent is posted when a batch of transactions enter the transaction pool.
type NewTxsEvent struct{ Txs []*model.Transaction }

// PendingLogsEvent is posted pre mining and notifies of pending logs.
type PendingLogsEvent struct {
	Logs []*model2.Log
}

// NewMinedBlockEvent is posted when a block has been imported.
type NewMinedBlockEvent struct{ Block *model.Block }

// RemovedLogsEvent is posted when a reorg happens
type RemovedLogsEvent struct{ Logs []*model2.Log }

type ChainEvent struct {
	Block *model.Block
	Hash  common.Hash
	Logs  []*model2.Log
}

type ChainSideEvent struct {
	Block *model.Block
}

type ChainHeadEvent struct{ Block *model.Block }
