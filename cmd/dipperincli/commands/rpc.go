// Copyright 2019, Keychain Foundation Ltd.
// This file is part of the dipperin-core library.
//
// The dipperin-core library is free software: you can redistribute
// it and/or modify it under the terms of the GNU Lesser General Public License
// as published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// The dipperin-core library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package commands

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/consts"
	"github.com/dipperin/dipperin-core/common/hexutil"
	"github.com/dipperin/dipperin-core/common/log"
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/core/accounts/accountsbase"
	"github.com/dipperin/dipperin-core/core/accounts/softwallet"
	"github.com/dipperin/dipperin-core/core/chain"
	"github.com/dipperin/dipperin-core/core/chainconfig"
	"github.com/dipperin/dipperin-core/core/dipperin"
	model2 "github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/core/rpcinterface"
	"github.com/dipperin/dipperin-core/third_party/rpc"
	"github.com/urfave/cli"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

var (
	l *zap.Logger

	client RpcClient

	defaultAccount common.Address
	defaultWallet  accountsbase.WalletIdentifier
	osExit         = os.Exit
)

const (
	VerifierStatusNoRegistered = "Not Registered"
	VerifierStatusRegistered   = "Registered"
	VerifiedStatusCanceled     = "Canceled"
	VerifiedStatusUnstaked     = "Unstaked"
)

//go:generate mockgen -destination=./rpc_client_mock_test.go -package=commands github.com/dipperin/dipperin-core/cmd/dipperincli/commands RpcClient
type RpcClient interface {
	Call(result interface{}, method string, args ...interface{}) error
	Subscribe(ctx context.Context, namespace string, channel interface{}, args ...interface{}) (*rpc.ClientSubscription, error)
}

func init() {
	l = log.NewLogger()
}

func InitRpcClient(info dipperin.NodeInfo) {
	l.Info("init inProc client")
	client = rpc.DialInProc(info.InProcHandler)
}

func InitAccountInfo(nodeType int, path, password, passPharse string) {
	// get default account
	if nodeType == chainconfig.NodeTypeOfNormal {
		err := initWallet(path, password, passPharse)
		if err != nil {
			l.Error("init Wallet Error", zap.Error(err))
			osExit(1)
		}
	}

	defaultAccount = getDefaultAccount()
	defaultWallet = getDefaultWallet()
}

// CsSubscribe registers a subscripion under the "cs" namespace.
//func csSubscribe(c RpcClient, ctx context.Context, channel interface{}, args ...interface{}) (*rpc.ClientSubscription, error) {
//	return c.Subscribe(ctx, "dipperin", channel, args...)
//}
//
//func Subscribe(nodeType string) {
//
//	if nodeType != "mine master" {
//		return
//	}
//
//	blockCh := make(chan string)
//
//	sub, err := csSubscribe(client, context.Background(), blockCh, "newBlock")
//
//	if err != nil {
//		l.Error("can't subscribe server", zap.Error(err))
//		return
//	}
//
//	l.Info("start subscribe node server msg")
//
//	for {
//		select {
//		case err := <-sub.Err():
//			l.Info("sub result", zap.Error(err))
//		case <-blockCh:
//			// todo Make another prompt, otherwise it is difficult to operate the command.
//			//l.Info(b)
//		}
//	}
//
//}

var (
	CheckSyncStatusDuration = time.Second * 3
)

var SyncStatus atomic.Value

//check downloader Sync Status
func CheckDownloaderSyncStatus() {
	timer := time.NewTimer(CheckSyncStatusDuration)
	SyncStatus.Store(false)

	var status bool
	if err := client.Call(&status, getDipperinRpcMethodByName("GetSyncStatus")); err != nil {
		l.Error("call GetSyncStatus", zap.Error(err))
		panic("call GetSyncStatus err")
	}

	//l.Info("the status is:","status",status)
	if !status {
		SyncStatus.Store(true)
		return
	}

	var resp bool
	for {
		select {
		case <-timer.C:
			if err := client.Call(&resp, getDipperinRpcMethodByName("GetSyncStatus")); err != nil {
				l.Error("call GetSyncStatus", zap.Error(err))
			}

			//l.Info("the GetSyncStatus resp is:","resp",resp)
			if !resp {
				SyncStatus.Store(true)
				return
			}

			timer.Reset(CheckSyncStatusDuration)
		}
	}
}

type rpcCaller struct{}

var callerRv = reflect.ValueOf(&rpcCaller{})

func RpcCall(c *cli.Context) {
	if client == nil {
		panic("rpc client not initialized")
	}
	// when use method := c.Args()[0],the command line `tx SendTransactionContract -p xxxx --abi` lead the node stop
	method := c.Args().First()
	if len(c.Args()) == 0 {
		l.Error("RpcCall params assign err, can't find the method")
		return
	}

	rvf := callerRv.MethodByName(method)
	if rvf.Kind() != reflect.Func {
		l.Error("not found method", zap.String("method_name", method))
		return
	}

	// call method
	rvf.Call([]reflect.Value{reflect.ValueOf(c)})
}

// get rpc method from name
func getDipperinRpcMethodByName(mName string) string {
	lm := strings.ToLower(string(mName[0])) + mName[1:]
	return "dipperin_" + lm
}

// get rpc parameters from map, return string list, delete space at end
func getRpcParamFromString(cParam string) []string {
	if cParam == "" {
		return []string{}
	}

	lm := strings.Split(cParam, ",")
	//l.Info("the lm is:", "lm", lm)
	return lm
}

func getRpcMethodAndParam(c *cli.Context) (mName string, cParams []string, err error) {
	mName = c.Args().First()
	//l.Info("the method name is:", "mName", mName)
	if mName == "" {
		return "", []string{}, errors.New("the method name is nil")
	}
	params := c.String("p")
	//l.Info("the params is:", "params", params)

	cParams = getRpcParamFromString(params)
	return mName, cParams, nil
}

func getRpcParamValue(c *cli.Context, paramName string) (path string, err error) {
	path = c.String(paramName)
	if path == "" {
		return "", errors.New("the " + paramName + " path is nil")
	}
	return
}

func getRpcSpecialParam(c *cli.Context, paramName string) (value string) {
	return c.String(paramName)
}

func checkSync() bool {
	if !SyncStatus.Load().(bool) {
		l.Error("the block downloader isn't finished")
		var respBlock rpcinterface.BlockResp
		if err := client.Call(&respBlock, getDipperinRpcMethodByName("CurrentBlock")); err != nil {
			l.Error("get current Block error", zap.Error(err))
			return true
		}
		l.Info("the current block number is:", zap.Uint64("number", respBlock.Header.Number))
		return true
	}
	return false
}

