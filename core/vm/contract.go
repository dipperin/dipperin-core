// Copyright 2019, Keychain Foundation Ltd.
// This file is part of the Dipperin-core library.
//
// The Dipperin-core library is free software: you can redistribute
// it and/or modify it under the terms of the GNU Lesser General Public License
// as published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// The Dipperin-core library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

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

	DelegateCall bool
}

func NewContract(caller, object resolver.ContractRef, value *big.Int, gas uint64, input []byte) *Contract {
	return &Contract{
		CallerAddress: caller.Address(),
		caller:        caller,
		self:          object,
		value:         value,
		Gas:           gas,
		Input:         input,
	}
}

func (c *Contract) Caller() resolver.ContractRef {
	return c.caller
}

func (c *Contract) Self() resolver.ContractRef {
	return c.self
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

// UseGas attempts the use gas and subtracts it and returns true on success
func (c *Contract) UseGas(gas uint64) (ok bool) {
	if c.Gas < gas {
		return false
	}
	c.Gas -= gas
	return true
}

func (c *Contract) SetCode(addr *common.Address, hash common.Hash, code []byte) {
	c.Code = code
	c.CodeHash = hash
	c.CodeAddr = addr
}

func (c *Contract) SetAbi(addr *common.Address, hash common.Hash, abi []byte) {
	c.ABI = abi
	c.ABIHash = hash
	c.ABIAddr = addr
}

// AsDelegate sets the contract to be a delegate call and returns the current
// contract (for chaining calls)
func (c *Contract) AsDelegate() *Contract {
	c.DelegateCall = true
	// NOTE: caller must, at all times be a contract. It should never happen
	// that caller is something other than a Contract.
	parent := c.caller.(*Contract)
	c.CallerAddress = parent.CallerAddress
	c.value = parent.value
	return c
}
