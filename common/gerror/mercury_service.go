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

package gerror

import "errors"

var (
	BlockHashNotFound            = errors.New("the block hash is not found")
	BlockNumberError             = errors.New("the block number is smaller than 2")
	BlockIsNilError              = errors.New("the block is nil")
	BeginNumLargerError          = errors.New("begin num is larger than end num")
	ErrEmptyTxData               = errors.New("empty tx data")
	ErrInvalidContractType       = errors.New("invalid contract type")
	ErrFunctionCalledConstant    = errors.New("function called is constant, no need to send transaction")
	ErrFunctionCalledNotConstant = errors.New("function called isn't constant, need to send a transaction")
	ErrFunctionInitCanNotCalled  = errors.New("function init can't be called")
	ErrFuncNameNotFoundInABI     = errors.New("funcName not found in abi")
	ErrBloombitsNotFound         = errors.New("can't find the bloombits")
	ErrReceiptIsNil              = errors.New("the transaction receipt is nil")
	ErrReceiptNotFound           = errors.New("the transaction receipt not found")
)