func (caller *rpcCaller) GetDefaultAccountBalance(c *cli.Context) {
	var resp rpcinterface.CurBalanceResp
	if err := client.Call(&resp, getDipperinRpcMethodByName("CurrentBalance"), defaultAccount); err != nil {
		l.Error("call current balance error", zap.Error(err))
		return
	}

	l.Debug("GetDefaultAccountBalance", zap.Any("resp", resp))
	balance, err := CSCoinToMoneyValue(resp.Balance)
	if err != nil {
		l.Error("the address isn't on the block chain, balance=0")
	} else {
		l.Info("address current Balance is:", zap.String("balance", balance))
	}
}

func (caller *rpcCaller) CurrentBalance(c *cli.Context) {
	mName, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error")
		return
	}

	if len(cParams) != 0 && len(cParams) != 1 {
		l.Error("parameter error")
		return
	}

	var addr common.Address
	if len(cParams) == 0 {
		addr = getDefaultAccount()
	} else {
		addr, err = CheckAndChangeHexToAddress(cParams[0])
		if err != nil {
			l.Error("the input address is invalid", zap.Error(err))
			return
		}
	}

	var resp rpcinterface.CurBalanceResp
	if err = client.Call(&resp, getDipperinRpcMethodByName(mName), addr); err != nil {

		l.Error("call current balance error", zap.Error(err))
		return
	}

	balance, err := CSCoinToMoneyValue(resp.Balance)
	if err != nil {
		l.Error("the address isn't on the block chain balance=0")
	} else {
		l.Info("address current Balance is:", zap.String("balance", balance))
	}
}

func printBlockInfo(respBlock rpcinterface.BlockResp) {
	fmt.Printf("Block info:\r\n")
	fmt.Printf(respBlock.Header.String())
	//l.Info("the current Block header is:","header",respBlock.Header.String())
	fmt.Printf("\r\n[the current Block txs is]:")
	for _, tx := range respBlock.Body.Txs {
		fmt.Printf("\r\ntxID:%v", tx.CalTxId().Hex())
	}

	fmt.Printf("\r\n[the current Block commit address is]:")
	for _, ver := range respBlock.Body.Vers {
		fmt.Printf("\r\n[commit address]:%v", ver.GetAddress().Hex())
	}

	fmt.Printf("\r\n")
}

func printTransactionInfo(respTx rpcinterface.TransactionResp) {
	if respTx.Transaction == nil {
		fmt.Printf("the tx isn't on the block chain\r\n")
		return
	}

	fmt.Printf("\r\n[the tx info is:]")
	fmt.Printf("\r\n%v", respTx.Transaction.String())
	fmt.Printf("\r\nthe BlockHash is:%v", respTx.BlockHash.Hex())
	fmt.Printf("\r\nthe BlockNumber is:%v", respTx.BlockNumber)
	fmt.Printf("\r\nthe TxIndex is:%v\r\n", respTx.TxIndex)
}

func (caller *rpcCaller) CurrentBlock(c *cli.Context) {
	mName, _, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error")
		return
	}
	var respBlock rpcinterface.BlockResp
	if err = client.Call(&respBlock, getDipperinRpcMethodByName(mName)); err != nil {
		l.Error("call current block error", zap.Error(err))
		return
	}
	printBlockInfo(respBlock)
}

func (caller *rpcCaller) GetGenesis(c *cli.Context) {
	mName, _, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error")
		return
	}

	var respBlock rpcinterface.BlockResp
	if err = client.Call(&respBlock, getDipperinRpcMethodByName(mName)); err != nil {

		l.Error("call genesis block error", zap.Error(err))
		return
	}
	printBlockInfo(respBlock)
}

// GetBlockByNumber get block information according to block num
func (caller *rpcCaller) GetBlockByNumber(c *cli.Context) {
	mName, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error")
		return
	}

	if len(cParams) == 0 {
		l.Error("parameter includes：blockNumber")
		return
	}

	blockNum, err := strconv.Atoi(cParams[0])
	if err != nil {
		l.Error("the blockNumber error")
	}
	l.Debug("the blockNum is:", zap.Int("blockNum", blockNum))

	var respBlock rpcinterface.BlockResp
	if err = client.Call(&respBlock, getDipperinRpcMethodByName(mName), blockNum); err != nil {
		l.Error("call block error", zap.Error(err))
		return
	}
	printBlockInfo(respBlock)
}

// GetSlotByNumber get block information according to block num
func (caller *rpcCaller) GetSlotByNumber(c *cli.Context) {
	mName, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error")
		return
	}

	if len(cParams) == 0 {
		l.Error("parameter includes：blockNumber")
		return
	}

	blockNum, err := strconv.Atoi(cParams[0])
	if err != nil {
		l.Error("the blockNumber error")
	}
	l.Debug("the blockNum is:", zap.Int("blockNum", blockNum))

	var resp uint64
	if err = client.Call(&resp, getDipperinRpcMethodByName(mName), blockNum); err != nil {
		l.Error("call GetSlotByNumber error", zap.Error(err))
		return
	}
	l.Info("GetSlotByNumber result", zap.Uint64("slot", resp))
}

// GetBlockByHash get block information based on block hash
func (caller *rpcCaller) GetBlockByHash(c *cli.Context) {
	mName, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error")
		return
	}

	if len(cParams) == 0 {
		l.Error("parameter includes：BlockHash")
		return
	}

	var respBlock rpcinterface.BlockResp
	if err = client.Call(&respBlock, getDipperinRpcMethodByName(mName), cParams[0]); err != nil {
		l.Error("call block error", zap.Error(err))
		return
	}
	printBlockInfo(respBlock)
}

func (caller *rpcCaller) StartMine(c *cli.Context) {
	if checkSync() {
		return
	}

	mName, _, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error")
		return
	}
	var resp interface{}
	if err = client.Call(&resp, getDipperinRpcMethodByName(mName)); err != nil {
		l.Error("start mining error", zap.Error(err))
		return
	}
	l.Debug("Mining Started")
}

func (caller *rpcCaller) StopMine(c *cli.Context) {
	mName, _, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error")
		return
	}

	var resp interface{}
	if err = client.Call(&resp, getDipperinRpcMethodByName(mName)); err != nil {
		l.Error("stop mining error", zap.Error(err))
		return
	}
	l.Debug("stop mining")
}

// SetMineCoinBase set the mining coinbase address
func (caller *rpcCaller) SetMineCoinBase(c *cli.Context) {
	mName, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error")
		return
	}

	if len(cParams) < 1 {
		l.Error("setting CoinBase needs：address")
		return
	}

	address, err := CheckAndChangeHexToAddress(cParams[0])
	if err != nil {
		l.Error("the input address is invalid", zap.Error(err))
		return
	}
	var resp interface{}
	if err = client.Call(&resp, getDipperinRpcMethodByName(mName), address); err != nil {
		l.Error("setting CoinBase", zap.Error(err))
		return
	}
	l.Debug("setting CoinBase　success")
}

