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

package minemsg

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third-party/crypto/cs-crypto"
)

// mine work
type Work interface {
	GetWorkerCoinbaseAddress() common.Address
	SetWorkerCoinbaseAddress(address common.Address)
	//BlockNumber() *big.Int
	// fill seal result to current block
	FillSealResult(curBlock model.AbstractBlock) error
}

type DefaultWork struct {
	WorkerCoinbaseAddress common.Address

	// add block header for mining
	// the first 4 bytes of the nonce value is assigned by the sharding server,
	// and the next 4 assigned by the worker. The rest is accomplished by the miner.
	BlockHeader model.Header
	ResultNonce common.BlockNonce
	//pre-calculate rlp
	RlpPreCal []byte
	// TODO: add diff baseline
}

func (work *DefaultWork) CalHash() (common.Hash, error) {
	if len(work.RlpPreCal) == 0 {
		return common.Hash{}, errors.New("DefaultWork rlp be not calculated yet")
	}

	//extract nonce
	nonce := work.BlockHeader.Nonce
	raw := append(work.RlpPreCal, nonce[:]...)
	//calculate hash
	return cs_crypto.Keccak256Hash(raw), nil
}

func (work *DefaultWork) CalBlockRlpWithoutNonce() {
	work.RlpPreCal = work.BlockHeader.RlpBlockWithoutNonce()
}

func (work *DefaultWork) GetWorkerCoinbaseAddress() common.Address {
	return work.WorkerCoinbaseAddress
}

func (work *DefaultWork) FillSealResult(curBlock model.AbstractBlock) error {
	if work.BlockHeader.Number != curBlock.Number() {
		return errors.New(fmt.Sprintf("work fill result, but height not match, work h: %v, block h: %v", work.BlockHeader.Number, curBlock.Number()))
	}
	curBlock.SetNonce(work.ResultNonce)
	return nil
}

//func (work *DefaultWork) BlockNumber() *big.Int {
//	return big.NewInt(0)
//}

func (work *DefaultWork) SetWorkerCoinbaseAddress(address common.Address) {
	work.WorkerCoinbaseAddress = address
}

// split the work into several slices
func (work *DefaultWork) Split(count int) (result []*DefaultWork) {
	for i := 0; i < count; i++ {
		newWork := *work
		binary.BigEndian.PutUint32(newWork.BlockHeader.Nonce[4:8], uint32(i))
		result = append(result, &newWork)
	}
	return
}
