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

package tx_pool

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/consts"
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/core/bloom"
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/dipperin/dipperin-core/core/chain/state-processor"
	"github.com/dipperin/dipperin-core/core/economy-model"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/tests/g-testData"
	"github.com/dipperin/dipperin-core/third-party/crypto/cs-crypto"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/sha3"
	"hash"
	"math/big"
	"runtime"
	"sync"
	"testing"
	"time"
)

var testTxFee = economy_model.GetMinimumTxFee(200)

var threshold = new(big.Int).Div(new(big.Int).Mul(testTxFee, big.NewInt(100+int64(DefaultTxPoolConfig.FeeBump))), big.NewInt(100))

var testRoot = "0x54bbe8ffddc42dd501ab37438c2496d1d3be51d9c562531d56b48ea3bea66708"

// testTxPoolConfig is a transaction pool configuration without stateful disk
// side effects used during testing.
var testTxPoolConfig TxPoolConfig
var ms = model.NewMercurySigner(big.NewInt(1))

func init() {
	testTxPoolConfig = DefaultTxPoolConfig
	testTxPoolConfig.NoLocals = true
	testTxPoolConfig.GlobalSlots = 4096
	testTxPoolConfig.AccountQueue = 512
	testTxPoolConfig.AccountSlots = 16
	testTxPoolConfig.GlobalQueue = 4096
	testTxPoolConfig.Journal = "./locals.out"
}

type testBlockChain struct {
	statedb *state_processor.AccountStateDB
}

func (bc *testBlockChain) CurrentBlock() model.AbstractBlock {
	header := model.NewHeader(0, 0, common.Hash{}, common.Hash{}, common.Difficulty{}, big.NewInt(0), common.Address{}, common.BlockNonce{})
	return model.NewBlock(header, nil, nil)
}

func (bc *testBlockChain) GetBlockByNumber(number uint64) model.AbstractBlock {
	return bc.CurrentBlock()
}

func (bc *testBlockChain) StateAtByStateRoot(root common.Hash) (*state_processor.AccountStateDB, error) {
	return bc.statedb, nil
}

func transaction(nonce uint64, to common.Address, amount *big.Int, gasPrice *big.Int, gasLimit uint64, key *ecdsa.PrivateKey) model.AbstractTransaction {

	uTx := model.NewTransaction(nonce, to, amount, gasPrice, gasLimit, nil)
	tx, _ := uTx.SignTx(key, ms)
	return tx
}

//type fakeValidator struct {
//}
//
//func newFakeValidator() *fakeValidator {
//	return &fakeValidator{}
//}
//
//func (v fakeValidator) Valid(tx model.AbstractTransaction) error {
//	return nil
//}

func createTestStateDB() (ethdb.Database, common.Hash) {
	db := ethdb.NewMemDatabase()
	tdb := state_processor.NewStateStorageWithCache(db)
	teststatedb, _ := state_processor.NewAccountStateDB(common.Hash{}, tdb)

	key1, key2, key3 := createKey()
	aliceAddr := cs_crypto.GetNormalAddress(key1.PublicKey)
	bobAddr := cs_crypto.GetNormalAddress(key2.PublicKey)
	chalieAddr := cs_crypto.GetNormalAddress(key3.PublicKey)
	teststatedb.NewAccountState(aliceAddr)
	teststatedb.NewAccountState(bobAddr)
	teststatedb.NewAccountState(chalieAddr)
	teststatedb.SetNonce(aliceAddr, uint64(20))
	teststatedb.SetNonce(bobAddr, uint64(30))
	teststatedb.SetNonce(chalieAddr, uint64(30))
	teststatedb.SetBalance(aliceAddr, big.NewInt(8400003000))
	teststatedb.SetBalance(bobAddr, big.NewInt(8400003000))
	teststatedb.SetBalance(chalieAddr, big.NewInt(8400003000))
	root, _ := teststatedb.Commit()
	tdb.TrieDB().Commit(root, false)
	return db, root
}

func createTestAddrs(num int) ([]common.Address, []*ecdsa.PrivateKey) {
	keys := createKeyBatch(num)
	addrs := []common.Address{}
	for i := 0; i < num; i++ {
		addr := cs_crypto.GetNormalAddress(keys[i].PublicKey)
		addrs = append(addrs, addr)
	}
	return addrs, keys
}

func createTestStateDBWithBatch(num int) (ethdb.Database, common.Hash) {
	db := ethdb.NewMemDatabase()
	tdb := state_processor.NewStateStorageWithCache(db)
	teststatedb, _ := state_processor.NewAccountStateDB(common.Hash{}, tdb)

	addrs, _ := createTestAddrs(num)
	for i := 0; i < num; i++ {
		teststatedb.NewAccountState(addrs[i])
		teststatedb.SetNonce(addrs[i], uint64(20))
		teststatedb.SetBalance(addrs[i], big.NewInt(1000000))
	}
	root, _ := teststatedb.Commit()
	tdb.TrieDB().Commit(root, false)
	return db, root
}

func setupTxPoolBatch(num int) *TxPool {
	db, root := createTestStateDBWithBatch(num)
	teststatedb, _ := state_processor.NewAccountStateDB(root, state_processor.NewStateStorageWithCache(db))

	blockchain := &testBlockChain{statedb: teststatedb}

	pool := NewTxPool(testTxPoolConfig, chain_config.ChainConfig{ChainId: big.NewInt(1)}, blockchain)

	pool.signer = ms

	return pool
}

func setupTxPool() *TxPool {
	db, root := createTestStateDB()
	teststatedb, _ := state_processor.NewAccountStateDB(root, state_processor.NewStateStorageWithCache(db))
	//con := newFakeValidator()
	blockchain := &testBlockChain{statedb: teststatedb}

	pool := NewTxPool(testTxPoolConfig, chain_config.ChainConfig{ChainId: big.NewInt(1)}, blockchain)

	pool.signer = ms

	return pool
}

// validateTxPoolInternals checks various consistency invariants within the pool.
func validTxPoolInternals(pool *TxPool, t *testing.T) {
	pool.mu.RLock()
	defer pool.mu.RUnlock()

	pending, queued := pool.stats()

	// Ensure the total transaction set is consistent with pending + queued
	assert.Equal(t, pending+queued, pool.all.Count())
	assert.Equal(t, pending+queued, pool.feeList.items.Len()-pool.feeList.stales)

	// Ensure the next nonce to assign is the correct one
	for addr, txs := range pool.pending {
		// Find the last transaction
		txs.txs.Sort()

		last := txs.txs.cache[len(txs.txs.items)-1].Nonce()
		nonce := pool.pendingState.GetNonce(addr)

		assert.NotEqual(t, nonce, last+1)
	}
}

