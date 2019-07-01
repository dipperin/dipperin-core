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


package tests

import (
	"crypto/ecdsa"
	"github.com/dipperin/dipperin-core/tests/g-testData"
	"github.com/dipperin/dipperin-core/third-party/crypto/cs-crypto"
	"math/big"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third-party/crypto"
	"time"
	"github.com/dipperin/dipperin-core/core/chain-config"
)

var txSigner = model.NewMercurySigner(chain_config.GetChainConfig().ChainId)

// 交易builder
type TxBuilder struct {
	Nonce  uint64
	To     common.Address
	Amount *big.Int
	Fee    *big.Int
	Data   []byte
	Pk *ecdsa.PrivateKey
}

func (b *TxBuilder) From() common.Address {
	return cs_crypto.GetNormalAddress(b.Pk.PublicKey)
}

func (b *TxBuilder) Build() *model.Transaction {
	tx := model.NewTransaction(b.Nonce, b.To, b.Amount,g_testData.TestGasPrice,g_testData.TestGasLimit, b.Data)
	tx.SignTx(b.Pk, txSigner)
	return tx
}

func (b *TxBuilder) BuildAbs() model.AbstractTransaction {
	return b.Build()
}

type VerBuilder struct {
	Round uint64
	VoteType model.VoteMsgType
	Block model.AbstractBlock
	Pk *ecdsa.PrivateKey
}

func (b *VerBuilder) Build() model.AbstractVerification {
	msg := &model.VoteMsg{
		Height:    b.Block.Number(),
		Round:     b.Round,
		BlockID:   b.Block.Hash(),
		VoteType:  b.VoteType,
		Timestamp: time.Now(),
	}

	// sign msg
	sign, err := crypto.Sign(msg.Hash().Bytes(), b.Pk)
	errPanic(err)

	msg.Witness = &model.WitMsg{
		Address: cs_crypto.GetNormalAddress(b.Pk.PublicKey),
		Sign:    sign,
	}

	return msg
}