package vm

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/dipperin/dipperin-core/tests/node-cluster"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/tests/g-testData"
)

func Test_ConvertContractCall(t *testing.T) {
	cluster, err := node_cluster.CreateNodeCluster()
	assert.NoError(t, err)

	nodeName := "default_v0"
	client := cluster.NodeClient[nodeName]

	WASMConvertPath := g_testData.GetWasmPath("convert")
	ABIConvertPath := g_testData.GetAbiPath("convert")
	contractHash := SendCreateContract(t, cluster, nodeName, WASMConvertPath, ABIConvertPath)
	checkTransactionOnChain(client, []common.Hash{contractHash})

	data := g_testData.GetCallExtraData(t, "getBlockInfo", "")
	txHash := SendCallContract(t, cluster, nodeName, contractHash, data)
	checkTransactionOnChain(client, []common.Hash{txHash})

	data = g_testData.GetCallExtraData(t, "printTest", "")
	txHash = SendCallContract(t, cluster, nodeName, contractHash, data)
	checkTransactionOnChain(client, []common.Hash{txHash})
}
