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
	"errors"
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/config"
	"github.com/dipperin/dipperin-core/common/g-error"
	"github.com/dipperin/dipperin-core/common/g-event"
	"github.com/dipperin/dipperin-core/common/g-metrics"
	"github.com/dipperin/dipperin-core/common/g-timer"
	"github.com/dipperin/dipperin-core/common/hexutil"
	"github.com/dipperin/dipperin-core/common/math"
	"github.com/dipperin/dipperin-core/core/accounts"
	"github.com/dipperin/dipperin-core/core/accounts/soft-wallet"
	"github.com/dipperin/dipperin-core/core/chain-communication"
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/dipperin/dipperin-core/core/chain/state-processor"
	"github.com/dipperin/dipperin-core/core/contract"
	"github.com/dipperin/dipperin-core/core/cs-chain/chain-writer/middleware"
	"github.com/dipperin/dipperin-core/core/economy-model"
	"github.com/dipperin/dipperin-core/core/mine/minemaster"
	"github.com/dipperin/dipperin-core/core/mine/mineworker"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/core/vm"
	"github.com/dipperin/dipperin-core/core/vm/common/utils"
	model2 "github.com/dipperin/dipperin-core/core/vm/model"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/dipperin/dipperin-core/third-party/log/pbft_log"
	"github.com/dipperin/dipperin-core/third-party/p2p"
	"github.com/dipperin/dipperin-core/third-party/p2p/enode"
	"github.com/dipperin/dipperin-core/third-party/rpc"
	"math/big"
	"os"
	"strings"
	"time"
)

type NodeConf interface {
	GetNodeType() int
	GetIsStartMine() bool
	SoftWalletName() string
	SoftWalletDir() string
	GetUploadURL() string
	GetNodeName() string
	GetNodeP2PPort() string
	GetNodeHTTPPort() string
}

type Chain interface {
	CurrentBlock() model.AbstractBlock
	GetBlockByHash(hash common.Hash) model.AbstractBlock
	GetBlockByNumber(number uint64) model.AbstractBlock
	Genesis() model.AbstractBlock
	GetBody(hash common.Hash) model.AbstractBody
	GetBlockNumber(hash common.Hash) *uint64
	CurrentState() (*state_processor.AccountStateDB, error)

	CurrentSeed() (common.Hash, uint64)
	NumBeforeLastBySlot(slot uint64) *uint64
	StateAtByBlockNumber(num uint64) (*state_processor.AccountStateDB, error)
	GetTransaction(txHash common.Hash) (model.AbstractTransaction, common.Hash, uint64, uint64)
	GetVerifiers(round uint64) []common.Address
	GetSlot(block model.AbstractBlock) *uint64
	GetCurrVerifiers() []common.Address
	GetNextVerifiers() []common.Address
	CurrentHeader() model.AbstractHeader

	GetEconomyModel() economy_model.EconomyModel

	GetReceipts(hash common.Hash, number uint64) model2.Receipts
}

type TxPool interface {
	AddRemotes(txs []model.AbstractTransaction) []error
	AddLocals(txs []model.AbstractTransaction) []error
	AddRemote(tx model.AbstractTransaction) error
	Stats() (int, int)
}

type Node interface {
	Start() error
	Stop()
}

type TxValidator interface {
	Valid(tx model.AbstractTransaction) error
}

type Broadcaster interface {
	BroadcastTx(txs []model.AbstractTransaction)
}

type MsgSigner interface {
	SetBaseAddress(address common.Address)
	GetAddress() common.Address
}

func MakeFullChainService(config *DipperinConfig) *MercuryFullChainService {
	return &MercuryFullChainService{
		DipperinConfig: config,
		TxValidator:    middleware.NewTxValidatorForRpcService(config.ChainReader),
	}
}

type DipperinConfig struct {
	PbftPm chain_communication.AbstractPbftProtocolManager

	Broadcaster    Broadcaster
	ChainReader    middleware.ChainInterface
	TxPool         TxPool
	MineMaster     minemaster.Master
	WalletManager  *accounts.WalletManager
	DefaultAccount common.Address

	NodeConf           NodeConf
	GetMineCoinBase    common.Address
	MsgSigner          MsgSigner
	ChainConfig        chain_config.ChainConfig
	PriorityCalculator model.PriofityCalculator
	MineMasterServer   minemaster.MasterServer
	P2PServer          *p2p.Server
	NormalPm           chain_communication.PeerManager

	Node Node
}

// deal mercury api things
type MercuryFullChainService struct {
	*DipperinConfig
	localWorker mineworker.Worker
	TxValidator TxValidator
}

func (service *MercuryFullChainService) RemoteHeight() uint64 {
	_, h := service.NormalPm.BestPeer().GetHead()
	return h
}

func (service *MercuryFullChainService) GetSyncStatus() bool {
	return service.NormalPm.IsSync()
}

func (service *MercuryFullChainService) CurrentBlock() model.AbstractBlock {
	return service.ChainReader.CurrentBlock()
}

func (service *MercuryFullChainService) GetBlockByNumber(number uint64) (model.AbstractBlock, error) {
	return service.ChainReader.GetBlockByNumber(number), nil
}

func (service *MercuryFullChainService) GetBlockByHash(hash common.Hash) (model.AbstractBlock, error) {
	return service.ChainReader.GetBlockByHash(hash), nil
}

func (service *MercuryFullChainService) GetBlockNumber(hash common.Hash) *uint64 {
	return service.ChainReader.GetBlockNumber(hash)
}

func (service *MercuryFullChainService) GetGenesis() (model.AbstractBlock, error) {
	return service.ChainReader.Genesis(), nil
}

func (service *MercuryFullChainService) GetBlockBody(hash common.Hash) model.AbstractBody {
	return service.ChainReader.GetBody(hash)
}

func (service *MercuryFullChainService) CurrentBalance(address common.Address) *big.Int {
	curState, err := service.ChainReader.CurrentState()
	if err != nil {
		log.Warn("get current state failed", "err", err)
		return nil
	}
	balance, err := curState.GetBalance(address)
	if err != nil {
		log.Info("get current balance failed", "err", err)
		return nil
	}
	log.Info("call current balance", "address", address.Hex(), "balance", balance)
	return balance
}
func (service *MercuryFullChainService) CurrentStake(address common.Address) *big.Int {
	log.Debug("call current balance", "address", address.Hex())
	curState, err := service.ChainReader.CurrentState()
	if err != nil {
		log.Warn("get current state failed", "err", err)
		return nil
	}
	stake, err := curState.GetStake(address)
	if err != nil {
		log.Info("get current balance failed", "err", err)
		return nil
	}
	log.Info("CurrentStake the stake is:", "stake", stake)
	return stake
}

func (service *MercuryFullChainService) Start() error {
	if service.MineMaster != nil && !service.MineMaster.CurrentCoinbaseAddress().IsEmpty() {
		if service.localWorker == nil {
			time.Sleep(500 * time.Millisecond)
			service.localWorker = mineworker.MakeLocalWorker(service.MineMaster.CurrentCoinbaseAddress(), 1, service.MineMasterServer)
			log.Info("start local worker")
			service.localWorker.Start()
		}

		log.Info("start mine master")
		log.Info("the service.nodeContext.nodeConf.IsStartMine is:", "isStartMine", service.NodeConf.GetIsStartMine())
		if service.NodeConf.GetIsStartMine() {
			service.MineMaster.Start()
		}
	}

	service.startTxsMetrics()

	log.Info("full chain service start success")
	return nil
}

func (service *MercuryFullChainService) startTxsMetrics() {
	g_timer.SetPeriodAndRun(func() {
		pending, queued := service.TxPool.Stats()
		g_metrics.Set(g_metrics.PendingTxCountInPool, "", float64(pending))
		g_metrics.Set(g_metrics.QueuedTxCountInPool, "", float64(queued))
	}, 5*time.Second)
}

func (service *MercuryFullChainService) Stop() {
	if service.MineMaster != nil {
		service.MineMaster.Stop()
	}
}

func (service *MercuryFullChainService) checkWalletIdentifier(walletIdentifier *accounts.WalletIdentifier) error {
	if walletIdentifier.WalletType != accounts.SoftWallet {
		return errors.New("wallet type error")
	}

	if walletIdentifier.WalletName == "" {
		walletIdentifier.WalletName = service.NodeConf.SoftWalletName()
	}

	if walletIdentifier.Path == "" {
		walletIdentifier.Path = service.NodeConf.SoftWalletDir()
	}

	return nil
}

