package state_processor

import (
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/core/vm"
	"github.com/dipperin/dipperin-core/third-party/log"
)

type CallCode struct {
	Func  []byte
	Input []byte `json:"Input"`
}

func (state *AccountStateDB) ProcessContract(tx model.AbstractTransaction, block model.AbstractBlock, blockGasLimit *uint64, create bool) (model.ReceiptPara, error) {
	context := vm.NewVMContext(tx, block)
	fullState := &Fullstate{
		state: state,
	}
	msg, err := tx.AsMessage()
	if err != nil {
		return model.ReceiptPara{}, err
	}
	dvm := vm.NewVM(context, fullState, vm.DEFAULT_VM_CONFIG)
	_, usedGas, failed, err := ApplyMessage(dvm, msg, blockGasLimit)
	if err != nil {
		log.Error("AccountStateDB#ProcessContract", "ApplyMessage err", err)
		return model.ReceiptPara{},err
	}

	root, err := state.Finalise()
	if err != nil {
		return model.ReceiptPara{}, err
	}
	return model.ReceiptPara{
		Root:          root[:],
		HandlerResult: !failed,
		//todo CumulativeGasUsed暂时使用usedGas,不考虑在apply交易前已有gas使用的情景
		CumulativeGasUsed: usedGas,
		GasUsed:           usedGas,
		Logs:              fullState.GetLogs(tx.CalTxId()),
	}, nil
}
