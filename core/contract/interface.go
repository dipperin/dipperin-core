// Copyright 2019, Keychain Foundation Ltd.
// This file is part of the dipperin-core library.
//
// The dipperin-core library is free software: you can redistribute
// it and/or modify it under the terms of the GNU Lesser General Public License
// as published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// The dipperin-core library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.


package contract

import (
	"github.com/dipperin/dipperin-core/common"
	"math/big"
	"reflect"
)

type ContractDB interface {
	PutContract(addr common.Address, v reflect.Value) error
	GetContract(addr common.Address, vType reflect.Type) (v reflect.Value, err error)
	ContractExist(addr common.Address) bool
}

type AccountDB interface {
	GetBalance(addr common.Address) (*big.Int, error)
	AddBalance(addr common.Address, amount *big.Int) error
	SubBalance(addr common.Address, amount *big.Int) error
}