// SetMinerGasConfig set gasFloor and gasCeil
func (caller *rpcCaller) SetMineGasConfig(c *cli.Context) {
	mName, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error")
		return
	}

	if len(cParams) < 2 {
		l.Error("parameter includes：gasFloor, gasCeil")
		return
	}

	gasFloor, err := strconv.Atoi(cParams[0])
	if err != nil {
		l.Error("parse gasFloor error", zap.Error(err))
		return
	}

	gasCeil, err := strconv.Atoi(cParams[1])
	if err != nil {
		l.Error("parse gasCeil error", zap.Error(err))
		return
	}

	l.Debug("the gasFloor is:", zap.Int("gasFloor", gasFloor))
	l.Debug("the gasCeil is:", zap.Int("gasCeil", gasCeil))

	var resp interface{}
	if err = client.Call(&resp, getDipperinRpcMethodByName(mName), uint64(gasFloor), uint64(gasCeil)); err != nil {
		l.Error("setting MinerGasConfig failed", zap.Error(err))
		return
	}
	l.Info("setting MinerGasConfig success")
}

func (caller *rpcCaller) SendTx(c *cli.Context) {
	if checkSync() {
		return
	}

	_, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error", zap.Error(err))
		return
	}

	if len(cParams) != 4 && len(cParams) != 5 {
		l.Error("parameter includes：to value gasPrice gasLimit extraData, extraData is optional")
		return
	}

	toAddress, err := CheckAndChangeHexToAddress(cParams[0])
	if err != nil {
		l.Error("the to address is invalid", zap.Error(err))
		return
	}

	value, err := MoneyValueToCSCoin(cParams[1])
	if err != nil {
		l.Error("the parameter value invalid", zap.Error(err))
		return
	}

	gasPrice, err := MoneyValueToCSCoin(cParams[2])
	if err != nil {
		l.Error("the parameter gasPrice invalid", zap.Error(err))
		return
	}

	gasLimit, err := strconv.ParseUint(cParams[3], 10, 64)
	if err != nil {
		l.Error("the parameter gasLimit invalid", zap.Error(err))
		return
	}

	extraData := make([]byte, 0)
	if len(cParams) == 5 {
		extraData = []byte(cParams[4])
	}

	var resp common.Hash
	l.Debug("the from is: ", zap.String("from", defaultAccount.Hex()))
	l.Debug("the to is: ", zap.String("to", toAddress.Hex()))
	l.Debug("the value is:", zap.String("value", MoneyWithUnit(cParams[1])))
	l.Debug("the gasPrice is:", zap.String("gasPrice", MoneyWithUnit(cParams[2])))
	l.Debug("the ExtraData is: ", zap.Uint8s("ExtraData", extraData))
	if err = client.Call(&resp, getDipperinRpcMethodByName("SendTransaction"), defaultAccount, toAddress, value, gasPrice, gasLimit, extraData, nil); err != nil {
		l.Error("call send transaction error", zap.Error(err))
		return
	}
	l.Info("SendTransaction result", zap.String("txId", resp.Hex()))
}

//send transaction
func (caller *rpcCaller) SendTransaction(c *cli.Context) {
	if checkSync() {
		return
	}

	mName, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error", zap.Error(err))
		return
	}

	if len(cParams) != 5 && len(cParams) != 6 {
		l.Error("parameter includes：from to value gasPrice gasLimit extraData, extraData is optional")
		return
	}

	from, err := CheckAndChangeHexToAddress(cParams[0])
	if err != nil {
		l.Error("the from address is invalid", zap.Error(err))
		return
	}

	to, err := CheckAndChangeHexToAddress(cParams[1])
	if err != nil {
		l.Error("the to address is invalid", zap.Error(err))
		return
	}

	value, err := MoneyValueToCSCoin(cParams[2])
	if err != nil {
		l.Error("the parameter value invalid", zap.Error(err))
		return
	}

	gasPrice, err := MoneyValueToCSCoin(cParams[3])
	if err != nil {
		l.Error("the parameter gasPrice is invalid", zap.Error(err))
		return
	}

	gasLimit, err := strconv.ParseUint(cParams[4], 10, 64)
	if err != nil {
		l.Error("the parameter gasLimit is invalid", zap.Error(err))
		return
	}

	ExtraData := make([]byte, 0)
	if len(cParams) == 6 {
		ExtraData = []byte(cParams[5])
	}

	var resp common.Hash
	l.Debug("the from is: ", zap.String("from", from.Hex()))
	l.Debug("the to is: ", zap.String("to", to.Hex()))
	l.Debug("the value is:", zap.String("value", MoneyWithUnit(cParams[2])))
	l.Debug("the gasPrice is:", zap.String("gasPrice", MoneyWithUnit(cParams[3])))
	l.Debug("the gasLimit is:", zap.Uint64("gasLimit", gasLimit))
	l.Debug("the ExtraData is: ", zap.Uint8s("ExtraData", ExtraData))
	if err = client.Call(&resp, getDipperinRpcMethodByName(mName), from, to, value, gasPrice, gasLimit, ExtraData, nil); err != nil {
		l.Error("call send transaction", zap.Error(err))
		return
	}
	l.Info("SendTransaction result", zap.String("txId", resp.Hex()))
}

//check transaction from transaction hash
func (caller *rpcCaller) Transaction(c *cli.Context) {
	if checkSync() {
		return
	}

	mName, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error")
		return
	}

	if len(cParams) < 1 {
		l.Error("Transaction need：txHash")
		return
	}

	tmpHash, err := hexutil.Decode(cParams[0])
	if err != nil {
		l.Error("the err is:", zap.Error(err))
		l.Error("Transaction decode error")
		return
	}

	var hash common.Hash
	_ = copy(hash[:], tmpHash)

	var resp rpcinterface.TransactionResp
	if err = client.Call(&resp, getDipperinRpcMethodByName(mName), hash); err != nil {
		l.Error("Call Transaction", zap.Error(err))
		return
	}

	printTransactionInfo(resp)
}

func (caller *rpcCaller) GetReceiptByTxHash(c *cli.Context) {
	_, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("GetReceiptByTxHash  getRpcMethodAndParam error")
		return
	}

	if len(cParams) < 1 {
		l.Error("GetReceiptByTxHash need：txHash")
		return
	}

	tmpHash, err := hexutil.Decode(cParams[0])
	if err != nil {
		l.Error("GetReceiptByTxHash decode error")
		return
	}

	var hash common.Hash
	copy(hash[:], tmpHash)

	var resp model2.Receipt
	if err = client.Call(&resp, getDipperinRpcMethodByName("GetReceiptByTxHash"), hash); err != nil {
		l.Error("Call GetReceiptByTxHash", zap.Error(err))
		return
	}

	fmt.Printf("ReceiptInfo:\r\n")
	fmt.Println(resp.String())
}

