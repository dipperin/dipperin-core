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
	"errors"
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/log"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third-party/crypto/cs-crypto"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
	"go.uber.org/zap"
	"math/big"
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
	return nil
}

func (chainDB *ChainDB) GetBlockHashByNumber(number uint64) common.Hash {
	//log.Middleware.Info("chainDB GetBlockHashByNumber start","number",number)
	data, _ := chainDB.db.Get(headerHashKey(number))
	if len(data) == 0 {
		return common.Hash{}
	}
	//log.Middleware.Info("chainDB GetBlockHashByNumber success")
	return common.BytesToHash(data)
}

func (chainDB *ChainDB) SaveBlockHash(hash common.Hash, number uint64) {
	if err := chainDB.db.Put(headerHashKey(number), hash.Bytes()); err != nil {
		log.DLogger.Error("Failed to store number to hash mapping", zap.Error(err))
	}
}

func (chainDB *ChainDB) DeleteBlockHashByNumber(number uint64) {
	if err := chainDB.db.Delete(headerHashKey(number)); err != nil {
		log.DLogger.Error("Failed to delete number to hash mapping", zap.Error(err))
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
		log.DLogger.Error("Failed to store hash to number mapping", zap.Error(err))
	}
}

func (chainDB *ChainDB) DeleteHeaderNumber(hash common.Hash) {
	if err := chainDB.db.Delete(headerNumberKey(hash)); err != nil {
		log.DLogger.Error("Failed to delete hash to number mapping", zap.Error(err))
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
		log.DLogger.Error("Failed to store last block's hash", zap.Error(err))
	}
}

func (chainDB *ChainDB) GetHeaderRLP(hash common.Hash, number uint64) rlp.RawValue {
	data, _ := chainDB.db.Get(headerKey(number, hash))
	return data
}

func (chainDB *ChainDB) SaveHeaderRLP(hash common.Hash, number uint64, rlp rlp.RawValue) {
	if err := chainDB.db.Put(headerKey(number, hash), rlp); err != nil {
		log.DLogger.Error("Failed to store header", zap.Error(err))
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
		log.DLogger.Error("Invalid block header RLP", zap.Any("hash", hash), zap.Error(err))
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
		log.DLogger.Error("Failed to RLP encode header", zap.Error(err))
		return
	}

	chainDB.SaveHeaderRLP(hash, number, rlpData)
}

func (chainDB *ChainDB) DeleteHeader(hash common.Hash, number uint64) {
	if err := chainDB.db.Delete(headerKey(number, hash)); err != nil {
		log.DLogger.Error("Failed to delete header", zap.Error(err))
	}
	chainDB.DeleteHeaderNumber(hash)
}

func (chainDB *ChainDB) GetBodyRLP(hash common.Hash, number uint64) rlp.RawValue {
	data, _ := chainDB.db.Get(blockBodyKey(number, hash))
	return data
}

func (chainDB *ChainDB) SaveBodyRLP(hash common.Hash, number uint64, rlp rlp.RawValue) {
	if err := chainDB.db.Put(blockBodyKey(number, hash), rlp); err != nil {
		log.DLogger.Error("Failed to store block body", zap.Error(err))
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
		log.DLogger.Error("Invalid block body RLP", zap.Any("hash", hash), zap.Error(err))
		return nil
	}

	return body
}

func (chainDB *ChainDB) SaveBody(hash common.Hash, number uint64, body model.AbstractBody) {
	data, err := body.EncodeRlpToBytes()
	if err != nil {
		log.DLogger.Error("Failed to RLP encode body", zap.Error(err))
		return
	}

	chainDB.SaveBodyRLP(hash, number, data)
}

func (chainDB *ChainDB) DeleteBody(hash common.Hash, number uint64) {
	if err := chainDB.db.Delete(blockBodyKey(number, hash)); err != nil {
		log.DLogger.Error("Failed to delete block body", zap.Error(err))
	}
}

