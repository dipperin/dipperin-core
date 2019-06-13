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

func Test_WASMContractCall(t *testing.T) {
	cluster, err := node_cluster.CreateNodeCluster()
	assert.NoError(t, err)

	nodeName := "default_v0"
	client := cluster.NodeClient[nodeName]
	txHashList := CreateContract(t, cluster, nodeName, 5)
	checkTransactionOnChain(client, txHashList)

	// 根据交易ID获取合约地址
	var addrList []common.Address
	for i := 0; i < len(txHashList); i++ {
		addr := GetContractAddressByTxHash(client, txHashList[i])
		addrList = append(addrList, addr)
	}

	txHashList = CallContract(t, cluster, nodeName, addrList)
	checkTransactionOnChain(client, txHashList)
}

func CreateContract(t *testing.T, cluster *node_cluster.NodeCluster, nodeName string, times int) []common.Hash {
	client := cluster.NodeClient[nodeName]
	from, err := cluster.GetNodeMainAddress(nodeName)
	LogTestPrint("Test", "From", "addr", from.Hex())
	assert.NoError(t, err)

	to := common.HexToAddress(common.AddressContractCreate)
	value := big.NewInt(100)
	gasLimit := big.NewInt(2 * consts.DIP)
	gasPrice := big.NewInt(2)
	//txFee := big.NewInt(0).Mul(gasLimit, gasPrice)

	data := getCreateExtraData(t, AbiPath, WASMPath, nil)
	var txHashList []common.Hash
	for i := 0; i < times; i++ {
		txHash, innerErr := SendTransactionContract(client, from, to, value, gasLimit, gasPrice, data)
		assert.NoError(t, innerErr)
		txHashList = append(txHashList, txHash)

		/*		txHash, innerErr = SendTransaction(client, from, factory.AliceAddrV, value, txFee, nil)
				assert.NoError(t, innerErr)
				txHashList = append(txHashList, txHash)*/
	}
	return txHashList
}

func CallContract(t *testing.T, cluster *node_cluster.NodeCluster, nodeName string, addrList []common.Address) []common.Hash {
	client := cluster.NodeClient[nodeName]
	from, err := cluster.GetNodeMainAddress(nodeName)
	LogTestPrint("Test", "From", "addr", from.Hex())
	assert.NoError(t, err)

	value := big.NewInt(100)
	gasLimit := big.NewInt(2 * consts.DIP)
	gasPrice := big.NewInt(1)

	var txHashList []common.Hash
	for i := 0; i < len(addrList); i++ {
		input := getCallExtraData(t, "hello", fmt.Sprintf("Event,%v", 100*i))
		txHash, innerErr := SendTransactionContract(client, from, addrList[i], value, gasLimit, gasPrice, input)
		assert.NoError(t, innerErr)
		txHashList = append(txHashList, txHash)
	}
	return txHashList
}
