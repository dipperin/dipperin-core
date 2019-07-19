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
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/core/accounts"
	"github.com/dipperin/dipperin-core/core/accounts/soft-wallet"
	"github.com/dipperin/dipperin-core/core/chain"
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/dipperin/dipperin-core/core/rpc-interface"
	"github.com/dipperin/dipperin-core/core/vm/model"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/dipperin/dipperin-core/third-party/rpc"
	"github.com/urfave/cli"
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
	l log.Logger

	client RpcClient

	defaultAccount common.Address
	defaultWallet  accounts.WalletIdentifier
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
	l = log.New()
	l.SetHandler(log.MultiHandler(log.CliOutHandler))
}

func InitRpcClient(port int) {
	l.Info("init rpc client", "port", port)
	var err error
	//if client, err = rpc.Dial(fmt.Sprintf("http://%v:%d", "127.0.0.1", port)); err != nil {
	//	panic("init rpc client failed: " + err.Error())
	//}
	wsURL := fmt.Sprintf("ws://%v:%d", "127.0.0.1", port)
	//l.Info("init rpc client", "wsURL", wsURL)
	if client, err = rpc.Dial(wsURL); err != nil {
		panic("init rpc client failed: " + err.Error())
	}
}

func InitAccountInfo(nodeType int, path, password, passPharse string) {
	// get default account
	if nodeType == chain_config.NodeTypeOfNormal {
		err := initWallet(path, password, passPharse)
		if err != nil {
			l.Error("init Wallet Error", "err", err)
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
//		l.Error("can't subscribe server", "err", err)
//		return
//	}
//
//	l.Info("start subscribe node server msg")
//
//	for {
//		select {
//		case err := <-sub.Err():
//			l.Info("sub result", "err", err)
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
		l.Error("call GetSyncStatus", "err", err)
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
				l.Error("call GetSyncStatus", "err", err)
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
		l.Info("RpcCall params assign err, can't find the method")
	}

	rvf := callerRv.MethodByName(method)
	if rvf.Kind() != reflect.Func {
		l.Error("not found method", "method_name", method)
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
	l.Info("the lm is:", "lm", lm)
	return lm
}

func getRpcMethodAndParam(c *cli.Context) (mName string, cParams []string, err error) {
	mName = c.Args().First()
	l.Info("the method name is:", "mName", mName)
	if mName == "" {
		return "", []string{}, errors.New("the method name is nil")
	}
	params := c.String("p")
	l.Info("the params is:", "params", params)

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
		var respBlock rpc_interface.BlockResp
		if err := client.Call(&respBlock, getDipperinRpcMethodByName("CurrentBlock")); err != nil {
			l.Error("get current Block error", "err", err)
			return true
		}
		l.Info("the current block number is:", "number", respBlock.Header.Number)
		return true
	}
	return false
}

func (caller *rpcCaller) GetDefaultAccountBalance(c *cli.Context) {
	var resp rpc_interface.CurBalanceResp
	if err := client.Call(&resp, getDipperinRpcMethodByName("CurrentBalance"), defaultAccount); err != nil {
		l.Error("call current balance error", "err", err)
		return
	}

	l.Info("GetDefaultAccountBalance", "resp", resp)
	balance, err := CSCoinToMoneyValue(resp.Balance)
	if err != nil {
		l.Error("the address isn't on the block chain balance=0")
	} else {
		l.Info("address current Balance is:", "balance", balance+consts.CoinDIPName)
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
			l.Error("the input address is invalid", "err", err)
			return
		}
	}

	var resp rpc_interface.CurBalanceResp
	if err = client.Call(&resp, getDipperinRpcMethodByName(mName), addr); err != nil {

		l.Error("call current balance error", "err", err)
		return
	}

	balance, err := CSCoinToMoneyValue(resp.Balance)
	if err != nil {
		l.Error("the address isn't on the block chain balance=0")
	} else {
		l.Info("address current Balance is:", "balance", balance+consts.CoinDIPName)
	}
}

