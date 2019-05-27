package state_processor

import (
	"encoding/json"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/core/vm"
	model2 "github.com/dipperin/dipperin-core/core/vm/model"
)
type CodeAbi struct {
	Code   []byte `json:"code"`
	Abi []byte `json:"abi"`
	Input []byte `json:"Input"`
}

type CallCode struct {
	Func []byte
	Input []byte `json:"Input"`
}

func (state *AccountStateDB) ProcessContract(tx model.AbstractTransaction, blockHeight uint64, create bool) (err error) {
	context := vm.NewVMContext(tx)
	fullState := &Fullstate{
		state,

	}

	newVm := vm.NewVM(context, fullState, vm.DEFAULT_VM_CONFIG)
	var gas uint64
	if create{
		data := tx.ExtraData()
		var ca *CodeAbi
		err := json.Unmarshal(data,ca)
		if err!= nil{
			return err
		}
		_, _,gas,err = newVm.Create(&vm.Caller{context.Origin},ca.Code,ca.Abi,ca.Input)
		if err != nil {
			return err
		}
	}else{
		data := tx.ExtraData()
		_, gas,err = newVm.Call(&vm.Caller{context.Origin},tx.To(),data)
		if err != nil {
			return err
		}
	}

	root, err := state.Finalise()
	if err != nil {
		return err
	}

	receipt := model2.NewReceipt(root.Bytes(), false, gas)

    return nil
}