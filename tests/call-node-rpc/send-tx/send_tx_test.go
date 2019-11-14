package send_tx

import (
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/consts"
	"github.com/dipperin/dipperin-core/tests/g-testData"
	"github.com/dipperin/dipperin-core/tests/node-cluster"
	"github.com/dipperin/dipperin-core/tests/vm"
	"github.com/dipperin/dipperin-core/third-party/crypto"
	"github.com/dipperin/dipperin-core/third-party/crypto/cs-crypto"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/dipperin/dipperin-core/third-party/rpc"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
	"time"
)

func Test_SendNormalTx(t *testing.T) {
	t.Skip()
	cluster, err := node_cluster.CreateNodeCluster()
	assert.NoError(t, err)

	nodeName := "default_v0"
	client := cluster.NodeClient[nodeName]

	nodeAddr, err := cluster.GetNodeMainAddress(nodeName)
	assert.NoError(t, err)
	testKey, err := crypto.GenerateKey()
	assert.NoError(t, err)
	toAddr := cs_crypto.GetNormalAddress(testKey.PublicKey)

	var resp common.Hash
	err = client.Call(&resp, vm.GetRpcTXMethod("SendTransaction"), nodeAddr, toAddr, 1000, 1, 21000, nil, nil)
	assert.NoError(t, err)
}

func Test_SendRegisterTx(t *testing.T) {
	cluster, err := node_cluster.CreateIpcNodeCluster()
	assert.NoError(t, err)

	// get verifier info
	number := 20
	var verifierAddr []common.Address
	var verifierClient []*rpc.Client
	for i := 0; i < number; i++ {
		verifierName := fmt.Sprintf("default_v%v", i)
		client := cluster.NodeClient[verifierName]
		verifier, innerErr := cluster.GetNodeMainAddress(verifierName)
		assert.NoError(t, innerErr)

		verifierAddr = append(verifierAddr, verifier)
		verifierClient = append(verifierClient, client)
	}

	// verifiers[0] send money to others
	for i := 0; i < number; i++ {
		_, err = vm.SendTransaction(verifierClient[0], verifierAddr[0], verifierAddr[i], big.NewInt(1000000), g_testData.TestGasPrice, g_testData.TestGasLimit, nil)
		assert.NoError(t, err)
	}

	time.Sleep(3 * time.Second)

	// verifiers send register tx
	for i := 0; i < number; i++ {
		_, err = vm.SendRegisterTransaction(verifierClient[i], verifierAddr[i], big.NewInt(10000), g_testData.TestGasPrice, g_testData.TestGasLimit)
		assert.NoError(t, err)
	}
}

func Test_CurrentBalance(t *testing.T) {
	cluster, err := node_cluster.CreateIpcNodeCluster()
	assert.NoError(t, err)

	number := 20
	var (
		verifierAddr   []common.Address
		verifierClient []*rpc.Client
		verifierName   []string
	)

	for i := 0; i < number; i++ {
		name := fmt.Sprintf("default_v%v", i)
		client := cluster.NodeClient[name]
		verifier, innerErr := cluster.GetNodeMainAddress(name)
		assert.NoError(t, innerErr)

		verifierAddr = append(verifierAddr, verifier)
		verifierClient = append(verifierClient, client)
		verifierName = append(verifierName, name)
	}

	// current balance
	for i := 0; i < number; i++ {
		resp := vm.CurrentBalance(verifierClient[i], verifierAddr[i])
		log.Info("balance", "node", verifierName[i], "balance", resp.Balance.ToInt(), "address", verifierAddr[i])
	}
}

func Test_SendMoneyFromV0(t *testing.T) {
	cluster, err := node_cluster.CreateNodeCluster()
	assert.NoError(t, err)

	nodeName := "default_v0"
	client := cluster.NodeClient[nodeName]

	nodeAddr, err := cluster.GetNodeMainAddress(nodeName)
	assert.NoError(t, err)

	toAddr := common.HexToAddress("0x000075023A165a587a9fBc81E9D65830338348141A44")
	var resp common.Hash
	err = client.Call(&resp, vm.GetRpcTXMethod("SendTransaction"), nodeAddr, toAddr, 10*consts.DIP, 1, 21000, nil, nil)
	assert.NoError(t, err)
}
