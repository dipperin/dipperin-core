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

package service

import (
	"context"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/g-error"
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/core/accounts"
	"github.com/dipperin/dipperin-core/core/accounts/soft-wallet"
	"github.com/dipperin/dipperin-core/core/chain"
	"github.com/dipperin/dipperin-core/core/chain-communication"
	"github.com/dipperin/dipperin-core/core/chain-config"
	contract2 "github.com/dipperin/dipperin-core/core/contract"
	"github.com/dipperin/dipperin-core/core/economy-model"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/core/vm/common/utils"
	"github.com/dipperin/dipperin-core/tests"
	"github.com/dipperin/dipperin-core/tests/g-testData"
	"github.com/dipperin/dipperin-core/third-party/crypto"
	"github.com/dipperin/dipperin-core/third-party/p2p"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"math/big"
	"os"
	"testing"
	"time"
)

var testFee = economy_model.GetMinimumTxFee(1000)

func TestMercuryFullChainService_NormalPM(t *testing.T) {
	config := &DipperinConfig{NormalPm: fakePeerManager{}}
	service := MakeFullChainService(config)

	height := service.RemoteHeight()
	assert.Equal(t, uint64(0), height)

	result := service.GetSyncStatus()
	assert.True(t, result)
}

func TestMercuryFullChainService_AddAccount(t *testing.T) {
	manager := createWalletManager(t)
	defer os.RemoveAll(util.HomeDir() + testPath)
	config := &DipperinConfig{
		NodeConf:      fakeNodeConfig{},
		WalletManager: manager,
	}
	service := MakeFullChainService(config)
	identifier := createWalletIdentifier()

	account, err := service.AddAccount(*identifier, "m/44'/709394'/0'/0")
	assert.NoError(t, err)
	assert.NotNil(t, account)

	account, err = service.AddAccount(*identifier, "")
	assert.NoError(t, err)
	assert.NotNil(t, account)

	// tmpWallet.Derive error
	account, err = service.AddAccount(*identifier, "123")
	assert.Equal(t, accounts.ErrInvalidDerivedPath, err)
	assert.Equal(t, accounts.Account{}, account)

	// ParseDerivationPath error
	account, err = service.AddAccount(*identifier, "m")
	assert.Equal(t, "empty derivation path", err.Error())
	assert.Equal(t, accounts.Account{}, account)

	// FindWalletFromIdentifier error
	identifier.Path = "t"
	account, err = service.AddAccount(*identifier, "")
	assert.Equal(t, accounts.ErrNotFindWallet, err)
	assert.Equal(t, accounts.Account{}, account)

	// checkWalletIdentifier error
	identifier.WalletType = 123
	account, err = service.AddAccount(*identifier, "")
	assert.Equal(t, "wallet type error", err.Error())
	assert.Equal(t, accounts.Account{}, account)
}

func TestMercuryFullChainService_AddPeer(t *testing.T) {
	config := &DipperinConfig{}
	service := MakeFullChainService(config)

	err := service.AddPeer("")
	assert.Error(t, err)
	err = service.AddTrustedPeer("")
	assert.Error(t, err)
	peer, err := service.Peers()
	assert.Error(t, err)
	assert.Nil(t, peer)

	config = &DipperinConfig{P2PServer: &p2p.Server{}}
	service = MakeFullChainService(config)

	err = service.AddPeer(url_wrong)
	assert.Error(t, err)
	err = service.AddTrustedPeer(url_wrong)
	assert.Error(t, err)
}

func TestMercuryFullChainService_RemovePeer(t *testing.T) {
	config := &DipperinConfig{}
	service := MakeFullChainService(config)

	err := service.RemovePeer("")
	assert.Error(t, err)
	err = service.RemoveTrustedPeer("")
	assert.Error(t, err)

	config = &DipperinConfig{P2PServer: &p2p.Server{}}
	service = MakeFullChainService(config)

	err = service.RemovePeer(url_wrong)
	assert.Error(t, err)
	err = service.RemoveTrustedPeer(url_wrong)
	assert.Error(t, err)
}

func TestMercuryFullChainService_GetCurrentConnectPeers(t *testing.T) {
	config := &DipperinConfig{}
	service := MakeFullChainService(config)

	peers := service.GetCurrentConnectPeers()
	assert.Equal(t, make(map[string]common.Address, 0), peers)
}

func TestMercuryFullChainService_CurrentBalance(t *testing.T) {
	csChain := createCsChain(nil)
	config := &DipperinConfig{ChainReader: csChain}
	service := MakeFullChainService(config)
	verifiers := service.ChainReader.GetCurrVerifiers()

	assert.NotNil(t, service.CurrentBalance(verifiers[0]))
	assert.Nil(t, service.CurrentBalance(common.HexToAddress("123")))

	block := csChain.CurrentBlock()
	csChain.ChainDB.DeleteBlock(block.Hash(), block.Number())
	assert.Nil(t, service.CurrentBalance(verifiers[0]))
}

func TestMercuryFullChainService_CurrentBlock(t *testing.T) {
	csChain := createCsChain(nil)
	config := &DipperinConfig{ChainReader: csChain}
	service := MakeFullChainService(config)
	assert.Equal(t, uint64(0), service.CurrentBlock().Number())
}

func TestMercuryFullChainService_CurrentStake(t *testing.T) {
	csChain := createCsChain(nil)
	config := &DipperinConfig{ChainReader: csChain}
	service := MakeFullChainService(config)
	verifiers := service.ChainReader.GetCurrVerifiers()

	assert.Equal(t, big.NewInt(0), service.CurrentStake(verifiers[0]))
	assert.NotNil(t, service.CurrentStake(verifiers[0]))
	assert.Nil(t, service.CurrentStake(common.HexToAddress("123")))

	block := csChain.CurrentBlock()
	csChain.ChainDB.DeleteBlock(block.Hash(), block.Number())
	assert.Nil(t, service.CurrentStake(verifiers[0]))

}