func printBlockInfo(respBlock rpc_interface.BlockResp) {
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

func printTransactionInfo(respTx rpc_interface.TransactionResp) {
	if respTx.Transaction == nil {
		fmt.Printf("the tx isn't on the block chain\r\n")
		return
	}

	fmt.Printf("\r\n[the tx info is:]")
	fmt.Printf("\r\n%v", respTx.Transaction.String())
	fmt.Printf("\r\nthe blockHash is:%v", respTx.BlockHash.Hex())
	fmt.Printf("\r\nthe BlockNumber is:%v", respTx.BlockNumber)
	fmt.Printf("\r\nthe TxIndex is:%v\r\n", respTx.TxIndex)
}

func (caller *rpcCaller) CurrentBlock(c *cli.Context) {
	mName, _, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error")
		return
	}
	var respBlock rpc_interface.BlockResp
	if err = client.Call(&respBlock, getDipperinRpcMethodByName(mName)); err != nil {
		l.Error("call current block error", "err", err)
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

	var respBlock rpc_interface.BlockResp
	if err = client.Call(&respBlock, getDipperinRpcMethodByName(mName)); err != nil {

		l.Error("call genesis block error", "err", err)
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
	l.Info("the blockNum is:", "blockNum", blockNum)

	var respBlock rpc_interface.BlockResp
	if err = client.Call(&respBlock, getDipperinRpcMethodByName(mName), blockNum); err != nil {
		l.Error("call block error", "err", err)
		return
	}
	printBlockInfo(respBlock)
}

// GetBlockByHash get block information based on block hash
func (caller *rpcCaller) GetBlockByHash(c *cli.Context) {
	mName, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error")
		return
	}

	if len(cParams) == 0 {
		l.Error("parameter includes：blockHash")
		return
	}

	var respBlock rpc_interface.BlockResp
	if err = client.Call(&respBlock, getDipperinRpcMethodByName(mName), cParams[0]); err != nil {
		l.Error("call block error", "err", err)
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
		l.Error("start mining error", "err", err)
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
		l.Error("stop mining error", "err", err)
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
		l.Error("the input address is invalid", "err", err)
		return
	}
	var resp interface{}
	if err = client.Call(&resp, getDipperinRpcMethodByName(mName), address); err != nil {
		l.Error("setting CoinBase", "err", err)
		return
	}
	l.Debug("setting CoinBase　complete")
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
		l.Error("parse gasFloor error", "err", err)
		return
	}

	gasCeil, err := strconv.Atoi(cParams[1])
	if err != nil {
		l.Error("parse gasCeil error", "err", err)
		return
	}

	l.Info("the gasFloor is:", "gasFloor", gasFloor)
	l.Info("the gasCeil is:", "gasCeil", gasCeil)

	var resp interface{}
	if err = client.Call(&resp, getDipperinRpcMethodByName(mName), uint64(gasFloor), uint64(gasCeil)); err != nil {
		l.Error("setting MinerGasConfig failed", "err", err)
		return
	}
	l.Info("setting MinerGasConfig complete")
}

