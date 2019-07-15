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
	"encoding/binary"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/model"
	model2 "github.com/dipperin/dipperin-core/core/vm/model"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
)

type ChainDB struct {
	db      ethdb.Database
	decoder model.BlockDecoder
}

func NewChainDB(db ethdb.Database, decoder model.BlockDecoder) *ChainDB {
	return &ChainDB{
		db:      db,
		decoder: decoder,
	}
}

func (chainDB *ChainDB) DB() ethdb.Database {
	return chainDB.db
}

func (chainDB *ChainDB) InsertBlock(block model.AbstractBlock) error {

	chainDB.SaveBlock(block)
	chainDB.SaveTxLookupEntries(block)

	chainDB.SaveBlockHash(block.Hash(), block.Number())
	chainDB.SaveHeadBlockHash(block.Hash())
	chainDB.SaveHeadHeaderHash(block.Header().Hash())
	return nil
}

func (chainDB *ChainDB) GetBlockHashByNumber(number uint64) common.Hash {
	data, _ := chainDB.db.Get(headerHashKey(number))
	if len(data) == 0 {
		return common.Hash{}
	}
	return common.BytesToHash(data)
}

func (chainDB *ChainDB) SaveBlockHash(hash common.Hash, number uint64) {
	if err := chainDB.db.Put(headerHashKey(number), hash.Bytes()); err != nil {
		log.Crit("Failed to store number to hash mapping", "err", err)
	}
}

func (chainDB *ChainDB) DeleteBlockHashByNumber(number uint64) {
	if err := chainDB.db.Delete(headerHashKey(number)); err != nil {
		log.Crit("Failed to delete number to hash mapping", "err", err)
	}
}

func (chainDB *ChainDB) GetHeaderNumber(hash common.Hash) *uint64 {
	data, _ := chainDB.db.Get(headerNumberKey(hash))
	if len(data) != 8 {
		return nil
	}
	number := binary.BigEndian.Uint64(data)
	return &number
}

func (chainDB *ChainDB) SaveHeaderNumber(hash common.Hash, number uint64) {
	if err := chainDB.db.Put(headerNumberKey(hash), encodeBlockNumber(number)); err != nil {
		log.Crit("Failed to store hash to number mapping", "err", err)
	}
}

func (chainDB *ChainDB) DeleteHeaderNumber(hash common.Hash) {
	if err := chainDB.db.Delete(headerNumberKey(hash)); err != nil {
		log.Crit("Failed to delete hash to number mapping", "err", err)
	}
}

func (chainDB *ChainDB) GetHeadHeaderHash() common.Hash {
	data, _ := chainDB.db.Get(headHeaderKey)
	if len(data) == 0 {
		return common.Hash{}
	}
	return common.BytesToHash(data)
}

func (chainDB *ChainDB) SaveHeadHeaderHash(hash common.Hash) {
	if err := chainDB.db.Put(headHeaderKey, hash.Bytes()); err != nil {
		log.Crit("Failed to store last header's hash", "err", err)
	}
}

func (chainDB *ChainDB) GetHeadBlockHash() common.Hash {
	data, _ := chainDB.db.Get(headBlockKey)
	if len(data) == 0 {
		return common.Hash{}
	}
	return common.BytesToHash(data)
}

func (chainDB *ChainDB) SaveHeadBlockHash(hash common.Hash) {
	if err := chainDB.db.Put(headBlockKey, hash.Bytes()); err != nil {
		log.Crit("Failed to store last block's hash", "err", err)
	}
}

func (chainDB *ChainDB) GetHeaderRLP(hash common.Hash, number uint64) rlp.RawValue {
	data, _ := chainDB.db.Get(headerKey(number, hash))
	return data
}

func (chainDB *ChainDB) SaveHeaderRLP(hash common.Hash, number uint64, rlp rlp.RawValue) {
	if err := chainDB.db.Put(headerKey(number, hash), rlp); err != nil {
		log.Crit("Failed to store header", "err", err)
	}
}

func (chainDB *ChainDB) HasHeader(hash common.Hash, number uint64) bool {
	if has, err := chainDB.db.Has(headerKey(number, hash)); !has || err != nil {
		return false
	}
	return true
}

func (chainDB *ChainDB) GetHeader(hash common.Hash, number uint64) model.AbstractHeader {
	data := chainDB.GetHeaderRLP(hash, number)
	if len(data) == 0 {
		return nil
	}

	header, err := chainDB.decoder.DecodeRlpHeaderFromBytes(data)
	if err != nil {
		log.Error("Invalid block header RLP", "hash", hash, "err", err)
		return nil
	}

	return header
}

func (chainDB *ChainDB) SaveHeader(header model.AbstractHeader) {
	// Write the hash -> number mapping
	var (
		hash   = header.Hash()
		number = header.GetNumber()
	)
	chainDB.SaveHeaderNumber(hash, number)

	rlpData, err := header.EncodeRlpToBytes()
	if err != nil {
		log.Crit("Failed to RLP encode header", "err", err)
		return
	}

	chainDB.SaveHeaderRLP(hash, number, rlpData)
}

