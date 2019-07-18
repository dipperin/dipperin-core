package state_processor

import (
	"bytes"
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/tests/g-testData"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"math/big"
	"testing"
)

// Test cannot revert after commit
func TestAccountStateDB_RevertToSnapshot1(t *testing.T) {
	db := ethdb.NewMemDatabase()
	tdb := NewStateStorageWithCache(db)

	processor, err := NewAccountStateDB(common.Hash{}, tdb)
	assert.NoError(t, err)

	root := processor.PreStateRoot()
	assert.Equal(t, common.Hash{}, root)

	snapshot := processor.Snapshot()
	err = processor.NewAccountState(aliceAddr)
	assert.NoError(t, err)

	err = processor.AddBalance(aliceAddr, big.NewInt(2000))
	assert.NoError(t, err)
	err = processor.AddNonce(aliceAddr, 10)
	assert.NoError(t, err)

	before, _ := processor.GetBalance(aliceAddr)
	assert.Equal(t, big.NewInt(2000), before)

	processor.RevertToSnapshot(snapshot)
	ba, _ := processor.GetBalance(aliceAddr)
	var nilBigInt *big.Int
	assert.Equal(t, nilBigInt, ba)

	fRoot, err := processor.Finalise()
	assert.NoError(t, err)
	savedRoot, err := processor.Commit()
	assert.Equal(t, fRoot, savedRoot)

}

// Test modification of code and abi
func TestAccountStateDB_RevertToSnapshot2(t *testing.T) {
	db := ethdb.NewMemDatabase()
	tdb := NewStateStorageWithCache(db)

	processor, err := NewAccountStateDB(common.Hash{}, tdb)
	assert.NoError(t, err)

	root := processor.PreStateRoot()
	assert.Equal(t, common.Hash{}, root)

	snapshot := processor.Snapshot()
	err = processor.NewAccountState(aliceAddr)
	assert.NoError(t, err)

	err = processor.AddBalance(aliceAddr, big.NewInt(2000))
	assert.NoError(t, err)
	err = processor.AddNonce(aliceAddr, 10)
	assert.NoError(t, err)

	err = processor.SetCode(aliceAddr, []byte("coooode"))
	code, err := processor.GetCode(aliceAddr)
	assert.Equal(t, []byte("coooode"), code)

	err = processor.SetAbi(aliceAddr, []byte("{input:int}"))
	abi, err := processor.GetAbi(aliceAddr)
	assert.Equal(t, []byte("{input:int}"), abi)

	root, err = processor.GetDataRoot(aliceAddr)
	assert.NoError(t, err)
	assert.Equal(t, common.Hash{}, root)

	processor.SetData(aliceAddr, "tkey", []byte("value"))

	assert.Equal(t, []byte("value"), processor.smartContractData[aliceAddr]["tkey"])

	formerRoot, _ := processor.GetDataRoot(aliceAddr)
	err = processor.finalSmartData()
	root, err = processor.GetDataRoot(aliceAddr)

	assert.NotEqual(t, formerRoot, root)

	processor.RevertToSnapshot(snapshot)
	reverted, _ := processor.GetDataRoot(aliceAddr)

	assert.Equal(t, formerRoot, reverted)
	assert.NotEqual(t, reverted, root)

	abi, err = processor.GetAbi(aliceAddr)
	assert.NotEqual(t, []byte("{input:int}"), abi)

	code, err = processor.GetCode(aliceAddr)
	assert.NotEqual(t, []byte("coooode"), code)
}

// Test the modification of contract data
func TestAccountStateDB_RevertToSnapshot3(t *testing.T) {
	db := ethdb.NewMemDatabase()
	tdb := NewStateStorageWithCache(db)

	processor, err := NewAccountStateDB(common.Hash{}, tdb)
	assert.NoError(t, err)

	root := processor.PreStateRoot()
	assert.Equal(t, common.Hash{}, root)

	err = processor.NewAccountState(aliceAddr)
	assert.NoError(t, err)

	err = processor.AddBalance(aliceAddr, big.NewInt(2000))
	assert.NoError(t, err)
	err = processor.AddNonce(aliceAddr, 10)
	assert.NoError(t, err)

	err = processor.SetCode(aliceAddr, []byte("coooode"))
	err = processor.SetAbi(aliceAddr, []byte("{input:int}"))
	processor.SetData(aliceAddr, "tkey", []byte("value"))

	assert.Equal(t, []byte("value"), processor.smartContractData[aliceAddr]["tkey"])

	err = processor.finalSmartData()
	root, err = processor.GetDataRoot(aliceAddr)

	tr, err := processor.getContractTrie(aliceAddr)
	v, err := tr.TryGet(GetContractFieldKey(aliceAddr, "tkey"))
	assert.Equal(t, v, []byte("value"))
}

