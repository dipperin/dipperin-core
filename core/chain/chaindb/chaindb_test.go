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
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/dipperin/dipperin-core/core/vm/model"
)

func TestChainDB_InsertBlock(t *testing.T) {
	db := newChainDB()
	b := createBlock(22)

	err := db.InsertBlock(b)
	assert.NoError(t, err)

	assert.Equal(t, true, db.GetBlockHashByNumber(b.Number()).IsEqual(b.Hash()))
	assert.Equal(t, true, db.GetHeadBlockHash().IsEqual(b.Hash()))

	tx, hash, _, _ := db.GetTransaction(b.Body().GetTxByIndex(0).CalTxId())

	assert.NotNil(t, tx)
	assert.Equal(t, true, hash.IsEqual(b.Hash()))
}

func TestCanonicalHash(t *testing.T) {
	db := newChainDB()
	b := createBlock(22)

	db.SaveBlockHash(b.Hash(), b.Number())
	assert.Equal(t, b.Hash(), db.GetBlockHashByNumber(b.Number()))

	db.DeleteBlockHashByNumber(b.Number())
	assert.Equal(t, common.Hash{}, db.GetBlockHashByNumber(b.Number()))
}

func TestCanonicalHash_Error(t *testing.T) {
	fakeDB := NewChainDB(fakeDataBase{}, newDecoder())
	b := createBlock(22)

	fakeDB.SaveBlockHash(b.Hash(), b.Number())
	fakeDB.DeleteBlockHashByNumber(22)
}

func TestWriteHeaderNumber(t *testing.T) {
	db := newChainDB()
	b := createBlock(22)

	db.SaveHeaderNumber(b.Hash(), b.Number())
	assert.Equal(t, b.Number(), *db.GetHeaderNumber(b.Hash()))

	db.DeleteHeaderNumber(b.Hash())
	assert.Nil(t, db.GetHeaderNumber(b.Hash()))
}

func TestWriteHeaderNumber_Error(t *testing.T) {
	fakeDB := NewChainDB(fakeDataBase{}, newDecoder())
	b := createBlock(22)

	fakeDB.SaveHeaderNumber(b.Hash(), b.Number())
	fakeDB.DeleteHeaderNumber(b.Hash())
}

func TestWriteHeadHeaderHash(t *testing.T) {
	db := newChainDB()
	b := createBlock(22)

	assert.Equal(t, common.Hash{}, db.GetHeadHeaderHash())

	db.SaveHeadHeaderHash(b.Hash())
	assert.Equal(t, b.Hash(), db.GetHeadHeaderHash())
}

func TestWriteHeadHeaderHash_Error(t *testing.T) {
	fakeDB := NewChainDB(fakeDataBase{}, newDecoder())
	b := createBlock(22)

	fakeDB.SaveHeadHeaderHash(b.Hash())
}

func TestWriteHeadBlockHash(t *testing.T) {
	db := newChainDB()
	b := createBlock(22)

	assert.Equal(t, common.Hash{}, db.GetHeadBlockHash())

	db.SaveHeadBlockHash(b.Hash())
	assert.Equal(t, b.Hash(), db.GetHeadBlockHash())
}

func TestWriteHeadBlockHash_Error(t *testing.T) {
	fakeDB := NewChainDB(fakeDataBase{}, newDecoder())
	b := createBlock(22)

	fakeDB.SaveHeadBlockHash(b.Hash())
}

func TestWriteHeader(t *testing.T) {
	db := newChainDB()
	b := createBlock(22)

	db.SaveHeader(b.Header())
	assert.Equal(t, b.Hash(), db.GetHeader(b.Hash(), b.Number()).Hash())

	has := db.HasHeader(b.Hash(), b.Number())
	assert.True(t, has)

	db.DeleteHeader(b.Hash(), b.Number())
	assert.Nil(t, db.GetHeader(b.Hash(), b.Number()))

	has = db.HasHeader(b.Hash(), b.Number())
	assert.False(t, has)
}

func TestWriteHeader_Error(t *testing.T) {
	fakeDB := NewChainDB(fakeDataBase{}, fakeDecoder{})
	b := createBlock(22)

	fakeDB.SaveHeader(b.Header())
	fakeDB.DeleteHeader(b.Hash(), b.Number())

	fakeDB = NewChainDB(newDb(), fakeDecoder{})
	fakeDB.SaveHeader(b.Header())
	fakeDB.GetHeader(b.Hash(), b.Number())

	fakeDB.SaveHeader(fakeHeader{})
}

func TestWriteBody(t *testing.T) {
	db := newChainDB()
	b := createBlock(22)

	db.SaveBody(b.Hash(), b.Number(), b.Body())
	assert.Equal(t, 2, db.GetBody(b.Hash(), b.Number()).GetTxsSize())

	has := db.HasBody(b.Hash(), b.Number())
	assert.True(t, has)

	db.DeleteBody(b.Hash(), b.Number())
	assert.Nil(t, db.GetHeader(b.Hash(), b.Number()))

	has = db.HasBody(b.Hash(), b.Number())
	assert.False(t, has)
}