func TestMercuryFullChainService_CurrentReputation(t *testing.T) {
	csChain := createCsChain(nil)
	config := &DipperinConfig{
		ChainReader:        csChain,
		PriorityCalculator: model.TestCalculator{},
	}
	service := MakeFullChainService(config)
	verifiers := service.ChainReader.GetCurrVerifiers()

	reputation, err := service.CurrentReputation(verifiers[0])
	assert.Equal(t, "stake not sufficient", err.Error())
	assert.Equal(t, uint64(0), reputation)

	// current state error
	block := csChain.CurrentBlock()
	csChain.ChainDB.DeleteBlock(block.Hash(), block.Number())
	service = MakeFullChainService(config)
	reputation, err = service.CurrentReputation(verifiers[0])
	assert.Equal(t, "current block is nil", err.Error())
	assert.Equal(t, uint64(0), reputation)
}

func TestMercuryFullChainService_CurrentElectPriority(t *testing.T) {
	manager := createWalletManager(t)
	defer os.Remove(util.HomeDir() + testPath)
	account, err := manager.Wallets[0].Accounts()
	assert.NoError(t, err)

	address := account[0].Address
	pk, err := manager.Wallets[0].GetSKFromAddress(address)
	testAccount := tests.NewAccount(pk, address)
	testAccounts := []tests.Account{*testAccount}
	config := &DipperinConfig{
		NodeConf:           fakeNodeConfig{},
		WalletManager:      manager,
		ChainReader:        createCsChain(testAccounts),
		PriorityCalculator: model.TestCalculator{},
	}

	service := MakeFullChainService(config)
	verifiers := service.ChainReader.GetCurrVerifiers()

	reputation, err := service.CurrentElectPriority(verifiers[0])
	assert.Equal(t, "stake not sufficient", err.Error())
	assert.Equal(t, uint64(0), reputation)

	// getLuckProof error
	reputation, err = service.CurrentElectPriority(aliceAddr)
	assert.Equal(t, accounts.ErrNotFindWallet, err)
	assert.Equal(t, uint64(0), reputation)
}

func TestMercuryFullChainService_GetAddressNonceFromWallet(t *testing.T) {
	manager := createWalletManager(t)
	defer os.Remove(util.HomeDir() + testPath)
	config := &DipperinConfig{
		NodeConf:      fakeNodeConfig{},
		WalletManager: manager,
	}
	service := MakeFullChainService(config)

	account, err := config.WalletManager.Wallets[0].Accounts()
	address := account[0].Address
	nonce, err := service.GetAddressNonceFromWallet(address)
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), nonce)

	nonce, err = service.GetAddressNonceFromWallet(common.HexToAddress("123"))
	assert.Equal(t, accounts.ErrNotFindWallet, err)
	assert.Equal(t, uint64(0), nonce)
}

func TestMercuryFullChainService_GetTransactionNonce(t *testing.T) {
	csChain := createCsChain(nil)
	config := &DipperinConfig{ChainReader: csChain}
	service := MakeFullChainService(config)

	nonce, err := service.GetTransactionNonce(chain.VerifierAddress[0])
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), nonce)

	block := csChain.CurrentBlock()
	csChain.ChainDB.DeleteBlock(block.Hash(), block.Number())
	config.ChainReader = csChain
	service.DipperinConfig = config
	nonce, err = service.GetTransactionNonce(chain.VerifierAddress[0])
	assert.Equal(t, "current block is nil", err.Error())
	assert.Equal(t, uint64(0), nonce)
}

func TestMercuryFullChainService_GetBlockBody(t *testing.T) {
	csChain := createCsChain(nil)
	config := &DipperinConfig{ChainReader: csChain}
	service := MakeFullChainService(config)

	hash := service.CurrentBlock().Hash()
	assert.NotNil(t, service.GetBlockBody(hash))
}

func TestMercuryFullChainService_GetBlockByHash(t *testing.T) {
	csChain := createCsChain(nil)
	config := &DipperinConfig{ChainReader: csChain}
	service := MakeFullChainService(config)

	hash := service.CurrentBlock().Hash()
	block, err := service.GetBlockByHash(hash)
	assert.NoError(t, err)
	assert.Equal(t, hash, block.Hash())
}

func TestMercuryFullChainService_GetBlockByNumber(t *testing.T) {
	csChain := createCsChain(nil)
	config := &DipperinConfig{ChainReader: csChain}
	service := MakeFullChainService(config)

	block, err := service.GetBlockByNumber(0)
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), block.Number())
}

func TestMercuryFullChainService_GetBlockNumber(t *testing.T) {
	csChain := createCsChain(nil)
	config := &DipperinConfig{ChainReader: csChain}
	service := MakeFullChainService(config)

	hash := service.CurrentBlock().Hash()
	number := uint64(0)
	assert.Equal(t, &number, service.GetBlockNumber(hash))
}

func TestMercuryFullChainService_GetChainConfig(t *testing.T) {
	config := &DipperinConfig{ChainConfig: *chain_config.GetChainConfig()}
	service := MakeFullChainService(config)

	assert.NotNil(t, service.GetChainConfig())
}

func TestMercuryFullChainService_GetBlockYear(t *testing.T) {
	csChain := createCsChain(nil)
	config := &DipperinConfig{ChainReader: csChain}
	service := MakeFullChainService(config)

	year, err := service.GetBlockYear(0)
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), year)
}

func TestMercuryFullChainService_GetContract(t *testing.T) {
	csChain := createCsChain(nil)
	tx, contractID := createERC20()
	block := createBlock(csChain, []*model.Transaction{tx}, nil)
	votes := createVerifiersVotes(block, csChain.ChainConfig.VerifierNumber)
	err := csChain.SaveBftBlock(block, votes)
	assert.NoError(t, err)

	config := &DipperinConfig{ChainReader: csChain}
	service := MakeFullChainService(config)

	contract, err := service.GetContract(contractID)
	assert.Error(t, err)
	assert.Nil(t, contract)
}

