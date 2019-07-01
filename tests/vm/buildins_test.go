package vm

import (
	"github.com/dipperin/dipperin-core/third-party/log"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/dipperin/dipperin-core/tests/node-cluster"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/tests/g-testData"
)

func Test_BuildInsContractCall(t *testing.T) {
	log.Info("start")
	cluster, err := node_cluster.CreateNodeCluster()
	assert.NoError(t, err)

	nodeName := "default_v0"
	client := cluster.NodeClient[nodeName]

	WASMBuildPath := g_testData.GetWasmPath("buildins")
	ABIBuildPath := g_testData.GetAbiPath("buildins")
	contractHash := SendCreateContract(t, cluster, nodeName, WASMBuildPath, ABIBuildPath)
	checkTransactionOnChain(client, []common.Hash{contractHash})

	data ,err:= g_testData.GetCallExtraData("arithmeticTest", "")
	assert.NoError(t,err)
	txHash := SendCallContract(t, cluster, nodeName, contractHash, data)
	checkTransactionOnChain(client, []common.Hash{txHash})
}
