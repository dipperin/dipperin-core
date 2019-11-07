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
	"crypto/ecdsa"
	"errors"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/address-util"
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/core/accounts"
	"github.com/dipperin/dipperin-core/core/accounts/soft-wallet"
	"github.com/dipperin/dipperin-core/core/chain"
	"github.com/dipperin/dipperin-core/core/chain-communication"
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/dipperin/dipperin-core/core/chain/cachedb"
	"github.com/dipperin/dipperin-core/core/contract"
	"github.com/dipperin/dipperin-core/core/cs-chain"
	"github.com/dipperin/dipperin-core/core/cs-chain/chain-state"
	"github.com/dipperin/dipperin-core/core/cs-chain/chain-writer"
	"github.com/dipperin/dipperin-core/core/mine/minemaster"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/core/tx-pool"
	"github.com/dipperin/dipperin-core/core/vm/common/utils"
	"github.com/dipperin/dipperin-core/tests"
	"github.com/dipperin/dipperin-core/tests/g-testData"
	"github.com/dipperin/dipperin-core/third-party/crypto"
	"github.com/dipperin/dipperin-core/third-party/p2p"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
	"math/big"
	"net"
	"testing"
)

const (
	testPath   = "/tmp/testSoftWallet"
	url_wrong  = "enode://01010101@123.124.125.126:3"
	Password   = "12345678"
	PassPhrase = "12345678"
)

var (
	alicePriv = "289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232031"
	aliceAddr = common.HexToAddress("0x00005586B883Ec6dd4f8c26063E18eb4Bd228e59c3E9")
	testErr   = errors.New("test error")
)

func createBlock(chain *chain_state.ChainState, txs []*model.Transaction, votes []model.AbstractVerification) model.AbstractBlock {
	key1, _ := crypto.HexToECDSA(alicePriv)
	bb := &tests.BlockBuilder{
		ChainState: chain,
		PreBlock:   chain.CurrentBlock(),
		Txs:        txs,
		Vers:       votes,
		MinerPk:    key1,
	}
	return bb.Build()
}

func createVerifiersVotes(block model.AbstractBlock, votesNum int, testAccounts []tests.Account) (votes []model.AbstractVerification) {
	testVerifierAccounts, _ := tests.ChangeVerifierAddress(testAccounts)
	for i := 0; i < votesNum; i++ {
		voteA := model.NewVoteMsg(block.Number(), uint64(0), block.Hash(), model.VoteMessage)
		sign, _ := crypto.Sign(voteA.Hash().Bytes(), testVerifierAccounts[i].Pk)
		voteA.Witness.Address = testVerifierAccounts[i].Address()
		voteA.Witness.Sign = sign
		votes = append(votes, voteA)
	}
	return
}

func insertBlockToChain(t *testing.T, chain *chain_state.ChainState, num int, txs []*model.Transaction) {
	curNum := int(chain.CurrentBlock().Number())
	config := chain_config.GetChainConfig()
	for i := curNum; i < curNum+num; i++ {
		curBlock := chain.CurrentBlock()
		var block model.AbstractBlock
		if curBlock.Number() == 0 {
			block = createBlock(chain, txs, nil)
		} else {

			// votes for curBlock on chain
			curBlockVotes := createVerifiersVotes(curBlock, config.VerifierNumber*2/3+1, nil)
			block = createBlock(chain, txs, curBlockVotes)
		}

		// votes for build block
		votes := createVerifiersVotes(block, config.VerifierNumber*2/3+1, nil)
		err := chain.SaveBftBlock(block, votes)
		assert.NoError(t, err)
		assert.Equal(t, uint64(i+1), chain.CurrentBlock().Number())
		assert.Equal(t, false, chain.CurrentBlock().IsSpecial())
	}
}

