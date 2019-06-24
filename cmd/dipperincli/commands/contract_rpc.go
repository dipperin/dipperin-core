package commands

import (
	"io/ioutil"
	"strconv"
	"github.com/urfave/cli"
	"math/big"
	"github.com/dipperin/dipperin-core/common/consts"
	"github.com/dipperin/dipperin-core/common/hexutil"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/config"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/dipperin/dipperin-core/third-party/log"
	"errors"
	"reflect"
)

func (caller *rpcCaller) GetContractAddressByTxHash(c *cli.Context) {
	_, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("GetContractAddressByTxHash  getRpcMethodAndParam error")
		return
	}
	if len(cParams) < 1 {
		l.Error("GetContractAddressByTxHash need：txHash")
		return
	}
	tmpHash, err := hexutil.Decode(cParams[0])
	if err != nil {
		l.Error("GetContractAddressByTxHash decode error")
		return
	}
	var hash common.Hash
	_ = copy(hash[:], tmpHash)

	var resp common.Address
	if err = client.Call(&resp, getDipperinRpcMethodByName("GetContractAddressByTxHash"), hash); err != nil {
		l.Error("Call GetContractAddressByTxHash failed", "err", err)
		return
	}
	l.Info("Call GetContractAddressByTxHash", "Contract Address", resp)
}

func (caller *rpcCaller) CallContract(c *cli.Context) {
	if checkSync() {
		return
	}
	mName, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error", "err", err)
		return
	}
	if len(cParams) != 2 && len(cParams) != 3 {
		l.Error("parameter includes：from to blockNum, blockNum is optional")
		return
	}
	From, err := CheckAndChangeHexToAddress(cParams[0])
	if err != nil {
		l.Error("call  the from address is invalid", "err", err)
		return
	}
	to, err := CheckAndChangeHexToAddress(cParams[1])
	if err != nil {
		l.Error("call the to address is invalid", "err", err)
		return
	}

	var blockNum uint64
	if len(cParams) == 3 {
		blockNum, err = strconv.ParseUint(cParams[2], 10, 64)
		if err != nil {
			l.Error("ParseUint failed", "err", err)
			return
		}
	}

	funcName, err := getCalledFuncName(c)
	if err != nil {
		l.Error(err.Error())
		return
	}
	input := getRpcSpecialParam(c, "input")
	// RLP([funcName][params])
	inputRlp, err := rlp.EncodeToBytes([]interface{}{
		funcName, input,
	})
	if err != nil {
		log.Error("input rlp err")
		return
	}

	l.Info("the From is: ", "From", From.Hex())
	l.Info("the To is: ", "To", to.Hex())
	l.Info("the BlockNum is: ", "Num", blockNum)
	l.Info("the funcName is:", "funcName", funcName)
	l.Info("the ExtraData is: ", "ExtraData", inputRlp)

	var resp string
	if err = client.Call(&resp, getDipperinRpcMethodByName(mName), From, to, inputRlp, blockNum); err != nil {
		l.Error("CallContract failed", "err", err)
		return
	}
	l.Info(" CallContract", "resp", resp)
}

func (caller *rpcCaller) EstimateGas(c *cli.Context) {
	if checkSync() {
		return
	}

	var (
		resp interface{}
		err  error
	)

	if isCreate(c) {
		resp, err = contractCreate(c)
		if err != nil {
			l.Error("EstimateGas failed", "err", err.Error())
			return
		}
	} else {
		resp, err = contractCall(c)
		if err != nil {
			l.Error("EstimateGas failed", "err", err.Error())
			return
		}
	}

	value, err := hexutil.DecodeUint64(reflect.ValueOf(resp).String())
	if err != nil {
		l.Error("EstimateGas failed", "err", err.Error())
	}
	l.Info(" EstimateGas successful", "resp", value)
}

func (caller *rpcCaller) SendTransactionContract(c *cli.Context) {
	if checkSync() {
		return
	}

	var (
		resp interface{}
		err  error
	)

	if isCreate(c) {
		resp, err = contractCreate(c)
		if err != nil {
			l.Error("SendTransactionContract failed", "err", err.Error())
			return
		}
	} else {
		resp, err = contractCall(c)
		if err != nil {
			l.Error("SendTransactionContract failed", "err", err.Error())
			return
		}
	}
	l.Info(" SendTransactionContract successful", "resp", reflect.ValueOf(resp).String())
}