func (caller *rpcCaller) SendTx(c *cli.Context) {
	if checkSync() {
		return
	}

	_, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error", "err", err)
		return
	}

	if len(cParams) < 4 {
		l.Error("parameter includes：to value gasPrice gasLimit")
		return
	}

	toAddress, err := CheckAndChangeHexToAddress(cParams[0])
	if err != nil {
		l.Error("the to address is invalid", "err", err)
		return
	}

	value, err := MoneyValueToCSCoin(cParams[1])
	if err != nil {
		l.Error("the parameter value invalid", "err", err)
		return
	}

	gasPrice, err := MoneyValueToCSCoin(cParams[2])
	if err != nil {
		l.Error("the parameter gasPrice invalid", "err", err)
		return
	}

	gasLimit, err := strconv.ParseUint(cParams[3], 10, 64)
	if err != nil {
		l.Error("the parameter gasLimit invalid", "err", err)
		return
	}

	extraData := make([]byte, 0)
	if len(cParams) >= 4 {
		extraData = []byte(cParams[3])
	}

	var resp common.Hash
	l.Info("the from is: ", "from", defaultAccount.Hex())
	l.Info("the to is: ", "to", toAddress.Hex())
	l.Info("the value is:", "value", cParams[1]+consts.CoinDIPName)
	l.Info("the gasPrice is:", "gasPrice", cParams[2]+consts.CoinDIPName)
	l.Info("the ExtraData is: ", "ExtraData", extraData)
	if err = client.Call(&resp, getDipperinRpcMethodByName("SendTransaction"), defaultAccount, toAddress, value, gasPrice, gasLimit, extraData, nil); err != nil {
		l.Error("call send transaction error", "err", err)
		return
	}
	l.Info("SendTransaction result", "txId", resp.Hex())
}

//send transaction
func (caller *rpcCaller) SendTransaction(c *cli.Context) {
	if checkSync() {
		return
	}

	mName, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error", "err", err)
		return
	}

	if len(cParams) < 6 {
		l.Error("parameter includes：from to value gasPrice gasLimit extraData")
		return
	}

	From, err := CheckAndChangeHexToAddress(cParams[0])
	if err != nil {
		l.Error("the from address is invalid", "err", err)
		l.Error(err.Error())
		return
	}

	To, err := CheckAndChangeHexToAddress(cParams[1])
	if err != nil {
		l.Error("the to address is invalid", "err", err)
		return
	}

	Value, err := MoneyValueToCSCoin(cParams[2])
	if err != nil {
		l.Error("the parameter value invalid")
		return
	}

	gasPrice, err := MoneyValueToCSCoin(cParams[3])
	if err != nil {
		l.Error("the gasPrice value is wrong")
		return
	}

	gasLimit, err := strconv.ParseUint(cParams[4], 10, 64)
	if err != nil {
		l.Error("the gasLimit value is wrong")
		return
	}

	ExtraData := make([]byte, 0)
	if len(cParams) == 6 {
		ExtraData = []byte(cParams[5])
	}

	var resp common.Hash
	l.Info("the From is: ", "From", From.Hex())
	l.Info("the To is: ", "To", To.Hex())
	l.Info("the Value is:", "Value", cParams[2]+consts.CoinDIPName)
	l.Info("the gasPrice is:", "gasPrice", cParams[3]+consts.CoinDIPName)
	l.Info("the gasLimit is:", "gasLimit", gasLimit)
	l.Info("the ExtraData is: ", "ExtraData", ExtraData)
	if err = client.Call(&resp, getDipperinRpcMethodByName(mName), From, To, Value, gasPrice, gasLimit, ExtraData, nil); err != nil {

		l.Error("call send transaction", "err", err)
		return
	}
	l.Info("SendTransaction result", "txId", resp.Hex())
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
		l.Error("the err is:", "err", err)
		l.Error("Transaction decode error")
		return
	}

	var hash common.Hash
	_ = copy(hash[:], tmpHash)

	var resp rpc_interface.TransactionResp
	if err = client.Call(&resp, getDipperinRpcMethodByName(mName), hash); err != nil {
		l.Error("Call Transaction", "err", err)
		return
	}
	printTransactionInfo(resp)
}

