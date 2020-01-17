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

package stateprocessor

import (
	"bytes"
	"crypto/ecdsa"
	"errors"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/core/vm"
	common2 "github.com/dipperin/dipperin-core/core/vm/common"
	"github.com/dipperin/dipperin-core/core/vm/common/utils"
	"github.com/dipperin/dipperin-core/third_party/crypto"
	"github.com/dipperin/dipperin-core/third_party/trie"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"math/big"
	"strings"
	"testing"
	"time"
)

var (
	testPriv1 = "289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232031"
	testPriv2 = "289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232032"

	aliceAddr   = common.HexToAddress("0x00005586B883Ec6dd4f8c26063E18eb4Bd228e59c3E9")
	bobAddr     = common.HexToAddress("0x0000970e8128aB834E8EAC17aB8E3812f010678CF791")
	charlieAddr = common.HexToAddress("0x00007dbbf084F4a6CcC070568f7674d4c2CE8CD2709E")

	TrieError = errors.New("trie error")

	testGasPrice = model.TestGasPrice
	testGasLimit = model.TestGasLimit * 100
)

func createKey() (*ecdsa.PrivateKey, *ecdsa.PrivateKey) {
	key1, err1 := crypto.HexToECDSA(testPriv1)
	key2, err2 := crypto.HexToECDSA(testPriv2)
	if err1 != nil || err2 != nil {
		return nil, nil
	}
	return key1, key2
}

func createContractTx(code, abi string, nonce, gasLimit uint64) *model.Transaction {
	key, _ := createKey()
	fs := model.NewSigner(big.NewInt(1))
	data := getContractCode(code, abi)
	to := common.HexToAddress(common.AddressContractCreate)
	tx := model.NewTransaction(nonce, to, big.NewInt(2000000), testGasPrice, gasLimit, data)
	tx.SignTx(key, fs)
	return tx
}

func callContractTx(to *common.Address, funcName string, param [][]byte, nonce uint64) *model.Transaction {
	key, _ := createKey()
	fs := model.NewSigner(big.NewInt(1))
	data := getContractInput(funcName, param)
	tx := model.NewTransaction(nonce, *to, big.NewInt(2000000), testGasPrice, testGasLimit, data)
	tx.SignTx(key, fs)
	return tx
}

func CreateBlock(num uint64, preHash common.Hash, txList []*model.Transaction, limit uint64) *model.Block {
	header := model.NewHeader(1, num, preHash, common.HexToHash("123456"), common.HexToDiff("1fffffff"), big.NewInt(time.Now().UnixNano()), bobAddr, common.BlockNonce{})

	// vote
	var voteList []model.AbstractVerification
	header.GasLimit = limit
	block := model.NewBlock(header, txList, voteList)

	// calculate block nonce
	model.CalNonce(block)
	block.RefreshHashCache()
	return block
}

func CreateTestStateDB() (ethdb.Database, common.Hash) {
	db := ethdb.NewMemDatabase()

	//todo The new method does not take the tree from the underlying database
	tdb := NewStateStorageWithCache(db)
	processor, _ := NewAccountStateDB(common.Hash{}, tdb)
	processor.NewAccountState(aliceAddr)
	processor.NewAccountState(bobAddr)
	processor.AddBalance(aliceAddr, big.NewInt(9e6))

	root, _ := processor.Commit()
	tdb.TrieDB().Commit(root, false)
	return db, root
}

func CreateTestStateDBWithMutableBalance(balance *big.Int) (ethdb.Database, common.Hash) {
	db := ethdb.NewMemDatabase()

	tdb := NewStateStorageWithCache(db)
	processor, _ := NewAccountStateDB(common.Hash{}, tdb)
	processor.NewAccountState(aliceAddr)
	processor.NewAccountState(bobAddr)
	processor.AddBalance(aliceAddr, balance)

	root, _ := processor.Commit()
	tdb.TrieDB().Commit(root, false)
	return db, root
}

func createStateProcessor(t *testing.T) *AccountStateDB {
	db, root := CreateTestStateDB()
	processor, _ := NewAccountStateDB(root, NewStateStorageWithCache(db))
	aliceOriginalStake, _ := processor.GetStake(aliceAddr)
	aliceOriginalBalance, _ := processor.GetBalance(aliceAddr)
	aliceOriginalNonce, _ := processor.GetNonce(aliceAddr)
	aliceOriginalLastElect, _ := processor.GetLastElect(aliceAddr)
	processor.newAccountState(bobAddr)

	assert.EqualValues(t, big.NewInt(0), aliceOriginalStake)
	assert.EqualValues(t, big.NewInt(9000000), aliceOriginalBalance)
	assert.EqualValues(t, uint64(0), aliceOriginalNonce)
	assert.EqualValues(t, uint64(0), aliceOriginalLastElect)
	return processor
}

