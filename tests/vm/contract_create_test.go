package vm

import (
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"math/big"
	"github.com/dipperin/dipperin-core/common/consts"
	"io/ioutil"
	"github.com/dipperin/dipperin-core/common/util"
	"testing"
	"github.com/dipperin/dipperin-core/tests/node-cluster"
	"github.com/dipperin/dipperin-core/common"
	"github.com/ethereum/go-ethereum/rlp"
	"time"
	"fmt"
	"github.com/dipperin/dipperin-core/tests/factory"
)

var (
	AbiPath  = filepath.Join(util.HomeDir(), "go/src/github.com/dipperin/dipperin-core/core/vm/event/event.cpp.abi.json")
	WASMPath = filepath.Join(util.HomeDir(), "go/src/github.com/dipperin/dipperin-core/core/vm/event/event.wasm")
)

func Test_WASMContractCreate(t *testing.T) {
	cluster, err := node_cluster.CreateNodeCluster()
	assert.NoError(t, err)

	nodeName := "default_v0"
	client := cluster.NodeClient[nodeName]
	from, err := cluster.GetNodeMainAddress(nodeName)
	LogTestPrint("Test", "From", "addr", from.Hex())
	assert.NoError(t, err)

	to := common.HexToAddress(common.AddressContractCreate)
	Value := big.NewInt(100)
	gasLimit := big.NewInt(2 * consts.DIP)
	gasPrice := big.NewInt(5)
	txFee := big.NewInt(0).Mul(gasLimit, gasPrice)

	balance, err := cluster.GetAddressBalance(nodeName, from)
	assert.NoError(t, err)
	assert.Equal(t, 1, balance.Cmp(txFee))

	abiBytes, err := ioutil.ReadFile(AbiPath)
	assert.NoError(t, err)
	WASMBytes, err := ioutil.ReadFile(WASMPath)
	assert.NoError(t, err)
	ExtraData, err := rlp.EncodeToBytes([]interface{}{WASMBytes, abiBytes})
	assert.NoError(t, err)

	var txHashList []common.Hash
	for i := 0; i < 1; i++ {
		txHash, innerErr := SendTransactionContract(client, from, to, Value, gasLimit, gasPrice, ExtraData)
		assert.NoError(t, innerErr)
		txHashList = append(txHashList, txHash)

		txHash, innerErr = SendTransaction(client, from, factory.AliceAddrV, Value, txFee, nil)
		assert.NoError(t, innerErr)
		txHashList = append(txHashList, txHash)
	}

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

	receipts := GetReceiptsByBlockNum(client, 75)
	fmt.Println(receipts)
}

func TestGetContractAddressByTxHash(t *testing.T) {
	cluster, err := node_cluster.CreateNodeCluster()
	assert.NoError(t, err)

	nodeName := "default_v0"
	client := cluster.NodeClient[nodeName]

	txHash := common.HexToHash("0x7a11365d934a0e583a26a70773b1fb46fb0d2c7307f802e10d6348a2ae78d0ac")
	contractAddr := GetContractAddressByTxHash(client, txHash)
	receipt := GetReceiptByTxHash(client, txHash)
	assert.Equal(t, receipt.ContractAddress, contractAddr)
}

func TestGetBlockByNumber(t *testing.T) {
	cluster, err := node_cluster.CreateNodeCluster()
	assert.NoError(t, err)

	nodeName := "default_v0"
	client := cluster.NodeClient[nodeName]

	block := GetBlockByNumber(client, 75)
	fmt.Println(block.Header.String())
}