func TestMercuryFullChainService_GetContractInfo(t *testing.T) {
	csChain := createCsChain(nil)
	config := &DipperinConfig{ChainReader: csChain}
	service := MakeFullChainService(config)

	contract, err := service.GetContractInfo(&contract2.ExtraDataForContract{})
	assert.Error(t, err)
	assert.Nil(t, contract)
}

func TestMercuryFullChainService_GetCurVerifiers(t *testing.T) {
	csChain := createCsChain(nil)
	config := &DipperinConfig{ChainReader: csChain}
	service := MakeFullChainService(config)

	verifiers := service.GetCurVerifiers()
	assert.Equal(t, chain.VerifierAddress[0], verifiers[0])
}

func TestMercuryFullChainService_GetNextVerifiers(t *testing.T) {
	csChain := createCsChain(nil)
	config := &DipperinConfig{ChainReader: csChain}
	service := MakeFullChainService(config)

	verifiers := service.GetNextVerifiers()
	assert.Equal(t, chain.VerifierAddress[0], verifiers[0])
}

func TestMercuryFullChainService_GetVerifiers(t *testing.T) {
	csChain := createCsChain(nil)
	config := &DipperinConfig{ChainReader: csChain}
	service := MakeFullChainService(config)

	verifiers := service.GetVerifiers(0)
	assert.Equal(t, chain.VerifierAddress[0], verifiers[0])
}

func TestMercuryFullChainService_GetGenesis(t *testing.T) {
	csChain := createCsChain(nil)
	config := &DipperinConfig{ChainReader: csChain}
	service := MakeFullChainService(config)

	genesis, err := service.GetGenesis()
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), genesis.Number())
}

func TestMercuryFullChainService_GetSlot(t *testing.T) {
	csChain := createCsChain(nil)
	config := &DipperinConfig{ChainReader: csChain}
	service := MakeFullChainService(config)

	slot := service.GetSlot(service.CurrentBlock())
	expect := uint64(0)
	assert.Equal(t, &expect, slot)
}

func TestMercuryFullChainService_ListWallet(t *testing.T) {
	manager := createWalletManager(t)
	defer os.Remove(util.HomeDir() + testPath)
	config := &DipperinConfig{
		NodeConf:      fakeNodeConfig{},
		WalletManager: manager,
	}
	service := MakeFullChainService(config)

	wallet, err := service.ListWallet()
	assert.NoError(t, err)
	assert.NotNil(t, wallet)

	// close wallet
	go func() {
		<-service.WalletManager.Event
		service.WalletManager.HandleResult <- true
	}()
	identifier := createWalletIdentifier()
	err = service.CloseWallet(*identifier)
	assert.NoError(t, err)

	// ListWallet error
	wallet, err = service.ListWallet()
	assert.Equal(t, accounts.ErrWalletNotOpen, err)
	assert.Equal(t, []accounts.WalletIdentifier{}, wallet)
}

func TestMercuryFullChainService_ListWalletAccount(t *testing.T) {
	manager := createWalletManager(t)
	defer os.Remove(util.HomeDir() + testPath)
	config := &DipperinConfig{
		NodeConf:      fakeNodeConfig{},
		WalletManager: manager,
	}
	service := MakeFullChainService(config)

	identifier := createWalletIdentifier()
	wallet, err := service.ListWalletAccount(*identifier)
	assert.NoError(t, err)
	assert.NotNil(t, wallet)

	// FindWalletFromIdentifier error
	identifier.Path = "t"
	wallet, err = service.ListWalletAccount(*identifier)
	assert.Equal(t, accounts.ErrNotFindWallet, err)
	assert.Equal(t, []accounts.Account{}, wallet)

	// checkWalletIdentifier error
	identifier.WalletType = 123
	wallet, err = service.ListWalletAccount(*identifier)
	assert.Equal(t, "wallet type error", err.Error())
	assert.Equal(t, []accounts.Account{}, wallet)
}

func TestMakeFullChainService_checkWalletIdentifier(t *testing.T) {
	config := &DipperinConfig{NodeConf: fakeNodeConfig{}}
	service := MakeFullChainService(config)

	identifier := &accounts.WalletIdentifier{
		WalletType: accounts.LedgerWallet,
	}

	err := service.checkWalletIdentifier(identifier)
	assert.Error(t, err)

	identifier = &accounts.WalletIdentifier{
		WalletType: accounts.SoftWallet,
	}

	err = service.checkWalletIdentifier(identifier)
	assert.NoError(t, err)
}

func TestMercuryFullChainService_GetVerifierReward(t *testing.T) {
	csChain := createCsChain(nil)
	insertBlockToChain(t, csChain, 1)

	config := &DipperinConfig{ChainReader: csChain}
	service := MakeFullChainService(config)

	reward, err := service.GetVerifierDIPReward(0)
	assert.Error(t, err)
	assert.NotNil(t, reward)

	reward, err = service.GetVerifierEDIPReward(0, 2)
	assert.Error(t, err)
	assert.NotNil(t, reward)

	reward, err = service.GetVerifierEDIPReward(1, 2)
	assert.NoError(t, err)
	assert.NotNil(t, reward)
}

func TestMercuryFullChainService_VerifierStatus(t *testing.T) {
	csChain := createCsChain(nil)
	config := DipperinConfig{
		ChainReader:        csChain,
		PriorityCalculator: model.TestCalculator{},
	}
	service := MakeFullChainService(&config)

	verifiers := chain.VerifierAddress
	state, stake, balance, reputation, isCurrent, err := service.VerifierStatus(verifiers[0])
	assert.NoError(t, err)
	assert.Equal(t, "Not Registered", state)
	assert.Equal(t, big.NewInt(0), stake)
	assert.Equal(t, big.NewInt(9999000000000), balance)
	assert.Equal(t, uint64(0), reputation)
	assert.Equal(t, true, isCurrent)
}

