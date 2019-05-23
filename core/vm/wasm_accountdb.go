package vm

import (
	"math/big"
	"github.com/dipperin/dipperin-core/common"
)

type WasmAccountDB struct {

}

func (db WasmAccountDB) GetBalance(addr common.Address) (*big.Int, error) {
	panic("implement me")
}

func (db WasmAccountDB) AddBalance(addr common.Address, amount *big.Int) error {
	panic("implement me")
}

func (db WasmAccountDB) SubBalance(addr common.Address, amount *big.Int) error {
	panic("implement me")
}


