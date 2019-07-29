package vm

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/tests/g-testData"
	"github.com/dipperin/dipperin-core/tests/node-cluster"
	"github.com/stretchr/testify/assert"
	"testing"
)

func CreatAndCallContract(t *testing.T, parameter *g_testData.ContractTestParameter) (common.Hash, error) {
	cluster, err := node_cluster.CreateNodeCluster()
	assert.NoError(t, err)
	client := cluster.NodeClient[parameter.NodeName]

	contractHash := SendCreateContract(t, cluster, parameter.NodeName, parameter.WASMPath, parameter.AbiPath, parameter.InitInputPara)
	checkTransactionOnChain(client, []common.Hash{contractHash})

	data, err := g_testData.GetCallExtraData(parameter.CallFuncName, parameter.CallInputPara)
	assert.NoError(t, err)
	txHash := SendCallContract(t, cluster, parameter.NodeName, contractHash, data)
	checkTransactionOnChain(client, []common.Hash{txHash})

	return contractHash, nil
}

func Test_EventContractCall(t *testing.T) {
	cluster, err := node_cluster.CreateNodeCluster()
	assert.NoError(t, err)

	nodeName := "default_v0"
	client := cluster.NodeClient[nodeName]

	WASMEventPath := g_testData.GetWASMPath("event", g_testData.CoreVmTestData)
	AbiEventPath := g_testData.GetAbiPath("event", g_testData.CoreVmTestData)
	contractHash := SendCreateContract(t, cluster, nodeName, WASMEventPath, AbiEventPath, "")
	checkTransactionOnChain(client, []common.Hash{contractHash})

	data, err := g_testData.GetCallExtraData("hello", "money,100")
	assert.NoError(t, err)
	txHash := SendCallContract(t, cluster, nodeName, contractHash, data)
	checkTransactionOnChain(client, []common.Hash{txHash})
}

func Test_DIPCLibContract(t *testing.T) {
	CreatAndCallContract(t, &g_testData.ContractTestPar)
}
