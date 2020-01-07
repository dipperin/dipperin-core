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
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/dipperin/dipperin-core/core/chain/state-processor"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/tests/factory"
	"github.com/dipperin/dipperin-core/tests/factory/g-testData"
	"github.com/dipperin/dipperin-core/third-party/crypto"
	"github.com/dipperin/dipperin-core/third-party/trie"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
	"time"

	"bytes"
	"errors"
)

var (
	SlotError  = errors.New("slot test error")
	PointError = errors.New("point test error")
	TestError  = errors.New("test error")
	alicePriv  = "289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232031"
)

func TestMakeGenesisRegisterProcessor(t *testing.T) {
	db := ethdb.NewMemDatabase()
	storage := state_processor.NewStateStorageWithCache(db)
	processor, err := MakeGenesisRegisterProcessor(storage)
	assert.NoError(t, err)
	assert.NotNil(t, processor)

	register, err := NewRegisterDB(common.HexToHash("123"), storage, nil)
	assert.Nil(t, register)
	assert.Error(t, err)
}

func TestRegisterDB_PrepareRegisterDB(t *testing.T) {
	registerDB := createRegisterDb(0)
	err := registerDB.PrepareRegisterDB()
	assert.NoError(t, err)

	slot, err := registerDB.GetSlot()
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), slot)

	point, err := registerDB.GetLastChangePoint()
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), point)
}

func TestRegisterDB_GetSlot(t *testing.T) {
	registerDB := createRegisterDb(0)

	err := registerDB.saveSlotData(uint64(5))
	assert.NoError(t, err)
	err = registerDB.saveSlotData(uint64(10))
	assert.NoError(t, err)
	err = registerDB.saveSlotData(uint64(7))
	assert.NoError(t, err)

	_, err = registerDB.Commit()
	assert.NoError(t, err)
	time.Sleep(time.Millisecond * 100)

	slot, err := registerDB.GetSlot()
	assert.NoError(t, err)
	assert.Equal(t, slot, uint64(7))

	err = registerDB.deleteSlotData()
	assert.NoError(t, err)

	_, err = registerDB.Commit()
	assert.NoError(t, err)

	slot, err = registerDB.GetSlot()
	assert.NoError(t, err)
	assert.Equal(t, slot, uint64(6))
}

func TestRegisterDB_GetSlot_Error(t *testing.T) {
	registerDB := createRegisterDb(0)
	registerDB.trie = fakeTrie{}

	_, err := registerDB.GetSlot()
	assert.Error(t, err)
}

func TestRegisterDB_GetLastChangePoint(t *testing.T) {
	registerDB := createRegisterDb(0)

	err := registerDB.saveChangePointData(uint64(5))
	assert.NoError(t, err)
	err = registerDB.saveChangePointData(uint64(10))
	assert.NoError(t, err)
	err = registerDB.saveChangePointData(uint64(7))
	assert.NoError(t, err)

	_, err = registerDB.Commit()
	assert.NoError(t, err)
	time.Sleep(time.Millisecond * 100)

	slot, err := registerDB.GetLastChangePoint()
	assert.NoError(t, err)
	assert.Equal(t, slot, uint64(7))
}

func TestRegisterDB_GetRegisterData(t *testing.T) {
	registerDB := createRegisterDb(0)

	err := registerDB.saveRegisterData(common.HexToAddress("0x123"))
	assert.NoError(t, err)
	err = registerDB.saveRegisterData(common.HexToAddress("0x124"))
	assert.NoError(t, err)
	err = registerDB.deleteRegisterData(common.HexToAddress("0x124"))
	assert.NoError(t, err)
	err = registerDB.saveRegisterData(common.HexToAddress("0x125"))
	assert.NoError(t, err)
	err = registerDB.saveSlotData(0)
	assert.NoError(t, err)

	_, err = registerDB.Commit()
	assert.NoError(t, err)
	time.Sleep(time.Millisecond * 100)

	data1 := registerDB.GetRegisterData()
	assert.Len(t, data1, 2)

	err = registerDB.deleteRegisterData(common.HexToAddress("0x125"))
	assert.NoError(t, err)

	_, err = registerDB.Commit()
	assert.NoError(t, err)
	data2 := registerDB.GetRegisterData()
	assert.Len(t, data2, 1)
}

