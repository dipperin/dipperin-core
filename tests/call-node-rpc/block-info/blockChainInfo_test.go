package block_info

import (
	"github.com/dipperin/dipperin-core/common/log"
	"github.com/dipperin/dipperin-core/core/rpc-interface"
	"github.com/dipperin/dipperin-core/tests/node-cluster"
	"github.com/dipperin/dipperin-core/tests/vm"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_GetCurrentBlock(t *testing.T) {
	t.Skip()
	cluster, err := node_cluster.CreateNodeCluster()
	assert.NoError(t, err)

	nodeName := "default_n0"
	client := cluster.NodeClient[nodeName]
	var respBlock rpc_interface.BlockResp
	err = client.Call(&respBlock, vm.GetRpcTXMethod("CurrentBlock"))
	assert.NoError(t, err)

	log.DLogger.Info("the current Block is:", "blockNumber", respBlock.Header.Number)
}

func Test_GetBlockByNumber(t *testing.T) {
	t.Skip()
	cluster, err := node_cluster.CreateNodeCluster()
	assert.NoError(t, err)

	assert.NoError(t, err)

	nodeName := "default_m0"
	client := cluster.NodeClient[nodeName]
	var respBlock rpc_interface.BlockResp
	err = client.Call(&respBlock, vm.GetRpcTXMethod("GetBlockByNumber"), 660)
	assert.NoError(t, err)
	log.DLogger.Info("the respBlock txs numbers is:", "txNumber", len(respBlock.Body.Txs))
}