func (caller *rpcCaller) GetTxActualFee(c *cli.Context) {
	_, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("GetTxActualFee  getRpcMethodAndParam error")
		return
	}

	if len(cParams) < 1 {
		l.Error("GetTxActualFee need：txHash")
		return
	}

	tmpHash, err := hexutil.Decode(cParams[0])
	if err != nil {
		l.Error("GetTxActualFee decode error")
		return
	}

	var hash common.Hash
	copy(hash[:], tmpHash)

	var resp rpcinterface.CurBalanceResp
	if err = client.Call(&resp, getDipperinRpcMethodByName("GetTxActualFee"), hash); err != nil {
		l.Error("Call GetTxActualFee", zap.Error(err))
		return
	}

	txFee, err := InterToDecimal(resp.Balance, consts.UnitDecimalBits)
	if err != nil {
		l.Error("can't get tx actual fee", zap.Error(err))
	} else {
		l.Info("the tx actual fee is", zap.String("txActualFee", txFee+consts.CoinWuName))
	}
}

func (caller *rpcCaller) SuggestGasPrice(c *cli.Context) {
	mName, _, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error")
		return
	}
	var resp rpcinterface.CurBalanceResp
	if err = client.Call(&resp, getDipperinRpcMethodByName(mName)); err != nil {
		l.Error("call SuggestGasPrice failed", zap.Error(err))
		return
	}

	gasPrice, err := InterToDecimal(resp.Balance, consts.UnitDecimalBits)
	if err != nil {
		l.Error("can't get suggest gas price", zap.Error(err))
	} else {
		l.Info("the gas price is", zap.String("gasPrice", gasPrice+consts.CoinWuName))
	}
}

func (caller *rpcCaller) GetLogs(c *cli.Context) {
	//BlockHash *common.Hash, FromBlock *big.Int, ToBlock *big.Int, Addresses []common.Address, Topics [][]common.Hash
	params := c.String("p")
	l.Debug("GetLogs the params is:", zap.String("params", params))
	var filterParams FilterParams
	if err := json.Unmarshal([]byte(params), &filterParams); err != nil {
		l.Error("json.Unmarshal failed", zap.Error(err))
		l.Info("json needed")
		fmt.Println("example1:", `{"block_hash":"0x000023e18421a0abfceea172867b9b4a3bcf593edd0b504554bb7d1cf5f5e7b7","addresses":["0x0014049F835be46352eD0Ec6B819272A2c8cF4feA10f"],"topics":[["0x0b5d2220daf8f0dfd95983d2ce625affbb7183c991271f49d818b4a64a268dbb"]]}`)
		fmt.Println("example2:", `{"from_block":10,"to_block":500,"addresses":["0x0014049F835be46352eD0Ec6B819272A2c8cF4feA10f"],"topics":[["0x0b5d2220daf8f0dfd95983d2ce625affbb7183c991271f49d818b4a64a268dbb"]]}`)
		return
	}

	if filterParams.BlockHash.IsEmpty() {
		l.Debug("the blockHash is", zap.String("blockHash", "empty"))
	} else {
		l.Debug("the blockHash is", zap.Any("blockHash", filterParams.BlockHash))
	}
	l.Debug("the fromBlock is", zap.Uint64("num", filterParams.FromBlock))
	if filterParams.ToBlock == uint64(0) {
		l.Debug("the toBlock is", zap.String("num", "currentBlock"))
	} else {
		l.Debug("the toBlock is", zap.Uint64("num", filterParams.ToBlock))
	}
	l.Debug("the contractAddresses is", zap.Any("addresses", filterParams.Addresses))
	l.Debug("the topics is", zap.Any("topics", filterParams.Topics))

	var resp []model2.Log
	if err := client.Call(&resp, getDipperinRpcMethodByName("GetLogs"), filterParams.BlockHash, filterParams.FromBlock, filterParams.ToBlock, filterParams.Addresses, filterParams.Topics); err != nil {
		l.Error("Call GetLogs failed", zap.Error(err))
		return
	}

	if len(resp) == 0 {
		l.Info("logs not found")
		return
	}
	fmt.Println("found logs:")
	for _, lg := range resp {
		fmt.Println(lg.String())
	}
	return
}

func (caller *rpcCaller) GetReceiptsByBlockNum(c *cli.Context) {
	_, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("GetReceiptsByBlockNum  getRpcMethodAndParam error")
		return
	}
	if len(cParams) < 1 {
		l.Error("GetReceiptsByBlockNum need：blockNum")
		return
	}
	blockNum, err := strconv.Atoi(cParams[0])
	if err != nil {
		l.Error("the blockNumber error")
		return
	}

	var resp model2.Receipts
	if err = client.Call(&resp, getDipperinRpcMethodByName("GetReceiptsByBlockNum"), blockNum); err != nil {
		l.Error("Call GetReceiptsByBlockNum", zap.Error(err))
		return
	}
	fmt.Printf("ReceiptInfos:\r\n")
	fmt.Println(resp)
}

//List Wallet
func (caller *rpcCaller) ListWallet(c *cli.Context) {
	mName, _, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error")
		return
	}
	l.Debug(getDipperinRpcMethodByName(mName))
	var resp []accountsbase.WalletIdentifier
	if err = client.Call(&resp, getDipperinRpcMethodByName(mName)); err != nil {
		l.Error("Call ListWallet", zap.Error(err))
		return
	}
	fmt.Println("  Call ListWallet Result:")
	for _, w := range resp {
		fmt.Println("    Wallet Info:", w.String())
	}
	//l.Info("Call ListWallet", "wallet list", resp)
}

//List Wallet Account
func (caller *rpcCaller) ListWalletAccount(c *cli.Context) {
	mName, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error")
		return
	}
	var identifier accountsbase.WalletIdentifier
	if len(cParams) == 0 {
		identifier = defaultWallet
	} else if len(cParams) == 2 {
		identifier.Path, identifier.WalletName = ParseWalletPathAndName(cParams[1])
		l.Debug("ListWalletAccount", zap.String("walletPath", identifier.Path), zap.String("walletName", identifier.WalletName))
		if cParams[0] == "SoftWallet" {
			identifier.WalletType = accountsbase.SoftWallet
		} else if cParams[0] == "LedgerWallet" {
			identifier.WalletType = accountsbase.TrezorWallet
		} else if cParams[0] == "TrezorWallet" {
			identifier.WalletType = accountsbase.TrezorWallet
		} else {
			l.Error("Wallet Type error")
			return
		}
	} else {
		l.Error("list Wallet Account need：Type Path")
		return
	}

	l.Debug(getDipperinRpcMethodByName(mName))
	var resp []accountsbase.Account

	if err = client.Call(&resp, getDipperinRpcMethodByName(mName), identifier); err != nil {
		l.Error("Call ListWallet", zap.Error(err))
		return
	}

	l.Info("Call ListWalletAccount result:")
	for _, account := range resp {
		fmt.Println("\taddress:", account.Address.Hex())
	}

}

