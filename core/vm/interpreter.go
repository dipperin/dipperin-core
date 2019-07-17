package vm

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/g-error"
	"github.com/dipperin/dipperin-core/common/math"
	"github.com/dipperin/dipperin-core/core/vm/common/utils"
	"github.com/dipperin/dipperin-core/core/vm/resolver"
	"github.com/dipperin/dipperin-core/third-party/life/exec"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
	"reflect"
	"runtime"
	"strings"
)

var (
	errReturnInvalidRlpFormat   = errors.New("interpreter_life: invalid rlp format")
	errReturnInsufficientParams = errors.New("interpreter_life: invalid input. ele must greater than 1")
	errReturnInvalidAbi         = errors.New("interpreter_life: invalid abi, encoded fail")
	errReturnInputAbiNotMatch   = errors.New("interpreter_life: length of input and abi not match")
)

const (
	CALL_CONTRACT_FLAG = 9
)

type Interpreter interface {
	// Run loops and evaluates the contract's code with the given input data and returns
	// the return byte-slice and an error if one occurred.
	Run(vm *VM, contract *Contract, create bool) ([]byte, error)
	// CanRun tells if the contract, passed as an argument, can be
	CanRun([]byte) bool
}

type WASMInterpreter struct {
	state    StateDB
	context  *Context
	config   exec.VMConfig
	resolver exec.ImportResolver
}

// NewWASMInterpreter returns a new instance of the Interpreter
func NewWASMInterpreter(state StateDB, context Context, vmConfig exec.VMConfig) *WASMInterpreter {
	return &WASMInterpreter{
		state,
		&context,
		vmConfig,
		&resolver.Resolver{},
	}
}

func (in *WASMInterpreter) Run(vm *VM, contract *Contract, create bool) (ret []byte, err error) {
	defer func() {
		if er := recover(); er != nil {
			fmt.Println(stack())
			ret, err = nil, fmt.Errorf("VM execute fail: %v", er)
		}
	}()
	vm.depth++
	defer func() {
		vm.depth--
		if vm.depth == 0 {
			log.Info("VM depth = 0")
		}
	}()

	if len(contract.Code) == 0 || len(contract.ABI) == 0 {
		log.Debug("Code or ABI Length is 0")
		return nil, nil
	}

	//　life方法注入新建虚拟机
	solver := resolver.NewResolver(vm, contract, in.state)
	lifeVm, err := exec.NewVirtualMachine(contract.Code, in.config, solver, nil)
	if err != nil {
		log.Info("NewVirtualMachine failed", "err", err)
		return nil, err
	}
	lifeVm.GasLimit = contract.Gas
	defer func() {
		lifeVm.Stop()
	}()

	var (
		funcName   string
		params     []int64
		returnType string
	)

	if create {
		// init function.
		funcName, params, returnType, err = parseInitFunctionByABI(lifeVm, contract.Input, contract.ABI)
		if err != nil {
			log.Error("parseInitFunctionByABI failed", "err", err)
			return nil, err
		}
	} else {
		// parse input
		funcName, params, returnType, err = parseCallExtraDataByABI(lifeVm, contract.Input, contract.ABI)
		if err != nil {
			if err == errReturnInsufficientParams { // transfer to contract address.
				return nil, nil
			}
			log.Error("parseCallExtraDataByABI failed", "err", err)
			return nil, err
		}
	}
	log.Info("WASMInterpreter Run", "funcName", funcName, "params", params, "return", returnType, "err", err)

	//　获取entryID
	entryID, ok := lifeVm.GetFunctionExport(funcName)
	if !ok {
		return nil, fmt.Errorf("entryId not found")
	}

	res, err := lifeVm.Run(entryID, params...)
	log.Info("Run lifeVm", "gasUsed", lifeVm.GasUsed, "gasLimit", lifeVm.GasLimit)
	if err != nil {
		log.Error("throw exception:", "err", err)
		return nil, err
	}

	if contract.Gas > lifeVm.GasUsed {
		contract.Gas = contract.Gas - lifeVm.GasUsed
	} else {
		return nil, g_error.ErrOutOfGas
	}

	if create {
		return contract.Code, nil
	}

	switch returnType {
	case "void", "int8", "int16", "int32", "int64":
		bigRes := new(big.Int)
		bigRes.SetInt64(res)
		finalRes := utils.Align32Bytes(math.U256(bigRes).Bytes())
		return finalRes, nil
	case "uint8", "uint16", "uint32", "uint64":
		finalRes := utils.Align32Bytes(utils.Uint64ToBytes(uint64(res)))
		return finalRes, nil
	case "string":
		returnBytes := make([]byte, 0)
		copyData := lifeVm.Memory.Memory[res:]
		for _, v := range copyData {
			if v == 0 {
				break
			}
			returnBytes = append(returnBytes, v)
		}
		strHash := common.BytesToHash(utils.Int32ToBytes(32))
		sizeHash := common.BytesToHash(utils.Int64ToBytes(int64(len(returnBytes))))
		var dataRealSize = len(returnBytes)
		if (dataRealSize % 32) != 0 {
			dataRealSize = dataRealSize + (32 - (dataRealSize % 32))
		}
		dataByt := make([]byte, dataRealSize)
		copy(dataByt[0:], returnBytes)

		finalData := make([]byte, 0)
		finalData = append(finalData, strHash.Bytes()...)
		finalData = append(finalData, sizeHash.Bytes()...)
		finalData = append(finalData, dataByt...)
		return finalData, nil
	}
	return nil, nil
}

