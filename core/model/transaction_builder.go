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

package model

import (
	"github.com/dipperin/dipperin-core/common"
	"math/big"
)

type chainStateReader interface {
	GetAccountNonce(address common.Address) uint64
}

type TransactionBuilder struct {
	chainStateReader chainStateReader
}

func (builder *TransactionBuilder) NewNormalTransaction(from common.Address, to common.Address, amount *big.Int, fee *big.Int, data []byte) AbstractTransaction {
	//return NewTransaction(builder.chainStateReader.GetAccountNonce(from), to, amount, fee, data)
	return nil
}
