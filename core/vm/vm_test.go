// Copyright 2019, Keychain Foundation Ltd.
// This file is part of the Dipperin-core library.
//
// The Dipperin-core library is free software: you can redistribute
// it and/or modify it under the terms of the GNU Lesser General Public License
// as published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// The Dipperin-core library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package vm

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/core/vm/common/utils"
	"github.com/dipperin/dipperin-core/tests/mock/model"
	"github.com/dipperin/dipperin-core/tests/util"
	"github.com/dipperin/dipperin-core/third-party/crypto"
	"github.com/dipperin/dipperin-core/third-party/crypto/cs-crypto"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/vntchain/go-vnt/rlp"
	"math/big"
	"testing"
	"time"
)

var (
	context Context
)

func init() {

}

func Test_NewVMContext(t *testing.T) {
	ctrl, db, _ := GetBaseVmInfo(t)
	defer ctrl.Finish()

	tx := model_mock.NewMockAbstractTransaction(ctrl)
	header := model_mock.NewMockAbstractHeader(ctrl)
	singer := model_mock.NewMockSigner(ctrl)
	tx.EXPECT().Sender(singer).Return(model.AliceAddr, nil).AnyTimes()
	tx.EXPECT().GetSigner().Return(singer).AnyTimes()
	tx.EXPECT().GetGasLimit().Return(uint64(chain_config.BlockGasLimit)).AnyTimes()
	tx.EXPECT().GetGasPrice().Return(model.TestGasPrice).AnyTimes()
	tx.EXPECT().CalTxId().Return(common.Hash{}).AnyTimes()
	header.EXPECT().GetNumber().Return(uint64(1)).AnyTimes()
	header.EXPECT().CoinBaseAddress().Return(model.AliceAddr).AnyTimes()
	header.EXPECT().GetTimeStamp().Return(big.NewInt(time.Now().UnixNano())).AnyTimes()

	context = NewVMContext(tx, header, getTestHashFunc())
	assert.Equal(t, tx.GetGasLimit(), context.GetGasLimit())
	assert.Equal(t, header.GetNumber(), context.GetBlockNumber().Uint64())
	assert.Equal(t, tx.GetGasPrice(), context.GetGasPrice())
	assert.Equal(t, common.Hash{}, context.GetBlockHash(0))
	assert.Equal(t, tx.CalTxId(), context.GetTxHash())
	assert.Equal(t, header.CoinBaseAddress(), context.GetCoinBase())
	assert.Equal(t, model.AliceAddr, context.GetOrigin())
	assert.Equal(t, header.GetTimeStamp(), context.GetTime())

	db.EXPECT().GetBalance(model.AliceAddr).Return(big.NewInt(400)).Times(2)
	db.EXPECT().GetBalance(model.ContractAddr).Return(big.NewInt(100)).Times(1)
	db.EXPECT().SubBalance(model.AliceAddr, big.NewInt(100)).Return(nil).Times(1)
	db.EXPECT().AddBalance(model.ContractAddr, big.NewInt(100)).Return(nil).Times(1)

	result := context.CanTransfer(db, model.AliceAddr, big.NewInt(100))
	assert.Equal(t, true, result)

	context.Transfer(db, model.AliceAddr, model.ContractAddr, big.NewInt(100))
	assert.Equal(t, big.NewInt(400), db.GetBalance(model.AliceAddr))
	assert.Equal(t, big.NewInt(100), db.GetBalance(model.ContractAddr))
}

