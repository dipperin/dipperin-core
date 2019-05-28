package vm

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/vm/resolver"
	"math/big"
)

type Contract struct {
	CallerAddress common.Address
	caller        resolver.ContractRef
	self          resolver.ContractRef

	Code     []byte
	CodeHash common.Hash
	CodeAddr *common.Address
	Input    []byte

	ABI     []byte
	ABIHash common.Hash
	ABIAddr *common.Address

	value *big.Int
	Gas   uint64
}

func NewContract(caller resolver.ContractRef, object resolver.ContractRef, value *big.Int, gas uint64) *Contract {
	return &Contract{
		CallerAddress: caller.Address(),
		caller:        caller,
		self:          object,
		value:         value,
		Gas:           gas,
	}
}

func (c *Contract) Caller() common.Address {
	return c.CallerAddress
}

func (c *Contract) Address() common.Address {
	return c.self.Address()
}

func (c *Contract) CallValue() *big.Int {
	return c.value
}

func (c *Contract) GetGas() uint64 {
	return c.Gas
}

func (c *Contract) SetCallCode(addr *common.Address, hash common.Hash, code []byte) {
	c.Code = code
	c.CodeHash = hash
	c.CodeAddr = addr
}

func (c *Contract) SetCallAbi(addr *common.Address, hash common.Hash, abi []byte) {
	c.ABI = abi
	c.ABIHash = hash
	c.ABIAddr = addr
}

/*
func (c *Contract)GetState(Key []byte) (value []byte){
	return
}
func (c *Contract)SetState(Key []byte, Value []byte) (err error){
	return
}
*/
