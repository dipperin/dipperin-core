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
	"github.com/dipperin/dipperin-core/tests/g-testData"
	"testing"
	"math/big"
	"github.com/stretchr/testify/assert"

)

func TestNewRegisterTransaction(t *testing.T) {
	k1, _ := CreateKey()
	trans := NewRegisterTransaction(1, big.NewInt(50),  g_testData.TestGasPrice,g_testData.TestGasLimit)
	fs := NewMercurySigner(big.NewInt(1))
	signedTx, _ := trans.SignTx(k1, fs)
	assert.EqualValues(t, signedTx.GetType(), common.AddressTypeStake)
}

func TestNewUnStakeTransaction(t *testing.T) {
	trans := NewUnStakeTransaction(3,  g_testData.TestGasPrice,g_testData.TestGasLimit)
	assert.EqualValues(t, trans.GetType(), common.AddressTypeUnStake)
}

func TestNewCancelTransaction(t *testing.T) {
	trans := NewCancelTransaction(3,  g_testData.TestGasPrice,g_testData.TestGasLimit)
	assert.EqualValues(t, trans.GetType(), common.AddressTypeCancel)
}

func TestNewEvidenceTransaction(t *testing.T) {
	target := common.HexToAddress("target")
	voteA := CreateSignedVote(1, 2, common.HexToHash("0x123456"), VoteMessage)
	trans := NewEvidenceTransaction(3,  g_testData.TestGasPrice,g_testData.TestGasLimit, &target, voteA, voteA)
	assert.EqualValues(t, trans.GetType(), common.AddressTypeEvidence)
}

func TestNewUnNormalTransaction(t *testing.T) {
	trans := NewUnNormalTransaction(3, big.NewInt(5), g_testData.TestGasPrice,g_testData.TestGasLimit)
	assert.EqualValues(t, trans.GetType(), common.TxType(9))
}
