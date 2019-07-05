package send_tx

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/tests/node-cluster"
	"github.com/dipperin/dipperin-core/tests/vm"
	"github.com/dipperin/dipperin-core/third-party/crypto"
	"github.com/dipperin/dipperin-core/third-party/crypto/cs-crypto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_SendNormalTx(t *testing.T){
	t.Skip()
	cluster, err := node_cluster.CreateNodeCluster()
	assert.NoError(t, err)

	nodeName := "default_v0"
	client := cluster.NodeClient[nodeName]

	nodeAddr ,err:= cluster.GetNodeMainAddress(nodeName)
	assert.NoError(t,err)
	testKey,err := crypto.GenerateKey()
	assert.NoError(t,err)
	toAddr := cs_crypto.GetNormalAddress(testKey.PublicKey)

	var resp common.Hash
	err = client.Call(&resp, vm.GetRpcTXMethod("SendTransaction"), nodeAddr, toAddr, 1000, 1, 21000, nil, nil)
	assert.NoError(t,err)
}
