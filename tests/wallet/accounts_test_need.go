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

package wallet

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/core/accounts"
	"github.com/dipperin/dipperin-core/core/accounts/soft-wallet"
	"math/big"
	"os"
	"path/filepath"
)

var ErrSeed = []byte{0x01, 0x02, 0x01}
var TestSeed = []byte{0x01, 0x02, 0x01, 0x02, 0x01, 0x02, 0x01, 0x02, 0x01, 0x02, 0x01, 0x02,
	0x01, 0x02, 0x01, 0x02, 0x01, 0x02, 0x01, 0x02, 0x01, 0x02, 0x01, 0x02}
var TestAddress = common.Address{0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12,
	0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12}

//签名hash值
var TestHashData = [32]byte{0x01, 0x02, 0x03, 0x04, 0x01, 0x02, 0x03, 0x04,
	0x01, 0x02, 0x03, 0x04, 0x01, 0x02, 0x03, 0x04,
	0x01, 0x02, 0x03, 0x04, 0x01, 0x02, 0x03, 0x04,
	0x01, 0x02, 0x03, 0x04, 0x01, 0x02, 0x03, 0x04}

var TestWalletName = "testSoftWallet1"
var Path = filepath.Join(util.HomeDir(), "/tmp/testSoftWallet1")

type testGetAccountInfo struct {
}

func (*testGetAccountInfo) CurrentBalance(address common.Address) *big.Int {
	return big.NewInt(0)
}

func (*testGetAccountInfo) GetTransactionNonce(addr common.Address) (nonce uint64, err error) {
	return 0, nil
}

var tmpAccountStatus testGetAccountInfo

func EstablishSoftWallet(path, walletName string) (*soft_wallet.SoftWallet, error) {
	testWallet, err := soft_wallet.NewSoftWallet()
	if err != nil {
		return nil, err
	}

	testWallet.Identifier.WalletName = walletName
	testWallet.Identifier.Path = path

	os.Remove(path)

	_, err = testWallet.Establish(path, walletName, "12345678", "12345678")
	if err != nil {
		return nil, err
	}

	return testWallet, nil
}

func GetTestWalletManager() (testWallet *soft_wallet.SoftWallet, manager *accounts.WalletManager, err error) {

	testWallet, err = EstablishSoftWallet(Path, TestWalletName)
	if err != nil {
		return nil, nil, err
	}

	walletManager, err := accounts.NewWalletManager(&tmpAccountStatus, testWallet)
	if err != nil {
		return nil, nil, err
	}
	return testWallet, walletManager, nil
}

func GetTestWalletSigner() (*accounts.WalletSigner, error) {
	testWallet, testWalletManager, err := GetTestWalletManager()
	if err != nil {
		return nil, err
	}

	testAccounts, err := testWallet.Accounts()
	if err != nil {
		return nil, err
	}

	testSigner := accounts.MakeWalletSigner(testAccounts[0].Address, testWalletManager)
	return testSigner, nil
}
