package stateprocessor

import (
	"github.com/dipperin/dipperin-core/common/log"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/core/vm"
	"github.com/dipperin/dipperin-core/core/vm/common"
	"go.uber.org/zap"
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
	msg, err := conf.Tx.AsMessage(true)
	if err != nil {
		log.DLogger.Error("AccountStateDB#ProcessContract", zap.Error(err))
		return model.ReceiptPara{}, err
	}
	dvm := vm.NewVM(context, fullState, common.DEFAULT_VM_CONFIG)
	_, usedGas, failed, fee, err := ApplyMessage(dvm, &msg, conf.GasLimit)
	if err != nil {
		log.DLogger.Error("AccountStateDB#ProcessContract", zap.Error(err))
		return model.ReceiptPara{}, err
	}

	root, err := state.IntermediateRoot()
	if err != nil {
		log.DLogger.Error("AccountStateDB#ProcessContract", zap.Error(err))
		return model.ReceiptPara{}, err
	}

	//padding fee and add block gasUsed
	conf.TxFee = fee
	*conf.GasUsed += usedGas
	log.DLogger.Info("ProcessContract", zap.Uint64("CumulativeGasUsed", *conf.GasUsed), zap.Uint64("usedGas", usedGas))
	return model.ReceiptPara{
		Root:              root[:],
		HandlerResult:     failed,
		CumulativeGasUsed: *conf.GasUsed,
		Logs:              fullState.GetLogs(conf.Tx.CalTxId()),
	}, nil
}
