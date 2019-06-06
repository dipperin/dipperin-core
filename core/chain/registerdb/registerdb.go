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

package registerdb

import (
	"errors"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/chain/state-processor"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/dipperin/dipperin-core/third-party/log/pbft_log"
	"github.com/dipperin/dipperin-core/third-party/trie"
	"github.com/ethereum/go-ethereum/rlp"
	"bytes"
	"github.com/dipperin/dipperin-core/core/chain-config"
)

var (
	TxIteratorError = errors.New("tx iterator has error")
)

const (
	slotKey            = "slot_key"
	lastChangePointKey = "last_change_point_key"
)

type RegisterDB struct {
	trie state_processor.StateTrie
	// nodeContext NodeContext
	storage state_processor.StateStorage
	//chainReader state_processor.ChainReader
	chainReader ChainReader
}

func MakeGenesisRegisterProcessor(storage state_processor.StateStorage) (RegisterProcessor, error) {
	preStateRoot := common.Hash{}
	return NewRegisterDB(preStateRoot, storage, nil)
}

func NewRegisterDB(preRoot common.Hash, storage state_processor.StateStorage, chainReader ChainReader) (*RegisterDB, error) {
	if t, err := storage.OpenTrie(preRoot); err != nil {
		return nil, err
	} else {
		registerDB := &RegisterDB{
			trie:        t,
			storage:     storage,
			chainReader: chainReader,
		}
		return registerDB, nil
	}
}

func (register RegisterDB) PrepareRegisterDB() error {

	if err := register.saveSlotData(uint64(0)); err != nil {
		return err
	}
	if err := register.saveChangePointData(uint64(0)); err != nil {
		return err
	}
	return nil
}

//　get block slot
func (register RegisterDB) GetSlot() (uint64, error) {

	var res uint64
	enc, err1 := register.trie.TryGet([]byte(slotKey))
	if err1 != nil {
		return res, err1
	}

	err2 := rlp.DecodeBytes(enc, &res)
	if err2 != nil {
		return res, err2
	}
	return res, nil
}

// the block height of last change point
func (register RegisterDB) GetLastChangePoint() (uint64, error) {

	var res uint64
	enc, err1 := register.trie.TryGet([]byte(lastChangePointKey))
	if err1 != nil {
		return res, err1
	}
	err2 := rlp.DecodeBytes(enc, &res)
	if err2 != nil {
		return res, err2
	}

	return res, nil
}

//　all register data
func (register RegisterDB) GetRegisterData() []common.Address {

	it := trie.NewIterator(register.trie.NodeIterator(nil))
	var list []common.Address
	for it.Next() {
		key := register.trie.GetKey(it.Key)
		if bytes.Equal(key, []byte(lastChangePointKey)) || bytes.Equal(key, []byte(slotKey)) {
			continue
		} else {
			list = append(list, common.BytesToAddress(key))
		}
	}
	return list
}

func (register RegisterDB) Process(block model.AbstractBlock) (err error) {
	log.Debug("r db process", "tx len", block.TxCount(), "block num", block.Number())
	//　get all register data
	if err := block.TxIterator(func(index int, tx model.AbstractTransaction) (error) {
		switch tx.GetType() {
		case common.AddressTypeCancel:
			sender, innerError := tx.Sender(tx.GetSigner())
			if innerError != nil {
				log.Error("tx.Sender(tx.GetSigner()) has error", "tx index", index, "err", innerError)
				return innerError
			}
			if err := register.deleteRegisterData(sender); err != nil {
				return err
			}
			pbft_log.Info("deletedRegisterData", "sender", sender)
			return nil

		case common.AddressTypeStake:
			sender, innerError := tx.Sender(tx.GetSigner())
			log.Debug("register db deal stake", "sender", sender)
			if innerError != nil {
				log.Error("tx.Sender(tx.GetSigner()) has error", "tx index", index, "err", innerError)
				return innerError
			}
			if err := register.saveRegisterData(sender); err != nil {
				return err
			}
			pbft_log.Info("savedRegisterData", "sender", sender)
			return nil

		default:
			return nil
		}

	}); err != nil {
		return TxIteratorError
	}

	if block.Number() <= 1 {
		return nil
	}

	preBlock := register.chainReader.GetBlockByNumber(block.Number() - 1)
	if err = register.processSlot(preBlock); err != nil {
		return err
	}

	return nil
}

func (register RegisterDB) processSlot(preBlock model.AbstractBlock) error {

	slot, err := register.GetSlot()
	if err != nil {
		log.Info("get slot failed", "err", err)
		return err
	}

	//update slot and last change point info if the block is special or the change point
	if register.IsChangePoint(preBlock, false) {
		if err = register.saveSlotData(slot + 1); err != nil {
			return err
		}
		if err = register.saveChangePointData(preBlock.Number()); err != nil {
			return err
		}
		log.Info("save slot successful", "cur slot", slot+1, "change point", preBlock.Number())
	}
	return nil
}

// judge if block is change point
func (register RegisterDB) IsChangePoint(block model.AbstractBlock, isProcessPackageBlock bool) bool {
	point, err := register.GetLastChangePoint()
	if err != nil {
		log.Debug("GetLastChangePoint failed", "err", err)
		return false
	}

	diff := block.Number() - point
	if point == 0 {
		diff += 1
	}

	slotSize := chain_config.GetChainConfig().SlotSize

	//the mineMaster packaged Block diff and nonce is 0
	if isProcessPackageBlock {
		if diff == slotSize {
			return true
		}
	} else {
		if block.IsSpecial() || diff == slotSize {
			return true
		}
	}

	return false
}

func (register RegisterDB) Finalise() common.Hash {
	return register.trie.Hash()
}

func (register RegisterDB) Commit() (root common.Hash, err error) {
	if root, err = register.trie.Commit(nil); err != nil {
		return root, err
	}
	err = register.storage.TrieDB().Commit(root, false)
	return root, err
}

func (register RegisterDB) saveRegisterData(addr common.Address) error {
	return register.trie.TryUpdate(addr.Bytes(), []byte{0})
}

func (register RegisterDB) deleteRegisterData(addr common.Address) error {
	return register.trie.TryDelete(addr.Bytes())
}

func (register RegisterDB) saveSlotData(slot uint64) error {
	newEnc, _ := rlp.EncodeToBytes(slot)
	return register.trie.TryUpdate([]byte(slotKey), newEnc)
}

func (register RegisterDB) deleteSlotData() error {
	slot, err := register.GetSlot()
	if err != nil {
		return err
	}

	return register.saveSlotData(slot - 1)
}

func (register RegisterDB) saveChangePointData(lastChangePoint uint64) error {
	newEnc, _ := rlp.EncodeToBytes(lastChangePoint)
	return register.trie.TryUpdate([]byte(lastChangePointKey), newEnc)
}
