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

package middleware

import (
	"crypto/ecdsa"
	"errors"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/bloom"
	"github.com/dipperin/dipperin-core/core/chain"
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/dipperin/dipperin-core/core/chain/chaindb"
	"github.com/dipperin/dipperin-core/core/chain/registerdb"
	"github.com/dipperin/dipperin-core/core/economy-model"
	"github.com/dipperin/dipperin-core/third-party/crypto"
	"github.com/dipperin/dipperin-core/third-party/crypto/cs-crypto"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"

	"github.com/dipperin/dipperin-core/core/chain/state-processor"
	"github.com/dipperin/dipperin-core/core/model"
)

func TestNewTxValidatorForRpcService(t *testing.T) {
	v := NewTxValidatorForRpcService(&fakeChainInterface{})
	assert.NotNil(t, v)
	assert.Panics(t, func() {
		v.Valid(&fakeTx{ sender: common.Address{0x11} })
	})

	assert.Error(t, ValidateBlockTxs(&BlockContext{Block: &fakeBlock{}, Chain: &fakeChainInterface{}})())
	assert.NoError(t, ValidateBlockTxs(&BlockContext{Block: &fakeBlock{
		txRoot: common.HexToHash("0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421"),
		isSpecial: true,
	}, Chain: &fakeChainInterface{}})())
	assert.Error(t, ValidateBlockTxs(&BlockContext{Block: &fakeBlock{
		txRoot: common.HexToHash("0xd76a8eabd6e80cb0bcac287d629cc69b498e995847eae3057bc2b36d752d6c63"),
		isSpecial: true,
		txs: []model.AbstractTransaction{&fakeTx{}},
	}, Chain: &fakeChainInterface{}})())

	assert.Error(t, ValidateBlockTxs(&BlockContext{Block: &fakeBlock{
		txRoot: common.HexToHash("0xd76a8eabd6e80cb0bcac287d629cc69b498e995847eae3057bc2b36d752d6c63"),
		isSpecial: false,
		txs: []model.AbstractTransaction{&fakeTx{ fee: big.NewInt(1) }},
	}, Chain: &fakeChainInterface{}})())
}

func TestTxValidatorForRpcService_Valid(t *testing.T) {
	assert.Error(t, ValidTxSize(&fakeTx{size: chain_config.MaxTxSize + 1}))
	assert.NoError(t, ValidTxSize(&fakeTx{}))

	assert.Error(t, ValidTxSender(&fakeTx{
		sender: common.Address{0x11},
		fee: big.NewInt(1),
	}, &fakeChainInterface{}, 0))

	assert.Error(t, ValidTxSender(&fakeTx{
		sender: common.Address{0x11},
		fee: big.NewInt(100000),
	}, &fakeChainInterface{}, 0))

	adb, _ := NewEmptyAccountDB()
	assert.Error(t, ValidTxSender(&fakeTx{
		sender: common.Address{0x11},
		fee: big.NewInt(100000),
	}, &fakeChainInterface{
		state: adb,
	}, 1))

	sender := common.Address{0x11}
	assert.NoError(t, adb.NewAccountState(sender))
	assert.NoError(t, adb.AddBalance(sender, big.NewInt(10000011)))
	assert.Error(t, ValidTxSender(&fakeTx{
		sender: sender,
		fee: big.NewInt(100000),
	}, &fakeChainInterface{
		state: adb,
		block: &fakeBlock{},
		em: &fakeEconomyModel{},
	}, 1))
	assert.Error(t, ValidTxSender(&fakeTx{
		sender: sender,
		fee: big.NewInt(100000),
		amount: big.NewInt(10),
	}, &fakeChainInterface{
		state: adb,
		block: &fakeBlock{},
		em: &fakeEconomyModel{ lockM: big.NewInt(10000011) },
	}, 1))
	assert.NoError(t, ValidTxSender(&fakeTx{
		sender: sender,
		fee: big.NewInt(100000),
		amount: big.NewInt(10),
	}, &fakeChainInterface{
		state: adb,
		block: &fakeBlock{},
		em: &fakeEconomyModel{ lockM: big.NewInt(0) },
	}, 1))
}