func (caller *rpcCaller) GetConvertReceiptByTxHash(c *cli.Context) {
	_, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("GetConvertReceiptByTxHash  getRpcMethodAndParam error")
		return
	}
	if len(cParams) < 1 {
		l.Error("GetConvertReceiptByTxHash need：txHash")
		return
	}
	tmpHash, err := hexutil.Decode(cParams[0])
	if err != nil {
		l.Error("GetConvertReceiptByTxHash decode error")
		return
	}
	var hash common.Hash
	_ = copy(hash[:], tmpHash)

	var resp model.Receipt
	if err = client.Call(&resp, getDipperinRpcMethodByName("GetConvertReceiptByTxHash"), hash); err != nil {
		l.Error("Call GetConvertReceiptByTxHash failed", "err", err)
		return
	}
	l.Info("Call GetConvertReceiptByTxHash")
	fmt.Println(resp.String())
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
	_ = copy(hash[:], tmpHash)

	var resp model.Receipt
	if err = client.Call(&resp, getDipperinRpcMethodByName("GetReceiptByTxHash"), hash); err != nil {
		l.Error("Call GetReceiptByTxHash", "err", err)
		return
	}
	fmt.Println(resp.String())
}

func (caller *rpcCaller) GetLogs(c *cli.Context) {
	//blockHash *common.Hash, fromBlock *big.Int, toBlock *big.Int, addresses []common.Address, topics [][]common.Hash
	params := c.String("p")
	l.Info("GetLogs the params is:", "params", params)
	var filterParams FilterParams
	if err := json.Unmarshal([]byte(params), &filterParams); err != nil {
		l.Error("rpcCaller#GetLogs", "err", err)
	}
	fmt.Println(filterParams)

	var resp []model.Log
	if err := client.Call(&resp, getDipperinRpcMethodByName("GetLogs"), filterParams.blockHash, filterParams.fromBlock, filterParams.toBlock, filterParams.addresses, filterParams.topics); err != nil {
		l.Error("Call GetReceiptByTxHash", "err", err)
		return
	}
	for _, lg := range resp {
		fmt.Println(lg)
	}
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
	}

	var resp model.Receipts
	if err = client.Call(&resp, getDipperinRpcMethodByName("GetReceiptsByBlockNum"), blockNum); err != nil {
		l.Error("Call GetReceiptsByBlockNum", "err", err)
		return
	}
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
	var resp []accounts.WalletIdentifier
	if err = client.Call(&resp, getDipperinRpcMethodByName(mName)); err != nil {
		l.Error("Call ListWallet", "err", err)
		return
	}
	l.Info("Call ListWallet", "resp wallet", resp)
}

//List Wallet Account
func (caller *rpcCaller) ListWalletAccount(c *cli.Context) {
	mName, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error")
		return
	}
	var identifier accounts.WalletIdentifier
	if len(cParams) == 0 {
		identifier = defaultWallet
	} else if len(cParams) == 2 {
		identifier.Path, identifier.WalletName = ParseWalletPathAndName(cParams[1])
		l.Info("ListWalletAccount", "walletPath", identifier.Path, "walletName", identifier.WalletName)
		if cParams[0] == "SoftWallet" {
			identifier.WalletType = accounts.SoftWallet
		} else if cParams[0] == "LedgerWallet" {
			identifier.WalletType = accounts.TrezorWallet
		} else if cParams[0] == "TrezorWallet" {
			identifier.WalletType = accounts.TrezorWallet
		} else {
			l.Error("Wallet Type error")
			return
		}
	} else {
		l.Error("list Wallet Account need：Type Path")
		return
	}

	l.Debug(getDipperinRpcMethodByName(mName))
	var resp []accounts.Account

	if err = client.Call(&resp, getDipperinRpcMethodByName(mName), identifier); err != nil {
		l.Error("Call ListWallet", "err", err)
		return
	}

	l.Info("Call ListWalletAccount", "resp wallet account", resp)
	for _, account := range resp {
		l.Info("the account address is", "address", account.Address.Hex())
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

	var identifier accounts.WalletIdentifier
	identifier.Path, identifier.WalletName = ParseWalletPathAndName(cParams[1])

	if cParams[0] == "SoftWallet" {
		identifier.WalletType = accounts.SoftWallet
	} else if cParams[0] == "LedgerWallet" {
		identifier.WalletType = accounts.TrezorWallet
	} else if cParams[0] == "TrezorWallet" {
		identifier.WalletType = accounts.TrezorWallet
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
		l.Error("Call EstablishWallet", "err", err)
		return
	}
	resp = strings.Replace(resp, " ", ",", -1)
	l.Info("Call EstablishWallet", "resp mnemonic", resp)
}