//Establish Wallet
func (caller *rpcCaller) EstablishWallet(c *cli.Context) {
	mName, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error")
		return
	}

	if len(cParams) != 3 {
		l.Error("EstablishWallet need：Type Path password")
		return
	}

	var identifier accountsbase.WalletIdentifier
	identifier.Path, identifier.WalletName = ParseWalletPathAndName(cParams[1])

	if cParams[0] == "SoftWallet" {
		identifier.WalletType = accountsbase.SoftWallet
	} else if cParams[0] == "LedgerWallet" {
		identifier.WalletType = accountsbase.TrezorWallet
	} else if cParams[0] == "TrezorWallet" {
		identifier.WalletType = accountsbase.TrezorWallet
	} else {
		l.Error("Wallet Type error")
		return
	}

	password := cParams[2]
	/*	var passPhrase string
		if len(cParams) == 4 {
			passPhrase = cParams[3]
		} else {
			passPhrase = ""
		}*/

	var resp string

	/*	l.Info("the identifier is: ","identifier",identifier)
		l.Info("the password is: ","password",password)
		l.Info("the passPhrase is: ","passPhrase",passPhrase)

		l.Info("call EstablishWallet")*/
	if err = client.Call(&resp, getDipperinRpcMethodByName(mName), password, "", identifier); err != nil {
		l.Error("Call EstablishWallet", zap.Error(err))
		return
	}
	resp = strings.Replace(resp, " ", ",", -1)
	l.Info("Call EstablishWallet", zap.String("mnemonic", resp))
}

//Restore Wallet
func (caller *rpcCaller) RestoreWallet(c *cli.Context) {
	mName, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error")
		return
	}
	if len(cParams) < 4 {
		l.Error("RestoreWallet need：type, walletIdentifier, password, mnemonic")
		return
	}

	var identifier accountsbase.WalletIdentifier
	identifier.Path, identifier.WalletName = ParseWalletPathAndName(cParams[1])

	if cParams[0] == "SoftWallet" {
		identifier.WalletType = accountsbase.SoftWallet
	} else if cParams[0] == "LedgerWallet" {
		identifier.WalletType = accountsbase.TrezorWallet
	} else if cParams[0] == "TrezorWallet" {
		identifier.WalletType = accountsbase.TrezorWallet
	} else {
		l.Error("Wallet Type error")
		return
	}

	var mnemonic string
	for i := 3; i < len(cParams); i++ {
		mnemonic += cParams[i]
		if i != len(cParams)-1 {
			mnemonic += " "
		}
	}
	password := cParams[2]
	//passPhrase := cParams[3]

	l.Debug(getDipperinRpcMethodByName(mName))
	var resp interface{}
	l.Debug("the identifier is: ", zap.Any("identifier", identifier))
	l.Debug("the password is: ", zap.String("password", password))
	//l.Debug("the passPhrase is: ", "passPhrase", passPhrase)
	l.Debug("the mnemonic is: ", zap.String("mnemonic", mnemonic))

	if err = client.Call(resp, getDipperinRpcMethodByName(mName), password, mnemonic, "", identifier); err != nil {
		l.Error("Call RestoreWallet", zap.Error(err))
		return
	}
	l.Info("Call RestoreWallet success")
}

//Open Wallet
func (caller *rpcCaller) OpenWallet(c *cli.Context) {
	mName, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error")
		return
	}

	var identifier accountsbase.WalletIdentifier
	var password string
	if len(cParams) == 1 {
		identifier = defaultWallet
		password = cParams[0]
	} else if len(cParams) == 3 {
		identifier.Path, identifier.WalletName = ParseWalletPathAndName(cParams[1])
		if cParams[0] == "SoftWallet" {
			identifier.WalletType = accountsbase.SoftWallet
		} else if cParams[0] == "LedgerWallet" {
			identifier.WalletType = accountsbase.LedgerWallet
		} else if cParams[0] == "TrezorWallet" {
			identifier.WalletType = accountsbase.TrezorWallet
		} else {
			l.Error("Wallet Type error")
			return
		}

		password = cParams[2]
	} else {
		l.Error("OpenWallet need：Type Path password")
		return
	}

	l.Debug(getDipperinRpcMethodByName(mName))
	var resp interface{}

	if err = client.Call(&resp, getDipperinRpcMethodByName(mName), password, identifier); err != nil {
		l.Error("Call OpenWallet", zap.Error(err))
		return
	}
	l.Info("Call OpenWallet success")

}

//Close Wallet
func (caller *rpcCaller) CloseWallet(c *cli.Context) {
	mName, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error")
		return
	}

	var identifier accountsbase.WalletIdentifier
	if len(cParams) == 0 {
		identifier = defaultWallet
	} else if len(cParams) == 2 {
		identifier.Path, identifier.WalletName = ParseWalletPathAndName(cParams[1])
		if cParams[0] == "SoftWallet" {
			identifier.WalletType = accountsbase.SoftWallet
		} else if cParams[0] == "LedgerWallet" {
			identifier.WalletType = accountsbase.TrezorWallet
		} else if cParams[0] == "TrezorWallet" {
			identifier.WalletType = accountsbase.TrezorWallet
		} else {
			l.Error("Wallet Type error")
			return
		}
	} else {
		l.Error("CloseWallet need：Type Path")
		return
	}
	l.Debug(getDipperinRpcMethodByName(mName))
	var resp interface{}

	if err = client.Call(&resp, getDipperinRpcMethodByName(mName), identifier); err != nil {
		l.Error("Call CloseWallet", zap.Error(err))
		return
	}
	l.Info("Call CloseWallet success")
}

//AddAccount
func (caller *rpcCaller) AddAccount(c *cli.Context) {
	mName, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error")
		return
	}

	var identifier accountsbase.WalletIdentifier
	if len(cParams) == 0 {
		identifier = defaultWallet
	} else if len(cParams) == 2 {
		identifier.Path, identifier.WalletName = ParseWalletPathAndName(cParams[1])
		if cParams[0] == "SoftWallet" {
			identifier.WalletType = accountsbase.SoftWallet
		} else if cParams[0] == "LedgerWallet" {
			identifier.WalletType = accountsbase.TrezorWallet
		} else if cParams[0] == "TrezorWallet" {
			identifier.WalletType = accountsbase.TrezorWallet
		} else {
			l.Error("Wallet Type error")
			return
		}
	} else {
		l.Error("AddAccount need：Type Path")
		return
	}

	derivationPath := ""
	l.Debug(getDipperinRpcMethodByName(mName))
	var resp accountsbase.Account
	if err = client.Call(&resp, getDipperinRpcMethodByName(mName), derivationPath, identifier); err != nil {
		l.Error("Call AddAccount", zap.Error(err))
		return
	}
	l.Info("Call AddAccount", zap.String("Added Account Address", resp.Address.Hex()))
}

