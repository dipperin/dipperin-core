package vm

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/tests/g-testData"
	"github.com/dipperin/dipperin-core/tests/node-cluster"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ConvertContractCall(t *testing.T) {
	cluster, err := node_cluster.CreateNodeCluster()
	assert.NoError(t, err)

	nodeName := "default_v0"
	client := cluster.NodeClient[nodeName]

	WASMConvertPath := g_testData.GetWASMPath("convert", g_testData.CoreVmTestData)
	ABIConvertPath := g_testData.GetAbiPath("convert", g_testData.CoreVmTestData)
	contractHash := SendCreateContract(t, cluster, nodeName, WASMConvertPath, ABIConvertPath, "")
	checkTransactionOnChain(client, []common.Hash{contractHash})

	data, err := g_testData.GetCallExtraData("getBlockInfo", "")
	assert.NoError(t, err)
	txHash := SendCallContract(t, cluster, nodeName, contractHash, data)
	checkTransactionOnChain(client, []common.Hash{txHash})
}
