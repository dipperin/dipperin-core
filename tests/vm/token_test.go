package vm

import (
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/tests/node-cluster"
	"github.com/stretchr/testify/assert"
	"testing"
	"github.com/dipperin/dipperin-core/tests/g-testData"
	"path/filepath"
	"github.com/dipperin/dipperin-core/common/util"
)

func Test_TokenContractCall(t *testing.T) {
	cluster, err := node_cluster.CreateNodeCluster()
	assert.NoError(t, err)

	nodeName := "default_v0"
	client := cluster.NodeClient[nodeName]

	WASMTokenPath := g_testData.GetWasmPath("token-const")
	AbiTokenPath := g_testData.GetAbiPath("token-const")
	contractHash := SendCreateContract(t, cluster, nodeName, WASMTokenPath, AbiTokenPath)
	checkTransactionOnChain(client, []common.Hash{contractHash})

	// Transfer money
	aliceAddr := "0x00005586B883Ec6dd4f8c26063E18eb4Bd228e59c3E9"
	data := g_testData.GetCallExtraData(t, "transfer", fmt.Sprintf("%s,1000", aliceAddr))
	txHash := SendCallContract(t, cluster, nodeName, contractHash, data)
	checkTransactionOnChain(client, []common.Hash{txHash})

	// Get Balance
	from, err := cluster.GetNodeMainAddress(nodeName)
	assert.NoError(t, err)
	to := GetContractAddressByTxHash(client, contractHash)
	input := g_testData.GetCallExtraData(t, "getBalance", aliceAddr)
	err = Call(client, from, to, input)
	assert.NoError(t, err)
}
