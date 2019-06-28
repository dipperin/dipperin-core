package vm

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"github.com/dipperin/dipperin-core/tests/node-cluster"
	"github.com/dipperin/dipperin-core/common"
	"fmt"
	"math/big"
	"github.com/dipperin/dipperin-core/common/consts"
)

var (
	value = big.NewInt(100)
	gasLimit = big.NewInt(2 * consts.DIP)
	gasPrice = big.NewInt(1)
)

func TestGetReceiptsByBlockNum(t *testing.T) {
	cluster, err := node_cluster.CreateNodeCluster()
	assert.NoError(t, err)

	nodeName := "default_v0"
	client := cluster.NodeClient[nodeName]

	receipts := GetReceiptsByBlockNum(client, 2589)
	fmt.Println(receipts)
}

func TestGetContractAddressByTxHash(t *testing.T) {
	cluster, err := node_cluster.CreateNodeCluster()
	assert.NoError(t, err)

	nodeName := "default_v0"
	client := cluster.NodeClient[nodeName]

	txHash := common.HexToHash("0x380699fc9660dffaa6d33e95194f72cc0cffdfe7bef9e3397fccd3fe182985d8")
	contractAddr := GetContractAddressByTxHash(client, txHash)
	receipt := GetReceiptByTxHash(client, txHash)
	assert.Equal(t, receipt.ContractAddress, contractAddr)
}

func TestGetBlockByNumber(t *testing.T) {
	cluster, err := node_cluster.CreateNodeCluster()
	assert.NoError(t, err)

	nodeName := "default_v0"
	client := cluster.NodeClient[nodeName]

	block := GetBlockByNumber(client, 992)
	fmt.Println(block.Header.String())
}
