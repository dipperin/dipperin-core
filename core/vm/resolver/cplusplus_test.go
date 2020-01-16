package resolver

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/vm/base"
	"github.com/dipperin/dipperin-core/tests/factory/vminfo"
	"github.com/dipperin/dipperin-core/third_party/life/exec"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func Test_envCallTransferUDIP(t *testing.T)  {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	state := NewMockStateDBService(ctrl)
	vmValue := NewMockVmContextService(ctrl)
	contractService := NewMockContractService(ctrl)
	solver := NewResolver(vmValue, contractService, state)

	code, _ := vminfo.GetTestData("event")
	vm, err := exec.NewVirtualMachine(code, base.DEFAULT_VM_CONFIG, solver, nil)
	assert.NoError(t, err)

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

	contractService.EXPECT().Self().Return(base.AccountRef(common.BytesToAddress(address))).AnyTimes()
	contractService.EXPECT().GetGas().Return(uint64(0)).AnyTimes()
	contractService.EXPECT().CallValue().Return(big.NewInt(100)).AnyTimes()
	vmValue.EXPECT().Call(contractService, common.BytesToAddress(address),params,uint64(0),big.NewInt(100)).Return([]byte(nil),uint64(0),nil)

	pos := solver.(*Resolver).envDipperCallString(vm)
	assert.Equal(t, int64(131072), pos)
}

func Test_envDipperCallAndDelegateCall(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	state := NewMockStateDBService(ctrl)
	vmValue := NewMockVmContextService(ctrl)
	contractService := NewMockContractService(ctrl)
	solver := NewResolver(vmValue, contractService, state)

	code, _ := vminfo.GetTestData("event")
	vm, err := exec.NewVirtualMachine(code, base.DEFAULT_VM_CONFIG, solver, nil)
	assert.NoError(t, err)

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

	contractService.EXPECT().Self().Return(base.AccountRef(common.BytesToAddress(address))).AnyTimes()
	contractService.EXPECT().GetGas().Return(uint64(0)).AnyTimes()
	contractService.EXPECT().CallValue().Return(transferValue).AnyTimes()

	testCases := []struct{
		name string
		given func() int64
		expect int64
	}{
		{
			name:"Test_envDipperDelegateCallString",
			given: func() int64 {
				vmValue.EXPECT().DelegateCall(contractService, common.BytesToAddress(address),params,uint64(0)).Return([]byte(nil),uint64(0),nil)

				pos := solver.(*Resolver).envDipperDelegateCallString(vm)
				return pos
			},
			expect:int64(131072),
		},
		{
			name:"Test_envDipperDelegateCallInt64",
			given: func() int64 {
				vmValue.EXPECT().DelegateCall(contractService, common.BytesToAddress(address),params,uint64(0)).Return([]byte(nil),uint64(0),nil)

				result := solver.(*Resolver).envDipperDelegateCallInt64(vm)
				return result
			},
			expect:int64(0),
		},
		{
			name:"test_envDipperCallInt64",
			given: func() int64 {
				vmValue.EXPECT().Call(contractService, common.BytesToAddress(address),params,uint64(0), transferValue).Return([]byte(nil),uint64(0),nil).AnyTimes()

				result := solver.(*Resolver).envDipperCallInt64(vm)
				return result
			},
			expect:int64(0),
		},
		{
			name:"test_envDipperDelegateCall",
			given: func() int64 {
				vmValue.EXPECT().DelegateCall(contractService, common.BytesToAddress(address),params,uint64(0)).Return([]byte(nil),uint64(0),nil).AnyTimes()

				result := solver.(*Resolver).envDipperCallInt64(vm)
				return result
			},
			expect:int64(0),
		},
		{
			name:"test_envDipperCall",
			given: func() int64 {
				vmValue.EXPECT().Call(contractService, common.BytesToAddress(address),params,uint64(0), transferValue).Return([]byte(nil),uint64(0),nil).AnyTimes()

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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	state := NewMockStateDBService(ctrl)
	vmValue := NewMockVmContextService(ctrl)
	contractService := NewMockContractService(ctrl)
	solver := NewResolver(vmValue, contractService, state)

	code, _ := vminfo.GetTestData("event")
	vm, err := exec.NewVirtualMachine(code, base.DEFAULT_VM_CONFIG, solver, nil)
	assert.NoError(t, err)

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

	contractService.EXPECT().Self().Return(base.AccountRef(common.BytesToAddress(address))).AnyTimes()
	contractService.EXPECT().GetGas().Return(uint64(0)).AnyTimes()
	contractService.EXPECT().CallValue().Return(big.NewInt(100)).AnyTimes()
	vmValue.EXPECT().Call(contractService, common.BytesToAddress(address),params,uint64(0),big.NewInt(100)).Return([]byte(nil),uint64(0),nil)

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
