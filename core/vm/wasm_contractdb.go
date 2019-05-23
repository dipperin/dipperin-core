package vm

import (
	"reflect"
	"github.com/dipperin/dipperin-core/common"
)

type WasmContractDB struct {

}

func (db WasmContractDB) ContractExist(addr common.Address) bool {
	panic("implement me")
}

func (db WasmContractDB) GetContract(addr common.Address, vType reflect.Type) (v reflect.Value, err error) {
	panic("implement me")
}

func (db WasmContractDB) FinalizeContract(addr common.Address, data reflect.Value) error {
	panic("implement me")
}