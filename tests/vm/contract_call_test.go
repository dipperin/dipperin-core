package vm

import (
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/consts"
	"github.com/dipperin/dipperin-core/core/vm/model"
	"github.com/dipperin/dipperin-core/tests/node-cluster"
	"github.com/dipperin/dipperin-core/third-party/rpc"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func Test_WASMContractCall(t *testing.T) {
	cluster, err := node_cluster.CreateNodeCluster()
	assert.NoError(t, err)

	nodeName := "default_v0"
	client := cluster.NodeClient[nodeName]
	txHashList := CreateContract(t, cluster, nodeName, 1)
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

func Test_GetReceipt(t *testing.T){
	cluster, err := node_cluster.CreateNodeCluster()
	assert.NoError(t, err)

	nodeName := "default_v0"
	client := cluster.NodeClient[nodeName]
	txId:= common.HexToHash("0x9ef3eb906c570de048a463ca168a299c67261d00bf4152a9452f91a0ba907dcc")
	//get receipt
	receipt:=GetContractReceipt(t,client,txId)
	fmt.Print("the receipt is:\r\n",receipt.String())
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

	data := getCreateExtraData(t, WASMPath, AbiPath, "")
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
	gasPrice := big.NewInt(2)

	var txHashList []common.Hash
	for i := 0; i < len(addrList); i++ {
		input := getCallExtraData(t, "hello", fmt.Sprintf("Event,%v", 100*i))
		txHash, innerErr := SendTransactionContract(client, from, addrList[i], value, gasLimit, gasPrice, input)
		assert.NoError(t, innerErr)
		txHashList = append(txHashList, txHash)
	}
	return txHashList
}

func GetContractReceipt(t *testing.T,client *rpc.Client,txId common.Hash) *model.Receipt{
	receipt := model.Receipt{}
	err := client.Call(&receipt, GetRpcTXMethod("GetConvertReceiptByTxHash"), txId)
	assert.NoError(t,err)

	return &receipt
}
