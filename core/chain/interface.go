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


package chain

import (
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/common"
)

type Chain interface {
	GetBlockByNumber(number uint64) model.AbstractBlock
	CurrentBlock() model.AbstractBlock
	Genesis() model.AbstractBlock
	GetCurrVerifiers() []common.Address
	GetNextVerifiers() []common.Address
	GetHeaderByHash(hash common.Hash) model.AbstractHeader
	GetHeaderByNumber(number uint64) model.AbstractHeader
	IsChangePoint(block model.AbstractBlock, isProcessPackageBlock bool) bool
}