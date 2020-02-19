package vm

import (
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/tests/g-testData"
	"github.com/dipperin/dipperin-core/tests/node-cluster"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

var (
	aliceAddr = "0x00005586B883Ec6dd4f8c26063E18eb4Bd228e59c3E9"
)

func Test_TokenConstantCall(t *testing.T) {
	cluster, err := node_cluster.CreateNodeCluster()
	assert.NoError(t, err)

	nodeName := "default_v0"
	client := cluster.NodeClient[nodeName]

	WASMTokenPath := g_testData.GetWASMPath("token-const", g_testData.CoreVmTestData)
	AbiTokenPath := g_testData.GetAbiPath("token-const", g_testData.CoreVmTestData)
	contractHash := SendCreateContract(t, cluster, nodeName, WASMTokenPath, AbiTokenPath, "dipp,DIPP,1000000")
	checkTransactionOnChain(client, []common.Hash{contractHash})

	// Transfer money
	data, err := g_testData.GetCallExtraData("transfer", fmt.Sprintf("%s,1000", aliceAddr))
	assert.NoError(t, err)

	txHash := SendCallContract(t, cluster, nodeName, contractHash, data, big.NewInt(0))
	checkTransactionOnChain(client, []common.Hash{txHash})

	// Get Balance
	from, err := cluster.GetNodeMainAddress(nodeName)
	assert.NoError(t, err)
	to := GetContractAddressByTxHash(client, contractHash)
	input, err := g_testData.GetCallExtraData("getBalance", aliceAddr)
	assert.NoError(t, err)
	resp, err := Call(client, from, to, input)
	assert.NoError(t, err)
	assert.Equal(t, "1000", resp)
}

func Test_TokenPayableCall(t *testing.T) {
	cluster, err := node_cluster.CreateNodeCluster()
	assert.NoError(t, err)

	nodeName := "default_v0"
	client := cluster.NodeClient[nodeName]

	WASMTokenPath := g_testData.GetWASMPath("token-payable", g_testData.CoreVmTestData)
	AbiTokenPath := g_testData.GetAbiPath("token-payable", g_testData.CoreVmTestData)
	contractHash := SendCreateContract(t, cluster, nodeName, WASMTokenPath, AbiTokenPath, "dipp,DIPP,1000000")
	checkTransactionOnChain(client, []common.Hash{contractHash})

	// transfer dip to contract and transfer token to alice
	data, err := g_testData.GetCallExtraData("transfer", fmt.Sprintf("%s,1000", aliceAddr))
	assert.NoError(t, err)
	txHash := SendCallContract(t, cluster, nodeName, contractHash, data, big.NewInt(500))
	checkTransactionOnChain(client, []common.Hash{txHash})

	// get contract balance
	contractAddr := GetContractAddressByTxHash(client, contractHash)
	balance := CurrentBalance(client, contractAddr)
	assert.Equal(t, uint64(500), balance.Balance.ToInt().Uint64())

	// get alice balance
	from, err := cluster.GetNodeMainAddress(nodeName)
	assert.NoError(t, err)
	input, err := g_testData.GetCallExtraData("getBalance", aliceAddr)
	assert.NoError(t, err)
	resp, err := Call(client, from, contractAddr, input)
	assert.NoError(t, err)
	assert.Equal(t, "1000", resp)

	// withdraw
	data, err = g_testData.GetCallExtraData("withdraw", "")
	assert.NoError(t, err)
	txHash = SendCallContract(t, cluster, nodeName, contractHash, data, big.NewInt(0))
	checkTransactionOnChain(client, []common.Hash{txHash})

	// get contract balance
	balance = CurrentBalance(client, contractAddr)
	assert.Equal(t, uint64(0), balance.Balance.ToInt().Uint64())

	// get alice balance
	resp, err = Call(client, from, contractAddr, input)
	assert.NoError(t, err)
	assert.Equal(t, "1000", resp)
}