//Restore Wallet
func (caller *rpcCaller) RestoreWallet(c *cli.Context) {
	mName, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error")
		return
	}
	if len(cParams) < 4 {
		l.Error("RestoreWallet need：Type Path password mnemonic")
		return
	}

	var identifier accounts.WalletIdentifier
	identifier.Path, identifier.WalletName = ParseWalletPathAndName(cParams[1])

	if cParams[0] == "SoftWallet" {
		identifier.WalletType = accounts.SoftWallet
	} else if cParams[0] == "LedgerWallet" {
		identifier.WalletType = accounts.TrezorWallet
	} else if cParams[0] == "TrezorWallet" {
		identifier.WalletType = accounts.TrezorWallet
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
	l.Info("the identifier is: ", "identifier", identifier)
	l.Info("the password is: ", "password", password)
	//l.Info("the passPhrase is: ", "passPhrase", passPhrase)
	l.Info("the mnemonic is: ", "mnemonic", mnemonic)

	if err = client.Call(resp, getDipperinRpcMethodByName(mName), password, mnemonic, "", identifier); err != nil {
		l.Error("Call RestoreWallet", "err", err)
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

	var identifier accounts.WalletIdentifier
	var password string
	if len(cParams) == 1 {
		identifier = defaultWallet
		password = cParams[0]
	} else if len(cParams) == 3 {
		identifier.Path, identifier.WalletName = ParseWalletPathAndName(cParams[1])
		if cParams[0] == "SoftWallet" {
			identifier.WalletType = accounts.SoftWallet
		} else if cParams[0] == "LedgerWallet" {
			identifier.WalletType = accounts.LedgerWallet
		} else if cParams[0] == "TrezorWallet" {
			identifier.WalletType = accounts.TrezorWallet
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
		l.Error("Call OpenWallet", "err", err)
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

	var identifier accounts.WalletIdentifier
	if len(cParams) == 0 {
		identifier = defaultWallet
	} else if len(cParams) == 2 {
		identifier.Path, identifier.WalletName = ParseWalletPathAndName(cParams[1])
		if cParams[0] == "SoftWallet" {
			identifier.WalletType = accounts.SoftWallet
		} else if cParams[0] == "LedgerWallet" {
			identifier.WalletType = accounts.TrezorWallet
		} else if cParams[0] == "TrezorWallet" {
			identifier.WalletType = accounts.TrezorWallet
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
		l.Error("Call CloseWallet", "err", err)
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

	var identifier accounts.WalletIdentifier
	if len(cParams) == 0 {
		identifier = defaultWallet
	} else if len(cParams) == 2 {
		identifier.Path, identifier.WalletName = ParseWalletPathAndName(cParams[1])
		if cParams[0] == "SoftWallet" {
			identifier.WalletType = accounts.SoftWallet
		} else if cParams[0] == "LedgerWallet" {
			identifier.WalletType = accounts.TrezorWallet
		} else if cParams[0] == "TrezorWallet" {
			identifier.WalletType = accounts.TrezorWallet
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
	var resp accounts.Account
	if err = client.Call(&resp, getDipperinRpcMethodByName(mName), derivationPath, identifier); err != nil {
		l.Error("Call AddAccount", "err", err)
		return
	}
	l.Info("Call AddAccount", "resp AddAccount", resp.Address.Hex())
}

func (caller *rpcCaller) SendRegisterTx(c *cli.Context) {
	if checkSync() {
		return
	}

	_, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error", "err", err)
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
		l.Error("the parameter gasPrice invalid", "err", err)
		return
	}

	gasLimit, err := strconv.ParseUint(cParams[2], 10, 64)
	if err != nil {
		l.Error("the parameter gasLimit invalid", "err", err)
		return
	}
	var resp common.Hash
	if err = client.Call(&resp, getDipperinRpcMethodByName("SendRegisterTransaction"), defaultAccount, stake, gasPrice, gasLimit, nil); err != nil {

		l.Error("call send transaction", "err", err)
		return
	}
	l.Info("SendRegisterTransaction result", "txId", resp.Hex())
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
	From, err := CheckAndChangeHexToAddress(cParams[0])
	if err != nil {
		l.Error("the from address is invalid", "err", err)
		return
	}

	stake, err := MoneyValueToCSCoin(cParams[1])
	if err != nil {
		l.Error("the parameter stake invalid")
		return
	}

	gasPrice, err := MoneyValueToCSCoin(cParams[2])
	if err != nil {
		l.Error("the parameter gasPrice invalid", "err", err)
		return
	}

	gasLimit, err := strconv.ParseUint(cParams[3], 10, 64)
	if err != nil {
		l.Error("the parameter gasLimit invalid", "err", err)
		return
	}

	if err = client.Call(&resp, getDipperinRpcMethodByName(mName), From, stake, gasPrice, gasLimit, nil); err != nil {

		l.Error("call send transaction", "err", err)
		return
	}
	l.Info("SendRegisterTransaction result", "txId", resp.Hex())
	addTrackingAccount(From)
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
		l.Error("the parameter gasPrice invalid", "err", err)
		return
	}

	gasLimit, err := strconv.ParseUint(cParams[1], 10, 64)
	if err != nil {
		l.Error("the parameter gasLimit invalid", "err", err)
		return
	}

	if err = client.Call(&resp, getDipperinRpcMethodByName("SendUnStakeTransaction"), defaultAccount, gasPrice, gasLimit, nil); err != nil {
		l.Error("call send transaction", "err", err)
		return
	}
	l.Info("SendUnStakeTransaction result", "txId", resp.Hex())
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
	From, err := CheckAndChangeHexToAddress(cParams[0])
	if err != nil {
		l.Error("the from address is invalid", "err", err)
		return
	}

	gasPrice, err := MoneyValueToCSCoin(cParams[1])
	if err != nil {
		l.Error("the parameter gasPrice invalid", "err", err)
		return
	}

	gasLimit, err := strconv.ParseUint(cParams[2], 10, 64)
	if err != nil {
		l.Error("the parameter gasLimit invalid", "err", err)
		return
	}

	if err = client.Call(&resp, getDipperinRpcMethodByName(mName), From, gasPrice, gasLimit, nil); err != nil {
		l.Error("call send transaction", "err", err)
		return
	}
	l.Info("SendUnStakeTransaction result", "txId", resp.Hex())
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
		l.Error("the parameter gasPrice invalid", "err", err)
		return
	}

	gasLimit, err := strconv.ParseUint(cParams[1], 10, 64)
	if err != nil {
		l.Error("the parameter gasLimit invalid", "err", err)
		return
	}

	if err = client.Call(&resp, getDipperinRpcMethodByName("SendCancelTransaction"), defaultAccount, gasPrice, gasLimit, nil); err != nil {
		l.Error("call send transaction", "err", err)
		return
	}
	l.Info("SendCancelTransaction result", "txId", resp.Hex())
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
	From, err := CheckAndChangeHexToAddress(cParams[0])
	if err != nil {
		l.Error("the from address is invalid", "err", err)
		return
	}

	gasPrice, err := MoneyValueToCSCoin(cParams[1])
	if err != nil {
		l.Error("the parameter gasPrice invalid", "err", err)
		return
	}

	gasLimit, err := strconv.ParseUint(cParams[2], 10, 64)
	if err != nil {
		l.Error("the parameter gasLimit invalid", "err", err)
		return
	}

	if err = client.Call(&resp, getDipperinRpcMethodByName(mName), From, gasPrice, gasLimit, nil); err != nil {
		l.Error("call send transaction", "err", err)
		return
	}
	l.Info("SendCancelTransaction result", "txId", resp.Hex())
	removeTrackingAccount(From)
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
		l.Error("GetVerifiersBySlot", "err", err)
		return
	}
	l.Info("GetVerifiersBySlot result")
	for _, tmpAddress := range resp {
		l.Info("verifier address is:", "verifier", tmpAddress.Hex())
	}
}

// VerifierStatus call to get the verifier status
func (caller *rpcCaller) VerifierStatus(c *cli.Context) {
	//params is a comma separated list of addresses
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
			l.Error("the input address is invalid", "err", err)
			return
		}
	}

	var resp rpc_interface.VerifierStatus
	l.Debug(getDipperinRpcMethodByName(mName))

	if err = client.Call(&resp, getDipperinRpcMethodByName(mName), addr); err != nil {
		l.Error("call verifier's status error", "err", err)
		return
	}

	balance, err := CSCoinToMoneyValue(resp.Balance)
	if err != nil {
		l.Error("The address has no balance, balance = 0 DIP")
		balance = "0"
	}

	if resp.Status == VerifierStatusNoRegistered {
		l.Info("Verifier status", "status", resp.Status, "balance", balance+" DIP")
		return
	}

	if resp.Status == VerifiedStatusUnstaked {
		l.Info("Verifier status", "status", resp.Status, "balance", balance+" DIP")
		return
	}

	stake, err := CSCoinToMoneyValue(resp.Stake)
	if err != nil {
		l.Error("The address has no stake, stake = 0 DIP")
	}
	if resp.Status == VerifierStatusRegistered {
		l.Info("Verifier status", "status", resp.Status, "balance", balance+" DIP", "stake", stake+" DIP", "reputation", resp.Reputation, "is current verifier", resp.IsCurrentVerifier)
	}

	if resp.Status == VerifiedStatusCanceled {
		l.Info("Verifier status", "status", resp.Status, "balance", balance+" DIP", "stake", stake+" DIP", "reputation", resp.Reputation)
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
		l.Error("the from address is invalid", "err", err)
		return
	}

	l.Debug(getDipperinRpcMethodByName(mName))
	var resp interface{}
	if err = client.Call(&resp, getDipperinRpcMethodByName(mName), addr); err != nil {
		l.Info("Set wallet signer（default account）", "err", err)
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
	//params is a comma separated list of addresses
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
			l.Error("the input address is invalid", "err", err)
			return
		}
	}

	var resp rpc_interface.CurBalanceResp
	l.Debug(getDipperinRpcMethodByName(mName))

	if err = client.Call(&resp, getDipperinRpcMethodByName(mName), addr); err != nil {
		l.Error("lookup current stake", "err", err)
		return
	}

	stake, err := CSCoinToMoneyValue(resp.Balance)
	if err != nil {
		l.Error("the address isn't on the block chain stake=0")
	} else {
		l.Info("address current stake is:", "stake", stake+consts.CoinDIPName)
	}
}