func Test_Run(t *testing.T) {
	ctrl, db, vm := GetBaseVmInfo(t)
	defer ctrl.Finish()

	type result struct {
		res []byte
		err error
	}
	code, abi := test_util.GetTestData("event")
	paramInput, err := rlp.EncodeToBytes([]interface{}{"winner"})
	param := "winner"
	assert.NoError(t, err)

	testCases := []struct {
		name   string
		given  func() ([]byte, error)
		expect result
	}{
		{
			name: "abiOrCodeErr",
			given: func() ([]byte, error) {
				contract := getContract(code, abi, nil)
				contract.ABI = []byte{}
				return run(vm, contract, true)
			},
			expect: result{nil, nil},
		},
		{
			name: "newVMErr",
			given: func() ([]byte, error) {
				contract := getContract(code, abi, nil)
				contract.Code = []byte{12, 23}
				return run(vm, contract, true)
			},
			expect: result{nil, errors.New("unexpected EOF")},
		},
		{
			name: "runCreateRight",
			given: func() ([]byte, error) {
				contract := getContract(code, abi, nil)
				contract.value = big.NewInt(0)
				return run(vm, contract, true)
			},
			expect: result{code, nil},
		},
		{
			name: "runCallRight",
			given: func() ([]byte, error) {

				input, err := rlp.EncodeToBytes([]interface{}{"returnString", "winner"})
				assert.NoError(t, err)

				contract := getContract(code, abi, input)
				contract.value = big.NewInt(0)
				log := model.Log{
					Address:     contract.Address(),
					Topics:      []common.Hash{common.HexToHash("0x5ef0c22ad5a85e4c701253956114eeac26c27503bd523ad6fbc3ac2d4553e69c")},
					TopicName:   "topic",
					BlockNumber: 0,
					Data:        paramInput,
					TxHash:      common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
					TxIndex:     uint(0),
					BlockHash:   common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
					Index:       uint(0),
					Removed:     false,
				}
				db.EXPECT().AddLog(&log).Return().Times(1)

				ret, err := run(vm, contract, false)
				if len(ret) > len(param) {
					ret = ret[:len(param)]
				}
				return ret, err
			},
			expect: result{[]byte(param), nil},
		},
	}

	for _, tc := range testCases {
		ret, err := tc.given()
		if err != nil {
			assert.Equal(t, tc.expect.err.Error(), err.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tc.expect.res, ret)
		}
	}

	/*	_, err := run(vm, contract, true)
		assert.NoError(t, err)

		ch := make(chan int, 1)
		go func() {
			vm.Cancel()
			ch <- 0
		}()
		_, err = run(vm, contract, true)
		assert.NoError(t, err)
		<-ch*/
}

