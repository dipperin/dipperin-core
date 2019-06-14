package vm

import (
	"math/big"
	"github.com/dipperin/dipperin-core/third-party/rpc"
	"github.com/dipperin/dipperin-core/common"
	"strings"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/dipperin/dipperin-core/core/rpc-interface"
	"github.com/dipperin/dipperin-core/core/vm/model"
	"io/ioutil"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/ethereum/go-ethereum/rlp"
	"path/filepath"
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/common/vmcommon"
	"github.com/dipperin/dipperin-core/core/vm/common/utils"
	"time"
	"fmt"
)

var (
	AbiPath  = filepath.Join(util.HomeDir(), "c++/src/dipc/testcontract/cppfile/event/event.cpp.abi.json")
	WASMPath = filepath.Join(util.HomeDir(), "c++/src/dipc/testcontract/cppfile/event/event.wasm")
)

var (
	AbiTokenPath  = filepath.Join(util.HomeDir(), "go/src/github.com/dipperin/dipperin-core/core/vm/event/token/token.cpp.abi.json")
	WASMTokenPath = filepath.Join(util.HomeDir(), "go/src/github.com/dipperin/dipperin-core/core/vm/event/token/token-jw.wasm")
)

func LogTestPrint(function, msg string, ctx ...interface{}) {
	printMsg := "[~wjw~" + function + "]" + msg
	log.Info(printMsg, ctx...)
}

func GetRpcTXMethod(methodName string) string {
	return "dipperin_" + strings.ToLower(methodName[0:1]) + methodName[1:]
}

func SendTransaction(client *rpc.Client, from, to common.Address, value, fee *big.Int, data []byte) (common.Hash, error) {
	var resp common.Hash
	if err := client.Call(&resp, GetRpcTXMethod("SendTransaction"), from, to, value, fee, data, nil); err != nil {
		LogTestPrint("Test", "SendTransaction failed", "err", err)
		return common.Hash{}, err
	}
	LogTestPrint("Test", "SendTransaction Successful", "txId", resp.Hex())
	return resp, nil
}

func SendTransactionContract(client *rpc.Client, from, to common.Address, value, gasLimit, gasPrice *big.Int, data []byte) (common.Hash, error) {
	var resp common.Hash
	if err := client.Call(&resp, GetRpcTXMethod("SendTransactionContract"), from, to, value, gasLimit, gasPrice, data, nil); err != nil {
		LogTestPrint("Test", "SendContract failed", "err", err)
		return common.Hash{}, err
	}
	LogTestPrint("Test", "SendContract Successful", "txId", resp.Hex())
	return resp, nil
}

func Transaction(client *rpc.Client, hash common.Hash) (bool, uint64) {
	var resp *rpc_interface.TransactionResp
	if err := client.Call(&resp, GetRpcTXMethod("Transaction"), hash); err != nil {
		return false, 0
	}
	if resp.BlockNumber == 0 {
		return false, 0
	}
	return true, resp.BlockNumber
}

func GetConvertReceiptByTxHash(client *rpc.Client, hash common.Hash) *model.Receipt {
	var resp *model.Receipt
	if err := client.Call(&resp, GetRpcTXMethod("GetConvertReceiptByTxHash"), hash); err != nil {
		LogTestPrint("Test", "call GetConvertReceiptByTxHash failed", "err", err)
		return nil
	}
	return resp
}

func GetReceiptByTxHash(client *rpc.Client, hash common.Hash) *model.Receipt {
	var resp *model.Receipt
	if err := client.Call(&resp, GetRpcTXMethod("GetReceiptByTxHash"), hash); err != nil {
		LogTestPrint("Test", "call GetReceiptByTxHash failed", "err", err)
		return nil
	}
	return resp
}

func GetReceiptsByBlockNum(client *rpc.Client, num uint64) model.Receipts {
	var resp model.Receipts
	if err := client.Call(&resp, GetRpcTXMethod("GetReceiptsByBlockNum"), num); err != nil {
		LogTestPrint("Test", "call GetReceiptsByBlockNum failed", "err", err)
		return nil
	}
	return resp
}

func GetContractAddressByTxHash(client *rpc.Client, hash common.Hash) common.Address {
	var resp common.Address
	if err := client.Call(&resp, GetRpcTXMethod("GetContractAddressByTxHash"), hash); err != nil {
		LogTestPrint("Test", "call GetContractAddressByTxHash failed", "err", err)
		return common.Address{}
	}
	return resp
}

func GetBlockByNumber(client *rpc.Client, num uint64) rpc_interface.BlockResp {
	var respBlock rpc_interface.BlockResp
	if err := client.Call(&respBlock, GetRpcTXMethod("GetBlockByNumber"), num); err != nil {
		LogTestPrint("Test", "call GetBlockByNumber failed", "err", err)
		return rpc_interface.BlockResp{}
	}
	return respBlock
}

func getCallExtraData(t *testing.T, funcName, param string) []byte {
	input := []interface{}{
		funcName,
		param,
	}

	result, err := rlp.EncodeToBytes(input)
	assert.NoError(t, err)
	return result
}

func getCreateExtraData(t *testing.T, abiPath, wasmPath string, params []string) []byte {
	// GetContractExtraData
	log.Info("the abiPath is:%v","abiPath",abiPath)
	abiBytes, err := ioutil.ReadFile(abiPath)
	assert.NoError(t, err)

	//log.Info("the abiBytes is:","abiBytes",hexutil.Encode(abiBytes))

	var wasmAbi utils.WasmAbi
	err = wasmAbi.FromJson(abiBytes)
	assert.NoError(t, err)
	var args []utils.InputParam
	for _, v := range wasmAbi.AbiArr {
		if strings.EqualFold("init", v.Name) && strings.EqualFold(v.Type, "function") {
			args = v.Inputs
		}
	}
	//params := []string{"dipp", "DIPP", "100000000"}
	wasmBytes, err := ioutil.ReadFile(wasmPath)
	//log.Info("the wasmBytes is:","wasmBytes",hexutil.Encode(wasmBytes))
	assert.NoError(t, err)
	rlpParams := []interface{}{
		wasmBytes, abiBytes,
	}

	//log.Info("the params is:","params",params)
	assert.Equal(t, len(params), len(args))
	for i, v := range args {
		bts := params[i]
		re, err := vmcommon.StringConverter(bts, v.Type)
		assert.NoError(t, err)
		rlpParams = append(rlpParams, re)
		//inputParams = append(inputParams, re)
	}
	data, err := rlp.EncodeToBytes(rlpParams)
	//input, err := rlp.EncodeToBytes(inputParams)
	assert.NoError(t, err)

	//log.Info("the generate extra data is:","extraData",hexutil.Encode(data))
	return data
}

func checkTransactionOnChain(client *rpc.Client, txHashList []common.Hash) {
	for i := 0; i < len(txHashList); i++ {
		for {
			result, num := Transaction(client, txHashList[i])
			if result {
				receipts := GetConvertReceiptByTxHash(client, txHashList[i])
				LogTestPrint("Test", "CallTransaction", "blockNum", num)
				fmt.Println(receipts)
				break
			}
			time.Sleep(time.Second * 2)
		}
		time.Sleep(time.Millisecond * 100)
	}
}