// CurrentReputation call the method to get the current reputation value
func (caller *rpcCaller) CurrentReputation(c *cli.Context) {
	//params is a comma separated list of addresses
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
			l.Error("the input address is invalid", "err", err)
			return
		}
	}

	var resp uint64
	if err = client.Call(&resp, getDipperinRpcMethodByName(mName), addr); err != nil {
		l.Error("lookup current reputation error", "err", err)
		return
	}
	l.Info("address current reputation is:", "reputation", resp)
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
		l.Error("call failed", "err", err)
		return
	}

	for _, a := range resp {
		l.Info("current verifier", "address", a.Hex(), "is_default", inDefaultVs(a))
	}
}

func inDefaultVs(a common.Address) bool {
	for _, v := range chain.VerifierAddress {
		if a.IsEqual(v) {
			return true
		}
	}
	return false
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
		l.Error("call failed", "err", err)
		return
	}
	//l.Info("get next verifiers3", "resp", resp)
	for _, a := range resp {
		l.Info("next verifier", "address", " "+a.Hex())
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
			l.Error("the input address is invalid", "err", err)
			return
		}
	}

	var resp uint64
	if err = client.Call(&resp, getDipperinRpcMethodByName(mName), addr); err != nil {
		l.Error("call GetTransactionNonce", "err", err)
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
		l.Error("GetTransactionNonce error", "err", err)
		return
	}
	l.Info("the address nonce from chain is:", "nonce", nonce)
}