func TestValidTxByType(t *testing.T) {
	assert.NoError(t, ValidTxByType(&fakeTx{}, &fakeChainInterface{}, 0)())
	assert.Error(t, ValidTxByType(&fakeTx{ txType: 0x9999 }, &fakeChainInterface{}, 0)())
	assert.Error(t, ValidTxByType(&fakeTx{ txType: common.AddressTypeUnStake }, &fakeChainInterface{}, 0)())
}

func Test_validTx(t *testing.T) {
	_, _, passTx, passChain := getTxTestEnv(t)
	assert.NoError(t, validTx(passTx, passChain, 0))
	passTx.size = chain_config.MaxTxSize + 1
	assert.Error(t, validTx(passTx, passChain, 0))
	passTx.size = 1
	passTx.txType = 0x9999
	assert.Error(t, validTx(passTx, passChain, 0))
	passTx.txType = common.AddressTypeUnStake
	assert.Error(t, validTx(passTx, passChain, 0))
}

func Test_validRegisterTx(t *testing.T) {
	assert.Error(t, validRegisterTx(nil, nil, 0))
}

func Test_validUnStakeTx(t *testing.T) {
	s, adb, passTx, passChain := getTxTestEnv(t)
	assert.NoError(t, adb.AddStake(s, big.NewInt(100)))
	assert.Error(t, validUnStakeTx(passTx, passChain, 0))
}

func Test_validCancelTx(t *testing.T) {
	s, adb, passTx, passChain := getTxTestEnv(t)
	assert.Error(t, validCancelTx(passTx, passChain, 0))
	passTx.sender = common.Address{}
	assert.Error(t, validCancelTx(passTx, passChain, 0))
	passTx.sender = s
	passChain.state = nil
	assert.Error(t, validCancelTx(passTx, passChain, 0))
	passChain.state, _ = NewEmptyAccountDB()
	assert.Error(t, validCancelTx(passTx, passChain, 0))
	passChain.state = adb
	assert.NoError(t, adb.AddStake(s, big.NewInt(11)))
	assert.NoError(t, validCancelTx(passTx, passChain, 0))
	assert.NoError(t, adb.SetLastElect(s, 1))
	assert.Error(t, validCancelTx(passTx, passChain, 0))
}

func Test_validContractTx(t *testing.T) {
	assert.Error(t, validContractTx(&fakeTx{}, &fakeChainInterface{}, 0))
	s, _ := NewEmptyAccountDB()
	assert.Error(t, validContractTx(&fakeTx{}, &fakeChainInterface{ state: s }, 0))
}

func Test_validEarlyTokenTx(t *testing.T) {
	assert.Nil(t, validEarlyTokenTx(nil, nil, 0))
}

func Test_validEvidenceTx(t *testing.T) {
	assert.Error(t, validEvidenceTx(&fakeTx{ extraData: []byte{} }, &fakeChainInterface{}, 0))

	a, p := getPassConflictVote(t)
	pb, err := rlp.EncodeToBytes(p)
	assert.NoError(t, err)
	tmpAddr := a.Address()
	assert.Error(t, validEvidenceTx(&fakeTx{ extraData: pb, to: &tmpAddr }, &fakeChainInterface{}, 0))
}