func (chainDB *ChainDB) DeleteHeader(hash common.Hash, number uint64) {
	if err := chainDB.db.Delete(headerKey(number, hash)); err != nil {
		log.Crit("Failed to delete header", "err", err)
	}
	chainDB.DeleteHeaderNumber(hash)
}

func (chainDB *ChainDB) GetBodyRLP(hash common.Hash, number uint64) rlp.RawValue {
	data, _ := chainDB.db.Get(blockBodyKey(number, hash))
	return data
}

func (chainDB *ChainDB) SaveBodyRLP(hash common.Hash, number uint64, rlp rlp.RawValue) {
	if err := chainDB.db.Put(blockBodyKey(number, hash), rlp); err != nil {
		log.Crit("Failed to store block body", "err", err)
	}
}

func (chainDB *ChainDB) HasBody(hash common.Hash, number uint64) bool {
	if has, err := chainDB.db.Has(blockBodyKey(number, hash)); !has || err != nil {
		return false
	}
	return true
}

func (chainDB *ChainDB) GetBody(hash common.Hash, number uint64) model.AbstractBody {
	data := chainDB.GetBodyRLP(hash, number)
	if len(data) == 0 {
		return nil
	}

	body, err := chainDB.decoder.DecodeRlpBodyFromBytes(data)
	if err != nil {
		log.Error("Invalid block body RLP", "hash", hash, "err", err)
		return nil
	}

	return body
}

func (chainDB *ChainDB) SaveBody(hash common.Hash, number uint64, body model.AbstractBody) {
	data, err := body.EncodeRlpToBytes()
	if err != nil {
		log.Crit("Failed to RLP encode body", "err", err)
		return
	}

	chainDB.SaveBodyRLP(hash, number, data)
}

func (chainDB *ChainDB) DeleteBody(hash common.Hash, number uint64) {
	if err := chainDB.db.Delete(blockBodyKey(number, hash)); err != nil {
		log.Crit("Failed to delete block body", "err", err)
	}
}

func (chainDB *ChainDB) GetBlock(hash common.Hash, number uint64) model.AbstractBlock {
	//log.Debug("chainDB get block", "hash", hash.Hex(), "num", number)
	headerRlp := chainDB.GetHeaderRLP(hash, number)
	if len(headerRlp) == 0 {
		log.Debug("block header not found")
		return nil
	}
	bodyRlp := chainDB.GetBodyRLP(hash, number)
	if len(bodyRlp) == 0 {
		log.Debug("block body not found")
		return nil
	}

	block, err := chainDB.decoder.DecodeRlpBlockFromHeaderAndBodyBytes(headerRlp, bodyRlp)

	if err != nil {
		log.Warn("decode block failed", "err", err)
		return nil
	}

	return block
}

func (chainDB *ChainDB) SaveBlock(block model.AbstractBlock) {
	chainDB.SaveHeader(block.Header())
	chainDB.SaveBody(block.Hash(), block.Number(), block.Body())
}

func (chainDB *ChainDB) DeleteBlock(hash common.Hash, number uint64) {
	chainDB.DeleteHeader(hash, number)
	chainDB.DeleteBody(hash, number)
}

//add receipts save get and delete
func (chainDB *ChainDB) SaveReceipts(hash common.Hash, number uint64, receipts model2.Receipts) error {
	// Convert the receipts into their storage form and serialize them
	storageReceipts := make([]*model2.ReceiptForStorage, len(receipts))
	for i, receipt := range receipts {
		storageReceipts[i] = (*model2.ReceiptForStorage)(receipt)
	}
	bytes, err := rlp.EncodeToBytes(storageReceipts)
	if err != nil {
		log.Crit("Failed to encode block receipts", "err", err)
		return err
	}
	// Store the flattened receipt slice
	if err := chainDB.db.Put(blockReceiptsKey(number, hash), bytes); err != nil {
		log.Crit("Failed to store block receipts", "err", err)
		return err
	}
	return nil
}

func (chainDB *ChainDB) DeleteReceipts(hash common.Hash, number uint64) {
	if err := chainDB.db.Delete(blockReceiptsKey(number, hash)); err != nil {
		log.Crit("Failed to delete block receipt", "err", err)
	}
}

func (chainDB *ChainDB) GetReceipts(hash common.Hash, number uint64) model2.Receipts {
	// Retrieve the flattened receipt slice
	data, _ := chainDB.db.Get(blockReceiptsKey(number, hash))
	if len(data) == 0 {
		return nil
	}

	// Convert the receipts from their storage form to their internal representation
	storageReceipts := []*model2.ReceiptForStorage{}
	if err := rlp.DecodeBytes(data, &storageReceipts); err != nil {
		log.Error("Invalid receipt array RLP", "hash", hash, "err", err)
		return nil
	}

	receipts := make(model2.Receipts, len(storageReceipts))
	for i, receipt := range storageReceipts {
		receipts[i] = (*model2.Receipt)(receipt)
	}
	return receipts
}

