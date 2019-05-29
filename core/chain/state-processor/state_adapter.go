package state_processor

import (
	"github.com/dipperin/dipperin-core/common"
	model2 "github.com/dipperin/dipperin-core/core/vm/model"
	"github.com/dipperin/dipperin-core/third-party/log"
	"math/big"
)

type Fullstate struct {
	state *AccountStateDB
}

func (f *Fullstate) CreateAccount(address common.Address) {
	f.state.newContractAccount(address)
}

func (f *Fullstate) GetBalance(addr common.Address) *big.Int {
	balance, err := f.state.GetBalance(addr)
	if err != nil {
		return big.NewInt(0)
	}
	return balance
}

func (f *Fullstate) GetNonce(addr common.Address) uint64 {
	nonce, err := f.state.GetNonce(addr)
	if err != nil {
		return uint64(0)
	}
	return nonce
}

func (f *Fullstate) AddNonce(addr common.Address, add uint64) {
	err := f.state.AddNonce(addr, add)
	if err != nil {
		panic("add nonce error")
	}
}

func (f *Fullstate) AddBalance(addr common.Address, amount *big.Int){
	err := f.state.AddBalance(addr, amount)
	if err != nil{
		panic("Can not add balance")
	}
}
func (f *Fullstate) SubBalance(addr common.Address, amount *big.Int){
	err := f.state.SubBalance(addr,amount)
	if err != nil{
		panic("Can not sub balance")
	}
}

func (f  *Fullstate) GetCodeHash(addr common.Address) common.Hash {
	code, err := f.state.GetCode(addr)
	if err!=nil{
		return common.Hash{}
	}
	return common.RlpHashKeccak256(code)
}

func (f *Fullstate) GetCode(addr common.Address) (result []byte) {
	//f.state.contractTrieCache
	code, err := f.state.GetCode(addr)
	if err!=nil{
		return
	}
	return code
}

func (f  *Fullstate) SetCode(addr common.Address, code []byte) {
	err := f.state.SetCode(addr,code)
	if err!=nil{
		panic("set code error")
	}
}

func (f  *Fullstate) GetCodeSize(addr common.Address) (size int) {
	code, err := f.state.GetCode(addr)
	if err!=nil{
		return
	}
	return len(code)
}

func (f  *Fullstate) GetAbiHash(addr common.Address) common.Hash {
	abi,err := f.state.GetAbi(addr)
	if err != nil{
		return common.Hash{}
	}
	return common.RlpHashKeccak256(abi)
}

func (f  *Fullstate) GetAbi(addr common.Address) (abi []byte) {
	abi,err := f.state.GetAbi(addr)
	if err != nil{
		return
	}
	return abi
}

func (f  *Fullstate) SetAbi(addr common.Address, abi []byte) {
	err := f.state.SetAbi(addr,abi)
	if err != nil{
		panic("set abi error")
	}
}

func (f *Fullstate) GetCommittedState(common.Address, []byte) []byte {
	panic("implement me")
}

func (f *Fullstate) GetState(addr common.Address, key []byte) (data []byte) {
	ct, err := f.state.getContractTrie(addr)
	if err != nil {
		return
	}
	return ct.GetKey(GetContractFieldKey(addr, string(key)))
}

func (f *Fullstate) SetState(addr common.Address, key []byte, value []byte) {
	ct, err := f.state.getContractTrie(addr)
	if err != nil {
		return
	}
	err = ct.TryUpdate(GetContractFieldKey(addr, string(key)), value)
	if err != nil {
		panic("can not update contract field")
	}
}

func (f *Fullstate) AddLog(addedLog *model2.Log) {
	log.Info("AddLog Called")

	txHash := addedLog.TxHash
	contractLogs := f.GetLogs(txHash)
	addedLog.Index = uint(len(contractLogs) + 1)
	f.state.logs[txHash] = append(contractLogs, addedLog)

	log.Info("Log Added", "txHash", txHash, "logs", f.state.logs[txHash])
}

func (f *Fullstate) GetLogs(txHash common.Hash) []*model2.Log {
	return f.state.logs[txHash]
}
func (f *Fullstate) Suicide(common.Address) bool {
	panic("implement me")
}

func (f *Fullstate) HasSuicided(common.Address) bool {
	panic("implement me")
}

func (f *Fullstate) Exist(common.Address) bool {
	panic("implement me")
}

func (f *Fullstate) Empty(common.Address) bool {
	panic("implement me")
}

func (f *Fullstate) RevertToSnapshot(int) {
	panic("implement me")
}

func (f *Fullstate) Snapshot() int {
	panic("implement me")
}

func (f *Fullstate) TxHash() common.Hash {
	panic("implement me")
}

func (f *Fullstate) TxIdx() uint32 {
	panic("implement me")
}
