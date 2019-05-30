package state_processor

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestAccountStateDB_RevertToSnapshot1(t *testing.T){
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

	before,_:= processor.GetBalance(aliceAddr)
	assert.Equal(t, big.NewInt(2000), before)

	processor.RevertToSnapshot(snapshot)
	ba,_ := processor.GetBalance(aliceAddr)
	var nilBigInt *big.Int
	assert.Equal(t,nilBigInt,ba)

	fRoot, err := processor.Finalise()
	assert.NoError(t, err)
	savedRoot, err := processor.Commit()
	assert.Equal(t, fRoot, savedRoot)

}

func TestAccountStateDB_RevertToSnapshot2(t *testing.T){
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

	err = processor.SetCode(aliceAddr,[]byte("coooode"))
	code,err := processor.GetCode(aliceAddr)
	assert.Equal(t,[]byte("coooode"),code)

	err = processor.SetAbi(aliceAddr,[]byte("{input:int}"))
	abi,err := processor.GetAbi(aliceAddr)
	assert.Equal(t,[]byte("{input:int}"),abi)

	root,err =processor.GetDataRoot(aliceAddr)
	assert.NoError(t,err)
	assert.Equal(t, common.Hash{}, root)

	processor.SetData(aliceAddr,"tkey",[]byte("value"))

	assert.Equal(t,[]byte("value"),processor.smartContractData[aliceAddr]["tkey"])

	formerRoot,_ := processor.GetDataRoot(aliceAddr)
	err = processor.finalSmartData()
	root,err =processor.GetDataRoot(aliceAddr)

	assert.NotEqual(t,formerRoot,root)

	processor.RevertToSnapshot(snapshot)
	reverted,_:=processor.GetDataRoot(aliceAddr)

	assert.Equal(t,formerRoot,reverted)
	assert.NotEqual(t,reverted,root)

	abi,err = processor.GetAbi(aliceAddr)
	assert.NotEqual(t,[]byte("{input:int}"),abi)

	code,err = processor.GetCode(aliceAddr)
	assert.NotEqual(t,[]byte("coooode"),code)
}

// To revert the trie
func TestAccountStateDB_RevertToSnapshot3(t *testing.T){

}