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
package stateprocessor

import (
	"errors"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/gerror"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third_party/crypto/cs-crypto"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestApplyMessage_Error(t *testing.T) {
	WASMPath := model.GetWASMPath("event", model.CoreVmTestData)
	AbiPath := model.GetAbiPath("event", model.CoreVmTestData)
	gasLimit := model.TestGasLimit
	tx := createContractTx(WASMPath, AbiPath, 0, gasLimit)
	msg, err := tx.AsMessage(true)
	assert.NoError(t, err)

	db := ethdb.NewMemDatabase()
	testVm := getTestVm(db, common.Hash{})
	gasPool := gasLimit / 2

	tdb := NewStateStorageWithCache(db)
	processor, _ := NewAccountStateDB(common.Hash{}, tdb)
	processor.NewAccountState(aliceAddr)
	processor.AddNonce(aliceAddr, 1)
	root, _ := processor.Commit()

	type result struct {
		err error
	}

	testCases := []struct {
		name   string
		given  func() error
		expect result
	}{
		{
			name:"account does not exist",
			given: func() error {
				_, _, _, _, err = ApplyMessage(testVm, &msg, &gasPool)
				return err
			},
			expect:result{errors.New("account does not exist")},
		},
		{
			name:"ErrNonceTooLow",
			given: func() error {
				testVm = getTestVm(db, root)
				_, _, _, _, err = ApplyMessage(testVm, &msg, &gasPool)
				return err
			},
			expect:result{gerror.ErrNonceTooLow},
		},
		{
			name:"ErrNonceTooHigh",
			given: func() error {
				tx = createContractTx(WASMPath, AbiPath, 2, gasLimit)
				msg, err = tx.AsMessage(true)
				assert.NoError(t, err)
				_, _, _, _, err = ApplyMessage(testVm, &msg, &gasPool)
				return err
			},
			expect:result{gerror.ErrNonceTooHigh},
		},
		{
			name:"ErrInsufficientBalanceForGas",
			given: func() error {
				tx = createContractTx(WASMPath, AbiPath, 1, gasLimit)
				msg, err = tx.AsMessage(true)
				assert.NoError(t, err)
				_, _, _, _, err = ApplyMessage(testVm, &msg, &gasPool)
				return err
			},
			expect:result{gerror.ErrInsufficientBalanceForGas},
		},
		{
			name:"ErrGasLimitReached",
			given: func() error {
				processor, _ = NewAccountStateDB(root, tdb)
				processor.AddBalance(aliceAddr, big.NewInt(0).SetUint64(7000000))
				root, _ = processor.Commit()
				testVm = getTestVm(db, root)
				_, _, _, _, err = ApplyMessage(testVm, &msg, &gasPool)
				return err
			},
			expect:result{gerror.ErrGasLimitReached},
		},
		{
			name:"ErrOutOfGas",
			given: func() error {
				gasPool = gasLimit
				_, _, _, _, err = ApplyMessage(testVm, &msg, &gasPool)
				return err
			},
			expect:result{gerror.ErrOutOfGas},
		},
		{
			name:"VM execute fail",
			given: func() error {
				testVm = getTestVm(db, root)
				gasLimit = model.TestGasLimit * 50
				gasPool = gasLimit * 10
				tx = createContractTx(WASMPath, AbiPath, 1, gasLimit)
				msg, err = tx.AsMessage(true)
				assert.NoError(t, err)
				_, _, _, _, err := ApplyMessage(testVm, &msg, &gasPool)
				return err
			},
			expect:result{errors.New("VM execute fail: abort")},
		},
		{
			name:"ErrInsufficientBalance",
			given: func() error {
				name := []byte("ApplyMsg")
				params := [][]byte{name}
				to := cs_crypto.CreateContractAddress(aliceAddr, 1)
				tx = callContractTx(&to, "returnString", params, 2)
				msg, err = tx.AsMessage(true)
				assert.NoError(t, err)
				_, _, _, _, err = ApplyMessage(testVm, &msg, &gasPool)
				return err
			},
			expect:result{gerror.ErrInsufficientBalance},
		},
	}

	for _,tc:=range testCases{
		err:=tc.given()
		if err!=nil{
			assert.Equal(t,tc.expect.err,err)
		}else {
			assert.NoError(t,err)
		}
	}
}

func BenchmarkApplyMessage_Create(b *testing.B) {
	WASMPath := model.GetWASMPath("event", model.CoreVmTestData)
	AbiPath := model.GetAbiPath("event", model.CoreVmTestData)
	tx := createContractTx(WASMPath, AbiPath, 0, testGasLimit)
	msg, err := tx.AsMessage(true)
	assert.NoError(b, err)
	for i := 0; i < b.N; i++ {
		db, root := CreateTestStateDB()
		testVm := getTestVm(db, root)
		gasPool := uint64(5 * testGasLimit)
		_, usedGas, failed, _, err := ApplyMessage(testVm, &msg, &gasPool)
		assert.NoError(b, err)
		assert.False(b, failed)
		assert.NotNil(b, usedGas)
	}
}

func BenchmarkApplyMessage_Call(b *testing.B) {
	WASMPath := model.GetWASMPath("event", model.CoreVmTestData)
	AbiPath := model.GetAbiPath("event", model.CoreVmTestData)

	// create tx
	tx1 := createContractTx(WASMPath, AbiPath, 0, testGasLimit)
	msg1, err := tx1.AsMessage(true)
	assert.NoError(b, err)

	// call tx
	name := []byte("ApplyMsg")
	params := [][]byte{name}
	to := cs_crypto.CreateContractAddress(aliceAddr, 0)
	tx2 := callContractTx(&to, "returnString", params, 1)
	msg2, err := tx2.AsMessage(true)
	assert.NoError(b, err)

	for i := 0; i < b.N; i++ {
		db, root := CreateTestStateDB()
		testVm := getTestVm(db, root)
		gasPool := uint64(5 * testGasLimit)

		_, usedGas, failed, _, innerErr := ApplyMessage(testVm, &msg1, &gasPool)
		assert.NoError(b, innerErr)
		assert.False(b, failed)
		assert.NotNil(b, usedGas)

		_, usedGas, failed, _, innerErr = ApplyMessage(testVm, &msg2, &gasPool)
		assert.NoError(b, innerErr)
		assert.False(b, failed)
		assert.NotNil(b, usedGas)
	}
}