func TestNewTxPool(t *testing.T) {
	db, root := createTestStateDB()
	teststatedb, _ := state_processor.NewAccountStateDB(root, state_processor.NewStateStorageWithCache(db))
	//con := newFakeValidator()
	blockchain := &testBlockChain{statedb: teststatedb}

	config := DefaultTxPoolConfig
	NewTxPool(config, chain_config.ChainConfig{ChainId: big.NewInt(1)}, blockchain)
	assert.NoError(t, nil)

	config.NoLocals = true
	NewTxPool(config, chain_config.ChainConfig{ChainId: big.NewInt(1)}, blockchain)
	assert.NoError(t, nil)
}

func TestEnqueueTx(t *testing.T) {
	//t.Skip()
	pool := setupTxPool()
	key1, key2, _ := createKey()
	bobAddr := cs_crypto.GetNormalAddress(key2.PublicKey)

	//enqueueTx add tx1 in pool queue list
	tx1 := transaction(1, bobAddr, big.NewInt(5000), testTxFee, g_testData.TestGasLimit, key1)
	pool.enqueueTx(tx1.CalTxId(), tx1)
	tx1Get := pool.Get(tx1.CalTxId())
	assert.Equal(t, tx1Get.CalTxId(), tx1.CalTxId())

	//enqueueTx try add tx2 in pool queue list but the fee does not exceed the feeBump ,so enqueueTx failed
	tx2 := transaction(1, bobAddr, big.NewInt(4000), big.NewInt(0).Add(testTxFee, big.NewInt(9)), g_testData.TestGasLimit, key1)
	ok, err := pool.enqueueTx(tx1.CalTxId(), tx2)
	assert.Equal(t, ok, false)
	assert.Error(t, err)
	tx2Get := pool.Get(tx1.CalTxId())
	assert.Equal(t, tx2Get.CalTxId(), tx1.CalTxId())
	tx3Get := pool.Get(tx2.CalTxId())
	assert.Equal(t, tx3Get, nil)

	//enqueueTx try add tx3 in pool queue list and the fee exceeds the feeBump, so the enqueueTx success.
	tx3 := transaction(1, bobAddr, big.NewInt(4000), big.NewInt(0).Add(testTxFee, testTxFee), g_testData.TestGasLimit, key1)
	ok, err = pool.enqueueTx(tx1.CalTxId(), tx3)
	assert.Equal(t, ok, true)
	assert.NoError(t, err)
	tx4Get := pool.Get(tx1.CalTxId())
	assert.Equal(t, tx4Get, nil)
	tx5Get := pool.Get(tx3.CalTxId())
	assert.Equal(t, tx5Get.CalTxId(), tx3.CalTxId())

}

func TestTxPool_validateTx(t *testing.T) {
	pool := setupTxPool()
	key1, key2, _ := createKey()
	aliceAddr := cs_crypto.GetNormalAddress(key1.PublicKey)
	bobAddr := cs_crypto.GetNormalAddress(key2.PublicKey)

	overLoad := make([]byte, 512*1024*2)
	for i := 0; i < 512*1024*2; i++ {
		overLoad[i] = byte(i)
	}
	unsignedTx1 := model.NewTransaction(1, bobAddr, big.NewInt(5000), g_testData.TestGasPrice, g_testData.TestGasLimit, overLoad)
	signedTx1, err := unsignedTx1.SignTx(key1, ms)
	assert.NoError(t, err)
	err = pool.validateTx(signedTx1, false)
	assert.Error(t, err)

	signedTx2 := transaction(1, bobAddr, big.NewInt(5000), testTxFee, g_testData.TestGasLimit, key1)
	err = pool.validateTx(signedTx2, false)
	assert.EqualError(t, err, "tx nonce is invalid")

	unsignedTx3 := model.NewTransaction(1, bobAddr, big.NewInt(5000), g_testData.TestGasPrice, g_testData.TestGasLimit, nil)
	signedTx3, err := unsignedTx3.SignTx(key1, model.NewMercurySigner(big.NewInt(2)))
	assert.NoError(t, err)
	err = pool.validateTx(signedTx3, false)
	assert.EqualError(t, err, "invalid sender")

	signedTx4 := transaction(1, bobAddr, big.NewInt(5000), big.NewInt(1), g_testData.TestGasLimit, key1)
	err = pool.validateTx(signedTx4, true)
	assert.EqualError(t, err, "tx nonce is invalid")

	signedTx5 := transaction(40, bobAddr, big.NewInt(1000000), big.NewInt(0).Mul(testTxFee, big.NewInt(10000)), g_testData.TestGasLimit, key1)
	curBalance, e := pool.currentState.GetBalance(aliceAddr)
	err = pool.validateTx(signedTx5, false)
	assert.EqualError(t, err, fmt.Sprintf("tx exceed balance limit, from:%v, cur balance:%v, cost:%v, err:%v", aliceAddr.Hex(), curBalance.String(), signedTx5.Cost().String(), e))
}

func TestTxPool_Pending(t *testing.T) {
	pool := setupTxPool()
	_, err := pool.Pending()
	assert.NoError(t, err)
}

func TestTxPool_TxDifference(t *testing.T) {
	key1, key2, _ := createKey()
	aliceAddr := cs_crypto.GetNormalAddress(key1.PublicKey)
	bobAddr := cs_crypto.GetNormalAddress(key2.PublicKey)

	alicetx1 := transaction(1, bobAddr, big.NewInt(5000), testTxFee, g_testData.TestGasLimit, key1)
	alicetx2 := transaction(2, bobAddr, big.NewInt(4000), big.NewInt(0).Add(testTxFee, big.NewInt(9)), g_testData.TestGasLimit, key1)
	alicetx3 := transaction(3, bobAddr, big.NewInt(6000), big.NewInt(0).Add(testTxFee, big.NewInt(10)), g_testData.TestGasLimit, key1)

	bobtx1 := transaction(4, aliceAddr, big.NewInt(5000), testTxFee, g_testData.TestGasLimit, key2)
	bobtx2 := transaction(5, aliceAddr, big.NewInt(4000), big.NewInt(0).Add(testTxFee, big.NewInt(9)), g_testData.TestGasLimit, key2)
	bobtx3 := transaction(6, aliceAddr, big.NewInt(6000), big.NewInt(0).Add(testTxFee, big.NewInt(10)), g_testData.TestGasLimit, key2)

	txListA := []model.AbstractTransaction{alicetx1, alicetx2, alicetx3, bobtx1, bobtx2}
	txListB := []model.AbstractTransaction{alicetx2, alicetx3, bobtx1, bobtx2, bobtx3}

	difference := model.TxDifference(txListA, txListB)
	//fmt.Println(difference)
	assert.Equal(t, 1, len(difference))
	assert.Equal(t, alicetx1, difference[0])
	difference = model.TxDifference(txListB, txListA)
	//fmt.Println(difference)
	assert.Equal(t, 1, len(difference))
	assert.Equal(t, bobtx3, difference[0])
}

