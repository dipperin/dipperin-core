package commands

import (
	"strings"
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
	"github.com/dipperin/dipperin-core/core/vm/common/utils"
	"fmt"
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

	var resp interface{}
	if isCreate(c) {
		resp = contractCreate(c)
	} else {
		resp = contractCall(c)
	}

	value, innerErr := hexutil.DecodeUint64(reflect.ValueOf(resp).String())
	if innerErr != nil {
		l.Info(" EstimateGas decode failed","err", innerErr)
	}
	l.Info(" EstimateGas", "resp", value)
}

func (caller *rpcCaller) SendTransactionContract(c *cli.Context) {
	if checkSync() {
		return
	}

	var resp interface{}
	if isCreate(c) {
		resp = contractCreate(c)
	} else {
		resp = contractCall(c)
	}
	l.Info(" SendTransactionContract", "resp", reflect.ValueOf(resp).String())
}

func contractCreate(c *cli.Context) interface{} {
	mName, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error", "err", err)
		return nil
	}
	if len(cParams) != 3 && len(cParams) != 4 {
		l.Error("parameter includes：from value gasLimit gasPrice, gasPrice is optional")
		return nil
	}

	From, err := CheckAndChangeHexToAddress(cParams[0])
	if err != nil {
		l.Error("the from address is invalid", "err", err)
		l.Error(err.Error())
		return nil
	}
	to := common.HexToAddress(common.AddressContractCreate)

	Value, err := MoneyValueToCSCoin(cParams[1])
	if err != nil {
		l.Error("the parameter value invalid")
		return nil
	}
	gasLimit, err := MoneyValueToCSCoin(cParams[2])
	if err != nil {
		l.Error("the parameter value invalid")
		return nil
	}

	var gasPrice *big.Int
	if len(cParams) == 3 {
		gasPrice.SetInt64(config.DEFAULT_GAS_PRICE)
	} else {
		gasPriceVal, innerErr := strconv.ParseInt(cParams[3], 10, 64)
		if innerErr != nil {
			l.Error("the parameter value invalid")
			return nil
		}
		gasPrice = new(big.Int).SetInt64(gasPriceVal)
	}

	ExtraData, err := generateExtraData(c)
	if err != nil {
		l.Error("generate extraData err", "err", err)
		return nil
	}

	var resp interface{}
	l.Info("the From is: ", "From", From.Hex())
	l.Info("the Value is:", "Value", cParams[1]+consts.CoinDIPName)
	l.Info("the gasLimit is:", "gasLimit", gasLimit)
	l.Info("the gasPrice is:", "gasPrice", gasPrice)
	l.Info("the ExtraData is: ", "ExtraData", ExtraData)
	if err = client.Call(&resp, getDipperinRpcMethodByName(mName), From, to, Value, gasLimit, gasPrice, ExtraData, nil); err != nil {
		l.Error(fmt.Sprintf("%s Create failed", mName), "err", err)
		return nil
	}
	return resp
}

func contractCall(c *cli.Context) interface{} {
	mName, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error", "err", err)
		return nil
	}
	if len(cParams) != 3 && len(cParams) != 4 {
		l.Error("parameter includes：from to gasLimit gasPrice, gasPrice is optional")
		return nil
	}
	From, err := CheckAndChangeHexToAddress(cParams[0])
	if err != nil {
		l.Error("call  the from address is invalid", "err", err)
		return nil
	}
	to, err := CheckAndChangeHexToAddress(cParams[1])
	if err != nil {
		l.Error("call the to address is invalid", "err", err)
		return nil
	}

	gasLimit, err := MoneyValueToCSCoin(cParams[2])
	if err != nil {
		l.Error("call the parameter value invalid", "err", err)
		return nil
	}

	var gasPrice *big.Int
	if len(cParams) == 2 {
		gasPrice.SetInt64(config.DEFAULT_GAS_PRICE)
	} else {
		gasPriceVal, innerErr := strconv.ParseInt(cParams[3], 10, 64)
		if innerErr != nil {
			l.Error("call the parameter value invalid")
			return nil
		}
		gasPrice = new(big.Int).SetInt64(gasPriceVal)
	}

	funcName, err := getCalledFuncName(c)
	if err != nil {
		l.Error(err.Error())
		return nil
	}
	input := getRpcSpecialParam(c, "input")
	// RLP([funcName][params])
	inputRlp, err := rlp.EncodeToBytes([]interface{}{
		funcName, input,
	})
	if err != nil {
		log.Error("input rlp err")
	}
	var resp interface{}

	l.Info("the From is: ", "From", From.Hex())
	l.Info("the gasLimit is:", "gasLimit", gasLimit)
	l.Info("the gasPrice is:", "gasPrice", gasPrice)
	l.Info("the funcName is:", "funcName", funcName)
	l.Info("the ExtraData is: ", "ExtraData", inputRlp)

	//SendTransactionContract(from, to common.Address,value,gasLimit, gasPrice *big.Int, data []byte, nonce *uint64 )
	if err = client.Call(&resp, getDipperinRpcMethodByName(mName), From, to, nil, gasLimit, gasPrice, inputRlp, nil); err != nil {
		l.Error(fmt.Sprintf("%s Call failed", mName), "err", err)
		return nil
	}
	return resp
}

func getCalledFuncName(c *cli.Context) (funcName string, err error) {
	funcName = c.String("funcName")
	if funcName == "" {
		return "", errors.New("function name is need")
	}
	return funcName, nil
}

func isCreate(c *cli.Context) bool {
	return c.Bool("isCreate")
}

func generateExtraData(c *cli.Context) (ExtraData []byte, err error) {
	abiPath, err := getRpcParamValue(c, "abi")
	if err != nil {
		return nil, errors.New("the abi path value invalid")
	}
	abiBytes, err := ioutil.ReadFile(abiPath)
	if err != nil {
		return nil, errors.New("the abi file read err")
	}
	var wasmAbi utils.WasmAbi
	err = wasmAbi.FromJson(abiBytes)
	//err = json.Unmarshal(abiBytes, &wasmAbi.AbiArr)
	if err != nil {
		return nil, errors.New("abi file is err")
	}

	var args []utils.InputParam
	for _, v := range wasmAbi.AbiArr {
		if strings.EqualFold("init", v.Name) && strings.EqualFold(v.Type, "function") {
			args = v.Inputs
			if len(v.Outputs) != 0 {
				return nil, errors.New("invalid init function outputs length")
			}
			break
		}
	}

	input := getRpcSpecialParam(c, "input")
	params := getRpcParamFromString(input)
	if len(params) != len(args) {
		l.Error("not enough create contract params")
	}

	wasmPath, err := getRpcParamValue(c, "wasm")
	if err != nil {
		l.Error("the wasm path value invalid")
		return
	}
	wasmBytes, err := ioutil.ReadFile(wasmPath)
	if err != nil {
		l.Error("the abi file read err")
		return
	}

	rlpParams := []interface{}{
		//strconv.Itoa(common.AddressTypeContractCreate),wasmBytes, abiBytes,
		wasmBytes, abiBytes,
	}
	for i, v := range args {
		bts := params[i]
		re, err := utils.StringConverter(bts, v.Type)
		if err != nil {
			return re, err
		}
		rlpParams = append(rlpParams, re)
	}
	return rlp.EncodeToBytes(rlpParams)
}

func geneteInputRlpBytes(input string) (result []byte, err error) {
	return rlp.EncodeToBytes([]interface{}{strconv.Itoa(common.AddressTypeContractCreate), "init", input})
}
