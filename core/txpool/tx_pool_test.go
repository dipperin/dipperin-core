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

package txpool

import (
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/consts"
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/core/bloom"
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/dipperin/dipperin-core/core/chain/state-processor"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third-party/crypto/cs-crypto"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/sha3"
	"math/big"
	"runtime"
	"sync"
	"testing"
	"time"
)

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
	tx1 := transaction(1, bobAddr, big.NewInt(5000), testTxFee, model.TestGasLimit, key1)
	pool.enqueueTx(tx1.CalTxId(), tx1)
	expect1 := pool.Get(tx1.CalTxId())
	assert.Equal(t, expect1.CalTxId(), tx1.CalTxId())

	//enqueueTx try add tx2 in pool queue list but the fee does not exceed the feeBump ,so enqueueTx failed
	tx2 := transaction(1, bobAddr, big.NewInt(4000), big.NewInt(0).Add(testTxFee, big.NewInt(9)), model.TestGasLimit, key1)
	expectOK, err := pool.enqueueTx(tx1.CalTxId(), tx2)
	assert.Equal(t, expectOK, false)
	assert.Error(t, err)
	expect2 := pool.Get(tx1.CalTxId())
	assert.Equal(t, expect2.CalTxId(), tx1.CalTxId())
	expect2Get := pool.Get(tx2.CalTxId())
	assert.Equal(t, expect2Get, nil)

	//enqueueTx try add tx3 in pool queue list and the fee exceeds the feeBump, so the enqueueTx success.
	tx3 := transaction(1, bobAddr, big.NewInt(4000), big.NewInt(0).Add(testTxFee, testTxFee), model.TestGasLimit, key1)
	expectOK, err = pool.enqueueTx(tx1.CalTxId(), tx3)
	assert.Equal(t, expectOK, true)
	assert.NoError(t, err)
	expect3 := pool.Get(tx1.CalTxId())
	assert.Equal(t, expect3, nil)
	expect3Get := pool.Get(tx3.CalTxId())
	assert.Equal(t, expect3Get.CalTxId(), tx3.CalTxId())

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
	unsignedTx1 := model.NewTransaction(1, bobAddr, big.NewInt(5000), model.TestGasPrice, model.TestGasLimit, overLoad)
	signedTx1, err := unsignedTx1.SignTx(key1, ms)
	assert.NoError(t, err)
	err = pool.validateTx(signedTx1, false)
	assert.Error(t, err)

	signedTx2 := transaction(1, bobAddr, big.NewInt(5000), testTxFee, model.TestGasLimit, key1)
	err = pool.validateTx(signedTx2, false)
	expectErr := "tx nonce is invalid"
	assert.EqualError(t, err, expectErr)

	unsignedTx3 := model.NewTransaction(1, bobAddr, big.NewInt(5000), model.TestGasPrice, model.TestGasLimit, nil)
	signedTx3, err := unsignedTx3.SignTx(key1, model.NewSigner(big.NewInt(2)))
	assert.NoError(t, err)
	err = pool.validateTx(signedTx3, false)
	expectErr = "invalid sender"
	assert.EqualError(t, err, expectErr)

	signedTx4 := transaction(1, bobAddr, big.NewInt(5000), big.NewInt(1), model.TestGasLimit, key1)
	err = pool.validateTx(signedTx4, true)
	expectErr = "tx nonce is invalid"
	assert.EqualError(t, err, expectErr)

	signedTx5 := transaction(40, bobAddr, big.NewInt(1000000), big.NewInt(0).Mul(testTxFee, big.NewInt(10000)), model.TestGasLimit, key1)
	curBalance, e := pool.currentState.GetBalance(aliceAddr)
	err = pool.validateTx(signedTx5, false)
	expectErr = fmt.Sprintf("tx exceed balance limit, from:%v, cur balance:%v, cost:%v, err:%v", aliceAddr.Hex(), curBalance.String(), signedTx5.Cost().String(), e)
	assert.EqualError(t, err, expectErr)
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

	alicetx1 := transaction(1, bobAddr, big.NewInt(5000), testTxFee, model.TestGasLimit, key1)
	alicetx2 := transaction(2, bobAddr, big.NewInt(4000), big.NewInt(0).Add(testTxFee, big.NewInt(9)), model.TestGasLimit, key1)
	alicetx3 := transaction(3, bobAddr, big.NewInt(6000), big.NewInt(0).Add(testTxFee, big.NewInt(10)), model.TestGasLimit, key1)

	bobtx1 := transaction(4, aliceAddr, big.NewInt(5000), testTxFee, model.TestGasLimit, key2)
	bobtx2 := transaction(5, aliceAddr, big.NewInt(4000), big.NewInt(0).Add(testTxFee, big.NewInt(9)), model.TestGasLimit, key2)
	bobtx3 := transaction(6, aliceAddr, big.NewInt(6000), big.NewInt(0).Add(testTxFee, big.NewInt(10)), model.TestGasLimit, key2)

	txListA := []model.AbstractTransaction{alicetx1, alicetx2, alicetx3, bobtx1, bobtx2}
	txListB := []model.AbstractTransaction{alicetx2, alicetx3, bobtx1, bobtx2, bobtx3}

	difference := model.TxDifference(txListA, txListB)
	assert.Equal(t, 1, len(difference))
	assert.Equal(t, alicetx1, difference[0])

	difference = model.TxDifference(txListB, txListA)
	assert.Equal(t, 1, len(difference))
	assert.Equal(t, bobtx3, difference[0])
}

