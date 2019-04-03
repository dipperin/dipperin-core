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
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
	"github.com/dipperin/dipperin-core/third-party/crypto/cs-crypto"
)

type Proofs struct {
	VoteA    *VoteMsg
	VoteB    *VoteMsg
	VRFHash  common.Hash
	Proof    []byte
	Priority uint64
}

/*
Name
CalledBy
Parameters
Return
 */
func NewRegisterTransaction(nonce uint64, amount *big.Int, fee *big.Int) *Transaction {
	target := common.HexToAddress(common.AddressStake)
	txData := txData{
		AccountNonce: nonce,
		Recipient:    &target,
		//HashLock:    common.HexToHash(""),
		TimeLock:  new(big.Int),
		Amount:    new(big.Int),
		Fee:       new(big.Int),
		ExtraData: []byte{},
	}
	wit := witness{
		R:       new(big.Int),
		S:       new(big.Int),
		V:       new(big.Int),
		HashKey: nil,
	}
	if amount != nil {
		txData.Amount.Set(amount)
	}
	if fee != nil {
		txData.Fee.Set(fee)
	}
	return &Transaction{data: txData, wit: wit}
}

func NewEvidenceTransaction(nonce uint64, fee *big.Int, target *common.Address, voteA *VoteMsg, voteB *VoteMsg) *Transaction {
	var emptyHash common.Hash
	proofs := Proofs{voteA, voteB, emptyHash, nil, 0}
	data, _ := rlp.EncodeToBytes(proofs)
	to := cs_crypto.GetEvidenceAddress(*target)
	txData := txData{
		AccountNonce: nonce,
		Recipient:    &to,
		TimeLock:     new(big.Int),
		Amount:       new(big.Int),
		Fee:          new(big.Int),
		ExtraData:    data,
	}
	wit := witness{
		R:       new(big.Int),
		S:       new(big.Int),
		V:       new(big.Int),
		HashKey: nil,
	}
	if fee != nil {
		txData.Fee.Set(fee)
	}
	return &Transaction{data: txData, wit: wit}
}

func NewUnStakeTransaction(nonce uint64, fee *big.Int) *Transaction {
	target := common.HexToAddress(common.AddressUnStake)
	txData := txData{
		AccountNonce: nonce,
		Recipient:    &target,
		TimeLock:     new(big.Int),
		Amount:       new(big.Int),
		Fee:          new(big.Int),
		ExtraData:    []byte{},
	}
	wit := witness{
		R:       new(big.Int),
		S:       new(big.Int),
		V:       new(big.Int),
		HashKey: nil,
	}
	if fee != nil {
		txData.Fee.Set(fee)
	}
	return &Transaction{data: txData, wit: wit}
}

func NewCancelTransaction(nonce uint64, fee *big.Int) *Transaction {
	target := common.HexToAddress(common.AddressCancel)
	txData := txData{
		AccountNonce: nonce,
		Recipient:    &target,
		TimeLock:     new(big.Int),
		Amount:       new(big.Int),
		Fee:          new(big.Int),
		ExtraData:    []byte{},
	}
	wit := witness{
		R:       new(big.Int),
		S:       new(big.Int),
		V:       new(big.Int),
		HashKey: nil,
	}
	if fee != nil {
		txData.Fee.Set(fee)
	}
	return &Transaction{data: txData, wit: wit}
}

func NewUnNormalTransaction(nonce uint64, amount *big.Int, fee *big.Int) *Transaction {
	target := common.HexToAddress("0x00090000000000000000000000000000000000000000")
	txData := txData{
		AccountNonce: nonce,
		Recipient:    &target,
		//HashLock:    common.HexToHash(""),
		TimeLock:  new(big.Int),
		Amount:    new(big.Int),
		Fee:       new(big.Int),
		ExtraData: []byte{},
	}
	wit := witness{
		R:       new(big.Int),
		S:       new(big.Int),
		V:       new(big.Int),
		HashKey: nil,
	}
	if amount != nil {
		txData.Amount.Set(amount)
	}
	if fee != nil {
		txData.Fee.Set(fee)
	}
	return &Transaction{data: txData, wit: wit}
}