func (chainDB *ChainDB) GetBloomBits(head common.Hash, bit uint, section uint64) ([]byte){
	bloomBits, err :=  chainDB.db.Get(bloomBitsKey(bit, section, head))
	if err != nil {
		log.Error("ChainDB#GetBloomBits err", "hash", head, "err", err)
		return nil
	}
	return bloomBits
}

func BatchSaveBloomBits(db DatabaseWriter, head common.Hash, bit uint, section uint64,  bits []byte) error {
	if err := db.Put(bloomBitsKey(bit, section, head), bits); err != nil {
		log.Error("Failed to store bloom bits", "err", err)
		return err
	}
	return nil
}

/*func (chainDB *ChainDB) FindCommonAncestor(a, b model.AbstractHeader) model.AbstractHeader {
	for bn := b.GetNumber(); a.GetNumber() > bn; {
		a = chainDB.GetHeader(a.GetPreHash(), a.GetNumber()-1)
		if a == nil {
			return nil
		}
	}
	for an := a.GetNumber(); an < b.GetNumber(); {
		b = chainDB.GetHeader(b.GetPreHash(), b.GetNumber()-1)
		if b == nil {
			return nil
		}
	}
	for a.Hash() != b.Hash() {
		a = chainDB.GetHeader(a.GetPreHash(), a.GetNumber()-1)
		if a == nil {
			return nil
		}
		b = chainDB.GetHeader(b.GetPreHash(), b.GetNumber()-1)
		if b == nil {
			return nil
		}
	}
	return a
}*/

func (chainDB *ChainDB) GetTxLookupEntry(txHash common.Hash) (common.Hash, uint64, uint64) {
	data, _ := chainDB.db.Get(txLookupKey(txHash))
	if len(data) == 0 {
		return common.Hash{}, 0, 0
	}
	var entry TxLookupEntry
	if err := rlp.DecodeBytes(data, &entry); err != nil {
		log.Error("Invalid transaction lookup entry RLP", "hash", txHash, "err", err)
		return common.Hash{}, 0, 0
	}
	return entry.BlockHash, entry.BlockIndex, entry.Index
}

func (chainDB *ChainDB) SaveTxLookupEntries(block model.AbstractBlock) {
	if block.TxCount() > 0 {
		batch := chainDB.db.NewBatch()

		if err := block.TxIterator(func(index int, tx model.AbstractTransaction) error {
			entry := TxLookupEntry{
				BlockHash:  block.Hash(),
				BlockIndex: block.Number(),
				Index:      uint64(index),
			}

			data, _ := rlp.EncodeToBytes(entry)
			if err := batch.Put(txLookupKey(tx.CalTxId()), data); err != nil {
				log.Crit("Failed to store transaction lookup entry", "err", err)
				return err
			}

			return nil

		}); err != nil {
			log.Error("block tx iterator failed", "err", err)
			return
		}

		if err := batch.Write(); err != nil {
			log.Error("tx batch write failed", "err", err)
			return
		}
	}
}

func (chainDB *ChainDB) DeleteTxLookupEntry(block model.AbstractBlock) {

	if block.TxCount() > 0 {
		if err := block.TxIterator(func(index int, tx model.AbstractTransaction) error {
			if err := chainDB.db.Delete(txLookupKey(tx.CalTxId())); err != nil {
				log.Error("tx lookup entry delete failed", "err", err)
				return err
			}
			return nil

		}); err != nil {
			log.Error("block tx iterator failed", "err", err)
			return
		}
	}
}

func (chainDB *ChainDB) GetTransaction(txHash common.Hash) (model.AbstractTransaction, common.Hash, uint64, uint64) {
	blockHash, blockNumber, txIndex := chainDB.GetTxLookupEntry(txHash)
	if blockHash == (common.Hash{}) {
		return nil, common.Hash{}, 0, 0
	}

	body := chainDB.GetBody(blockHash, blockNumber)

	if body == nil || body.GetTxsSize() <= int(txIndex) {
		log.Error("Transaction referenced missing", "number", blockNumber, "hash", blockHash, "index", txIndex)
		return nil, common.Hash{}, 0, 0
	}
	return body.GetTxByIndex(int(txIndex)), blockHash, blockNumber, txIndex
}

/*func (chainDB *ChainDB) GetInterLink(root common.Hash) (model.InterLink, error) {
	return nil, nil
}

func (chainDB *ChainDB) SaveInterLink(root common.Hash, link model.InterLink) {
	rlpData, err := rlp.EncodeToBytes(link)

	if err != nil {
		log.Crit("Failed to RLP encode interlinks", "err", err)
		return
	}

	if err := chainDB.db.Put(interLinkKey(root), rlpData); err != nil {
		log.Crit("Failed to store interlink", "err", err)
	}
}*/
