package state_processor

import (
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/core/vm"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/dipperin/dipperin-core/third-party/crypto/cs-crypto"
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
	//log.Debug("Called create","data",fullState.GetState(addrs,[]byte{7, 98, 97, 108, 97, 110 ,99, 101}),"err",vmerr)
	signer := conf.Tx.GetSigner()
	caller,err := conf.Tx.Sender(signer)
	addr := cs_crypto.CreateContractAddress(caller, conf.Tx.Nonce())
	byteKey := []byte{7, 98, 97, 108, 97, 110 ,99, 101}
	log.Debug("Called process contract","data",state.GetData(addr,string(byteKey)),"err",err)

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
