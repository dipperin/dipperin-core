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

package g_error

import "errors"

var (
	/*Account state errors*/
	ErrTxNonceNotMatch         = errors.New("tx nonce not match")
	ErrTxGasUsedIsOverGasLimit = errors.New("the tx gasUsed is over the gasLimit")
	ErrSenderOrReceiverIsEmpty = errors.New("sender or receiver is empty")
	ErrSenderNotExist          = errors.New("sender not exist")
	ErrInvalidContractData     = errors.New("invalid contract data")
	ErrContractNotExist        = errors.New("contract KV map not exist")
	ErrProhibitFunctionCalled  = errors.New("prohibit function not allow to call")
	ErrAccountNotExist         = errors.New("account does not exist")
	ErrBalanceNegative         = errors.New("balance can not be negative")
	ErrUnknownTxType           = errors.New("unknown tx type")
	ErrAddedLogIsNil           = errors.New("added log is nil")
	ErrTxNotSupported          = errors.New("tx not supported")
	ErrAddressTypeNotMatch     = errors.New("sender address type should be normal")

	/*Verifier processor errors*/
	ErrTxTypeNotMatch        = errors.New("tx type not match with processor function")
	ErrBalanceNotEnough      = errors.New("target balance not enough")
	ErrStakeNotEnough        = errors.New("target stake not enough")
	ErrReceiverNotExist      = errors.New("receiver not exist")
	StateSendRegisterTxFirst = errors.New("processor: need to send register tx first")
	StateSendCancelTxFirst   = errors.New("processor: need to send cancel tx first")

	/*Block processor errors*/
	NotHavePreBlockErr        = errors.New("not have pre block")
	InvalidCoinBaseAddressErr = errors.New("invalid coinBase address")

	/*State transaction errors*/
	ErrNonceTooHigh    = errors.New("nonce too high")
	ErrNonceTooLow     = errors.New("nonce too low")
	ErrGasLimitReached = errors.New("gas limit reached")
)
