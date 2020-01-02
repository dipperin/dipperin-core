package commands

import (
	"errors"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/hexutil"
	"github.com/dipperin/dipperin-core/common/log"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/urfave/cli"
	"go.uber.org/zap"
	"io/ioutil"
	"reflect"
	"strconv"
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
		l.Error("Call GetContractAddressByTxHash failed", zap.Error(err))
		return
	}
	l.Info("Call GetContractAddressByTxHash", zap.Any("Contract Address", resp))
}

func (caller *rpcCaller) CallContract(c *cli.Context) {
	mName, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error", zap.Error(err))
		return
	}
	if len(cParams) != 2 && len(cParams) != 3 {
		l.Error("parameter includes：from to blockNum, blockNum is optional")
		return
	}
	from, err := CheckAndChangeHexToAddress(cParams[0])
	if err != nil {
		l.Error("call the from address is invalid", zap.Error(err))
		return
	}
	to, err := CheckAndChangeHexToAddress(cParams[1])
	if err != nil {
		l.Error("call the to address is invalid", zap.Error(err))
		return
	}

	var blockNum uint64
	if len(cParams) == 3 {
		blockNum, err = strconv.ParseUint(cParams[2], 10, 64)
		if err != nil {
			l.Error("ParseUint failed", zap.Error(err))
			return
		}
	}

	input := getRpcSpecialParam(c, "input")
	funcName, err := getCalledFuncName(c)
	if err != nil {
		l.Error(err.Error())
		return
	}
	// RLP([funcName][param1,param2,param3...])
	inputRlp, err := rlp.EncodeToBytes([]interface{}{funcName, input})
	if err != nil {
		log.DLogger.Error("input rlp err")
		return
	}

	l.Debug("the from is: ", zap.String("from", from.Hex()))
	l.Debug("the to is: ", zap.String("to", to.Hex()))
	l.Debug("the blockNum is: ", zap.Uint64("num", blockNum))
	l.Debug("the funcName is:", zap.String("funcName", funcName))
	//l.Debug("the ExtraData is: ", "ExtraData", inputRlp)

	var resp string
	if err = client.Call(&resp, getDipperinRpcMethodByName(mName), from, to, inputRlp, blockNum); err != nil {
		l.Error("CallContract failed", zap.Error(err))
		return
	}
	l.Info("CallContract", zap.String("resp", resp))
}

func (caller *rpcCaller) EstimateGas(c *cli.Context) {
	var (
		resp interface{}
		err  error
	)

	if isCreate(c) {
		resp, err = contractCreate(c)
		if err != nil {
			l.Error("EstimateGas failed", zap.Error(err))
			return
		}
	} else {
		resp, err = contractCall(c)
		if err != nil {
			l.Error("EstimateGas failed", zap.Error(err))
			return
		}
	}

	value, err := hexutil.DecodeUint64(reflect.ValueOf(resp).String())
	if err != nil {
		l.Error("EstimateGas failed", zap.Error(err))
	}
	l.Info(" EstimateGas successful", zap.Uint64("estimated gas", value))
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
			l.Error("SendTransactionContract failed", zap.Error(err))
			return
		}
	} else {
		resp, err = contractCall(c)
		if err != nil {
			l.Error("SendTransactionContract failed", zap.Error(err))
			return
		}
	}
	l.Info("SendTransactionContract successful", zap.String("txId", reflect.ValueOf(resp).String()))
}

func contractCreate(c *cli.Context) (resp interface{}, err error) {
	mName, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		return
	}
	if len(cParams) != 4 {
		err = errors.New("parameter includes：from value gasPrice gasLimit")
		return
	}

	from, err := CheckAndChangeHexToAddress(cParams[0])
	if err != nil {
		return
	}

	to := common.HexToAddress(common.AddressContractCreate)
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

	ExtraData, err := getCreateExtraData(c)
	if err != nil {
		return
	}

	l.Debug("the from is: ", zap.String("from", from.Hex()))
	l.Debug("the value is:", zap.String("value", MoneyWithUnit(cParams[1])))
	l.Debug("the gasPrice is:", zap.String("gasPrice", MoneyWithUnit(cParams[2])))
	l.Debug("the gasLimit is:", zap.Uint64("gasLimit", gasLimit))
	//l.Debug("the ExtraData is: ", "ExtraData", ExtraData)
	err = client.Call(&resp, getDipperinRpcMethodByName(mName), from, to, value, gasPrice, gasLimit, ExtraData, nil)
	return
}

func contractCall(c *cli.Context) (resp interface{}, err error) {
	mName, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		return
	}

	if len(cParams) != 5 {
		err = errors.New("parameter includes：from to value gasPrice gasLimit")
		return
	}

	from, err := CheckAndChangeHexToAddress(cParams[0])
	if err != nil {
		return
	}

	to, err := CheckAndChangeHexToAddress(cParams[1])
	if err != nil {
		return
	}

	value, err := MoneyValueToCSCoin(cParams[2])
	if err != nil {
		l.Error("the parameter value invalid", zap.Error(err))
		return
	}

	gasPrice, err := MoneyValueToCSCoin(cParams[3])
	if err != nil {
		l.Error("the parameter gasPrice invalid", zap.Error(err))
		return
	}

	gasLimit, err := strconv.ParseUint(cParams[4], 10, 64)
	if err != nil {
		l.Error("the parameter gasLimit invalid", zap.Error(err))
		return
	}

	input := getRpcSpecialParam(c, "input")
	funcName, err := getCalledFuncName(c)
	if err != nil {
		return
	}

	// RLP([funcName][param1,param2,param3...])
	inputRlp, err := rlp.EncodeToBytes([]interface{}{funcName, input})
	if err != nil {
		log.DLogger.Error("input rlp err")
		return
	}

	l.Debug("the from is: ", zap.String("from", from.Hex()))
	l.Debug("the to is: ", zap.String("to", to.Hex()))
	l.Debug("the value is:", zap.String("value", MoneyWithUnit(cParams[2])))
	l.Debug("the gasPrice is:", zap.String("gasPrice", MoneyWithUnit(cParams[3])))
	l.Debug("the gasLimit is:", zap.Uint64("gasLimit", gasLimit))
	l.Debug("the funcName is:", zap.String("funcName", funcName))
	//l.Debug("the ExtraData is: ", "ExtraData", inputRlp)

	err = client.Call(&resp, getDipperinRpcMethodByName(mName), from, to, value, gasPrice, gasLimit, inputRlp, nil)
	return
}

func getCalledFuncName(c *cli.Context) (funcName string, err error) {
	funcName = c.String("func-name")
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