func TestPromoteTx(t *testing.T) {
	//t.Skip()
	pool := setupTxPool()
	key1, key2, _ := createKey()
	aliceAddr := cs_crypto.GetNormalAddress(key1.PublicKey)
	bobAddr := cs_crypto.GetNormalAddress(key2.PublicKey)

	alicetx1 := transaction(0, bobAddr, big.NewInt(5000), testTxFee, model.TestGasLimit, key1)
	alicetx2 := transaction(0, bobAddr, big.NewInt(4000), big.NewInt(0).Add(testTxFee, big.NewInt(9)), model.TestGasLimit, key1)
	alicetx3 := transaction(0, bobAddr, big.NewInt(6000), big.NewInt(0).Add(testTxFee, testTxFee), model.TestGasLimit, key1)

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

	key1, key2, _ := createKey()
	aliceAddr := cs_crypto.GetNormalAddress(key1.PublicKey)
	bobAddr := cs_crypto.GetNormalAddress(key2.PublicKey)

	alicetx1 := transaction(19, bobAddr, big.NewInt(100000), testTxFee, model.TestGasLimit, key1)
	alicetx2 := transaction(20, bobAddr, big.NewInt(4000), testTxFee, model.TestGasLimit, key1)
	bobtx1 := transaction(30, aliceAddr, big.NewInt(3000), testTxFee, model.TestGasLimit, key2)
	bobtx2 := transaction(31, aliceAddr, big.NewInt(3000), testTxFee, model.TestGasLimit, key2)
	bobtx3 := transaction(33, aliceAddr, big.NewInt(3000), testTxFee, model.TestGasLimit, key2)

	pool.AddRemote(alicetx1)
	pool.AddRemote(alicetx2)
	pool.AddRemote(bobtx1)
	pool.AddRemote(bobtx3)
	pool.AddRemote(bobtx2)

	pendTx, _ := pool.Pending()
	except := 3
	assert.Equal(t, except, len(pendTx[aliceAddr])+len(pendTx[bobAddr]))

	queueTx, _ := pool.Queueing()
	except = 1
	assert.Equal(t, except, len(queueTx[aliceAddr])+len(queueTx[bobAddr]))
}

func TestTxPool_removeTx(t *testing.T) {
	pool := setupTxPool()
	key1, key2, _ := createKey()
	aliceAddr := cs_crypto.GetNormalAddress(key1.PublicKey)
	bobAddr := cs_crypto.GetNormalAddress(key2.PublicKey)

	bobtx1 := transaction(30, aliceAddr, big.NewInt(3000), testTxFee, model.TestGasLimit, key2)
	bobtx2 := transaction(31, aliceAddr, big.NewInt(3000), testTxFee, model.TestGasLimit, key2)
	bobtx3 := transaction(32, aliceAddr, big.NewInt(3000), testTxFee, model.TestGasLimit, key2)

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
	pending, _ := pool.Pending()
	queueing, _ := pool.Queueing()

	assert.Equal(t, 0, len(pending[bobAddr]))
	assert.Equal(t, bobtx2.CalTxId(), queueing[bobAddr][0].CalTxId())
	assert.Equal(t, bobtx3.CalTxId(), queueing[bobAddr][1].CalTxId())
}

