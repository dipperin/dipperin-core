package stateprocessor

import (
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/log"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third_party/crypto/cs-crypto"
	"go.uber.org/zap"
	"math/big"
)

type Fullstate struct {
	state *AccountStateDB
}

func NewFullState(state *AccountStateDB) *Fullstate {
	return &Fullstate{
		state: state,
	}
}

func (f *Fullstate) CreateAccount(address common.Address) error {
	_, err := f.state.newContractAccount(address)
	return err
}

func (f *Fullstate) GetBalance(addr common.Address) *big.Int {
	balance, err := f.state.GetBalance(addr)
	if err != nil {
		panic(fmt.Sprintf("GetBalance failed, err=%v", err))
	}
	return balance
}

func (f *Fullstate) GetNonce(addr common.Address) (uint64, error) {
	return f.state.GetNonce(addr)
}

func (f *Fullstate) AddNonce(addr common.Address, add uint64) {
	err := f.state.AddNonce(addr, add)
	if err != nil {
		log.DLogger.Error("Fullstate#AddBalance", zap.Error(err))
		//panic(fmt.Sprintf("AddNonce failed, err=%v", err))
	}
}

func (f *Fullstate) AddBalance(addr common.Address, amount *big.Int) error {
	err := f.state.AddBalance(addr, amount)
	if err != nil {
		log.DLogger.Error("Fullstate#AddBalance", zap.Error(err))
		//panic(fmt.Sprintf("AddBalance failed, err=%v", err))
	}
	return err
}
func (f *Fullstate) SubBalance(addr common.Address, amount *big.Int) error {
	err := f.state.SubBalance(addr, amount)
	if err != nil {
		log.DLogger.Error("Fullstate#AddBalance", zap.Error(err))
		//panic(fmt.Sprintf("SubBalance failed, err=%v", err))
	}
	return err
}

func (f *Fullstate) GetCodeHash(addr common.Address) common.Hash {
	code, err := f.state.GetCode(addr)
	if err != nil {
		return common.Hash{}
	}
	return cs_crypto.Keccak256Hash(code)
}

func (f *Fullstate) GetCode(addr common.Address) []byte {
	//f.state.contractTrieCache
	code, err := f.state.GetCode(addr)
	if err != nil {
		return nil
	}
	return code
}

func (f *Fullstate) SetCode(addr common.Address, code []byte) {
	err := f.state.SetCode(addr, code)
	if err != nil {
		panic(fmt.Sprintf("SetCode failed, err=%v", err))
	}
}

func (f *Fullstate) GetAbiHash(addr common.Address) common.Hash {
	abi, err := f.state.GetAbi(addr)
	if err != nil {
		return common.Hash{}
	}
	return cs_crypto.Keccak256Hash(abi)
}

func (f *Fullstate) GetAbi(addr common.Address) []byte {
	abi, err := f.state.GetAbi(addr)
	if err != nil {
		return nil
	}
	return abi
}

func (f *Fullstate) SetAbi(addr common.Address, abi []byte) {
	err := f.state.SetAbi(addr, abi)
	if err != nil {
		panic(fmt.Sprintf("SetAbi failed, err=%v", err))
	}
}

func (f *Fullstate) GetState(addr common.Address, key []byte) (data []byte) {
	return f.state.GetData(addr, string(key))
}

func (f *Fullstate) SetState(addr common.Address, key []byte, value []byte) {
	f.state.SetData(addr, string(key), value)
}

func (f *Fullstate) AddLog(addedLog *model.Log) {
	err := f.state.AddLog(addedLog)
	if err != nil {
		panic(fmt.Sprintf("SetAbi failed, err=%v", err))
	}
}

func (f *Fullstate) GetLogs(txHash common.Hash) []*model.Log {
	return f.state.GetLogs(txHash)
}

/*func (f *Fullstate) Suicide(common.Address) bool {
	panic("implement me")
}

func (f *Fullstate) HasSuicided(common.Address) bool {
	panic("implement me")
}*/

func (f *Fullstate) Exist(addr common.Address) bool {
	return !f.state.IsEmptyAccount(addr)
}

func (f *Fullstate) RevertToSnapshot(id int) {
	log.DLogger.Debug("State Reverted", zap.Int("id", id))
	f.state.RevertToSnapshot(id)
}

func (f *Fullstate) Snapshot() int {
	return f.state.Snapshot()
}
