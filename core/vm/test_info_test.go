package vm

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"math/big"
)

var contractAddr = common.HexToAddress("0x00005586B883Ec6dd4f8c26063E18eb4Bd228e59c3E9")

type fakeStateDB struct {

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
	panic("implement me")
}

func (state fakeStateDB) TxIdx() uint32 {
	panic("implement me")
}
