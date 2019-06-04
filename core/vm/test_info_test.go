package vm

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/third-party/log"
	"math/big"
	"github.com/dipperin/dipperin-core/core/vm/model"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
	"github.com/ethereum/go-ethereum/rlp"
)

var contractAddr = common.HexToAddress("0x00005586B883Ec6dd4f8c26063E18eb4Bd228e59c3E9")

type fakeContractRef struct {
	addr common.Address
}

func (ref fakeContractRef) Address() common.Address {
	return ref.addr
}

type fakeStateDB struct {
}

func (state fakeStateDB) GetLogs(txHash common.Hash) []*model.Log {
	panic("implement me")
}

func (state fakeStateDB) AddLog(addedLog *model.Log) {
	log.Info("add log success")
	return
}

func (state fakeStateDB) CreateAccount(common.Address) {
	panic("implement me")
}

func (state fakeStateDB) SubBalance(common.Address, *big.Int) {
	panic("implement me")
}

func (state fakeStateDB) AddBalance(common.Address, *big.Int) {
	panic("implement me")
}

func (state fakeStateDB) GetBalance(common.Address) *big.Int {
	panic("implement me")
}

func (state fakeStateDB) GetNonce(common.Address) uint64 {
	panic("implement me")
}

func (state fakeStateDB) SetNonce(common.Address, uint64) {
	panic("implement me")
}

func (state fakeStateDB) AddNonce(common.Address, uint64) {
	panic("implement me")
}

func (state fakeStateDB) GetCodeHash(common.Address) common.Hash {
	panic("implement me")
}

func (state fakeStateDB) GetCode(common.Address) []byte {
	panic("implement me")
}

func (state fakeStateDB) SetCode(common.Address, []byte) {
	panic("implement me")
}

func (state fakeStateDB) GetCodeSize(common.Address) int {
	panic("implement me")
}

func (state fakeStateDB) GetAbiHash(common.Address) common.Hash {
	panic("implement me")
}

func (state fakeStateDB) GetAbi(common.Address) []byte {
	panic("implement me")
}

func (state fakeStateDB) SetAbi(common.Address, []byte) {
	panic("implement me")
}

func (state fakeStateDB) AddRefund(uint64) {
	panic("implement me")
}

func (state fakeStateDB) SubRefund(uint64) {
	panic("implement me")
}

func (state fakeStateDB) GetRefund() uint64 {
	panic("implement me")
}

func (state fakeStateDB) GetCommittedState(common.Address, []byte) []byte {
	panic("implement me")
}

func (state fakeStateDB) GetState(common.Address, []byte) []byte {
	fmt.Println("fake stateDB get state sucessful")
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.LittleEndian, int32(123))
	return bytesBuffer.Bytes()
}

func (state fakeStateDB) SetState(common.Address, []byte, []byte) {
	fmt.Println("fake stateDB set state sucessful")
}

func (state fakeStateDB) Suicide(common.Address) bool {
	panic("implement me")
}

func (state fakeStateDB) HasSuicided(common.Address) bool {
	panic("implement me")
}

func (state fakeStateDB) Exist(common.Address) bool {
	panic("implement me")
}

func (state fakeStateDB) Empty(common.Address) bool {
	panic("implement me")
}

func (state fakeStateDB) RevertToSnapshot(int) {
	panic("implement me")
}

func (state fakeStateDB) Snapshot() int {
	panic("implement me")
}

func (state fakeStateDB) AddPreimage(common.Hash, []byte) {
	panic("implement me")
}

func (state fakeStateDB) ForEachStorage(common.Address, func(common.Hash, common.Hash) bool) {
	panic("implement me")
}

func (state fakeStateDB) TxHash() common.Hash {
	return common.Hash{}
}

func (state fakeStateDB) TxIdx() uint32 {
	return 0
}

func genInput(t *testing.T, funcName string, param [][]byte) []byte {
	var input [][]byte
	input = make([][]byte, 0)
	// func name
	input = append(input, []byte(funcName))
	// func parameter
	for _, v := range (param) {
		input = append(input, v)
	}

	buffer := new(bytes.Buffer)
	err := rlp.Encode(buffer, input)
	assert.NoError(t, err)
	return buffer.Bytes()
}

func getCodeWithABI(t *testing.T, code, abi []byte) []byte {
	var input [][]byte
	input = make([][]byte, 0)
	// code
	input = append(input, code)
	// abi
	input = append(input, abi)

	buffer := new(bytes.Buffer)
	err := rlp.Encode(buffer, input)
	assert.NoError(t, err)
	return buffer.Bytes()
}

func getContract(t *testing.T, addr common.Address, code, abi string) *Contract {
	fileCode, err := ioutil.ReadFile(code)
	assert.NoError(t, err)

	fileABI, err := ioutil.ReadFile(abi)
	assert.NoError(t, err)

	ca := getCodeWithABI(t, fileCode, fileABI)
	return &Contract{
		self: fakeContractRef{addr: addr},
		Code: ca,
		Gas:  model.TxGas,
	}
}

func getTestVm() *VM {
	return NewVM(Context{
		BlockNumber: big.NewInt(1),
		GasLimit:    model.TxGas,
		GetHash:     getTestHashFunc(),
	}, fakeStateDB{}, DEFAULT_VM_CONFIG)
}

func getTestHashFunc() func(num uint64) common.Hash {
	return func(num uint64) common.Hash {
		return common.Hash{}
	}
}