func TestMercuryFullChainService_GetBlockDiffVerifierInfo(t *testing.T) {
	csChain := createCsChain(nil)

	insertBlockToChain(t, csChain, 2)

	config := DipperinConfig{ChainReader: csChain}
	service := MakeFullChainService(&config)

	info, err := service.GetBlockDiffVerifierInfo(0)
	assert.Error(t, err)
	assert.Equal(t, map[economy_model.VerifierType][]common.Address{}, info)

	info, err = service.GetBlockDiffVerifierInfo(2)
	assert.NoError(t, err)
	assert.NotNil(t, info)
}

func TestMercuryFullChainService_SetBftSigner(t *testing.T) {
	manager := createWalletManager(t)
	defer os.Remove(util.HomeDir() + testPath)
	config := &DipperinConfig{
		NodeConf:      fakeNodeConfig{},
		WalletManager: manager,
	}
	service := MakeFullChainService(config)

	account, err := config.WalletManager.Wallets[0].Accounts()
	assert.NoError(t, err)
	address := account[0].Address
	signer := accounts.MakeWalletSigner(address, config.WalletManager)
	config.MsgSigner = signer

	err = service.SetBftSigner(address)
	assert.NoError(t, err)
}

func TestMercuryFullChainService_OpenWallet(t *testing.T) {
	manager := createWalletManager(t)
	defer os.Remove(util.HomeDir() + testPath)
	config := &DipperinConfig{
		NodeConf:      fakeNodeConfig{},
		WalletManager: manager,
	}
	service := MakeFullChainService(config)

	go func() {
		<-service.WalletManager.Event
		service.WalletManager.HandleResult <- true
	}()

	err := service.OpenWallet(*createWalletIdentifier(), "123")
	assert.NoError(t, err)

	identifier := &accounts.WalletIdentifier{WalletType: 123}
	err = service.OpenWallet(*identifier, "123")
	assert.Equal(t, "wallet type error", err.Error())

	identifier = &accounts.WalletIdentifier{
		WalletType: accounts.SoftWallet,
		Path:       "t",
		WalletName: "name",
	}
	err = service.OpenWallet(*identifier, "123")
	assert.Equal(t, accounts.ErrWalletPathError, err)
}

type fakeMsgSigner struct{ addr common.Address }

func (f *fakeMsgSigner) SetBaseAddress(address common.Address) {}

func (f *fakeMsgSigner) GetAddress() common.Address {
	return f.addr
}

func TestMercuryFullChainService_CloseWallet(t *testing.T) {
	manager := createWalletManager(t)
	config := &DipperinConfig{
		NodeConf:      fakeNodeConfig{nodeType: chain_config.NodeTypeOfMineMaster},
		WalletManager: manager,
		MsgSigner:     &fakeMsgSigner{},
	}
	service := MakeFullChainService(config)

	// No error
	go func() {
		<-service.WalletManager.Event
		service.WalletManager.HandleResult <- true
	}()
	identifier := createWalletIdentifier()
	err := service.CloseWallet(*identifier)
	assert.NoError(t, err)
	assert.NoError(t, os.RemoveAll(util.HomeDir()+testPath))

	// wallet contains coinbase
	manager = createWalletManager(t)
	defer os.RemoveAll(util.HomeDir() + testPath)
	account, err := manager.Wallets[0].Accounts()
	assert.NoError(t, err)
	config.WalletManager = manager
	config.MsgSigner = &fakeMsgSigner{addr: account[0].Address}
	service = MakeFullChainService(config)
	err = service.CloseWallet(*identifier)
	assert.Equal(t, "this wallet contains coinbase, can not close", err.Error())

	// FindWalletFromIdentifier error
	identifier.WalletName = "name"
	err = service.CloseWallet(*identifier)
	assert.Equal(t, accounts.ErrNotFindWallet, err)

	// wallet type error
	identifier.WalletType = 123
	err = service.CloseWallet(*identifier)
	assert.Equal(t, "wallet type error", err.Error())
}

func TestMercuryFullChainService_RestoreWallet(t *testing.T) {
	wallet, err := soft_wallet.NewSoftWallet()
	assert.NoError(t, err)

	memory, err := wallet.Establish(util.HomeDir()+testPath, "testSoftWallet", "123", "")
	assert.NoError(t, err)

	manager, err := accounts.NewWalletManager(&fakeGetAccountInfo{}, wallet)
	assert.NoError(t, err)

	config := &DipperinConfig{
		NodeConf:      fakeNodeConfig{},
		WalletManager: manager,
		ChainReader:   createCsChain(nil),
	}
	service := MakeFullChainService(config)
	assert.NoError(t, os.RemoveAll(util.HomeDir()+testPath))

	// No error
	go func() {
		<-service.WalletManager.Event
		service.WalletManager.HandleResult <- true
		<-service.WalletManager.Event
		service.WalletManager.HandleResult <- true
	}()

	err = service.RestoreWallet(*createWalletIdentifier(), "123", "", memory)
	assert.NoError(t, err)
	assert.NoError(t, os.RemoveAll(util.HomeDir()+testPath))
}

func TestMercuryFullChainService_RestoreWallet_Error(t *testing.T) {
	manager := createWalletManager(t)
	defer os.RemoveAll(util.HomeDir() + testPath)
	config := &DipperinConfig{
		NodeConf:      fakeNodeConfig{nodeType: chain_config.NodeTypeOfMineMaster},
		WalletManager: manager,
		MsgSigner:     &fakeMsgSigner{},
	}
	service := MakeFullChainService(config)

	identifier := *createWalletIdentifier()
	identifier.Path = "t"
	err := service.RestoreWallet(identifier, "123", "", "")
	assert.Equal(t, accounts.ErrWalletPathError, err)

	identifier.WalletType = 123
	err = service.RestoreWallet(identifier, "123", "", "")
	assert.Equal(t, "wallet type error", err.Error())
}

