package vm

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/tests/g-testData"
	"github.com/dipperin/dipperin-core/tests/node-cluster"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_MapContractCall(t *testing.T) {
	cluster, err := node_cluster.CreateNodeCluster()
	assert.NoError(t, err)

	nodeName := "default_v0"
	client := cluster.NodeClient[nodeName]

	WASMMapPath := g_testData.GetWasmPath("map-string")
	ABIMapPath := g_testData.GetAbiPath("map-string")
	contractHash := SendCreateContract(t, cluster, nodeName, WASMMapPath, ABIMapPath, "")
	checkTransactionOnChain(client, []common.Hash{contractHash})

	data, err := g_testData.GetCallExtraData("setBalance", "balance,100")
	assert.NoError(t, err)
	txHash := SendCallContract(t, cluster, nodeName, contractHash, data)

	checkTransactionOnChain(client, []common.Hash{txHash})

	from, err := cluster.GetNodeMainAddress(nodeName)
	assert.NoError(t, err)
	to := GetContractAddressByTxHash(client, contractHash)
	data, err = g_testData.GetCallExtraData("getBalance", "balance")
	assert.NoError(t, err)

	err = Call(client, from, to, data)
	assert.NoError(t, err)
}