//set CoinBase Address
func (service *MercuryFullChainService) SetMineCoinBase(addr common.Address) error {
	if service.NodeConf.GetNodeType() != chain_config.NodeTypeOfMineMaster {
		return errors.New("the node isn't mineMaster")
	}
	tmpWallet, err := service.WalletManager.FindWalletFromAddress(addr)
	if err != nil {
		return errors.New("can not find the target wallet of this address, or the wallet is not open")
	}
	state, _ := tmpWallet.Status()
	if state == "close" {
		return errors.New("target wallet is closed")
	}

	service.MsgSigner.SetBaseAddress(addr)
	service.MineMaster.SetCoinbaseAddress(addr)
	return nil
}

func (service *MercuryFullChainService) SetMineGasConfig(gasFloor, gasCeil uint64) error {
	if gasFloor < gasCeil {
		return errors.New("gasFloor should greater than gasCeil")
	}
	service.MineMaster.SetMineGasConfig(gasFloor, gasCeil)
	return nil
}

func (service *MercuryFullChainService) EstablishWallet(walletIdentifier accounts.WalletIdentifier, password, passPhrase string) (string, error) {
	err := service.checkWalletIdentifier(&walletIdentifier)
	if err != nil {
		log.Info("the err1 is :", "err", err)
		return "", err
	}

	//establish softWallet
	wallet, _ := soft_wallet.NewSoftWallet()
	mnemonic, err := wallet.Establish(walletIdentifier.Path, walletIdentifier.WalletName, password, passPhrase)
	if err != nil {
		log.Info("the err3 is :", "err", err)
		return "", err
	}

	//add softWallet to wallet manager
	testEvent := accounts.WalletEvent{
		Wallet: wallet,
		Type:   accounts.WalletArrived,
	}

	log.Info("send wallet manager event")

	service.WalletManager.Event <- testEvent

	log.Info("wait for the wallet manager handle result")
	select {
	case <-service.WalletManager.HandleResult:
	}
	return mnemonic, nil
}

func (service *MercuryFullChainService) OpenWallet(walletIdentifier accounts.WalletIdentifier, password string) error {
	err := service.checkWalletIdentifier(&walletIdentifier)
	if err != nil {
		return err
	}

	//Open according to the path
	//establish softWallet
	wallet, _ := soft_wallet.NewSoftWallet()
	err = wallet.Open(walletIdentifier.Path, walletIdentifier.WalletName, password)
	if err != nil {
		return err
	}

	//add wallet to the manager
	WalletEvent := accounts.WalletEvent{
		Wallet: wallet,
		Type:   accounts.WalletArrived,
	}

	service.WalletManager.Event <- WalletEvent

	select {
	case <-service.WalletManager.HandleResult:
	}
	return nil
}

func (service *MercuryFullChainService) CloseWallet(walletIdentifier accounts.WalletIdentifier) error {
	err := service.checkWalletIdentifier(&walletIdentifier)
	if err != nil {
		return err
	}

	//find wallet according to walletIdentifier
	tmpWallet, err := service.WalletManager.FindWalletFromIdentifier(walletIdentifier)
	if err != nil {
		return err
	}

	if service.NodeConf.GetNodeType() == chain_config.NodeTypeOfMineMaster {
		addr := service.MsgSigner.GetAddress()
		isInclude, err := tmpWallet.Contains(accounts.Account{Address: addr})
		if err != nil {
			return err
		}
		if isInclude == true {
			return errors.New("this wallet contains coinbase, can not close")
		}

	}
	WalletEvent := accounts.WalletEvent{
		Wallet: tmpWallet,
		Type:   accounts.WalletDropped,
	}

	err = tmpWallet.Close()
	if err != nil {
		return err
	}

	service.WalletManager.Event <- WalletEvent
	select {
	case <-service.WalletManager.HandleResult:
	}
	return nil
}

func (service *MercuryFullChainService) RestoreWallet(walletIdentifier accounts.WalletIdentifier, password, passPhrase, mnemonic string) error {
	err := service.checkWalletIdentifier(&walletIdentifier)
	if err != nil {
		return err
	}

	//Check if the wallet to be restored is in the walletManager, remove if it is
	findWallet, _ := service.WalletManager.FindWalletFromIdentifier(walletIdentifier)
	if findWallet != nil {
		removeEvent := accounts.WalletEvent{
			Wallet: findWallet,
			Type:   accounts.WalletDropped,
		}
		service.WalletManager.Event <- removeEvent
		select {
		case <-service.WalletManager.HandleResult:
		}
	}

	//establish softWallet
	wallet, _ := soft_wallet.NewSoftWallet()
	err = wallet.RestoreWallet(walletIdentifier.Path, walletIdentifier.WalletName, password, passPhrase, mnemonic, service)
	if err != nil {
		return err
	}

	//add the restored wallet to manager
	testEvent := accounts.WalletEvent{
		Wallet: wallet,
		Type:   accounts.WalletArrived,
	}

	service.WalletManager.Event <- testEvent

	select {
	case <-service.WalletManager.HandleResult:
	}
	return nil
}

func (service *MercuryFullChainService) ListWallet() ([]accounts.WalletIdentifier, error) {
	walletIdentifiers, err := service.WalletManager.ListWalletIdentifier()
	if err != nil {
		log.Info("the listWallet err is:", "err", err)
		return []accounts.WalletIdentifier{}, err
	}
	return walletIdentifiers, nil
}

func (service *MercuryFullChainService) ListWalletAccount(walletIdentifier accounts.WalletIdentifier) ([]accounts.Account, error) {
	err := service.checkWalletIdentifier(&walletIdentifier)
	if err != nil {
		return []accounts.Account{}, err
	}

	//find wallet according to walletIdentifier
	tmpWallet, err := service.WalletManager.FindWalletFromIdentifier(walletIdentifier)
	if err != nil {
		return []accounts.Account{}, err
	}
	return tmpWallet.Accounts()
}

func (service *MercuryFullChainService) SetBftSigner(address common.Address) error {
	log.Info("MercuryFullChainService SetWalletAccountAddress run")
	service.MsgSigner.SetBaseAddress(address)
	/*if service.nodeContext.NodeConf().NodeType == chain_config.NodeTypeOfMineMaster {
		service.nodeContext.SetMineCoinBase(address)
	}*/
	return nil
}

func (service *MercuryFullChainService) AddAccount(walletIdentifier accounts.WalletIdentifier, derivationPath string) (accounts.Account, error) {
	err := service.checkWalletIdentifier(&walletIdentifier)
	if err != nil {
		return accounts.Account{}, err
	}
	//find wallet according to walletIdentifier
	tmpWallet, err := service.WalletManager.FindWalletFromIdentifier(walletIdentifier)
	if err != nil {
		return accounts.Account{}, err
	}

	log.Info("AddAccount the path is:", "derivationPath", derivationPath)

	var path accounts.DerivationPath
	if derivationPath == "" {
		path = nil
	} else {
		path, err = accounts.ParseDerivationPath(derivationPath)
		if err != nil {
			return accounts.Account{}, err
		}
	}
	//derive new account and save
	account, err := tmpWallet.Derive(path, true)
	if err != nil {
		return accounts.Account{}, err
	}
	return account, nil
}

/*func (service *MercuryFullChainService) SyncUsedAccounts(walletIdentifier accounts.WalletIdentifier, MaxChangeValue, MaxIndex uint32) error {
	err := service.checkWalletIdentifier(&walletIdentifier)
	if err != nil {
		return err
	}
	return nil
}*/

func (service *MercuryFullChainService) getSendTxInfo(from common.Address, nonce *uint64) (accounts.Wallet, uint64, error) {
	//find wallet according to address
	tmpWallet, err := service.WalletManager.FindWalletFromAddress(from)
	if err != nil {
		log.Error("MercuryFullChainService#getSendTxInfo FindWalletFromAddress", "err", err)
		return nil, 0, err
	}
	//generate transaction
	state, err := service.ChainReader.CurrentState()
	if err != nil {
		log.Error("MercuryFullChainService#getSendTxInfo  CurrentState", "err", err)
		return nil, 0, err
	}

	//get nonce from blockChain
	chainNonce, err := state.GetNonce(from)
	if err != nil {
		log.Info("the address is:", "address", from.Hex())
		log.Info("~~~~~~~~~~~~~~~~~~~~~~~~get nonce fail", "err", err)
		return nil, 0, err
	}

	//get nonce from wallet
	walletNonce, _ := tmpWallet.GetAddressNonce(from)
	var spendableNonce uint64
	if walletNonce < chainNonce {
		spendableNonce = chainNonce
	} else {
		spendableNonce = walletNonce
	}

	//log.Info("the nonce is:", "nonce", nonce)
	var txNonce uint64
	if nonce == nil {
		txNonce = spendableNonce
	} else {
		txNonce = *nonce
	}
	return tmpWallet, txNonce, nil
}