func TestTxPool_RemoveTx(t *testing.T) {
	pool := setupTxPool()
	key1, key2, _ := createKey()
	aliceAddr := cs_crypto.GetNormalAddress(key1.PublicKey)

	bobtx1 := transaction(30, aliceAddr, big.NewInt(3000), testTxFee, model.TestGasLimit, key2)
	bobtx2 := transaction(31, aliceAddr, big.NewInt(3000), testTxFee, model.TestGasLimit, key2)
	bobtx3 := transaction(32, aliceAddr, big.NewInt(3000), testTxFee, model.TestGasLimit, key2)

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
		tx := transaction(uint64(30+i), aliceAddr, big.NewInt(1), testTxFee, model.TestGasLimit, key2)
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

	key1, key2, _ := createKey()
	aliceAddr := cs_crypto.GetNormalAddress(key1.PublicKey)
	bobAddr := cs_crypto.GetNormalAddress(key2.PublicKey)

	alicetx1 := transaction(19, bobAddr, big.NewInt(100000), testTxFee, model.TestGasLimit, key1)
	alicetx2 := transaction(20, bobAddr, big.NewInt(4000), testTxFee, model.TestGasLimit, key1)
	bobtx1 := transaction(30, aliceAddr, big.NewInt(3000), testTxFee, model.TestGasLimit, key2)
	bobtx2 := transaction(31, aliceAddr, big.NewInt(3000), testTxFee, model.TestGasLimit, key2)
	bobtx3 := transaction(33, aliceAddr, big.NewInt(3000), testTxFee, model.TestGasLimit, key2)

	pool.AddLocal(alicetx1)
	pool.AddLocal(alicetx2)
	pool.AddLocal(bobtx1)
	pool.AddLocal(bobtx3)
	pool.AddLocal(bobtx2)

	pendingTX, _ := pool.Pending()
	queueTX, _ := pool.Queueing()
	assert.Equal(t, 3, len(pendingTX[aliceAddr])+len(pendingTX[bobAddr]))
	assert.Equal(t, 1, len(queueTX[aliceAddr])+len(queueTX[bobAddr]))
}

func TestStatus(t *testing.T) {
	pool := setupTxPool()
	key1, key2, _ := createKey()
	aliceAddr := cs_crypto.GetNormalAddress(key1.PublicKey)
	bobtx1 := transaction(30, aliceAddr, big.NewInt(3000), testTxFee, model.TestGasLimit, key2)
	bobtx2 := transaction(31, aliceAddr, big.NewInt(3000), testTxFee, model.TestGasLimit, key2)
	bobtx3 := transaction(33, aliceAddr, big.NewInt(3000), testTxFee, model.TestGasLimit, key2)
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

	tx := transaction(uint64(30), aliceAddr, big.NewInt(1), testTxFee, model.TestGasLimit, key2)
	err := pool.AddLocal(tx)
	assert.NoError(t, err)
}