func (caller *rpcCaller) SendRegisterTx(c *cli.Context) {
	if checkSync() {
		return
	}

	_, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error", zap.Error(err))
		return
	}

	if len(cParams) != 3 {
		l.Error("SendRegisterTransaction need：stake gasPrice gasLimit")
		return
	}

	stake, err := MoneyValueToCSCoin(cParams[0])
	if err != nil {
		l.Error("the parameter stake invalid")
		return
	}

	gasPrice, err := MoneyValueToCSCoin(cParams[1])
	if err != nil {
		l.Error("the parameter gasPrice invalid", zap.Error(err))
		return
	}

	gasLimit, err := strconv.ParseUint(cParams[2], 10, 64)
	if err != nil {
		l.Error("the parameter gasLimit invalid", zap.Error(err))
		return
	}

	l.Debug("the stake is:", zap.String("stake", MoneyWithUnit(cParams[0])))
	l.Debug("the gasPrice is:", zap.String("gasPrice", MoneyWithUnit(cParams[1])))
	l.Debug("the gasLimit is:", zap.Uint64("gasLimit", gasLimit))
	var resp common.Hash
	if err = client.Call(&resp, getDipperinRpcMethodByName("SendRegisterTransaction"), defaultAccount, stake, gasPrice, gasLimit, nil); err != nil {
		l.Error("call send transaction", zap.Error(err))
		return
	}
	l.Info("SendRegisterTransaction result", zap.String("txId", resp.Hex()))
	addTrackingAccount(defaultAccount)
	RecordRegistration(resp.Hex())
}

//send Register transaction
func (caller *rpcCaller) SendRegisterTransaction(c *cli.Context) {
	if checkSync() {
		return
	}

	mName, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error")
		return
	}
	if len(cParams) != 4 {
		l.Error("SendRegisterTransaction need：from stake gasPrice gasLimit")
		return
	}

	var resp common.Hash
	from, err := CheckAndChangeHexToAddress(cParams[0])
	if err != nil {
		l.Error("the from address is invalid", zap.Error(err))
		return
	}

	stake, err := MoneyValueToCSCoin(cParams[1])
	if err != nil {
		l.Error("the parameter stake invalid")
		return
	}

	gasPrice, err := MoneyValueToCSCoin(cParams[2])
	if err != nil {
		l.Error("the parameter gasPrice invalid", zap.Error(err))
		return
	}

	gasLimit, err := strconv.ParseUint(cParams[3], 10, 64)
	if err != nil {
		l.Error("the parameter gasLimit invalid", zap.Error(err))
		return
	}

	l.Debug("the stake is:", zap.String("stake", MoneyWithUnit(cParams[1])))
	l.Debug("the gasPrice is:", zap.String("gasPrice", MoneyWithUnit(cParams[2])))
	l.Debug("the gasLimit is:", zap.Uint64("gasLimit", gasLimit))
	if err = client.Call(&resp, getDipperinRpcMethodByName(mName), from, stake, gasPrice, gasLimit, nil); err != nil {
		l.Error("call send transaction", zap.Error(err))
		return
	}
	l.Info("SendRegisterTransaction result", zap.String("txId", resp.Hex()))
	addTrackingAccount(from)
	RecordRegistration(resp.Hex())
}

func (caller *rpcCaller) SendUnStakeTx(c *cli.Context) {
	if checkSync() {
		return
	}

	_, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error")
		return
	}

	if len(cParams) != 2 {
		l.Error("SendUnStakeTransaction need gasPrice and gasLimit")
		return
	}

	var resp common.Hash
	gasPrice, err := MoneyValueToCSCoin(cParams[0])
	if err != nil {
		l.Error("the parameter gasPrice invalid", zap.Error(err))
		return
	}

	gasLimit, err := strconv.ParseUint(cParams[1], 10, 64)
	if err != nil {
		l.Error("the parameter gasLimit invalid", zap.Error(err))
		return
	}

	//l.Info("the gasPrice is:", "gasPrice", MoneyWithUnit(cParams[0]))
	//l.Info("the gasLimit is:", "gasLimit", gasLimit)
	if err = client.Call(&resp, getDipperinRpcMethodByName("SendUnStakeTransaction"), defaultAccount, gasPrice, gasLimit, nil); err != nil {
		l.Error("call send transaction", zap.Error(err))
		return
	}
	l.Info("SendUnStakeTransaction result", zap.String("txId", resp.Hex()))
}

//send UnStake transaction
func (caller *rpcCaller) SendUnStakeTransaction(c *cli.Context) {
	if checkSync() {
		return
	}

	mName, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error")
		return
	}

	if len(cParams) != 3 {
		l.Error("SendUnStakeTransaction need：from gasPrice gasLimit")
		return
	}

	var resp common.Hash
	from, err := CheckAndChangeHexToAddress(cParams[0])
	if err != nil {
		l.Error("the from address is invalid", zap.Error(err))
		return
	}

	gasPrice, err := MoneyValueToCSCoin(cParams[1])
	if err != nil {
		l.Error("the parameter gasPrice invalid", zap.Error(err))
		return
	}

	gasLimit, err := strconv.ParseUint(cParams[2], 10, 64)
	if err != nil {
		l.Error("the parameter gasLimit invalid", zap.Error(err))
		return
	}

	l.Debug("the gasPrice is:", zap.String("gasPrice", MoneyWithUnit(cParams[1])))
	l.Debug("the gasLimit is:", zap.Uint64("gasLimit", gasLimit))
	if err = client.Call(&resp, getDipperinRpcMethodByName(mName), from, gasPrice, gasLimit, nil); err != nil {
		l.Error("call send transaction", zap.Error(err))
		return
	}
	l.Info("SendUnStakeTransaction result", zap.String("txId", resp.Hex()))
}

func (caller *rpcCaller) SendCancelTx(c *cli.Context) {
	if checkSync() {
		return
	}

	_, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error")
		return
	}

	if len(cParams) != 2 {
		l.Error("SendCancelTransaction need gasPrice gasLimit")
		return
	}

	var resp common.Hash
	gasPrice, err := MoneyValueToCSCoin(cParams[0])
	if err != nil {
		l.Error("the parameter gasPrice invalid", zap.Error(err))
		return
	}

	gasLimit, err := strconv.ParseUint(cParams[1], 10, 64)
	if err != nil {
		l.Error("the parameter gasLimit invalid", zap.Error(err))
		return
	}

	l.Debug("the gasPrice is:", zap.String("gasPrice", MoneyWithUnit(cParams[0])))
	l.Debug("the gasLimit is:", zap.Uint64("gasLimit", gasLimit))
	if err = client.Call(&resp, getDipperinRpcMethodByName("SendCancelTransaction"), defaultAccount, gasPrice, gasLimit, nil); err != nil {
		l.Error("call send transaction", zap.Error(err))
		return
	}
	l.Info("SendCancelTransaction result", zap.String("txId", resp.Hex()))
	removeTrackingAccount(defaultAccount)
	RemoveRegistration()
}

