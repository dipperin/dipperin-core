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
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/core/vm/base"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

/*

const (
	abi1 = `[{
        "name": "init",
        "inputs": [],
        "outputs": [],
        "constant": "false",
        "type": "function"
    }]`

	abi2 = `[{
        "name": "init",
        "inputs": [
            {
                "name": "inputName",
                "type": "uint64"
            }
        ],
        "outputs": [
            {
                "name": "outputName",
                "type": "string"
            }
        ],
        "constant": "false",
        "type": "function"
    }]`

	abi3 = `[{
        "name": "init",
        "inputs": [
            {
                "name": "inputName",
                "type": "uint32"
            }
        ],
        "outputs": [],
        "constant": "false",
        "type": "function"
    }]`

	abi4 = `[{
        "name": "init",
        "inputs": [],
        "outputs": [
            {
                "name": "outputName",
                "type": "string"
            }
        ],
        "constant": "false",
        "type": "function"
    }]`

	abi5 = `[{
        "name": "init",
        "inputs": [
            {
                "name": "inputName",
                "type": "uint64"
            },
			{
                "name": "inputName",
                "type": "uint32"
            },
			{
                "name": "inputName",
                "type": "uint16"
            },
			{
                "name": "inputName",
                "type": "uint8"
            },
			{
                "name": "inputName",
                "type": "bool"
            }
        ],
        "outputs": [],
        "constant": "false",
        "type": "function"
    }]`
)

type fakeStateDB struct {
	balanceMap map[common.Address]*big.Int
	nonceMap   map[common.Address]uint64
	codeMap    map[common.Address][]byte
	abiMap     map[common.Address][]byte
	stateMap   map[common.Address]map[string][]byte
}

func NewFakeStateDB() *fakeStateDB {
	return &fakeStateDB{
		balanceMap: make(map[common.Address]*big.Int, 0),
		nonceMap:   make(map[common.Address]uint64, 0),
		codeMap:    make(map[common.Address][]byte, 0),
		abiMap:     make(map[common.Address][]byte, 0),
		stateMap:   make(map[common.Address]map[string][]byte, 0),
	}
}

func (state *fakeStateDB) GetLogs(txHash common.Hash) []*model.Log {
	panic("implement me")
}

func (state *fakeStateDB) AddLog(addedLog *model.Log) {
	log.DLogger.Info("add log success")
	return
}

func (state *fakeStateDB) CreateAccount(addr common.Address) error {
	state.balanceMap[addr] = big.NewInt(0)
	state.nonceMap[addr] = uint64(0)
	state.codeMap[addr] = []byte{}
	state.abiMap[addr] = []byte{}
	state.stateMap[addr] = make(map[string][]byte)
	return nil
}

func (state *fakeStateDB) SubBalance(addr common.Address, amount *big.Int) error {
	balance := state.balanceMap[addr]
	state.balanceMap[addr] = big.NewInt(0).Sub(balance, amount)
	return nil
}

func (state *fakeStateDB) AddBalance(addr common.Address, amount *big.Int) error {
	balance := state.balanceMap[addr]
	state.balanceMap[addr] = big.NewInt(0).Add(balance, amount)
	return nil
}

func (state *fakeStateDB) GetBalance(addr common.Address) *big.Int {
	return state.balanceMap[addr]
}

func (state *fakeStateDB) GetNonce(addr common.Address) (uint64, error) {
	if nonce, ok := state.nonceMap[addr]; !ok {
		return uint64(0), errors.New("empty account")
	} else {
		return nonce, nil
	}
}

func (state *fakeStateDB) AddNonce(addr common.Address, nonce uint64) {
	curNonce, _ := state.GetNonce(addr)
	state.nonceMap[addr] = curNonce + nonce
}

func (state *fakeStateDB) GetCodeHash(addr common.Address) common.Hash {
	code := state.GetCode(addr)
	return cs_crypto.Keccak256Hash(code)
}

func (state *fakeStateDB) GetCode(addr common.Address) []byte {
	return state.codeMap[addr]
}

func (state *fakeStateDB) SetCode(addr common.Address, code []byte) {
	if state.codeMap == nil {
		state.codeMap = make(map[common.Address][]byte)
	}
	state.codeMap[addr] = code
}

func (state *fakeStateDB) GetAbiHash(addr common.Address) common.Hash {
	abi := state.GetCode(addr)
	return cs_crypto.Keccak256Hash(abi)
}

func (state *fakeStateDB) GetAbi(addr common.Address) []byte {
	return state.abiMap[addr]
}

func (state *fakeStateDB) SetAbi(addr common.Address, abi []byte) {
	state.abiMap[addr] = abi
}

func (state *fakeStateDB) AddRefund(uint64) {
	panic("implement me")
}

func (state *fakeStateDB) SubRefund(uint64) {
	panic("implement me")
}

func (state *fakeStateDB) GetRefund() uint64 {
	panic("implement me")
}

func (state *fakeStateDB) GetState(addr common.Address, key []byte) []byte {
	return state.stateMap[addr][string(key)]
}

func (state *fakeStateDB) SetState(addr common.Address, key []byte, value []byte) {
	state.stateMap[addr][string(key)] = value
}

func (state *fakeStateDB) Exist(addr common.Address) bool {
	if _, ok := state.nonceMap[addr]; ok {
		return true
	} else {
		return false
	}
}

func (state *fakeStateDB) RevertToSnapshot(int) {
	return
}

func (state *fakeStateDB) Snapshot() int {
	return 0
}

func (state *fakeStateDB) AddPreimage(common.Hash, []byte) {
	panic("implement me")
}

func (state *fakeStateDB) ForEachStorage(common.Address, func(common.Hash, common.Hash) bool) {
	panic("implement me")
}

func (state *fakeStateDB) TxHash() common.Hash {
	return common.Hash{}
}

func (state *fakeStateDB) TxIdx() uint32 {
	return 0
}

func genInput(t *testing.T, funcName string, param [][]byte) []byte {
	input := make([][]byte, 0)

	// func name
	if funcName != "" {
		input = append(input, []byte(funcName))
	}

	// func parameter
	for _, v := range param {
		input = append(input, v)
	}

	result, err := rlp.EncodeToBytes(input)
	assert.NoError(t, err)
	return result
}






*/