func TestRegisterDB_IsChangePoint(t *testing.T) {
	registerDB := createRegisterDb(0)
	slotSize := chain_config.GetChainConfig().SlotSize
	block := factory.CreateBlock(20)
	assert.Equal(t, registerDB.IsChangePoint(block, false), false)

	err := registerDB.saveChangePointData(uint64(15))
	assert.NoError(t, err)

	_, err = registerDB.Commit()
	assert.NoError(t, err)
	time.Sleep(time.Millisecond * 100)

	assert.Equal(t, registerDB.IsChangePoint(block, false), false)

	block = factory.CreateSpecialBlock(20)
	assert.Equal(t, registerDB.IsChangePoint(block, false), true)

	block = factory.CreateBlock(15 + slotSize)
	assert.Equal(t, registerDB.IsChangePoint(block, false), true)
	assert.Equal(t, registerDB.IsChangePoint(block, true), true)
}

func TestRegisterDB_Process(t *testing.T) {
	registerDB := createRegisterDb(0)

	//　config environment
	slotSize := chain_config.GetChainConfig().SlotSize
	err := registerDB.PrepareRegisterDB()
	_, err = registerDB.Commit()
	assert.NoError(t, err)
	time.Sleep(time.Millisecond * 100)

	// insert the block that the block number is 1
	block := factory.CreateBlock(1)
	err = registerDB.Process(block)
	assert.NoError(t, err)

	//　insert the block that is change point
	block = factory.CreateBlock(slotSize)
	err = registerDB.Process(block)
	assert.NoError(t, err)

	// insert the block that contain unRegister transaction
	tx1 := createRegisterTX(0, big.NewInt(10000))
	tx2 := createCannelTX(1)
	block = createBlock(2, []*model.Transaction{tx1, tx2})
	err = registerDB.Process(block)
	assert.NoError(t, err)

	slot, err := registerDB.GetSlot()
	assert.NoError(t, err)
	assert.Equal(t, slot, uint64(1))
	point, err := registerDB.GetLastChangePoint()
	assert.NoError(t, err)
	assert.Equal(t, point, uint64(slotSize-1))

	//　insert the block that is change point
	block = factory.CreateBlock(2 * slotSize)
	err = registerDB.Process(block)
	assert.NoError(t, err)

	slot, err = registerDB.GetSlot()
	assert.NoError(t, err)
	assert.Equal(t, slot, uint64(2))
	point, err = registerDB.GetLastChangePoint()
	assert.NoError(t, err)
	assert.Equal(t, point, uint64(2*slotSize-1))

	//　insert special block
	specialBlock := factory.CreateSpecialBlock(4*slotSize + 2)
	err = registerDB.Process(specialBlock)
	assert.NoError(t, err)
	slot, err = registerDB.GetSlot()
	assert.NoError(t, err)
	assert.Equal(t, slot, uint64(3))
	point, err = registerDB.GetLastChangePoint()
	assert.NoError(t, err)
	assert.Equal(t, point, uint64(4*slotSize+1))
}

func TestRegisterDB_Finalise(t *testing.T) {
	registerDB := createRegisterDb(0)
	hash := registerDB.Finalise()
	assert.NotNil(t, hash)
}

func TestRegisterDB_Commit(t *testing.T) {
	registerDB := createRegisterDb(0)
	root, err := registerDB.Commit()
	assert.NoError(t, err)
	assert.NotNil(t, root)
}

func TestRegisterDB_Error(t *testing.T) {
	registerDB := createRegisterDb(0)
	registerDB.trie = fakeTrie{
		slotErr:  SlotError,
		pointErr: PointError,
	}

	// insert the block that contain register transaction
	tx1 := createRegisterTX(0, big.NewInt(10000))
	block := createBlock(2, []*model.Transaction{tx1})
	err := registerDB.Process(block)
	assert.Error(t, err)

	// insert the block that contain unRegister transaction
	tx2 := createCannelTX(0)
	block = createBlock(2, []*model.Transaction{tx2})
	err = registerDB.Process(block)
	assert.Error(t, err)

	// insert the block that not contain transaction
	block = createBlock(2, nil)
	err = registerDB.Process(block)
	assert.Error(t, err)

	// insert the block that contain error transaction
	tx1 = model.NewRegisterTransaction(0, big.NewInt(10000), g_testData.TestGasPrice, g_testData.TestGasLimit)
	block = createBlock(2, []*model.Transaction{tx1})
	err = registerDB.Process(block)
	assert.Error(t, err)

	tx2 = model.NewCancelTransaction(0, g_testData.TestGasPrice, g_testData.TestGasLimit)
	block = createBlock(2, []*model.Transaction{tx2})
	err = registerDB.Process(block)
	assert.Error(t, err)
}