func TestPromoteTx(t *testing.T) {
	//t.Skip()
	pool := setupTxPool()
	key1, key2, _ := createKey()
	aliceAddr := cs_crypto.GetNormalAddress(key1.PublicKey)
	bobAddr := cs_crypto.GetNormalAddress(key2.PublicKey)

	alicetx1 := transaction(0, bobAddr, big.NewInt(5000), testTxFee, g_testData.TestGasLimit, key1)
	alicetx2 := transaction(0, bobAddr, big.NewInt(4000), big.NewInt(0).Add(testTxFee, big.NewInt(9)), g_testData.TestGasLimit, key1)
	alicetx3 := transaction(0, bobAddr, big.NewInt(6000), big.NewInt(0).Add(testTxFee, testTxFee), g_testData.TestGasLimit, key1)

	ok := pool.promoteTx(aliceAddr, alicetx1.CalTxId(), alicetx1)
	assert.Equal(t, ok, true)
	tx1Get := pool.Get(alicetx1.CalTxId())
	assert.Equal(t, tx1Get.CalTxId(), alicetx1.CalTxId())

	ok = pool.promoteTx(aliceAddr, alicetx2.CalTxId(), alicetx2)
	assert.Equal(t, ok, false)
	tx2Get := pool.Get(alicetx1.CalTxId())
	assert.Equal(t, tx2Get.CalTxId(), alicetx1.CalTxId())

	ok = pool.promoteTx(aliceAddr, alicetx3.CalTxId(), alicetx3)
	assert.Equal(t, ok, true)
	tx3Get := pool.Get(alicetx3.CalTxId())
	assert.Equal(t, tx3Get.CalTxId(), alicetx3.CalTxId())

	assert.Equal(t, pool.pendingState.GetNonce(aliceAddr), uint64(20))
}

func TestTxPool_Add(t *testing.T) {
	pool := setupTxPool()
	fmt.Println(pool.all.Count(), pool.config.GlobalSlots, pool.config.GlobalQueue)

	key1, key2, _ := createKey()
	aliceAddr := cs_crypto.GetNormalAddress(key1.PublicKey)
	bobAddr := cs_crypto.GetNormalAddress(key2.PublicKey)
	fmt.Println(pool.currentState.GetNonce(aliceAddr))
	fmt.Println(pool.currentState.GetNonce(bobAddr))

	alicetx1 := transaction(19, bobAddr, big.NewInt(100000), testTxFee, g_testData.TestGasLimit, key1)
	alicetx2 := transaction(20, bobAddr, big.NewInt(4000), testTxFee, g_testData.TestGasLimit, key1)
	bobtx1 := transaction(30, aliceAddr, big.NewInt(3000), testTxFee, g_testData.TestGasLimit, key2)
	bobtx2 := transaction(31, aliceAddr, big.NewInt(3000), testTxFee, g_testData.TestGasLimit, key2)
	bobtx3 := transaction(33, aliceAddr, big.NewInt(3000), testTxFee, g_testData.TestGasLimit, key2)

	//pool.add(alicetx1,false)
	//pool.add(alicetx2,false)
	//pool.add(bobtx1,false)
	//pool.add(bobtx2,false)
	//pool.add(bobtx3,false)
	//
	//fmt.Println(pool.stats())
	//pool.promoteExecutables(nil)
	//fmt.Println(pool.stats())

	pool.AddRemote(alicetx1)
	pool.AddRemote(alicetx2)
	pool.AddRemote(bobtx1)

	pool.AddRemote(bobtx3)

	pool.AddRemote(bobtx2)

	fmt.Println(pool.stats())
	pendTx, _ := pool.Pending()
	fmt.Println(pendTx)
	qTx, _ := pool.Queueing()
	fmt.Println(qTx)

}

func TestTxPool_removeTx(t *testing.T) {
	pool := setupTxPool()
	key1, key2, _ := createKey()
	aliceAddr := cs_crypto.GetNormalAddress(key1.PublicKey)
	bobAddr := cs_crypto.GetNormalAddress(key2.PublicKey)

	bobtx1 := transaction(30, aliceAddr, big.NewInt(3000), testTxFee, g_testData.TestGasLimit, key2)
	bobtx2 := transaction(31, aliceAddr, big.NewInt(3000), testTxFee, g_testData.TestGasLimit, key2)
	bobtx3 := transaction(32, aliceAddr, big.NewInt(3000), testTxFee, g_testData.TestGasLimit, key2)

	err := pool.AddRemote(bobtx1)
	assert.NoError(t, err)
	err = pool.AddRemote(bobtx2)
	assert.NoError(t, err)
	err = pool.AddRemote(bobtx3)
	assert.NoError(t, err)
	pend, queue := pool.stats()
	assert.Equal(t, 3, pend)
	assert.Equal(t, 0, queue)

	pool.removeTx(bobtx1.CalTxId(), true)
	pending, err := pool.Pending()
	assert.NoError(t, err)
	queueing, err := pool.Queueing()
	assert.NoError(t, err)
	assert.Equal(t, []model.AbstractTransaction(nil), pending[bobAddr])
	assert.Equal(t, bobtx2.CalTxId(), queueing[bobAddr][0].CalTxId())
	assert.Equal(t, bobtx3.CalTxId(), queueing[bobAddr][1].CalTxId())
}

func TestTxPool_RemoveTx(t *testing.T) {
	pool := setupTxPool()
	key1, key2, _ := createKey()
	aliceAddr := cs_crypto.GetNormalAddress(key1.PublicKey)

	bobtx1 := transaction(30, aliceAddr, big.NewInt(3000), testTxFee, g_testData.TestGasLimit, key2)
	bobtx2 := transaction(31, aliceAddr, big.NewInt(3000), testTxFee, g_testData.TestGasLimit, key2)
	bobtx3 := transaction(32, aliceAddr, big.NewInt(3000), testTxFee, g_testData.TestGasLimit, key2)

	pool.AddRemote(bobtx1)
	pool.AddRemote(bobtx2)
	pool.AddRemote(bobtx3)
	pend, queue := pool.stats()
	assert.Equal(t, 3, pend)
	assert.Equal(t, 0, queue)

	header := model.NewHeader(0, 0, common.Hash{}, common.Hash{}, common.Difficulty{}, big.NewInt(0), common.Address{}, common.BlockNonce{})
	trans := make([]*model.Transaction, 3)
	util.InterfaceSliceCopy(trans, []model.AbstractTransaction{bobtx1, bobtx2, bobtx3})
	block := model.NewBlock(header, trans, []model.AbstractVerification{})

	pool.RemoveTxs(block)
	pend, queue = pool.stats()
	assert.Equal(t, 0, pend)
	assert.Equal(t, 0, queue)
}