func genInput(t *testing.T, funcName string, param [][]byte) []byte {
	input := make([][]byte, 0)

	// func name
	if funcName != "" {
		input = append(input, []byte(funcName))
	}

	// func parameter
	for _, v := range param {
		input = append(input, v)
	}

	result, err := rlp.EncodeToBytes(input)
	assert.NoError(t, err)
	return result
}

func GetBaseVmInfo(t *testing.T) (*gomock.Controller, *base.MockStateDB, *VM) {
	ctrl := gomock.NewController(t)

	db := base.NewMockStateDB(ctrl)
	vm := NewVM(base.Context{
		Origin:      model.AliceAddr,
		BlockNumber: big.NewInt(1),
		CanTransfer: base.CanTransfer,
		Transfer:    base.Transfer,
		GetHash:     getTestHashFunc(),
	}, db, base.DEFAULT_VM_CONFIG)
	return ctrl, db, vm
}

func getContract(code, abi []byte, input []byte) *Contract {
	//fileCode, fileABI := g_testData.GetCodeAbi(code, abi)
	caller := base.AccountRef(model.AliceAddr)
	self := base.AccountRef(model.ContractAddr)
	value := model.TestZeroValue
	gasLimit := model.TestGasLimit
	contract := NewContract(caller, self, value, gasLimit, input)
	contract.SetCode(&model.AliceAddr, common.Hash{}, code)
	contract.SetAbi(&model.AliceAddr, common.Hash{}, abi)
	return contract
}

/*func getTestVm() *VM {
	return NewVM(Context{
		Origin:      model.AliceAddr,
		BlockNumber: big.NewInt(1),
		CanTransfer: CanTransfer,
		Transfer:    Transfer,
		GetHash:     getTestHashFunc(),
	}, NewFakeStateDB(), DEFAULT_VM_CONFIG)
}*/

func getTestHashFunc() func(num uint64) common.Hash {
	return func(num uint64) common.Hash {
		return common.Hash{}
	}
}
