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
	"encoding/hex"
	"github.com/dipperin/dipperin-core/common/gerror"
	"github.com/dipperin/dipperin-core/common/log"
	"github.com/dipperin/dipperin-core/core/accounts/accountsbase"
	"github.com/dipperin/dipperin-core/third_party/go-bip39"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
	"testing"
)

var testWalletPlain = [12]byte{0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12}
var errIv = []byte{0x11, 0x11, 0x11}
var testIv = [16]byte{
	0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12}
var testKey = [32]byte{
	0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12,
	0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12}
var encKey = EncryptKey{
	encryptKey: testKey,
	macKey:     testKey,
}
var testWalletCipher = WalletCipher{
	Cipher:    "6eb6e9b5743e41847cdb27bb428bbf60",
	MacCipher: "2511f7e98fad50a19aac15e208a8e3de3e6c2dd4e482d953282877e49ece3177",
}

func TestGenerateMnemonic(t *testing.T) {
	testCases := []struct{
		name string
		given int
		expect error
	}{
		{
			name:"err",
			given:33,
			expect:bip39.ErrEntropyLengthInvalid,
		},
		{
			name:"GenerateMnemonicRight",
			given:128,
			expect:nil,
		},
	}

	for _,tc := range testCases{
		bitSize := tc.given
		_, err := GenerateMnemonic(bitSize)
		assert.Equal(t,tc.expect, err )
	}
}

func TestGenSymKeyFromPassword(t *testing.T) {
	testKdfPara := KDFParameter{
		KDFParams: map[string]interface{}{
			"n":       WalletLightScryptN,
			"p":       WalletLightScryptP,
			"kdfType": KDF,
			"r":       WalletscryptR,
			"keyLen":  WalletscryptDKLen,
		},
	}

	testCases := []struct{
		name string
		given func() (string, KDFParameter)
		expect error
	}{
		{
			name:"GenSymKeyFromPasswordRight",
			given: func() (string, KDFParameter) {
				testPassword := "123"
				return testPassword,testKdfPara
			},
			expect:nil,
		},
		{
			name:"ErrNotSupported",
			given: func() (string, KDFParameter) {
				testPassword := "123"
				testKdfPara.KDFParams["kdfType"] = "PBKDF2"
				return testPassword,testKdfPara
			},
			expect:gerror.ErrNotSupported,
		},
		{
			name:"ErrInvalidKDFParameter",
			given: func() (string, KDFParameter) {
				testPassword := "123"
				testKdfPara.KDFParams["kdfType"] = KDF
				testKdfPara.KDFParams["keyLen"] = 12
				return testPassword,testKdfPara
			},
			expect:gerror.ErrInvalidKDFParameter,
		},
		{
			name:"err",
			given: func() (string, KDFParameter) {
				testPassword := "123"
				testKdfPara.KDFParams["keyLen"] = WalletscryptDKLen
				testKdfPara.KDFParams["salt"] = "123"
				return testPassword,testKdfPara
			},
			expect:gerror.ErrOddLenghtHexString,
		},
		{
			name:"ErrDeriveKey",
			given: func() (string, KDFParameter) {
				testPassword := "123"
				testKdfPara.KDFParams["salt"] = ""
				testKdfPara.KDFParams["n"] = 1
				return testPassword,testKdfPara
			},
			expect:gerror.ErrDeriveKey,
		},
	}

	for _,tc := range testCases{
		passwd, kdfPara := tc.given()
		_, err := GenSymKeyFromPassword(passwd, kdfPara)
		assert.Equal(t, tc.expect, err)
	}

}

func TestGetAccountFromExtendedKey(t *testing.T) {
	testKey, err := NewMaster(testSeed, &DipperinChainCfg)
	assert.NoError(t, err)

	testKey.isPrivate = false
	_, err = GetAccountFromExtendedKey(testKey)
	assert.Error(t, err)
}