//send single tx

func (service *MercuryFullChainService) signTxAndSend(tmpWallet accounts.Wallet, from common.Address, tx *model.Transaction, usedNonce uint64) (*model.Transaction, error) {
	fromAccount := accounts.Account{Address: from}
	//get chainId
	signedTx, err := tmpWallet.SignTx(fromAccount, tx, service.ChainConfig.ChainId)
	if err != nil {
		return nil, err
	}
	pbft_log.Debug("Sign and send transaction", "txid", signedTx.CalTxId().Hex())
	if err := service.TxValidator.Valid(signedTx); err != nil {
		log.Error("Transaction not valid", "error", err)
		return nil, err
	}

	tsx := []model.AbstractTransaction{signedTx}

	errs := service.TxPool.AddRemotes(tsx)

	for i := range errs {
		if errs[i] != nil {
			return nil, errs[i]
		}
	}

	service.Broadcaster.BroadcastTx(tsx)

	err = tmpWallet.SetAddressNonce(from, usedNonce+1)
	if err != nil {
		return nil, err
	}
	return signedTx, nil
}

//send multiple-txs
func (service *MercuryFullChainService) SendTransactions(from common.Address, rpcTxs []model.RpcTransaction) (int, error) {
	//start := time.Now()
	tmpWallet, err := service.WalletManager.FindWalletFromAddress(from)
	if err != nil {
		return 0, err
	}
	fromAccount := accounts.Account{Address: from}

	txs := make([]model.AbstractTransaction, 0)
	for _, item := range rpcTxs {
		tx := model.NewTransaction(item.Nonce, item.To, item.Value, item.GasPrice, item.GasLimit, item.Data)
		signedTx, err := tmpWallet.SignTx(fromAccount, tx, service.ChainConfig.ChainId)
		if err != nil {
			log.Info("send Transactions SignTx:", "err", err)
			return 0, err
		}

		if err := service.TxValidator.Valid(signedTx); err != nil {
			log.Info("send Transactions ValidTx:", "err", err)
			return 0, err
		}
		log.Info("the SendTransaction txId is: ", "txId", tx.CalTxId().Hex(), "txSize", tx.Size())
		txs = append(txs, tx)
	}
	errs := service.TxPool.AddLocals(txs)

	for i := range errs {
		if errs[i] != nil {
			return 0, errs[i]
		}
	}
	return len(txs), nil
}

//
func (service *MercuryFullChainService) NewSendTransactions(txs []model.Transaction) (int, error) {
	temtxs := make([]model.AbstractTransaction, 0)
	for _, item := range txs {
		temtx := item
		temtxs = append(temtxs, &temtx)
	}
	errs := service.TxPool.AddLocals(temtxs)

	for i := range errs {
		if errs[i] != nil {
			return 0, errs[i]
		}
	}
	return len(txs), nil
}

//send a normal transaction
func (service *MercuryFullChainService) SendTransaction(from, to common.Address, value, gasPrice *big.Int, gasLimit uint64, data []byte, nonce *uint64) (common.Hash, error) {
	//start:=time.Now()
	// automatic transfer need this
	if from.IsEqual(common.Address{}) {
		from = service.DefaultAccount
		if from.IsEqual(common.Address{}) {
			return common.Hash{}, errors.New("no default account in this node")
		}
	}

	//log.Info("send Transaction the nonce is:", "nonce", nonce)

	tmpWallet, usedNonce, err := service.getSendTxInfo(from, nonce)
	if err != nil {
		return common.Hash{}, err
	}

	tx := model.NewTransaction(usedNonce, to, value, gasPrice, gasLimit, data)
	signTx, err := service.signTxAndSend(tmpWallet, from, tx, usedNonce)
	if err != nil {
		pbft_log.Error("send tx error", "txid", tx.CalTxId().Hex(), "err", err)
		return common.Hash{}, err
	}

	pbft_log.Info("send transaction", "txId", signTx.CalTxId().Hex())
	txHash := signTx.CalTxId()
	log.Info("the Sendnot enough balance errorTransaction txId is: ", "txId", txHash.Hex(), "txSize", signTx.Size())
	return txHash, nil
}

//send a register transaction
func (service *MercuryFullChainService) SendRegisterTransaction(from common.Address, stake, gasPrice *big.Int, gasLimit uint64, nonce *uint64) (common.Hash, error) {
	if service.NodeConf.GetNodeType() != chain_config.NodeTypeOfVerifier {
		return common.Hash{}, errors.New("the node isn't verifier")
	}

	tmpWallet, usedNonce, err := service.getSendTxInfo(from, nonce)
	if err != nil {
		return common.Hash{}, err
	}

	tx := model.NewRegisterTransaction(usedNonce, stake, gasPrice, gasLimit)
	signTx, err := service.signTxAndSend(tmpWallet, from, tx, usedNonce)
	if err != nil {
		return common.Hash{}, err
	}

	txHash := signTx.CalTxId()
	log.Info("the SendRegisterTransaction txId is: ", "txId", txHash.Hex())
	return txHash, nil
}

func (service *MercuryFullChainService) getLuckProof(addr common.Address) (common.Hash, []byte, uint64, error) {
	tmpWallet, err := service.WalletManager.FindWalletFromAddress(addr)
	if err != nil {
		return common.Hash{}, []byte{}, 0, err
	}

	//current seed is last block num by slot's seed
	seed, blockNumber := service.ChainReader.CurrentSeed()
	fromAccount := accounts.Account{Address: addr}

	log.Info("the seed is:", "seed", seed.String())
	luck, proof, err := tmpWallet.Evaluate(fromAccount, seed.Bytes())
	if err != nil {
		return common.Hash{}, []byte{}, 0, err
	}
	return luck, proof, blockNumber, nil
}

func (service *MercuryFullChainService) CurrentElectPriority(addr common.Address) (uint64, error) {
	luck, _, _, err := service.getLuckProof(addr)
	if err != nil {
		return 0, err
	}

	slot := service.GetSlot(service.CurrentBlock())
	num := service.ChainReader.NumBeforeLastBySlot(*slot)
	if num == nil {
		log.Debug("CurrentElectPriority error", "slot", slot, "num", num)
		return 0, errors.New("number before last is nil")
	}
	log.Info("LastNumberBySlot:", "num", num)
	state, err := service.ChainReader.StateAtByBlockNumber(*num)
	if err != nil {
		return 0, err
	}

	accountNonce, err := state.GetNonce(addr)
	if err != nil {
		return 0, err
	}

	stake, err := state.GetStake(addr)
	if err != nil {
		return 0, err
	}

	performance, err := state.GetPerformance(addr)
	if err != nil {
		return 0, err
	}

	priority, err := service.PriorityCalculator.GetElectPriority(luck, accountNonce, stake, performance)
	if err != nil {
		return 0, err
	}
	return priority, nil
}

func (service *MercuryFullChainService) CurrentReputation(addr common.Address) (uint64, error) {
	state, err := service.ChainReader.CurrentState()
	if err != nil {
		return 0, err
	}
	stake, err := state.GetStake(addr)
	performance, err := state.GetPerformance(addr)

	reputation, err := service.PriorityCalculator.GetReputation(0, stake, performance)
	if err != nil {
		return 0, err
	}
	return reputation, nil
}

func (service *MercuryFullChainService) MineTxCount() int {
	if service.MineMaster != nil {
		return service.MineMaster.MineTxCount()
	}
	return 0
}

//send a evidence transaction
func (service *MercuryFullChainService) SendEvidenceTransaction(from, target common.Address, gasPrice *big.Int, gasLimit uint64, voteA *model.VoteMsg, voteB *model.VoteMsg, nonce *uint64) (common.Hash, error) {
	if service.NodeConf.GetNodeType() != chain_config.NodeTypeOfVerifier {
		return common.Hash{}, errors.New("the node isn't verifier")
	}

	tmpWallet, usedNonce, err := service.getSendTxInfo(from, nonce)
	if err != nil {
		return common.Hash{}, err
	}

	tx := model.NewEvidenceTransaction(usedNonce, gasPrice, gasLimit, &target, voteA, voteB)
	//log.Debug("SendEvidenceTransaction size", "tx size", tx.Size().String())
	signTx, err := service.signTxAndSend(tmpWallet, from, tx, usedNonce)
	if err != nil {
		return common.Hash{}, err
	}

	txHash := signTx.CalTxId()
	log.Info("the SendEvidenceTransaction txId is: ", "txId", txHash.Hex())
	return txHash, nil
}

