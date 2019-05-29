package state_processor

import (
	"encoding/json"
	vm2 "github.com/dipperin/dipperin-core/common/vmcommon"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/core/vm"
	model2 "github.com/dipperin/dipperin-core/core/vm/model"
)

type CallCode struct {
	Func  []byte
	Input []byte `json:"Input"`
}


func (state *AccountStateDB) ProcessContract(tx model.AbstractTransaction, block model.AbstractBlock, create bool) (err error) {
	context := vm.NewVMContext(tx, block)
	fullState := &Fullstate{
		state: state,
	}
	dvm := vm.NewVM(context, fullState, vm.DEFAULT_VM_CONFIG)
	err = dvm.PreCheck()
	if create {
		data := tx.ExtraData()
		var ca *vm2.CodeAbi
		err := json.Unmarshal(data,ca)
		if err!= nil{
			return err
		}
		_, _, _, err = dvm.Create(vm.AccountRef(context.Origin), ca.Code, 0, nil)
		if err != nil {
			return err
		}
	} else {
		data := tx.ExtraData()
		_, _, err = dvm.Call(vm.AccountRef(context.Origin), *tx.To(), data, 0, tx.Amount())
		if err != nil {
			return err
		}
	}

	root, err := state.Finalise()
	if err != nil {
		return err
	}
	model2.NewReceipt(root.Bytes(), false, 0)

	return nil
}