func TestTxPool_LocalAdd(t *testing.T) {
	pool := setupTxPool()
	key1, key2, _ := createKey()
	aliceAddr := cs_crypto.GetNormalAddress(key1.PublicKey)

	//test nonce < record nonce
	for i := 0; i < 10; i++ {
		tx := transaction(uint64(20+i), aliceAddr, big.NewInt(1), testTxFee, model.TestGasLimit, key2)
		_, err := pool.add(tx, false)
		//assert.True(t, ok)
		assert.Error(t, err)
	}

	//make pool full
	for i := 0; i < 5120; i++ {
		tx := transaction(uint64(30+i), aliceAddr, big.NewInt(1), testTxFee, model.TestGasLimit, key2)
		_, err := pool.add(tx, false)
		//assert.True(t, ok)
		assert.NoError(t, err)
	}
	pendN, queueN := pool.stats()
	assert.Equal(t, 0, pendN)
	assert.Equal(t, 5120, queueN)

	//add one more transaction
	txMore := transaction(uint64(30+5120), aliceAddr, big.NewInt(1), testTxFee, model.TestGasLimit, key2)
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
		bobtx := transaction(uint64(30+i), aliceAddr, big.NewInt(1), testTxFee, model.TestGasLimit, key2)
		_, err := pool.add(bobtx, false)
		assert.NoError(t, err)
	}

	for i := 0; i < 33; i++ {
		alicetx := transaction(uint64(20+i), chalieAddr, big.NewInt(1), testTxFee, model.TestGasLimit, key1)
		_, err := pool.add(alicetx, false)
		assert.NoError(t, err)
	}

	for i := 0; i < 33; i++ {
		chalietx := transaction(uint64(30+i), bobAddr, big.NewInt(1), testTxFee, model.TestGasLimit, key3)
		_, err := pool.add(chalietx, false)
		assert.NoError(t, err)
	}

	pendN, queueN := pool.stats()
	assert.Equal(t, 0, pendN)
	assert.Equal(t, 96, queueN)

	pool.promoteExecutables(nil)
	pendN, queueN = pool.stats()
	assert.Equal(t, 96, pendN)
	assert.Equal(t, 0, queueN)
}

func TestTxPool_promptExecutables(t *testing.T) {
	pool := setupTxPool()
	key1, key2, _ := createKey()
	aliceAddr := cs_crypto.GetNormalAddress(key1.PublicKey)

	//transaction enqueue first
	for i := 0; i < 10; i++ {
		bobtx := transaction(uint64(30+i), aliceAddr, big.NewInt(1), testTxFee, model.TestGasLimit, key2)
		_, err := pool.add(bobtx, false)
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
		bobtx := transaction(uint64(40+i), aliceAddr, big.NewInt(1), testTxFee, model.TestGasLimit, key2)
		_, err := pool.add(bobtx, false)
		assert.NoError(t, err)
	}
	pool.promoteExecutables(nil)
	pendN, queueN = pool.stats()
	assert.Equal(t, 4096, pendN)
	assert.Equal(t, 0, queueN)

	//add 1024 more tx
	for i := 0; i < 1024; i++ {
		bobtx := transaction(uint64(4136+i), aliceAddr, big.NewInt(1), testTxFee, model.TestGasLimit, key2)
		_, err := pool.add(bobtx, false)
		assert.NoError(t, err)
	}

	pendN, queueN = pool.stats()
	assert.Equal(t, 4096, pendN)
	assert.Equal(t, 1024, queueN)

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
		bobtx := transaction(uint64(30+i), aliceAddr, big.NewInt(1), testTxFee, model.TestGasLimit, key2)
		//add bob address to local too
		_, err := pool.add(bobtx, true)
		assert.NoError(t, err)
	}
	pool.promoteExecutables(nil)

	//add 4096 more tx
	for i := 0; i < 4096; i++ {
		bobtx := transaction(uint64(40+i), aliceAddr, big.NewInt(1), testTxFee, model.TestGasLimit, key2)
		_, err := pool.add(bobtx, false)
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
		bobtx := transaction(uint64(30+i), aliceAddr, big.NewInt(1), model.TestGasPrice, model.TestGasLimit, key2)
		pool.add(bobtx, false)
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
		pool.currentState.SetBalance(addr, big.NewInt(0).Mul(testTxFee, big.NewInt(int64(model.TestGasLimit))))
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
	FatTx := model.NewTransaction(30, aliceAddr, big.NewInt(3000), model.TestGasPrice, model.TestGasLimit, extraData)
	err := pool.AddRemote(FatTx)
	assert.Error(t, err)

	//error nonce < current nonce
	bobtx := transaction(29, aliceAddr, big.NewInt(3000), testTxFee, model.TestGasLimit, key2)
	err = pool.AddRemote(bobtx)
	assert.Error(t, err)

	bobtx = transaction(30, aliceAddr, big.NewInt(3000), big.NewInt(1), model.TestGasLimit, key2)
	err = pool.AddRemote(bobtx)
	assert.NoError(t, err)

	//error amount > balance
	bobtx = transaction(30, aliceAddr, big.NewInt(3000000), big.NewInt(0).Mul(testTxFee, big.NewInt(10000)), model.TestGasLimit, key2)
	err = pool.AddRemote(bobtx)
	assert.Error(t, err)

	//for pending
	//same nonce for the 2nd tx, fee < 1st tx
	bobtx = transaction(30, aliceAddr, big.NewInt(300), testTxFee, model.TestGasLimit, key2)
	err = pool.AddRemote(bobtx)
	bobtx = transaction(30, aliceAddr, big.NewInt(200), big.NewInt(0).Sub(testTxFee, big.NewInt(10)), model.TestGasLimit, key2)
	err = pool.AddRemote(bobtx)
	assert.Error(t, err)
	// fee = 1st tx
	bobtx = transaction(30, aliceAddr, big.NewInt(210), testTxFee, model.TestGasLimit, key2)
	err = pool.AddRemote(bobtx)
	assert.Error(t, err)
	//fee > 1st tx
	bobtx = transaction(30, aliceAddr, big.NewInt(220), threshold, model.TestGasLimit, key2)
	err = pool.AddRemote(bobtx)
	assert.NoError(t, err)

	//for future queue
	//same nonce for the 2nd tx, fee < 1st tx
	bobtx = transaction(33, aliceAddr, big.NewInt(300), testTxFee, model.TestGasLimit, key2)
	err = pool.AddRemote(bobtx)
	bobtx = transaction(33, aliceAddr, big.NewInt(200), big.NewInt(0).Sub(testTxFee, big.NewInt(10)), model.TestGasLimit, key2)
	err = pool.AddRemote(bobtx)
	assert.Error(t, err)
	// fee = 1st tx
	bobtx = transaction(33, aliceAddr, big.NewInt(210), testTxFee, model.TestGasLimit, key2)
	err = pool.AddRemote(bobtx)
	assert.Error(t, err)
	//fee > 1st tx
	bobtx = transaction(33, aliceAddr, big.NewInt(220), threshold, model.TestGasLimit, key2)
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
		tx := transaction(uint64(30+i), aliceAddr, big.NewInt(1), testTxFee, model.TestGasLimit, key2)
		tx1 := transaction(uint64(30+i), bobAddr, big.NewInt(1), testTxFee, model.TestGasLimit, key1)
		txs[i] = tx
		txs1[i] = tx1
	}

	st := time.Now()
	for i := 0; i < TXNUM; i++ {
		err := pool.AddRemote(txs[i])
		assert.NoError(t, err)
	}
	assert.Equal(t, true, time.Now().Sub(st).Milliseconds() < 1000)

	st = time.Now()
	pool.AddRemotes(txs1)
	assert.Equal(t, true, time.Now().Sub(st).Milliseconds() < 1000)
}

