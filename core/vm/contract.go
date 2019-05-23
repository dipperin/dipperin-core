package vm

import (
	"github.com/dipperin/dipperin-core/common"
)

type Contract struct{
	CallerAddress common.Address
	caller ContractRef
	self ContractRef

	ABI []byte
	Code []byte
}

/*
func (c *Contract)GetState(Key []byte) (value []byte){
	return
}
func (c *Contract)SetState(Key []byte, Value []byte) (err error){
	return
}
*/