func (caller *rpcCaller) GetAddressNonceFromWallet(c *cli.Context) {
	nonce, err := getNonceInfo(c)
	if err != nil {
		l.Error("GetTransactionNonce error", "err", err)
		return
	}

	l.Info("the address nonce from wallet is:", "nonce", nonce)
}

func initWallet(path, password, passPhrase string) (err error) {
	var identifier accounts.WalletIdentifier
	identifier.WalletType = accounts.SoftWallet
	identifier.Path, identifier.WalletName = ParseWalletPathAndName(path)

	//open
	exit, _ := soft_wallet.PathExists(identifier.Path)
	if exit {
		l.Info("open wallet", "identifier", identifier)
		var resp interface{}
		if err = client.Call(&resp, getDipperinRpcMethodByName("OpenWallet"), password, identifier); err != nil {
			l.Error("open Wallet err", "err", err)
			return err
		}
	} else {
		l.Info("establish wallet", "identifier", identifier)
		var mnemonic string
		if err = client.Call(&mnemonic, getDipperinRpcMethodByName("EstablishWallet"), password, passPhrase, identifier); err != nil {
			l.Error("Call EstablishWallet", "err", err)
			return err
		}
		mnemonic = strings.Replace(mnemonic, " ", ",", -1)
		l.Info("EstablishWallet mnemonic is:", "mnemonic", mnemonic)
	}
	return nil
}