func TestAccountStateDB_RevertToSnapshot4(t *testing.T) {
	db := ethdb.NewMemDatabase()
	tdb := NewStateStorageWithCache(db)

	processor, err := NewAccountStateDB(common.Hash{}, tdb)
	err = processor.NewAccountState(aliceAddr)
	assert.NoError(t, err)

	snapshot := processor.Snapshot()
	processor.AddBalance(aliceAddr, big.NewInt(2000))
	processor.AddNonce(aliceAddr, 10)

	processor.SetCode(aliceAddr, []byte("coooode"))
	processor.SetAbi(aliceAddr, []byte("{input:int}"))

	processor.SetData(aliceAddr, "tkey", []byte("value"))
	processor.SetData(aliceAddr, "taaa", []byte("vaaaa"))

	assert.Equal(t, []byte("value"), processor.smartContractData[aliceAddr]["tkey"])
	assert.Equal(t, []byte("vaaaa"), processor.smartContractData[aliceAddr]["taaa"])

	processor.RevertToSnapshot(snapshot)
	assert.Equal(t, true, processor.smartContractData[aliceAddr] == nil)
}

func fakeGetBlockHash(number uint64) common.Hash {
	return common.Hash{}
}
func TestContractCreate(t *testing.T) {
	db := ethdb.NewMemDatabase()
	tdb := NewStateStorageWithCache(db)

	processor, _ := NewAccountStateDB(common.Hash{}, tdb)
	processor.NewAccountState(aliceAddr)
	processor.AddBalance(aliceAddr, big.NewInt(200000000))
	processor.AddNonce(aliceAddr, 11)

	nonce, _ := processor.GetNonce(aliceAddr)
	fmt.Println(nonce)
	tx := FakeContract(t)
	assert.Equal(t, common.AddressTypeContractCreate, int(tx.GetType()))

	//blockGas := uint64(100000000)
	block := model.NewBlock(model.NewHeader(1, 10, common.Hash{}, common.HexToHash("1111"), common.HexToDiff("0x20ffffff"), big.NewInt(324234), common.Address{}, common.BlockNonceFromInt(432423)), nil, nil)
	gasLimit := block.GasLimit()
	gasUsed := block.GasUsed()
	conf := TxProcessConfig{
		Tx:       tx,
		Header:   block.Header().(*model.Header),
		GetHash:  fakeGetBlockHash,
		GasLimit: &gasLimit,
		GasUsed:  &gasUsed,
		TxFee:    big.NewInt(0),
	}

	err := processor.ProcessTxNew(&conf)
	assert.NoError(t, err)
}

func FakeContract(t *testing.T) *model.Transaction {
	codePath := g_testData.GetWasmPath("map-string")
	abiPath := g_testData.GetAbiPath("map-string")
	fileCode, err := ioutil.ReadFile(codePath)
	assert.NoError(t, err)

	fileABI, err := ioutil.ReadFile(abiPath)
	assert.NoError(t, err)
	var input [][]byte
	input = make([][]byte, 0)
	// code
	input = append(input, fileCode)
	// abi
	input = append(input, fileABI)

	buffer := new(bytes.Buffer)
	err = rlp.Encode(buffer, input)

	fs := model.NewMercurySigner(big.NewInt(1))
	to := common.HexToAddress("0x00120000000000000000000000000000000000000000")
	tx := model.NewTransactionSc(uint64(11), &to, big.NewInt(0), big.NewInt(1), uint64(20000000), buffer.Bytes())
	key, _ := createKey()

	log.Info("the tx receipt is:", "to", tx.To().Hex())
	tx.SignTx(key, fs)
	return tx
}