func TestTxPool_RemoveTxsBatch(t *testing.T) {
	pool := setupTxPool()
	key1, key2, _ := createKey()
	aliceAddr := cs_crypto.GetNormalAddress(key1.PublicKey)
	txs := []model.AbstractTransaction{}
	txIds := make([]common.Hash, 0)
	for i := 0; i < 400; i++ {
		tx := transaction(uint64(30+i), aliceAddr, big.NewInt(1), testTxFee, g_testData.TestGasLimit, key2)
		txs = append(txs, tx)
		txIds = append(txIds, tx.CalTxId())
	}
	pool.AddRemotes(txs)
	pend, queue := pool.stats()
	assert.Equal(t, 400, pend)
	assert.Equal(t, 0, queue)

	pool.RemoveTxsBatch(txIds[:200])
	pend, queue = pool.stats()
	//fmt.Println(pend, queue)
	assert.Equal(t, 0, pend)
	assert.Equal(t, 200, queue)
}

func TestTxPool_loop(t *testing.T) {
	pool := setupTxPool()
	pool.Start()
	time.Sleep(3 * time.Second)
	pool.Stop()
}

func TestTxPoolAdd(t *testing.T) {
	pool := setupTxPool()
	fmt.Println(pool.all.Count(), pool.config.GlobalSlots, pool.config.GlobalQueue)

	key1, key2, _ := createKey()
	aliceAddr := cs_crypto.GetNormalAddress(key1.PublicKey)
	bobAddr := cs_crypto.GetNormalAddress(key2.PublicKey)
	fmt.Println(pool.currentState.GetNonce(aliceAddr))
	fmt.Println(pool.currentState.GetNonce(bobAddr))

	alicetx1 := transaction(19, bobAddr, big.NewInt(100000), testTxFee, g_testData.TestGasLimit, key1)
	alicetx2 := transaction(20, bobAddr, big.NewInt(4000), testTxFee, g_testData.TestGasLimit, key1)
	bobtx1 := transaction(30, aliceAddr, big.NewInt(3000), testTxFee, g_testData.TestGasLimit, key2)
	bobtx2 := transaction(31, aliceAddr, big.NewInt(3000), testTxFee, g_testData.TestGasLimit, key2)
	bobtx3 := transaction(33, aliceAddr, big.NewInt(3000), testTxFee, g_testData.TestGasLimit, key2)

	//pool.add(alicetx1,false)
	//pool.add(alicetx2,false)
	//pool.add(bobtx1,false)
	//pool.add(bobtx2,false)
	//pool.add(bobtx3,false)
	//
	//fmt.Println(pool.stats())
	//pool.promoteExecutables(nil)
	//fmt.Println(pool.stats())

	pool.AddLocal(alicetx1)
	pool.AddLocal(alicetx2)
	pool.AddLocal(bobtx1)

	pool.AddLocal(bobtx3)

	pool.AddLocal(bobtx2)

}

func TestStatus(t *testing.T) {
	pool := setupTxPool()
	key1, key2, _ := createKey()
	aliceAddr := cs_crypto.GetNormalAddress(key1.PublicKey)
	bobtx1 := transaction(30, aliceAddr, big.NewInt(3000), testTxFee, g_testData.TestGasLimit, key2)
	bobtx2 := transaction(31, aliceAddr, big.NewInt(3000), testTxFee, g_testData.TestGasLimit, key2)
	bobtx3 := transaction(33, aliceAddr, big.NewInt(3000), testTxFee, g_testData.TestGasLimit, key2)
	pool.AddRemote(bobtx1)
	pool.AddRemote(bobtx3)
	bobTxHashes := []common.Hash{bobtx2.CalTxId(), bobtx1.CalTxId(), bobtx3.CalTxId()}
	result := pool.Status(bobTxHashes)
	assert.Equal(t, TxStatusUnknown, result[0])
	assert.Equal(t, TxStatusPending, result[1])
	assert.Equal(t, TxStatusQueued, result[2])
}

func TestTxPool_AddLocal(t *testing.T) {
	pool := setupTxPool()
	key1, key2, _ := createKey()
	aliceAddr := cs_crypto.GetNormalAddress(key1.PublicKey)

	tx := transaction(uint64(30), aliceAddr, big.NewInt(1), testTxFee, g_testData.TestGasLimit, key2)
	err := pool.AddLocal(tx)
	assert.NoError(t, err)
}

func TestTxPool_LocalAdd(t *testing.T) {
	pool := setupTxPool()
	key1, key2, _ := createKey()
	aliceAddr := cs_crypto.GetNormalAddress(key1.PublicKey)

	//test nonce < record nonce
	for i := 0; i < 10; i++ {
		tx := transaction(uint64(20+i), aliceAddr, big.NewInt(1), testTxFee, g_testData.TestGasLimit, key2)
		_, err := pool.add(tx, false)
		//assert.True(t, ok)
		assert.Error(t, err)
	}

	//make pool full
	for i := 0; i < 5120; i++ {
		tx := transaction(uint64(30+i), aliceAddr, big.NewInt(1), testTxFee, g_testData.TestGasLimit, key2)
		_, err := pool.add(tx, false)
		//assert.True(t, ok)
		assert.NoError(t, err)
	}
	pendN, queueN := pool.stats()
	assert.Equal(t, 0, pendN)
	assert.Equal(t, 5120, queueN)

	//add one more transaction
	txMore := transaction(uint64(30+5120), aliceAddr, big.NewInt(1), testTxFee, g_testData.TestGasLimit, key2)
	_, err := pool.add(txMore, false)
	assert.NoError(t, err)
}

func TestTxPool_promptExecutables_2addr(t *testing.T) {
	pool := setupTxPool()
	key1, key2, key3 := createKey()
	aliceAddr := cs_crypto.GetNormalAddress(key1.PublicKey)
	bobAddr := cs_crypto.GetNormalAddress(key2.PublicKey)
	chalieAddr := cs_crypto.GetNormalAddress(key3.PublicKey)

	for i := 0; i < 30; i++ {
		bobtx := transaction(uint64(30+i), aliceAddr, big.NewInt(1), testTxFee, g_testData.TestGasLimit, key2)
		//fmt.Println(bobtx.Sender(pool.signer))
		_, err := pool.add(bobtx, false)
		//assert.True(t, ok)
		assert.NoError(t, err)
	}

	for i := 0; i < 33; i++ {
		alicetx := transaction(uint64(20+i), chalieAddr, big.NewInt(1), testTxFee, g_testData.TestGasLimit, key1)
		//fmt.Println(bobtx.Sender(pool.signer))
		_, err := pool.add(alicetx, false)
		//assert.True(t, ok)
		assert.NoError(t, err)
	}

	for i := 0; i < 33; i++ {
		chalietx := transaction(uint64(30+i), bobAddr, big.NewInt(1), testTxFee, g_testData.TestGasLimit, key3)
		//fmt.Println(bobtx.Sender(pool.signer))
		_, err := pool.add(chalietx, false)
		//assert.True(t, ok)
		assert.NoError(t, err)
	}
	fmt.Println(pool.stats())

	pool.promoteExecutables(nil)
	pendN, queueN := pool.stats()
	fmt.Println(pendN, queueN)
	//fmt.Println(pool.Pending())
}

