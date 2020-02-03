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

package softwallet

import (
	"errors"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/gerror"
	"github.com/dipperin/dipperin-core/core/accounts/accountsbase"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

var testJsonData = `{"accounts":[{"Address":"0x01010101020202020101010102020202030303030303"}],"paths":{"0x01010101020202020101010102020202030303030303":[]},"extend_keys":{"0x01010101020202020101010102020202030303030303":{"Key":"d6L9PvGFt2A/27IeIsVJnqf2+S9eR3u6vWqpZU9vXQ0=","PubKey":null,"ChainCode":"5204jONIM90o0r6810AQsM1WINcnKQEBzZK57t3ldZI=","Depth":0,"ParentFP":"AAAAAA==","ChildNum":0,"Version":"BIit5A==","IsPrivate":true}},"balances":{"0x01010101020202020101010102020202030303030303":0},"Nonce":{"0x01010101020202020101010102020202030303030303":0},"DerivedPathIndex":{},"seed":"AQEBAQICAgIBAQEBAgICAgEBAQECAgICAQEBAQICAgIBAQEBAgICAgEBAQECAgIC"}`



func TestWalletInfo_HdWalletInfoEncodeJson(t *testing.T) {

	encodeData := NewHdWalletInfo()
	tmpAccount := accountsbase.Account{
		Address: [22]byte{0x01, 0x01, 0x01, 0x01, 0x02, 0x02, 0x02, 0x02, 0x01, 0x01, 0x01, 0x01, 0x02, 0x02, 0x02, 0x02,
			0x03, 0x03, 0x03, 0x03, 0x03, 0x03},
	}
	encodeData.Accounts = append(encodeData.Accounts, tmpAccount)
	encodeData.Paths[encodeData.Accounts[0].Address] = accountsbase.DerivationPath{}

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
	testCases := []struct{
		name string
		given []byte
		expect error
	}{
		{
			name:"HdWalletInfoDecodeJsonRight",
			given:[]byte(testJsonData),
			expect:nil,
		},
		{
			name:"HdWalletInfoDecodeJsonErr",
			given:[]byte(testJsonData)[:10],
			expect:errors.New("unexpected end of JSON input"),
		},
	}

	for _,tc := range testCases{
		input := tc.given

		//decoding
		decodeData := NewHdWalletInfo()

		err := decodeData.HdWalletInfoDecodeJson([]byte(input))
		if err != nil {
			assert.Equal(t, tc.expect.Error(), err.Error())
		}
	}
}

func TestWalletInfo_GenerateKeyFromSeedAndPath(t *testing.T) {
	errDrivePath := ""

	testCases := []struct{
		name string
		given func() (*WalletInfo, string, uint32)
		expect error
	}{
		{
			name:"ErrInvalidSeedLen",
			given: func() (*WalletInfo, string, uint32) {
				testWalletInfo := NewHdWalletInfo()
				testWalletInfo.Seed = errSeed
				return testWalletInfo, DefaultDerivedPath, AddressIndexStartValue
			},
			expect:gerror.ErrInvalidSeedLen,
		},
		{
			name:"ErrDerivedPath",
			given: func() (*WalletInfo, string, uint32) {
				testWalletInfo := NewHdWalletInfo()
				testWalletInfo.Seed = testSeed
				return testWalletInfo, errDrivePath, AccountValueIndex
			},
			expect:gerror.ErrDerivedPath,
		},
		{
			name:"GenerateKeyFromSeedAndPathRight",
			given: func() (*WalletInfo, string, uint32) {
				testWalletInfo := NewHdWalletInfo()
				testWalletInfo.Seed = testSeed
				return testWalletInfo, DefaultDerivedPath, AddressIndexStartValue
			},
			expect:nil,
		},
	}

	for _, tc := range testCases {
		wallet, derivedPath, index := tc.given()
		_,_, err := wallet.GenerateKeyFromSeedAndPath(derivedPath, index)
		assert.Equal(t, tc.expect, err)
	}
}


func TestWalletInfo_paddingUsedAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	accountStatus := accountsbase.NewMockAddressInfoReader(ctrl)

	errSeed := []byte{0x01, 0x02, 0x01}

	testCases := []struct{
		name string
		given func() (*WalletInfo, *accountsbase.MockAddressInfoReader)
		expect error
	}{
		{
			name:"ErrInvalidSeedLen",
			given: func() (*WalletInfo, *accountsbase.MockAddressInfoReader) {
				testWalletInfo := NewHdWalletInfo()
				testWalletInfo.Seed = errSeed
				return testWalletInfo,accountStatus
			},
			expect:gerror.ErrInvalidSeedLen,
		},
		{
			name:"paddingUsedAccountRight",
			given: func() (*WalletInfo, *accountsbase.MockAddressInfoReader) {
				testWalletInfo := NewHdWalletInfo()
				testWalletInfo.Seed = testSeed
				accountStatus.EXPECT().GetTransactionNonce(common.HexToAddress("0x00003e406C9dE46907A53A0C39b683747874fb6e9EDC")).Return(uint64(0), nil).AnyTimes()
				accountStatus.EXPECT().GetTransactionNonce(common.HexToAddress("0x0000C3FE396BEc36673626A8dA154161044CfD289A41")).Return(uint64(0), gerror.ErrAccountNotExist).AnyTimes()
				accountStatus.EXPECT().CurrentBalance(common.HexToAddress("0x00003e406C9dE46907A53A0C39b683747874fb6e9EDC")).Return(big.NewInt(100)).AnyTimes()
				accountStatus.EXPECT().CurrentBalance(common.HexToAddress("0x0000C3FE396BEc36673626A8dA154161044CfD289A41")).Return(big.NewInt(100)).AnyTimes()
				return testWalletInfo,accountStatus
			},
			expect:nil,
		},
	}

	for _,tc := range testCases {
		t.Log(tc.name)
		wallet,acc := tc.given()
		err := wallet.paddingUsedAccount(acc)
		assert.Equal(t, tc.expect, err)
	}
}