//Send redemption transaction
func (service *MercuryFullChainService) SendUnStakeTransaction(from common.Address, gasPrice *big.Int, gasLimit uint64, nonce *uint64) (common.Hash, error) {
	if service.NodeConf.GetNodeType() != chain_config.NodeTypeOfVerifier {
		return common.Hash{}, errors.New("the node isn't verifier")
	}

	tmpWallet, usedNonce, err := service.getSendTxInfo(from, nonce)
	if err != nil {
		return common.Hash{}, err
	}

	tx := model.NewUnStakeTransaction(usedNonce, gasPrice, gasLimit)
	signTx, err := service.signTxAndSend(tmpWallet, from, tx, usedNonce)
	if err != nil {
		return common.Hash{}, err
	}

	txHash := signTx.CalTxId()
	log.Info("the SendCancelTransaction txId is: ", "txId", txHash.Hex())
	return txHash, nil
}

//send a cancellation transaction
func (service *MercuryFullChainService) SendCancelTransaction(from common.Address, gasPrice *big.Int, gasLimit uint64, nonce *uint64) (common.Hash, error) {
	if service.NodeConf.GetNodeType() != chain_config.NodeTypeOfVerifier {
		return common.Hash{}, errors.New("the node isn't verifier")
	}

	tmpWallet, usedNonce, err := service.getSendTxInfo(from, nonce)
	if err != nil {
		return common.Hash{}, err
	}

	tx := model.NewCancelTransaction(usedNonce, gasPrice, gasLimit)
	signTx, err := service.signTxAndSend(tmpWallet, from, tx, usedNonce)
	if err != nil {
		return common.Hash{}, err
	}

	txHash := signTx.CalTxId()
	log.Info("the SendCancelTransaction txId is: ", "txId", txHash.Hex())
	return txHash, nil
}

//get address nonce from chain
func (service *MercuryFullChainService) GetTransactionNonce(addr common.Address) (nonce uint64, err error) {
	state, err := service.ChainReader.CurrentState()
	if err != nil {
		return 0, err
	}
	nonce, err = state.GetNonce(addr)
	if err != nil {
		return 0, err
	}
	return nonce, nil
}

//get address nonce from wallet
func (service *MercuryFullChainService) GetAddressNonceFromWallet(address common.Address) (nonce uint64, err error) {
	//find wallet according to address
	tmpWallet, err := service.WalletManager.FindWalletFromAddress(address)
	if err != nil {
		return 0, err
	}
	return tmpWallet.GetAddressNonce(address)
}

// wallet initiates a transaction
func (service *MercuryFullChainService) NewTransaction(transaction model.Transaction) (txHash common.Hash, err error) {

	//todo delete after test
	log.Info("NewTransaction ~~~~~~~~~~~~~~~~~~~~~~~~~~~~", "txId", transaction.CalTxId().Hex())
	if err = service.TxValidator.Valid(&transaction); err != nil {
		log.Info("NewTransaction validTx result is:", "err", err)
		return
	}

	err = service.TxPool.AddRemote(&transaction)
	if err != nil {
		return common.Hash{}, err
	}

	//todo: Here the local wallet Nonce maintains the nonce value used by the wallet. Therefore, when the wallet and the command line are used to send the transaction at the same time, the nonce may be invalid and the transaction may not be packaged in the transaction pool.
	//broadcast  transaction
	log.Info("[NewTransaction] broadcast transaction~~~~~~~~~~~~~~~~")
	service.Broadcaster.BroadcastTx([]model.AbstractTransaction{&transaction})

	txHash = transaction.CalTxId()
	return txHash, nil
}

// consult a transaction
func (service *MercuryFullChainService) Transaction(hash common.Hash) (transaction *model.Transaction, blockHash common.Hash, blockNumber uint64, txIndex uint64, err error) {

	tx, blockHash, blockNum, txIndex := service.ChainReader.GetTransaction(hash)

	if tx != nil {
		transaction = tx.(*model.Transaction)
	}
	return transaction, blockHash, blockNum, txIndex, nil
}

//Test get verifiers of this round
func (service *MercuryFullChainService) GetVerifiers(slotNum uint64) (addresses []common.Address) {
	addresses = service.ChainReader.GetVerifiers(slotNum)
	log.Debug("Get verifiers addresses", "slot", slotNum, "Length", len(addresses), "addresses", addresses)
	return addresses
}

func (service *MercuryFullChainService) GetSlot(block model.AbstractBlock) *uint64 {
	return service.ChainReader.GetSlot(block)
}

func (service *MercuryFullChainService) GetCurVerifiers() []common.Address {
	return service.ChainReader.GetCurrVerifiers()
}

func (service *MercuryFullChainService) GetNextVerifiers() []common.Address {
	return service.ChainReader.GetNextVerifiers()
}

func (service *MercuryFullChainService) VerifierStatus(addr common.Address) (verifierState string, stake *big.Int, balance *big.Int, reputation uint64, isCurrentVerifier bool, err error) {
	status := []string{"Not Registered", "Registered", "Canceled", "Unstaked"}
	verifierState = status[0]
	state, err := service.ChainReader.CurrentState()
	if err != nil {
		return
	}
	stake, err = state.GetStake(addr)
	if err != nil {
		if err.Error() != "account does not exist" && err.Error() != "stake not sufficient" {
			return
		}
	}

	balance, err = state.GetBalance(addr)
	if err != nil {
		if err.Error() != "account does not exist" {
			return
		}
	}

	lastElect, err := state.GetLastElect(addr)
	if err != nil {
		if err.Error() != "account does not exist" {
			return
		}
	}

	//Not Registered
	if lastElect == 0 && stake.Cmp(big.NewInt(0)) == 0 {
		verifierState = status[0]
	}

	//Registered
	if lastElect == 0 && stake.Cmp(big.NewInt(0)) != 0 {
		verifierState = status[1]
	}

	//Canceled
	if lastElect != 0 && stake.Cmp(big.NewInt(0)) != 0 {
		verifierState = status[2]
	}

	//Unstaked
	if lastElect != 0 && stake.Cmp(big.NewInt(0)) == 0 {
		verifierState = status[3]
	}

	isCurrentVerifier = service.isCurrentVerifier(addr)

	reputation, err = service.CurrentReputation(addr)
	if err != nil {
		if err.Error() == "account does not exist" || err.Error() == "stake not sufficient" {
			err = nil
		}
	}
	return
}

func (service *MercuryFullChainService) isCurrentVerifier(address common.Address) bool {
	vers := service.ChainReader.GetCurrVerifiers()
	for v := range vers {
		if vers[v].IsEqual(address) {
			return true
		}
	}
	return false
}

func (service *MercuryFullChainService) GetCurrentConnectPeers() map[string]common.Address {
	if service.PbftPm != nil {
		return service.PbftPm.GetCurrentConnectPeers()
	} else {
		return make(map[string]common.Address, 0)
	}
}

// start mine
func (service *MercuryFullChainService) StartMine() error {
	if service.MineMaster == nil {
		return errors.New("current node is not mine master")
	}

	if service.Mining() {
		return errors.New("miner is mining")
	}

	service.MineMaster.Start()
	return nil
}

// stop mine
func (service *MercuryFullChainService) StopMine() error {
	if service.MineMaster == nil {
		return errors.New("current node is not mine master")
	}

	if !service.Mining() {
		return errors.New("mining had been stopped")
	}

	service.MineMaster.Stop()
	return nil
}

// check if is mining
func (service *MercuryFullChainService) Mining() bool {
	if service.MineMaster != nil {
		return service.MineMaster.Mining()
	}
	return false
}