//send Cancel transaction
func (caller *rpcCaller) SendCancelTransaction(c *cli.Context) {
	if checkSync() {
		return
	}

	mName, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error")
		return
	}

	if len(cParams) != 3 {
		l.Error("SendCancelTransaction need：from gasPrice gasLimit")
		return
	}

	var resp common.Hash
	from, err := CheckAndChangeHexToAddress(cParams[0])
	if err != nil {
		l.Error("the from address is invalid", zap.Error(err))
		return
	}

	gasPrice, err := MoneyValueToCSCoin(cParams[1])
	if err != nil {
		l.Error("the parameter gasPrice invalid", zap.Error(err))
		return
	}

	gasLimit, err := strconv.ParseUint(cParams[2], 10, 64)
	if err != nil {
		l.Error("the parameter gasLimit invalid", zap.Error(err))
		return
	}

	l.Debug("the gasPrice is:", zap.String("gasPrice", MoneyWithUnit(cParams[1])))
	l.Debug("the gasLimit is:", zap.Uint64("gasLimit", gasLimit))
	if err = client.Call(&resp, getDipperinRpcMethodByName(mName), from, gasPrice, gasLimit, nil); err != nil {
		l.Error("call send transaction", zap.Error(err))
		return
	}
	l.Info("SendCancelTransaction result", zap.String("txId", resp.Hex()))
	removeTrackingAccount(from)
	RemoveRegistration()
}

func (caller *rpcCaller) GetVerifiersBySlot(c *cli.Context) {
	if checkSync() {
		return
	}

	mName, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error")
		return
	}

	if len(cParams) < 1 {
		l.Error("GetVerifiersBySlot need：slotNum")
		return
	}

	var resp []common.Address
	slotNum, err := strconv.Atoi(cParams[0])
	if err != nil {
		l.Error("the parameter slotNum invalid")
		return
	}
	if err = client.Call(&resp, getDipperinRpcMethodByName(mName), slotNum); err != nil {
		l.Error("GetVerifiersBySlot", zap.Error(err))
		return
	}
	l.Info("GetVerifiersBySlot result")
	for _, tmpAddress := range resp {
		fmt.Println("\t", "address:", tmpAddress.Hex())
		//l.Info("verifier address is:", "verifier", tmpAddress.Hex())
	}
}

// VerifierStatus call to get the verifier status
func (caller *rpcCaller) VerifierStatus(c *cli.Context) {
	//params is a comma separated list of Addresses
	mName, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error")
		return
	}
	if len(cParams) != 0 && len(cParams) != 1 {
		l.Error("parameter error")
		return
	}

	var addr common.Address
	if len(cParams) == 0 {
		addr = getDefaultAccount()
	} else {
		addr, err = CheckAndChangeHexToAddress(cParams[0])
		if err != nil {
			l.Error("the input address is invalid", zap.Error(err))
			return
		}
	}

	var resp rpcinterface.VerifierStatus
	l.Debug(getDipperinRpcMethodByName(mName))

	if err = client.Call(&resp, getDipperinRpcMethodByName(mName), addr); err != nil {
		l.Error("call verifier's status error", zap.Error(err))
		return
	}

	balance, err := CSCoinToMoneyValue(resp.Balance)
	if err != nil {
		l.Error("The address has no balance, balance = 0 DIP")
		balance = "0"
	}

	if resp.Status == VerifierStatusNoRegistered {
		l.Info("Verifier status", zap.String("status", resp.Status), zap.String("balance", balance))
		return
	}

	if resp.Status == VerifiedStatusUnstaked {
		l.Info("Verifier status", zap.String("status", resp.Status), zap.String("balance", balance))
		return
	}

	stake, err := CSCoinToMoneyValue(resp.Stake)
	if err != nil {
		l.Error("The address has no stake, stake = 0 DIP")
	}
	if resp.Status == VerifierStatusRegistered {
		l.Info("Verifier status", zap.String("status", resp.Status), zap.String("balance", balance), zap.String("stake", stake), zap.Uint64("reputation", resp.Reputation), zap.Bool("is current verifier", resp.IsCurrentVerifier))
	}

	if resp.Status == VerifiedStatusCanceled {
		l.Info("Verifier status", zap.String("status", resp.Status), zap.String("balance", balance), zap.String("stake", stake), zap.Uint64("reputation", resp.Reputation))
	}
}

// SetBftSigner set the signer of the wallet
func (caller *rpcCaller) SetBftSigner(c *cli.Context) {
	mName, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error")
		return
	}

	if len(cParams) != 1 {
		l.Error("need parameter: address")
		return
	}

	addr, err := CheckAndChangeHexToAddress(cParams[0])
	if err != nil {
		l.Error("the from address is invalid", zap.Error(err))
		return
	}

	l.Debug(getDipperinRpcMethodByName(mName))
	var resp interface{}
	if err = client.Call(&resp, getDipperinRpcMethodByName(mName), addr); err != nil {
		l.Error("Set wallet signer（default account）", zap.Error(err))
		return
	}
	l.Info("Set wallet signer（default account）succeed")
}

func (caller *rpcCaller) GetDefaultAccountStake(c *cli.Context) {
	loadDefaultAccountStake()
	PrintDefaultAccountStake()
}

// CurrentStake call to get the current account stake
func (caller *rpcCaller) CurrentStake(c *cli.Context) {
	//params is a comma separated list of Addresses
	mName, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error")
		return
	}
	if len(cParams) != 0 && len(cParams) != 1 {
		l.Error("parameter error")
		return
	}

	var addr common.Address
	if len(cParams) == 0 {
		addr = getDefaultAccount()
	} else {
		addr, err = CheckAndChangeHexToAddress(cParams[0])
		if err != nil {
			l.Error("the input address is invalid", zap.Error(err))
			return
		}
	}

	var resp rpcinterface.CurBalanceResp
	l.Debug(getDipperinRpcMethodByName(mName))

	if err = client.Call(&resp, getDipperinRpcMethodByName(mName), addr); err != nil {
		l.Error("lookup current stake", zap.Error(err))
		return
	}

	stake, err := CSCoinToMoneyValue(resp.Balance)
	if err != nil {
		l.Error("the address isn't on the block chain stake=0")
	} else {
		l.Info("address current stake is:", zap.String("stake", stake))
	}
}

// CurrentReputation call the method to get the current reputation value
func (caller *rpcCaller) CurrentReputation(c *cli.Context) {
	//params is a comma separated list of Addresses
	mName, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error")
		return
	}
	if len(cParams) != 0 && len(cParams) != 1 {
		l.Error("parameter error")
		return
	}

	var addr common.Address
	if len(cParams) == 0 {
		addr = getDefaultAccount()
	} else {
		addr, err = CheckAndChangeHexToAddress(cParams[0])
		if err != nil {
			l.Error("the input address is invalid", zap.Error(err))
			return
		}
	}

	var resp uint64
	if err = client.Call(&resp, getDipperinRpcMethodByName(mName), addr); err != nil {
		l.Error("lookup current reputation error", zap.Error(err))
		return
	}
	l.Info("address current reputation is:", zap.Uint64("reputation", resp))
}