func createERC20() (*model.Transaction, common.Address) {
	erc20 := contract.BuiltInERC20Token{}
	erc20.Owner = chain.VerifierAddress[0]
	erc20.TokenDecimals = 2
	erc20.TokenName = "name"
	erc20.TokenSymbol = "symbol"
	erc20.TokenTotalSupply = big.NewInt(1e4)

	es := util.StringifyJson(erc20)

	extra := contract.ExtraDataForContract{}
	extra.Action = "create"
	extra.Params = es
	contractAdr, _ := address_util.GenERC20Address()
	extra.ContractAddress = contractAdr

	tx := createSignedTx(0, aliceAddr, big.NewInt(0), util.StringifyJsonToBytes(extra), nil)

	return tx, contractAdr
}

func createCsChain(accounts []tests.Account) *chain_state.ChainState {
	f := chain_writer.NewChainWriterFactory()
	chainState := chain_state.NewChainState(&chain_state.ChainStateConfig{
		DataDir:       "",
		WriterFactory: f,
		ChainConfig:   chain_config.GetChainConfig(),
	})
	f.SetChain(chainState)

	// Mainly the default initial verifier is in it, the outer test need to call it to vote on the block
	tests.NewGenesisEnv(chainState.GetChainDB(), chainState.GetStateStorage(), accounts)
	return chainState
}

func createCsChainService(accounts []tests.Account) *cs_chain.CsChainService {
	f := chain_writer.NewChainWriterFactory()
	stateConfig := &chain_state.ChainStateConfig{
		DataDir:       "",
		WriterFactory: f,
		ChainConfig:   chain_config.GetChainConfig(),
	}
	chainState, _ := cs_chain.NewCacheChainState(chain_state.NewChainState(stateConfig))
	f.SetChain(chainState)

	// Mainly the default initial verifier is in it, the outer test need to call it to vote on the block
	tests.NewGenesisEnv(chainState.GetChainDB(), chainState.GetStateStorage(), accounts)

	serviceConfig := &cs_chain.CsChainServiceConfig{
		CacheDB: cachedb.NewCacheDB(chainState.GetDB()),
	}
	return &cs_chain.CsChainService{
		CsChainServiceConfig: serviceConfig,
		CacheChainState:      chainState,
	}
}

func createWalletManager(t *testing.T) *accounts.WalletManager {
	wallet, err := soft_wallet.NewSoftWallet()
	assert.NoError(t, err)

	_, err = wallet.Establish(util.HomeDir()+testPath, "testSoftWallet", Password, PassPhrase)
	assert.NoError(t, err)

	manager, err := accounts.NewWalletManager(&fakeGetAccountInfo{}, wallet)
	assert.NoError(t, err)

	return manager
}

func createWalletIdentifier() *accounts.WalletIdentifier {
	return &accounts.WalletIdentifier{
		WalletType: accounts.SoftWallet,
		Path:       util.HomeDir() + testPath,
		WalletName: "testSoftWallet",
	}
}

func createTxPool(csChain *chain_state.ChainState) *tx_pool.TxPool {
	txConfig := tx_pool.DefaultTxPoolConfig
	txConfig.NoLocals = true
	config := chain_config.GetChainConfig()
	return tx_pool.NewTxPool(txConfig, *config, csChain)
}

func createSignedTx(nonce uint64, to common.Address, amount *big.Int, extraData []byte, testAccounts []tests.Account) *model.Transaction {
	verifiers, _ := tests.ChangeVerifierAddress(testAccounts)
	fs1 := model.NewSigner(big.NewInt(1))
	gasLimit := g_testData.TestGasLimit * 500
	tx := model.NewTransaction(nonce, to, amount, g_testData.TestGasPrice, gasLimit, extraData)
	signedTx, _ := tx.SignTx(verifiers[0].Pk, fs1)
	return signedTx
}

func createSignedTx2(nonce uint64, from *ecdsa.PrivateKey, to common.Address, amount *big.Int) *model.Transaction {
	fs1 := model.NewSigner(big.NewInt(1))
	tx := model.NewTransaction(nonce, to, amount, g_testData.TestGasPrice, g_testData.TestGasLimit, []byte{})
	signedTx, _ := tx.SignTx(from, fs1)
	return signedTx
}