func TestVM_CreateAndCall(t *testing.T) {
	ctrl, db, vm := GetBaseVmInfo(t)
	defer ctrl.Finish()

	ref := AccountRef(model.AliceAddr)
	gasLimit := model.TestGasLimit * 100
	value := big.NewInt(0)
	code, abi := test_util.GetTestData("event")
	rlpParams := []interface{}{
		code, abi,
	}
	data, err := rlp.EncodeToBytes(rlpParams)
	assert.NoError(t, err)

	// vm.Create
	contractAddr := cs_crypto.CreateContractAddress(ref.Address(), uint64(0))
	log := model.Log{
		//Address:common.HexToAddress("0x0014B5Df12F50295469Fe33951403b8f4E63231Ef488"),
		Address:     model.ContractAddr,
		TopicName:   "topic",
		Topics:      []common.Hash{common.BytesToHash(crypto.Keccak256([]byte("topic")))},
		Data:        []byte("ƅparam"),
		BlockNumber: 0,
		TxHash:      common.Hash{},
		TxIndex:     0,
		BlockHash:   common.Hash{},
		Index:       0,
		Removed:     false,
	}

	log2 := log
	log2.Address = contractAddr
	t.Log("log", log.Address)
	t.Log("log2", log2.Address)

	db.EXPECT().GetNonce(ref.Address()).Return(uint64(0), nil).AnyTimes()
	db.EXPECT().GetBalance(ref.Address()).Return(new(big.Int).Mul(value, big.NewInt(10))).AnyTimes()
	db.EXPECT().AddNonce(ref.Address(), uint64(1)).Return().AnyTimes()
	db.EXPECT().GetCodeHash(contractAddr).Return(common.Hash{}).AnyTimes()
	db.EXPECT().GetAbiHash(contractAddr).Return(common.Hash{}).AnyTimes()
	db.EXPECT().GetNonce(contractAddr).Return(uint64(0), nil).AnyTimes()
	db.EXPECT().Snapshot().Return(1).AnyTimes()
	db.EXPECT().CreateAccount(contractAddr).Return(nil).AnyTimes()
	db.EXPECT().SubBalance(ref.Address(), value).Return(nil).AnyTimes()
	db.EXPECT().AddBalance(contractAddr, value).Return(nil).AnyTimes()
	db.EXPECT().SetCode(contractAddr, code).AnyTimes()
	db.EXPECT().SetAbi(contractAddr, abi).AnyTimes()
	db.EXPECT().Exist(contractAddr).Return(true).AnyTimes()
	db.EXPECT().GetCode(contractAddr).Return(code).AnyTimes()
	db.EXPECT().GetAbi(contractAddr).Return(abi).AnyTimes()
	db.EXPECT().AddLog(&log2).Times(1)
	db.EXPECT().AddLog(&log).Times(1)

	testCases := []struct {
		name   string
		given  func() error
		expect error
	}{
		{
			name: "TestContractCreate",
			given: func() error {
				resp, addr, _, err := vm.Create(ref, data, gasLimit, value)
				//expectAddr := cs_crypto.CreateContractAddress(ref.Address(), uint64(0))
				assert.Equal(t, code, resp)
				assert.Equal(t, contractAddr, addr)
				return err
			},
			expect: nil,
		},
		{
			name: "TestContractCall",
			given: func() error {
				resp, addr, gasLimit, err := vm.Create(ref, data, gasLimit, value)
				assert.NoError(t, err)

				input := []interface{}{
					"returnString",
					"param",
				}

				inputData, err := rlp.EncodeToBytes(input)
				assert.NoError(t, err)

				resp, _, err = vm.Call(ref, addr, inputData, gasLimit, big.NewInt(0))
				expectResp := utils.Align32BytesConverter(resp, "string")
				assert.Equal(t, "param", expectResp)
				assert.Equal(t, contractAddr, addr)
				return err
			},
			expect: nil,
		},
		{
			name: "TestContractDelegateCall",
			given: func() error {
				resp, addr, gasLimit, err := vm.Create(ref, data, gasLimit, value)
				assert.NoError(t, err)

				input := []interface{}{
					"returnString",
					"param",
				}

				inputData, err := rlp.EncodeToBytes(input)
				assert.NoError(t, err)

				parentContract := getContract(code, abi, inputData)
				parentContract.value = big.NewInt(0)
				resp, _, err = vm.DelegateCall(parentContract, addr, inputData, gasLimit)
				expectResp := utils.Align32BytesConverter(resp, "string")
				assert.Equal(t, "param", expectResp)
				assert.Equal(t, contractAddr, addr)
				return err
			},
			expect: nil,
		},
	}

	for _, tc := range testCases {
		err := tc.given()
		assert.Equal(t, err, tc.expect)
	}

}