func contractCreate(c *cli.Context) (resp interface{}, err error) {
	mName, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		return
	}
	if len(cParams) != 3 && len(cParams) != 4 {
		err = errors.New("parameter includes：from value gasLimit gasPrice, gasPrice is optional")
		return
	}

	From, err := CheckAndChangeHexToAddress(cParams[0])
	if err != nil {
		return
	}

	to := common.HexToAddress(common.AddressContractCreate)
	Value, err := MoneyValueToCSCoin(cParams[1])
	if err != nil {
		return
	}
	gasLimit, err := MoneyValueToCSCoin(cParams[2])
	if err != nil {
		return
	}

	var gasPrice *big.Int
	if len(cParams) == 3 {
		gasPrice.SetInt64(config.DEFAULT_GAS_PRICE)
	} else {
		var gasPriceVal int64
		gasPriceVal, err = strconv.ParseInt(cParams[3], 10, 64)
		if err != nil {
			return
		}
		gasPrice = new(big.Int).SetInt64(gasPriceVal)
	}

	ExtraData, err := getCreateExtraData(c)
	if err != nil {
		return
	}

	l.Info("the From is: ", "From", From.Hex())
	l.Info("the Value is:", "Value", cParams[1]+consts.CoinDIPName)
	l.Info("the gasLimit is:", "gasLimit", gasLimit)
	l.Info("the gasPrice is:", "gasPrice", gasPrice)
	l.Info("the ExtraData is: ", "ExtraData", ExtraData)

	err = client.Call(&resp, getDipperinRpcMethodByName(mName), From, to, Value, gasLimit, gasPrice, ExtraData, nil)
	return
}

func contractCall(c *cli.Context) (resp interface{}, err error) {
	mName, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		return
	}
	if len(cParams) != 3 && len(cParams) != 4 {
		err = errors.New("parameter includes：from to gasLimit gasPrice, gasPrice is optional")
		return
	}
	From, err := CheckAndChangeHexToAddress(cParams[0])
	if err != nil {
		return
	}
	to, err := CheckAndChangeHexToAddress(cParams[1])
	if err != nil {
		return
	}

	gasLimit, err := MoneyValueToCSCoin(cParams[2])
	if err != nil {
		return
	}

	var gasPrice *big.Int
	if len(cParams) == 2 {
		gasPrice.SetInt64(config.DEFAULT_GAS_PRICE)
	} else {
		var gasPriceVal int64
		gasPriceVal, err = strconv.ParseInt(cParams[3], 10, 64)
		if err != nil {
			return
		}
		gasPrice = new(big.Int).SetInt64(gasPriceVal)
	}

	funcName, err := getCalledFuncName(c)
	if err != nil {
		return
	}

	input := getRpcSpecialParam(c, "input")
	// RLP([funcName][params])
	inputRlp, err := rlp.EncodeToBytes([]interface{}{
		funcName, input,
	})
	if err != nil {
		return
	}

	l.Info("the From is: ", "From", From.Hex())
	l.Info("the gasLimit is:", "gasLimit", gasLimit)
	l.Info("the gasPrice is:", "gasPrice", gasPrice)
	l.Info("the funcName is:", "funcName", funcName)
	l.Info("the ExtraData is: ", "ExtraData", inputRlp)

	//SendTransactionContract(from, to common.Address,value,gasLimit, gasPrice *big.Int, data []byte, nonce *uint64 )
	err = client.Call(&resp, getDipperinRpcMethodByName(mName), From, to, nil, gasLimit, gasPrice, inputRlp, nil)
	return
}

func getCalledFuncName(c *cli.Context) (funcName string, err error) {
	funcName = c.String("funcName")
	if funcName == "" {
		return "", errors.New("function name is need")
	}
	return funcName, nil
}

func isCreate(c *cli.Context) bool {
	return c.Bool("is-create")
}

func getCreateExtraData(c *cli.Context) (ExtraData []byte, err error) {
	// Get wasm
	wasmPath, err := getRpcParamValue(c, "wasm")
	if err != nil {
		return nil, errors.New("the wasm path value invalid")
	}
	wasmBytes, err := ioutil.ReadFile(wasmPath)
	if err != nil {
		return nil, errors.New("the wasm file read err")
	}

	// Get abi
	abiPath, err := getRpcParamValue(c, "abi")
	if err != nil {
		return nil, errors.New("the abi path value invalid")
	}
	abiBytes, err := ioutil.ReadFile(abiPath)
	if err != nil {
		return nil, errors.New("the abi file read err")
	}

	input := getRpcSpecialParam(c, "input")
	rlpParams := []interface{}{
		wasmBytes, abiBytes, input,
	}
	return rlp.EncodeToBytes(rlpParams)
}