func Test_conflictVote(t *testing.T) {
	a, p := getPassConflictVote(t)
	pb, err := rlp.EncodeToBytes(p)
	assert.NoError(t, err)
	tmpAddr := common.Address{0x12}
	assert.Error(t, conflictVote(&fakeTx{ extraData: pb, to: &tmpAddr }, &fakeChainInterface{}, 0))

	p.VoteB = a.getVoteMsg(0, 1, common.Hash{}, model.VoteMessage)
	pb, err = rlp.EncodeToBytes(p)
	assert.Error(t, conflictVote(&fakeTx{ extraData: pb }, &fakeChainInterface{}, 0))

	p.VoteB.Height = 3
	pb, err = rlp.EncodeToBytes(p)
	assert.Error(t, conflictVote(&fakeTx{ extraData: pb }, &fakeChainInterface{}, 0))

	p.VoteA.Height = 2
	pb, err = rlp.EncodeToBytes(p)
	assert.Error(t, conflictVote(&fakeTx{ extraData: pb }, &fakeChainInterface{}, 0))

	assert.Error(t, conflictVote(&fakeTx{ extraData: []byte{} }, &fakeChainInterface{}, 0))
}

func Test_validEvidenceTime(t *testing.T) {
	_, adb, passTx, passChain := getTxTestEnv(t)
	to := common.Address{0x12}
	passTx.to = &to
	assert.Error(t, validEvidenceTime(passTx, passChain, 0))

	norTo := cs_crypto.GetNormalAddressFromEvidence(to)
	assert.NoError(t, adb.NewAccountState(norTo))
	assert.NoError(t, adb.SetLastElect(norTo, 1))
	assert.Error(t, validEvidenceTime(passTx, passChain, 0))
}

func Test_validTargetStake(t *testing.T) {
	s, adb, passTx, passChain := getTxTestEnv(t)
	passTx.to = &s
	assert.Error(t, validTargetStake(passTx, passChain, 0))

	target := cs_crypto.GetNormalAddressFromEvidence(s)
	assert.NoError(t, adb.NewAccountState(target))
	assert.Error(t, validTargetStake(passTx, passChain, 0))

	assert.NoError(t, adb.AddStake(target, big.NewInt(10)))
	assert.NoError(t, validTargetStake(passTx, passChain, 0))
}

func Test_validUnStakeTime(t *testing.T) {
	assert.Error(t, validUnStakeTime(&fakeTx{}, &fakeChainInterface{}, 0))
	assert.Error(t, validUnStakeTime(&fakeTx{ sender: common.Address{0x12} }, &fakeChainInterface{}, 0))
	adb, _ := NewEmptyAccountDB()
	assert.Error(t, validUnStakeTime(&fakeTx{ sender: common.Address{0x12} }, &fakeChainInterface{ state: adb }, 0))
	assert.NoError(t, adb.NewAccountState(common.Address{0x12}))
	assert.Error(t, validUnStakeTime(&fakeTx{ sender: common.Address{0x12} }, &fakeChainInterface{ state: adb }, 0))

	assert.NoError(t, adb.AddStake(common.Address{0x12}, big.NewInt(12)))
	assert.NoError(t, adb.SetLastElect(common.Address{0x12}, 12))
	assert.Error(t, validUnStakeTime(&fakeTx{ sender: common.Address{0x12} }, &fakeChainInterface{
		state: adb,
		block: &fakeBlock{},
	}, 0))
}

func getPassConflictVote(t *testing.T) (*Account, model.Proofs) {
	a := NewAccount()
	va := a.getVoteMsg(0, 1, common.Hash{}, model.VoteMessage)
	vb := a.getVoteMsg(0, 1, common.Hash{0x12}, model.VoteMessage)
	p := model.Proofs{
		VoteA: va,
		VoteB: vb,
		VRFHash: common.Hash{0x12},
		Proof: []byte{},
		Priority: 12,
	}
	return a, p
}

func getTxTestEnv(t *testing.T) (common.Address, *state_processor.AccountStateDB, *fakeTx, *fakeChainInterface) {
	s := common.Address{0x11}
	f := big.NewInt(2000)
	adb, storage := NewEmptyAccountDB()
	assert.NoError(t, adb.NewAccountState(s))
	assert.NoError(t, adb.AddBalance(s, big.NewInt(10000011)))
	passTx := &fakeTx{
		sender: s,
		fee: f,
		amount: big.NewInt(10),
	}
	passChain := &fakeChainInterface{
		state: adb,
		block: &fakeBlock{},
		em: &fakeEconomyModel{ lockM: big.NewInt(0) },
		storage: storage,
	}
	return s, adb, passTx, passChain
}