func TestMercuryFullChainService_EstablishWallet(t *testing.T) {
	manager := &accounts.WalletManager{
		Event:        make(chan accounts.WalletEvent, 0),
		HandleResult: make(chan bool, 0),
	}

	config := &DipperinConfig{WalletManager: manager}
	service := MakeFullChainService(config)

	go func() {
		<-service.WalletManager.Event
		service.WalletManager.HandleResult <- true
	}()

	identifier := createWalletIdentifier()
	memory, err := service.EstablishWallet(*identifier, "123", "")
	defer os.RemoveAll(util.HomeDir() + testPath)
	assert.NoError(t, err)
	assert.NotNil(t, memory)

	identifier = &accounts.WalletIdentifier{WalletType: 123}
	memory, err = service.EstablishWallet(*identifier, "123", "")
	assert.Equal(t, "wallet type error", err.Error())
	assert.Equal(t, "", memory)

	identifier = &accounts.WalletIdentifier{
		WalletType: accounts.SoftWallet,
		Path:       "t",
		WalletName: "name",
	}
	memory, err = service.EstablishWallet(*identifier, "123", "")
	assert.Equal(t, accounts.ErrWalletPathError, err)
	assert.Equal(t, "", memory)
}

func TestMercuryFullChainService_Start(t *testing.T) {
	config := DipperinConfig{
		MineMaster:       fakeMaster{},
		MineMasterServer: fakeMasterServer{},
		NodeConf:         fakeNodeConfig{},
	}
	service := MakeFullChainService(&config)
	err := service.Start()
	assert.NoError(t, err)
}

func TestMercuryFullChainService_Stop(t *testing.T) {
	config := DipperinConfig{MineMaster: fakeMaster{}}
	service := MakeFullChainService(&config)
	service.Stop()
}

func TestMercuryFullChainService_Mining(t *testing.T) {
	config := DipperinConfig{}
	service := MakeFullChainService(&config)
	assert.Equal(t, false, service.Mining())

	config = DipperinConfig{MineMaster: fakeMaster{}}
	service = MakeFullChainService(&config)
	assert.Equal(t, false, service.Mining())
}

func TestMercuryFullChainService_MineTxCount(t *testing.T) {
	config := DipperinConfig{}
	service := MakeFullChainService(&config)
	assert.Equal(t, 0, service.MineTxCount())

	config = DipperinConfig{MineMaster: fakeMaster{}}
	service = MakeFullChainService(&config)
	assert.Equal(t, 1, service.MineTxCount())
}

func TestMercuryFullChainService_StartMine(t *testing.T) {
	config := DipperinConfig{}
	service := MakeFullChainService(&config)
	err := service.StartMine()
	assert.Error(t, err)

	config = DipperinConfig{MineMaster: fakeMaster{isMine: true}}
	service = MakeFullChainService(&config)
	err = service.StartMine()
	assert.Error(t, err)

	config = DipperinConfig{MineMaster: fakeMaster{isMine: false}}
	service = MakeFullChainService(&config)
	err = service.StartMine()
	assert.NoError(t, err)
}
func TestMercuryFullChainService_StopMine(t *testing.T) {
	config := DipperinConfig{}
	service := MakeFullChainService(&config)
	err := service.StopMine()
	assert.Error(t, err)

	config = DipperinConfig{MineMaster: fakeMaster{isMine: false}}
	service = MakeFullChainService(&config)
	err = service.StopMine()
	assert.Error(t, err)

	config = DipperinConfig{MineMaster: fakeMaster{isMine: true}}
	service = MakeFullChainService(&config)
	err = service.StopMine()
	assert.NoError(t, err)
}

func TestMercuryFullChainService_SetMineCoinBase(t *testing.T) {
	config := &DipperinConfig{
		NodeConf: &fakeNodeConfig{nodeType: chain_config.NodeTypeOfVerifier},
	}
	service := MakeFullChainService(config)
	err := service.SetMineCoinBase(aliceAddr)
	assert.Equal(t, "the node isn't mineMaster", err.Error())

	manager := createWalletManager(t)
	defer os.RemoveAll(util.HomeDir() + testPath)
	account, err := manager.Wallets[0].Accounts()
	assert.NoError(t, err)

	config = &DipperinConfig{
		WalletManager: manager,
		NodeConf:      fakeNodeConfig{nodeType: chain_config.NodeTypeOfMineMaster},
		MineMaster:    fakeMaster{isMine: false},
		MsgSigner:     accounts.MakeWalletSigner(account[0].Address, manager),
	}
	service = MakeFullChainService(config)
	err = service.SetMineCoinBase(aliceAddr)
	assert.Equal(t, "can not find the target wallet of this address, or the wallet is not open", err.Error())

	err = service.SetMineCoinBase(account[0].Address)
	assert.NoError(t, err)
}

func TestMercuryFullChainService_GetMineMasterReward(t *testing.T) {
	csChain := createCsChain(nil)
	insertBlockToChain(t, csChain, 1)

	config := DipperinConfig{ChainReader: csChain}
	service := MakeFullChainService(&config)

	reward, err := service.GetMineMasterEDIPReward(0, 2)
	assert.Error(t, err)
	assert.Nil(t, reward)

	reward, err = service.GetMineMasterEDIPReward(1, 2)
	assert.NoError(t, err)
	assert.NotNil(t, reward)

	reward, err = service.GetMineMasterDIPReward(1)
	assert.NoError(t, err)
	assert.NotNil(t, reward)
}

func TestMercuryFullChainService_getSendTxInfo_Error(t *testing.T) {
	manager := createWalletManager(t)
	defer os.Remove(util.HomeDir() + testPath)
	account, err := manager.Wallets[0].Accounts()
	assert.NoError(t, err)

	csChain := createCsChain(nil)
	block := csChain.CurrentBlock()
	csChain.ChainDB.DeleteBlock(block.Hash(), block.Number())
	config := &DipperinConfig{
		NodeConf:      fakeNodeConfig{},
		WalletManager: manager,
		ChainReader:   csChain,
	}
	service := MakeFullChainService(config)

	// FindWalletFromAddress error
	wallet, nonce, err := service.getSendTxInfo(common.HexToAddress("123"), nil)
	assert.Equal(t, accounts.ErrNotFindWallet, err)
	assert.Equal(t, uint64(0), nonce)
	assert.Nil(t, wallet)

	// CurrentState error
	wallet, nonce, err = service.getSendTxInfo(account[0].Address, nil)
	assert.Equal(t, "current block is nil", err.Error())
	assert.Equal(t, uint64(0), nonce)
	assert.Nil(t, wallet)

	// GetNonce error
	csChain.ChainDB.InsertBlock(block)
	config.ChainReader = csChain
	service = MakeFullChainService(config)
	wallet, nonce, err = service.getSendTxInfo(account[0].Address, nil)
	assert.Equal(t, g_error.AccountNotExist, err)
	assert.Equal(t, uint64(0), nonce)
	assert.Nil(t, wallet)
}

