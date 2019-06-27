package vm

import (
	"github.com/dipperin/dipperin-core/third-party/log"
	"testing"
	"path/filepath"
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/stretchr/testify/assert"
	"github.com/dipperin/dipperin-core/tests/node-cluster"
	"github.com/dipperin/dipperin-core/common"
)

var (
	WASMConvertPath = filepath.Join(util.HomeDir(), "go/src/dipperin-c/dipc/testcontract/convert/convert.wasm")
	ABIConvertPath  = filepath.Join(util.HomeDir(), "go/src/dipperin-c/dipc/testcontract/convert/convert.cpp.abi.json")
)

func Test_ConvertContractCall(t *testing.T) {
	//log.InitLogger(log.LvlDebug)
	log.Info("start")
	cluster, err := node_cluster.CreateNodeCluster()
	assert.NoError(t, err)

	nodeName := "default_v0"
	client := cluster.NodeClient[nodeName]

	contractHash := CreateConvertContract(t, cluster, nodeName)
	checkTransactionOnChain(client, []common.Hash{contractHash})

/*	data := getCallExtraData(t, "getBlockInfo", "")
	txHash := CallConvertContract(t, cluster, nodeName, contractHash, data)
	checkTransactionOnChain(client, []common.Hash{txHash})*/

	data := getCallExtraData(t, "printTest", "")
	txHash := CallConvertContract(t, cluster, nodeName, contractHash, data)
	checkTransactionOnChain(client, []common.Hash{txHash})
}

func CreateConvertContract(t *testing.T, cluster *node_cluster.NodeCluster, nodeName string) common.Hash {
	client := cluster.NodeClient[nodeName]
	from, err := cluster.GetNodeMainAddress(nodeName)
	LogTestPrint("Test", "From", "addr", from.Hex())
	assert.NoError(t, err)

	to := common.HexToAddress(common.AddressContractCreate)
	data := GetCreateExtraData(t, WASMConvertPath, ABIConvertPath, "")
	txHash, innerErr := SendTransactionContract(client, from, to, value, gasLimit, gasPrice, data)
	assert.NoError(t, innerErr)
	return txHash
}

func CallConvertContract(t *testing.T, cluster *node_cluster.NodeCluster, nodeName string, txHash common.Hash, input []byte) common.Hash {
	client := cluster.NodeClient[nodeName]
	from, err := cluster.GetNodeMainAddress(nodeName)
	LogTestPrint("Test", "From", "addr", from.Hex())
	assert.NoError(t, err)

	to := GetContractAddressByTxHash(client, txHash)
	txHash, innerErr := SendTransactionContract(client, from, to, value, gasLimit, gasPrice, input)
	assert.NoError(t, innerErr)
	return txHash
}
