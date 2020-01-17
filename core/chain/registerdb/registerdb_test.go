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
	"github.com/dipperin/dipperin-core/core/chain/stateprocessor"
	"github.com/dipperin/dipperin-core/core/chainconfig"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
	"time"
)

func TestMakeGenesisRegisterProcessor(t *testing.T) {
	db := ethdb.NewMemDatabase()
	storage := stateprocessor.NewStateStorageWithCache(db)
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

	type result struct {
		slot uint64
		err error
	}

	testCases := []struct {
		name   string
		given  func() (uint64,error)
		expect result
	}{
		{
			name:"slot 7",
			given: func() (uint64, error) {
				_, err = registerDB.Commit()
				slot, err := registerDB.GetSlot()
				return slot,err
			},
			expect:result{uint64(7),nil},
		},
		{
			name:"slot 6",
			given:func() (uint64, error) {
				err = registerDB.deleteSlotData()
				_, err = registerDB.Commit()
				slot, err := registerDB.GetSlot()
				return slot,err
			},
			expect:result{uint64(6),nil},
		},
	}

	for _,tc:=range testCases{
		slot,err:=tc.given()
		assert.Equal(t,tc.expect.slot,slot)
		assert.Equal(t,tc.expect.err,err)
	}
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
	err = registerDB.saveRegisterData(common.HexToAddress("0x125"))
	assert.NoError(t, err)
	err = registerDB.saveSlotData(0)
	assert.NoError(t, err)

	_, err = registerDB.Commit()
	assert.NoError(t, err)
	//time.Sleep(time.Millisecond * 100)

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
	slotSize := chainconfig.GetChainConfig().SlotSize
	block := CreateBlock(20)

	type result struct {
		isNot bool
	}

	testCases := []struct {
		name   string
		given  func() bool
		expect result
	}{
		{
			name:"not change point",
			given: func() bool {
				err := registerDB.saveChangePointData(uint64(15))
				_, err = registerDB.Commit()
				assert.NoError(t, err)
				return registerDB.IsChangePoint(block, false)
			},
			expect:result{false},
		},
		{
			name:"change point",
			given: func() bool {
				block = CreateSpecialBlock(20)
				return registerDB.IsChangePoint(block, false)
			},
			expect:result{true},
		},
		{
			name:"block big than slot size,not package block",
			given: func() bool {
				block = CreateBlock(15 + slotSize)
				return registerDB.IsChangePoint(block, false)
			},
			expect:result{true},
		},
		{
			name:"block big than slot size,is package block",
			given: func() bool {
				block = CreateBlock(15 + slotSize)
				return registerDB.IsChangePoint(block, true)
			},
			expect:result{true},
		},
	}

	for _,tc:=range testCases{
		ret:=tc.given()
		assert.Equal(t,tc.expect.isNot,ret)
	}
}

func TestRegisterDB_Process(t *testing.T) {
	registerDB := createRegisterDb(0)

	//　config environment
	slotSize := chainconfig.GetChainConfig().SlotSize
	err := registerDB.PrepareRegisterDB()
	_, err = registerDB.Commit()
	assert.NoError(t, err)
	time.Sleep(time.Millisecond * 100)

	type result struct {
		slot uint64
		point uint64
		err error
	}

	testCases := []struct {
		name   string
		given  func() (uint64,uint64,error)
		expect result
	}{
		{
			name:"insert the block that the block number is 1",
			given: func() (uint64, uint64, error) {
				block := CreateBlock(1)
				err = registerDB.Process(block)
				return 0, 0, err
			},
			expect:result{0,0,nil},
		},
		{
			name:"insert the block that is change point,block num=slot size",
			given: func() (uint64, uint64, error) {
				block := CreateBlock(slotSize)
				err = registerDB.Process(block)
				return 0, 0, err
			},
			expect:result{0,0,nil},
		},
		{
			name:"insert the block that contain unRegister transaction",
			given:func() (uint64, uint64, error) {
				tx1 := createRegisterTX(0, big.NewInt(10000))
				tx2 := createCannelTX(1)
				block := createBlock(2, []*model.Transaction{tx1, tx2})
				err = registerDB.Process(block)
				slot, err := registerDB.GetSlot()
				point, err := registerDB.GetLastChangePoint()
				return slot, point, err
			},
			expect:result{uint64(1),uint64(slotSize-1),nil},
		},
		{
			name:"insert the block that is change point,block number big than slot size",
			given:func() (uint64, uint64, error) {
				block := CreateBlock(2 * slotSize)
				err = registerDB.Process(block)
				slot, err := registerDB.GetSlot()
				point, err := registerDB.GetLastChangePoint()
				return slot, point, err
			},
			expect:result{uint64(2),uint64(2*slotSize-1),nil},
		},
		{
			name:"insert special block",
			given:func() (uint64, uint64, error) {
				specialBlock := CreateSpecialBlock(4*slotSize + 2)
				err := registerDB.Process(specialBlock)
				slot, err := registerDB.GetSlot()
				point, err := registerDB.GetLastChangePoint()
				return slot, point, err
			},
			expect:result{uint64(3),uint64(4*slotSize+1),nil},
		},
	}

	for _,tc:=range testCases{
		slot,point,err:=tc.given()
		assert.Equal(t,tc.expect.slot,slot)
		assert.Equal(t,tc.expect.point,point)
		assert.Equal(t,tc.expect.err,err)
	}
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

	type result struct {
		err error
	}

	testCases := []struct {
		name   string
		given  func() error
		expect result
	}{
		{
			name:"insert the block that contain register transaction",
			given: func() error {
				tx1 := createRegisterTX(0, big.NewInt(10000))
				block := createBlock(2, []*model.Transaction{tx1})
				err := registerDB.Process(block)
				return err
			},
			expect:result{errors.New("tx iterator has error")},
		},
		{
			name:"insert the block that contain unRegister transaction",
			given: func() error {
				tx2 := createCannelTX(0)
				block := createBlock(2, []*model.Transaction{tx2})
				err := registerDB.Process(block)
				return err
			},
			expect:result{errors.New("tx iterator has error")},
		},
		{
			name:"insert the block that not contain transaction",
			given: func() error {
				block := createBlock(2, nil)
				err := registerDB.Process(block)
				return err
			},
			expect:result{errors.New("slot test error")},
		},
		{
			name:"insert the block that contain error register transaction",
			given: func() error {
				tx1 := model.NewRegisterTransaction(0, big.NewInt(10000), model.TestGasPrice, model.TestGasLimit)
				block := createBlock(2, []*model.Transaction{tx1})
				err := registerDB.Process(block)
				return err
			},
			expect:result{errors.New("tx iterator has error")},
		},
		{
			name:"insert the block that contain error cancel transaction",
			given: func() error {
				tx2 := model.NewCancelTransaction(0, model.TestGasPrice, model.TestGasLimit)
				block := createBlock(2, []*model.Transaction{tx2})
				err := registerDB.Process(block)
				return err
			},
			expect:result{errors.New("tx iterator has error")},
		},
	}

	for _,tc:=range testCases{
		err:=tc.given()
		assert.Equal(t,tc.expect.err,err)
	}
}


