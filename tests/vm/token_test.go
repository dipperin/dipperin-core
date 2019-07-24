package vm

import (
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/tests/g-testData"
	"github.com/dipperin/dipperin-core/tests/node-cluster"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_TokenContractCall(t *testing.T) {
	cluster, err := node_cluster.CreateNodeCluster()
	assert.NoError(t, err)

	nodeName := "default_v0"
	client := cluster.NodeClient[nodeName]

	WASMTokenPath := g_testData.GetWASMPath("token-const",g_testData.CoreVmTestData)
	AbiTokenPath := g_testData.GetAbiPath("token-const",g_testData.CoreVmTestData)
	contractHash := SendCreateContract(t, cluster, nodeName, WASMTokenPath, AbiTokenPath, "dipp,DIPP,1000000")
	checkTransactionOnChain(client, []common.Hash{contractHash})

	// Transfer money
	aliceAddr := "0x00005586B883Ec6dd4f8c26063E18eb4Bd228e59c3E9"
	data, err := g_testData.GetCallExtraData("transfer", fmt.Sprintf("%s,1000", aliceAddr))
	assert.NoError(t, err)

	txHash := SendCallContract(t, cluster, nodeName, contractHash, data)
	checkTransactionOnChain(client, []common.Hash{txHash})

	// Get Balance
	from, err := cluster.GetNodeMainAddress(nodeName)
	assert.NoError(t, err)
	to := GetContractAddressByTxHash(client, contractHash)
	input, err := g_testData.GetCallExtraData("getBalance", aliceAddr)
	assert.NoError(t, err)

	err = Call(client, from, to, input)
	assert.NoError(t, err)
}

/*func Test_dial_server(t *testing.T){
	//if client, err = rpc.Dial(fmt.Sprintf("http://%v:%d", "127.0.0.1", port)); err != nil {
	//	panic("init rpc client failed: " + err.Error())
	//}
	wsURL := fmt.Sprintf("ws://%v:%d", "172.16.5.183", 7002)
	//l.Info("init rpc client", "wsURL", wsURL)
	if _, err := rpc.Dial(wsURL); err != nil {
		panic("init rpc client failed: " + err.Error())
	}

	log.Info("dial success")
}*/