type Account struct {
	Pk      *ecdsa.PrivateKey
	address common.Address
}

func NewAccount() *Account {
	sk, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}
	return &Account{Pk: sk, address: cs_crypto.GetNormalAddress(sk.PublicKey)}
}

func (a *Account) Address() common.Address {
	if !a.address.IsEmpty() {
		return a.address
	}

	a.address = cs_crypto.GetNormalAddress(a.Pk.PublicKey)
	return a.address
}

func (a *Account) SignHash(hash []byte) ([]byte, error) {
	return crypto.Sign(hash, a.Pk)
}

func (a *Account) getVoteMsg(height, round uint64, blockID common.Hash, voteType model.VoteMsgType) *model.VoteMsg {
	v, err := model.NewVoteMsgWithSign(height, round, blockID, voteType, a.SignHash, a.Address())
	if err != nil {
		panic(err)
	}
	return v
}

type fakeChainInterface struct {
	state *state_processor.AccountStateDB
	block *fakeBlock
	em *fakeEconomyModel
	storage state_processor.StateStorage
	slot uint64
	verifiers []common.Address
	cf *chain_config.ChainConfig
}

func (ci *fakeChainInterface) Genesis() model.AbstractBlock {
	panic("implement me")
}

func (ci *fakeChainInterface) CurrentBlock() model.AbstractBlock {
	return ci.block
}

func (ci *fakeChainInterface) CurrentHeader() model.AbstractHeader {
	panic("implement me")
}

func (ci *fakeChainInterface) GetBlock(hash common.Hash, number uint64) model.AbstractBlock {
	panic("implement me")
}

func (ci *fakeChainInterface) GetBlockByHash(hash common.Hash) model.AbstractBlock {
	panic("implement me")
}

func (ci *fakeChainInterface) GetBlockByNumber(number uint64) model.AbstractBlock {
	return ci.block
}

func (ci *fakeChainInterface) HasBlock(hash common.Hash, number uint64) bool {
	panic("implement me")
}

func (ci *fakeChainInterface) GetBody(hash common.Hash) model.AbstractBody {
	panic("implement me")
}

func (ci *fakeChainInterface) GetBodyRLP(hash common.Hash) rlp.RawValue {
	panic("implement me")
}

func (ci *fakeChainInterface) GetHeader(hash common.Hash, number uint64) model.AbstractHeader {
	panic("implement me")
}

func (ci *fakeChainInterface) GetHeaderByHash(hash common.Hash) model.AbstractHeader {
	panic("implement me")
}

func (ci *fakeChainInterface) GetHeaderByNumber(number uint64) model.AbstractHeader {
	panic("implement me")
}

func (ci *fakeChainInterface) GetHeaderRLP(hash common.Hash) rlp.RawValue {
	panic("implement me")
}

func (ci *fakeChainInterface) HasHeader(hash common.Hash, number uint64) bool {
	panic("implement me")
}

func (ci *fakeChainInterface) GetBlockNumber(hash common.Hash) *uint64 {
	panic("implement me")
}

func (ci *fakeChainInterface) GetTransaction(txHash common.Hash) (model.AbstractTransaction, common.Hash, uint64, uint64) {
	panic("implement me")
}

func (ci *fakeChainInterface) GetLatestNormalBlock() model.AbstractBlock {
	return ci.block
}

func (ci *fakeChainInterface) BlockProcessor(root common.Hash) (*chain.BlockProcessor, error) {
	if ci.state == nil {
		return nil, errors.New("failed")
	}
	return chain.NewBlockProcessor(ci, root, ci.storage)
}

