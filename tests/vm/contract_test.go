package vm

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"github.com/dipperin/dipperin-core/tests/node-cluster"
	"github.com/dipperin/dipperin-core/common"
	"time"
	"fmt"
)

func Test_WASMContractCreate(t *testing.T) {
	cluster, err := node_cluster.CreateNodeCluster()
	assert.NoError(t, err)

	nodeName := "default_v0"
	client := cluster.NodeClient[nodeName]
	txHashList := CreateContract(t, cluster, nodeName, 5)

	// 检查交易是否上链
	for i := 0; i < len(txHashList); i++ {
		for {
			result, num := Transaction(client, txHashList[i])
			if result {
				receipts := GetReceiptByTxHash(client, txHashList[i])
				fmt.Println(receipts)
				LogTestPrint("Test", "CallTransaction", "blockNum", num)
				break
			}
			time.Sleep(time.Second * 2)
		}
		time.Sleep(time.Millisecond * 100)
	}
}

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

	txHash := common.HexToHash("0x9d553401af38bbbe348947b94a7cf0881e4454307e2b092622048b336e6d0f98")
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
