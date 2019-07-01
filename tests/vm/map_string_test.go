package vm

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/dipperin/dipperin-core/tests/node-cluster"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/tests/g-testData"
)

func Test_MapContractCall(t *testing.T) {
	cluster, err := node_cluster.CreateNodeCluster()
	assert.NoError(t, err)

	nodeName := "default_v0"
	client := cluster.NodeClient[nodeName]

	WASMMapPath := g_testData.GetWasmPath("map-string")
	ABIMapPath := g_testData.GetAbiPath("map-string")
	contractHash := SendCreateContract(t, cluster, nodeName, WASMMapPath, ABIMapPath)
	checkTransactionOnChain(client, []common.Hash{contractHash})

	data := g_testData.GetCallExtraData(t, "setBalance", "balance,100")
	txHash := SendCallContract(t, cluster, nodeName, contractHash, data)
	checkTransactionOnChain(client, []common.Hash{txHash})

	from, err := cluster.GetNodeMainAddress(nodeName)
	assert.NoError(t, err)
	to := GetContractAddressByTxHash(client, contractHash)
	data = g_testData.GetCallExtraData(t, "getBalance", "balance")
	err = Call(client, from, to, data)
	assert.NoError(t, err)
}