func TestMercuryFullChainService_Transaction(t *testing.T) {
	csChain := createCsChain(nil)
	config := DipperinConfig{ChainReader: csChain}
	service := MercuryFullChainService{
		DipperinConfig: &config,
		TxValidator:    fakeValidator{},
	}

	tx, blockHash, blockNum, txIndex, err := service.Transaction(common.Hash{})
	var expect *model.Transaction
	assert.NoError(t, err)
	assert.Equal(t, expect, tx)
	assert.Equal(t, common.Hash{}, blockHash)
	assert.Equal(t, uint64(0), blockNum)
	assert.Equal(t, uint64(0), txIndex)
}

func TestMercuryFullChainService_NewSendTransactions(t *testing.T) {
	csChain := createCsChain(nil)
	config := DipperinConfig{ChainReader: csChain, TxPool: fakeTxPool{}}
	service := MercuryFullChainService{
		DipperinConfig: &config,
		TxValidator:    fakeValidator{},
	}

	tx := createSignedTx(0, aliceAddr, big.NewInt(1000), []byte{})
	txLen, err := service.NewSendTransactions([]model.Transaction{*tx})
	assert.NoError(t, err)
	assert.Equal(t, 1, txLen)

	// tx pool AddLocals failed
	config.TxPool = fakeTxPool{err: testErr}
	service.DipperinConfig = &config
	txLen, err = service.NewSendTransactions([]model.Transaction{*tx})
	assert.Equal(t, testErr, err)
	assert.Equal(t, 0, txLen)
}

func TestMercuryFullChainService_signTxAndSend_Error(t *testing.T) {
	manager := createWalletManager(t)
	defer os.Remove(util.HomeDir() + testPath)
	account, err := manager.Wallets[0].Accounts()

	config := DipperinConfig{
		WalletManager: manager,
		TxPool:        fakeTxPool{err: testErr},
	}
	service := MercuryFullChainService{
		DipperinConfig: &config,
		TxValidator:    fakeValidator{},
	}

	sk, err := crypto.HexToECDSA(alicePriv)
	assert.NoError(t, err)
	tx := createSignedTx2(0, sk, common.HexToAddress("123"), big.NewInt(1000))

	// SignTx error
	result, err := service.signTxAndSend(manager.Wallets[0], aliceAddr, tx, 0)
	assert.Nil(t, result)
	assert.Equal(t, accounts.ErrInvalidAddress, err)

	// SignTx
	//tx2 := createSignedTx2(0, sk, aliceAddr, big.NewInt(1000))
	//result,err = service.signTxAndSend(manager.Wallets[0], aliceAddr, tx2, 0)
	//assert.NoError(t, err)

	// AddRemotes error
	result, err = service.signTxAndSend(manager.Wallets[0], account[0].Address, tx, 0)
	assert.Nil(t, result)
	assert.Equal(t, testErr, err)
}

func TestMercuryFullChainService_SendTransactions(t *testing.T) {
	manager := createWalletManager(t)
	defer os.Remove(util.HomeDir() + testPath)
	account, err := manager.Wallets[0].Accounts()
	assert.NoError(t, err)

	address := account[0].Address
	pk, err := manager.Wallets[0].GetSKFromAddress(address)
	testAccount := tests.NewAccount(pk, address)
	testAccounts := []tests.Account{*testAccount}

	csChain := createCsChain(testAccounts)
	config := &DipperinConfig{
		WalletManager: manager,
		ChainReader:   csChain,
		TxPool:        createTxPool(csChain),
		ChainConfig:   *chain_config.GetChainConfig(),
	}

	service := MercuryFullChainService{
		DipperinConfig: config,
		TxValidator:    fakeValidator{},
	}

	tx := model.RpcTransaction{
		To:       aliceAddr,
		Value:    big.NewInt(100),
		Nonce:    uint64(0),
		GasPrice: g_testData.TestGasPrice,
		GasLimit: g_testData.TestGasLimit,
	}

	// No error
	num, err := service.SendTransactions(address, []model.RpcTransaction{tx})
	assert.NoError(t, err)
	assert.Equal(t, 1, num)

	// tx pool AddLocals failed
	num, err = service.SendTransactions(address, []model.RpcTransaction{tx, tx})
	assert.Equal(t, "this transaction already in tx pool", err.Error())
	assert.Equal(t, 0, num)

	// Valid tx error
	service = MercuryFullChainService{
		DipperinConfig: config,
		TxValidator:    fakeValidator{err: testErr},
	}
	num, err = service.SendTransactions(address, []model.RpcTransaction{tx})
	assert.Equal(t, testErr, err)
	assert.Equal(t, 0, num)

	// FindWalletFromAddress error
	num, err = service.SendTransactions(common.HexToAddress("123"), []model.RpcTransaction{tx})
	assert.Equal(t, accounts.ErrNotFindWallet, err)
	assert.Equal(t, 0, num)
}