// debug
func (service *MercuryFullChainService) Metrics(raw bool) (map[string]interface{}, error) {
	/*// Create a rate formatter
	units := []string{"", "K", "M", "G", "T", "E", "P"}
	round := func(value float64, prec int) string {
		unit := 0
		for value >= 1000 {
			unit, value, prec = unit+1, value/1000, 2
		}
		return fmt.Sprintf(fmt.Sprintf("%%.%df%s", prec, units[unit]), value)
	}
	format := func(total float64, rate float64) string {
		return fmt.Sprintf("%s (%s/s)", round(total, 0), round(rate, 2))
	}
	// Iterate over all the metrics, and just dump for now
	counters := make(map[string]interface{})
	metrics.DefaultRegistry.Each(func(name string, metric interface{}) {
		// Create or retrieve the counter hierarchy for this metric
		root, parts := counters, strings.Split(name, "/")
		for _, part := range parts[:len(parts)-1] {
			if _, ok := root[part]; !ok {
				root[part] = make(map[string]interface{})
			}
			root = root[part].(map[string]interface{})
		}
		name = parts[len(parts)-1]

		// Fill the counter with the metric details, formatting if requested
		if raw {
			switch metric := metric.(type) {
			case metrics.Counter:
				root[name] = map[string]interface{}{
					"Overall": float64(metric.Count()),
				}

			case metrics.Meter:
				root[name] = map[string]interface{}{
					"AvgRate01Min": metric.Rate1(),
					"AvgRate05Min": metric.Rate5(),
					"AvgRate15Min": metric.Rate15(),
					"MeanRate":     metric.RateMean(),
					"Overall":      float64(metric.Count()),
				}

			case metrics.Timer:
				root[name] = map[string]interface{}{
					"AvgRate01Min": metric.Rate1(),
					"AvgRate05Min": metric.Rate5(),
					"AvgRate15Min": metric.Rate15(),
					"MeanRate":     metric.RateMean(),
					"Overall":      float64(metric.Count()),
					"Percentiles": map[string]interface{}{
						"5":  metric.Percentile(0.05),
						"20": metric.Percentile(0.2),
						"50": metric.Percentile(0.5),
						"80": metric.Percentile(0.8),
						"95": metric.Percentile(0.95),
					},
				}

			case metrics.ResettingTimer:
				t := metric.Snapshot()
				ps := t.Percentiles([]float64{5, 20, 50, 80, 95})
				root[name] = map[string]interface{}{
					"Measurements": len(t.Values()),
					"Mean":         t.Mean(),
					"Percentiles": map[string]interface{}{
						"5":  ps[0],
						"20": ps[1],
						"50": ps[2],
						"80": ps[3],
						"95": ps[4],
					},
				}

			default:
				root[name] = "Unknown metric type"
			}
		} else {
			switch metric := metric.(type) {
			case metrics.Counter:
				root[name] = map[string]interface{}{
					"Overall": float64(metric.Count()),
				}

			case metrics.Meter:
				root[name] = map[string]interface{}{
					"Avg01Min": format(metric.Rate1()*60, metric.Rate1()),
					"Avg05Min": format(metric.Rate5()*300, metric.Rate5()),
					"Avg15Min": format(metric.Rate15()*900, metric.Rate15()),
					"Overall":  format(float64(metric.Count()), metric.RateMean()),
				}

			case metrics.Timer:
				root[name] = map[string]interface{}{
					"Avg01Min": format(metric.Rate1()*60, metric.Rate1()),
					"Avg05Min": format(metric.Rate5()*300, metric.Rate5()),
					"Avg15Min": format(metric.Rate15()*900, metric.Rate15()),
					"Overall":  format(float64(metric.Count()), metric.RateMean()),
					"Maximum":  time.Duration(metric.Max()).String(),
					"Minimum":  time.Duration(metric.Min()).String(),
					"Percentiles": map[string]interface{}{
						"5":  time.Duration(metric.Percentile(0.05)).String(),
						"20": time.Duration(metric.Percentile(0.2)).String(),
						"50": time.Duration(metric.Percentile(0.5)).String(),
						"80": time.Duration(metric.Percentile(0.8)).String(),
						"95": time.Duration(metric.Percentile(0.95)).String(),
					},
				}

			case metrics.ResettingTimer:
				t := metric.Snapshot()
				ps := t.Percentiles([]float64{5, 20, 50, 80, 95})
				root[name] = map[string]interface{}{
					"Measurements": len(t.Values()),
					"Mean":         time.Duration(t.Mean()).String(),
					"Percentiles": map[string]interface{}{
						"5":  time.Duration(ps[0]).String(),
						"20": time.Duration(ps[1]).String(),
						"50": time.Duration(ps[2]).String(),
						"80": time.Duration(ps[3]).String(),
						"95": time.Duration(ps[4]).String(),
					},
				}

			default:
				root[name] = "Unknown metric type"
			}
		}
	})
	return counters, nil*/
	return nil, nil
}

// add peer
func (service *MercuryFullChainService) AddPeer(url string) error {
	server := service.P2PServer
	if server == nil {
		return errors.New("no p2p server running")
	}

	node, err := enode.ParseV4(url)

	if err != nil {
		return fmt.Errorf("invalid url: %v", err)
	}
	server.AddPeer(node)
	return nil
}

// remove peer
func (service *MercuryFullChainService) RemovePeer(url string) error {
	server := service.P2PServer
	if server == nil {
		return errors.New("no p2p server running")
	}

	node, err := enode.ParseV4(url)

	if err != nil {
		return fmt.Errorf("invalid url: %v", err)
	}
	server.RemovePeer(node)
	return nil
}

func (service *MercuryFullChainService) CsPmInfo() (*p2p.CsPmPeerInfo, error) {
	pm := service.NormalPm.(*chain_communication.CsProtocolManager)
	return pm.ShowPmInfo(), nil
}

// AddTrustedPeer allows a remote node to always connect, even if slots are full
func (service *MercuryFullChainService) AddTrustedPeer(url string) error {
	server := service.P2PServer
	if server == nil {
		return errors.New("no p2p server running")
	}

	node, err := enode.ParseV4(url)

	if err != nil {
		return fmt.Errorf("invalid url: %v", err)
	}
	server.AddTrustedPeer(node)
	return nil
}

// RemoveTrustedPeer removes a remote node from the trusted peer set, but it
// does not disconnect it automatically.
func (service *MercuryFullChainService) RemoveTrustedPeer(url string) error {
	server := service.P2PServer
	if server == nil {
		return errors.New("no p2p server running")
	}

	node, err := enode.ParseV4(url)

	if err != nil {
		return fmt.Errorf("invalid url: %v", err)
	}
	server.RemoveTrustedPeer(node)
	return nil
}

func (service *MercuryFullChainService) Peers() ([]*p2p.PeerInfo, error) {
	server := service.P2PServer
	if server == nil {
		return nil, errors.New("no p2p server running")
	}
	return server.PeersInfo(), nil
}

func (service *MercuryFullChainService) GetChainConfig() chain_config.ChainConfig {
	return service.ChainConfig
}

func (service *MercuryFullChainService) GetContractInfo(eData *contract.ExtraDataForContract) (interface{}, error) {
	state, err := service.ChainReader.CurrentState()
	if err != nil {
		return nil, err
	}
	blockHeight := service.ChainReader.CurrentHeader().GetNumber()

	cProcessor := contract.NewProcessor(state, blockHeight)
	//cProcessor := contract.NewProcessor(service.nodeContext.ChainReader(), blockHeight)

	info, err := cProcessor.GetContractReadOnlyInfo(eData)
	return info, err
}

func (service *MercuryFullChainService) GetContract(contractAddr common.Address) (interface{}, error) {
	state, err := service.ChainReader.CurrentState()
	if err != nil {
		return nil, err
	}

	// get contract type
	contractType := contractAddr.GetAddressTypeStr()
	ct, ctErr := contract.GetContractTempByType(contractType)
	if ctErr != nil {
		return nil, ctErr
	}
	nContractV, err := state.GetContract(contractAddr, ct)
	//cb, err := service.nodeContext.ChainReader().GetContract(contractAddr)
	if err != nil {
		return nil, err
	}
	return nContractV.Interface(), nil
}

func (service *MercuryFullChainService) GetBlockDiffVerifierInfo(blockNumber uint64) (map[economy_model.VerifierType][]common.Address, error) {
	if blockNumber < 2 {
		return map[economy_model.VerifierType][]common.Address{}, g_error.BlockNumberError
	}

	block, _ := service.GetBlockByNumber(blockNumber)
	preBlock, _ := service.GetBlockByNumber(blockNumber - 1)
	return service.ChainReader.GetEconomyModel().GetDiffVerifierAddress(preBlock, block)
}

func (service *MercuryFullChainService) GetVerifierDIPReward(blockNumber uint64) (map[economy_model.VerifierType]*big.Int, error) {
	block, _ := service.GetBlockByNumber(blockNumber)
	return service.ChainReader.GetEconomyModel().GetVerifierDIPReward(block)
}

