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
	"github.com/dipperin/dipperin-core/common/log"
	"github.com/dipperin/dipperin-core/core/accounts"
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
	errBitSize := 33
	_, err := GenerateMnemonic(errBitSize)
	assert.Error(t, err)

	testBitSize := 128
	_, err = GenerateMnemonic(testBitSize)
	assert.NoError(t, err)
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
	testPassword := "123"
	_, err := GenSymKeyFromPassword(testPassword, testKdfPara)
	assert.NoError(t, err)

	testKdfPara.KDFParams["kdfType"] = "PBKDF2"
	_, err = GenSymKeyFromPassword(testPassword, testKdfPara)
	assert.Equal(t, accounts.ErrNotSupported, err)

	testKdfPara.KDFParams["kdfType"] = KDF
	testKdfPara.KDFParams["keyLen"] = 12
	_, err = GenSymKeyFromPassword(testPassword, testKdfPara)
	assert.Equal(t, accounts.ErrInvalidKDFParameter, err)

	testKdfPara.KDFParams["keyLen"] = WalletscryptDKLen
	testKdfPara.KDFParams["salt"] = "123"
	_, err = GenSymKeyFromPassword(testPassword, testKdfPara)
	assert.Error(t, err)

	testKdfPara.KDFParams["salt"] = ""
	testKdfPara.KDFParams["n"] = 1
	_, err = GenSymKeyFromPassword(testPassword, testKdfPara)
	assert.Equal(t, accounts.ErrDeriveKey, err)
}

func TestGetAccountFromExtendedKey(t *testing.T) {
	testKey, err := NewMaster(testSeed, &DipperinChainCfg)
	assert.NoError(t, err)

	testKey.isPrivate = false
	_, err = GetAccountFromExtendedKey(testKey)
	assert.Error(t, err)
}

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
	plain, err := DecryptWalletContent(testWalletCipher, testIv[:], encKey)
	assert.NoError(t, err)
	assert.Equal(t, testWalletPlain[:], plain)

	errWalletCipher := WalletCipher{
		Cipher:    "12324",
		MacCipher: "qerwqer",
	}
	_, err = DecryptWalletContent(errWalletCipher, testIv[:], encKey)
	assert.Error(t, err)

	errWalletCipher.Cipher = "123244"
	_, err = DecryptWalletContent(errWalletCipher, testIv[:], encKey)
	assert.Equal(t, accounts.ErrAESDecryption, err)

	errWalletCipher.Cipher = testWalletCipher.Cipher
	_, err = DecryptWalletContent(errWalletCipher, testIv[:], encKey)
	assert.Error(t, err)

	errWalletCipher.MacCipher = "123244"
	_, err = DecryptWalletContent(errWalletCipher, testIv[:], encKey)
	assert.Equal(t, accounts.ErrAESDecryption, err)
}

func TestCheckPassword(t *testing.T) {
	//err := CheckPassword("")
	//assert.Equal(t, errors.New("password is nil"), err)

	err := CheckPassword("19abc```")
	assert.NoError(t, err)

	err = CheckPassword("国1234567")
	assert.Error(t, err)

	err = CheckPassword("1234567")
	assert.Error(t, err)

	err = CheckPassword("1234567890asertyuiopasdfh")
	assert.Error(t, err)

	err = CheckPassword("234567890~!@#$%^&*()_+<")
	assert.NoError(t, err)
}

func TestCheckDerivedPathValid(t *testing.T) {
	result, _ := CheckDerivedPathValid(accounts.DerivationPath{0x12, 0x12})
	assert.Equal(t, false, result)
}

func TestCheckWalletPath(t *testing.T) {
	err := CheckWalletPath("/tmp")
	assert.Error(t, err)
}
