package vm

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/tests/g-testData"
	"github.com/dipperin/dipperin-core/tests/node-cluster"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_EventContractCall(t *testing.T) {
	cluster, err := node_cluster.CreateNodeCluster()
	assert.NoError(t, err)

	nodeName := "default_v0"
	client := cluster.NodeClient[nodeName]

	WASMEventPath := g_testData.GetWasmPath("event")
	AbiEventPath := g_testData.GetAbiPath("event")
	contractHash := SendCreateContract(t, cluster, nodeName, WASMEventPath, AbiEventPath, "")
	checkTransactionOnChain(client, []common.Hash{contractHash})

	data, err := g_testData.GetCallExtraData("hello", "money,100")
	assert.NoError(t, err)
	txHash := SendCallContract(t, cluster, nodeName, contractHash, data)
	checkTransactionOnChain(client, []common.Hash{txHash})
}