func (service *MercuryFullChainService) GetMineMasterDIPReward(blockNumber uint64) (*big.Int, error) {
	block, _ := service.GetBlockByNumber(blockNumber)
	return service.ChainReader.GetEconomyModel().GetMineMasterDIPReward(block)
}

func (service *MercuryFullChainService) GetBlockYear(blockNumber uint64) (uint64, error) {
	return service.ChainReader.GetEconomyModel().GetBlockYear(blockNumber)
}

func (service *MercuryFullChainService) GetOneBlockTotalDIPReward(blockNumber uint64) (*big.Int, error) {
	if blockNumber == 0 {
		return big.NewInt(0), nil
	}
	return service.ChainReader.GetEconomyModel().GetOneBlockTotalDIPReward(blockNumber)
}

func (service *MercuryFullChainService) GetInvestorInfo() map[common.Address]*big.Int {
	return service.ChainReader.GetEconomyModel().GetInvestorInitBalance()
}

func (service *MercuryFullChainService) GetDeveloperInfo() map[common.Address]*big.Int {
	return service.ChainReader.GetEconomyModel().GetDeveloperInitBalance()
}

func (service *MercuryFullChainService) GetAddressLockMoney(address common.Address) (*big.Int, error) {
	currentBlock := service.CurrentBlock()
	if currentBlock == nil {
		return big.NewInt(0), g_error.BlockIsNilError
	}

	return service.ChainReader.GetEconomyModel().GetAddressLockMoney(address, currentBlock.Number())
}

func (service *MercuryFullChainService) GetInvestorLockDIP(address common.Address, blockNumber uint64) (*big.Int, error) {
	return service.ChainReader.GetEconomyModel().GetInvestorLockDIP(address, blockNumber)
}

func (service *MercuryFullChainService) GetDeveloperLockDIP(address common.Address, blockNumber uint64) (*big.Int, error) {
	return service.ChainReader.GetEconomyModel().GetDeveloperLockDIP(address, blockNumber)
}

func (service *MercuryFullChainService) GetFoundationInfo(usage economy_model.FoundationDIPUsage) map[common.Address]*big.Int {
	return service.ChainReader.GetEconomyModel().GetFoundation().GetFoundationInfo(usage)
}

func (service *MercuryFullChainService) GetMaintenanceLockDIP(address common.Address, blockNumber uint64) (*big.Int, error) {
	return service.ChainReader.GetEconomyModel().GetFoundation().GetMaintenanceLockDIP(address, blockNumber)
}

func (service *MercuryFullChainService) GetReMainRewardLockDIP(address common.Address, blockNumber uint64) (*big.Int, error) {
	return service.ChainReader.GetEconomyModel().GetFoundation().GetReMainRewardLockDIP(address, blockNumber)
}

func (service *MercuryFullChainService) GetEarlyTokenLockDIP(address common.Address, blockNumber uint64) (*big.Int, error) {
	return service.ChainReader.GetEconomyModel().GetFoundation().GetEarlyTokenLockDIP(address, blockNumber)
}

func (service *MercuryFullChainService) GetMineMasterEDIPReward(blockNumber uint64, tokenDecimals int) (*big.Int, error) {
	block, _ := service.GetBlockByNumber(blockNumber)
	DIPReward, err := service.ChainReader.GetEconomyModel().GetMineMasterDIPReward(block)
	if err != nil {
		return nil, err
	}
	return service.ChainReader.GetEconomyModel().GetFoundation().GetMineMasterEDIPReward(DIPReward, blockNumber, tokenDecimals)
}

func (service *MercuryFullChainService) GetVerifierEDIPReward(blockNumber uint64, tokenDecimals int) (map[economy_model.VerifierType]*big.Int, error) {
	block, _ := service.GetBlockByNumber(blockNumber)
	DIPReward, err := service.ChainReader.GetEconomyModel().GetVerifierDIPReward(block)
	if err != nil {
		return map[economy_model.VerifierType]*big.Int{}, err
	}
	return service.ChainReader.GetEconomyModel().GetFoundation().GetVerifierEDIPReward(DIPReward, blockNumber, tokenDecimals)
}

// notify wallet
func (service *MercuryFullChainService) NewBlock(ctx context.Context) (*rpc.Subscription, error) {
	notifier, supported := rpc.NotifierFromContext(ctx)
	if !supported {
		return &rpc.Subscription{}, rpc.ErrNotificationsUnsupported
	}

	rpcSub := notifier.CreateSubscription()

	go func() {
		blockCh := make(chan model.Block)
		//blockSub := service.nodeContext.ChainReader().SubscribeBlockEvent(blockCh)
		blockSub := g_event.Subscribe(g_event.NewBlockInsertEvent, blockCh)

		for {
			select {
			case b := <-blockCh:
				addr := service.GetMineCoinBase
				if !addr.IsEmpty() {
					if b.CoinBaseAddress().IsEqual(addr) {
						if err := notifier.Notify(rpcSub.ID, fmt.Sprintf("mined block: %v", b.Number())); err != nil {
							log.Error("can't notify cli app", "err", err)
						}
					}
				}

			case <-rpcSub.Err():
				blockSub.Unsubscribe()
				return
			case <-notifier.Closed():
				blockSub.Unsubscribe()
				return
			}
		}

	}()
	return rpcSub, nil
}

type SubBlockResp struct {
	Number       uint64            `json:"number"`
	Hash         common.Hash       `json:"hash"`
	CoinBase     common.Address    `json:"coin_base"`
	TimeStamp    *big.Int          `json:"timestamp"  gencodec:"required"`
	Transactions []*SubBlockTxResp `json:"transactions"`
}

type SubBlockTxResp struct {
	TxID         common.Hash     `json:"tx_id"`
	From         common.Address  `json:"from"`
	AccountNonce uint64          `json:"nonce"    gencodec:"required"`
	Recipient    *common.Address `json:"to"       rlp:"nil"` // nil means contract creation
	//HashLock     *common.Hash    `json:"hashLock" rlp:"nil"`
	//TimeLock     *big.Int        `json:"timeLock" gencodec:"required"`
	Amount       *big.Int `json:"value"    gencodec:"required"`
	Fee          *big.Int `json:"fee"      gencodec:"required"`
	ExtraData    []byte   `json:"input"    gencodec:"required"`
	ExtraDataStr string   `json:"input_str"    gencodec:"required"`

	// Signature values
	//R *big.Int `json:"r" gencodec:"required"`
	//S *big.Int `json:"s" gencodec:"required"`
	//V *big.Int `json:"v" gencodec:"required"`
	//// hash_key
	//HashKey []byte `json:"hashKey"    gencodec:"required"`
}

// notify wallet
func (service *MercuryFullChainService) SubscribeBlock(ctx context.Context) (*rpc.Subscription, error) {
	notifier, supported := rpc.NotifierFromContext(ctx)
	if !supported {
		return &rpc.Subscription{}, rpc.ErrNotificationsUnsupported
	}

	rpcSub := notifier.CreateSubscription()

	go func() {
		blockCh := make(chan model.Block)
		//blockSub := service.nodeContext.ChainReader().SubscribeBlockEvent(blockCh)
		blockSub := g_event.Subscribe(g_event.NewBlockInsertEvent, blockCh)

		for {
			select {
			case b := <-blockCh:
				var respTxs []*SubBlockTxResp
				_ = b.TxIterator(func(i int, transaction model.AbstractTransaction) error {
					from, _ := transaction.Sender(nil)
					respTxs = append(respTxs, &SubBlockTxResp{
						TxID:         transaction.CalTxId(),
						From:         from,
						AccountNonce: transaction.Nonce(),
						Recipient:    transaction.To(),
						Amount:       transaction.Amount(),
						ExtraData:    transaction.ExtraData(),
						ExtraDataStr: hexutil.Encode(transaction.ExtraData()),
					})
					return nil
				})

				if err := notifier.Notify(rpcSub.ID, &SubBlockResp{
					Number:       b.Number(),
					Hash:         b.Hash(),
					CoinBase:     b.CoinBaseAddress(),
					TimeStamp:    b.Timestamp(),
					Transactions: respTxs,
				}); err != nil {
					log.Error("can't notify wallet", "err", err)
				}

			case <-rpcSub.Err():
				blockSub.Unsubscribe()
				return
			case <-notifier.Closed():
				blockSub.Unsubscribe()
				return
			}
		}

	}()
	return rpcSub, nil
}

// stop this node service
func (service *MercuryFullChainService) StopDipperin() {
	service.Node.Stop()
	go func() {
		time.Sleep(1 * time.Second)
		os.Exit(0)
	}()
}