func createSignedVote(num uint64, blockId common.Hash, voteType model.VoteMsgType, testPriv string, address common.Address) *model.VoteMsg {
	voteA := model.NewVoteMsg(num, num, blockId, voteType)
	hash := common.RlpHashKeccak256(voteA)
	key, _ := crypto.HexToECDSA(testPriv)
	sign, _ := crypto.Sign(hash.Bytes(), key)
	voteA.Witness.Address = address
	voteA.Witness.Sign = sign
	return voteA
}

func getTestVm(db ethdb.Database, root common.Hash) *vm.VM {
	processor, _ := NewAccountStateDB(root, NewStateStorageWithCache(db))
	state := NewFullState(processor)
	return vm.NewVM(vm.Context{
		BlockNumber: big.NewInt(1),
		CanTransfer: vm.CanTransfer,
		Transfer:    vm.Transfer,
		GasLimit:    model.TxGas,
		GetHash:     getTestHashFunc(),
	}, state, common2.DEFAULT_VM_CONFIG)
}

func getTestHashFunc() func(num uint64) common.Hash {
	return func(num uint64) common.Hash {
		return common.Hash{}
	}
}

func getContractCode(code, abi string) []byte {
	fileCode, err := ioutil.ReadFile(code)
	if err != nil {
		log.Error("Read code failed", "err", err)
		return nil
	}

	fileABI, err := ioutil.ReadFile(abi)
	if err != nil {
		log.Error("Read abi failed", "err", err)
		return nil
	}

	var input [][]byte
	input = make([][]byte, 0)
	// code
	input = append(input, fileCode)
	// abi
	input = append(input, fileABI)
	// params

	buffer := new(bytes.Buffer)
	if err = rlp.Encode(buffer, input); err != nil {
		log.Error("RLP encode failed", "err", err)
		return nil
	}
	return buffer.Bytes()
}

func getContractInput(funcName string, param [][]byte) []byte {
	var input [][]byte
	input = make([][]byte, 0)
	// func name
	input = append(input, []byte(funcName))
	// func parameter
	for _, v := range param {
		input = append(input, v)
	}

	buffer := new(bytes.Buffer)
	if err := rlp.Encode(buffer, input); err != nil {
		log.Error("RLP encode failed", "err", err)
		return nil
	}
	return buffer.Bytes()
}

// get Contract data
func getCreateExtraData(wasmPath, abiPath string, params []string) (extraData []byte, err error) {
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
	var args []utils.InputParam
	for _, v := range wasmAbi.AbiArr {
		if strings.EqualFold("init", v.Name) && strings.EqualFold(v.Type, "function") {
			args = v.Inputs
		}
	}
	//params := []string{"dipp", "DIPP", "100000000"}
	wasmBytes, err := ioutil.ReadFile(wasmPath)
	if err != nil {
		return
	}
	rlpParams := []interface{}{
		wasmBytes, abiBytes,
	}
	if len(params) != len(args) {
		return nil, errors.New("vm_utils: length of input and abi not match")
	}
	for i, v := range args {
		bts := params[i]
		re, innerErr := utils.StringConverter(bts, v.Type)
		if innerErr != nil {
			return nil, innerErr
		}
		rlpParams = append(rlpParams, re)
	}
	return rlp.EncodeToBytes(rlpParams)
}

//Get a test transaction
func getTestRegisterTransaction(nonce uint64, key *ecdsa.PrivateKey, amount *big.Int) *model.Transaction {
	trans := model.NewRegisterTransaction(nonce, amount, model.TestGasPrice, model.TestGasLimit)
	fs := model.NewSigner(big.NewInt(1))
	signedTx, _ := trans.SignTx(key, fs)
	return signedTx
}

func getTestCancelTransaction(nonce uint64, key *ecdsa.PrivateKey) *model.Transaction {
	trans := model.NewCancelTransaction(nonce, model.TestGasPrice, model.TestGasLimit)
	fs := model.NewSigner(big.NewInt(1))
	signedTx, _ := trans.SignTx(key, fs)
	return signedTx
}