func createContractTx(nonce uint64, WASMPath, AbiPath, input string, testAccounts []tests.Account) *model.Transaction {
	codeByte, abiByte := g_testData.GetCodeAbi(WASMPath, AbiPath)
	var data []byte
	if input == "" {
		data, _ = rlp.EncodeToBytes([]interface{}{codeByte, abiByte})
	} else {
		data, _ = rlp.EncodeToBytes([]interface{}{codeByte, abiByte, input})
	}
	extraData, _ := utils.ParseCreateContractData(data)
	to := common.HexToAddress(common.AddressContractCreate)
	return createSignedTx(nonce, to, g_testData.TestValue, extraData, testAccounts)
}

type fakeValidator struct {
	err error
}

func (v fakeValidator) Valid(tx model.AbstractTransaction) error {
	return v.err
}

type fakeNodeConfig struct {
	nodeType int
}

func (nc fakeNodeConfig) GetNodeType() int {
	return nc.nodeType
}

func (nc fakeNodeConfig) GetIsStartMine() bool {
	return true
}

func (nc fakeNodeConfig) SoftWalletName() string {
	return "name"
}

func (nc fakeNodeConfig) SoftWalletDir() string {
	return "dir"
}

func (nc fakeNodeConfig) GetUploadURL() string {
	panic("implement me")
}

func (nc fakeNodeConfig) GetNodeName() string {
	panic("implement me")
}

func (nc fakeNodeConfig) GetNodeP2PPort() string {
	panic("implement me")
}

func (nc fakeNodeConfig) GetNodeHTTPPort() string {
	panic("implement me")
}

type fakeGetAccountInfo struct{}

func (getAddressInfo *fakeGetAccountInfo) CurrentBalance(address common.Address) *big.Int {
	return big.NewInt(0)
}

func (getAddressInfo *fakeGetAccountInfo) GetTransactionNonce(addr common.Address) (nonce uint64, err error) {
	return 0, nil
}

type fakePeerManager struct{}

func (pm fakePeerManager) GetPeers() map[string]chain_communication.PmAbstractPeer {
	return nil
}

func (pm fakePeerManager) BestPeer() chain_communication.PmAbstractPeer {
	return fakePeer{}
}

func (pm fakePeerManager) IsSync() bool {
	return true
}

func (pm fakePeerManager) GetPeer(id string) chain_communication.PmAbstractPeer {
	panic("implement me")
}

func (pm fakePeerManager) RemovePeer(id string) {
	panic("implement me")
}

type fakePbftNode struct{}

func (pbft fakePbftNode) OnNewWaitVerifyBlock(block model.AbstractBlock, id string) {
	panic("implement me")
}

func (pbft fakePbftNode) OnNewMsg(msg interface{}) error {
	panic("implement me")
}

func (pbft fakePbftNode) ChangePrimary(primary string) {
	panic("implement me")
}

func (pbft fakePbftNode) OnNewP2PMsg(msg p2p.Msg, p chain_communication.PmAbstractPeer) error {
	panic("implement me")
}

func (pbft fakePbftNode) AddPeer(p chain_communication.PmAbstractPeer) error {
	panic("implement me")
}

func (pbft fakePbftNode) OnEnterNewHeight(h uint64) {
	panic("implement me")
}

type fakePeer struct{}

func (peer fakePeer) NodeName() string {
	panic("implement me")
}

func (peer fakePeer) NodeType() uint64 {
	panic("implement me")
}

func (peer fakePeer) SendMsg(msgCode uint64, msg interface{}) error {
	panic("implement me")
}

func (peer fakePeer) ID() string {
	panic("implement me")
}

func (peer fakePeer) ReadMsg() (p2p.Msg, error) {
	panic("implement me")
}

func (peer fakePeer) GetHead() (common.Hash, uint64) {
	return common.HexToHash("123"), uint64(0)
}

func (peer fakePeer) SetHead(head common.Hash, height uint64) {
	panic("implement me")
}

