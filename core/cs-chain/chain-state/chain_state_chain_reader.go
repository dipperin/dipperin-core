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
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/model"
	model2 "github.com/dipperin/dipperin-core/core/vm/model"
	"github.com/ethereum/go-ethereum/rlp"
)

/*
	all the operations here do not consider the cache, just take it directly from the db
*/

func (cs *ChainState) Genesis() model.AbstractBlock {
	// here do not consider what to cache, go directly to the db, in the outer layer to consider the cache Genesis
	return cs.GetBlockByNumber(0)
}

func (cs *ChainState) CurrentBlock() model.AbstractBlock {
	headHash := cs.ChainDB.GetHeadBlockHash()
	if headHash.IsEmpty() {
		return nil
	}

	return cs.GetBlockByHash(headHash)
}

func (cs *ChainState) CurrentHeader() model.AbstractHeader {
	headHash := cs.ChainDB.GetHeadBlockHash()

	if headHash.IsEmpty() {
		return nil
	}

	if tmpBlock := cs.GetBlockByHash(headHash); tmpBlock != nil {
		return tmpBlock.Header()
	}

	return nil
}

func (cs *ChainState) GetBlock(hash common.Hash, number uint64) model.AbstractBlock {
	return cs.ChainDB.GetBlock(hash, number)
}

func (cs *ChainState) GetBlockByHash(hash common.Hash) model.AbstractBlock {
	number := cs.GetBlockNumber(hash)
	if number == nil {
		return nil
	}

	return cs.GetBlock(hash, *number)
}

func (cs *ChainState) GetBlockByNumber(number uint64) model.AbstractBlock {
	hash := cs.ChainDB.GetBlockHashByNumber(number)
	if hash == (common.Hash{}) {
		return nil
	}
	return cs.GetBlock(hash, number)
}

func (cs *ChainState) HasBlock(hash common.Hash, number uint64) bool {
	return cs.ChainDB.HasBody(hash, number)
}

func (cs *ChainState) GetBody(hash common.Hash) model.AbstractBody {
	number := cs.GetBlockNumber(hash)

	if number == nil {
		return nil
	}

	body := cs.ChainDB.GetBody(hash, *number)

	if body == nil {
		return nil
	}

	return body
}

func (cs *ChainState) GetBodyRLP(hash common.Hash) rlp.RawValue {
	number := cs.GetBlockNumber(hash)
	if number == nil {
		return nil
	}

	bodyRLP := cs.ChainDB.GetBodyRLP(hash, *number)

	if len(bodyRLP) == 0 {
		return nil
	}

	return bodyRLP
}

func (cs *ChainState) GetHeader(hash common.Hash, number uint64) model.AbstractHeader {
	header := cs.ChainDB.GetHeader(hash, number)

	if header == nil {
		return nil
	}

	return header
}

func (cs *ChainState) GetHeaderByHash(hash common.Hash) model.AbstractHeader {
	number := cs.GetBlockNumber(hash)
	if number == nil {
		return nil
	}
	return cs.GetHeader(hash, *number)
}

func (cs *ChainState) GetHeaderByNumber(number uint64) model.AbstractHeader {
	hash := cs.ChainDB.GetBlockHashByNumber(number)
	if hash == (common.Hash{}) {
		return nil
	}

	return cs.GetHeader(hash, number)
}

func (cs *ChainState) GetHeaderRLP(hash common.Hash) rlp.RawValue {
	number := cs.GetBlockNumber(hash)
	if number == nil {
		return nil
	}

	headerRLP := cs.ChainDB.GetHeaderRLP(hash, *number)

	if len(headerRLP) == 0 {
		return nil
	}

	return headerRLP
}

func (cs *ChainState) HasHeader(hash common.Hash, number uint64) bool {
	return cs.ChainDB.HasHeader(hash, number)
}

func (cs *ChainState) GetBlockNumber(hash common.Hash) *uint64 {
	return cs.ChainDB.GetHeaderNumber(hash)
}

func (cs *ChainState) GetTransaction(txHash common.Hash) (model.AbstractTransaction, common.Hash, uint64, uint64) {
	return cs.ChainDB.GetTransaction(txHash)
}

func (cs *ChainState) GetReceipts(hash common.Hash, number uint64) model2.Receipts {
	return cs.ChainDB.GetReceipts(hash, number)
}


func (cs *ChainState)GetBloomBits(head common.Hash, bit uint, section uint64) []byte {
	return cs.ChainDB.GetBloomBits(head , bit , section)
}

func (cs *ChainState) GetLatestNormalBlock() model.AbstractBlock {
	findBlock := cs.CurrentBlock()
	for {
		if findBlock.IsSpecial() {
			findBlock = cs.GetBlockByNumber(findBlock.Number() - 1)
		} else {
			break
		}
	}

	return findBlock
}