func getTestUnStakeTransaction(nonce uint64, key *ecdsa.PrivateKey) *model.Transaction {
	trans := model.NewUnStakeTransaction(nonce, model.TestGasPrice, model.TestGasLimit)
	fs := model.NewSigner(big.NewInt(1))
	signedTx, _ := trans.SignTx(key, fs)
	return signedTx
}

func getTestEvidenceTransaction(nonce uint64, key *ecdsa.PrivateKey, target common.Address, voteA, voteB *model.VoteMsg) *model.Transaction {
	trans := model.NewEvidenceTransaction(nonce, model.TestGasPrice, model.TestGasLimit, &target, voteA, voteB)
	fs := model.NewSigner(big.NewInt(1))
	signedTx, _ := trans.SignTx(key, fs)
	return signedTx
}

type fakeStateStorage struct {
	getErr    error
	setErr    error
	passKey   string
	decodeErr bool
}

func (storage fakeStateStorage) OpenTrie(root common.Hash) (StateTrie, error) {
	return fakeTrie{
		getErr:    storage.getErr,
		setErr:    storage.setErr,
		passKey:   storage.passKey,
		decodeErr: storage.decodeErr,
	}, nil
}

func (storage fakeStateStorage) OpenStorageTrie(addrHash, root common.Hash) (StateTrie, error) {
	panic("implement me")
}

func (storage fakeStateStorage) CopyTrie(StateTrie) StateTrie {
	panic("implement me")
}

func (storage fakeStateStorage) TrieDB() *trie.Database {
	panic("implement me")
}

func (storage fakeStateStorage) DiskDB() ethdb.Database {
	return ethdb.NewMemDatabase()
}

type fakeTrie struct {
	getErr    error
	setErr    error
	passKey   string
	decodeErr bool
}

func (trie fakeTrie) TryGet(key []byte) ([]byte, error) {
	if trie.passKey == string(key[22:]) {
		return []byte{128}, nil
	}

	if trie.getErr != nil {
		return []byte{128}, trie.getErr
	}

	if trie.decodeErr {
		return []byte{1, 3}, nil
	}
	return []byte{128}, nil
}

func (trie fakeTrie) TryUpdate(key, value []byte) error {
	return trie.setErr
}

func (trie fakeTrie) TryDelete(key []byte) error {
	return TrieError
}

func (trie fakeTrie) Commit(onleaf trie.LeafCallback) (common.Hash, error) {
	return common.Hash{}, TrieError
}

func (trie fakeTrie) Hash() common.Hash {
	panic("implement me")
}

func (trie fakeTrie) NodeIterator(startKey []byte) trie.NodeIterator {
	panic("implement me")
}

func (trie fakeTrie) GetKey([]byte) []byte {
	panic("implement me")
}

func (trie fakeTrie) Prove(key []byte, fromLevel uint, proofDb ethdb.Putter) error {
	panic("implement me")
}

type erc20 struct {
	// todo special characters cause conversion errors
	Owners  []string            `json:"owne.rs"`
	Balance map[string]*big.Int `json:"balance"`
	Name    string              `json:"name"`
	Name2   string              `json:"name2"`
	Dis     uint64              `json:"dis"`
}

type fakeTransaction struct {
	txType *common.TxType
	nonce  *uint64
	err    error
	sender common.Address
}

func (tx fakeTransaction) PaddingReceipt(parameters model.ReceiptPara) {
	return
}

func (tx fakeTransaction) PaddingActualTxFee(fee *big.Int) {
	return
}

func (tx fakeTransaction) GetReceipt() *model.Receipt {
	panic("implement me")
}

func (tx fakeTransaction) GetActualTxFee() (fee *big.Int) {
	panic("implement me")
}

func (tx fakeTransaction) GetGasLimit() uint64 {
	return model.TestGasLimit
}

func (tx fakeTransaction) AsMessage(checkNonce bool) (model.Message, error) {
	panic("implement me")
}

func (tx fakeTransaction) Size() common.StorageSize {
	panic("implement me")
}

func (tx fakeTransaction) GetGasPrice() *big.Int {
	return model.TestGasPrice
}

func (tx fakeTransaction) Amount() *big.Int {
	return big.NewInt(10000)
}

func (tx fakeTransaction) CalTxId() common.Hash {
	return common.HexToHash("123")
}

func (tx fakeTransaction) Nonce() uint64 {
	if tx.nonce == nil {
		return uint64(0)
	}
	return *tx.nonce
}

