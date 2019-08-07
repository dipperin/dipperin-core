package vm

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/vm/model"
	"math/big"
)

type StateDB interface {
	GetBalance(common.Address) *big.Int
	CreateAccount(common.Address)
	GetNonce(common.Address) (uint64, error)
	AddNonce(common.Address, uint64)
	GetCodeHash(common.Address) common.Hash
	GetCode(common.Address) []byte
	SetCode(common.Address, []byte)
	GetAbiHash(common.Address) common.Hash
	GetAbi(common.Address) []byte
	SetAbi(common.Address, []byte)
	AddBalance(addr common.Address, amount *big.Int)
	SubBalance(addr common.Address, amount *big.Int)

	/*AddRefund(uint64)
	SubRefund(uint64)
	GetRefund() uint64*/

	// todo: hash -> bytes
	GetState(common.Address, []byte) []byte
	SetState(common.Address, []byte, []byte)
	AddLog(addedLog *model.Log)
	GetLogs(txHash common.Hash) []*model.Log

	/*	Suicide(common.Address) bool
		HasSuicided(common.Address) bool*/

	// Exist reports whether the given account exists in state.
	// Notably this should also return true for suicided accounts.
	Exist(common.Address) bool
	RevertToSnapshot(int)
	Snapshot() int
	/*AddPreimage(common.Hash, []byte)
	ForEachStorage(common.Address, func(common.Hash, common.Hash) bool)*/

	//ppos add
	/*	TxHash() common.Hash
		TxIdx() uint32*/
}

/*type AccountDB interface {
	GetBalance(addr common.Address) (*big.Int, error)
	AddBalance(addr common.Address, amount *big.Int) error
	SubBalance(addr common.Address, amount *big.Int) error
}

type ContractDB interface {
	ContractExist(addr common.Address) bool
	GetContract(addr common.Address, vType reflect.Type) (v reflect.Value, err error)
	FinalizeContract(addr common.Address, data reflect.Value) error
}
*/
