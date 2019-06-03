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

func (state *AccountStateDB) ProcessContract(tx model.AbstractTransaction,header *model.Header, create bool,GetHash vm.GetHashFunc) (model.ReceiptPara, error) {
	context := vm.NewVMContext(tx, header,GetHash)
	fullState := &Fullstate{
		state: state,
	}
	msg, err := tx.AsMessage()
	if err != nil {
		log.Error("AccountStateDB#ProcessContract", "as Message err", err)
		return model.ReceiptPara{}, err
	}
	dvm := vm.NewVM(context, fullState, vm.DEFAULT_VM_CONFIG)
	//gasLimit := header.GasLimit
	gasLimit := uint64(2100000) * 10000000000
	_, usedGas, failed, err := ApplyMessage(dvm, msg, &gasLimit)
	if err != nil {
		log.Error("AccountStateDB#ProcessContract", "ApplyMessage err", err)
		return model.ReceiptPara{},err
	}

	root, err := state.Finalise()
	if err != nil {
		log.Error("AccountStateDB#ProcessContract", "state finalise err", err)
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