func (ci *fakeChainInterface) BlockProcessorByNumber(num uint64) (*chain.BlockProcessor, error) {
	panic("implement me")
}

func (ci *fakeChainInterface) Rollback(target uint64) error {
	return nil
}

func (ci *fakeChainInterface) CurrentSeed() (common.Hash, uint64) {
	panic("implement me")
}

func (ci *fakeChainInterface) IsChangePoint(block model.AbstractBlock, isProcessPackageBlock bool) bool {
	return false
}

func (ci *fakeChainInterface) GetLastChangePoint(block model.AbstractBlock) *uint64 {
	panic("implement me")
}

func (ci *fakeChainInterface) GetSlotByNum(num uint64) *uint64 {
	panic("implement me")
}

func (ci *fakeChainInterface) GetSlot(block model.AbstractBlock) *uint64 {
	return &ci.slot
}

func (ci *fakeChainInterface) GetCurrVerifiers() []common.Address {
	panic("implement me")
}

func (ci *fakeChainInterface) GetVerifiers(round uint64) []common.Address {
	return ci.verifiers
}

func (ci *fakeChainInterface) GetNextVerifiers() []common.Address {
	panic("implement me")
}

func (ci *fakeChainInterface) NumBeforeLastBySlot(slot uint64) *uint64 {
	panic("implement me")
}

func (ci *fakeChainInterface) BuildRegisterProcessor(preRoot common.Hash) (*registerdb.RegisterDB, error) {
	if preRoot.IsEqual(common.Hash{0x12}) {
		return nil, errors.New("failed")
	}
	return registerdb.NewRegisterDB(preRoot, ci.storage, ci)
}

func (ci *fakeChainInterface) GetStateStorage() state_processor.StateStorage {
	panic("implement me")
}

func (ci *fakeChainInterface) CurrentState() (*state_processor.AccountStateDB, error) {
	if ci.state == nil {
		return nil, errors.New("no state")
	}
	return ci.state, nil
}

func (ci *fakeChainInterface) StateAtByBlockNumber(num uint64) (*state_processor.AccountStateDB, error) {
	if ci.state == nil {
		return nil, errors.New("no state")
	}
	return ci.state, nil
}

func (ci *fakeChainInterface) StateAtByStateRoot(root common.Hash) (*state_processor.AccountStateDB, error) {
	panic("implement me")
}

func (ci *fakeChainInterface) BuildStateProcessor(preAccountStateRoot common.Hash) (*state_processor.AccountStateDB, error) {
	panic("implement me")
}

func (ci *fakeChainInterface) GetChainConfig() *chain_config.ChainConfig {
	if ci.cf == nil {
		return chain_config.GetChainConfig()
	}
	return ci.cf
}

func (ci *fakeChainInterface) GetEconomyModel() economy_model.EconomyModel {
	return ci.em
}

func (ci *fakeChainInterface) GetChainDB() chaindb.Database {
	return chaindb.NewChainDB(ethdb.NewMemDatabase(), model.MakeDefaultBlockDecoder())
}

type fakeEconomyModel struct {
	lockM *big.Int
}

func (fe *fakeEconomyModel) GetMineMasterDIPReward(block model.AbstractBlock) (*big.Int, error) {
	panic("implement me")
}

func (fe *fakeEconomyModel) GetVerifierDIPReward(block model.AbstractBlock) (map[economy_model.VerifierType]*big.Int, error) {
	panic("implement me")
}

func (fe *fakeEconomyModel) GetInvestorInitBalance() map[common.Address]*big.Int {
	panic("implement me")
}

func (fe *fakeEconomyModel) GetDeveloperInitBalance() map[common.Address]*big.Int {
	panic("implement me")
}

func (fe *fakeEconomyModel) GetInvestorLockDIP(address common.Address, blockNumber uint64) (*big.Int, error) {
	panic("implement me")
}

func (fe *fakeEconomyModel) GetDeveloperLockDIP(address common.Address, blockNumber uint64) (*big.Int, error) {
	panic("implement me")
}