func (in *WASMInterpreter) CanRun([]byte) bool {
	return true
}

// input = RLP([funcName][params])
func ParseInputForFuncName(rlpData []byte) (funcName string, err error) {
	ptr := new(interface{})
	err = rlp.Decode(bytes.NewReader(rlpData), &ptr)
	if err != nil {
		return "", err
	}

	rlpList := reflect.ValueOf(ptr).Elem().Interface()
	if _, ok := rlpList.([]interface{}); !ok {
		return "", errReturnInvalidRlpFormat
	}

	iRlpList := rlpList.([]interface{})
	if len(iRlpList) == 0 {
		return "", errReturnInsufficientParams
	}
	fmt.Println("rlpList", rlpList, iRlpList[0].([]byte) )

	if v, ok := iRlpList[0].([]byte); ok {
		funcName = string(v)
	}
	return
}

// input = RLP([params])
// returnType must void
func parseInitFunctionByABI(vm *exec.VirtualMachine, input []byte, abi []byte) (funcName string, params []int64, returnType string, err error) {
	funcName = "init"
	if input == nil {
		log.Info("InitFunc has no input")
		return
	}

	// rlp decode
	ptr := new(interface{})
	err = rlp.Decode(bytes.NewReader(input), &ptr)
	if err != nil {
		return funcName, nil, "", err
	}

	rlpList := reflect.ValueOf(ptr).Elem().Interface()
	if _, ok := rlpList.([]interface{}); !ok {
		return funcName, nil, "", errReturnInvalidRlpFormat
	}

	params, returnType, err = findParams(vm, abi, funcName, rlpList.([]interface{}))
	if returnType != "void" {
		return funcName, nil, returnType, errors.New("InitFunc returnType must be void")
	}
	return
}

// input = RLP([funcName][params])
// get returnType[0] if more than 1 return
func parseCallExtraDataByABI(vm *exec.VirtualMachine, input []byte, abi []byte) (funcName string, params []int64, returnType string, err error) {
	if input == nil || len(input) <= 1 {
		return "", nil, "", fmt.Errorf("invalid input")
	}

	// rlp decode
	ptr := new(interface{})
	err = rlp.Decode(bytes.NewReader(input), &ptr)
	if err != nil {
		return "", nil, "", err
	}

	rlpList := reflect.ValueOf(ptr).Elem().Interface()
	if _, ok := rlpList.([]interface{}); !ok {
		return "", nil, "", errReturnInvalidRlpFormat
	}

	iRlpList := rlpList.([]interface{})
	if len(iRlpList) == 0 {
		return "", nil, "", errReturnInsufficientParams
	}

	if v, ok := iRlpList[0].([]byte); ok {
		funcName = string(v)
	}

	if len(iRlpList) > 1 {
		var inputList []interface{}
		for _, value := range iRlpList[1:] {
			inputList = append(inputList, value)
		}
		params, returnType, err = findParams(vm, abi, funcName, inputList)
	}
	return
}