func TestTxPool_promptExecutables(t *testing.T) {
	pool := setupTxPool()
	key1, key2, _ := createKey()
	aliceAddr := cs_crypto.GetNormalAddress(key1.PublicKey)

	//transaction enqueue first
	for i := 0; i < 10; i++ {
		bobtx := transaction(uint64(30+i), aliceAddr, big.NewInt(1), testTxFee, g_testData.TestGasLimit, key2)
		//fmt.Println(bobtx.Sender(pool.signer))
		_, err := pool.add(bobtx, false)
		//assert.True(t, ok)
		assert.NoError(t, err)
	}
	pendN, queueN := pool.stats()
	assert.Equal(t, 0, pendN)
	assert.Equal(t, 10, queueN)

	pool.promoteExecutables(nil)
	pendN, queueN = pool.stats()
	assert.Equal(t, 10, pendN)
	assert.Equal(t, 0, queueN)

	//add 4096 more tx
	for i := 0; i < 4096; i++ {
		bobtx := transaction(uint64(40+i), aliceAddr, big.NewInt(1), testTxFee, g_testData.TestGasLimit, key2)
		//fmt.Println(bobtx.Sender(pool.signer))
		_, err := pool.add(bobtx, false)
		//assert.True(t, ok)
		assert.NoError(t, err)
	}
	pool.promoteExecutables(nil)
	pendN, queueN = pool.stats()
	assert.Equal(t, 4096, pendN)
	assert.Equal(t, 0, queueN)

	//add 1024 more tx
	for i := 0; i < 1024; i++ {
		bobtx := transaction(uint64(4136+i), aliceAddr, big.NewInt(1), testTxFee, g_testData.TestGasLimit, key2)
		//fmt.Println(bobtx.Sender(pool.signer))
		_, err := pool.add(bobtx, false)
		//assert.True(t, ok)
		assert.NoError(t, err)
	}
	pendN, queueN = pool.stats()
	fmt.Println(pendN, queueN)
	pool.promoteExecutables(nil)
	pendN, queueN = pool.stats()
	assert.Equal(t, int(testTxPoolConfig.GlobalQueue), pendN)
	assert.Equal(t, int(testTxPoolConfig.AccountQueue), queueN)
}

func TestTxPool_AddWithLocal(t *testing.T) {
	pool := setupTxPool()
	key1, key2, _ := createKey()
	aliceAddr := cs_crypto.GetNormalAddress(key1.PublicKey)

	//transaction enqueue first
	for i := 0; i < 10; i++ {
		bobtx := transaction(uint64(30+i), aliceAddr, big.NewInt(1), testTxFee, g_testData.TestGasLimit, key2)
		//fmt.Println(bobtx.Sender(pool.signer))
		//add bob address to local too
		_, err := pool.add(bobtx, true)
		//assert.True(t, ok)
		assert.NoError(t, err)
	}
	pool.promoteExecutables(nil)

	//add 4096 more tx
	for i := 0; i < 4096; i++ {
		bobtx := transaction(uint64(40+i), aliceAddr, big.NewInt(1), testTxFee, g_testData.TestGasLimit, key2)
		//fmt.Println(bobtx.Sender(pool.signer))
		_, err := pool.add(bobtx, false)
		//assert.True(t, ok)
		assert.NoError(t, err)
	}
	pool.promoteExecutables(nil)
	pendN, queueN := pool.stats()
	assert.Equal(t, 4106, pendN)
	assert.Equal(t, 0, queueN)
}

func TestTxPool_demoteUnExecutables(t *testing.T) {
	pool := setupTxPool()
	key1, key2, _ := createKey()
	aliceAddr := cs_crypto.GetNormalAddress(key1.PublicKey)

	//transaction enqueue first
	for i := 0; i < 20; i++ {
		bobtx := transaction(uint64(30+i), aliceAddr, big.NewInt(1), g_testData.TestGasPrice, g_testData.TestGasLimit, key2)
		//fmt.Println(bobtx.Sender(pool.signer))
		pool.add(bobtx, false)
		//assert.True(t, ok)
		//assert.NoError(t, err)
	}
	pool.promoteExecutables(nil)
	for addr, _ := range pool.pending {
		pool.currentState.SetNonce(addr, 40)
	}
	pool.demoteUnexecutables()
	pendN, queueN := pool.stats()
	assert.Equal(t, 10, pendN)
	assert.Equal(t, 0, queueN)

	for addr, _ := range pool.pending {
		pool.currentState.SetBalance(addr, big.NewInt(0).Mul(testTxFee, big.NewInt(int64(g_testData.TestGasLimit))))
	}
	pool.demoteUnexecutables()
	pendN, queueN = pool.stats()
	assert.Equal(t, 10, pendN)
	assert.Equal(t, 0, queueN)

	for addr, _ := range pool.pending {
		pool.currentState.SetBalance(addr, testTxFee)
	}
	pool.demoteUnexecutables()
	pendN, queueN = pool.stats()
	assert.Equal(t, 0, pendN)
	assert.Equal(t, 0, queueN)
}

