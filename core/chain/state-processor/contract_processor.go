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

func (state *AccountStateDB) ProcessContract(conf *TxProcessConfig, create bool) (model.ReceiptPara, error) {
	context := vm.NewVMContext(conf.Tx, conf.Header, conf.GetHash)
	fullState := &Fullstate{
		state: state,
	}
	msg, err := conf.Tx.AsMessage()
	if err != nil {
		log.Error("AccountStateDB#ProcessContract", "as Message err", err)
		return model.ReceiptPara{}, err
	}
	dvm := vm.NewVM(context, fullState, vm.DEFAULT_VM_CONFIG)
	_, usedGas, failed, fee, err := ApplyMessage(dvm, msg, conf.GasLimit)

	if err != nil {
		log.Error("AccountStateDB#ProcessContract", "ApplyMessage err", err)
		return model.ReceiptPara{}, err
	}

	root, err := state.IntermediateRoot()
	if err != nil {
		log.Error("AccountStateDB#ProcessContract", "state finalise err", err)
		return model.ReceiptPara{}, err
	}

	//padding fee and add block gasUsed
	conf.TxFee = fee
	*conf.GasUsed += usedGas
	log.Info("ProcessContract", "CumulativeGasUsed", *conf.GasUsed, "usedGas", usedGas)
	return model.ReceiptPara{
		Root:              root[:],
		HandlerResult:     failed,
		CumulativeGasUsed: *conf.GasUsed,
		GasUsed:           usedGas,
		Logs:              fullState.GetLogs(conf.Tx.CalTxId()),
	}, nil
}
