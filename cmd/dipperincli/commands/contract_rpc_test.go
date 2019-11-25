package commands

import (
	"fmt"
	"github.com/dipperin/dipperin-core/tests/g-testData"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
	"os"
	"testing"
)

func TestRpcCaller_GetContractAddressByTxHash(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := getRpcTestApp()
	app.Action = func(context *cli.Context) {
		c := &rpcCaller{}
		c.GetContractAddressByTxHash(context)
	}
	assert.NoError(t, app.Run([]string{os.Args[0]}))

	app.Action = func(context *cli.Context) {
		c := &rpcCaller{}
		context.Set("p", "")
		c.GetContractAddressByTxHash(context)

		context.Set("p", "txHash")
		c.GetContractAddressByTxHash(context)

		context.Set("p", "txHash,123")
		c.GetContractAddressByTxHash(context)

		context.Set("p", txHash)
		client = NewMockRpcClient(ctrl)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		c.GetContractAddressByTxHash(context)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)
		c.GetContractAddressByTxHash(context)
	}
	assert.NoError(t, app.Run([]string{os.Args[0], "GetContractAddressByTxHash"}))
	client = nil
}

func TestRpcCaller_CallContract(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := getRpcTestApp()
	app.Action = func(context *cli.Context) {
		c := &rpcCaller{}
		c.CallContract(context)
	}
	assert.NoError(t, app.Run([]string{os.Args[0]}))

	app.Action = func(context *cli.Context) {
		c := &rpcCaller{}
		context.Set("p", "")
		c.CallContract(context)

		context.Set("p", "from")
		c.CallContract(context)

		context.Set("p", "from,to")
		c.CallContract(context)

		context.Set("p", fmt.Sprintf("%s,to", from))
		c.CallContract(context)

		context.Set("p", fmt.Sprintf("%s,%s", from, to))
		c.CallContract(context)

		context.Set("p", fmt.Sprintf("%s,%s,blockNum", from, to))
		c.CallContract(context)

		context.Set("p", fmt.Sprintf("%s,%s,%v", from, to, "64"))
		context.Set("func-name", "name")
		context.Set("input", "input1")
		client = NewMockRpcClient(ctrl)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		c.CallContract(context)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)
		context.Set("input", "input1,input2")
		c.CallContract(context)

		context.Set("p", "from,to,blockNum,extraData")
		c.CallContract(context)
	}
	assert.NoError(t, app.Run([]string{os.Args[0], "CallContract"}))
	client = nil
}

func TestRpcCaller_EstimateGas(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := getRpcTestApp()
	app.Action = func(context *cli.Context) {
		c := &rpcCaller{}
		c.EstimateGas(context)
	}
	assert.NoError(t, app.Run([]string{os.Args[0]}))

	app.Action = func(context *cli.Context) {
		c := &rpcCaller{}
		context.Set("p", "")
		c.EstimateGas(context)

		context.Set("p", "from,to")
		c.EstimateGas(context)

		context.Set("p", "from,to,value,gasPrice,gasLimit")
		c.EstimateGas(context)

		context.Set("p", fmt.Sprintf("%s,to,value,gasPrice,gasLimit", from))
		c.EstimateGas(context)

		context.Set("p", fmt.Sprintf("%s,%s,value,gasPrice,gasLimit", from, to))
		c.EstimateGas(context)

		context.Set("p", fmt.Sprintf("%s,%s,%v,gasPrice,gasLimit", from, to, "10dip"))
		c.EstimateGas(context)

		context.Set("p", fmt.Sprintf("%s,%s,%v,%v,gasLimit", from, to, "10dip", "1wu"))
		c.EstimateGas(context)

		context.Set("p", fmt.Sprintf("%s,%s,%v,%v,%v", from, to, "10dip", "1wu", "1000"))
		c.EstimateGas(context)

		context.Set("func-name", "name")
		context.Set("input", "input1")
		client = NewMockRpcClient(ctrl)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		c.EstimateGas(context)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)
		c.EstimateGas(context)

		context.Set("is-create", "true")
		c.EstimateGas(context)
	}
	assert.NoError(t, app.Run([]string{os.Args[0], "CallContract"}))
	client = nil
}

func TestRpcCaller_SendTransactionContract(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := getRpcTestApp()
	app.Action = func(context *cli.Context) {
		c := &rpcCaller{}
		SyncStatus.Store(true)
		context.Set("is-create", "true")
		c.SendTransactionContract(context)
	}
	assert.NoError(t, app.Run([]string{os.Args[0]}))

	wasm := g_testData.GetWASMPath("event", g_testData.CoreVmTestData)
	abi := g_testData.GetAbiPath("event", g_testData.CoreVmTestData)
	app.Action = func(context *cli.Context) {
		client = NewMockRpcClient(ctrl)
		c := &rpcCaller{}
		SyncStatus.Store(false)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		c.SendTransactionContract(context)

		SyncStatus.Store(true)
		c.SendTransactionContract(context)

		context.Set("is-create", "true")
		c.SendTransactionContract(context)

		context.Set("p", "from,value,gasPrice,gasLimit")
		c.SendTransactionContract(context)

		context.Set("p", fmt.Sprintf("%s,value,gasPrice,gasLimit", from))
		c.SendTransactionContract(context)

		context.Set("p", fmt.Sprintf("%s,%v,gasPrice,gasLimit", from, "10dip"))
		c.SendTransactionContract(context)

		context.Set("p", fmt.Sprintf("%s,%v,%v,gasLimit", from, "10dip", "1wu"))
		c.SendTransactionContract(context)

		context.Set("p", fmt.Sprintf("%s,%v,%v,%v", from, "10dip", "1wu", "1000"))
		c.SendTransactionContract(context)

		context.Set("wasm", wasm)
		c.SendTransactionContract(context)

		context.Set("abi", abi)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		c.SendTransactionContract(context)
	}

	app.Run([]string{"xxx", "SendTransactionContract"})
	client = nil
}

func TestRlpBool(t *testing.T)  {
	exist := true
	notExist := false
	rlpParam := []interface{}{
		exist, notExist,
	}

	rlp.EncodeToBytes(rlpParam)
}