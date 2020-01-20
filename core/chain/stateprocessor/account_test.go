package stateprocessor

import (
	"bytes"
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"math/big"
	"testing"
)
func TestAccountStateDB_RevertToSnapshot(t *testing.T) {
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

	type result struct {
		isNot bool
	}

	testCases := []struct {
		name   string
		given  func() bool
		expect result
	}{
		{
			name:"Test cannot revert after commit",
			given: func() bool {
				processor.RevertToSnapshot(snapshot)
				fRoot, _ := processor.Finalise()
				savedRoot, _ := processor.Commit()
				return fRoot==savedRoot
			},
			expect:result{true},
		},
		{
			name:"Test modification of code and abi",
			given: func() bool {
				snapshot = processor.Snapshot()
				err = processor.SetCode(aliceAddr, []byte("coooode"))
				err = processor.SetAbi(aliceAddr, []byte("{input:int}"))
				processor.SetData(aliceAddr, "tkey", []byte("value"))
				err = processor.finalSmartData()
				root, err = processor.GetDataRoot(aliceAddr)
				processor.RevertToSnapshot(snapshot)
				abi, _ := processor.GetAbi(aliceAddr)
				code, _ := processor.GetCode(aliceAddr)
				return bytes.Equal([]byte("{input:int}"),abi)&&bytes.Equal([]byte("coooode"), code)
			},
			expect:result{false},
		},
		{
			name:"Test the modification of contract data",
			given: func() bool {
				processor, _ := NewAccountStateDB(common.Hash{}, tdb)
				err = processor.NewAccountState(aliceAddr)
				err = processor.SetCode(aliceAddr, []byte("coooode"))
				err = processor.SetAbi(aliceAddr, []byte("{input:int}"))
				processor.SetData(aliceAddr, "tkey", []byte("value"))
				err = processor.finalSmartData()
				root, err = processor.GetDataRoot(aliceAddr)
				tr, _ := processor.getContractTrie(aliceAddr)
				v, _ := tr.TryGet(GetContractFieldKey(aliceAddr, "tkey"))
				return bytes.Equal(v, []byte("value"))
			},
			expect:result{true},
		},
		{
			name:"Test smartContractData after revert",
			given: func() bool {
				snapshot = processor.Snapshot()
				processor.SetCode(aliceAddr, []byte("coooode"))
				processor.SetAbi(aliceAddr, []byte("{input:int}"))
				processor.SetData(aliceAddr, "tkey", []byte("value"))
				processor.SetData(aliceAddr, "taaa", []byte("vaaaa"))
				processor.RevertToSnapshot(snapshot)
				return processor.smartContractData[aliceAddr] == nil
			},
			expect:result{true},
		},
	}

	for _,tc:=range testCases{
		ret:=tc.given()
		assert.Equal(t,tc.expect.isNot,ret)
	}
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
	codePath := model.GetWASMPath("map-string", model.CoreVmTestData)
	abiPath := model.GetAbiPath("map-string", model.CoreVmTestData)
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

	fs := model.NewSigner(big.NewInt(1))
	to := common.HexToAddress("0x00120000000000000000000000000000000000000000")
	tx := model.NewTransaction(uint64(11), to, big.NewInt(0), big.NewInt(1), uint64(20000000), buffer.Bytes())
	key, _ := createKey()

	tx.SignTx(key, fs)
	return tx
}