/*
func TestVM_CreateAndCallWithdraw(t *testing.T) {
	vm := getTestVm()
	aliceRef := AccountRef(aliceAddr)
	gasLimit := g_testData.TestGasLimit * 100
	value := g_testData.TestValue
	WASMPath := g_testData.GetWASMPath("token-payable", g_testData.CoreVmTestData)
	AbiPath := g_testData.GetAbiPath("token-payable", g_testData.CoreVmTestData)
	code, abi := g_testData.GetCodeAbi(WASMPath, AbiPath)

	// create alice and bob account
	vm.GetStateDB().CreateAccount(aliceAddr)
	vm.GetStateDB().AddBalance(aliceAddr, big.NewInt(500))
	vm.GetStateDB().CreateAccount(bobAddr)
	vm.GetStateDB().AddBalance(bobAddr, big.NewInt(1000))
	expectAddr := cs_crypto.CreateContractAddress(aliceAddr, uint64(0))

	tokenName := []byte("tokenName")
	symbolName := []byte("symbolName")
	supply := utils.Uint64ToBytes(500)
	data, err := rlp.EncodeToBytes([]interface{}{code, abi, tokenName, symbolName, supply})
	assert.NoError(t, err)

	// init() is not payable function, couldn't transfer DIP
	resp, addr, _, err := vm.Create(aliceRef, data, gasLimit, value)
	assert.Equal(t, []byte(nil), resp)
	assert.Equal(t, expectAddr, addr)
	assert.Equal(t, "VM execute fail: abort", err.Error())
	fmt.Println("-----------------------------------------------------")

	// init() is not payable function, run with value == 0
	expectAddr = cs_crypto.CreateContractAddress(aliceAddr, uint64(1))
	resp, addr, _, err = vm.Create(aliceRef, data, gasLimit, big.NewInt(0))
	assert.Equal(t, code, resp)
	assert.Equal(t, expectAddr, addr)
	assert.NoError(t, err)

	// transfer token to bob and transfer DIP to contract
	funcName := []byte("transfer")
	to := []byte(bobAddr.Hex())
	amount := utils.Uint64ToBytes(100)
	data, err = rlp.EncodeToBytes([]interface{}{funcName, to, amount})
	assert.NoError(t, err)
	resp, _, err = vm.Call(aliceRef, addr, data, gasLimit, value)
	expectResp := utils.Align32BytesConverter(resp, "bool")
	assert.Equal(t, true, expectResp)
	assert.NoError(t, err)
	assert.Equal(t, value, vm.GetStateDB().GetBalance(addr))

	// call withdraw
	funcName = []byte("withdraw")
	data, err = rlp.EncodeToBytes([]interface{}{funcName})
	assert.NoError(t, err)

	// withdraw() is not payable function, couldn't transfer DIP
	resp, _, err = vm.Call(AccountRef(aliceAddr), addr, data, gasLimit, value)
	assert.Equal(t, []byte(nil), resp)
	assert.Equal(t, "VM execute fail: abort", err.Error())

	// bob can't withdraw alice's contract
	resp, _, err = vm.Call(AccountRef(bobAddr), addr, data, gasLimit, big.NewInt(0))
	assert.Equal(t, []byte(nil), resp)
	assert.Equal(t, "VM execute fail: abort", err.Error())

	// alice withdraw the contract
	resp, _, err = vm.Call(aliceRef, addr, data, gasLimit, big.NewInt(0))
	expectResp = utils.Align32BytesConverter(resp, "bool")
	assert.Equal(t, true, expectResp)
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), vm.GetStateDB().GetBalance(addr).Uint64())
}

func TestVM_CreateAndCallToken_Transfer(t *testing.T) {
	vm := getTestVm()
	aliceRef := AccountRef(aliceAddr)
	gasLimit := g_testData.TestGasLimit * 100
	value := g_testData.TestValue
	WASMPath := g_testData.GetWASMPath("token-payable", g_testData.CoreVmTestData)
	AbiPath := g_testData.GetAbiPath("token-payable", g_testData.CoreVmTestData)
	code, abi := g_testData.GetCodeAbi(WASMPath, AbiPath)

	tokenName := []byte("tokenName")
	symbolName := []byte("symbolName")
	supply := utils.Uint64ToBytes(500)
	data, err := rlp.EncodeToBytes([]interface{}{code, abi, tokenName, symbolName, supply})
	assert.NoError(t, err)

	vm.GetStateDB().CreateAccount(aliceRef.Address())
	vm.GetStateDB().AddBalance(aliceRef.Address(), big.NewInt(10000))
	expectAddr := cs_crypto.CreateContractAddress(aliceRef.Address(), uint64(0))
	resp, addr, gasLimit, err := vm.Create(aliceRef, data, gasLimit, big.NewInt(0))
	assert.Equal(t, code, resp)
	assert.Equal(t, expectAddr, addr)
	assert.NoError(t, err)

	// vm.Call transfer to zero address(failed)
	funcName := []byte("transfer")
	to := []byte("0x0")
	amount := utils.Uint64ToBytes(100)
	data, err = rlp.EncodeToBytes([]interface{}{funcName, to, amount})
	assert.NoError(t, err)
	resp, _, err = vm.Call(aliceRef, addr, data, gasLimit, value)
	assert.Equal(t, []byte(nil), resp)
	assert.Equal(t, "VM execute fail: abort", err.Error())

	// vm.Call transfer to bob
	to = []byte(bobAddr.Hex())
	data, err = rlp.EncodeToBytes([]interface{}{funcName, to, amount})
	assert.NoError(t, err)
	resp, _, err = vm.Call(aliceRef, addr, data, gasLimit, value)
	expectResp := utils.Align32BytesConverter(resp, "bool")
	assert.Equal(t, true, expectResp)
	assert.NoError(t, err)

	// vm.Call bob getBalance
	funcName = []byte("getBalance")
	data, err = rlp.EncodeToBytes([]interface{}{funcName, to})
	assert.NoError(t, err)
	resp, _, err = vm.Call(aliceRef, addr, data, gasLimit, big.NewInt(0))
	expectResp = utils.Align32BytesConverter(resp, "uint64")
	assert.Equal(t, uint64(100), expectResp)
	assert.NoError(t, err)

	// vm.Call alice getBalance
	from := []byte(aliceAddr.Hex())
	data, err = rlp.EncodeToBytes([]interface{}{funcName, from})
	assert.NoError(t, err)

	resp, _, err = vm.Call(aliceRef, addr, data, gasLimit, big.NewInt(0))
	expectResp = utils.Align32BytesConverter(resp, "uint64")
	assert.Equal(t, uint64(400), expectResp)
	assert.NoError(t, err)
}

func TestVM_CreateAndCallToken_TransferFrom(t *testing.T) {
	vm := getTestVm()
	aliceRef := AccountRef(aliceAddr)
	bobRef := AccountRef(bobAddr)
	charlieRef := AccountRef(charlieAddr)
	gasLimit := g_testData.TestGasLimit * 100
	WASMPath := g_testData.GetWASMPath("token-payable", g_testData.CoreVmTestData)
	AbiPath := g_testData.GetAbiPath("token-payable", g_testData.CoreVmTestData)
	code, abi := g_testData.GetCodeAbi(WASMPath, AbiPath)

	tokenName := []byte("tokenName")
	symbolName := []byte("symbolName")
	supply := utils.Uint64ToBytes(500)
	data, err := rlp.EncodeToBytes([]interface{}{code, abi, tokenName, symbolName, supply})
	assert.NoError(t, err)

	vm.GetStateDB().CreateAccount(aliceRef.Address())
	vm.GetStateDB().CreateAccount(bobRef.Address())
	vm.GetStateDB().CreateAccount(charlieRef.Address())
	vm.GetStateDB().AddBalance(aliceRef.Address(), big.NewInt(10000))
	expectAddr := cs_crypto.CreateContractAddress(aliceRef.Address(), uint64(0))
	resp, addr, gasLimit, err := vm.Create(aliceRef, data, gasLimit, big.NewInt(0))
	assert.Equal(t, code, resp)
	assert.Equal(t, expectAddr, addr)
	assert.NoError(t, err)

	// vm.Call alice approve 1000 to bob (not enough token)
	funcName := []byte("approve")
	spender := []byte(bobAddr.Hex())
	amount1 := utils.Uint64ToBytes(1000)
	data, err = rlp.EncodeToBytes([]interface{}{funcName, spender, amount1})
	assert.NoError(t, err)
	resp, _, err = vm.Call(aliceRef, addr, data, gasLimit, big.NewInt(0))
	assert.Equal(t, []byte(nil), resp)
	assert.Equal(t, "VM execute fail: abort", err.Error())

	// vm.Call alice approve 400 to bob
	amount2 := utils.Uint64ToBytes(400)
	data, err = rlp.EncodeToBytes([]interface{}{funcName, spender, amount2})
	assert.NoError(t, err)
	resp, _, err = vm.Call(aliceRef, addr, data, gasLimit, big.NewInt(0))
	expectResp := utils.Align32BytesConverter(resp, "bool")
	assert.Equal(t, true, expectResp)
	assert.NoError(t, err)

	// vm.Call charlie getApproveBalance
	funcName = []byte("getApproveBalance")
	from := []byte(aliceAddr.Hex())
	getApproveData, err := rlp.EncodeToBytes([]interface{}{funcName, from, spender})
	assert.NoError(t, err)
	resp, _, err = vm.Call(charlieRef, addr, getApproveData, gasLimit, big.NewInt(0))
	expectResp = utils.Align32BytesConverter(resp, "uint64")
	assert.Equal(t, uint64(400), expectResp)
	assert.NoError(t, err)

	// vm.Call bob transferFrom 1000 (not enough approval token)
	funcName = []byte("transferFrom")
	data, err = rlp.EncodeToBytes([]interface{}{funcName, from, spender, amount1})
	assert.NoError(t, err)
	resp, _, err = vm.Call(bobRef, addr, data, gasLimit, big.NewInt(0))
	assert.Equal(t, []byte(nil), resp)
	assert.Equal(t, "VM execute fail: abort", err.Error())

	// vm.Call charlie transferFrom the token which alice approve to bob (failed)
	data, err = rlp.EncodeToBytes([]interface{}{funcName, from, spender, amount2})
	assert.NoError(t, err)
	resp, _, err = vm.Call(charlieRef, addr, data, gasLimit, big.NewInt(0))
	assert.Equal(t, []byte(nil), resp)
	assert.Equal(t, "VM execute fail: abort", err.Error())

	// vm.Call bob transferFrom 400 which alice approve to bob
	resp, _, err = vm.Call(bobRef, addr, data, gasLimit, big.NewInt(0))
	expectResp = utils.Align32BytesConverter(resp, "bool")
	assert.Equal(t, true, expectResp)
	assert.NoError(t, err)

	// vm.Call charlie getApproveBalance
	resp, _, err = vm.Call(charlieRef, addr, getApproveData, gasLimit, big.NewInt(0))
	expectResp = utils.Align32BytesConverter(resp, "uint64")
	assert.Equal(t, uint64(0), expectResp)
	assert.NoError(t, err)

	// vm.Call charlie get alice balance
	funcName = []byte("getBalance")
	aliceBalanceData, err := rlp.EncodeToBytes([]interface{}{funcName, from})
	assert.NoError(t, err)
	resp, _, err = vm.Call(charlieRef, addr, aliceBalanceData, gasLimit, big.NewInt(0))
	expectResp = utils.Align32BytesConverter(resp, "uint64")
	assert.Equal(t, uint64(100), expectResp)
	assert.NoError(t, err)

	// vm.Call charlie get bob balance
	bobBalanceData, err := rlp.EncodeToBytes([]interface{}{funcName, spender})
	assert.NoError(t, err)
	resp, _, err = vm.Call(charlieRef, addr, bobBalanceData, gasLimit, big.NewInt(0))
	expectResp = utils.Align32BytesConverter(resp, "uint64")
	assert.Equal(t, uint64(400), expectResp)
	assert.NoError(t, err)

	// vm.Call charlie burn the token (not enough token)
	funcName = []byte("burn")
	data, err = rlp.EncodeToBytes([]interface{}{funcName, amount2})
	assert.NoError(t, err)
	resp, _, err = vm.Call(charlieRef, addr, data, gasLimit, big.NewInt(0))
	assert.Equal(t, []byte(nil), resp)
	assert.Equal(t, "VM execute fail: abort", err.Error())

	// vm.Call alice burn the token(failed)
	resp, _, err = vm.Call(aliceRef, addr, data, gasLimit, big.NewInt(0))
	assert.Equal(t, []byte(nil), resp)
	assert.Equal(t, "VM execute fail: abort", err.Error())

	// vm.Call bob burn the token
	resp, _, err = vm.Call(bobRef, addr, data, gasLimit, big.NewInt(0))
	expectResp = utils.Align32BytesConverter(resp, "bool")
	assert.Equal(t, true, expectResp)
	assert.NoError(t, err)

	// vm.Call charlie get alice balance
	resp, _, err = vm.Call(charlieRef, addr, aliceBalanceData, gasLimit, big.NewInt(0))
	expectResp = utils.Align32BytesConverter(resp, "uint64")
	assert.Equal(t, uint64(500), expectResp)
	assert.NoError(t, err)

	// vm.Call charlie get bob balance
	resp, _, err = vm.Call(charlieRef, addr, bobBalanceData, gasLimit, big.NewInt(0))
	expectResp = utils.Align32BytesConverter(resp, "uint64")
	assert.Equal(t, uint64(0), expectResp)
	assert.NoError(t, err)
}

func TestVM_Call_Error(t *testing.T) {
	vm := getTestVm()
	ref := AccountRef(aliceAddr)
	gasLimit := g_testData.TestGasLimit
	value := g_testData.TestValue

	vm.GetStateDB().CreateAccount(ref.Address())
	vm.GetStateDB().AddBalance(ref.Address(), big.NewInt(10000))

	vm.GetStateDB().CreateAccount(contractAddr)
	vm.GetStateDB().SetCode(contractAddr, []byte{123})
	vm.GetStateDB().SetAbi(contractAddr, []byte{123})
	_, _, err := vm.Call(ref, contractAddr, nil, gasLimit, value)
	assert.Equal(t, "unexpected EOF", err.Error())

	vm = getTestVm()
	vm.GetStateDB().CreateAccount(ref.Address())
	_, _, err = vm.Call(ref, contractAddr, nil, gasLimit, value)
	assert.Equal(t, g_error.ErrInsufficientBalance, err)

	vm.GetStateDB().AddBalance(ref.Address(), big.NewInt(10000))
	_, _, err = vm.Call(ref, contractAddr, nil, gasLimit, value)
	assert.NoError(t, err)

	vm.depth = int(model2.CallCreateDepth + 1)
	_, _, err = vm.Call(ref, contractAddr, nil, gasLimit, value)
	assert.Equal(t, g_error.ErrDepth, err)

	vm.depth = 1
	vm.vmConfig.NoRecursion = true
	_, _, err = vm.Call(ref, contractAddr, nil, gasLimit, value)
	assert.NoError(t, err)
}

func TestVM_DelegateCall_Error(t *testing.T) {
	vm := getTestVm()
	caller := AccountRef(aliceAddr)
	self := AccountRef(contractAddr)
	value := g_testData.TestValue
	gasLimit := g_testData.TestGasLimit
	contract := NewContract(caller, self, value, gasLimit, nil)

	vm.GetStateDB().CreateAccount(contractAddr)
	vm.GetStateDB().SetCode(contractAddr, []byte{123})
	vm.GetStateDB().SetAbi(contractAddr, []byte{123})

	_, _, err := vm.DelegateCall(contract, contractAddr, nil, gasLimit)
	assert.Equal(t, "unexpected EOF", err.Error())

	vm.depth = int(model2.CallCreateDepth + 1)
	_, _, err = vm.DelegateCall(contract, contractAddr, nil, gasLimit)
	assert.Equal(t, g_error.ErrDepth, err)

	vm.depth = 1
	vm.vmConfig.NoRecursion = true
	_, _, err = vm.DelegateCall(contract, contractAddr, nil, gasLimit)
	assert.NoError(t, err)
}

func TestVM_Create_Error(t *testing.T) {
	vm := getTestVm()
	caller := AccountRef(aliceAddr)
	value := g_testData.TestValue
	gasLimit := g_testData.TestGasLimit

	_, _, _, err := vm.Create(caller, nil, gasLimit, value)
	assert.Equal(t, "empty account", err.Error())

	vm.GetStateDB().CreateAccount(caller.Address())
	_, _, _, err = vm.Create(caller, nil, gasLimit, value)
	assert.Equal(t, g_error.ErrInsufficientBalance, err)

	vm.GetStateDB().AddBalance(caller.Address(), big.NewInt(10000))
	_, _, _, err = vm.Create(caller, nil, gasLimit, value)
	assert.Equal(t, errEmptyInput, err)

	input, _ := rlp.EncodeToBytes([]interface{}{"code", "abi"})
	_, _, _, err = vm.Create(caller, input, gasLimit, value)
	assert.Equal(t, "wasm: Invalid magic number", err.Error())

	vm.depth = 1
	vm.vmConfig.NoRecursion = true
	_, _, _, err = vm.Create(caller, input, gasLimit, value)
	assert.NoError(t, err)

	nonce, _ := vm.GetStateDB().GetNonce(caller.Address())
	expectAddr := cs_crypto.CreateContractAddress(caller.Address(), nonce)
	vm.GetStateDB().CreateAccount(expectAddr)
	vm.GetStateDB().AddNonce(expectAddr, uint64(1))
	_, _, _, err = vm.Create(caller, input, gasLimit, value)
	assert.Equal(t, g_error.ErrContractAddressCollision, err)

	vm.depth = int(model2.CallCreateDepth + 1)
	_, _, _, err = vm.Create(caller, nil, gasLimit, value)
	assert.Equal(t, g_error.ErrDepth, err)
}*/
