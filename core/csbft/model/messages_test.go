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
	"github.com/dipperin/dipperin-core/tests/factory/g-testData"
	"testing"

	"errors"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third-party/crypto"
	"math/big"
	"time"
)

func createBlock(number uint64, txs []*model.Transaction) *model.Block {
	header := model.NewHeader(1, number, common.HexToHash("0x12312fa0929348"), common.HexToHash("1111"), common.HexToDiff("1fffffff"), big.NewInt(time.Now().UnixNano()), common.HexToAddress("032f14"), common.BlockNonceFromInt(432423))
	return model.NewBlock(header, txs, []model.AbstractVerification{})
}

func createRegisterTX(nonce uint64, amount *big.Int) *model.Transaction {
	fs1 := model.NewSigner(big.NewInt(1))
	tx := model.NewRegisterTransaction(nonce, amount, g_testData.TestGasPrice, g_testData.TestGasLimit)
	alicePriv := "289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232031"
	key, _ := crypto.HexToECDSA(alicePriv)
	signedTx, _ := tx.SignTx(key, fs1)
	return signedTx
}

func createCannelTX(nonce uint64) *model.Transaction {
	fs1 := model.NewSigner(big.NewInt(1))
	tx := model.NewCancelTransaction(nonce, g_testData.TestGasPrice, g_testData.TestGasLimit)
	alicePriv := "289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232031"
	key, _ := crypto.HexToECDSA(alicePriv)
	signedTx, _ := tx.SignTx(key, fs1)
	return signedTx
}

func TestNewRoundMsgWithSign2(t *testing.T) {
	addr := common.HexToAddress("0x3f3d")
	NewRoundMsgWithSign(1, 1, func(hash []byte) (b []byte, e error) { return }, addr)
}

func TestNewRoundMsgWithSign3(t *testing.T) {
	addr := common.HexToAddress("0x3f3d")
	NewRoundMsgWithSign(1, 1, func(hash []byte) (b []byte, e error) { return b, errors.New("Error") }, addr)
}

func TestNewRoundMsg_Valid2(t *testing.T) {
	addr := common.HexToAddress("0x3f3d")
	nrmws := NewRoundMsgWithSign(1, 1, func(hash []byte) (b []byte, e error) { return }, addr)
	nrmws.Valid()
}

func TestNewRoundMsg_Valid3(t *testing.T) {
	addr := common.HexToAddress("0x3f3d")
	nrmws := NewRoundMsgWithSign(1, 1, func(hash []byte) (b []byte, e error) { return }, addr)
	nrmws.Witness = nil
	nrmws.Valid()
}

func TestNewProposalWithSign2(t *testing.T) {
	addr := common.HexToAddress("0x3f3d")
	blockID := common.HexToHash("0x9d7d96bfb791080316de884d1f43947764742a7cda226d076b4d42964d00ac92")
	npws := NewProposalWithSign(1, 1, blockID, func(hash []byte) (b []byte, e error) { return }, addr)
	npws.Hash()

	tx1 := createRegisterTX(0, big.NewInt(10000))
	tx2 := createCannelTX(1)
	block := createBlock(2, []*model.Transaction{tx1, tx2})
	npws.ValidBlock(block)
}

func TestNewProposalWithSign3(t *testing.T) {
	addr := common.HexToAddress("0x3f3d")
	blockID := common.HexToHash("0x9d7d96bfb791080316de884d1f43947764742a7cda226d076b4d42964d00ac92")
	NewProposalWithSign(1, 1, blockID, func(hash []byte) (b []byte, e error) { return b, errors.New("Error") }, addr)

}
