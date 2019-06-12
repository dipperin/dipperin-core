package vm

import (
	"github.com/stretchr/testify/assert"
	"math/big"
	"github.com/dipperin/dipperin-core/common/consts"
	"testing"
	"github.com/dipperin/dipperin-core/tests/node-cluster"
	"github.com/dipperin/dipperin-core/common"
	"time"
	"fmt"
)

func Test_TokenContractCall(t *testing.T) {
	cluster, err := node_cluster.CreateNodeCluster()
	assert.NoError(t, err)

	nodeName := "default_v0"
	client := cluster.NodeClient[nodeName]
	txHashList := CreateTokenContract(t, cluster, nodeName, 1)

	// 检查交易是否上链
	for i := 0; i < len(txHashList); i++ {
		for {
			result, num := Transaction(client, txHashList[i])
			if result {
				receipts := GetReceiptByTxHash(client, txHashList[i])
				LogTestPrint("Test", "CallTransaction", "blockNum", num)
				fmt.Println(receipts)
				break
			}
			time.Sleep(time.Second * 2)
		}
		time.Sleep(time.Millisecond * 100)
	}

	// 根据交易ID获取合约地址
	var addrList []common.Address
	for i := 0; i < len(txHashList); i++ {
		addr := GetContractAddressByTxHash(client, txHashList[i])
		addrList = append(addrList, addr)
	}

	txHashList = CallTokenContract(t, cluster, nodeName, addrList)

	// 检查交易是否上链
	for i := 0; i < len(txHashList); i++ {
		for {
			result, num := Transaction(client, txHashList[i])
			if result {
				receipts := GetReceiptByTxHash(client, txHashList[i])
				LogTestPrint("Test", "CallTransaction", "blockNum", num)
				fmt.Println(receipts)
				break
			}
			time.Sleep(time.Second * 2)
		}
		time.Sleep(time.Millisecond * 100)
	}
}

func CreateTokenContract(t *testing.T, cluster *node_cluster.NodeCluster, nodeName string, times int) []common.Hash {
	client := cluster.NodeClient[nodeName]
	from, err := cluster.GetNodeMainAddress(nodeName)
	LogTestPrint("Test", "From", "addr", from.Hex())
	assert.NoError(t, err)

	to := common.HexToAddress(common.AddressContractCreate)
	value := big.NewInt(100)
	gasLimit := big.NewInt(2 * consts.DIP)
	gasPrice := big.NewInt(2)

	params := []string{"dipp", "DIPP", "100000000"}
	data := getCreateExtraData(t, AbiTokenPath, WASMTokenPath, params)

	var txHashList []common.Hash
	for i := 0; i < times; i++ {
		txHash, innerErr := SendTransactionContract(client, from, to, value, gasLimit, gasPrice, data)
		assert.NoError(t, innerErr)
		txHashList = append(txHashList, txHash)
	}
	return txHashList
}

func CallTokenContract(t *testing.T, cluster *node_cluster.NodeCluster, nodeName string, addrList []common.Address) []common.Hash {
	client := cluster.NodeClient[nodeName]
	from, err := cluster.GetNodeMainAddress(nodeName)
	LogTestPrint("Test", "From", "addr", from.Hex())
	assert.NoError(t, err)

	value := big.NewInt(100)
	gasLimit := big.NewInt(2 * consts.DIP)
	gasPrice := big.NewInt(1)

	var txHashList []common.Hash
	for i := 0; i < len(addrList); i++ {
		input := getCallExtraData(t, "transfer", "0x00005586B883Ec6dd4f8c26063E18eb4Bd228e59c3E9,10000")
		txHash, innerErr := SendTransactionContract(client, from, addrList[i], value, gasLimit, gasPrice, input)
		assert.NoError(t, innerErr)
		txHashList = append(txHashList, txHash)
	}
	return txHashList
}