func (peer fakePeer) GetPeerRawUrl() string {
	panic("implement me")
}

func (peer fakePeer) DisconnectPeer() {
	panic("implement me")
}

func (peer fakePeer) RemoteVerifierAddress() (addr common.Address) {
	panic("implement me")
}

func (peer fakePeer) RemoteAddress() net.Addr {
	panic("implement me")
}

func (peer fakePeer) SetRemoteVerifierAddress(addr common.Address) {
	panic("implement me")
}

func (peer fakePeer) SetNodeName(name string) {
	panic("implement me")
}

func (peer fakePeer) SetNodeType(nt uint64) {
	panic("implement me")
}

func (peer fakePeer) SetPeerRawUrl(rawUrl string) {
	panic("implement me")
}

func (peer fakePeer) SetNotRunning() {
	panic("implement me")
}

func (peer fakePeer) IsRunning() bool {
	panic("implement me")
}

func (peer fakePeer) GetCsPeerInfo() *p2p.CsPeerInfo {
	panic("implement me")
}

type fakeMaster struct {
	isMine bool
}

func (m fakeMaster) SetMineGasConfig(gasFloor, gasCeil uint64) {
	panic("implement me")
}

func (m fakeMaster) Start() {
	return
}

func (m fakeMaster) Stop() {
	return
}

func (m fakeMaster) CurrentCoinbaseAddress() common.Address {
	return aliceAddr
}

func (m fakeMaster) SetCoinbaseAddress(addr common.Address) {
	return
}

func (m fakeMaster) OnNewBlock(block model.AbstractBlock) {
	panic("implement me")
}

func (m fakeMaster) Workers() map[minemaster.WorkerId]minemaster.WorkerForMaster {
	panic("implement me")
}

func (m fakeMaster) GetReward(address common.Address) *big.Int {
	panic("implement me")
}

func (m fakeMaster) GetPerformance(address common.Address) uint64 {
	panic("implement me")
}

func (m fakeMaster) Mining() bool {
	return m.isMine
}

func (m fakeMaster) MineTxCount() int {
	return 1
}

func (m fakeMaster) RetrieveReward(address common.Address) {
	panic("implement me")
}

type fakeMasterServer struct{}

func (s fakeMasterServer) RegisterWorker(worker minemaster.WorkerForMaster) {
	return
}

func (s fakeMasterServer) UnRegisterWorker(workerId minemaster.WorkerId) {
	panic("implement me")
}

func (s fakeMasterServer) ReceiveMsg(workerID minemaster.WorkerId, code uint64, msg interface{}) {
	panic("implement me")
}

func (s fakeMasterServer) OnNewMsg(msg p2p.Msg, p chain_communication.PmAbstractPeer) error {
	panic("implement me")
}

func (s fakeMasterServer) SetMineMasterPeer(peer chain_communication.PmAbstractPeer) {
	panic("implement me")
}

type fakeNode struct{}

func (fakeNode) Start() error {
	panic("implement me")
}

func (fakeNode) Stop() {
	return
}

type fakeTxPool struct {
	err error
}

func (pool fakeTxPool) AddRemotes(txs []model.AbstractTransaction) []error {
	var errs []error
	for i := 0; i < len(txs); i++ {
		errs = append(errs, pool.err)
	}
	return errs
}

func (pool fakeTxPool) AddLocals(txs []model.AbstractTransaction) []error {
	var errs []error
	for i := 0; i < len(txs); i++ {
		errs = append(errs, pool.err)
	}
	return errs
}

func (pool fakeTxPool) AddRemote(tx model.AbstractTransaction) error {
	panic("implement me")
}

func (pool fakeTxPool) Stats() (int, int) {
	panic("implement me")
}

type fakeMsgSigner struct{ addr common.Address }

func (f *fakeMsgSigner) SetBaseAddress(address common.Address) {}

func (f *fakeMsgSigner) GetAddress() common.Address {
	return f.addr
}