func (fe *fakeEconomyModel) GetFoundation() economy_model.Foundation {
	panic("implement me")
}

func (fe *fakeEconomyModel) CheckAddressType(address common.Address) economy_model.EconomyModelAddress {
	panic("implement me")
}

func (fe *fakeEconomyModel) GetDiffVerifierAddress(preBlock, block model.AbstractBlock) (map[economy_model.VerifierType][]common.Address, error) {
	panic("implement me")
}

func (fe *fakeEconomyModel) GetAddressLockMoney(address common.Address, blockNumber uint64) (*big.Int, error) {
	if fe.lockM == nil {
		return nil, errors.New("err")
	}
	return fe.lockM, nil
}

func (fe *fakeEconomyModel) GetBlockYear(blockNumber uint64) (uint64, error) {
	panic("implement me")
}

func (fe *fakeEconomyModel) GetOneBlockTotalDIPReward(blockNumber uint64) (*big.Int, error) {
	panic("implement me")
}

type fakeTx struct {
	sender common.Address
	fee *big.Int
	size common.StorageSize
	amount *big.Int
	txType common.TxType
	extraData []byte
	to *common.Address
}

func (ft *fakeTx) AsMessage() (model.Message, error) {
	panic("implement me")
}

func (ft *fakeTx) Size() common.StorageSize {
	return ft.size
}
func (ft *fakeTx) GetGasPrice() *big.Int {
	panic("implement me")
}

func (ft *fakeTx) Amount() *big.Int {
	return ft.amount
}

func (ft *fakeTx) CalTxId() common.Hash {
	return common.Hash{}
}

func (ft *fakeTx) Fee() *big.Int {
	return ft.fee
}

func (ft *fakeTx) Nonce() uint64 {
	panic("implement me")
}

func (ft *fakeTx) To() *common.Address {
	return ft.to
}

func (ft *fakeTx) Sender(singer model.Signer) (common.Address, error) {
	if ft.sender.IsEqual(common.Address{}) {
		return common.Address{}, errors.New("invalid sender")
	}
	return ft.sender, nil
}

func (ft *fakeTx) SenderPublicKey(signer model.Signer) (*ecdsa.PublicKey, error) {
	panic("implement me")
}

func (ft *fakeTx) EncodeRlpToBytes() ([]byte, error) {
	panic("implement me")
}

func (ft *fakeTx) GetSigner() model.Signer {
	return model.NewMercurySigner(big.NewInt(1))
}

func (ft *fakeTx) GetType() common.TxType {
	return ft.txType
}

func (ft *fakeTx) ExtraData() []byte {
	return ft.extraData
}

func (ft *fakeTx) Cost() *big.Int {
	panic("implement me")
}

func (ft *fakeTx) EstimateFee() *big.Int {
	panic("implement me")
}

type fakeBlock struct {
	txRoot common.Hash
	isSpecial bool
	txs []model.AbstractTransaction
	num uint64
	hash common.Hash
	stateRoot common.Hash
	ts *big.Int
	preHash common.Hash
	registerRoot common.Hash
	vs []model.AbstractVerification
	vRoot common.Hash
	diff common.Difficulty
	cb common.Address
	seed common.Hash
	proof []byte
	mPk []byte
	version uint64

	ExtraData []byte
}

func (fb *fakeBlock) Version() uint64 {
	return fb.version
}

func (fb *fakeBlock) Number() uint64 {
	return fb.num
}

func (fb *fakeBlock) IsSpecial() bool {
	return fb.isSpecial
}

func (fb *fakeBlock) Difficulty() common.Difficulty {
	if fb.diff.Equal(common.Difficulty{}) {
		return common.HexToDiff("0x1fffffff")
	}
	return fb.diff
}

func (fb *fakeBlock) PreHash() common.Hash {
	return fb.preHash
}

func (fb *fakeBlock) Seed() common.Hash {
	panic("implement me")
}

