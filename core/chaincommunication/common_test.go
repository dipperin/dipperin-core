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

package chaincommunication

import (
	"crypto/ecdsa"
	"github.com/dipperin/dipperin-core/common"
	iblt "github.com/dipperin/dipperin-core/core/bloom"
	model2 "github.com/dipperin/dipperin-core/core/csbft/model"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third_party/crypto"
	cs_crypto "github.com/dipperin/dipperin-core/third_party/crypto/cs-crypto"
	"math/big"
	"time"
)

func createTestKey() (*ecdsa.PrivateKey, *ecdsa.PrivateKey) {
	testPriv1 := "289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232031"
	testPriv2 := "289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232032"
	key1, err1 := crypto.HexToECDSA(testPriv1)
	key2, err2 := crypto.HexToECDSA(testPriv2)
	if err1 != nil || err2 != nil {
		return nil, nil
	}
	return key1, key2
}

func createTestTx() (*model.Transaction, *model.Transaction) {
	key1, key2 := createTestKey()
	fs1 := model.NewSigner(big.NewInt(1))
	fs2 := model.NewSigner(big.NewInt(3))
	gasPrice := big.NewInt(1)
	gasLimit := 2 * model.TxGas

	testTx1 := model.NewTransaction(10, common.HexToAddress("0121321432423534534534"), big.NewInt(10000), gasPrice, gasLimit, []byte{})
	gasUsed, _ := model.IntrinsicGas(testTx1.ExtraData(), false, false)
	testTx1.PaddingActualTxFee(big.NewInt(0).Mul(big.NewInt(int64(gasUsed)), testTx1.GetGasPrice()))
	testTx1.SignTx(key1, fs1)

	testTx2 := model.NewTransaction(10, common.HexToAddress("0121321432423534534535"), big.NewInt(20000), gasPrice, gasLimit, []byte{})
	gasUsed, _ = model.IntrinsicGas(testTx2.ExtraData(), false, false)
	testTx2.PaddingActualTxFee(big.NewInt(0).Mul(big.NewInt(int64(gasUsed)), testTx2.GetGasPrice()))
	testTx2.SignTx(key2, fs2)

	return testTx1, testTx2
}

func createTestBlock(diff common.Difficulty, number uint64) *model.Block {

	header := &model.Header{Number: number, PreHash: common.HexToHash("0x12312fa0929348"), Diff: diff, Bloom: iblt.NewBloom(model.DefaultBlockBloomConfig)}

	tx1, tx2 := createTestTx()

	txs := []*model.Transaction{tx1, tx2}
	vfy := []model.AbstractVerification{}

	b := model.NewBlock(header, txs, vfy)

	return b
}

func createTestBlockByNumber(number uint64) *model.Block {

	header := model.NewHeader(1, number, common.HexToHash("0x12312fa0929348"), common.HexToHash("1111"), common.HexToDiff("1fffffff"), big.NewInt(time.Now().UnixNano()), common.HexToAddress("032f14"), common.BlockNonceFromInt(432423))

	tx1, tx2 := createTestTx()
	txs := []*model.Transaction{tx1, tx2}
	vfy := []model.AbstractVerification{}

	b := model.NewBlock(header, txs, vfy)

	return b
}

func chainHeight() model.AbstractBlock {
	return createTestBlockByNumber(2)
}

func getBlockByHashReturnNil(hash common.Hash) model.AbstractBlock {
	return nil
}

func saveBlock(block model.AbstractBlock, seenCommits []model.AbstractVerification) error {
	return nil
}

func broadcast(b *model2.VerifyResult) {
}

func getFetcher() *BlockFetcher {
	return NewBlockFetcher(chainHeight, getBlockByHashReturnNil, saveBlock, broadcast)
}

func getVr(hash common.Hash) error {
	// mock get vr p2p send request
	return nil
}

var testAccFactory = &AccountFactory{}

var defaultAccounts = []Account{
	{Pk: crypto.HexToECDSAErrPanic("fe10ee89565549d616d43c4e71b61d46a963fdb69489093a57cacf06836ecd91")},
	{Pk: crypto.HexToECDSAErrPanic("fe20ee89565549d616d43c4e71b61d46a963fdb69489093a57cacf06836ecd91")},
	{Pk: crypto.HexToECDSAErrPanic("fe30ee89565549d616d43c4e71b61d46a963fdb69489093a57cacf06836ecd91")},
	{Pk: crypto.HexToECDSAErrPanic("fe40ee89565549d616d43c4e71b61d46a963fdb69489093a57cacf06836ecd91")},
	{Pk: crypto.HexToECDSAErrPanic("fe50ee89565549d616d43c4e71b61d46a963fdb69489093a57cacf06836ecd91")},
	{Pk: crypto.HexToECDSAErrPanic("fe60ee89565549d616d43c4e71b61d46a963fdb69489093a57cacf06836ecd91")},
	{Pk: crypto.HexToECDSAErrPanic("fe70ee89565549d616d43c4e71b61d46a963fdb69489093a57cacf06836ecd91")},
	{Pk: crypto.HexToECDSAErrPanic("fe80ee89565549d616d43c4e71b61d46a963fdb69489093a57cacf06836ecd91")},
	{Pk: crypto.HexToECDSAErrPanic("fe90ee89565549d616d43c4e71b61d46a963fdb69489093a57cacf06836ecd91")},
	{Pk: crypto.HexToECDSAErrPanic("fea0ee89565549d616d43c4e71b61d46a963fdb69489093a57cacf06836ecd91")},
	{Pk: crypto.HexToECDSAErrPanic("feb0ee89565549d616d43c4e71b61d46a963fdb69489093a57cacf06836ecd91")},
}

// get or gen account
type AccountFactory struct{}

// gen account
func (acc *AccountFactory) GenAccount() Account {
	pk, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}
	return Account{Pk: pk}
}

// gen x accounts
func (acc *AccountFactory) GenAccounts(x int) (r []Account) {
	for i := 0; i < x; i++ {
		pk, err := crypto.GenerateKey()
		if err != nil {
			panic(err)
		}
		r = append(r, Account{Pk: pk})
	}
	return
}

// get default account(have certain address)
func (acc *AccountFactory) GetAccount(index int) Account {
	return defaultAccounts[index]
}

type Account struct {
	Pk      *ecdsa.PrivateKey
	address common.Address
}

func NewAccount(pk *ecdsa.PrivateKey, address common.Address) *Account {
	return &Account{Pk: pk, address: address}
}

func (a *Account) Address() common.Address {
	if !a.address.IsEmpty() {
		return a.address
	}

	a.address = cs_crypto.GetNormalAddress(a.Pk.PublicKey)
	return a.address
}

func (a *Account) SignHash(hash []byte) ([]byte, error) {
	return crypto.Sign(hash, a.Pk)
}
