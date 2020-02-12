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
	"github.com/dipperin/dipperin-core/common/util"
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

func getWalletAndWalletManager(t *testing.T) (*gomock.Controller, *MockWallet, *WalletManager) {
	ctrl := gomock.NewController(t)
	infoReader := NewMockAddressInfoReader(ctrl)
	wallet := NewMockWallet(ctrl)
	wallet.EXPECT().GetWalletIdentifier().Return(accountsbase.WalletIdentifier{WalletType:accountsbase.SoftWallet}, nil).AnyTimes()
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

func TestWalletManager_SubScribeStartService(t *testing.T) {
	ctrl, _, walletManager := getWalletAndWalletManager(t)
	ctrl.Finish()

	event := walletManager.SubScribeStartService()
	go func() {
		select {
		case result := <-event:
			assert.Equal(t, true, result)
		}
	}()

	walletManager.StartOtherServices()
}


func TestWalletManager_backend(t *testing.T)  {
	ctrl, wallet, walletManager := getWalletAndWalletManager(t)
	ctrl.Finish()

	go func() {
		walletManager.backend()
	}()

	testCases := []struct{
		name string
		given func() WalletEvent
		expect bool
	} {
		{
			name:"WalletArrived",
			given: func() WalletEvent {
				testEvent := WalletEvent{
					Wallet: wallet,
					Type:   WalletArrived,
				}
				return testEvent
			},
			expect:true,
		},
		{
			name:"WalletDropped",
			given: func() WalletEvent {
				testEvent := WalletEvent{
					Wallet: wallet,
					Type:   WalletDropped,
				}
				return testEvent
			},
			expect:true,
		},
	}

	for _, tc := range testCases {
		t.Log(tc.name)
		testEvent := tc.given()
		wallet.EXPECT().GetWalletIdentifier().Return(accountsbase.WalletIdentifier{WalletType:accountsbase.SoftWallet}, nil).AnyTimes()
		walletManager.Event <- testEvent
		result := <- walletManager.HandleResult
		assert.Equal(t, true, result)
	}
	walletManager.ManagerClose <- true

}

func TestWalletManager_refreshWalletNonce(t *testing.T) {
	ctrl, wallet, walletManager := getWalletAndWalletManager(t)
	ctrl.Finish()

	go func() {
		walletManager.refreshWalletNonce()
	}()

	wallet.EXPECT().PaddingAddressNonce(walletManager.GetAddressRelatedInfo).Return(nil).AnyTimes()

	walletManager.ManagerClose <- true
}

func TestWalletManager_add(t *testing.T) {
	ctrl, wallet, walletManager := getWalletAndWalletManager(t)
	ctrl.Finish()

	testCases := []struct{
		name string
		given func () *MockWallet
		expect int
	}{
		{
			name:"theSameWallet",
			given: func() *MockWallet {
				return wallet
			},
			expect:len(walletManager.Wallets),
		},
		{
			name:"diffWallet",
			given: func() *MockWallet {
				wallet := NewMockWallet(ctrl)
				wallet.EXPECT().GetWalletIdentifier().Return(accountsbase.WalletIdentifier{WalletType:accountsbase.SoftWallet, Path:util.HomeDir()}, nil).AnyTimes()
				return wallet
			},
			expect:len(walletManager.Wallets) + 1,
		},
	}

	for _, tc := range testCases {
		tempWallet := tc.given()
		walletManager.add(tempWallet)
		assert.Equal(t, tc.expect, len(walletManager.Wallets))
	}

}

func TestWalletManager_remove(t *testing.T) {
	ctrl, wallet, walletManager := getWalletAndWalletManager(t)
	ctrl.Finish()

	testCases := []struct{
		name string
		given func () *MockWallet
		expect int
	}{
		{
			name:"diffWallet",
			given: func() *MockWallet {
				wallet := NewMockWallet(ctrl)
				wallet.EXPECT().GetWalletIdentifier().Return(accountsbase.WalletIdentifier{WalletType:accountsbase.SoftWallet, Path:util.HomeDir()}, nil).AnyTimes()
				return wallet
			},
			expect:len(walletManager.Wallets),
		},
		{
			name:"theSameWallet",
			given: func() *MockWallet {
				return wallet
			},
			expect:len(walletManager.Wallets) -1 ,
		},

	}

	for _, tc := range testCases {
		tempWallet := tc.given()
		walletManager.remove(tempWallet)
		assert.Equal(t, tc.expect, len(walletManager.Wallets))
	}
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
}


func TestWalletManager_FindWalletFromIdentifier(t *testing.T) {
	ctrl, wallet, walletManager := getWalletAndWalletManager(t)
	ctrl.Finish()

	testCases := []struct{
		name string
		given func () accountsbase.WalletIdentifier
		expect error
	}{
		{
			name:"diffWallet",
			given: func() accountsbase.WalletIdentifier {
				wallet := NewMockWallet(ctrl)
				wallet.EXPECT().GetWalletIdentifier().Return(accountsbase.WalletIdentifier{WalletType:accountsbase.SoftWallet, Path:util.HomeDir()}, nil).AnyTimes()
				return accountsbase.WalletIdentifier{WalletType:accountsbase.SoftWallet, Path:util.HomeDir()}
			},
			expect:gerror.ErrNotFindWallet,
		},
		{
			name:"theSameWallet",
			given: func() accountsbase.WalletIdentifier {
				return accountsbase.WalletIdentifier{WalletType:accountsbase.SoftWallet}
			},
			expect:nil,
		},

	}

	for _, tc := range testCases {
		identifier := tc.given()
		walletTemp, err := walletManager.FindWalletFromIdentifier(identifier)
		if err != nil {
			assert.Equal(t, tc.expect, err)
		} else {
			assert.Equal(t, wallet, walletTemp)
		}
	}
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

func TestWalletManager_GetMainAccount(t *testing.T) {

	testCases := []struct{
		name string
		given func() (*WalletManager, accountsbase.Account)
		expect error
	}{
		{
			name:"ErrWalletManagerNotRunning",
			given: func() (*WalletManager, accountsbase.Account) {
				_, _, walletManager := getWalletAndWalletManager(t)
				walletManager.serviceStatus.Store(false)
				return walletManager,accountsbase.Account{}
			},
			expect:gerror.ErrWalletManagerNotRunning,
		},
		{
			name:"ErrWalletManagerIsEmpty",
			given: func() (*WalletManager, accountsbase.Account) {
				_, _, walletManager := getWalletAndWalletManager(t)
				walletManager.Wallets = []accountsbase.Wallet{}
				walletManager.serviceStatus.Store(true)
				return walletManager,accountsbase.Account{}
			},
			expect:gerror.ErrWalletManagerIsEmpty,
		},
		{
			name:"GetMainAccountRight",
			given: func() (*WalletManager, accountsbase.Account) {
				_, wallet, walletManager := getWalletAndWalletManager(t)
				walletManager.serviceStatus.Store(true)
				wallet.EXPECT().Accounts().Return([]accountsbase.Account{accountsbase.Account{Address:common.HexToAddress("0x000")}},nil).AnyTimes()
				accounts ,err := walletManager.Wallets[0].Accounts()
				assert.NoError(t, err)
				return walletManager,accounts[0]
			},
			expect:nil,
		},
	}

	for _, tc := range testCases {
		manager, account := tc.given()
		tempAccount, err := manager.GetMainAccount()
		if err != nil {
			assert.Equal(t, tc.expect, err)
		} else {
			assert.Equal(t, account, tempAccount)
		}
	}

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