func TestTxPool_AddRemote(t *testing.T) {
	pool := setupTxPool()
	key1, key2, _ := createKey()
	aliceAddr := cs_crypto.GetNormalAddress(key1.PublicKey)

	//make pool full
	for i := 0; i < 4096; i++ {
		tx := transaction(uint64(30+i), aliceAddr, big.NewInt(1), testTxFee, model.TestGasLimit, key2)
		err := pool.AddRemote(tx)
		assert.NoError(t, err)
	}
	txMore := transaction(uint64(30+4096), aliceAddr, big.NewInt(1), testTxFee, model.TestGasLimit, key2)
	pool.AddRemote(txMore)

	txMore = transaction(uint64(30+4097), aliceAddr, big.NewInt(1), testTxFee, model.TestGasLimit, key2)
	pool.AddRemote(txMore)
	pendN, queueN := pool.stats()

	assert.Equal(t, 4096, pendN)
	assert.Equal(t, 2, queueN)
}

func TestTxFee_Pop(t *testing.T) {
	key1, key2, _ := createKey()
	aliceAddr := cs_crypto.GetNormalAddress(key1.PublicKey)
	bobAddr := cs_crypto.GetNormalAddress(key2.PublicKey)

	txs := make(map[common.Address][]model.AbstractTransaction)
	txlB := make([]model.AbstractTransaction, 0)
	txlA := make([]model.AbstractTransaction, 0)
	for i := 0; i < 6; i++ {
		//bob->alice
		tx := transaction(uint64(30+i), aliceAddr, big.NewInt(1), testTxFee, model.TestGasLimit, key2)
		txlB = append(txlB, tx)
	}

	for i := 0; i < 6; i++ {
		//alice->bob
		tx := transaction(uint64(30+i), bobAddr, big.NewInt(1), testTxFee, model.TestGasLimit, key1)
		txlA = append(txlA, tx)
	}
	txs[aliceAddr] = txlA
	txs[bobAddr] = txlB

	txsFee := model.NewTransactionsByFeeAndNonce(nil, txs)
	expect1 := txsFee.Peek()

	txsFee.Pop()
	expect2 := txsFee.Peek()

	assert.NotEqualf(t, expect1.CalTxId(), expect2.CalTxId(), "not equal")
}