func getDefaultAccount() common.Address {
	var resp []accounts.WalletIdentifier
	l.Info("getDefaultAccount")

	if err := client.Call(&resp, getDipperinRpcMethodByName("ListWallet")); err != nil {
		l.Error("Call ListWallet", "err", err)
		return common.Address{}
	}

	var respA []accounts.Account
	if err := client.Call(&respA, getDipperinRpcMethodByName("ListWalletAccount"), resp[0]); err != nil {
		l.Error("Call ListWallet", "err", err)
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
		l.Error("the input address is invalid", "err", err)
	}

	defaultAccount = addr
	l.Info("default account set success", "account", defaultAccount.Hex())
}*/

func getDefaultWallet() accounts.WalletIdentifier {
	var resp []accounts.WalletIdentifier
	if err := client.Call(&resp, getDipperinRpcMethodByName("ListWallet")); err != nil {
		l.Error("Call ListWallet", "err", err)
		return accounts.WalletIdentifier{}
	}
	return resp[0]
}

//if user applies for registering verifier, record,
// creating a file in $Home/.dipperin
// for startup check, quit if node type is not verifier
func RecordRegistration(txHash string) {
	confPath := filepath.Join(util.HomeDir(), ".dipperin", "registration")
	exist, _ := soft_wallet.PathExists(confPath)
	if !exist {
		if err := os.MkdirAll(filepath.Dir(confPath), 0766); err != nil {
			l.Error("can't make dir", "err", err)
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
	exist, _ := soft_wallet.PathExists(confPath)
	if !exist {
		return
	}

	if err := os.Remove(confPath); err != nil {
		l.Error("can't remove record registration")
		return
	}
}

func CheckRegistration() bool {
	confPath := filepath.Join(util.HomeDir(), ".dipperin", "registration")
	exist, _ := soft_wallet.PathExists(confPath)
	return exist
}