func TestTxPool_ValidateTx(t *testing.T) {
	pool := setupTxPool()
	key1, key2, _ := createKey()
	aliceAddr := cs_crypto.GetNormalAddress(key1.PublicKey)

	//tx size > 32k
	extraData := make([]byte, 32*1024+1)
	FatTx := model.NewTransaction(30, aliceAddr, big.NewInt(3000), g_testData.TestGasPrice, g_testData.TestGasLimit, extraData)
	err := pool.AddRemote(FatTx)
	assert.Error(t, err)

	//error nonce < current nonce
	bobtx := transaction(29, aliceAddr, big.NewInt(3000), testTxFee, g_testData.TestGasLimit, key2)
	//assert.NotNil(t, bobtx)
	err = pool.AddRemote(bobtx)
	//assert.NoError(t, err)
	assert.Error(t, err)

	bobtx = transaction(30, aliceAddr, big.NewInt(3000), big.NewInt(1), g_testData.TestGasLimit, key2)
	err = pool.AddRemote(bobtx)
	assert.NoError(t, err)

	//error amount > balance
	bobtx = transaction(30, aliceAddr, big.NewInt(3000000), big.NewInt(0).Mul(testTxFee, big.NewInt(10000)), g_testData.TestGasLimit, key2)
	err = pool.AddRemote(bobtx)
	//assert.NoError(t, err)
	assert.Error(t, err)

	//for pending
	//same nonce for the 2nd tx, fee < 1st tx
	bobtx = transaction(30, aliceAddr, big.NewInt(300), testTxFee, g_testData.TestGasLimit, key2)
	err = pool.AddRemote(bobtx)
	bobtx = transaction(30, aliceAddr, big.NewInt(200), big.NewInt(0).Sub(testTxFee, big.NewInt(10)), g_testData.TestGasLimit, key2)
	err = pool.AddRemote(bobtx)
	//assert.NoError(t, err)
	assert.Error(t, err)
	// fee = 1st tx
	bobtx = transaction(30, aliceAddr, big.NewInt(210), testTxFee, g_testData.TestGasLimit, key2)
	err = pool.AddRemote(bobtx)
	//assert.NoError(t, err)
	assert.Error(t, err)
	//fee > 1st tx
	bobtx = transaction(30, aliceAddr, big.NewInt(220), threshold, g_testData.TestGasLimit, key2)
	err = pool.AddRemote(bobtx)
	assert.NoError(t, err)

	//for future queue
	//same nonce for the 2nd tx, fee < 1st tx
	bobtx = transaction(33, aliceAddr, big.NewInt(300), testTxFee, g_testData.TestGasLimit, key2)
	err = pool.AddRemote(bobtx)
	bobtx = transaction(33, aliceAddr, big.NewInt(200), big.NewInt(0).Sub(testTxFee, big.NewInt(10)), g_testData.TestGasLimit, key2)
	err = pool.AddRemote(bobtx)
	//assert.NoError(t, err)
	assert.Error(t, err)
	// fee = 1st tx
	bobtx = transaction(33, aliceAddr, big.NewInt(210), testTxFee, g_testData.TestGasLimit, key2)
	err = pool.AddRemote(bobtx)
	//assert.NoError(t, err)
	assert.Error(t, err)
	//fee > 1st tx
	bobtx = transaction(33, aliceAddr, big.NewInt(220), threshold, g_testData.TestGasLimit, key2)
	err = pool.AddRemote(bobtx)
	assert.NoError(t, err)
	pendN, queueN := pool.stats()
	assert.Equal(t, 1, pendN)
	assert.Equal(t, 1, queueN)
}

func TestTxPool_AddTxPerf(t *testing.T) {
	pool := setupTxPool()
	key1, key2, _ := createKey()
	aliceAddr := cs_crypto.GetNormalAddress(key1.PublicKey)
	bobAddr := cs_crypto.GetNormalAddress(key2.PublicKey)

	TXNUM := 256
	txs := make([]model.AbstractTransaction, TXNUM)
	txs1 := make([]model.AbstractTransaction, TXNUM)
	for i := 0; i < TXNUM; i++ {
		tx := transaction(uint64(30+i), aliceAddr, big.NewInt(1), testTxFee, g_testData.TestGasLimit, key2)
		tx1 := transaction(uint64(30+i), bobAddr, big.NewInt(1), testTxFee, g_testData.TestGasLimit, key1)
		txs[i] = tx
		txs1[i] = tx1
	}

	st := time.Now()
	for i := 0; i < TXNUM; i++ {
		err := pool.AddRemote(txs[i])
		//assert.True(t, ok)
		assert.NoError(t, err)
	}
	fmt.Println("add remote using", time.Now().Sub(st))

	st = time.Now()
	pool.AddRemotes(txs1)
	fmt.Println("add remotes using", time.Now().Sub(st))

}

func TestTxPool_AddRemote(t *testing.T) {
	pool := setupTxPool()
	key1, key2, _ := createKey()
	aliceAddr := cs_crypto.GetNormalAddress(key1.PublicKey)
	//bobAddr := cs_crypto.GetNormalAddress(key2.PublicKey)

	//test nonce < record nonce
	//for i:=0; i<10; i++ {
	//	tx := transaction(uint64(20+i), aliceAddr, big.NewInt(1), testTxFee, key2)
	//	err := pool.AddRemote(tx)
	//	//assert.True(t, ok)
	//	assert.Error(t, err)
	//}

	//make pool full
	for i := 0; i < 4096; i++ {
		tx := transaction(uint64(30+i), aliceAddr, big.NewInt(1), testTxFee, g_testData.TestGasLimit, key2)
		err := pool.AddRemote(tx)
		//assert.True(t, ok)
		assert.NoError(t, err)
	}
	txMore := transaction(uint64(30+4096), aliceAddr, big.NewInt(1), testTxFee, g_testData.TestGasLimit, key2)
	pool.AddRemote(txMore)

	txMore = transaction(uint64(30+4097), aliceAddr, big.NewInt(1), testTxFee, g_testData.TestGasLimit, key2)
	pool.AddRemote(txMore)
	pendN, queueN := pool.stats()

	//qu, _ := pool.Queueing()
	//fmt.Println(qu[bobAddr])

	//pe, _ := pool.Pending()
	//fmt.Println(pe[bobAddr][len(pe)-1:])

	assert.Equal(t, 4096, pendN)
	assert.Equal(t, 2, queueN)

	//add one more transaction
	//txMore := transaction(uint64(30+4096), aliceAddr, big.NewInt(1), testTxFee, key2)
	//err := pool.AddRemote(txMore)
	//assert.NoError(t, err)
	//txMore = transaction(uint64(30+4097), aliceAddr, big.NewInt(1), testTxFee, key2)
	//err = pool.AddRemote(txMore)
}

func TestTxFee_Pop(t *testing.T) {
	key1, key2, _ := createKey()
	aliceAddr := cs_crypto.GetNormalAddress(key1.PublicKey)
	bobAddr := cs_crypto.GetNormalAddress(key2.PublicKey)
	//fmt.Println(aliceAddr.Hex(), bobAddr.Hex())

	txs := make(map[common.Address][]model.AbstractTransaction)
	txlB := make([]model.AbstractTransaction, 0)
	txlA := make([]model.AbstractTransaction, 0)
	for i := 0; i < 6; i++ {
		//bob->alice
		tx := transaction(uint64(30+i), aliceAddr, big.NewInt(1), testTxFee, g_testData.TestGasLimit, key2)
		txlB = append(txlB, tx)
	}

	for i := 0; i < 6; i++ {
		//alice->bob
		tx := transaction(uint64(30+i), bobAddr, big.NewInt(1), testTxFee, g_testData.TestGasLimit, key1)
		txlA = append(txlA, tx)
	}
	txs[aliceAddr] = txlA
	txs[bobAddr] = txlB
	//fmt.Println(txs)
	txsFee := model.NewTransactionsByFeeAndNonce(nil, txs)
	fmt.Println(txsFee.Peek())

	fmt.Println("------------")
	txsFee.Pop()
	fmt.Println(txsFee.Peek())
}

func TestAddressesByHeartbeat_Len(t *testing.T) {
	log.InitLogger(log.LvlError)
	pool := setupTxPool()
	key1, key2, _ := createKey()
	aliceAddr := cs_crypto.GetNormalAddress(key1.PublicKey)
	txs := []model.AbstractTransaction{}
	for i := 0; i < 5000; i++ {
		tx := transaction(uint64(30+i), aliceAddr, big.NewInt(1), testTxFee, g_testData.TestGasLimit, key2)
		txs = append(txs, tx)
	}

	start := time.Now()
	pool.AddRemotes(txs)
	fmt.Println(time.Now().Sub(start))
	fmt.Println(pool.Stats())
}

