package vm

import (
	"github.com/dipperin/dipperin-core/common"
	"math/big"
)

type Contract struct{
	CallerAddress common.Address
	caller ContractRef
	self ContractRef

	ABI []byte
	Code []byte

	value *big.Int
	Gas   uint64
}

func (c *Contract)Caller() common.Address{
	return c.CallerAddress
}

func (c *Contract)Address() common.Address{
	return c.self.Address()
}

func (c *Contract)CallValue() *big.Int{
	return c.value
}

func (c *Contract)GetGas() uint64{
	return c.Gas
}

/*
func (c *Contract)GetState(Key []byte) (value []byte){
	return
}
func (c *Contract)SetState(Key []byte, Value []byte) (err error){
	return
}
*/