// todo
func TestEncryptWalletContent(t *testing.T) {
	log.InitLogger(log.LoggerConfig{
		Lvl:         zapcore.DebugLevel,
		WithConsole: true,
	})

	cipher, err := EncryptWalletContent(testWalletPlain[:], testIv[:], encKey)
	assert.NoError(t, err)
	assert.Equal(t, testWalletCipher, cipher)

	assert.Panics(t,
		func() {
			EncryptWalletContent(testWalletPlain[:], errIv, encKey)
		},
	)
}

func TestDecryptWalletContent(t *testing.T) {

	testCases := []struct{
		name string
		given func() (walletCipher WalletCipher, iv []byte, sysKey EncryptKey)
		expect error
		expectPlain []byte
	}{
		{
			name:"DecryptWalletContentRight",
			given: func() (walletCipher WalletCipher, iv []byte, sysKey EncryptKey) {
				return testWalletCipher, testIv[:], encKey
			},
			expect:nil,
			expectPlain:testWalletPlain[:],
		},
		{
			name:"err",
			given: func() (walletCipher WalletCipher, iv []byte, sysKey EncryptKey) {
				errWalletCipher := WalletCipher{
					Cipher:    "12324",
					MacCipher: "qerwqer",
				}
				return errWalletCipher, testIv[:], encKey
			},
			expect:gerror.ErrOddLenghtHexString,
		},
		{
			name:"ErrAESDecryption",
			given: func() (walletCipher WalletCipher, iv []byte, sysKey EncryptKey) {
				errWalletCipher := WalletCipher{
					Cipher:    "123244",
					MacCipher: "qerwqer",
				}
				return errWalletCipher, testIv[:], encKey
			},
			expect:gerror.ErrAESDecryption,
		},
		{
			name:"err",
			given: func() (walletCipher WalletCipher, iv []byte, sysKey EncryptKey) {
				errWalletCipher := WalletCipher{
					Cipher:    testWalletCipher.Cipher,
					MacCipher: "qerwqer",
				}
				return errWalletCipher, testIv[:], encKey
			},
			expect:hex.InvalidByteError(0x71),
		},
		{
			name:"ErrAESDecryption",
			given: func() (walletCipher WalletCipher, iv []byte, sysKey EncryptKey) {
				errWalletCipher := WalletCipher{
					Cipher:    testWalletCipher.Cipher,
					MacCipher: "123244",
				}
				return errWalletCipher, testIv[:], encKey
			},
			expect:gerror.ErrAESDecryption,
		},
	}

	for _, tc := range testCases{
		walletCipher, iv, encKey := tc.given()
		plain, err := DecryptWalletContent(walletCipher, iv, encKey)
		assert.Equal(t, tc.expect, err)
		assert.Equal(t, tc.expectPlain, plain)
	}

}

func TestCheckPassword(t *testing.T) {
	testCases := []struct{
		name string
		given string
		expect error
	}{
		{
			name:"right one",
			given:"19abc```",
			expect:nil,
		},
		{
			name:"",
			given:"1234567",
			expect:gerror.ErrPasswordOrPassPhraseIllegal,
		},
		{
			name:"can not have chinese ",
			given:"国1234567",
			expect:gerror.ErrPasswordOrPassPhraseIllegal,
		},
		{
			name:"too long",
			given:"1234567890asertyuiopasdfh",
			expect:gerror.ErrPasswordOrPassPhraseIllegal,
		},
		{
			name:"right two",
			given:"234567890~!@#$%^&*()_+<",
			expect:nil,
		},
	}

	for _,tc := range testCases{
		err := CheckPassword(tc.given)
		assert.Equal(t, tc.expect, err)
	}
}

func TestCheckDerivedPathValid(t *testing.T) {
	result, _ := CheckDerivedPathValid(accountsbase.DerivationPath{0x12, 0x12})
	assert.Equal(t, false, result)
}

func TestCheckWalletPath(t *testing.T) {
	err := CheckWalletPath("/tmp")
	assert.Error(t, err)
}