func TestWriteBody_Error(t *testing.T) {
	fakeDB := NewChainDB(fakeDataBase{}, fakeDecoder{})
	b := createBlock(22)

	fakeDB.SaveBody(b.Hash(), b.Number(), b.Body())
	fakeDB.DeleteBody(b.Hash(), b.Number())

	fakeDB = NewChainDB(newDb(), fakeDecoder{})
	fakeDB.SaveBody(b.Hash(), b.Number(), b.Body())
	fakeDB.GetBody(b.Hash(), b.Number())

	fakeDB.SaveBody(b.Hash(), b.Number(), fakeBody{})
}

func TestWriteBlock(t *testing.T) {
	db := newChainDB()
	b := createBlock(22)

	db.SaveBlock(b)
	assert.Equal(t, b.Number(), db.GetBlock(b.Hash(), b.Number()).Number())
	assert.Equal(t, b.Hash(), db.GetBlock(b.Hash(), b.Number()).Hash())

	db.DeleteBlock(b.Hash(), b.Number())
	assert.Nil(t, db.GetBlock(b.Hash(), b.Number()))
}

func TestWriteBlock_Error(t *testing.T) {
	db := newChainDB()
	b := createBlock(22)

	db.SaveBlock(b)
	db.db.Delete(blockBodyKey(b.Number(), b.Hash()))
	db.GetBlock(b.Hash(), b.Number())

	fakeDB := NewChainDB(newDb(), fakeDecoder{})
	fakeDB.SaveBlock(b)
	fakeDB.GetBlock(b.Hash(), b.Number())
}

func TestWriteTxLookupEntries(t *testing.T) {
	db := newChainDB()
	b := createBlock(22)

	db.SaveTxLookupEntries(b)
	txHash := b.GetTransactions()[0].CalTxId()
	bHash, bIndex, txIndex := db.GetTxLookupEntry(txHash)

	assert.Equal(t, b.Hash(), bHash)
	assert.Equal(t, b.Number(), bIndex)
	assert.Equal(t, uint64(0), txIndex)

	txHash = b.GetTransactions()[1].CalTxId()

	bHash, bIndex, txIndex = db.GetTxLookupEntry(txHash)

	assert.Equal(t, b.Hash(), bHash)
	assert.Equal(t, b.Number(), bIndex)
	assert.Equal(t, uint64(1), txIndex)

	db.DeleteTxLookupEntry(b)
	bHash, bIndex, txIndex = db.GetTxLookupEntry(txHash)

	assert.Equal(t, common.Hash{}, bHash)
	assert.Equal(t, uint64(0), bIndex)
	assert.Equal(t, uint64(0), txIndex)
}

func TestWriteTxLookupEntries_Error(t *testing.T) {
	fakeDB := NewChainDB(fakeDataBase{}, fakeDecoder{})
	b := createBlock(22)
	fakeDB.SaveTxLookupEntries(b)
	fakeDB.DeleteTxLookupEntry(b)

	fakeDB = NewChainDB(fakeDataBase{err: BatchErr}, fakeDecoder{})
	fakeDB.SaveTxLookupEntries(b)
	fakeDB.DeleteTxLookupEntry(b)

	db := newChainDB()
	db.db.Put(txLookupKey(common.HexToHash("123")), []byte{1})
	db.GetTxLookupEntry(common.HexToHash("123"))
}

func TestReadTransaction(t *testing.T) {
	db := newChainDB()
	b := createBlock(22)

	txHash := b.GetTransactions()[0].CalTxId()
	tx, bHash, bNumber, txIndex := db.GetTransaction(txHash)
	assert.Nil(t, tx)

	db.SaveBlock(b)
	db.SaveTxLookupEntries(b)
	tx, bHash, bNumber, txIndex = db.GetTransaction(txHash)
	assert.NotNil(t, tx)
	assert.Equal(t, b.Hash(), bHash)
	assert.Equal(t, b.Number(), bNumber)
	assert.Equal(t, uint64(0), txIndex)

	db.DeleteBody(b.Hash(), b.Number())
	tx, bHash, bNumber, txIndex = db.GetTransaction(txHash)
	assert.Nil(t, tx)
}

func TestChainDB_DB(t *testing.T) {
	db := newChainDB()
	assert.NotNil(t, db.DB())
}

func TestChainDB_SaveReceipts(t *testing.T) {
	db := newChainDB()
	receipt := model.NewReceipt(common.HexToHash("0xd36884f1f26ca542ac572832eb02316f3a665571c34653e104eb74bebd2d4335").Bytes(), false, uint64(100))
	receipts := []*model.Receipt{receipt}
	err := db.SaveReceipts(common.Hash{}, 1, receipts)
	assert.NoError(t, err)

	result := db.GetReceipts(common.Hash{}, 1)
	assert.Equal(t, receipts[0], result[0])
}