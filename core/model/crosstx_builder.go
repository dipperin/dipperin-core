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
	"github.com/dipperin/dipperin-core/third-party/crypto/cs-crypto"
	"math/big"
)

//
func CreateRawLockTx(nonce uint64, lock common.Hash, time *big.Int, amount *big.Int, gasPrice *big.Int, gasLimit uint64, alice common.Address, bob common.Address) *Transaction {
	to := cs_crypto.GetLockAddress(alice, bob)
	data := bob.Bytes()
	tempLock := common.CopyHash(&lock)
	txdata := txData{
		AccountNonce: nonce,
		Recipient:    &to,
		HashLock:     tempLock,
		TimeLock:     time,
		Amount:       new(big.Int).Set(amount),
		ExtraData:    data,
	}
	wit := witness{
		R:       new(big.Int),
		S:       new(big.Int),
		V:       new(big.Int),
		HashKey: nil,
	}
	return &Transaction{data: txdata, wit: wit}
}

func CreateRawRefundTx(nonce uint64, amount *big.Int, gasPrice *big.Int, gasLimit uint64, alice common.Address, bob common.Address) *Transaction {
	to := cs_crypto.GetLockAddress(alice, bob)
	data := bob.Bytes()
	txdata := txData{
		AccountNonce: nonce,
		Recipient:    &to,
		HashLock:     nil,
		TimeLock:     new(big.Int),
		Amount:       new(big.Int).Set(amount),
		ExtraData:    data,
	}
	wit := witness{
		R:       new(big.Int),
		S:       new(big.Int),
		V:       new(big.Int),
		HashKey: nil,
	}
	return &Transaction{data: txdata, wit: wit}
}

func CreateRawClaimTx(nonce uint64, key []byte, amount *big.Int, gasPrice *big.Int, gasLimit uint64, alice common.Address, bob common.Address) *Transaction {
	to := cs_crypto.GetLockAddress(alice, bob)
	data := alice.Bytes()
	tempkey := common.CopyBytes(key)
	txdata := txData{
		AccountNonce: nonce,
		Recipient:    &to,
		HashLock:     nil,
		TimeLock:     new(big.Int),
		Amount:       new(big.Int).Set(amount),
		ExtraData:    data,
	}
	wit := witness{
		R:       new(big.Int),
		S:       new(big.Int),
		V:       new(big.Int),
		HashKey: tempkey,
	}
	return &Transaction{data: txdata, wit: wit}

}
