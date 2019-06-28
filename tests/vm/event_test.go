package vm

import (
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/tests/node-cluster"
	"github.com/stretchr/testify/assert"
	"testing"
	"path/filepath"
	"github.com/dipperin/dipperin-core/common/util"
)

var (
	AbiPath  = filepath.Join(util.HomeDir(), "go/src/github.com/dipperin/dipperin-core/core/vm/event/event.cpp.abi.json")
	WASMPath = filepath.Join(util.HomeDir(), "go/src/github.com/dipperin/dipperin-core/core/vm/event/event.wasm")
)

func Test_EventContractCall(t *testing.T) {
	cluster, err := node_cluster.CreateNodeCluster()
	assert.NoError(t, err)

	nodeName := "default_v0"
	client := cluster.NodeClient[nodeName]
	txHashList := CreateContract(t, cluster, nodeName, 1)
	checkTransactionOnChain(client, txHashList)

	// 根据交易ID获取合约地址
	var addrList []common.Address
	for i := 0; i < len(txHashList); i++ {
		addr := GetContractAddressByTxHash(client, txHashList[i])
		addrList = append(addrList, addr)
	}

	txHashList = CallContract(t, cluster, nodeName, addrList)
	checkTransactionOnChain(client, txHashList)
}

func CreateContract(t *testing.T, cluster *node_cluster.NodeCluster, nodeName string, times int) []common.Hash {
	client := cluster.NodeClient[nodeName]
	from, err := cluster.GetNodeMainAddress(nodeName)
	LogTestPrint("Test", "From", "addr", from.Hex())
	assert.NoError(t, err)

	to := common.HexToAddress(common.AddressContractCreate)
	data := getCreateExtraData(t, WASMPath, AbiPath, "")
	var txHashList []common.Hash
	for i := 0; i < times; i++ {
		txHash, innerErr := SendTransactionContract(client, from, to, value, gasLimit, gasPrice, data)
		assert.NoError(t, innerErr)
		txHashList = append(txHashList, txHash)
	}
	return txHashList
}

func CallContract(t *testing.T, cluster *node_cluster.NodeCluster, nodeName string, addrList []common.Address) []common.Hash {
	client := cluster.NodeClient[nodeName]
	from, err := cluster.GetNodeMainAddress(nodeName)
	LogTestPrint("Test", "From", "addr", from.Hex())
	assert.NoError(t, err)

	var txHashList []common.Hash
	for i := 0; i < len(addrList); i++ {
		input := getCallExtraData(t, "hello", fmt.Sprintf("Event,%v", 100*i))
		txHash, innerErr := SendTransactionContract(client, from, addrList[i], value, gasLimit, gasPrice, input)
		assert.NoError(t, innerErr)
		txHashList = append(txHashList, txHash)
	}
	return txHashList
}