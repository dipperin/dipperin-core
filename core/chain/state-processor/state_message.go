package state_processor

import (
	"math/big"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/model"
)

type StateMessage struct {
	tx model.AbstractTransaction
}

func NewStateMessage(tx model.AbstractTransaction) *StateMessage {
	return &StateMessage{
		tx:tx,
	}
}

func (msg StateMessage) From() common.Address {
	panic("implement me")
}

func (msg StateMessage) To() *common.Address {
	return msg.tx.To()
}

func (msg StateMessage) GasPrice() *big.Int {
	return msg.tx.GetGasPrice()
}

func (msg StateMessage) Gas() uint64 {
	return msg.tx.Fee().Uint64()
}

func (msg StateMessage) Value() *big.Int {
	panic("implement me")
}

func (msg StateMessage) Nonce() uint64 {
	panic("implement me")
}

func (msg StateMessage) CheckNonce() bool {
	panic("implement me")
}

func (msg StateMessage) Data() []byte {
	panic("implement me")
}
