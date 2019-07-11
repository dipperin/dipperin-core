package vm

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/tests/g-testData"
	"github.com/dipperin/dipperin-core/tests/node-cluster"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_DipcLibTestCall(t *testing.T) {
	cluster, err := node_cluster.CreateNodeCluster()
	assert.NoError(t, err)

	nodeName := "default_v0"
	client := cluster.NodeClient[nodeName]

	WASMConvertPath := g_testData.GetWasmPath("dipclib_test")
	ABIConvertPath := g_testData.GetAbiPath("dipclib_test")
	contractHash := SendCreateContract(t, cluster, nodeName, WASMConvertPath, ABIConvertPath, "")
	checkTransactionOnChain(client, []common.Hash{contractHash})

	data, err := g_testData.GetCallExtraData("libTest", "")
	assert.NoError(t, err)
	txHash := SendCallContract(t, cluster, nodeName, contractHash, data)
	checkTransactionOnChain(client, []common.Hash{txHash})
}
