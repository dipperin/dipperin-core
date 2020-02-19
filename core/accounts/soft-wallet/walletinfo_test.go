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

package soft_wallet

import (
	"errors"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/g-error"
	"github.com/dipperin/dipperin-core/core/accounts"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

var testJsonData = `{"accounts":[{"Address":"0x01010101020202020101010102020202030303030303"}],"paths":{"0x01010101020202020101010102020202030303030303":[]},"extend_keys":{"0x01010101020202020101010102020202030303030303":{"Key":"d6L9PvGFt2A/27IeIsVJnqf2+S9eR3u6vWqpZU9vXQ0=","PubKey":null,"ChainCode":"5204jONIM90o0r6810AQsM1WINcnKQEBzZK57t3ldZI=","Depth":0,"ParentFP":"AAAAAA==","ChildNum":0,"Version":"BIit5A==","IsPrivate":true}},"balances":{"0x01010101020202020101010102020202030303030303":0},"Nonce":{"0x01010101020202020101010102020202030303030303":0},"DerivedPathIndex":{},"seed":"AQEBAQICAgIBAQEBAgICAgEBAQECAgICAQEBAQICAgIBAQEBAgICAgEBAQECAgIC"}`

type fakeAccountStatus struct {
}

func (*fakeAccountStatus) CurrentBalance(address common.Address) *big.Int {
	return big.NewInt(0)
}

func (*fakeAccountStatus) GetTransactionNonce(addr common.Address) (nonce uint64, err error) {
	return 0, g_error.ErrAccountNotExist
}

type errAccountStatus struct {
}

func (*errAccountStatus) CurrentBalance(address common.Address) *big.Int {
	return big.NewInt(0)
}

func (*errAccountStatus) GetTransactionNonce(addr common.Address) (nonce uint64, err error) {
	return 0, errors.New("get nonce other errors")
}

func TestWalletInfo_HdWalletInfoEncodeJson(t *testing.T) {

	encodeData := NewHdWalletInfo()
	tmpAccount := accounts.Account{
		Address: [22]byte{0x01, 0x01, 0x01, 0x01, 0x02, 0x02, 0x02, 0x02, 0x01, 0x01, 0x01, 0x01, 0x02, 0x02, 0x02, 0x02,
			0x03, 0x03, 0x03, 0x03, 0x03, 0x03},
	}
	encodeData.Accounts = append(encodeData.Accounts, tmpAccount)
	encodeData.Paths[encodeData.Accounts[0].Address] = accounts.DerivationPath{}

	encodeData.Balances[encodeData.Accounts[0].Address] = big.NewInt(0)

	encodeData.Seed = []byte{
		0x01, 0x01, 0x01, 0x01, 0x02, 0x02, 0x02, 0x02, 0x01, 0x01, 0x01, 0x01, 0x02, 0x02, 0x02, 0x02,
		0x01, 0x01, 0x01, 0x01, 0x02, 0x02, 0x02, 0x02, 0x01, 0x01, 0x01, 0x01, 0x02, 0x02, 0x02, 0x02,
		0x01, 0x01, 0x01, 0x01, 0x02, 0x02, 0x02, 0x02, 0x01, 0x01, 0x01, 0x01, 0x02, 0x02, 0x02, 0x02,
	}

	//generate master key according to seed
	extKey, err := NewMaster(encodeData.Seed, &DipperinChainCfg)
	assert.NoError(t, err)

	encodeData.ExtendKeys[encodeData.Accounts[0].Address] = *extKey

	//encode json
	encodeResult, err := encodeData.HdWalletInfoEncodeJson()
	assert.NoError(t, err)
	assert.Equal(t, testJsonData, string(encodeResult))

}

func TestWalletInfo_HdWalletInfoDecodeJson(t *testing.T) {

	//decoding
	decodeData := NewHdWalletInfo()

	err := decodeData.HdWalletInfoDecodeJson([]byte(testJsonData))

	assert.NoError(t, err)

	err = decodeData.HdWalletInfoDecodeJson([]byte(testJsonData)[:10])
	assert.Error(t, err)
}

func TestWalletInfo_GenerateKeyFromSeedAndPath(t *testing.T) {
	errDrivePath := ""
	testWalletInfo := NewHdWalletInfo()
	testWalletInfo.Seed = errSeed

	_, _, err := testWalletInfo.GenerateKeyFromSeedAndPath(DefaultDerivedPath, AddressIndexStartValue)
	assert.Error(t, err)

	testWalletInfo.Seed = testSeed
	_, _, err = testWalletInfo.GenerateKeyFromSeedAndPath(errDrivePath, AccountValueIndex)
	assert.Error(t, err)

	testWalletInfo.Seed = testSeed
	_, _, err = testWalletInfo.GenerateKeyFromSeedAndPath(DefaultDerivedPath, AddressIndexStartValue)
	assert.NoError(t, err)

}

func TestWalletInfo_paddingUsedAccount(t *testing.T) {
	errSeed := []byte{0x01, 0x02, 0x01}
	testWalletInfo := NewHdWalletInfo()
	testWalletInfo.Seed = errSeed

	err := testWalletInfo.paddingUsedAccount(&testAccountStatus{})
	assert.Error(t, err)

	testWalletInfo.Seed = testSeed
	err = testWalletInfo.paddingUsedAccount(&fakeAccountStatus{})
	assert.NoError(t, err)

	err = testWalletInfo.paddingUsedAccount(&errAccountStatus{})
	assert.Error(t, err)

	testWalletInfo.Seed = testSeed
	err = testWalletInfo.paddingUsedAccount(&testAccountStatus{})
	assert.NoError(t, err)
}

func TestWalletInfo_getSkFromAddress(t *testing.T) {
	testWalletInfo := NewHdWalletInfo()

	_, err := testWalletInfo.getSkFromAddress(common.Address{})
	assert.Equal(t, accounts.ErrInvalidAddress, err)

	testWalletInfo.ExtendKeys[testAddress] = ExtendedKey{
		isPrivate: false,
	}

	_, err = testWalletInfo.getSkFromAddress(testAddress)
	assert.Equal(t, ErrNotPrivExtKey, err)
}

func TestWalletInfo_PaddingAddressNonce(t *testing.T) {
	testWalletInfo := NewHdWalletInfo()

	testWalletInfo.Accounts = append(testWalletInfo.Accounts, accounts.Account{Address: testAddress})

	err := testWalletInfo.PaddingAddressNonce(&testAccountStatus{})
	assert.NoError(t, err)

	err = testWalletInfo.PaddingAddressNonce(&fakeAccountStatus{})
	assert.NoError(t, err)
}

func TestWalletInfo_SetAndGetAddressNonce(t *testing.T) {
	testWalletInfo := NewHdWalletInfo()

	testWalletInfo.SetAddressNonce(testAddress, 2)

	nonce, _ := testWalletInfo.GetAddressNonce(testAddress)
	assert.Equal(t, uint64(2), nonce)

}
