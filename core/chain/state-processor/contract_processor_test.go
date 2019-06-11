package state_processor

import (
	"encoding/json"
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/common/vmcommon"
	"github.com/dipperin/dipperin-core/core/accounts"
	"github.com/dipperin/dipperin-core/core/accounts/soft-wallet"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/core/vm/common/utils"
	"github.com/dipperin/dipperin-core/third-party/crypto/cs-crypto"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"math/big"
	"strings"
	"testing"
)

func TestAccountStateDB_ProcessContract(t *testing.T) {
	ownAddress := common.HexToAddress("0x000062be10f46b5d01Ecd9b502c4bA3d6131f6fc2e41")
	log.InitLogger(log.LvlDebug)
	transactionJson := "{\"TxData\":{\"nonce\":\"0x0\",\"to\":\"0x00120000000000000000000000000000000000000000\",\"hashlock\":null,\"timelock\":\"0x0\",\"value\":\"0x2540be400\",\"fee\":\"0x69db9c0\",\"gasPrice\":\"0xa\",\"gas\":\"0x1027127dc00\",\"input\":\"0xf9026db8eb0061736d01000000010d0360017f0060027f7f00600000021d0203656e76067072696e7473000003656e76087072696e74735f6c00010304030202000405017001010105030100020615037f01419088040b7f00419088040b7f004186080b073405066d656d6f727902000b5f5f686561705f6261736503010a5f5f646174615f656e64030204696e697400030568656c6c6f00040a450302000b02000b3d01017f230041106b220124004180081000200141203a000f2001410f6a41011001200010002001410a3a000e2001410e6a41011001200141106a24000b0b0d01004180080b0668656c6c6f00b9017d5b0a202020207b0a2020202020202020226e616d65223a2022696e6974222c0a202020202020202022696e70757473223a205b5d2c0a2020202020202020226f757470757473223a205b5d2c0a202020202020202022636f6e7374616e74223a202266616c7365222c0a20202020202020202274797065223a202266756e6374696f6e220a202020207d2c0a202020207b0a2020202020202020226e616d65223a202268656c6c6f222c0a202020202020202022696e70757473223a205b0a2020202020202020202020207b0a20202020202020202020202020202020226e616d65223a20226e616d65222c0a202020202020202020202020202020202274797065223a2022737472696e67220a2020202020202020202020207d0a20202020202020205d2c0a2020202020202020226f757470757473223a205b5d2c0a202020202020202022636f6e7374616e74223a202274727565222c0a20202020202020202274797065223a202266756e6374696f6e220a202020207d0a5d0a\"},\"Wit\":{\"r\":\"0xa1509f3efb1e632643c9972b9183234445c539a1b483ad0ea4b36a4edabf8d04\",\"s\":\"0xa7a16d72b826aea44e8f56247abbad367cf7e300d564949e66ac97098b9f234\",\"v\":\"0x39\",\"hashkey\":\"0x\"}}"

	var tx model.Transaction
	err := json.Unmarshal([]byte(transactionJson), &tx)
	if err != nil {
		log.Info("TestAccountStateDB_ProcessContract", "err", err)
	}
	log.Info("processContract", "Tx", tx)

	tx.PaddingTxIndex(0)
	gasLimit := gasLimit * 10000000000
	block := createBlock(1, common.Hash{}, []*model.Transaction{&tx}, gasLimit)

	db, root := createTestStateDB()
	processor, err := NewAccountStateDB(root, NewStateStorageWithCache(db))
	assert.NoError(t, err)

	processor.NewAccountState(ownAddress)
	err = processor.AddNonce(ownAddress, 0)
	processor.AddBalance(ownAddress, new(big.Int).SetInt64(int64(1000000000000000000)))

	assert.NoError(t, err)
	balance, err := processor.GetBalance(ownAddress)
	nonce, err := processor.GetNonce(ownAddress)
	log.Info("balance", "balance", balance.String())
	//log.Info("nonce", "nonce", nonce, "tx.nonce", tx.Nonce())

	log.Info("gasLimit", "gasLimit", gasLimit)

	gasUsed := uint64(0);
	txConfigCreate := &TxProcessConfig  {
		Tx:&tx,
		Header:  block.Header().(*model.Header),
		GetHash:fakeGetBlockHash,
		GasLimit:&gasLimit,
		GasUsed:&gasUsed,
	}

	receipt, err := processor.ProcessContract(txConfigCreate, true)
	assert.NoError(t, err)
	log.Info("result", "receipt", receipt)
	assert.Equal(t, true, receipt.HandlerResult)
	tx.PaddingReceipt(receipt)
	receiptResult, err := tx.GetReceipt()
	assert.NoError(t, err)
	contractNonce, err := processor.GetNonce(receiptResult.ContractAddress)
	log.Info("TestAccountStateDB_ProcessContract", "contractNonce", contractNonce, "receiptResult", receiptResult)
	code, err := processor.GetCode(receiptResult.ContractAddress)
	log.Info("TestAccountStateDB_ProcessContract", "code  get from state", code)
	assert.NoError(t, err)
	assert.Equal(t, code, tx.ExtraData())

	sw, err := soft_wallet.NewSoftWallet()
	err = sw.Open(util.HomeDir() +"/go/src/github.com/dipperin/dipperin-core/core/vm/event/CSWallet", "CSWallet", "123")
	assert.NoError(t, err)

	callTx, err := newContractCallTx(nil, &receiptResult.ContractAddress, new(big.Int).SetUint64(1), uint64(1500000), "hello", "name", nonce+1, code)
	account := accounts.Account{ownAddress}
	signCallTx, err := sw.SignTx(account, callTx, nil)

	assert.NoError(t, err)
	callTx.PaddingTxIndex(0)
	block2 := createBlock(2, common.Hash{}, []*model.Transaction{signCallTx}, gasLimit)
	log.Info("callTx info", "callTx", callTx)


	gasUsed2 := uint64(0);
	txConfig := &TxProcessConfig  {
		Tx:signCallTx,
		Header:  block2.Header().(*model.Header),
		GetHash:fakeGetBlockHash,
		GasLimit:&gasLimit,
		GasUsed:&gasUsed2,
	}


	callRecipt, err := processor.ProcessContract(txConfig, false)
	//assert.NoError(t, err)
	log.Info("TestAccountStateDB_ProcessContract++", "callRecipt", callRecipt, "err", err)

}

