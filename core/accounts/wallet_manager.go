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
	"github.com/dipperin/dipperin-core/common/gtimer"
	"github.com/dipperin/dipperin-core/common/log"
	"github.com/dipperin/dipperin-core/core/accounts/accountsbase"
	"github.com/ethereum/go-ethereum/event"
	"go.uber.org/zap"
	"sync"
	"sync/atomic"
	"time"
)

//wallet event type
type WalletEventType int

const (
	//establish wallet
	//default wallet event channel is 0, so there will be an problem if it is 0
	WalletArrived WalletEventType = 1 + iota

	//open wallet
	WalletOpened

	//remove wallet
	WalletDropped
)

//refresh wallet nonce in the wallet manager timely
const (
	RefreshWalletInfoDuration time.Duration = time.Second * 60
)

//record wallet backend event
type WalletEvent struct {
	Wallet accountsbase.Wallet // Wallet instance arrived or departed
	Type   WalletEventType     // Event type that happened in the system
}

type WalletManager struct {
	Wallets               []accountsbase.Wallet          //wallets
	GetAddressRelatedInfo accountsbase.AddressInfoReader //get nonce and balance of the account
	Event                 chan WalletEvent               //listen wallet event
	HandleResult          chan bool                      //event handle result
	ManagerClose          chan bool                      //listen the manger close
	StartService          chan bool

	serviceStatus    atomic.Value
	startServiceFeed event.Feed
	feed             event.Feed //subscribe managerClose channel
	Lock             sync.RWMutex
}

//new wallet manager
func NewWalletManager(getAddressInfo accountsbase.AddressInfoReader, wallets ...accountsbase.Wallet) (*WalletManager, error) {

	tmpWallets := make([]accountsbase.Wallet, 0)
	for _, tmpWallet := range wallets {
		walletIdentifier, err := tmpWallet.GetWalletIdentifier()
		if err != nil {
			return nil, err
		}
		if walletIdentifier.WalletType != accountsbase.SoftWallet {
			return nil, gerror.ErrNotSupportUsbWallet
		} else {
			tmpWallets = append(tmpWallets, tmpWallet)
		}
	}

	manager := &WalletManager{
		Wallets:               tmpWallets,
		GetAddressRelatedInfo: getAddressInfo,
		Event:                 make(chan WalletEvent, 0),
		HandleResult:          make(chan bool, 0),
		ManagerClose:          make(chan bool, 0),
		StartService:          make(chan bool, 0),
		startServiceFeed:      event.Feed{},
		feed:                  event.Feed{},
		Lock:                  sync.RWMutex{},
	}

	return manager, nil
}

func (manager *WalletManager) SubScribeStartService() <-chan bool {
	manager.startServiceFeed.Subscribe(manager.StartService)
	return manager.StartService
}

func (manager *WalletManager) StartOtherServices() {
	manager.startServiceFeed.Send(true)
}

//listen wallet　event
func (manager *WalletManager) backend() {
	log.DLogger.Info("backend start")
	sub := manager.feed.Subscribe(manager.ManagerClose)
	log.DLogger.Info("backend subscribe ManagerClose")
	for {
		select {
		case walletEvent := <-manager.Event:
			//new wallet event
			manager.Lock.Lock()
			if walletEvent.Type == WalletArrived {
				//add wallet to manager
				manager.add(walletEvent.Wallet)
				//handle end
				manager.HandleResult <- true
			} else if walletEvent.Type == WalletDropped {
				//remove wallet from manager
				manager.remove(walletEvent.Wallet)
				//handle result
				manager.HandleResult <- true
			}
			manager.Lock.Unlock()
		case <-manager.ManagerClose:
			sub.Unsubscribe()
			log.DLogger.Info("Wallet manager backend return")
			return
		}
	}
}

//refresh the account nonce timely in the manager
func (manager *WalletManager) refreshWalletNonce() {
	//subscribe　wallet manager　channelencryptWalletAndWriteFile(
	sub := manager.feed.Subscribe(manager.ManagerClose)

	timeoutHandler := func() {
		//refresh wallet balance
		manager.Lock.Lock()
		for _, wallet := range manager.Wallets {
			wallet.PaddingAddressNonce(manager.GetAddressRelatedInfo)
		}
		manager.Lock.Unlock()
	}
	timer := gtimer.SetPeriodAndRun(timeoutHandler, RefreshWalletInfoDuration)
	defer gtimer.StopWork(timer)

	for {
		select {
		case <-manager.ManagerClose:
			sub.Unsubscribe()
			log.DLogger.Info("refresh Wallet backend return")
			return
		}
	}
}