func TestTxBatch(t *testing.T) {
	log.InitLogger(log.LvlError)
	accountNum := 400
	pool := setupTxPoolBatch(accountNum)
	addrs, keys := createTestAddrs(accountNum)
	txs := []model.AbstractTransaction{}
	for i := 0; i < accountNum; i++ {
		for j := 0; j < 13; j++ {
			tx := transaction(uint64(20+j), addrs[i], big.NewInt(1), testTxFee, g_testData.TestGasLimit, keys[i])
			txs = append(txs, tx)
		}
	}
	start := time.Now()
	pool.AddRemotes(txs)
	fmt.Println(time.Now().Sub(start))
	fmt.Println(pool.Stats())
}

func TestTxCalId(t *testing.T) {
	log.InitLogger(log.LvlError)

	key1, key2, _ := createKey()
	aliceAddr := cs_crypto.GetNormalAddress(key1.PublicKey)
	txs := []model.AbstractTransaction{}
	for i := 0; i < 10000; i++ {
		tx := transaction(uint64(30+i), aliceAddr, big.NewInt(1), testTxFee, g_testData.TestGasLimit, key2)
		txs = append(txs, tx)
	}

	coreNum := runtime.NumCPU()
	count := 10000 / coreNum
	fmt.Println("core num = ", coreNum)

	st := time.Now()

	var wg sync.WaitGroup

	for j := 0; j < coreNum; j++ {

		wg.Add(1)

		go func(no int, txp []model.AbstractTransaction) {
			//defer wg.Add(-1)
			defer wg.Done()
			//fmt.Println(j)
			for _, tx := range txp {
				tx.CalTxId()
			}
			fmt.Printf("---%d %v\n", no, time.Now().Sub(st))
		}(j, txs[j*count:(j+1)*count])
	}

	wg.Wait()
	fmt.Println(time.Now().Sub(st))
}

type txHelper struct {
	cacher *model.TxCacher
	wg     sync.WaitGroup
}

func (th *txHelper) help_TxRecover(txs []model.AbstractTransaction) {
	th.cacher.TxRecover(txs)
	th.wg.Done()
}

func dummy() {
	i := 0
	j := 1
	for {
		if j == 0 {
			j++
		}
		i = 65 * 34
		i += 21
		i = i / j
		j++
		time.Sleep(time.Microsecond)
	}
}

func TestTxCacher_TxRecover(t *testing.T) {
	log.InitLogger(log.LvlError)

	key1, key2, _ := createKey()
	aliceAddr := cs_crypto.GetNormalAddress(key1.PublicKey)
	txs := []model.AbstractTransaction{}
	for i := 0; i < 10000; i++ {
		tx := transaction(uint64(30+i), aliceAddr, big.NewInt(1), testTxFee, g_testData.TestGasLimit, key2)
		txs = append(txs, tx)
	}

	txsCmp := []model.AbstractTransaction{}
	bobAddr := cs_crypto.GetNormalAddress(key2.PublicKey)
	for i := 0; i < 10000; i++ {
		tx := transaction(uint64(20+i), bobAddr, big.NewInt(1), testTxFee, g_testData.TestGasLimit, key1)
		txsCmp = append(txsCmp, tx)
	}

	cacher := model.NewTxCacher(runtime.NumCPU())
	helper := &txHelper{
		cacher: cacher,
	}
	//go dummy()
	//go dummy()
	//go dummy()
	//go dummy()
	//go dummy()
	//go dummy()
	st := time.Now()
	helper.wg.Add(1)
	go helper.help_TxRecover(txs[:2500])
	helper.wg.Add(1)
	go helper.help_TxRecover(txs[2500:5000])
	helper.wg.Add(1)
	go helper.help_TxRecover(txs[5000:7500])
	helper.wg.Add(1)
	go helper.help_TxRecover(txs[7500:])
	cacher.TxRecover(txsCmp)
	helper.wg.Wait()
	fmt.Printf("---%v\n", time.Now().Sub(st))

}

func TestTxPool_TxsCaching(t *testing.T) {
	log.InitLogger(log.LvlError)

	key1, key2, _ := createKey()
	aliceAddr := cs_crypto.GetNormalAddress(key1.PublicKey)
	txs := []model.AbstractTransaction{}
	for i := 0; i < 10000; i++ {
		tx := transaction(uint64(30+i), aliceAddr, big.NewInt(1), testTxFee, g_testData.TestGasLimit, key2)
		txs = append(txs, tx)
	}
	pool := setupTxPool()

	st := time.Now()

	pool.AbstractTxsCaching(txs)
	fmt.Printf("---%v\n", time.Now().Sub(st))

	st = time.Now()
	for i := 0; i < 10000; i++ {
		txs[i].Sender(nil)
	}
	fmt.Printf("---%v\n", time.Now().Sub(st))

	st = time.Now()
	for i := 0; i < 10000; i++ {
		txs[i].Sender(nil)
	}
	fmt.Printf("---%v\n", time.Now().Sub(st))
}

func TestTxPool_TxCacheAfterCopy(t *testing.T) {
	log.InitLogger(log.LvlError)

	pool := setupTxPool()
	key1, key2, _ := createKey()
	aliceAddr := cs_crypto.GetNormalAddress(key1.PublicKey)
	txs := []model.AbstractTransaction{}
	for i := 0; i < 4096; i++ {
		tx := transaction(uint64(30+i), aliceAddr, big.NewInt(1), testTxFee, g_testData.TestGasLimit, key2)
		txs = append(txs, tx)
	}
	pool.AddRemotes(txs)
	fmt.Println(pool.Stats())
	pending, _ := pool.Pending()
	st := time.Now()
	txs1 := model.NewTransactionsByFeeAndNonce(nil, pending)
	cnt := 0
	for {
		// Retrieve the next transaction and abort if all done
		tx := txs1.Peek()
		if tx == nil {
			break
		}
		tx.Sender(nil)
		txs1.Shift()
		cnt++
	}
	fmt.Printf("---%v\n", time.Now().Sub(st))
	fmt.Println("comparing...")

	txsCmp := []model.AbstractTransaction{}
	bobAddr := cs_crypto.GetNormalAddress(key2.PublicKey)
	for i := 0; i < 4096; i++ {
		tx := transaction(uint64(20+i), bobAddr, big.NewInt(1), testTxFee, g_testData.TestGasLimit, key1)
		txsCmp = append(txsCmp, tx)
	}
	st = time.Now()
	for i := 0; i < 4096; i++ {
		txsCmp[i].Sender(nil)
	}
	fmt.Printf("---%v\n", time.Now().Sub(st))
}

