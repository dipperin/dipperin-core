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

package chaindb

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
)

type Database interface {
	DB() ethdb.Database

	GetBlockHashByNumber(number uint64) common.Hash
	SaveBlockHash(hash common.Hash, number uint64)
	DeleteBlockHashByNumber(number uint64)

	GetHeaderNumber(hash common.Hash) *uint64
	SaveHeaderNumber(hash common.Hash, number uint64)
	DeleteHeaderNumber(hash common.Hash)

	GetHeadBlockHash() common.Hash
	SaveHeadBlockHash(hash common.Hash)

	GetHeaderRLP(hash common.Hash, number uint64) rlp.RawValue
	SaveHeaderRLP(hash common.Hash, number uint64, rlp rlp.RawValue)
	HasHeader(hash common.Hash, number uint64) bool
	GetHeader(hash common.Hash, number uint64) model.AbstractHeader
	SaveHeader(header model.AbstractHeader)
	DeleteHeader(hash common.Hash, number uint64)

	GetBodyRLP(hash common.Hash, number uint64) rlp.RawValue
	SaveBodyRLP(hash common.Hash, number uint64, rlp rlp.RawValue)
	HasBody(hash common.Hash, number uint64) bool
	GetBody(hash common.Hash, number uint64) model.AbstractBody
	SaveBody(hash common.Hash, number uint64, body model.AbstractBody)
	DeleteBody(hash common.Hash, number uint64)

	GetBlock(hash common.Hash, number uint64) model.AbstractBlock
	SaveBlock(block model.AbstractBlock)
	DeleteBlock(hash common.Hash, number uint64)

	//FindCommonAncestor(a, b model.AbstractHeader) model.AbstractHeader

	GetTxLookupEntry(txHash common.Hash) (common.Hash, uint64, uint64)
	SaveTxLookupEntries(block model.AbstractBlock)
	DeleteTxLookupEntry(block model.AbstractBlock)

	GetTransaction(txHash common.Hash) (model.AbstractTransaction, common.Hash, uint64, uint64)

	InsertBlock(block model.AbstractBlock) error

	SaveReceipts(hash common.Hash, number uint64, receipts model.Receipts) error
	GetReceipts(hash common.Hash, number uint64) model.Receipts
	//SaveBloomBits(head common.Hash, bit uint, section uint64,  bits []byte) error
	GetBloomBits(head common.Hash, bit uint, section uint64) []byte
}
