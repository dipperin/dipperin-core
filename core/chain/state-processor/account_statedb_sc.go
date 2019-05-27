package state_processor

import (
	"github.com/dipperin/dipperin-core/core/contract"
	"github.com/dipperin/dipperin-core/core/model"
)

func (state *AccountStateDB) ProcessSmartContract(tx model.AbstractTransaction, blockHeight uint64) (err error) {
	cProcessor := contract.NewProcessor(state, blockHeight)
	cProcessor.ProcessSmartContract(state, tx, blockHeight)
    return nil
}