func TestTxPool_TxsInBlockCache(t *testing.T) {
	log.InitLogger(log.LvlError)

	key1, key2, _ := createKey()
	aliceAddr := cs_crypto.GetNormalAddress(key1.PublicKey)
	txs := []model.AbstractTransaction{}
	for i := 0; i < 4096; i++ {
		tx := transaction(uint64(30+i), aliceAddr, big.NewInt(1), testTxFee, g_testData.TestGasLimit, key2)
		txs = append(txs, tx)
	}

	st := time.Now()
	//for i:=0; i<4096; i++ {
	//	txs[i].Sender(nil)
	//}
	//fmt.Printf("+++%v\n", time.Now().Sub(st))

	st = time.Now()
	header := model.NewHeader(0, 0, common.Hash{}, common.Hash{}, common.Difficulty{}, big.NewInt(0), common.Address{}, common.BlockNonce{})
	trans := make([]*model.Transaction, 4096)
	util.InterfaceSliceCopy(trans, txs)

	block := model.NewBlock(header, trans, []model.AbstractVerification{})
	fmt.Printf("---%v\n", time.Now().Sub(st))

	//pool := setupTxPool()

	fmt.Println(len(block.GetAbsTransactions()))

	st = time.Now()
	//pool.AbstractTxsCaching(block.GetAbsTransactions())
	//
	//fmt.Printf("+++%v\n", time.Now().Sub(st))

	txCmp := block.GetAbsTransactions()
	st = time.Now()
	for i := 0; i < 4096; i++ {
		txCmp[i].Sender(nil)
	}
	fmt.Printf("---%v\n", time.Now().Sub(st))
}

func rlpHash(x interface{}) (h common.Hash, err error) {
	hw := sha3.NewLegacyKeccak256()
	err = rlp.Encode(hw, x)
	if err != nil {
		return
	}
	hw.Sum(h[:0])
	return
}

func rlpHashNew(hw hash.Hash, data []byte) (h common.Hash, err error) {
	hw.Write(data)
	hw.Sum(h[:0])
	return
}

func TestTxPool_PoolSetup(t *testing.T) {
	//pool := setupTxPool()
	//assert.NotNil(t, pool)
	//
	//pending, queue := pool.Stats()
	//fmt.Println(pending, queue)

	log.InitLogger(log.LvlError)

	key1, key2, _ := createKey()
	aliceAddr := cs_crypto.GetNormalAddress(key1.PublicKey)
	txs := []model.AbstractTransaction{}
	for i := 0; i < 10000; i++ {
		tx := transaction(uint64(30+i), aliceAddr, big.NewInt(1), testTxFee, g_testData.TestGasLimit, key2)
		txs = append(txs, tx)
	}
	tx_rlps := [][]byte{}

	st := time.Now()
	for i := 0; i < 10000; i++ {
		ret, _ := txs[i].EncodeRlpToBytes()
		tx_rlps = append(tx_rlps, ret)
	}
	fmt.Printf("---%v\n", time.Now().Sub(st))

	st = time.Now()
	for i := 0; i < 10000; i++ {
		rlpHash(txs[i])
	}
	fmt.Printf("---%v\n", time.Now().Sub(st))

	nonce := make([]byte, 64)
	for i := 0; i < 64; i++ {
		nonce[i] = byte(i)
	}
	hw := sha3.NewLegacyKeccak256()
	st = time.Now()
	for i := 0; i < 10000; i++ {
		splice := append(tx_rlps[i], nonce...)
		rlpHashNew(hw, splice)
	}
	fmt.Printf("---%v\n", time.Now().Sub(st))
}

func TestTxPool(t *testing.T) {
	sDB, _ := state_processor.NewAccountStateDB(common.Hash{}, state_processor.NewStateStorageWithCache(ethdb.NewMemDatabase()))
	//blockchain := &testBlockChain{statedb: teststatedb}
	c := gomock.NewController(t)
	defer c.Finish()
	b1 := model.NewBlock(&model.Header{Number: 1, Bloom: iblt.NewBloom(model.DefaultBlockBloomConfig)}, nil, nil)
	b2 := model.NewBlock(&model.Header{Number: 2, Bloom: iblt.NewBloom(model.DefaultBlockBloomConfig)}, nil, nil)
	b3 := model.NewBlock(&model.Header{Number: 3, Bloom: iblt.NewBloom(model.DefaultBlockBloomConfig)}, nil, nil)
	b4 := model.NewBlock(&model.Header{Number: 4, Bloom: iblt.NewBloom(model.DefaultBlockBloomConfig)}, nil, nil)

	blockchain := NewMockBlockChain(c)
	blockchain.EXPECT().CurrentBlock().Return(b2).AnyTimes()
	blockchain.EXPECT().StateAtByStateRoot(gomock.Any()).Return(sDB, nil).AnyTimes()
	blockchain.EXPECT().GetBlockByNumber(uint64(4)).Return(b4).AnyTimes()
	blockchain.EXPECT().GetBlockByNumber(uint64(3)).Return(b3).AnyTimes()
	blockchain.EXPECT().GetBlockByNumber(uint64(2)).Return(nil).AnyTimes()
	blockchain.EXPECT().GetBlockByNumber(uint64(1)).Return(b1).AnyTimes()

	pool := NewTxPool(testTxPoolConfig, chain_config.ChainConfig{ChainId: big.NewInt(1)}, blockchain)
	pool.Reset(&model.Header{Number: 3, Bloom: iblt.NewBloom(model.DefaultBlockBloomConfig)}, &model.Header{Number: 1, Bloom: iblt.NewBloom(model.DefaultBlockBloomConfig)})
	pool.Reset(&model.Header{Number: 1, Bloom: iblt.NewBloom(model.DefaultBlockBloomConfig)}, &model.Header{Number: 3, Bloom: iblt.NewBloom(model.DefaultBlockBloomConfig)})

	txs := createTxListWithFee(2)
	sender0, err := txs[0].Sender(nil)
	assert.NoError(t, err)
	err = sDB.NewAccountState(sender0)
	assert.NoError(t, sDB.AddBalance(sender0, big.NewInt(consts.DIP)))
	assert.NoError(t, err)
	errs := pool.AddLocals([]model.AbstractTransaction{txs[0]})
	assert.NoError(t, errs[0])

	pool.TxsCaching(txs)
	assert.NotEmpty(t, pool.ConvertPoolToMap())
	assert.NotNil(t, pool.GetTxsEstimator(iblt.NewBloom(model.DefaultBlockBloomConfig)))
	pool.locals.accounts = map[common.Address]struct{}{sender0: {}}
	assert.NotEmpty(t, pool.local())

	pool.config.Rejournal = time.Millisecond
	evictionInterval = time.Millisecond
	go pool.loop()
	time.Sleep(2 * time.Millisecond)
}