func (service *MercuryFullChainService) GetContractAddressByTxHash(txHash common.Hash) (common.Address, error) {
	_, blockHash, blockNumber, _, err := service.Transaction(txHash)
	if err != nil {
		return common.Address{}, err
	}

	receipts := service.ChainReader.GetReceipts(blockHash, blockNumber)
	if receipts == nil {
		return common.Address{}, g_error.ErrReceiptIsNil
	}
	for _, value := range receipts {
		if txHash.IsEqual(value.TxHash) {
			return value.ContractAddress, nil
		}
	}
	return common.Address{}, g_error.ErrReceiptNotFound
}

func (service *MercuryFullChainService) GetABI(contractAddr common.Address) (*utils.WasmAbi, error) {
	stateRoot := service.CurrentBlock().StateRoot()
	stateDB, err := service.ChainReader.AccountStateDB(stateRoot)
	if err != nil {
		return nil, err
	}

	fullState := state_processor.NewFullState(stateDB)
	dataAbi := fullState.GetAbi(contractAddr)

	var abi utils.WasmAbi
	err = abi.FromJson(dataAbi)
	if err != nil {
		return nil, err
	}
	return &abi, nil
}

func (service *MercuryFullChainService) convertReceiptLog(src model2.Receipt) (*model2.Receipt, error) {
	abi, err := service.GetABI(src.ContractAddress)
	if err != nil {
		return nil, err
	}

	logs := make([]*model2.Log, 0)
	for _, value := range src.Logs {
		for _, v := range abi.AbiArr {
			if strings.EqualFold(v.Name, value.TopicName) && strings.EqualFold(v.Type, "event") {
				data, innerErr := utils.ConvertInputs(value.Data, v.Inputs)
				if innerErr != nil {
					return nil, innerErr
				}
				tmpLog := *value
				tmpLog.Data = data
				logs = append(logs, &tmpLog)
				break
			}
		}
	}
	src.Logs = logs
	return &src, nil
}

func (service *MercuryFullChainService) GetConvertReceiptByTxHash(txHash common.Hash) (*model2.Receipt, error) {
	_, blockHash, blockNumber, _, err := service.Transaction(txHash)
	if err != nil {
		return nil, err
	}

	receipts := service.ChainReader.GetReceipts(blockHash, blockNumber)
	if receipts == nil {
		return nil, g_error.ErrReceiptIsNil
	}
	for _, value := range receipts {
		if txHash.IsEqual(value.TxHash) {
			result, innerErr := service.convertReceiptLog(*value)
			if innerErr != nil {
				return nil, innerErr
			}
			return result, nil
		}
	}
	return nil, g_error.ErrReceiptNotFound
}

func (service *MercuryFullChainService) GetTxActualFee(txHash common.Hash) (*big.Int, error) {
	receipt, err := service.GetConvertReceiptByTxHash(txHash)
	if err != nil {
		return nil, err
	}

	tx, _, _, _, err := service.Transaction(txHash)
	if err != nil {
		return nil, err
	}

	actualFee := big.NewInt(0).Mul(big.NewInt(int64(receipt.GasUsed)), tx.GetGasPrice())
	return actualFee, nil
}

func (service *MercuryFullChainService) GetReceiptsByBlockNum(num uint64) (model2.Receipts, error) {
	block, err := service.GetBlockByNumber(num)
	if err != nil {
		return nil, err
	}

	receipts := service.ChainReader.GetReceipts(block.Hash(), block.Number())
	if receipts == nil {
		return nil, g_error.ErrReceiptIsNil
	}
	return receipts, nil
}

func (service *MercuryFullChainService) GetReceiptByTxHash(txHash common.Hash) (*model2.Receipt, error) {
	_, blockHash, blockNumber, _, err := service.Transaction(txHash)
	if err != nil {
		return nil, err
	}

	receipts := service.ChainReader.GetReceipts(blockHash, blockNumber)
	if receipts == nil {
		return nil, g_error.ErrReceiptIsNil
	}
	for _, value := range receipts {
		if txHash.IsEqual(value.TxHash) {
			return value, nil
		}
	}
	return nil, g_error.ErrReceiptNotFound
}

func (service *MercuryFullChainService) SendTransactionContract(from, to common.Address, value, gasLimit, gasPrice *big.Int, data []byte, nonce *uint64) (common.Hash, error) {

	// check Tx type
	if to.GetAddressType() != common.AddressTypeContractCall && to.GetAddressType() != common.AddressTypeContractCreate {
		return common.Hash{}, g_error.ErrInvalidContractType
	}

	extraData, err := service.getExtraData(to, data)
	if err != nil {
		return common.Hash{}, err
	}

	// check constant
	if to.GetAddressType() == common.AddressTypeContractCall {
		constant, _, _, innerErr := service.checkConstant(to, extraData)
		if innerErr != nil {
			return common.Hash{}, innerErr
		}

		if constant {
			return common.Hash{}, g_error.ErrFunctionCalledConstant
		}
	}

	// automatic transfer need this
	if from.IsEqual(common.Address{}) {
		from = service.DefaultAccount
		if from.IsEqual(common.Address{}) {
			return common.Hash{}, errors.New("no default account in this node")
		}
	}

	tmpWallet, usedNonce, err := service.getSendTxInfo(from, nonce)
	if err != nil {
		log.Error("MercuryFullChainService#SendTransactionContract", "err", err)
		return common.Hash{}, err
	}

	tx := model.NewTransactionSc(usedNonce, &to, value, gasPrice, gasLimit.Uint64(), extraData)
	signTx, err := service.signTxAndSend(tmpWallet, from, tx, usedNonce)
	if err != nil {
		pbft_log.Error("send tx error", "txid", tx.CalTxId().Hex(), "err", err)
		log.Error("send tx error", "txid", tx.CalTxId().Hex(), "err", err)
		return common.Hash{}, err
	}

	//log.Info("send transaction", "txId", signTx.CalTxId().Hex())
	//log.Info("send transaction", "gasPrice", signTx.GetGasPrice())
	//log.Info("send transaction", "gas limit", signTx.GetGasLimit())
	/*	signJson, _ := json.Marshal(signTx)
		pbft_log.Info("send transaction", "signTx json", string(signJson))*/
	txHash := signTx.CalTxId()
	log.Info("the SendTransaction txId is: ", "txId", txHash.Hex(), "txSize", signTx.Size())
	return txHash, nil
}

func (service *MercuryFullChainService) getExtraData(to common.Address, data []byte) ([]byte, error) {
	var extraData []byte
	if to.GetAddressType() == common.AddressTypeContractCall {
		state, err := service.ChainReader.CurrentState()
		if err != nil {
			return nil, err
		}

		abi, err := state.GetAbi(to)
		if err != nil {
			log.Error("getExtraData#GetABI failed", "err", err)
			return nil, err
		}

		log.Info("ParseCallContractData")
		extraData, err = utils.ParseCallContractData(abi, data)
		if err != nil {
			log.Error("getExtraData#ParseData failed", "err", err)
			return nil, err
		}
	} else {
		log.Info("ParseCreateContractData")
		var err error
		extraData, err = utils.ParseCreateContractData(data)
		if err != nil {
			return nil, err
		}
	}
	return extraData, nil
}

// CallArgs represents the arguments for a call.
type CallArgs struct {
	From     common.Address  `json:"from"`
	To       *common.Address `json:"to"`
	Gas      hexutil.Uint64  `json:"gas"`
	GasPrice hexutil.Big     `json:"gasPrice"`
	Value    hexutil.Big     `json:"value"`
	Data     hexutil.Bytes   `json:"data"`
}

// CallContract executes the given transaction on the state for the given block number.
// It doesn't make and changes in the state/block chain and is useful to execute and retrieve values.
func (service *MercuryFullChainService) CallContract(from, to common.Address, data []byte, blockNum uint64) (string, error) {

	// check Tx type
	if to.GetAddressType() != common.AddressTypeContractCall {
		return "", g_error.ErrInvalidContractType
	}

	constant, funcName, abi, err := service.checkConstant(to, data)
	if err != nil {
		return "", err
	}

	if !constant {
		return "", g_error.ErrFunctionCalledNotConstant
	}

	var gasLimit uint64
	curBlock := service.CurrentBlock()
	if blockNum == 0 || curBlock.Number() < blockNum {
		gasLimit = curBlock.Header().GetGasLimit()
	} else {
		block, _ := service.GetBlockByNumber(blockNum)
		gasLimit = block.Header().GetGasLimit()
	}

	args := CallArgs{
		From: from,
		To:   &to,
		Data: data,
		Gas:  hexutil.Uint64(gasLimit),
	}
	return service.callContract(args, blockNum, funcName, abi)
}