func inDefaultVs(addr common.Address) (bool, string) {
	for i, v := range chain.VerifierAddress {
		if addr.IsEqual(v) {
			name := fmt.Sprintf("default_v%v", i)
			return true, name
		}
	}
	return false, ""
}

func (caller *rpcCaller) GetCurVerifiers(c *cli.Context) {
	if checkSync() {
		return
	}

	mName, _, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error")
		return
	}

	var resp []common.Address
	if err = client.Call(&resp, getDipperinRpcMethodByName(mName)); err != nil {
		l.Error("call failed", zap.Error(err))
		return
	}

	for _, a := range resp {
		result, name := inDefaultVs(a)
		if result {
			l.Info("Current Verifiers:", zap.String("address", a.Hex()), zap.Bool("is_default", result), zap.String("name", name))
		} else {
			l.Info("Current Verifiers:", zap.String("address", a.Hex()), zap.Bool("is_default", result))
		}

	}
}

func (caller *rpcCaller) GetNextVerifiers(c *cli.Context) {
	if checkSync() {
		return
	}

	mName, _, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error")
		return
	}

	var resp []common.Address
	if err = client.Call(&resp, getDipperinRpcMethodByName(mName)); err != nil {
		l.Error("call failed", zap.Error(err))
		return
	}

	for _, a := range resp {
		result, name := inDefaultVs(a)
		if result {
			l.Info("Next Verifiers:", zap.String("address", a.Hex()), zap.Bool("is_default", result), zap.String("name", name))
		} else {
			l.Info("Next Verifiers:", zap.String("address", a.Hex()), zap.Bool("is_default", result))
		}

	}
}

func getNonceInfo(c *cli.Context) (nonce uint64, err error) {
	mName, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error")
		return 0, errors.New("getRpcMethodAndParam error")
	}

	if len(cParams) != 0 && len(cParams) != 1 {
		l.Error("parameter error")
		return 0, errors.New("parameter error")
	}

	var addr common.Address
	if len(cParams) == 0 {
		addr = getDefaultAccount()
	} else {
		addr, err = CheckAndChangeHexToAddress(cParams[0])
		if err != nil {
			l.Error("the input address is invalid", zap.Error(err))
			return
		}
	}

	var resp uint64
	if err = client.Call(&resp, getDipperinRpcMethodByName(mName), addr); err != nil {
		l.Error("call GetTransactionNonce", zap.Error(err))
		return 0, err
	}
	return resp, nil
}

func (caller *rpcCaller) GetTransactionNonce(c *cli.Context) {
	if checkSync() {
		return
	}

	nonce, err := getNonceInfo(c)
	if err != nil {
		l.Error("GetTransactionNonce error", zap.Error(err))
		return
	}
	l.Info("the address nonce from chain is:", zap.Uint64("nonce", nonce))
}

func (caller *rpcCaller) GetAddressNonceFromWallet(c *cli.Context) {
	nonce, err := getNonceInfo(c)
	if err != nil {
		l.Error("GetTransactionNonce error", zap.Error(err))
		return
	}

	l.Info("the address nonce from wallet is:", zap.Uint64("nonce", nonce))
}

func initWallet(path, password, passPhrase string) (err error) {
	var identifier accountsbase.WalletIdentifier
	identifier.WalletType = accountsbase.SoftWallet
	identifier.Path, identifier.WalletName = ParseWalletPathAndName(path)

	//open
	exit, _ := softwallet.PathExists(identifier.Path)
	if exit {
		l.Info("open wallet", zap.Any("identifier", identifier))
		var resp interface{}
		if err = client.Call(&resp, getDipperinRpcMethodByName("OpenWallet"), password, identifier); err != nil {
			l.Error("open Wallet err", zap.Error(err))
			return err
		}
	} else {
		l.Info("establish wallet", zap.Any("identifier", identifier))
		var mnemonic string
		if err = client.Call(&mnemonic, getDipperinRpcMethodByName("EstablishWallet"), password, passPhrase, identifier); err != nil {
			l.Error("Call EstablishWallet", zap.Error(err))
			return err
		}
		mnemonic = strings.Replace(mnemonic, " ", ",", -1)
		l.Info("EstablishWallet mnemonic is:", zap.String("mnemonic", mnemonic))
	}
	return nil
}

func getDefaultAccount() common.Address {
	var resp []accountsbase.WalletIdentifier
	l.Debug("getDefaultAccount")

	if err := client.Call(&resp, getDipperinRpcMethodByName("ListWallet")); err != nil {
		l.Error("Call ListWallet", zap.Error(err))
		return common.Address{}
	}

	var respA []accountsbase.Account
	if err := client.Call(&respA, getDipperinRpcMethodByName("ListWalletAccount"), resp[0]); err != nil {
		l.Error("Call ListWallet", zap.Error(err))
		return common.Address{}
	}
	return respA[0].Address
}

/*func (caller *rpcCaller) SetDefaultAccount(c *cli.Context) {
	_, cParams, err := getRpcMethodAndParam(c)

	if err != nil {
		l.Error("get rpc method and param error")
		return
	}

	addr, err := CheckAndChangeHexToAddress(cParams[0])
	if err != nil {
		l.Error("the input address is invalid", zap.Error(err))
	}

	defaultAccount = addr
	l.Info("default account set success", "account", defaultAccount.Hex())
}*/

func getDefaultWallet() accountsbase.WalletIdentifier {
	var resp []accountsbase.WalletIdentifier
	if err := client.Call(&resp, getDipperinRpcMethodByName("ListWallet")); err != nil {
		l.Error("Call ListWallet", zap.Error(err))
		return accountsbase.WalletIdentifier{}
	}
	return resp[0]
}

//if user applies for registering verifier, record,
// creating a file in $Home/.dipperin
// for startup check, quit if node type is not verifier
func RecordRegistration(txHash string) {
	confPath := filepath.Join(util.HomeDir(), ".dipperin", "registration")
	exist, _ := softwallet.PathExists(confPath)
	if !exist {
		if err := os.MkdirAll(filepath.Dir(confPath), 0766); err != nil {
			l.Error("can't make dir", zap.Error(err))
			return
		}
	}

	if err := ioutil.WriteFile(confPath, util.StringifyJsonToBytes(struct {
		TxHash string `json:"tx_hash"`
	}{
		TxHash: txHash,
	}), 0644); err != nil {
		l.Error("can't record registration")
	}
}

// remove file after unregister transaction
func RemoveRegistration() {
	confPath := filepath.Join(util.HomeDir(), ".dipperin", "registration")
	exist, _ := softwallet.PathExists(confPath)
	if !exist {
		return
	}

	if err := os.Remove(confPath); err != nil {
		l.Error("can't remove record registration")
		return
	}
}
