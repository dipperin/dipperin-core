package state_processor

import (
	"encoding/json"
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/vmcommon"
	"github.com/dipperin/dipperin-core/core/accounts"
	"github.com/dipperin/dipperin-core/core/accounts/soft-wallet"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestAccountStateDB_ProcessContract(t *testing.T) {
	ownAddress := common.HexToAddress("0x000062be10f46b5d01Ecd9b502c4bA3d6131f6fc2e41")
	log.InitLogger(log.LvlDebug)
	transactionJson := "{\"TxData\":{\"nonce\":\"0x0\",\"to\":\"0x00120000000000000000000000000000000000000000\",\"hashlock\":null,\"timelock\":\"0x0\",\"value\":\"0x2540be400\",\"fee\":\"0x69db9c0\",\"gasPrice\":\"0xa\",\"gas\":\"0x1027127dc00\",\"input\":\"0xf9027b823138b8eb0061736d01000000010d0360017f0060027f7f00600000021d0203656e76067072696e7473000003656e76087072696e74735f6c00010304030202000405017001010105030100020615037f01419088040b7f00419088040b7f004186080b073405066d656d6f727902000b5f5f686561705f6261736503010a5f5f646174615f656e64030204696e697400030568656c6c6f00040a450302000b02000b3d01017f230041106b220124004180081000200141203a000f2001410f6a41011001200010002001410a3a000e2001410e6a41011001200141106a24000b0b0d01004180080b0668656c6c6f00b901887b22616269417272223a5b0a202020207b0a2020202020202020226e616d65223a2022696e6974222c0a202020202020202022696e70757473223a205b5d2c0a2020202020202020226f757470757473223a205b5d2c0a202020202020202022636f6e7374616e74223a202266616c7365222c0a20202020202020202274797065223a202266756e6374696f6e220a202020207d2c0a202020207b0a2020202020202020226e616d65223a202268656c6c6f222c0a202020202020202022696e70757473223a205b0a2020202020202020202020207b0a20202020202020202020202020202020226e616d65223a20226e616d65222c0a202020202020202020202020202020202274797065223a2022737472696e67220a2020202020202020202020207d0a20202020202020205d2c0a2020202020202020226f757470757473223a205b5d2c0a202020202020202022636f6e7374616e74223a202274727565222c0a20202020202020202274797065223a202266756e6374696f6e220a202020207d0a5d7d0a\"},\"Wit\":{\"r\":\"0x30e173f7590a6e12bb4d51bbf6ae113ee668245d2e30a145d1845f55ae5a9f4a\",\"s\":\"0x7d9f36d62573ac09e1dd84d31650a8b5e20b5dffb34d3955dde224c61d299744\",\"v\":\"0x39\",\"hashkey\":\"0x\"}}"

	var tx model.Transaction
	err := json.Unmarshal([]byte(transactionJson), &tx)
	if err != nil {
		log.Info("TestAccountStateDB_ProcessContract", "err", err)
	}
	log.Info("processContract", "Tx", tx)

	tx.PaddingTxIndex(0)
	block := createBlock(1,common.Hash{},[]*model.Transaction{&tx} )

	db, root := createTestStateDB()
	processor, err := NewAccountStateDB(root, NewStateStorageWithCache(db))
	assert.NoError(t, err)

	processor.NewAccountState(ownAddress)
	err = processor.AddNonce(ownAddress, 0)
	processor.AddBalance(ownAddress,  new(big.Int).SetInt64(int64(1000000000000000000)))

	assert.NoError(t,err)
	balance, err := processor.GetBalance(ownAddress)
	nonce, err := processor.GetNonce(ownAddress)
	log.Info("balance", "balance", balance.String())
	//log.Info("nonce", "nonce", nonce, "tx.nonce", tx.Nonce())

	gasLimit := gasLimit * 10000000000
	log.Info("gasLimit", "gasLimit", gasLimit)


	receipt, err := processor.ProcessContract(&tx, block.Header().(*model.Header), true,fakeGetBlockHash)
	assert.NoError(t, err)
	log.Info("result", "receipt", receipt)
	assert.Equal(t, true, receipt.HandlerResult)
	tx.PaddingReceipt(receipt)
	receiptResult, err := tx.GetReceipt()
	assert.NoError(t, err)
	contractNonce, err := processor.GetNonce(receiptResult.ContractAddress)
	log.Info("TestAccountStateDB_ProcessContract", "contractNonce", contractNonce, "receiptResult", receiptResult)
	code, err := processor.GetCode(receiptResult.ContractAddress)
	assert.NoError(t, err)
	assert.Equal(t, code, tx.ExtraData())



    sw, err := soft_wallet.NewSoftWallet()
    sw.Open("/Users/konggan/tmp/dipperin_apps/node/CSWallet", "CSWallet","123")


	callTx, err := newContractCallTx(nil, &receiptResult.ContractAddress, new(big.Int).SetUint64(1),uint64(1500000), "hello", "name", nonce+1, code)
	account := accounts.Account{ownAddress}
	signCallTx, err := sw.SignTx(account, callTx, nil )



	assert.NoError(t, err)
	callTx.PaddingTxIndex(0)
	block2 := createBlock(2,common.Hash{},[]*model.Transaction{signCallTx} )
	log.Info("callTx info", "callTx", callTx)
	callRecipt, err := processor.ProcessContract(signCallTx, block2.Header().(*model.Header), false,fakeGetBlockHash)
	assert.NoError(t, err)
	log.Info("TestAccountStateDB_ProcessContract2", "callRecipt", callRecipt, "err", err)

}

func newContractCallTx(from *common.Address, to *common.Address, gasPrice *big.Int, gasLimit uint64, funcName string, input string, nonce uint64, code []byte) (tx *model.Transaction, err error)  {
	// RLP([funcName][params])
	inputRlp,err := rlp.EncodeToBytes([]interface{}{
		funcName,input,
	})
	if err != nil {
		log.Error("input rlp err")
		return
	}

	extraData, err := vmcommon.ParseAndGetRlpData(code, inputRlp)


	if err != nil {
		log.Error("ParseAndGetRlpData  inputRlp", "err", err)
		return
	}

	tx = model.NewTransactionSc(nonce,to,nil,gasPrice, gasLimit, extraData)
	return tx, nil
}



func TestAccountStateDB_ProcessContract2(t *testing.T) {
	var testPath = "../../vm/event"
	tx := createContractTx(t, testPath+"/event.wasm", testPath+"/event.cpp.abi.json")

	db, root := createTestStateDB()
	processor, err := NewAccountStateDB(root, NewStateStorageWithCache(db))
	assert.NoError(t, err)

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
	assert.NoError(t, err)
	nonce, err := processor.GetNonce(fullReceipt.ContractAddress)
	fmt.Println(nonce, err)
	code,err:=processor.GetCode(fullReceipt.ContractAddress)
	trueCode := getContractCode(t,testPath+"/event.wasm", testPath+"/event.cpp.abi.json")
	assert.Equal(t,trueCode,code)
}