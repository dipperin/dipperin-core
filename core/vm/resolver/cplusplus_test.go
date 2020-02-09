package resolver

import (
	"bytes"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/core/vm/base"
	"github.com/dipperin/dipperin-core/tests/factory/vminfo"
	"github.com/dipperin/dipperin-core/third_party/crypto"
	"github.com/dipperin/dipperin-core/third_party/life/exec"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func GetBaseVmRelatedInfo(t *testing.T) (*gomock.Controller, *MockStateDBService, VmContextService,  ContractService, exec.ImportResolver, *exec.VirtualMachine) {
	ctrl := gomock.NewController(t)
	state := NewMockStateDBService(ctrl)
	vmValue := NewMockVmContextService(ctrl)
	contractService := NewMockContractService(ctrl)
	solver := NewResolver(vmValue, contractService, state)

	code, _ := vminfo.GetTestData("event")
	base.NewMockStateDB(ctrl)
	vm, err := exec.NewVirtualMachine(code, base.DEFAULT_VM_CONFIG, solver, nil)
	assert.NoError(t, err)

	return ctrl, state,vmValue, contractService, solver, vm
}

func Test_envMemset(t *testing.T)  {
	ctrl, _, _,_, _, vm := GetBaseVmRelatedInfo(t)
	defer ctrl.Finish()

	dataStart := int64(0)
	value := int64(12)
	num := int64(100)
	t.Log("currentFrame", vm.CurrentFrame)
	vm.CurrentFrame++
	frame := exec.Frame{
		Locals: []int64{dataStart, value, num},
	}
	vm.CallStack[0] = frame
	data := bytes.Repeat([]byte{byte(value)}, int(num))

	result  := envMemset(vm)
	assert.Equal(t, int64(dataStart), result)
	assert.Equal(t, data, vm.Memory.Memory[dataStart:dataStart+num])

	gasCost, _ := envMemsetGasCost(vm)
	assert.Equal(t, uint64(num), gasCost)

}

func Test_envPrintsGasCost(t *testing.T)  {
	ctrl, _, _,_, _, vm := GetBaseVmRelatedInfo(t)
	defer ctrl.Finish()

	dataStart := int64(0)
	data := []byte("hello")
	t.Log("currentFrame", vm.CurrentFrame)
	vm.CurrentFrame++
	frame := exec.Frame{
		Locals: []int64{dataStart},
	}
	vm.CallStack[0] = frame
	vm.Memory.Memory = append(data, vm.Memory.Memory...)

	result , _ := envPrintsGasCost(vm)
	assert.Equal(t, uint64(len(data)), result)
}


func Test_envPrintslGasCost(t *testing.T)  {
	ctrl, _, _,_, _, vm := GetBaseVmRelatedInfo(t)
	defer ctrl.Finish()

	dataStart := int64(0)
	dataLen := int64(100)
	t.Log("currentFrame", vm.CurrentFrame)
	vm.CurrentFrame++
	frame := exec.Frame{
		Locals: []int64{dataStart,dataLen},
	}
	vm.CallStack[0] = frame

	result , _ := envPrintslGasCost(vm)
	assert.Equal(t, uint64(dataLen), result)
}

func Test_envGasPrice(t *testing.T)  {

	ctrl, _, vmValue,_, solver, vm := GetBaseVmRelatedInfo(t)
	defer ctrl.Finish()

	vm.CurrentFrame++
	gasPrice := int64(1000)

	vmValue.(*MockVmContextService).EXPECT().GetGasPrice().Return(big.NewInt(gasPrice))

	result := solver.(*Resolver).envGasPrice(vm)
	assert.Equal(t, gasPrice, result)
}


func Test_envBlockHash(t *testing.T)  {
	ctrl, _, vmValue,_, solver, vm := GetBaseVmRelatedInfo(t)
	defer ctrl.Finish()

	blockNum := int64(100)
	offsize := int64(20)
	t.Log("currentFrame", vm.CurrentFrame)
	vm.CurrentFrame++
	frame := exec.Frame{
		Locals: []int64{blockNum, offsize},
	}
	vm.CallStack[0] = frame
    blockHash := common.HexToHash("0x12345")

	vmValue.(*MockVmContextService).EXPECT().GetBlockHash(uint64(blockNum)).Return(blockHash)

	result := solver.(*Resolver).envBlockHash(vm)
	assert.Equal(t, int64(0), result)
	assert.Equal(t, blockHash.Bytes(), vm.Memory.Memory[offsize:offsize+int64(len(blockHash.Bytes()))])
}

func Test_envNumber(t *testing.T)  {
	ctrl, _, vmValue,_, solver, vm := GetBaseVmRelatedInfo(t)
	defer ctrl.Finish()

	vm.CurrentFrame++

	vmValue.(*MockVmContextService).EXPECT().GetBlockNumber().Return(big.NewInt(int64(1000)))

	result := solver.(*Resolver).envNumber(vm)
	assert.Equal(t, int64(1000), result)
}


func Test_envGasLimit(t *testing.T)  {

	ctrl, _, vmValue,_, solver, vm := GetBaseVmRelatedInfo(t)
	defer ctrl.Finish()

	vm.CurrentFrame++

	vmValue.(*MockVmContextService).EXPECT().GetGasLimit().Return(uint64(1000))

	result := solver.(*Resolver).envGasLimit(vm)
	assert.Equal(t, int64(1000), result)
}

func Test_envTimestamp(t *testing.T)  {

	ctrl, _, vmValue,_, solver, vm := GetBaseVmRelatedInfo(t)
	defer ctrl.Finish()

	vm.CurrentFrame++

	vmValue.(*MockVmContextService).EXPECT().GetTime().Return(big.NewInt(int64(1000)))

	result := solver.(*Resolver).envTimestamp(vm)
	assert.Equal(t, int64(1000), result)
}

func Test_envCoinbase(t *testing.T)  {
	ctrl, _, vmValue,_, solver, vm := GetBaseVmRelatedInfo(t)
	defer ctrl.Finish()

	dataStart := int64(0)
	t.Log("currentFrame", vm.CurrentFrame)
	vm.CurrentFrame++
	frame := exec.Frame{
		Locals: []int64{dataStart},
	}
	vm.CallStack[0] = frame

	vmValue.(*MockVmContextService).EXPECT().GetCoinBase().Return(model.AliceAddr)

	result := solver.(*Resolver).envCoinbase(vm)
	assert.Equal(t, int64(0), result)
	assert.Equal(t, model.AliceAddr.Bytes(), vm.Memory.Memory[dataStart:len(model.AliceAddr.Bytes())])
}

func Test_envBalance(t *testing.T)  {

}

func Test_envOrigin(t *testing.T)  {
	ctrl, _, vmValue,_, solver, vm := GetBaseVmRelatedInfo(t)
	defer ctrl.Finish()

	dataStart := int64(0)
	t.Log("currentFrame", vm.CurrentFrame)
	vm.CurrentFrame++
	frame := exec.Frame{
		Locals: []int64{dataStart},
	}
	vm.CallStack[0] = frame

	vmValue.(*MockVmContextService).EXPECT().GetOrigin().Return(model.AliceAddr)

	result := solver.(*Resolver).envOrigin(vm)
	assert.Equal(t, int64(0), result)
	assert.Equal(t, model.AliceAddr.Bytes(), vm.Memory.Memory[dataStart:len(model.AliceAddr.Bytes())])
}

func Test_envCaller(t *testing.T)  {
	ctrl, _, _,contractService, solver, vm := GetBaseVmRelatedInfo(t)
	defer ctrl.Finish()

	dataStart := int64(0)
	t.Log("currentFrame", vm.CurrentFrame)
	vm.CurrentFrame++
	frame := exec.Frame{
		Locals: []int64{dataStart},
	}
	vm.CallStack[0] = frame

	contractService.(*MockContractService).EXPECT().Caller().Return(base.AccountRef(model.AliceAddr))

	result := solver.(*Resolver).envCaller(vm)
	assert.Equal(t, int64(0), result)
	assert.Equal(t, model.AliceAddr.Bytes(), vm.Memory.Memory[dataStart:len(model.AliceAddr.Bytes())])
}

func Test_envCallValue(t *testing.T)  {

}

func Test_envCallValueUDIP(t *testing.T)  {
	ctrl, _, _,contractService, solver, vm := GetBaseVmRelatedInfo(t)
	defer ctrl.Finish()

	vm.CurrentFrame++
	contractService.(*MockContractService).EXPECT().CallValue().Return(new(big.Int).Exp(big.NewInt(10),big.NewInt(19),big.NewInt(0)))

	result := solver.(*Resolver).envCallValueUDIP(vm)
	assert.Equal(t, int64(10000), result)
}

func Test_envAddress(t *testing.T)  {
	ctrl, _, _,contractService, solver, vm := GetBaseVmRelatedInfo(t)
	defer ctrl.Finish()

	dataStart := int64(0)
	t.Log("currentFrame", vm.CurrentFrame)
	vm.CurrentFrame++
	frame := exec.Frame{
		Locals: []int64{dataStart},
	}
	vm.CallStack[0] = frame
	contractService.(*MockContractService).EXPECT().Address().Return(model.AliceAddr)

	pos := solver.(*Resolver).envAddress(vm)
	assert.Equal(t, int64(0), pos)
	assert.Equal(t, model.AliceAddr.Bytes(), vm.Memory.Memory[dataStart:dataStart+int64(len(common.Address{}))])
}

func Test_envSha3(t *testing.T)  {
	ctrl, _, _,_, solver, vm := GetBaseVmRelatedInfo(t)
	defer ctrl.Finish()

	data := []byte("hello")
	dataHash := crypto.Keccak256(data)

	dataStart := int64(0)
	dataLen := int64(len(data))
	destStart := dataLen
	destLen := int64(len(dataHash))


	t.Log("currentFrame", vm.CurrentFrame)
	vm.CurrentFrame++
	frame := exec.Frame{
		Locals: []int64{dataStart, dataLen, destStart, destLen},
	}
	vm.CallStack[0] = frame

	//address := vm.Memory.Memory[addr:addr+common.AddressLength]
	backend := vm.Memory.Memory[dataStart+dataLen:]
	vm.Memory.Memory = append(append(vm.Memory.Memory[:dataStart],  data...), backend...)

	pos := solver.(*Resolver).envSha3(vm)
	assert.Equal(t, int64(0), pos)
	assert.Equal(t, dataHash, vm.Memory.Memory[destStart:destStart+destLen])
}

func Test_envHexStringSameWithVM(t *testing.T)  {
	ctrl, _, _,_, solver, vm := GetBaseVmRelatedInfo(t)
	defer ctrl.Finish()

	data := []byte("hEllo")
	//dataHash := crypto.Keccak256(data)

	hexString := common.HexStringSameWithVM(string(data))

	dataStart := int64(0)
	dataLen := int64(len(data))
	destStart := dataLen


	t.Log("currentFrame", vm.CurrentFrame)
	vm.CurrentFrame++
	frame := exec.Frame{
		Locals: []int64{dataStart, dataLen, destStart},
	}
	vm.CallStack[0] = frame

	//address := vm.Memory.Memory[addr:addr+common.AddressLength]
	backend := vm.Memory.Memory[dataStart+dataLen:]
	vm.Memory.Memory = append(append(vm.Memory.Memory[:dataStart],  data...), backend...)

	pos := solver.(*Resolver).envHexStringSameWithVM(vm)
	assert.Equal(t, int64(0), pos)
	assert.Equal(t, []byte(hexString), vm.Memory.Memory[destStart:destStart+int64(len(hexString))])
}

func Test_envGetCallerNonce(t *testing.T)  {
	ctrl, state, _,contractService, solver, vm := GetBaseVmRelatedInfo(t)
	defer ctrl.Finish()
	expectNonce := 10
	contractService.(*MockContractService).EXPECT().Caller().Return(base.AccountRef(common.Address{}))
	state.EXPECT().GetNonce(common.Address{}).Return(uint64(expectNonce),nil)


	result := solver.(*Resolver).envGetCallerNonce(vm)
	assert.Equal(t, int64(expectNonce), result)
}

func Test_envGetSignerAddress(t *testing.T)  {
	ctrl, _, _,contractService, solver, vm := GetBaseVmRelatedInfo(t)
	defer ctrl.Finish()

	//params, err := rlp.EncodeToBytes([]interface{}{"returnString"})
	sha3Data := "f2abd116754ce66ff83fb14f3d7f8b104ab0201a74398cafc3f4719e90ca114e"
	signatureData := "8759548f229a2cebdc608a1496105906993850891b9e6e406e391d7442da2aef3330d5cb379188ad56e656c6aa3ea7c42968ca63a07527d3dbaae9932c218b8f00"

	sha3DataBytes := common.Hex2Bytes(sha3Data)

	sha3DataStart := int64(0)
	sha3DataLen := int64(len(sha3DataBytes))
	signatureStart := sha3DataLen
	signatureLen := int64(len(signatureData))
	returnStart := int64(len(sha3Data)+len(signatureData))

	//paramLen := int64(len(params))

	t.Log("currentFrame", vm.CurrentFrame)
	vm.CurrentFrame++
	frame := exec.Frame{
		Locals: []int64{sha3DataStart, sha3DataLen, signatureStart, signatureLen,  returnStart},
	}
	vm.CallStack[0] = frame

	//address := vm.Memory.Memory[addr:addr+common.AddressLength]
	backend := vm.Memory.Memory[sha3DataStart+sha3DataLen:]
	vm.Memory.Memory = append(append(vm.Memory.Memory[:sha3DataStart],  sha3DataBytes...), backend...)
	backendSign := vm.Memory.Memory[signatureStart+signatureLen:]
	vm.Memory.Memory = append(append(vm.Memory.Memory[:signatureStart],  signatureData...), backendSign...)


	//contractService.EXPECT().Self().Return(base.AccountRef(common.BytesToAddress(address))).AnyTimes()
	contractService.(*MockContractService).EXPECT().GetGas().Return(uint64(0)).AnyTimes()
	contractService.(*MockContractService).EXPECT().CallValue().Return(big.NewInt(100)).AnyTimes()
	contractService.(*MockContractService).EXPECT().Self().Return(base.AccountRef(common.HexToAddress("0x0000430d04D9db5aac60a88848B80e26991bBFd13C74")) )

	pos := solver.(*Resolver).envGetSignerAddress(vm)
	assert.Equal(t, int64(0), pos)
}

func Test_envDipperCallAndDelegateCall(t *testing.T) {
	ctrl, _, vmValue,contractService, solver, vm := GetBaseVmRelatedInfo(t)
	defer ctrl.Finish()

	params, err := rlp.EncodeToBytes([]interface{}{"returnInt"})
	assert.NoError(t, err)
	addr := int64(0)
	param := int64(50)

	paramLen := int64(len(params))

	t.Log("currentFrame", vm.CurrentFrame)
	vm.CurrentFrame++
	frame := exec.Frame{
		Locals: []int64{addr, param, paramLen},
	}
	vm.CallStack[0] = frame

	address := vm.Memory.Memory[addr:addr+common.AddressLength]
	backend := vm.Memory.Memory[param+paramLen:]
	vm.Memory.Memory = append(append(vm.Memory.Memory[:param],  params...), backend...)
	transferValue := big.NewInt(100)

	contractService.(*MockContractService).EXPECT().Self().Return(base.AccountRef(common.BytesToAddress(address))).AnyTimes()
	contractService.(*MockContractService).EXPECT().GetGas().Return(uint64(0)).AnyTimes()
	contractService.(*MockContractService).EXPECT().CallValue().Return(transferValue).AnyTimes()

	testCases := []struct{
		name string
		given func() int64
		expect int64
	}{
		{
			name:"Test_envDipperDelegateCallString",
			given: func() int64 {
				vmValue.(*MockVmContextService).EXPECT().DelegateCall(contractService, common.BytesToAddress(address),params,uint64(0)).Return([]byte(nil),uint64(0),nil)

				pos := solver.(*Resolver).envDipperDelegateCallString(vm)
				return pos
			},
			expect:int64(131072),
		},
		{
			name:"Test_envDipperDelegateCallInt64",
			given: func() int64 {
				vmValue.(*MockVmContextService).EXPECT().DelegateCall(contractService, common.BytesToAddress(address),params,uint64(0)).Return([]byte(nil),uint64(0),nil)

				result := solver.(*Resolver).envDipperDelegateCallInt64(vm)
				return result
			},
			expect:int64(0),
		},
		{
			name:"test_envDipperCallInt64",
			given: func() int64 {
				vmValue.(*MockVmContextService).EXPECT().Call(contractService, common.BytesToAddress(address),params,uint64(0), transferValue).Return([]byte(nil),uint64(0),nil).AnyTimes()

				result := solver.(*Resolver).envDipperCallInt64(vm)
				return result
			},
			expect:int64(0),
		},
		{
			name:"test_envDipperDelegateCallInt64",
			given: func() int64 {
				vmValue.(*MockVmContextService).EXPECT().DelegateCall(contractService, common.BytesToAddress(address),params,uint64(0)).Return([]byte(nil),uint64(0),nil).AnyTimes()

				result := solver.(*Resolver).envDipperDelegateCallInt64(vm)
				return result
			},
			expect:int64(0),
		},
		{
			name:"test_envDipperCallInt64",
			given: func() int64 {
				vmValue.(*MockVmContextService).EXPECT().Call(contractService, common.BytesToAddress(address),params,uint64(0), transferValue).Return([]byte(nil),uint64(0),nil).AnyTimes()

				result := solver.(*Resolver).envDipperCallInt64(vm)
				return result
			},
			expect:int64(0),
		},
	}

	for _,tc := range testCases{
		t.Log("Test_envDipperCallAndDelegateCall", tc.name)
		result := tc.given()
		assert.Equal(t, tc.expect, result)
	}


}

func Test_envDipperCallString(t *testing.T) {
	ctrl, _, vmValue,contractService, solver, vm := GetBaseVmRelatedInfo(t)
	defer ctrl.Finish()

	params, err := rlp.EncodeToBytes([]interface{}{"returnString"})
	assert.NoError(t, err)
	addr := int64(0)
	param := int64(50)

	paramLen := int64(len(params))

	t.Log("currentFrame", vm.CurrentFrame)
	vm.CurrentFrame++
	frame := exec.Frame{
		Locals: []int64{addr, param, paramLen},
	}
	vm.CallStack[0] = frame

	address := vm.Memory.Memory[addr:addr+common.AddressLength]
	backend := vm.Memory.Memory[param+paramLen:]
	vm.Memory.Memory = append(append(vm.Memory.Memory[:param],  params...), backend...)

	contractService.(*MockContractService).EXPECT().Self().Return(base.AccountRef(common.BytesToAddress(address))).AnyTimes()
	contractService.(*MockContractService).EXPECT().GetGas().Return(uint64(0)).AnyTimes()
	contractService.(*MockContractService).EXPECT().CallValue().Return(big.NewInt(100)).AnyTimes()
	vmValue.(*MockVmContextService).EXPECT().Call(contractService, common.BytesToAddress(address),params,uint64(0),big.NewInt(100)).Return([]byte(nil),uint64(0),nil)

	pos := solver.(*Resolver).envDipperCallString(vm)
	assert.Equal(t, int64(131072), pos)
}


//func Test_envDipperCallInt64(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//	state := NewMockStateDBService(ctrl)
//	vmValue := NewMockVmContextService(ctrl)
//	contractService := NewMockContractService(ctrl)
//	solver := NewResolver(vmValue, contractService, state)
//
//	code, _ := vminfo.GetTestData("event")
//	vm, err := exec.NewVirtualMachine(code, base.DEFAULT_VM_CONFIG, solver, nil)
//	assert.NoError(t, err)
//
//	params, err := rlp.EncodeToBytes([]interface{}{"returnInt"})
//	assert.NoError(t, err)
//	addr := int64(0)
//	param := int64(50)
//
//	paramLen := int64(len(params))
//
//	t.Log("currentFrame", vm.CurrentFrame)
//	vm.CurrentFrame++
//	frame := exec.Frame{
//		Locals: []int64{addr, param, paramLen},
//	}
//	vm.CallStack[0] = frame
//
//	address := vm.Memory.Memory[addr:addr+common.AddressLength]
//	backend := vm.Memory.Memory[param+paramLen:]
//	vm.Memory.Memory = append(append(vm.Memory.Memory[:param],  params...), backend...)
//	transferValue := big.NewInt(100)
//
//	contractService.EXPECT().Self().Return(base.AccountRef(common.BytesToAddress(address))).AnyTimes()
//	contractService.EXPECT().GetGas().Return(uint64(0)).AnyTimes()
//	contractService.EXPECT().CallValue().Return(transferValue).AnyTimes()
//
//	//testCases := []struct{
//	//	name string
//	//	given func() int64
//	//	expect int64
//	//}{
//	//	{
//	//		name:"Test_envDipperDelegateCallString",
//	//		given: func() int64 {
//	//			vmValue.EXPECT().DelegateCall(contractService, common.BytesToAddress(address),params,uint64(0)).Return([]byte(nil),uint64(0),nil)
//	//
//	//			pos := solver.(*Resolver).envDipperDelegateCallString(vm)
//	//			return pos
//	//		},
//	//		expect:int64(131072),
//	//	},
//	//	{
//	//		name:"Test_envDipperCallString",
//	//		given: func() int64 {
//	//			vmValue.EXPECT().Call(contractService, common.BytesToAddress(address),params,uint64(0),big.NewInt(100)).Return([]byte(nil),uint64(0),nil)
//	//
//	//			pos := solver.(*Resolver).envDipperCallString(vm)
//	//			return pos
//	//		},
//	//		expect:int64(131072),
//	//	},
//	//	{
//	//		name:"Test_envDipperDelegateCallInt64",
//	//		given: func() int64 {
//	//			vmValue.EXPECT().DelegateCall(contractService, common.BytesToAddress(address),params,uint64(0)).Return([]byte(nil),uint64(0),nil)
//	//
//	//			result := solver.(*Resolver).envDipperDelegateCallInt64(vm)
//	//			return result
//	//		},
//	//		expect:int64(0),
//	//	},
//	//	{
//	//		name:"test_envDipperCallInt64",
//	//		given: func() int64 {
//	//			vmValue.EXPECT().Call(contractService, common.BytesToAddress(address),params,uint64(0), transferValue).Return([]byte(nil),uint64(0),nil)
//	//
//	//			result := solver.(*Resolver).envDipperCallInt64(vm)
//	//			return result
//	//		},
//	//		expect:int64(0),
//	//	},
//	//}
//	//
//	//for _,tc := range testCases{
//	//	t.Log("Test_envDipperCallAndDelegateCall", tc.name)
//	//	result := tc.given()
//	//	assert.Equal(t, tc.expect, result)
//	//}
//
//		vmValue.EXPECT().Call(contractService, common.BytesToAddress(address),params,uint64(0), transferValue).Return([]byte(nil),uint64(0),nil)
//
//		result := solver.(*Resolver).envDipperCallInt64(vm)
//		assert.Equal(t, int64(0),result)
//
//}
//
//func Test_envDipperDelegateCallInt64(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//	state := NewMockStateDBService(ctrl)
//	vmValue := NewMockVmContextService(ctrl)
//	contractService := NewMockContractService(ctrl)
//	solver := NewResolver(vmValue, contractService, state)
//
//	code, _ := vminfo.GetTestData("event")
//	vm, err := exec.NewVirtualMachine(code, base.DEFAULT_VM_CONFIG, solver, nil)
//	assert.NoError(t, err)
//
//	params, err := rlp.EncodeToBytes([]interface{}{"returnInt"})
//	assert.NoError(t, err)
//	addr := int64(0)
//	param := int64(50)
//
//	paramLen := int64(len(params))
//
//	t.Log("currentFrame", vm.CurrentFrame)
//	vm.CurrentFrame++
//	frame := exec.Frame{
//		Locals: []int64{addr, param, paramLen},
//	}
//	vm.CallStack[0] = frame
//
//	address := vm.Memory.Memory[addr:addr+common.AddressLength]
//	backend := vm.Memory.Memory[param+paramLen:]
//	vm.Memory.Memory = append(append(vm.Memory.Memory[:param],  params...), backend...)
//
//	contractService.EXPECT().Self().Return(base.AccountRef(common.BytesToAddress(address))).AnyTimes()
//	contractService.EXPECT().GetGas().Return(uint64(0)).AnyTimes()
//	contractService.EXPECT().CallValue().Return(big.NewInt(100)).AnyTimes()
//	vmValue.EXPECT().DelegateCall(contractService, common.BytesToAddress(address),params,uint64(0)).Return([]byte(nil),uint64(0),nil)
//
//	result := solver.(*Resolver).envDipperDelegateCallInt64(vm)
//	assert.Equal(t, int64(0), result)
//}
//
//
//func Test_envDipperDelegateCallString(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//	state := NewMockStateDBService(ctrl)
//	vmValue := NewMockVmContextService(ctrl)
//	contractService := NewMockContractService(ctrl)
//	solver := NewResolver(vmValue, contractService, state)
//
//	code, _ := vminfo.GetTestData("event")
//	vm, err := exec.NewVirtualMachine(code, base.DEFAULT_VM_CONFIG, solver, nil)
//	assert.NoError(t, err)
//
//	params, err := rlp.EncodeToBytes([]interface{}{"returnString"})
//	assert.NoError(t, err)
//	addr := int64(0)
//	param := int64(50)
//
//	paramLen := int64(len(params))
//
//	t.Log("currentFrame", vm.CurrentFrame)
//	vm.CurrentFrame++
//	frame := exec.Frame{
//		Locals: []int64{addr, param, paramLen},
//	}
//	vm.CallStack[0] = frame
//
//	address := vm.Memory.Memory[addr:addr+common.AddressLength]
//	backend := vm.Memory.Memory[param+paramLen:]
//	vm.Memory.Memory = append(append(vm.Memory.Memory[:param],  params...), backend...)
//
//	contractService.EXPECT().Self().Return(base.AccountRef(common.BytesToAddress(address))).AnyTimes()
//	contractService.EXPECT().GetGas().Return(uint64(0)).AnyTimes()
//	contractService.EXPECT().CallValue().Return(big.NewInt(100)).AnyTimes()
//	vmValue.EXPECT().DelegateCall(contractService, common.BytesToAddress(address),params,uint64(0)).Return([]byte(nil),uint64(0),nil)
//
//	pos := solver.(*Resolver).envDipperDelegateCallString(vm)
//	assert.Equal(t, int64(131072), pos)
//}
