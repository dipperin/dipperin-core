package state_processor

import (
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/core/vm"
)

type CallCode struct {
	Func  []byte
	Input []byte `json:"Input"`
}

func (state *AccountStateDB) ProcessContract(tx model.AbstractTransaction, block model.AbstractBlock, create bool) (model.ReceiptPara, error) {
	gp := model.DefaultGasLimit
	context := vm.NewVMContext(tx, block)
	fullState := &Fullstate{
		state: state,
	}
	msg, err := tx.AsMessage()
	if err != nil {
		return model.ReceiptPara{}, err
	}
	dvm := vm.NewVM(context, fullState, vm.DEFAULT_VM_CONFIG)
	_, usedGas, failed, err := ApplyMessage(dvm, msg, &gp)
	if err != nil {
		return model.ReceiptPara{}, err
	}

	root, err := state.Finalise()
	if err != nil {
		return model.ReceiptPara{}, err
	}
	return model.ReceiptPara{
		Root:          root[:],
		HandlerResult: failed,
		//todo CumulativeGasUsed暂时使用usedGas,不考虑在apply交易前已有gas使用的情景
		CumulativeGasUsed: usedGas,
		GasUsed:           usedGas,
		Logs:              fullState.GetLogs(tx.CalTxId()),
	}, nil
}
