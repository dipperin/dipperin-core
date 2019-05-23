package state_processor

import (
	"github.com/dipperin/dipperin-core/core/vm"
	"github.com/dipperin/dipperin-core/core/model"
)

//Move these const to common type.go
const (
	AddressTypeContractCreate    = 0x0012
	AddressTypeContract = 0x0013
)

type AccountStateDB struct{

}

func (state *AccountStateDB) ProcessTx (tx model.AbstractTransaction) (err error){

	err = state.processBasicTx(tx)
	if err != nil {
		return
	}

	switch tx.GetType() {
	case AddressTypeContractCreate:
		err = state.processContractCreate(tx)
	case AddressTypeContract:
		err = state.processContractCall(tx)
	}
	return
}

func (state *AccountStateDB) processContractCreate (tx model.AbstractTransaction) (err error){
	context := vm.NewVMContext(tx)
	vm := vm.NewVM(context, nil, vm.DEFAULT_VM_CONFIG)

	// apply msg
	vm.Create(&vm.Caller{vm.Origin}, tx.ExtraData(),tx.Amount())

	// modify db

	// prepare receipt
	return
}

func (state *AccountStateDB) processContractCall (tx model.AbstractTransaction) (err error){
	return
}

func (state *AccountStateDB) processBasicTx(tx model.AbstractTransaction) (err error) {
	// Add nonce, sub balance.
	return
}

