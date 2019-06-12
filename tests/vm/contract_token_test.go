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

	params := []string{"dipp", "DIPP", "1000000"}
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
		input := getCallExtraData(t, "transfer", "0x00004179D57e45Cb3b54D6FAEF69e746bf240E287978,1000")
		txHash, innerErr := SendTransactionContract(client, from, addrList[i], value, gasLimit, gasPrice, input)
		assert.NoError(t, innerErr)
		txHashList = append(txHashList, txHash)
	}
	return txHashList
}

func TestGetReceiptByTxHash(t *testing.T) {
	//input := []byte{0, 225, 245, 5, 0, 0, 0, 0}
	//input := []byte{16, 39, 0, 0, 0, 0, 0, 0}
	input := []byte{1, 44, 48, 48, 48, 48, 54, 50, 98, 101, 49, 48, 102, 52, 54, 98, 53, 100, 48, 49, 101, 99, 100, 57, 98, 53, 48, 50, 99, 52, 98, 97, 51, 100, 54, 49, 51, 49, 102, 54, 102, 99, 50, 101, 52, 49, 16, 39, 0, 0, 0, 0, 0, 0,}
	fmt.Println(string(input))

	addr := common.HexToAddress("0x00004179D57e45Cb3b54D6FAEF69e746bf240E287978")
	fmt.Println(addr[:])
}