func TestMercuryFullChainService_SendTransaction(t *testing.T) {
	manager := createWalletManager(t)
	defer os.Remove(util.HomeDir() + testPath)
	account, err := manager.Wallets[0].Accounts()
	assert.NoError(t, err)

	address := account[0].Address
	pk, err := manager.Wallets[0].GetSKFromAddress(address)
	testAccount := tests.NewAccount(pk, address)
	testAccounts := []tests.Account{*testAccount}

	serviceChain := createCsChainService(testAccounts)
	txPool := createTxPool(serviceChain.ChainState)
	serviceChain.TxPool = txPool

	broadcaster := chain_communication.NewBroadcastDelegate(txPool, fakeNodeConfig{}, fakePeerManager{}, serviceChain, fakePbftNode{})
	config := &DipperinConfig{
		NodeConf:      fakeNodeConfig{nodeType: chain_config.NodeTypeOfVerifier},
		WalletManager: manager,
		ChainReader:   serviceChain,
		TxPool:        txPool,
		ChainConfig:   *chain_config.GetChainConfig(),
		Broadcaster:   broadcaster,
	}

	service := MercuryFullChainService{
		DipperinConfig: config,
		TxValidator:    fakeValidator{},
	}

	nonce := uint64(0)
	value := big.NewInt(100)
	txFee := testFee
	hash, err := service.SendRegisterTransaction(address, value, g_testData.TestGasPrice, g_testData.TestGasLimit, &nonce)
	assert.NoError(t, err)
	assert.NotNil(t, hash)

	nonce = uint64(1)
	hash, err = service.SendCancelTransaction(address, g_testData.TestGasPrice, g_testData.TestGasLimit, &nonce)
	assert.NoError(t, err)
	assert.NotNil(t, hash)

	nonce = uint64(2)
	hash, err = service.SendUnStakeTransaction(address, g_testData.TestGasPrice, g_testData.TestGasLimit, &nonce)
	assert.NoError(t, err)
	assert.NotNil(t, hash)

	nonce = uint64(3)
	vote := &model.VoteMsg{}
	hash, err = service.SendEvidenceTransaction(address, aliceAddr, g_testData.TestGasPrice, g_testData.TestGasLimit, vote, vote, &nonce)
	assert.NoError(t, err)
	assert.NotNil(t, hash)

	nonce = uint64(4)
	hash, err = service.SendTransaction(address, aliceAddr, value, g_testData.TestGasPrice, g_testData.TestGasLimit, []byte{}, &nonce)
	assert.NoError(t, err)
	assert.NotNil(t, hash)

	nonce = uint64(5)
	hash, err = service.SendTransaction(common.Address{}, aliceAddr, value, g_testData.TestGasPrice, g_testData.TestGasLimit, []byte{}, &nonce)
	assert.Equal(t, "no default account in this node", err.Error())
	assert.Equal(t, common.Hash{}, hash)

	nonce = uint64(6)
	to := common.HexToAddress(common.AddressContractCreate)
	abiPath := "../../vm/event/token/token.cpp.abi.json"
	wasmPath := "../../vm/event/token/token-wh.wasm"
	err, data := getCreateExtraData(abiPath, wasmPath, "dipp,DIPP,100000000")
	assert.NoError(t, err)
	hash, err = service.SendTransactionContract(address, to, value, new(big.Int).Mul(txFee, new(big.Int).SetInt64(int64(30))), new(big.Int).SetInt64(int64(1)), data, &nonce)
	assert.NoError(t, err)

	nonce = uint64(7)
	fs1 := model.NewMercurySigner(big.NewInt(1))
	tx := model.NewTransaction(nonce, aliceAddr, value, g_testData.TestGasPrice, g_testData.TestGasLimit, []byte{})
	signedTx, _ := tx.SignTx(pk, fs1)
	hash, err = service.NewTransaction(*signedTx)
	assert.NoError(t, err)
	assert.NotNil(t, hash)

	hash, err = service.NewTransaction(*signedTx)
	assert.Equal(t, "this transaction already in tx pool", err.Error())
	assert.Equal(t, common.Hash{}, hash)
}

func getCreateExtraData(abiPath, wasmPath string, params string) (err error, extraData []byte) {
	// GetContractExtraData
	abiBytes, err := ioutil.ReadFile(abiPath)
	if err != nil {
		return
	}
	var wasmAbi utils.WasmAbi
	err = wasmAbi.FromJson(abiBytes)
	if err != nil {
		return
	}

	//params := []string{"dipp", "DIPP", "100000000"}
	wasmBytes, err := ioutil.ReadFile(wasmPath)
	if err != nil {
		return
	}
	rlpParams := []interface{}{
		wasmBytes, abiBytes, params,
	}

	data, err := rlp.EncodeToBytes(rlpParams)
	if err != nil {
		return err, nil
	}
	return err, data
}

/*func TestGetTxData(t *testing.T)  {
	ownAddress := common.HexToAddress("0x000062be10f46b5d01Ecd9b502c4bA3d6131f6fc2e41")
	to := common.HexToAddress(common.AddressContractCreate)
	sw, err := soft_wallet.NewSoftWallet()
	assert.NoError(t, err)
	sw.Open("/Users/konggan/tmp/dipperin_apps/node/CSWallet", "CSWallet","123")

	tx := model.NewTransactionSc(0, &to, new(big.Int).SetInt64(10), new(big.Int).SetInt64(1110), 10, extraData)

	account := accounts.Account{ownAddress}
	signCallTx, err := sw.SignTx(account, callTx, nil )
}*/