func (fb *fakeBlock) RefreshHashCache() common.Hash {
	return fb.preHash
}

func (fb *fakeBlock) Hash() common.Hash {
	return fb.hash
}

func (fb *fakeBlock) EncodeRlpToBytes() ([]byte, error) {
	panic("implement me")
}

func (fb *fakeBlock) TxIterator(cb func(int, model.AbstractTransaction) error) error {
	return nil
}

func (fb *fakeBlock) TxRoot() common.Hash {
	return fb.txRoot
}

func (fb *fakeBlock) Timestamp() *big.Int {
	return fb.ts
}

func (fb *fakeBlock) Nonce() common.BlockNonce {
	panic("implement me")
}

func (fb *fakeBlock) StateRoot() common.Hash {
	return fb.stateRoot
}

func (fb *fakeBlock) SetStateRoot(root common.Hash) {
	panic("implement me")
}

func (fb *fakeBlock) GetRegisterRoot() common.Hash {
	return fb.registerRoot
}

func (fb *fakeBlock) SetRegisterRoot(root common.Hash) {
	panic("implement me")
}

func (fb *fakeBlock) FormatForRpc() interface{} {
	panic("implement me")
}

func (fb *fakeBlock) SetNonce(nonce common.BlockNonce) {
	panic("implement me")
}

func (fb *fakeBlock) CoinBaseAddress() common.Address {
	return fb.cb
}

func (fb *fakeBlock) GetTransactionFees() *big.Int {
	panic("implement me")
}

func (fb *fakeBlock) CoinBase() *big.Int {
	panic("implement me")
}

func (fb *fakeBlock) GetTransactions() []*model.Transaction {
	panic("implement me")
}

func (fb *fakeBlock) GetInterlinks() model.InterLink {
	panic("implement me")
}

func (fb *fakeBlock) SetInterLinkRoot(root common.Hash) {
	panic("implement me")
}

func (fb *fakeBlock) GetInterLinkRoot() (root common.Hash) {
	panic("implement me")
}

func (fb *fakeBlock) SetInterLinks(inter model.InterLink) {
	panic("implement me")
}

func (fb *fakeBlock) GetAbsTransactions() []model.AbstractTransaction {
	return fb.txs
}

func (fb *fakeBlock) GetBloom() iblt.Bloom {
	panic("implement me")
}

func (fb *fakeBlock) Header() model.AbstractHeader {
	return &model.Header{
		RegisterRoot: fb.registerRoot,
		Seed: fb.seed,
		Proof: fb.proof,
		MinerPubKey: fb.mPk,
		Bloom: iblt.NewBloom(model.DefaultBlockBloomConfig),
	}
}

func (fb *fakeBlock) Body() model.AbstractBody {
	return &model.Body{}
}

func (fb *fakeBlock) TxCount() int {
	return 0
}

func (fb *fakeBlock) GetEiBloomBlockData(reqEstimator *iblt.HybridEstimator) *model.BloomBlockData {
	panic("implement me")
}

func (fb *fakeBlock) GetBlockTxsBloom() *iblt.Bloom {
	panic("implement me")
}

func (fb *fakeBlock) VerificationRoot() common.Hash {
	return fb.vRoot
}

func (fb *fakeBlock) SetVerifications(vs []model.AbstractVerification) {
	panic("implement me")
}

func (fb *fakeBlock) VersIterator(func(int, model.AbstractVerification, model.AbstractBlock) error) (error) {
	panic("implement me")
}

func (fb *fakeBlock) GetVerifications() []model.AbstractVerification {
	return fb.vs
}

func NewEmptyAccountDB() (*state_processor.AccountStateDB, state_processor.StateStorage) {
	storage := state_processor.NewStateStorageWithCache(ethdb.NewMemDatabase())
	db, err := state_processor.NewAccountStateDB(common.Hash{}, storage)
	if err != nil {
		panic(err)
	}
	return db, storage
}