func (service *MercuryFullChainService) callContract(args CallArgs, blockNum uint64, funcName string, abi *utils.WasmAbi) (string, error) {
	result, _, err := service.doCall(args, blockNum, 5*time.Second)
	if err != nil {
		return "", err
	}

	// convert result by abi
	var resp string
	for _, v := range abi.AbiArr {
		if strings.EqualFold(v.Name, funcName) && strings.EqualFold(v.Type, "function") {
			if len(v.Outputs) != 0 {
				convertResult := utils.Align32BytesConverter(result, v.Outputs[0].Type)
				resp = fmt.Sprintf("%v", convertResult)
			} else {
				resp = "void"
			}
			break
		}

	}
	log.Info("CallContract test", "response", resp)
	return resp, err
}

func (service *MercuryFullChainService) EstimateGas(from, to common.Address, value, gasLimit, gasPrice *big.Int, data []byte, nonce *uint64) (hexutil.Uint64, error) {
	if value == nil {
		value = new(big.Int).SetUint64(0)
	}

	if gasLimit == nil || gasPrice == nil {
		return hexutil.Uint64(0), errors.New("gasLimit or gasPrice are nil")
	}

	extraData, err := service.getExtraData(to, data)
	if err != nil {
		return hexutil.Uint64(0), err
	}

	args := CallArgs{
		From:     from,
		To:       &to,
		Gas:      hexutil.Uint64(gasLimit.Uint64()),
		GasPrice: hexutil.Big(*gasPrice),
		Value:    hexutil.Big(*value),
		Data:     hexutil.Bytes(extraData),
	}
	return service.estimateGas(args)
}

// EstimateGas returns an estimate of the amount of gas needed to execute the
// given transaction against the current block.
func (service *MercuryFullChainService) estimateGas(args CallArgs) (hexutil.Uint64, error) {
	log.Info("estimateGas Start")

	// Binary search the gas requirement, as it may be higher than the amount used
	block := service.CurrentBlock()
	var (
		low      = model2.TxGas - 1
		high     uint64
		capacity uint64
	)
	if uint64(args.Gas) >= model2.TxGas {
		high = uint64(args.Gas)
	} else {
		// Retrieve the current pending block to act as the gas ceiling
		high = block.Header().GetGasLimit()
	}
	capacity = high

	// Create a helper to check if a gas allowance results in an executable transaction
	executable := func(gas uint64) bool {
		args.Gas = hexutil.Uint64(gas)
		_, failed, err := service.doCall(args, block.Number(), 0)
		if err != nil || failed {
			return false
		}
		return true
	}

	// Execute the binary search and hone in on an executable gas limit
	index := 0
	log.Info("executable Start", "low", low, "high", high, "cap", capacity)
	for low+1 < high {
		mid := (high + low) / 2
		if !executable(mid) {
			low = mid
		} else {
			high = mid
		}
		index++
	}
	log.Info("executable End", "times", index, "low", low, "high", high, "cap", capacity)

	// Reject the transaction as invalid if it still fails at the highest allowance
	if high == capacity {
		if !executable(high) {
			return 0, fmt.Errorf("gas required exceeds allowance or always failing transaction")
		}
	}
	return hexutil.Uint64(high), nil
}

func (service *MercuryFullChainService) GetBlockHashByNumber(number uint64) common.Hash {
	if number > service.CurrentBlock().Number() {
		log.Info("GetBlockHashByNumber failed, can't get future block")
		return common.Hash{}
	}
	block, _ := service.GetBlockByNumber(number)
	return block.Hash()
}

func (service *MercuryFullChainService) doCall(args CallArgs, blockNum uint64, timeout time.Duration) ([]byte, bool, error) {
	defer func(start time.Time) { log.Debug("Executing EVM call finished", "runtime", time.Since(start)) }(time.Now())

	var (
		state *state_processor.AccountStateDB
		err   error
	)
	if blockNum == 0 {
		state, err = service.ChainReader.CurrentState()
	} else {
		state, err = service.ChainReader.StateAtByBlockNumber(blockNum)
	}

	if state == nil || err != nil {
		return nil, false, err
	}

	// Set sender address or use a default if none specified
	from := args.From
	if from == (common.Address{}) {
		if wallets := service.WalletManager.Wallets; len(wallets) > 0 {
			a, innerErr := wallets[0].Accounts()
			if innerErr != nil {
				return nil, false, innerErr
			}
			if len(a) > 0 {
				from = a[0].Address
			}
		}
	}

	// Set to address or use a default if none specified
	to := args.To
	if to == nil {
		createAddr := common.HexToAddress(common.AddressContractCreate)
		to = &createAddr
	}

	// Set default gas & gas price if none were set
	gas, gasPrice, value := uint64(args.Gas), args.GasPrice.ToInt(), args.Value.ToInt()
	if gas == 0 {
		gas = math.MaxUint64 / 2
	}
	if gasPrice.Sign() == 0 {
		gasPrice = new(big.Int).SetUint64(uint64(config.DEFAULT_GAS_PRICE))
	}
	if value.Sign() == 0 {
		gasPrice = new(big.Int).SetUint64(uint64(1))
	}

	// Create tmpTransaction
	tmpWallet, usedNonce, err := service.getSendTxInfo(from, nil)
	if err != nil {
		log.Error("doCall#getSendTxInfo failed", "err", err)
		return nil, false, err
	}

	tmpTx := model.NewTransactionSc(usedNonce, to, value, gasPrice, gas, args.Data)
	fromAccount := accounts.Account{Address: from}
	signedTx, err := tmpWallet.SignTx(fromAccount, tmpTx, service.ChainConfig.ChainId)
	if err != nil {
		log.Error("doCall#SignTx failed", "err", err)
		return nil, false, err
	}

	// GetBlock
	block, _ := service.GetBlockByNumber(blockNum)

	// Setup context so it may be cancelled the call has completed
	// or, in case of unmetered gas, setup a context with a timeout.
	ctx := context.Background()
	var cancel context.CancelFunc
	if timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, timeout)
	} else {
		ctx, cancel = context.WithCancel(ctx)
	}
	// Make sure the context is cancelled when the call has completed
	// this makes sure resources are cleaned up.
	defer cancel()

	// Create NewVM
	signedTx.PaddingTxIndex(0)
	conText := vm.NewVMContext(signedTx, block.Header(), service.GetBlockHashByNumber)
	fullState := state_processor.NewFullState(state)
	dvm := vm.NewVM(conText, fullState, vm.DEFAULT_VM_CONFIG)

	/*	// Wait for the context to be done and cancel the evm. Even if the
		// EVM has finished, cancelling may be done (repeatedly)
		go func() {
			<-ctx.Done()
			dvm.Cancel()
		}()*/

	// Setup the gas pool (also for unmetered requests)
	// and apply the message.
	gp := uint64(math.MaxUint64)
	msg, err := signedTx.AsMessage()
	if err != nil {
		log.Error("doCall#AsMessage failed", "err", err)
		return nil, false, err
	}
	result, _, failed, _, err := state_processor.ApplyMessage(dvm, msg, &gp)
	if err != nil {
		log.Error("doCall#ApplyMessage failed", "err", err)
		return result, failed, err
	}
	if failed {
		log.Error("doCall#RunVm failed", "err", err)
		return result, failed, err
	}
	return result, failed, nil
}

func (service *MercuryFullChainService) checkConstant(to common.Address, data []byte) (bool, string, *utils.WasmAbi, error) {
	funcName, err := vm.ParseInputForFuncName(data)
	if err != nil {
		return false, "", nil, err
	}

	// check funcName
	if strings.EqualFold(funcName, "init") {
		log.Debug("checkConstant failed, can't call init function")
		return false, "", nil, g_error.ErrFunctionInitCanNotCalled
	}

	abi, err := service.GetABI(to)
	if err != nil {
		return false, funcName, nil, err
	}

	// check function constant by abi
	for _, v := range abi.AbiArr {
		if strings.EqualFold(v.Name, funcName) && strings.EqualFold(v.Type, "function") {
			if strings.EqualFold(v.Constant, "True") {
				return true, funcName, abi, nil
			} else {
				return false, funcName, abi, nil
			}
			break
		}
	}
	return false, funcName, abi, g_error.ErrFuncNameNotFoundInABI
}