func (chainDB *ChainDB) GetBlock(hash common.Hash, number uint64) model.AbstractBlock {
	//log.Middleware.Info("chainDB get block start", "hash", hash.Hex(), "num", number)
	headerRlp := chainDB.GetHeaderRLP(hash, number)
	if len(headerRlp) == 0 {
		log.DLogger.Debug("block header not found")
		return nil
	}
	bodyRlp := chainDB.GetBodyRLP(hash, number)
	if len(bodyRlp) == 0 {
		log.DLogger.Debug("block body not found")
		return nil
	}

	block, err := chainDB.decoder.DecodeRlpBlockFromHeaderAndBodyBytes(headerRlp, bodyRlp)

	if err != nil {
		log.DLogger.Warn("decode block failed", zap.Error(err))
		return nil
	}
	//log.Middleware.Info("chainDB get block end")
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
func (chainDB *ChainDB) SaveReceipts(hash common.Hash, number uint64, receipts model.Receipts) error {
	// Convert the receipts into their storage form and serialize them
	storageReceipts := make([]*model.ReceiptForStorage, len(receipts))
	for i, receipt := range receipts {
		storageReceipts[i] = (*model.ReceiptForStorage)(receipt)
	}
	bytes, err := rlp.EncodeToBytes(storageReceipts)
	if err != nil {
		log.DLogger.Error("Failed to encode block receipts", zap.Error(err))
		return err
	}
	// Store the flattened receipt slice
	if err := chainDB.db.Put(blockReceiptsKey(number, hash), bytes); err != nil {
		log.DLogger.Error("Failed to store block receipts", zap.Error(err))
		return err
	}
	return nil
}

func (chainDB *ChainDB) GetReceipts(hash common.Hash, number uint64) model.Receipts {
	// Retrieve the flattened receipt slice
	data, _ := chainDB.db.Get(blockReceiptsKey(number, hash))
	if len(data) == 0 {
		return nil
	}

	// Convert the receipts from their storage form to their internal representation
	var storageReceipts []*model.ReceiptForStorage
	if err := rlp.DecodeBytes(data, &storageReceipts); err != nil {
		log.DLogger.Error("Invalid receipt array RLP", zap.Any("hash", hash), zap.Error(err))
		return nil
	}

	receipts := make(model.Receipts, len(storageReceipts))
	for i, receipt := range storageReceipts {
		receipts[i] = (*model.Receipt)(receipt)
	}

	// complement receipts
	block := chainDB.GetBlock(hash, number)
	err := DeriveFields(receipts, block)
	if err != nil {
		log.DLogger.Error("DeriveFields failed", zap.Error(err))
		return nil
	}
	return receipts
}

// DeriveFields fills the receipts with their computed fields based on consensus
// data and contextual infos like containing block and transactions.
func DeriveFields(r model.Receipts, block model.AbstractBlock) error {
	logIndex := uint(0)
	txs := block.GetTransactions()
	number := block.Number()
	hash := block.Hash()
	if len(txs) != len(r) {
		return errors.New(fmt.Sprintf("length of txs and receipts not match txs:%v, receipts:%v", len(txs), len(r)))
	}
	for i := 0; i < len(r); i++ {
		// The transaction hash can be retrieved from the transaction itself
		r[i].TxHash = txs[i].CalTxId()

		// block location fields
		r[i].BlockHash = hash
		r[i].BlockNumber = new(big.Int).SetUint64(number)
		r[i].TransactionIndex = uint(i)

		// The contract address can be derived from the transaction itself
		if txs[i].GetType() == common.AddressTypeContractCreate {
			callerAddress, err := txs[i].Sender(nil)
			if err != nil {
				return err
			}
			r[i].ContractAddress = cs_crypto.CreateContractAddress(callerAddress, txs[i].Nonce())
		} else {
			r[i].ContractAddress = *txs[i].To()
		}

		// The used gas can be calculated based on previous r
		if i == 0 {
			r[i].GasUsed = r[i].CumulativeGasUsed
		} else {
			r[i].GasUsed = r[i].CumulativeGasUsed - r[i-1].CumulativeGasUsed
		}
		// The derived log fields can simply be set from the block and transaction
		for j := 0; j < len(r[i].Logs); j++ {
			r[i].Logs[j].BlockNumber = number
			r[i].Logs[j].BlockHash = hash
			r[i].Logs[j].TxHash = r[i].TxHash
			r[i].Logs[j].TxIndex = uint(i)
			r[i].Logs[j].Index = logIndex
			logIndex++
		}
	}
	return nil
}

func (chainDB *ChainDB) GetBloomBits(head common.Hash, bit uint, section uint64) []byte {
	bloomBits, err := chainDB.db.Get(bloomBitsKey(bit, section, head))
	if err != nil {
		log.DLogger.Error("ChainDB#GetBloomBits err", zap.Any("hash", head), zap.Error(err))
		return nil
	}
	return bloomBits
}

func BatchSaveBloomBits(db DatabaseWriter, head common.Hash, bit uint, section uint64, bits []byte) error {
	if err := db.Put(bloomBitsKey(bit, section, head), bits); err != nil {
		log.DLogger.Error("Failed to store bloom bits", zap.Error(err))
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
		log.DLogger.Error("Invalid transaction lookup entry RLP", zap.Any("hash", txHash), zap.Error(err))
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
				log.DLogger.Error("Failed to store transaction lookup entry", zap.Error(err))
				return err
			}

			return nil

		}); err != nil {
			log.DLogger.Error("block tx iterator failed", zap.Error(err))
			return
		}

		if err := batch.Write(); err != nil {
			log.DLogger.Error("tx batch write failed", zap.Error(err))
			return
		}
	}
}

func (chainDB *ChainDB) DeleteTxLookupEntry(block model.AbstractBlock) {

	if block.TxCount() > 0 {
		if err := block.TxIterator(func(index int, tx model.AbstractTransaction) error {
			if err := chainDB.db.Delete(txLookupKey(tx.CalTxId())); err != nil {
				log.DLogger.Error("tx lookup entry delete failed", zap.Error(err))
				return err
			}
			return nil

		}); err != nil {
			log.DLogger.Error("block tx iterator failed", zap.Error(err))
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
		log.DLogger.Error("Transaction referenced missing", zap.Uint64("number", blockNumber), zap.Any("hash", blockHash), zap.Uint64("index", txIndex))
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
		log.DLogger.Error("Failed to RLP encode interlinks", zap.Error(err))
		return
	}

	if err := chainDB.db.Put(interLinkKey(root), rlpData); err != nil {
		log.DLogger.Error("Failed to store interlink", zap.Error(err))
	}
}*/