func newContractCallTx(from *common.Address, to *common.Address, gasPrice *big.Int, gasLimit uint64, funcName string, input string, nonce uint64, code []byte) (tx *model.Transaction, err error) {
	// RLP([funcName][params])
	inputRlp, err := rlp.EncodeToBytes([]interface{}{
		funcName, input,
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

	tx = model.NewTransactionSc(nonce, to, nil, gasPrice, gasLimit, extraData)
	return tx, nil
}

func TestAccountStateDB_ProcessContract2(t *testing.T) {
	var testPath = "../../vm/event"
	tx1 := createContractTx(t, testPath+"/event.wasm", testPath+"/event.cpp.abi.json")
	contractAddr := cs_crypto.CreateContractAddress(aliceAddr, 0)
	name := []byte("ProcessContract")
	num := vmcommon.Int64ToBytes(456)
	param := [][]byte{name, num}
	tx2 := callContractTx(t, &contractAddr, "hello", param, 1)

	db, root := createTestStateDB()
	tdb := NewStateStorageWithCache(db)
	processor, err := NewAccountStateDB(root, tdb)
	assert.NoError(t, err)

	block := createBlock(1, common.Hash{}, []*model.Transaction{tx1, tx2}, 5*gasLimit)
	tmpGasLimit := block.GasLimit()
	gasUsed := block.GasUsed()
	config := &TxProcessConfig{
		Tx:      tx1,
		Header:  block.Header(),
		GetHash: getTestHashFunc(),
		GasLimit: &tmpGasLimit,
		GasUsed: &gasUsed,
		TxFee: big.NewInt(0),
	}
	err = processor.ProcessTxNew(config)
	assert.NoError(t, err)

	receipt1, err := tx1.GetReceipt()
	assert.NoError(t, err)
	assert.Equal(t, tx1.CalTxId(), receipt1.TxHash)
	assert.Equal(t, cs_crypto.CreateContractAddress(aliceAddr, 0), receipt1.ContractAddress)
	assert.Len(t, receipt1.Logs, 0)

	fmt.Println("---------------------------")

	config.Tx = tx2
	err = processor.ProcessTxNew(config)
	assert.NoError(t, err)

	receipt2, err := tx2.GetReceipt()
	assert.NoError(t, err)
	assert.Equal(t, tx2.CalTxId(), receipt2.TxHash)
	assert.Equal(t, receipt1.ContractAddress, receipt2.ContractAddress)
	assert.Len(t, receipt2.Logs, 1)

	log1 := receipt2.Logs[0]
	assert.Equal(t, tx2.CalTxId(), log1.TxHash)
	assert.Equal(t, common.Hash{}, log1.BlockHash)
	assert.Equal(t, receipt2.ContractAddress, log1.Address)
	assert.Equal(t, uint64(1), log1.BlockNumber)
}


func TestAccountStateDB_ProcessContract3(t *testing.T) {
	ownAddress := common.HexToAddress("0x000062be10f46b5d01Ecd9b502c4bA3d6131f6fc2e41")

	// GetContractExtraData
	abiBytes, err := ioutil.ReadFile("../../vm/event/token/token.cpp.abi.json")
	assert.NoError(t, err)
	var wasmAbi utils.WasmAbi
	err = wasmAbi.FromJson(abiBytes)
	assert.NoError(t, err)

	var args []utils.InputParam
	for _, v := range wasmAbi.AbiArr {
		if strings.EqualFold("init", v.Name) && strings.EqualFold(v.Type, "function") {
			args = v.Inputs
		}
	}

	params := []string{"dipp", "DIPP", "100000000"}
	wasmBytes, err := ioutil.ReadFile("../../vm/event/token/token5.wasm")
	assert.NoError(t, err)

	rlpParams := []interface{}{
		wasmBytes, abiBytes,
	}

	assert.Equal(t, len(params), len(args))
	for i, v := range args {
		bts := params[i]
		re, err := vmcommon.StringConverter(bts, v.Type)
		assert.NoError(t, err)
		rlpParams = append(rlpParams, re)
		//inputParams = append(inputParams, re)
	}

	data, err :=  rlp.EncodeToBytes(rlpParams)
	//input, err := rlp.EncodeToBytes(inputParams)
	assert.NoError(t, err)


	addr := common.HexToAddress(common.AddressContractCreate)

	tx := model.NewTransactionSc(0,&addr,new(big.Int).SetUint64(uint64(10)),new(big.Int).SetUint64(uint64(1)),26427000, data)


	sw, err := soft_wallet.NewSoftWallet()
	err = sw.Open(util.HomeDir() +"/go/src/github.com/dipperin/dipperin-core/core/vm/event/CSWallet", "CSWallet", "123")
	assert.NoError(t, err)

	account := accounts.Account{ownAddress}
	signCreateTx, err := sw.SignTx(account, tx, nil)
	signCreateTx.PaddingTxIndex(0)



	gasLimit := gasLimit * 10000000000
	block := createBlock(1, common.Hash{}, []*model.Transaction{signCreateTx}, gasLimit)

	db, root := createTestStateDB()
	processor, err := NewAccountStateDB(root, NewStateStorageWithCache(db))
	assert.NoError(t, err)

	processor.NewAccountState(ownAddress)
	err = processor.AddNonce(ownAddress, 0)
	processor.AddBalance(ownAddress, new(big.Int).SetInt64(int64(1000000000000000000)))


	tmpGasLimit := block.GasLimit()
	gasUsed := block.GasUsed()
	config := &TxProcessConfig{
		Tx:      tx,
		Header:  block.Header(),
		GetHash: getTestHashFunc(),
		GasLimit: &tmpGasLimit,
		GasUsed: &gasUsed,
		TxFee: big.NewInt(0),
	}
	err = processor.ProcessTxNew(config)
	assert.NoError(t, err)

	receipt1, err := tx.GetReceipt()
	fmt.Println(receipt1)
	assert.NoError(t, err)
	assert.Equal(t, tx.CalTxId(), receipt1.TxHash)
	//assert.Equal(t, cs_crypto.CreateContractAddress(aliceAddr, 0), receipt1.ContractAddress)
	//assert.Len(t, receipt1.Logs, 0)

	fmt.Println("---------------------------")

	//config.Tx = tx2
	//err = processor.ProcessTxNew(config)
	//assert.NoError(t, err)
	//
	//receipt2, err := tx2.GetReceipt()
	//assert.NoError(t, err)
	//assert.Equal(t, tx2.CalTxId(), receipt2.TxHash)
	//assert.Equal(t, receipt1.ContractAddress, receipt2.ContractAddress)
	//assert.Len(t, receipt2.Logs, 1)
	//
	//log1 := receipt2.Logs[0]
	//assert.Equal(t, tx2.CalTxId(), log1.TxHash)
	//assert.Equal(t, common.Hash{}, log1.BlockHash)
	//assert.Equal(t, receipt2.ContractAddress, log1.Address)
	//assert.Equal(t, uint64(1), log1.BlockNumber)
}