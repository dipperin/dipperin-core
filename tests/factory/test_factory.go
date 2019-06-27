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


package factory

import (
	"crypto/ecdsa"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/g-testData"
	"github.com/dipperin/dipperin-core/core/bloom"
	"github.com/dipperin/dipperin-core/core/model"
	"math/big"
	"time"
	"github.com/dipperin/dipperin-core/third-party/crypto"
)

var testPriv1 = "289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232031"
var testPriv2 = "289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232032"

func CreateKey() (*ecdsa.PrivateKey, *ecdsa.PrivateKey) {
	key1, err1 := crypto.HexToECDSA(testPriv1)
	key2, err2 := crypto.HexToECDSA(testPriv2)
	if err1 != nil || err2 != nil {
		return nil, nil
	}
	return key1, key2
}

func CreateTestTx() (*model.Transaction, *model.Transaction) {
	key1, key2 := CreateKey()
	fs1 := model.NewMercurySigner(big.NewInt(1))
	fs2 := model.NewMercurySigner(big.NewInt(3))
	//alice := crypto.GetNormalAddress(key1.PublicKey)
	//bob := crypto.GetNormalAddress(key2.PublicKey)
	//hashkey := []byte("123")
	//hashlock := cs_crypto.Keccak256Hash(hashkey)
	testtx1 := model.NewTransaction(10, common.HexToAddress("0121321432423534534534"), big.NewInt(10000),g_testData.TestGasPrice,g_testData.TestGasLimit, []byte{})
	testtx1.SignTx(key1, fs1)
	testtx2 := model.NewTransaction(10, common.HexToAddress("0121321432423534534535"), big.NewInt(20000),g_testData.TestGasPrice,g_testData.TestGasLimit, []byte{})
	//testtx2 := model.CreateRawLockTx(1, hashlock, big.NewInt(34564), big.NewInt(10000), big.NewInt(100), alice, bob)
	testtx2.SignTx(key2, fs2)
	return testtx1, testtx2
}

func CreateSpecialBlock(number uint64) *model.Block {

	header := model.NewHeader(1, number, common.HexToHash("0x12312fa0929348"), common.HexToHash("1111"), common.Difficulty{0}, big.NewInt(time.Now().UnixNano()), common.HexToAddress("032f14"), common.BlockNonce{0})
	b := model.NewBlock(header, []*model.Transaction{}, []model.AbstractVerification{})

	return b
}

func CreateBlock(number uint64) *model.Block {

	header := model.NewHeader(1, number, common.HexToHash("0x12312fa0929348"), common.HexToHash("1111"), common.HexToDiff("1fffffff"), big.NewInt(time.Now().UnixNano()), common.HexToAddress("032f14"), common.BlockNonceFromInt(432423))

	tx1, tx2 := CreateTestTx()
	txs := []*model.Transaction{tx1, tx2}
	vfy := []model.AbstractVerification{}

	b := model.NewBlock(header, txs, vfy)

	return b
}

func CreateBlock2(diff common.Difficulty,number uint64) *model.Block {

	header := &model.Header{Number: number, PreHash: common.HexToHash("0x12312fa0929348"),Diff:diff, Bloom: iblt.NewBloom(model.DefaultBlockBloomConfig)}

	tx1, tx2 := CreateTestTx()
	txs := []*model.Transaction{tx1, tx2}
	vfy := []model.AbstractVerification{}

	b := model.NewBlock(header, txs, vfy)

	return b
}

func CreateBlockByPH(number uint64, preHash common.Hash) *model.Block {
	header := model.NewHeader(1, number, preHash, common.HexToHash("1111"), common.HexToDiff("1fffffff"), big.NewInt(time.Now().UnixNano()), common.HexToAddress("032f14"), common.BlockNonceFromInt(432423))

	tx1, tx2 := CreateTestTx()
	txs := []*model.Transaction{tx1, tx2}
	vfy := []model.AbstractVerification{}
	b := model.NewBlock(header, txs, vfy)

	return b
}

func CreateBlockBySeed(number uint64, preHash, seed common.Hash, proof, pubKey []byte) *model.Block {
	header := model.NewHeader(1, number, preHash, seed, common.HexToDiff("1fffffff"), big.NewInt(time.Now().UnixNano()), common.HexToAddress("032f14"), common.BlockNonceFromInt(432423))
	header.Proof = proof
	header.MinerPubKey = pubKey

	tx1, tx2 := CreateTestTx()
	txs := []*model.Transaction{tx1, tx2}
	vfy := []model.AbstractVerification{}
	b := model.NewBlock(header, txs, vfy)
	return b
}