func findParams(vm *exec.VirtualMachine, abi []byte, funcName string, inputList []interface{}) (params []int64, returnType string, err error) {
	wasmAbi := new(utils.WasmAbi)
	// TODO
	//  err = json.Unmarshal(abi, wasmAbi)
	err = wasmAbi.FromJson(abi)
	if err != nil {
		return nil, "", errReturnInvalidAbi
	}

	var abiParam []utils.InputParam
	for _, v := range wasmAbi.AbiArr {
		if strings.EqualFold(v.Name, funcName) && strings.EqualFold(v.Type, "function") {
			abiParam = v.Inputs
			fmt.Println("len outputs", len(v.Outputs), "abiParam", len(abiParam), "inputlist", len(inputList) )
			log.Info("findParams", "len outputs", len(v.Outputs), "abiParam", len(abiParam), "inputlist", len(inputList))
			if len(v.Outputs) != 0 {
				returnType = v.Outputs[0].Type
			} else {
				returnType = "void"
			}
			break
		}
	}

	if len(abiParam) != len(inputList) {
		log.Error("findParams failed", "err", errReturnInputAbiNotMatch, "abiLen", len(abiParam), "inputLen", len(inputList))
		return nil, "", errReturnInputAbiNotMatch
	}

	// uint64 uint32  uint16 uint8 int64 int32  int16 int8 float32 float64 string void
	for i, v := range abiParam {
		input := inputList[i].([]byte)
		switch v.Type {
		case "string":
			pos := resolver.MallocString(vm, string(input))
			params = append(params, pos)
		case "int8":
			params = append(params, int64(input[0]))
		case "int16":
			params = append(params, int64(binary.BigEndian.Uint16(input)))
		case "int32", "int":
			params = append(params, int64(binary.BigEndian.Uint32(input)))
		case "int64":
			params = append(params, int64(binary.BigEndian.Uint64(input)))
		case "uint8":
			params = append(params, int64(input[0]))
		case "uint32", "uint":
			params = append(params, int64(binary.BigEndian.Uint32(input)))
		case "uint64":
			params = append(params, int64(binary.BigEndian.Uint64(input)))
		case "bool":
			params = append(params, int64(input[0]))
		}
	}
	return
}

// rlpData=RLP([code][abi][init params])
func parseCreateExtraData(rlpData []byte) (code, abi, rlpInit []byte, err error) {
	ptr := new(interface{})
	err = rlp.Decode(bytes.NewReader(rlpData), &ptr)
	if err != nil {
		return
	}

	rlpList := reflect.ValueOf(ptr).Elem().Interface()
	if _, ok := rlpList.([]interface{}); !ok {
		return nil, nil, nil, errReturnInvalidRlpFormat
	}

	iRlpList := rlpList.([]interface{})
	if len(iRlpList) < 2 {
		return nil, nil, nil, errReturnInsufficientParams
	}

	// parse code and abi
	if v, ok := iRlpList[0].([]byte); ok {
		code = v
	}

	if v, ok := iRlpList[1].([]byte); ok {
		abi = v
	}

	if len(iRlpList) == 2 {
		return
	}

	// encode init
	var init []interface{}
	for _, value := range iRlpList[2:] {
		init = append(init, value)
	}

	rlpInit, err = rlp.EncodeToBytes(init)
	if err != nil {
		return
	}
	return
}

func stack() string {
	var buf [2 << 10]byte
	return string(buf[:runtime.Stack(buf[:], true)])
}