func TestAddressesByHeartbeat_Len(t *testing.T) {
	pool := setupTxPool()
	key1, key2, _ := createKey()
	aliceAddr := cs_crypto.GetNormalAddress(key1.PublicKey)
	txs := []model.AbstractTransaction{}
	for i := 0; i < 5000; i++ {
		tx := transaction(uint64(30+i), aliceAddr, big.NewInt(1), testTxFee, model.TestGasLimit, key2)
		txs = append(txs, tx)
	}

	start := time.Now()
	pool.AddRemotes(txs)
	assert.Equal(t, true, time.Now().Sub(start).Seconds() < 10)

	pending, queueing := pool.Stats()
	assert.Equal(t, 4096, pending)
	assert.Equal(t, 0, queueing)
}

func TestTxBatch(t *testing.T) {
	accountNum := 400
	pool := setupTxPoolBatch(accountNum)
	addrs, keys := createTestAddrs(accountNum)
	txs := []model.AbstractTransaction{}
	for i := 0; i < accountNum; i++ {
		for j := 0; j < 13; j++ {
			tx := transaction(uint64(20+j), addrs[i], big.NewInt(1), testTxFee, model.TestGasLimit, keys[i])
			txs = append(txs, tx)
		}
	}

	start := time.Now()
	pool.AddRemotes(txs)
	assert.Equal(t, true, time.Now().Sub(start).Seconds() < 10)

	pending, queueing := pool.Stats()
	assert.Equal(t, 0, pending)
	assert.Equal(t, 0, queueing)
}

func TestTxCalId(t *testing.T) {
	key1, key2, _ := createKey()
	aliceAddr := cs_crypto.GetNormalAddress(key1.PublicKey)
	txs := []model.AbstractTransaction{}
	for i := 0; i < 10000; i++ {
		tx := transaction(uint64(30+i), aliceAddr, big.NewInt(1), testTxFee, model.TestGasLimit, key2)
		txs = append(txs, tx)
	}

	coreNum := runtime.NumCPU()
	count := 10000 / coreNum

	st := time.Now()
	var wg sync.WaitGroup

	for j := 0; j < coreNum; j++ {
		wg.Add(1)
		go func(no int, txp []model.AbstractTransaction) {
			//defer wg.Add(-1)
			defer wg.Done()
			for _, tx := range txp {
				tx.CalTxId()
			}
		}(j, txs[j*count:(j+1)*count])
	}

	wg.Wait()
	assert.Equal(t, true, time.Now().Sub(st).Seconds() < 10)
}

func TestTxCacher_TxRecover(t *testing.T) {
	key1, key2, _ := createKey()
	aliceAddr := cs_crypto.GetNormalAddress(key1.PublicKey)
	txs := []model.AbstractTransaction{}
	for i := 0; i < 10000; i++ {
		tx := transaction(uint64(30+i), aliceAddr, big.NewInt(1), testTxFee, model.TestGasLimit, key2)
		txs = append(txs, tx)
	}

	txsCmp := []model.AbstractTransaction{}
	bobAddr := cs_crypto.GetNormalAddress(key2.PublicKey)
	for i := 0; i < 10000; i++ {
		tx := transaction(uint64(20+i), bobAddr, big.NewInt(1), testTxFee, model.TestGasLimit, key1)
		txsCmp = append(txsCmp, tx)
	}

	cacher := model.NewTxCacher(runtime.NumCPU())
	helper := &txHelper{
		cacher: cacher,
	}

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

	assert.Equal(t, len(txs), len(txsCmp))
	assert.Equal(t, true, time.Now().Sub(st).Seconds() < 10)
}