func TestWalletInfo_getSkFromAddress(t *testing.T) {

	testCases := []struct{
		name string
		given  func()( *WalletInfo, common.Address)
		expect error
	}{
		{
			name:"ErrInvalidAddress",
			given: func() (*WalletInfo, common.Address) {
				testWalletInfo := NewHdWalletInfo()
				return testWalletInfo,  common.Address{}
			},
			expect:gerror.ErrInvalidAddress,
		},
		{
			name:"ErrNotPrivExtKey",
			given: func() (*WalletInfo, common.Address) {
				testWalletInfo := NewHdWalletInfo()
				testWalletInfo.ExtendKeys[testAddress] = ExtendedKey{
					isPrivate: false,
				}
				return testWalletInfo, testAddress
			},
			expect:gerror.ErrNotPrivExtKey,
		},
		{
			name:"getSkFromAddressRight",
			given: func() (*WalletInfo, common.Address) {
				testWalletInfo := NewHdWalletInfo()
				testWalletInfo.ExtendKeys[testAddress] = ExtendedKey{
					isPrivate: true,
				}
				return testWalletInfo, testAddress
			},
			expect:nil,
		},
	}

	for _, tc := range testCases {
		t.Log(tc.name)
		wallet, addr := tc.given()
		_, err := wallet.getSkFromAddress(addr)
		assert.Equal(t, tc.expect, err)
	}
}


// todo
func TestWalletInfo_PaddingAddressNonce(t *testing.T) {
	ctrl := gomock.NewController(t)
	accountStatus := accountsbase.NewMockAddressInfoReader(ctrl)

	testWalletInfo := NewHdWalletInfo()
	testWalletInfo.Accounts = append(testWalletInfo.Accounts, accountsbase.Account{Address: testAddress})

	testCases := []struct{
		name string
		given func() *accountsbase.MockAddressInfoReader
		expect error
	}{
		{
			name :"PaddingAddressNonce",
			given: func()  *accountsbase.MockAddressInfoReader {
				accountStatus.EXPECT().GetTransactionNonce(testAddress).Return(uint64(10), nil)
				return accountStatus
			},
			expect:nil,
		},
		{
			name :"PaddingAddressNonce two",
			given: func()  *accountsbase.MockAddressInfoReader {
				accountStatus.EXPECT().GetTransactionNonce(testAddress).Return(uint64(10),errors.New("error"))
				return accountStatus
			},
			expect:nil,
		},
	}

	for _, tc := range testCases {
		t.Log(tc.name)
		as := tc.given()
		err := testWalletInfo.PaddingAddressNonce(as)
		assert.Equal(t, tc.expect, err)
	}
}

func TestWalletInfo_SetAndGetAddressNonce(t *testing.T) {
	testWalletInfo := NewHdWalletInfo()

	testWalletInfo.SetAddressNonce(testAddress, 2)

	nonce, _ := testWalletInfo.GetAddressNonce(testAddress)
	assert.Equal(t, uint64(2), nonce)

}
