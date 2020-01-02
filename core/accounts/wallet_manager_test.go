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

package accounts_test

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/log"
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/core/accounts"
	"github.com/dipperin/dipperin-core/core/accounts/soft-wallet"
	"github.com/dipperin/dipperin-core/tests/wallet"
	"github.com/ethereum/go-ethereum/event"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"os"
	"sync/atomic"
	"testing"
	"time"
)

//test new wallet manager
func Test_NewWalletManager(t *testing.T) {
	log.DLogger.Info("Test_NewWalletManager start")
	testWallet, walletManager, err := wallet.GetTestWalletManager()
	assert.NoError(t, err)

	walletManager.Start()

	err = testWallet.Close()
	assert.NoError(t, err)
	walletManager.Stop()

	os.Remove(testWallet.Identifier.Path)
	log.DLogger.Info("Test_NewWalletManager end")
}

func Test_ListWalletIdentifier(t *testing.T) {
	log.DLogger.Info("Test_ListWalletIdentifier start")
	testWallet, walletManager, err := wallet.GetTestWalletManager()
	assert.NoError(t, err)

	identifiers, err := walletManager.ListWalletIdentifier()
	assert.NoError(t, err)

	assert.EqualValues(t, testWallet.Identifier, identifiers[0])

	err = testWallet.Close()
	assert.NoError(t, err)

	os.Remove(testWallet.Identifier.Path)
	log.DLogger.Info("Test_ListWalletIdentifier end")
}

func Test_FindWalletFromName(t *testing.T) {
	log.DLogger.Info("Test_FindWalletFromName start")
	testWallet, walletManager, err := wallet.GetTestWalletManager()
	assert.NoError(t, err)

	wallet, err := walletManager.FindWalletFromIdentifier(testWallet.Identifier)
	assert.NoError(t, err)

	assert.EqualValues(t, testWallet, wallet.(*soft_wallet.SoftWallet))

	err = testWallet.Close()
	assert.NoError(t, err)
	_, err = walletManager.FindWalletFromIdentifier(testWallet.Identifier)
	assert.Equal(t, accounts.ErrWalletNotOpen, err)

	os.Remove(testWallet.Identifier.Path)
	log.DLogger.Info("Test_FindWalletFromName end")
}

func TestWalletManager_FindWalletFromAddress(t *testing.T) {
	log.DLogger.Info("TestWalletManager_FindWalletFromAddress start")
	testWallet, walletManager, err := wallet.GetTestWalletManager()
	assert.NoError(t, err)

	testAccounts, err := testWallet.Accounts()
	assert.NoError(t, err)

	findWallet, err := walletManager.FindWalletFromAddress(testAccounts[0].Address)
	assert.NoError(t, err)
	assert.Equal(t, testWallet, findWallet)

	_, err = walletManager.FindWalletFromAddress(common.Address{})
	assert.Equal(t, accounts.ErrNotFindWallet, err)

	err = testWallet.Close()
	assert.NoError(t, err)
	_, err = walletManager.FindWalletFromAddress(testAccounts[0].Address)
	assert.Equal(t, accounts.ErrWalletNotOpen, err)
	log.DLogger.Info("TestWalletManager_FindWalletFromAddress end")
}

func Test_WalletManagerBackend(t *testing.T) {
	log.DLogger.Info("Test_WalletManagerBackend start")
	testWallet, walletManager, err := wallet.GetTestWalletManager()
	assert.NoError(t, err)

	walletManager.Start()

	WalletName2 := "testSoftWallet2"
	Path2 := util.HomeDir() + "/testSoftWallet2"

	testWallet2, err := wallet.EstablishSoftWallet(Path2, WalletName2)
	assert.NoError(t, err)

	testEvent := accounts.WalletEvent{
		Wallet: testWallet2,
		Type:   accounts.WalletArrived,
	}

	walletManager.Event <- testEvent

	select {
	case result := <-walletManager.HandleResult:
		assert.EqualValues(t, true, result)
	}

	assert.EqualValues(t, 2, len(walletManager.Wallets))

	testEvent.Type = accounts.WalletDropped
	walletManager.Event <- testEvent

	select {
	case result := <-walletManager.HandleResult:
		assert.EqualValues(t, true, result)
	}

	assert.EqualValues(t, 1, len(walletManager.Wallets))

	err = testWallet.Close()
	assert.NoError(t, err)

	err = testWallet2.Close()
	assert.NoError(t, err)
	os.Remove(testWallet.Identifier.Path)
	os.Remove(testWallet2.Identifier.Path)

	log.DLogger.Info("Test_WalletManagerBackend end")
}

func Test_ChannelTransport(t *testing.T) {
	type testAtomic struct {
		value atomic.Value
	}

	testData := testAtomic{
		value: atomic.Value{},
	}

	testData.value.Store(32)

	log.DLogger.Info("the atomic value is:", zap.Int("value", testData.value.Load().(int)))

	testChan := make(chan testAtomic)
	var feed event.Feed

	go func() {
		sub := feed.Subscribe(testChan)
		defer sub.Unsubscribe()
		log.DLogger.Info("subscribe success")
		for {
			select {
			case readData := <-testChan:
				log.DLogger.Info("the read atomic value is:", zap.Int("value", readData.value.Load().(int)))
				return
			}
		}
	}()

	time.Sleep(2)
	ret := feed.Send(testData)
	log.DLogger.Info("the ret is:", zap.Int("ret", ret))
	time.Sleep(2)
	//testChan <- testData
}
