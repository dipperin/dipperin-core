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
	From, err := cluster.GetNodeMainAddress(nodeName)
	LogTestPrint("Test_WASMContractCreat", "the from address is:", "addr", From.Hex())
	assert.NoError(t, err)

	accounts, err := (nodeName)
	assert.NoError(t, err)
	for _, account := range accounts {
		LogTestPrint("Test_WASMContractCreat", "the account addr is:", "addr", account.Address.Hex())
	}

	//return

	to := common.AddressContractCreate
	Value := big.NewInt(100)
	gasLimit := 2 * consts.DIP
	gasPrice := 1

	balance, err := cluster.GetAddressBalance(nodeName, From)
	assert.NoError(t, err)
	assert.Equal(t, 1, balance.Cmp(big.NewInt(int64(gasLimit*gasPrice))))

	abiBytes, err := ioutil.ReadFile(AbiPath)
	assert.NoError(t, err)
	WASMBytes, err := ioutil.ReadFile(WASMPath)
	assert.NoError(t, err)
	ExtraData, err := rlp.EncodeToBytes([]interface{}{WASMBytes, abiBytes})
	assert.NoError(t, err)
	var resp common.Hash
	if err := client.Call(&resp, GetRpcMethod("SendTransactionContract"), From, to, Value, gasLimit, gasPrice, ExtraData, nil); err != nil {

		LogTestPrint("Test_WASMContractCreat", "call send transaction", "err", err)
		return
	}
	LogTestPrint("Test_WASMContractCreat", "the contract creat transaction id is:", "txId", resp.Hex())

	cluster.CheckTxIsOnBlockChain(nodeName, resp, ContractTxType)
}

func Test_TransactionReceipt(t *testing.T) {
	cluster, err := CreateNodeCluster()
	assert.NoError(t, err)

	nodeName := "default_v0"
	client := cluster.NodeClient[nodeName]

	txHash := "0x47d4beb428771608aa2a77b5eae05aad695b5928d1ecba8c32ee37a7fddf5111"
	var resp model.Receipts
	if err := client.Call(&resp, GetRpcMethod("TransactionReceipt"), txHash); err != nil {
		LogTestPrint("Test_TransactionReceipt", "call TransactionReceipt", "err", err)
		return
	}
	for i := 0; i < len(resp); i++ {
		if common.HexToHash(txHash).IsEqual(resp[i].TxHash) {
			fmt.Println(resp[i].String())
		}
	}
}