func TestRegisterDB_Error2(t *testing.T) {
	registerDB := createRegisterDb(0)
	registerDB.trie = fakeTrie{
		slotErr:  SlotError,
		pointErr: PointError,
	}
	root, err := registerDB.Commit()
	assert.Equal(t, SlotError, err)
	assert.NotNil(t, root)

	slot, err := registerDB.GetSlot()
	assert.Equal(t, SlotError, err)
	assert.NotNil(t, slot)

	point, err := registerDB.GetLastChangePoint()
	assert.Equal(t, PointError, err)
	assert.NotNil(t, point)

	err = registerDB.deleteSlotData()
	assert.Equal(t, SlotError, err)
	assert.NotNil(t, slot)

	err = registerDB.PrepareRegisterDB()
	assert.Equal(t, SlotError, err)

	registerDB.trie = fakeTrie{pointErr: PointError}

	err = registerDB.PrepareRegisterDB()
	assert.Equal(t, PointError, err)
}

func createRegisterDb(blockNum uint64) *RegisterDB {
	db := ethdb.NewMemDatabase()
	storage := state_processor.NewStateStorageWithCache(db)
	block := factory.CreateBlock(blockNum)
	reader := factory.NewFakeReader(block)
	register, _ := NewRegisterDB(common.Hash{}, storage, reader)
	return register
}

func createBlock(number uint64, txs []*model.Transaction) *model.Block {
	header := model.NewHeader(1, number, common.HexToHash("0x12312fa0929348"), common.HexToHash("1111"), common.HexToDiff("1fffffff"), big.NewInt(time.Now().UnixNano()), common.HexToAddress("032f14"), common.BlockNonceFromInt(432423))
	return model.NewBlock(header, txs, []model.AbstractVerification{})
}

func createRegisterTX(nonce uint64, amount *big.Int) *model.Transaction {
	fs1 := model.NewSigner(big.NewInt(1))
	tx := model.NewRegisterTransaction(nonce, amount, g_testData.TestGasPrice, g_testData.TestGasLimit)
	key, _ := crypto.HexToECDSA(alicePriv)
	signedTx, _ := tx.SignTx(key, fs1)
	return signedTx
}

func createCannelTX(nonce uint64) *model.Transaction {
	fs1 := model.NewSigner(big.NewInt(1))
	tx := model.NewCancelTransaction(nonce, g_testData.TestGasPrice, g_testData.TestGasLimit)
	key, _ := crypto.HexToECDSA(alicePriv)
	signedTx, _ := tx.SignTx(key, fs1)
	return signedTx
}

type fakeTrie struct {
	slotErr  error
	pointErr error
}

func (f fakeTrie) TryGet(key []byte) ([]byte, error) {
	if bytes.Equal(key, []byte(slotKey)) {
		return nil, f.slotErr
	}

	if bytes.Equal(key, []byte(lastChangePointKey)) {
		return nil, f.pointErr
	}
	return nil, nil
}

func (f fakeTrie) TryUpdate(key, value []byte) error {
	if bytes.Equal(key, []byte(slotKey)) {
		return f.slotErr
	}

	if bytes.Equal(key, []byte(lastChangePointKey)) {
		return f.pointErr
	}
	return TestError
}

func (f fakeTrie) TryDelete(key []byte) error {
	return TestError
}

func (f fakeTrie) Commit(onleaf trie.LeafCallback) (common.Hash, error) {
	return common.Hash{}, SlotError
}

func (f fakeTrie) Hash() common.Hash {
	panic("implement me")
}

func (f fakeTrie) NodeIterator(startKey []byte) trie.NodeIterator {
	panic("implement me")
}

func (f fakeTrie) GetKey([]byte) []byte {
	panic("implement me")
}

func (f fakeTrie) Prove(key []byte, fromLevel uint, proofDb ethdb.Putter) error {
	panic("implement me")
}