func (tx fakeTransaction) To() *common.Address {
	if tx.txType == nil {
		return &bobAddr
	}
	switch *tx.txType {
	case common.AddressTypeStake:
		addr := common.HexToAddress(common.AddressStake)
		return &addr
	case common.AddressTypeUnStake:
		addr := common.HexToAddress(common.AddressUnStake)
		return &addr
	case common.AddressTypeCancel:
		addr := common.HexToAddress(common.AddressCancel)
		return &addr
	case common.AddressTypeEvidence:
		addr := common.HexToAddress("0x0005970e8128aB834E8EAC17aB8E3812f010678CF791")
		return &addr
	default:
		return &bobAddr
	}
}

func (tx fakeTransaction) Sender(singer model.Signer) (common.Address, error) {
	return tx.sender, tx.err
}

func (tx fakeTransaction) SenderPublicKey(signer model.Signer) (*ecdsa.PublicKey, error) {
	panic("implement me")
}

func (tx fakeTransaction) EncodeRlpToBytes() ([]byte, error) {
	panic("implement me")
}

func (tx fakeTransaction) GetSigner() model.Signer {
	panic("implement me")
}

func (tx fakeTransaction) GetType() common.TxType {
	return *tx.txType
}

func (tx fakeTransaction) ExtraData() []byte {
	c := erc20{}
	return util.StringifyJsonToBytes(c)
}

func (tx fakeTransaction) Cost() *big.Int {
	panic("implement me")
}

func (tx fakeTransaction) EstimateFee() *big.Int {
	panic("implement me")
}

//add at 2019.11.28
type fakeTransactionOpt struct {
	txType *common.TxType
	nonce  *uint64
	err    error
	sender common.Address
	amount *big.Int
}

func (tx fakeTransactionOpt) PaddingReceipt(parameters model.ReceiptPara) {
	return
}

func (tx fakeTransactionOpt) PaddingActualTxFee(fee *big.Int) {
	return
}

func (tx fakeTransactionOpt) GetReceipt() *model.Receipt {
	panic("implement me")
}

func (tx fakeTransactionOpt) GetActualTxFee() (fee *big.Int) {
	panic("implement me")
}

func (tx fakeTransactionOpt) GetGasLimit() uint64 {
	return model.TestGasLimit
}

func (tx fakeTransactionOpt) AsMessage(checkNonce bool) (model.Message, error) {
	panic("implement me")
}

func (tx fakeTransactionOpt) Size() common.StorageSize {
	panic("implement me")
}

func (tx fakeTransactionOpt) GetGasPrice() *big.Int {
	return model.TestGasPrice
}

func (tx fakeTransactionOpt) Amount() *big.Int {
	return tx.amount
}

func (tx fakeTransactionOpt) CalTxId() common.Hash {
	return common.HexToHash("123")
}

func (tx fakeTransactionOpt) Nonce() uint64 {
	if tx.nonce == nil {
		return uint64(0)
	}
	return *tx.nonce
}

func (tx fakeTransactionOpt) To() *common.Address {
	if tx.txType == nil {
		return &bobAddr
	}
	switch *tx.txType {
	case common.AddressTypeStake:
		addr := common.HexToAddress(common.AddressStake)
		return &addr
	case common.AddressTypeUnStake:
		addr := common.HexToAddress(common.AddressUnStake)
		return &addr
	case common.AddressTypeCancel:
		addr := common.HexToAddress(common.AddressCancel)
		return &addr
	case common.AddressTypeEvidence:
		addr := common.HexToAddress("0x0005970e8128aB834E8EAC17aB8E3812f010678CF791")
		return &addr
	default:
		return &bobAddr
	}
}

func (tx fakeTransactionOpt) Sender(singer model.Signer) (common.Address, error) {
	return tx.sender, tx.err
}

func (tx fakeTransactionOpt) SenderPublicKey(signer model.Signer) (*ecdsa.PublicKey, error) {
	panic("implement me")
}

func (tx fakeTransactionOpt) EncodeRlpToBytes() ([]byte, error) {
	panic("implement me")
}

func (tx fakeTransactionOpt) GetSigner() model.Signer {
	panic("implement me")
}

func (tx fakeTransactionOpt) GetType() common.TxType {
	return *tx.txType
}

func (tx fakeTransactionOpt) ExtraData() []byte {
	c := erc20{}
	return util.StringifyJsonToBytes(c)
}

func (tx fakeTransactionOpt) Cost() *big.Int {
	panic("implement me")
}

func (tx fakeTransactionOpt) EstimateFee() *big.Int {
	panic("implement me")
}

func fakeGetBlockHash(number uint64) common.Hash {
	return common.Hash{}
}