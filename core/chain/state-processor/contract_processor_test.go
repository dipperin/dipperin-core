package state_processor

import (
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAccountStateDB_ProcessContract(t *testing.T) {
	/*ownAddress := common.HexToAddress("0x000062be10f46b5d01Ecd9b502c4bA3d6131f6fc2e41")
	log.InitLogger(log.LvlDebug)
	transactionJson := "{\"TxData\":{\"nonce\":\"0x6\",\"to\":\"0x00000000000000000000000000000000000000000018\",\"hashlock\":null,\"timelock\":\"0x0\",\"value\":\"0x2540be400\",\"fee\":\"0x424ad340\",\"gasPrice\":\"0xa\",\"gas\":\"0xa1d8adbf400\",\"input\":\"0x7b22636f6465223a224147467a625145414141414244514e674158384159414a2f6677426741414143485149445a573532426e427961573530637741414132567564676877636d6c7564484e666241414241775144416749414241554263414542415155444151414342685544667746426b49674543333841515a43494241742f41454747434173484e4155476257567462334a354167414c5831396f5a57467758324a68633255444151706658325268644746665a57356b4177494561573570644141444257686c624778764141514b52514d4341417343414173394151462f4977424245477369415351415159414945414167415545674f674150494146424432704241524142494141514143414251516f36414134674155454f616b45424541456741554551616951414377734e415142426741674c426d686c6247787641413d3d222c22616269223a2257776f674943416765776f67494341674943416749434a755957316c496a6f67496d6c75615851694c416f67494341674943416749434a70626e423164484d694f6942625853774b494341674943416749434169623356306348563063794936494674644c416f67494341674943416749434a6a6232357a644746756443493649434a6d5957787a5a534973436941674943416749434167496e5235634755694f6941695a6e5675593352706232346943694167494342394c416f674943416765776f67494341674943416749434a755957316c496a6f67496d686c624778764969774b494341674943416749434169615735776458527a496a6f6757776f6749434167494341674943416749434237436941674943416749434167494341674943416749434169626d46745a53493649434a755957316c4969774b494341674943416749434167494341674943416749434a306558426c496a6f67496e4e30636d6c755a79494b4943416749434167494341674943416766516f67494341674943416749463073436941674943416749434167496d39316448423164484d694f6942625853774b4943416749434167494341695932397563335268626e51694f69416964484a315a534973436941674943416749434167496e5235634755694f6941695a6e567559335270623234694369416749434239436c303d222c22496e707574223a6e756c6c7d\"},\"Wit\":{\"r\":\"0x2f9d296eeeda2bfe075729dc6114b593550814695553ce895a16da21a714b7b5\",\"s\":\"0x76134a5681f704456ab9c875e3c78e3f49ce0fcdc73b2616792502e62c40f9e0\",\"v\":\"0x39\",\"hashkey\":\"0x\"}}"

	var tx model.Transaction
	err := json.Unmarshal([]byte(transactionJson), &tx)
	if err != nil {
		log.Info("TestAccountStateDB_ProcessContract", "err", err)
	}
	log.Info("processContract", "Tx", tx)*/

	/*block := model.NewBlock(model.NewHeader(1, 101, common.HexToHash("ss"), common.HexToHash("fdfs"), common.StringToDiff("000000000000000000000011"), big.NewInt(111), common.StringToAddress("fdsfds"), common.EncodeNonce(33)), nil, nil)*/

	/*block := model.NewBlock(model.NewHeader(1, 10, common.Hash{}, common.HexToHash("1111"), common.HexToDiff("0x20ffffff"), big.NewInt(324234), common.Address{}, common.BlockNonceFromInt(432423)),nil,nil)



	db := ethdb.NewMemDatabase()
	sdb := NewStateStorageWithCache(db)
	processor, _ := NewAccountStateDB(common.Hash{}, sdb)
	processor.NewAccountState(ownAddress)
	err = processor.AddNonce(ownAddress, 6)
	assert.NoError(t, err)
	processor.AddBalance(ownAddress,  new(big.Int).SetInt64(int64(100000000000000000)))
	balance, err := processor.GetBalance(ownAddress)
	log.Info("balance", "balance", balance.String())

	receipt, err := processor.ProcessContract(&tx, block.Header().(*model.Header), true,fakeGetBlockHash)
	assert.NoError(t, err)
	assert.Equal(t, true, receipt.HandlerResult)*/

}

func TestAccountStateDB_ProcessContract2(t *testing.T) {
	var testPath = "/home/qydev/go/src/github.com/dipperin/dipperin-core/core/vm/event"
	tx := createContractTx(t,testPath+"/event.wasm", testPath+"/event.cpp.abi.json")

	db, root := createTestStateDB()
	processor, err := NewAccountStateDB(root, NewStateStorageWithCache(db))
	assert.NoError(t, err)

	nonce,_:=processor.GetNonce(aliceAddr)
	fmt.Println(nonce)

	block := createBlock(1, common.Hash{}, []*model.Transaction{tx})
	//gasPool := gasLimit*5
	conf := TxProcessConfig{
		Tx:tx,
		TxIndex:0,
		Header:block.Header().(*model.Header),
		GetHash:fakeGetBlockHash,
	}
	err = processor.ProcessTxNew(&conf)
	assert.NoError(t, err)

	fullReceipt,err:= tx.GetReceipt()
	nonce, err = processor.GetNonce(fullReceipt.ContractAddress)
	fmt.Println(nonce, err)
	code,err:=processor.GetCode(fullReceipt.ContractAddress)
	trueCode := getContractCode(t,testPath+"/event.wasm", testPath+"/event.cpp.abi.json")
	assert.Equal(t,trueCode,code)
}