package send_tx

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/consts"
	"github.com/dipperin/dipperin-core/tests/node-cluster"
	"github.com/dipperin/dipperin-core/tests/vm"
	"github.com/dipperin/dipperin-core/third-party/crypto"
	"github.com/dipperin/dipperin-core/third-party/crypto/cs-crypto"
	"github.com/stretchr/testify/assert"
	"testing"
	"fmt"
	"github.com/dipperin/dipperin-core/third-party/rpc"
	"github.com/dipperin/dipperin-core/tests/g-testData"
	"math/big"
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
	cluster, err := node_cluster.CreateNodeCluster()
	assert.NoError(t, err)

	// get verifier info
	number := 3
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
		_, err = vm.SendTransaction(verifierClient[0], verifierAddr[0], verifierAddr[i], big.NewInt(10000), g_testData.TestGasPrice, g_testData.TestGasLimit, nil)
		assert.NoError(t, err)
	}

	/*time.Sleep(3*time.Second)

	// verifiers send register tx
	for i := 0; i < number; i++ {
		err = verifierClient[i].Call(&resp, vm.GetRpcTXMethod("SendRegisterTransaction"), verifierAddr[i], 1000, g_testData.TestGasPrice, g_testData.TestGasLimit, nil)
		assert.NoError(t, err)
	}*/
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