func TestTxPool_TxsCaching(t *testing.T) {
	key1, key2, _ := createKey()
	aliceAddr := cs_crypto.GetNormalAddress(key1.PublicKey)
	txs := []model.AbstractTransaction{}
	for i := 0; i < 10000; i++ {
		tx := transaction(uint64(30+i), aliceAddr, big.NewInt(1), testTxFee, model.TestGasLimit, key2)
		txs = append(txs, tx)
	}
	pool := setupTxPool()

	st := time.Now()

	pool.AbstractTxsCaching(txs)
	assert.Equal(t, 10000, len(txs))

	st = time.Now()
	for i := 0; i < 10000; i++ {
		txs[i].Sender(nil)
	}

	assert.Equal(t, true, time.Now().Sub(st).Seconds() < 10)
}

func TestTxPool_TxCacheAfterCopy(t *testing.T) {
	pool := setupTxPool()
	key1, key2, _ := createKey()
	aliceAddr := cs_crypto.GetNormalAddress(key1.PublicKey)
	txs := []model.AbstractTransaction{}
	for i := 0; i < 4096; i++ {
		tx := transaction(uint64(30+i), aliceAddr, big.NewInt(1), testTxFee, model.TestGasLimit, key2)
		txs = append(txs, tx)
	}

	pool.AddRemotes(txs)
	p, q := pool.Stats()
	assert.Equal(t, 4096, p)
	assert.Equal(t, 0, q)

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

	txsCmp := []model.AbstractTransaction{}
	bobAddr := cs_crypto.GetNormalAddress(key2.PublicKey)
	for i := 0; i < 4096; i++ {
		tx := transaction(uint64(20+i), bobAddr, big.NewInt(1), testTxFee, model.TestGasLimit, key1)
		txsCmp = append(txsCmp, tx)
	}
	st = time.Now()
	for i := 0; i < 4096; i++ {
		txsCmp[i].Sender(nil)
	}

	assert.Equal(t, len(txs), len(txsCmp))
	assert.Equal(t, true, time.Now().Sub(st).Seconds() < 10)
}

func TestTxPool_TxsInBlockCache(t *testing.T) {
	key1, key2, _ := createKey()
	aliceAddr := cs_crypto.GetNormalAddress(key1.PublicKey)
	txs := []model.AbstractTransaction{}
	for i := 0; i < 4096; i++ {
		tx := transaction(uint64(30+i), aliceAddr, big.NewInt(1), testTxFee, model.TestGasLimit, key2)
		txs = append(txs, tx)
	}

	header := model.NewHeader(0, 0, common.Hash{}, common.Hash{}, common.Difficulty{}, big.NewInt(0), common.Address{}, common.BlockNonce{})
	trans := make([]*model.Transaction, 4096)
	util.InterfaceSliceCopy(trans, txs)

	block := model.NewBlock(header, trans, []model.AbstractVerification{})
	assert.Equal(t, 4096, len(block.GetAbsTransactions()))

	st := time.Now()

	txCmp := block.GetAbsTransactions()
	st = time.Now()
	for i := 0; i < 4096; i++ {
		txCmp[i].Sender(nil)
	}
	assert.Equal(t,true,time.Now().Sub(st).Seconds()<10)
	assert.Equal(t,len(txs),len(txCmp))
}

func TestTxPool_PoolSetup(t *testing.T) {
	key1, key2, _ := createKey()
	aliceAddr := cs_crypto.GetNormalAddress(key1.PublicKey)
	txs := []model.AbstractTransaction{}
	for i := 0; i < 10000; i++ {
		tx := transaction(uint64(30+i), aliceAddr, big.NewInt(1), testTxFee, model.TestGasLimit, key2)
		txs = append(txs, tx)
	}
	txRlps := [][]byte{}

	st := time.Now()
	for i := 0; i < 10000; i++ {
		ret, _ := txs[i].EncodeRlpToBytes()
		txRlps = append(txRlps, ret)
	}
	fmt.Printf("---%v\n", time.Now().Sub(st))
	assert.Equal(t,len(txs),len(txRlps))

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
		splice := append(txRlps[i], nonce...)
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
