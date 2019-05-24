package state_processor

import (
	"encoding/json"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/core/vm"
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
	vm := vm.NewVM(context, fullState, vm.DEFAULT_VM_CONFIG)
	if create{
		data := tx.ExtraData()
		var ca *CodeAbi
		err := json.Unmarshal(data,ca)
		if err!= nil{
			return err
		}
		_, _,_,err = vm.Create(&vm.Caller{context.Origin},ca.Code,ca.Abi,ca.Input)
		if err != nil {
			return err
		}
	}else{
		data := tx.ExtraData()
		_, _,err = vm.Call(&vm.Caller{context.Origin},tx.To(),data)
		if err != nil {
			return err
		}
	}


    return nil
}