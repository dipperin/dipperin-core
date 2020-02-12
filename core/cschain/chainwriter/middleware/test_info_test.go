package middleware

import (
	"crypto/ecdsa"
	"errors"
	"github.com/dipperin/dipperin-core/common"
	iblt "github.com/dipperin/dipperin-core/core/bloom"
	"github.com/dipperin/dipperin-core/core/chain"
	"github.com/dipperin/dipperin-core/core/chain/chaindb"
	"github.com/dipperin/dipperin-core/core/chain/registerdb"
	"github.com/dipperin/dipperin-core/core/chain/stateprocessor"
	"github.com/dipperin/dipperin-core/core/chainconfig"
	"github.com/dipperin/dipperin-core/core/economymodel"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third_party/crypto"
	cs_crypto "github.com/dipperin/dipperin-core/third_party/crypto/cs-crypto"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

var testBlockNum = uint64(10)

func getPassConflictVote() (*Account, model.Proofs) {
	a := NewAccount()
	va := a.getVoteMsg(0, 1, common.Hash{}, model.VoteMessage)
	vb := a.getVoteMsg(0, 1, common.Hash{0x12}, model.VoteMessage)
	p := model.Proofs{
		VoteA:    va,
		VoteB:    vb,
		VRFHash:  common.Hash{0x12},
		Proof:    []byte{},
		Priority: 12,
	}
	return a, p
}

func getTxTestEnv(t *testing.T) (common.Address, *stateprocessor.AccountStateDB, *fakeTx, *fakeChainInterface) {
	s := common.Address{0x11}
	adb, storage := NewEmptyAccountDB()
	assert.NoError(t, adb.NewAccountState(s))
	assert.NoError(t, adb.AddBalance(s, big.NewInt(10000011)))
	passTx := &fakeTx{
		sender:   s,
		GasLimit: 2 * model.TxGas,
		amount:   big.NewInt(10),
	}
	passChain := &fakeChainInterface{
		state:   adb,
		block:   &fakeBlock{num: testBlockNum},
		em:      &fakeEconomyModel{lockM: big.NewInt(0)},
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
	state     *stateprocessor.AccountStateDB
	block     *fakeBlock
	em        *fakeEconomyModel
	storage   stateprocessor.StateStorage
	slot      uint64
	verifiers []common.Address
	cf        *chainconfig.ChainConfig
}

func (ci *fakeChainInterface) GetBloomLog(hash common.Hash, number uint64) model.Bloom {
	panic("implement me")
}

func (ci *fakeChainInterface) GetBloomBits(head common.Hash, bit uint, section uint64) []byte {
	panic("implement me")
}

func (ci *fakeChainInterface) GetReceipts(hash common.Hash, number uint64) model.Receipts {
	panic("implement me")
}

func (ci *fakeChainInterface) GetSeenCommit(height uint64) []model.AbstractVerification {
	panic("implement me")
}

func (ci *fakeChainInterface) SaveBlock(block model.AbstractBlock, seenCommits []model.AbstractVerification) error {
	panic("implement me")
}

func (ci *fakeChainInterface) Genesis() model.AbstractBlock {
	panic("implement me")
}

func (ci *fakeChainInterface) CurrentBlock() model.AbstractBlock {
	return ci.block
}

func (ci *fakeChainInterface) CurrentHeader() model.AbstractHeader {
	return ci.block.Header()
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

func (ci *fakeChainInterface) GetStateStorage() stateprocessor.StateStorage {
	panic("implement me")
}

func (ci *fakeChainInterface) CurrentState() (*stateprocessor.AccountStateDB, error) {
	if ci.state == nil {
		return &stateprocessor.AccountStateDB{}, nil
	}
	return ci.state, nil
}

func (ci *fakeChainInterface) StateAtByBlockNumber(num uint64) (*stateprocessor.AccountStateDB, error) {
	if ci.state == nil {
		return &stateprocessor.AccountStateDB{}, nil
	}
	return ci.state, nil
}

func (ci *fakeChainInterface) StateAtByStateRoot(root common.Hash) (*stateprocessor.AccountStateDB, error) {
	panic("implement me")
}

func (ci *fakeChainInterface) BuildStateProcessor(preAccountStateRoot common.Hash) (*stateprocessor.AccountStateDB, error) {
	panic("implement me")
}

func (ci *fakeChainInterface) GetChainConfig() *chainconfig.ChainConfig {
	if ci.cf == nil {
		return chainconfig.GetChainConfig()
	}
	return ci.cf
}

func (ci *fakeChainInterface) GetEconomyModel() economymodel.EconomyModel {
	return ci.em
}

func (ci *fakeChainInterface) GetChainDB() chaindb.Database {
	return chaindb.NewChainDB(ethdb.NewMemDatabase(), model.MakeDefaultBlockDecoder())
}

func (ci *fakeChainInterface) AccountStateDB(root common.Hash) (*stateprocessor.AccountStateDB, error) {
	aDB, err := stateprocessor.NewAccountStateDB(root, ci.storage)
	if err != nil {
		return nil, err
	}
	
	return aDB, nil
}

type fakeEconomyModel struct {
	lockM *big.Int
}

func (fe *fakeEconomyModel) GetMineMasterDIPReward(block model.AbstractBlock) (*big.Int, error) {
	panic("implement me")
}

func (fe *fakeEconomyModel) GetVerifierDIPReward(block model.AbstractBlock) (map[economymodel.VerifierType]*big.Int, error) {
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

func (fe *fakeEconomyModel) GetFoundation() economymodel.Foundation {
	panic("implement me")
}

func (fe *fakeEconomyModel) CheckAddressType(address common.Address) economymodel.EconomyModelAddress {
	panic("implement me")
}

func (fe *fakeEconomyModel) GetDiffVerifierAddress(preBlock, block model.AbstractBlock) (map[economymodel.VerifierType][]common.Address, error) {
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
	sender    common.Address
	fee       *big.Int
	size      common.StorageSize
	amount    *big.Int
	txType    common.TxType
	extraData []byte
	to        *common.Address
	Price     *big.Int
	GasLimit  uint64
	Receipt   *model.Receipt
}

func (ft *fakeTx) PaddingReceipt(parameters model.ReceiptPara) {
	panic("implement me")
}

func (ft *fakeTx) PaddingActualTxFee(fee *big.Int) {
	panic("implement me")
}

func (ft *fakeTx) GetReceipt() *model.Receipt {
	return ft.Receipt
}

func (ft *fakeTx) GetActualTxFee() (fee *big.Int) {
	panic("implement me")
}

func (ft *fakeTx) GetGasLimit() uint64 {
	return ft.GasLimit
}

func (ft *fakeTx) AsMessage(checkNonce bool) (model.Message, error) {
	panic("implement me")
}

func (ft *fakeTx) Size() common.StorageSize {
	return ft.size
}
func (ft *fakeTx) GetGasPrice() *big.Int {
	return big.NewInt(1)
}

func (ft *fakeTx) Amount() *big.Int {
	return ft.amount
}

func (ft *fakeTx) CalTxId() common.Hash {
	return common.Hash{}
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
	return model.NewSigner(big.NewInt(1))
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
	txRoot       common.Hash
	isSpecial    bool
	txs          []model.AbstractTransaction
	num          uint64
	hash         common.Hash
	stateRoot    common.Hash
	ts           *big.Int
	preHash      common.Hash
	registerRoot common.Hash
	vs           []model.AbstractVerification
	vRoot        common.Hash
	diff         common.Difficulty
	cb           common.Address
	seed         common.Hash
	proof        []byte
	mPk          []byte
	version      uint64
	GasLimit     uint64
	GasUsed      uint64
	ExtraData    []byte
	ReceiptHash  common.Hash
}

func (fb *fakeBlock) GetBloomLog() model.Bloom {
	panic("implement me")
}

func (fb *fakeBlock) SetBloomLog(bloom model.Bloom) {
	panic("implement me")
}

func (fb *fakeBlock) SetReceiptHash(receiptHash common.Hash) {
	panic("implement me")
}

func (fb *fakeBlock) GetReceiptHash() common.Hash {
	return fb.ReceiptHash
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
	for i, tx := range fb.txs {
		if err := cb(i, tx); err != nil {
			return err
		}
	}
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
		Seed:         fb.seed,
		Proof:        fb.proof,
		MinerPubKey:  fb.mPk,
		Bloom:        iblt.NewBloom(model.DefaultBlockBloomConfig),
		GasLimit:     fb.GasLimit,
		GasUsed:      fb.GasUsed,
		ReceiptHash:  fb.ReceiptHash,
		Number:       fb.num,
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
	fb.vRoot = model.DeriveSha(model.Verifications(vs))
}

func (fb *fakeBlock) VersIterator(func(int, model.AbstractVerification, model.AbstractBlock) error) error {
	panic("implement me")
}

func (fb *fakeBlock) GetVerifications() []model.AbstractVerification {
	return fb.vs
}

func NewEmptyAccountDB() (*stateprocessor.AccountStateDB, stateprocessor.StateStorage) {
	storage := stateprocessor.NewStateStorageWithCache(ethdb.NewMemDatabase())
	db, err := stateprocessor.NewAccountStateDB(common.Hash{}, storage)
	if err != nil {
		panic(err)
	}
	return db, storage
}