func TestMercuryFullChainService_SendTransaction_Error(t *testing.T) {
	manager := createWalletManager(t)
	defer os.Remove(util.HomeDir() + testPath)
	account, err := manager.Wallets[0].Accounts()
	assert.NoError(t, err)

	address := account[0].Address
	pk, err := manager.Wallets[0].GetSKFromAddress(address)
	testAccount := tests.NewAccount(pk, address)
	testAccounts := []tests.Account{*testAccount}
	csChain := createCsChainService(testAccounts)

	config := &DipperinConfig{
		NodeConf:      fakeNodeConfig{nodeType: chain_config.NodeTypeOfVerifier},
		WalletManager: manager,
		ChainReader:   csChain,
		ChainConfig:   *chain_config.GetChainConfig(),
		TxPool:        createTxPool(csChain.ChainState),
	}

	service := MercuryFullChainService{
		DipperinConfig: config,
		TxValidator:    fakeValidator{err: testErr},
	}

	// signTxAndSend-valid error
	nonce := uint64(0)
	value := big.NewInt(100)
	hash, err := service.SendRegisterTransaction(address, value, g_testData.TestGasPrice, g_testData.TestGasLimit, &nonce)
	assert.Equal(t, testErr, err)
	assert.Equal(t, common.Hash{}, hash)

	hash, err = service.SendCancelTransaction(address, g_testData.TestGasPrice, g_testData.TestGasLimit, &nonce)
	assert.Equal(t, testErr, err)
	assert.Equal(t, common.Hash{}, hash)

	hash, err = service.SendUnStakeTransaction(address, g_testData.TestGasPrice, g_testData.TestGasLimit, &nonce)
	assert.Equal(t, testErr, err)
	assert.Equal(t, common.Hash{}, hash)

	vote := &model.VoteMsg{}
	hash, err = service.SendEvidenceTransaction(address, aliceAddr, g_testData.TestGasPrice, g_testData.TestGasLimit, vote, vote, &nonce)
	assert.Equal(t, testErr, err)
	assert.Equal(t, common.Hash{}, hash)

	hash, err = service.SendTransaction(address, aliceAddr, value, g_testData.TestGasPrice, g_testData.TestGasLimit, []byte{}, &nonce)
	assert.Equal(t, testErr, err)
	assert.Equal(t, common.Hash{}, hash)

	hash, err = service.NewTransaction(*createSignedTx(nonce, aliceAddr, value, []byte{}))
	assert.Equal(t, testErr, err)
	assert.Equal(t, common.Hash{}, hash)

	// getSendTxInfo error
	fakeAddr := common.HexToAddress("123")
	hash, err = service.SendRegisterTransaction(fakeAddr, value, g_testData.TestGasPrice, g_testData.TestGasLimit, &nonce)
	assert.Equal(t, accounts.ErrNotFindWallet, err)
	assert.Equal(t, common.Hash{}, hash)

	hash, err = service.SendCancelTransaction(fakeAddr, g_testData.TestGasPrice, g_testData.TestGasLimit, &nonce)
	assert.Equal(t, accounts.ErrNotFindWallet, err)
	assert.Equal(t, common.Hash{}, hash)

	hash, err = service.SendUnStakeTransaction(fakeAddr, g_testData.TestGasPrice, g_testData.TestGasLimit, &nonce)
	assert.Equal(t, accounts.ErrNotFindWallet, err)
	assert.Equal(t, common.Hash{}, hash)

	hash, err = service.SendEvidenceTransaction(fakeAddr, aliceAddr, g_testData.TestGasPrice, g_testData.TestGasLimit, vote, vote, &nonce)
	assert.Equal(t, accounts.ErrNotFindWallet, err)
	assert.Equal(t, common.Hash{}, hash)

	hash, err = service.SendTransaction(fakeAddr, aliceAddr, value, g_testData.TestGasPrice, g_testData.TestGasLimit, []byte{}, &nonce)
	assert.Equal(t, accounts.ErrNotFindWallet, err)
	assert.Equal(t, common.Hash{}, hash)

	// Type error
	config.NodeConf = fakeNodeConfig{nodeType: chain_config.NodeTypeOfMineMaster}
	service.DipperinConfig = config
	hash, err = service.SendRegisterTransaction(address, value, g_testData.TestGasPrice, g_testData.TestGasLimit, &nonce)
	assert.Equal(t, "the node isn't verifier", err.Error())
	assert.Equal(t, common.Hash{}, hash)

	hash, err = service.SendCancelTransaction(address, g_testData.TestGasPrice, g_testData.TestGasLimit, &nonce)
	assert.Equal(t, "the node isn't verifier", err.Error())
	assert.Equal(t, common.Hash{}, hash)

	hash, err = service.SendUnStakeTransaction(address, g_testData.TestGasPrice, g_testData.TestGasLimit, &nonce)
	assert.Equal(t, "the node isn't verifier", err.Error())
	assert.Equal(t, common.Hash{}, hash)

	hash, err = service.SendEvidenceTransaction(address, aliceAddr, g_testData.TestGasPrice, g_testData.TestGasLimit, vote, vote, &nonce)
	assert.Equal(t, "the node isn't verifier", err.Error())
	assert.Equal(t, common.Hash{}, hash)
}

func TestMakeFullChainService_EconomyModel(t *testing.T) {
	csChain := createCsChain(nil)
	config := DipperinConfig{ChainReader: csChain}
	service := MakeFullChainService(&config)

	reward, err := service.GetOneBlockTotalDIPReward(0)
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(0), reward)

	reward, err = service.GetOneBlockTotalDIPReward(1)
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(2e10), reward)

	info := service.GetInvestorInfo()
	assert.NotNil(t, info)

	info = service.GetDeveloperInfo()
	assert.NotNil(t, info)

	info = service.GetFoundationInfo(0)
	assert.NotNil(t, info)

	address := chain.VerifierAddress[0]
	DIP, err := service.GetInvestorLockDIP(address, 0)
	assert.Error(t, err)
	assert.Equal(t, big.NewInt(0), DIP)

	DIP, err = service.GetDeveloperLockDIP(address, 0)
	assert.Error(t, err)
	assert.Equal(t, big.NewInt(0), DIP)

	DIP, err = service.GetMaintenanceLockDIP(address, 0)
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(0), DIP)

	DIP, err = service.GetReMainRewardLockDIP(address, 0)
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(0), DIP)

	DIP, err = service.GetEarlyTokenLockDIP(address, 0)
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(0), DIP)

}

func TestMercuryFullChainService_StopDipperin(t *testing.T) {
	config := DipperinConfig{Node: fakeNode{}}
	service := MakeFullChainService(&config)

	service.StopDipperin()
	time.Sleep(time.Millisecond * 100)
}

func TestMercuryFullChainService_Metrics(t *testing.T) {
	config := DipperinConfig{}
	service := MakeFullChainService(&config)

	service.Metrics(false)
	service.NewBlock(context.Background())
	service.SubscribeBlock(context.Background())
}
