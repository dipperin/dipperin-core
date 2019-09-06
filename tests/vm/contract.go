package vm

import (
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/hexutil"
	"github.com/dipperin/dipperin-core/core/rpc-interface"
	"github.com/dipperin/dipperin-core/core/vm/model"
	"github.com/dipperin/dipperin-core/tests/g-testData"
	"github.com/dipperin/dipperin-core/tests/node-cluster"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/dipperin/dipperin-core/third-party/rpc"
	"github.com/stretchr/testify/assert"
	"math/big"
	"strings"
	"testing"
	"time"
)

func LogTestPrint(function, msg string, ctx ...interface{}) {
	printMsg := "[~wjw~" + function + "]" + msg
	log.Info(printMsg, ctx...)
}

func GetRpcTXMethod(methodName string) string {
	return "dipperin_" + strings.ToLower(methodName[0:1]) + methodName[1:]
}

func SendTransaction(client *rpc.Client, from, to common.Address, value, gasPrice *big.Int, gasLimit uint64, data []byte) (common.Hash, error) {
	var resp common.Hash
	if err := client.Call(&resp, GetRpcTXMethod("SendTransaction"), from, to, value, gasPrice, gasLimit, data, nil); err != nil {
		LogTestPrint("Test", "SendTransaction failed", "err", err)
		return common.Hash{}, err
	}
	LogTestPrint("Test", "SendTransaction Successful", "txId", resp.Hex())
	return resp, nil
}

func SendTransactionContract(client *rpc.Client, from, to common.Address, value, gasPrice *big.Int, gasLimit uint64, data []byte) (common.Hash, error) {
	var resp common.Hash
	if err := client.Call(&resp, GetRpcTXMethod("SendTransactionContract"), from, to, value, gasPrice, gasLimit, data, nil); err != nil {
		LogTestPrint("Test", "SendTransactionContract failed", "err", err)
		return common.Hash{}, err
	}
	LogTestPrint("Test", "SendTransactionContract Successful", "txId", resp.Hex())
	return resp, nil
}

func Call(client *rpc.Client, from, to common.Address, data []byte) (string, error) {
	var resp string
	if err := client.Call(&resp, GetRpcTXMethod("CallContract"), from, to, data, uint64(0)); err != nil {
		LogTestPrint("Test", "CallContract failed", "err", err)
		return "", err
	}
	LogTestPrint("Test", "CallContract Successful", "resp", resp)
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

func CurrentBalance(client *rpc.Client, addr common.Address) *rpc_interface.CurBalanceResp {
	var resp *rpc_interface.CurBalanceResp
	if err := client.Call(&resp, GetRpcTXMethod("CurrentBalance"), addr); err != nil {
		LogTestPrint("Test", "call CurrentBalance failed", "err", err)
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

func SendCreateContract(t *testing.T, cluster *node_cluster.NodeCluster, nodeName, wasmPath, abiPath, params string) common.Hash {
	client := cluster.NodeClient[nodeName]
	from, err := cluster.GetNodeMainAddress(nodeName)
	LogTestPrint("Test", "From", "addr", from.Hex())
	assert.NoError(t, err)

	to := common.HexToAddress(common.AddressContractCreate)
	data, err := g_testData.GetCreateExtraData(wasmPath, abiPath, params)
	assert.NoError(t, err)

	log.Info("SendCreateContract the extraData is:", "data", hexutil.Encode(data))

	value := big.NewInt(0)
	gasLimit := g_testData.TestGasLimit * 10000
	gasPrice := big.NewInt(0).Mul(g_testData.TestGasPrice, big.NewInt(5))
	txHash, innerErr := SendTransactionContract(client, from, to, value, gasPrice, gasLimit, data)
	assert.NoError(t, innerErr)
	return txHash
}

func SendCallContract(t *testing.T, cluster *node_cluster.NodeCluster, nodeName string, txHash common.Hash, input []byte, value *big.Int) common.Hash {
	client := cluster.NodeClient[nodeName]
	from, err := cluster.GetNodeMainAddress(nodeName)
	LogTestPrint("Test", "From", "addr", from.Hex())
	assert.NoError(t, err)

	to := GetContractAddressByTxHash(client, txHash)
	gasLimit := g_testData.TestGasLimit * 10000
	gasPrice := big.NewInt(0).Mul(g_testData.TestGasPrice, big.NewInt(2))
	txHash, innerErr := SendTransactionContract(client, from, to, value, gasPrice, gasLimit, input)
	assert.NoError(t, innerErr)
	return txHash
}

func checkTransactionOnChain(client *rpc.Client, txHashList []common.Hash) {
	for i := 0; i < len(txHashList); i++ {
		for {
			result, num := Transaction(client, txHashList[i])
			if result {
				receipts := GetReceiptByTxHash(client, txHashList[i])
				LogTestPrint("Test", "CallTransaction", "blockNum", num)
				fmt.Println(receipts)
				break
			}
			time.Sleep(time.Second * 2)
		}
		time.Sleep(time.Millisecond * 100)
	}
}
