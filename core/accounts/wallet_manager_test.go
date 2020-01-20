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

package accounts

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/gerror"
	"github.com/dipperin/dipperin-core/common/log"
	"github.com/dipperin/dipperin-core/core/accounts/accountsbase"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/ethereum/go-ethereum/event"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"os"
	"sync/atomic"
	"testing"
	"time"
)

func getWalletAndWalletManager(t *testing.T) (*gomock.Controller, *accountsbase.MockWallet, *WalletManager) {
	ctrl := gomock.NewController(t)
	infoReader := accountsbase.NewMockAddressInfoReader(ctrl)
	wallet := accountsbase.NewMockWallet(ctrl)
	wallet.EXPECT().GetWalletIdentifier().Return(accountsbase.WalletIdentifier{WalletType:accountsbase.SoftWallet}, nil)
	wallet.EXPECT().Close().Return(nil).AnyTimes()
	walletManager,err := NewWalletManager(infoReader, wallet)
	assert.NoError(t, err)
	return ctrl, wallet, walletManager
}

//test new wallet manager
func Test_NewWalletManager(t *testing.T) {
	ctrl, wallet, walletManager := getWalletAndWalletManager(t)

	ctrl.Finish()

	walletManager.Start()

	err := wallet.Close()
	assert.NoError(t, err)
	walletManager.Stop()

	//os.Remove(wa.Identifier.Path)
	//log.DLogger.Info("Test_NewWalletManager end")
}

func Test_ListWalletIdentifier(t *testing.T) {
	ctrl, wallet, walletManager := getWalletAndWalletManager(t)
	ctrl.Finish()

	wallet.EXPECT().GetWalletIdentifier().Return(accountsbase.WalletIdentifier{WalletType:accountsbase.SoftWallet}, nil).AnyTimes()

	identifiers, err := walletManager.ListWalletIdentifier()
	assert.NoError(t, err)
	walletIndentifier, err := wallet.GetWalletIdentifier()
	assert.NoError(t, err)
	assert.EqualValues(t, walletIndentifier, identifiers[0])

	err = wallet.Close()
	assert.NoError(t, err)

	os.Remove(walletIndentifier.Path)
	log.DLogger.Info("Test_ListWalletIdentifier end")
}


// todo
func Test_FindWalletFromName(t *testing.T) {
	t.Skip()
	ctrl, wallet, walletManager := getWalletAndWalletManager(t)
	ctrl.Finish()
	wallet.EXPECT().GetWalletIdentifier().Return(accountsbase.WalletIdentifier{WalletType:accountsbase.SoftWallet}, nil).AnyTimes()
	walletIdentifier, err := wallet.GetWalletIdentifier()
	assert.NoError(t, err)
	wt, err := walletManager.FindWalletFromIdentifier(walletIdentifier)
	assert.NoError(t, err)

	assert.EqualValues(t, wallet, wt)

	err = wallet.Close()
	assert.NoError(t, err)
	_, err = walletManager.FindWalletFromIdentifier(walletIdentifier)
	assert.Equal(t, gerror.ErrWalletNotOpen, err)

	//os.Remove(testWallet.Identifier.Path)
	log.DLogger.Info("Test_FindWalletFromName end")
}

func TestWalletManager_FindWalletFromAddress(t *testing.T) {
	ctrl, wallet, walletManager := getWalletAndWalletManager(t)
	ctrl.Finish()
	wallet.EXPECT().Accounts().Return([]accountsbase.Account{accountsbase.Account{Address:model.AliceAddr},},nil)
	testAccounts, err := wallet.Accounts()
	assert.NoError(t, err)

	testCases := []struct{
		name string
		given func () common.Address
		expect error
		expectResult accountsbase.Wallet
	}{
		{
			name:"FindWalletFromAddressRight",
			given: func() common.Address {
				wallet.EXPECT().Contains(accountsbase.Account{Address:model.AliceAddr}).Return(true,nil)
				return testAccounts[0].Address
			},
			expect:nil,
			expectResult:wallet,
		},
		{
			name:"ErrNotFindWallet",
			given: func() common.Address {
				wallet.EXPECT().Contains(accountsbase.Account{Address:common.Address{}}).Return(false,gerror.ErrNotFindWallet)
				return common.Address{}
			},
			expect:gerror.ErrNotFindWallet,
		},
		{
			name:"ErrWalletNotOpen",
			given: func() common.Address {
				wallet.EXPECT().Contains(accountsbase.Account{Address:model.AliceAddr}).Return(false,gerror.ErrWalletNotOpen)
				return testAccounts[0].Address
			},
			expect:gerror.ErrWalletNotOpen,
		},
	}

	for _,tc := range testCases{
		address := tc.given()
		findWallet, err := walletManager.FindWalletFromAddress(address)
		if err != nil {
			assert.Equal(t, tc.expect, err)
		}else {
			assert.Equal(t, wallet, findWallet)
		}
	}

}

// todo
func Test_WalletManagerBackend(t *testing.T) {
	t.Skip()
	ctrl, wallet, walletManager := getWalletAndWalletManager(t)
	ctrl.Finish()

	walletManager.Start()

	testEvent := WalletEvent{
		Wallet: wallet,
		Type:   WalletArrived,
	}

	walletManager.Event <- testEvent

	go func() {walletManager.HandleResult <- true}()

	select {
	case result := <-walletManager.HandleResult:
		assert.EqualValues(t, true, result)
	}



	assert.EqualValues(t, 2, len(walletManager.Wallets))

	testEvent.Type = WalletDropped
	walletManager.Event <- testEvent

	go func() {walletManager.HandleResult <- true}()

	select {
	case result := <-walletManager.HandleResult:
		assert.EqualValues(t, true, result)
	}

	assert.EqualValues(t, 1, len(walletManager.Wallets))

	err := wallet.Close()
	assert.NoError(t, err)

	//err = testWallet2.Close()
	//assert.NoError(t, err)
	//os.Remove(testWallet.Identifier.Path)
	//os.Remove(testWallet2.Identifier.Path)
	//
	//log.DLogger.Info("Test_WalletManagerBackend end")
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
