package commands

import (
	"errors"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/hexutil"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/urfave/cli"
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
		l.Error("Call GetContractAddressByTxHash failed", "err", err)
		return
	}
	l.Info("Call GetContractAddressByTxHash", "Contract Address", resp)
}

func (caller *rpcCaller) CallContract(c *cli.Context) {
	mName, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error", "err", err)
		return
	}
	if len(cParams) != 2 && len(cParams) != 3 {
		l.Error("parameter includes：from to blockNum, blockNum is optional")
		return
	}
	from, err := CheckAndChangeHexToAddress(cParams[0])
	if err != nil {
		l.Error("call the from address is invalid", "err", err)
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

	input := getRpcSpecialParam(c, "input")
	funcName, err := getCalledFuncName(c)
	if err != nil {
		l.Error(err.Error())
		return
	}
	// RLP([funcName][param1,param2,param3...])
	inputRlp, err := rlp.EncodeToBytes([]interface{}{funcName, input})
	if err != nil {
		log.Error("input rlp err")
		return
	}

	l.Info("the from is: ", "from", from.Hex())
	l.Info("the to is: ", "to", to.Hex())
	l.Info("the blockNum is: ", "num", blockNum)
	l.Info("the funcName is:", "funcName", funcName)
	//l.Info("the ExtraData is: ", "ExtraData", inputRlp)

	var resp string
	if err = client.Call(&resp, getDipperinRpcMethodByName(mName), from, to, inputRlp, blockNum); err != nil {
		l.Error("CallContract failed", "err", err)
		return
	}
	l.Info("CallContract", "resp", resp)
}

func (caller *rpcCaller) EstimateGas(c *cli.Context) {
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
	l.Info("SendTransactionContract successful", "resp", reflect.ValueOf(resp).String())
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

	ExtraData, err := getCreateExtraData(c)
	if err != nil {
		return
	}

	l.Info("the from is: ", "from", from.Hex())
	l.Info("the value is:", "value", MoneyWithUnit(cParams[1]))
	l.Info("the gasPrice is:", "gasPrice", MoneyWithUnit(cParams[2]))
	l.Info("the gasLimit is:", "gasLimit", gasLimit)
	//l.Info("the ExtraData is: ", "ExtraData", ExtraData)
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
		l.Error("the parameter value invalid", "err", err)
		return
	}

	gasPrice, err := MoneyValueToCSCoin(cParams[3])
	if err != nil {
		l.Error("the parameter gasPrice invalid", "err", err)
		return
	}

	gasLimit, err := strconv.ParseUint(cParams[4], 10, 64)
	if err != nil {
		l.Error("the parameter gasLimit invalid", "err", err)
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
		log.Error("input rlp err")
		return
	}

	l.Info("the from is: ", "from", from.Hex())
	l.Info("the to is: ", "to", to.Hex())
	l.Info("the value is:", "value", MoneyWithUnit(cParams[2]))
	l.Info("the gasPrice is:", "gasPrice", MoneyWithUnit(cParams[3]))
	l.Info("the gasLimit is:", "gasLimit", gasLimit)
	l.Info("the funcName is:", "funcName", funcName)
	//l.Info("the ExtraData is: ", "ExtraData", inputRlp)

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