//add wallet to manager
func (manager *WalletManager) add(wallet accountsbase.Wallet) {
	for _, value := range manager.Wallets {
		identifier, _ := value.GetWalletIdentifier()
		tmpIdentifier, _ := wallet.GetWalletIdentifier()
		if identifier == tmpIdentifier {
			//there is the same wallet in the wallet manager
			return
		}
	}
	manager.Wallets = append(manager.Wallets, wallet)
}

//remove wallet
func (manager *WalletManager) remove(wallet accountsbase.Wallet) {
	for i, value := range manager.Wallets {
		identifier, _ := value.GetWalletIdentifier()
		tmpIdentifier, _ := wallet.GetWalletIdentifier()
		if identifier == tmpIdentifier {
			if i == (len(manager.Wallets) - 1) {
				manager.Wallets = manager.Wallets[:i]
			} else {
				manager.Wallets = append(manager.Wallets[:i], manager.Wallets[i+1:]...)
			}
			log.DLogger.Info("the manager.Wallets is: ", zap.Any("manager.Wallets", manager.Wallets))
		}
	}
}

//list all wallet identifier in the wallet manager
func (manager *WalletManager) ListWalletIdentifier() ([]accountsbase.WalletIdentifier, error) {

	manager.Lock.Lock()
	defer manager.Lock.Unlock()
	identifiers := make([]accountsbase.WalletIdentifier, 0)

	for _, wallet := range manager.Wallets {
		walletIdentifier, err := wallet.GetWalletIdentifier()
		if err != nil {
			return nil, err
		}
		identifiers = append(identifiers, walletIdentifier)
	}

	return identifiers, nil
}

//get wallet according to he wallet id
func (manager *WalletManager) FindWalletFromIdentifier(identifier accountsbase.WalletIdentifier) (accountsbase.Wallet, error) {
	manager.Lock.Lock()
	defer manager.Lock.Unlock()

	for _, wallet := range manager.Wallets {

		walletIdentifier, err := wallet.GetWalletIdentifier()
		if err != nil {
			return nil, err
		}
		if walletIdentifier == identifier {
			return wallet, nil
		}
	}
	return nil, gerror.ErrNotFindWallet
}

//get wallet from account address
func (manager *WalletManager) FindWalletFromAddress(address common.Address) (accountsbase.Wallet, error) {
	manager.Lock.Lock()
	defer manager.Lock.Unlock()

	tmpAccount := accountsbase.Account{Address: address}
	for _, wallet := range manager.Wallets {

		exist, err := wallet.Contains(tmpAccount)
		if err != nil {
			return nil, err
		}
		if exist == true {
			return wallet, nil
		}
	}
	return nil, gerror.ErrNotFindWallet
}

func (manager *WalletManager) GetMainAccount() (accountsbase.Account, error) {
	if !manager.ServiceStatus() {
		return accountsbase.Account{}, gerror.ErrWalletManagerNotRunning
	}
	identifiers, err := manager.ListWalletIdentifier()
	if err != nil {
		return accountsbase.Account{}, err
	}

	if len(identifiers) == 0 {
		return accountsbase.Account{}, gerror.ErrWalletManagerIsEmpty
	}

	wallet, err := manager.FindWalletFromIdentifier(identifiers[0])
	if err != nil {
		return accountsbase.Account{}, err
	}

	account, err := wallet.Accounts()
	if err != nil {
		return accountsbase.Account{}, err
	}

	return account[0], nil
}

func (manager *WalletManager) Start() error {
	go manager.backend()
	go manager.refreshWalletNonce()
	manager.serviceStatus.Store(true)
	return nil
}

//close wallet manager
func (manager *WalletManager) Stop() {
	log.DLogger.Info("WalletManager close")
	log.DLogger.Info("close request lock")
	log.DLogger.Info("close request end")

	manager.Lock.Lock()
	defer manager.Lock.Unlock()

	//close backend
	manager.feed.Send(true)
	//manager.ManagerClose <- true
	log.DLogger.Info("the feed send end")

	//close all wallets in the manager
	for _, wallet := range manager.Wallets {
		wallet.Close()
	}

	//clear all wallet in the manager
	manager.Wallets = []accountsbase.Wallet{}

	close(manager.ManagerClose)
	close(manager.HandleResult)
	close(manager.Event)

	manager.serviceStatus.Store(false)
}

func (manager *WalletManager) ServiceStatus() bool {
	status := manager.serviceStatus.Load()
	if status == nil {
		return false
	}

	return status.(bool)
}
