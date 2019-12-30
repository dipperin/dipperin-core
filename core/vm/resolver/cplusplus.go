// Copyright 2019, Keychain Foundation Ltd.
// This file is part of the Dipperin-core library.
//
// The Dipperin-core library is free software: you can redistribute
// it and/or modify it under the terms of the GNU Lesser General Public License
// as published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// The Dipperin-core library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package resolver

// #cgo CXXFLAGS: -std=c++14
// #cgo windows LDFLAGS: -static-libgcc
// #include "print128.h"
import "C"
import (
	"encoding/binary"
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/math"
	"github.com/dipperin/dipperin-core/core/vm/common/utils"
	"github.com/dipperin/dipperin-core/third-party/crypto"
	"github.com/dipperin/dipperin-core/third-party/life/exec"
	"github.com/dipperin/dipperin-core/third-party/log"
	"math/big"
)

/*func PrintTest(){
	lo := uint64(1232)
	ho := uint64(0)
	ret := C.printi128(C.uint64_t(lo), C.uint64_t(ho))

	num := C.GoString(ret)
	fmt.Printf("envPrinti128 called result is:%v \r\n", num)
}*/

func envMemcpyGasCost(vm *exec.VirtualMachine) (uint64, error) {
	len := int(uint32(vm.GetCurrentFrame().Locals[2]))
	return uint64(len), nil
}

func envMemmoveGasCost(vm *exec.VirtualMachine) (uint64, error) {
	len := int(uint32(vm.GetCurrentFrame().Locals[2]))
	return uint64(len), nil
}

/*//int memcmp ( const void * ptr1, const void * ptr2, size_t num );
func envMemcmp(vm *exec.VirtualMachine) int64 {
	ptr1 := int(uint32(vm.GetCurrentFrame().Locals[0]))
	ptr2 := int(uint32(vm.GetCurrentFrame().Locals[1]))
	num := int(uint32(vm.GetCurrentFrame().Locals[2]))

	return int64(bytes.Compare(vm.Memory.Memory[ptr1:ptr1+num], vm.Memory.Memory[ptr2:ptr2+num]))
}

func envMemcmpGasCost(vm *exec.VirtualMachine) (uint64, error) {
	len := int(uint32(vm.GetCurrentFrame().Locals[2]))
	return uint64(len), nil
}*/

//void * memset ( void * ptr, int value, size_t num );
func envMemset(vm *exec.VirtualMachine) int64 {
	ptr := int(uint32(vm.GetCurrentFrame().Locals[0]))
	value := int(uint32(vm.GetCurrentFrame().Locals[1]))
	num := int(uint32(vm.GetCurrentFrame().Locals[2]))

	pos := 0
	for pos < num {
		vm.Memory.Memory[ptr+pos] = byte(value)
		pos++
	}
	return int64(ptr)
}

func envMemsetGasCost(vm *exec.VirtualMachine) (uint64, error) {
	len := int(uint32(vm.GetCurrentFrame().Locals[2]))
	return uint64(len), nil
}

//libc prints()
func (r *Resolver) envPrints(vm *exec.VirtualMachine) int64 {
	start := int(uint32(vm.GetCurrentFrame().Locals[0]))
	end := 0
	for end = start; end < len(vm.Memory.Memory); end++ {
		if vm.Memory.Memory[end] == 0 {
			break
		}
	}
	str := vm.Memory.Memory[start:end]
	log.Vm.Debug(string(str))
	log.Info("envPrints called", "string", string(str))
	return 0
}

func envPrintsGasCost(vm *exec.VirtualMachine) (uint64, error) {
	start := int(uint32(vm.GetCurrentFrame().Locals[0]))
	end := 0
	for end = start; end < len(vm.Memory.Memory); end++ {
		if vm.Memory.Memory[end] == 0 {
			break
		}
	}
	return uint64(end - start), nil
}

//libc prints_l
func envPrintsl(vm *exec.VirtualMachine) int64 {
	ptr := int(uint32(vm.GetCurrentFrame().Locals[0]))
	msgLen := int(uint32(vm.GetCurrentFrame().Locals[1]))
	msg := vm.Memory.Memory[ptr : ptr+msgLen]
	log.Vm.Debug(string(msg))
	log.Info("envPrintsl called", "string", string(msg))
	return 0
}

func envPrintslGasCost(vm *exec.VirtualMachine) (uint64, error) {
	msgLen := int(uint32(vm.GetCurrentFrame().Locals[1]))
	return uint64(msgLen), nil
}

//libc printi()
func envPrinti(vm *exec.VirtualMachine) int64 {
	num := vm.GetCurrentFrame().Locals[0]
	log.Vm.Debug(fmt.Sprintf("%d", num))
	log.Info("envPrinti called", "int", fmt.Sprintf("%d", num))
	return 0
}

func envPrintiGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}

func envPrintui(vm *exec.VirtualMachine) int64 {
	num := vm.GetCurrentFrame().Locals[0]
	log.Vm.Debug(fmt.Sprintf("%d", num))
	log.Info("envPrintui called", "uint", fmt.Sprintf("%d", num))
	return 0
}

func envPrintuiGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}

func envPrinti128(vm *exec.VirtualMachine) int64 {
	pos := vm.GetCurrentFrame().Locals[0]
	buf := vm.Memory.Memory[pos : pos+16]
	lo := uint64(binary.LittleEndian.Uint64(buf[:8]))
	ho := uint64(binary.LittleEndian.Uint64(buf[8:]))
	ret := C.printi128(C.uint64_t(lo), C.uint64_t(ho))

	num := C.GoString(ret)
	log.Vm.Debug(fmt.Sprintf("%s", num))
	log.Info("envPrinti128 called", "int128", fmt.Sprintf("%s", num))
	return 0
}

func envPrinti128GasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}

func envPrintui128(vm *exec.VirtualMachine) int64 {
	pos := vm.GetCurrentFrame().Locals[0]
	buf := vm.Memory.Memory[pos : pos+16]
	lo := uint64(binary.LittleEndian.Uint64(buf[:8]))
	ho := uint64(binary.LittleEndian.Uint64(buf[8:]))
	ret := C.printui128(C.uint64_t(lo), C.uint64_t(ho))

	num := C.GoString(ret)
	log.Vm.Debug(fmt.Sprintf("%s", num))
	log.Info("envPrintui128 called", "uint128", fmt.Sprintf("%s", num))
	return 0
}

func envPrintui128GasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}

/*func envPrintsf(vm *exec.VirtualMachine) int64 {
	pos := vm.GetCurrentFrame().Locals[0]
	float := math.Float32frombits(uint32(pos))
	log.Vm.Debug(fmt.Sprintf("%g", float))
	log.Info("envPrintsf called", "float", fmt.Sprintf("%g", float))
	return 0
}

func envPrintsfGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}

func envPrintdf(vm *exec.VirtualMachine) int64 {
	pos := vm.GetCurrentFrame().Locals[0]
	double := math.Float64frombits(uint64(pos))
	log.Vm.Debug(fmt.Sprintf("%g", double))
	log.Info("envPrintdf called", "double", fmt.Sprintf("%g", double))
	return 0
}

func envPrintdfGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}*/

/*func envPrintqf(vm *exec.VirtualMachine) int64 {
	frame := vm.GetCurrentFrame()
	pos := frame.Locals[0]

	low := C.uint64_t(binary.LittleEndian.Uint64(vm.Memory.Memory[pos : pos+8]))
	high := C.uint64_t(binary.LittleEndian.Uint64(vm.Memory.Memory[pos+8 : pos+16]))

	buf := C.GoString(C.__printqf(low, high))
	log.Vm.Debug(fmt.Sprintf("%s", buf))
	log.Info("envPrintqf called", "longDouble", fmt.Sprintf("%s", buf))
	return 0
}

func envPrintqfGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}
*/

func envPrintn(vm *exec.VirtualMachine) int64 {
	data := fmt.Sprintf("%d", int(uint32(vm.GetCurrentFrame().Locals[0])))
	log.Vm.Debug(data)
	log.Info("envPrintn called", "envPrintn", data)
	return 0
}

func envPrintnGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}

func envPrinthex(vm *exec.VirtualMachine) int64 {
	data := int(uint32(vm.GetCurrentFrame().Locals[0]))
	dataLen := int(uint32(vm.GetCurrentFrame().Locals[1]))
	hex := vm.Memory.Memory[data : data+dataLen]
	log.Vm.Debug(fmt.Sprintf("%x", hex))
	log.Info("envPrinthex called", "hex", fmt.Sprintf("%x", hex))
	return 0
}

func envPrinthexGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}

func envMallocGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 0, nil
}

func envFreeGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 0, nil
}

//libc calloc()
func envCalloc(vm *exec.VirtualMachine) int64 {
	mem := vm.Memory
	num := int(int32(vm.GetCurrentFrame().Locals[0]))
	size := int(int32(vm.GetCurrentFrame().Locals[1]))
	total := num * size

	pos := mem.Malloc(total)
	return int64(pos)
}

func envCallocGasCost(vm *exec.VirtualMachine) (uint64, error) {
	num := int(int32(vm.GetCurrentFrame().Locals[0]))
	size := int(int32(vm.GetCurrentFrame().Locals[1]))
	total := num * size
	return uint64(total), nil
}

func envRealloc(vm *exec.VirtualMachine) int64 {
	mem := vm.Memory
	size := int(int32(vm.GetCurrentFrame().Locals[1]))

	if size == 0 {
		return 0
	}

	pos := mem.Malloc(size)

	return int64(pos)
}

func envReallocGasCost(vm *exec.VirtualMachine) (uint64, error) {
	size := int(int32(vm.GetCurrentFrame().Locals[1]))
	return uint64(size), nil
}

func envAbort(vm *exec.VirtualMachine) int64 {
	panic("abort")
}

func envAbortGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 0, nil
}

// define: int64_t gasPrice();
func (r *Resolver) envGasPrice(vm *exec.VirtualMachine) int64 {
	gasPrice := r.Service.GetGasPrice().Int64()
	log.Info("envGasPrice", "price", gasPrice)
	return gasPrice
}

// define: void blockHash(int num, char hash[20]);
func (r *Resolver) envBlockHash(vm *exec.VirtualMachine) int64 {
	num := int(int32(vm.GetCurrentFrame().Locals[0]))
	offset := int(int32(vm.GetCurrentFrame().Locals[1]))
	blockHash := r.Service.GetBlockHash(uint64(num))
	copy(vm.Memory.Memory[offset:], blockHash.Bytes())
	return 0
}

// define: int64_t number();
func (r *Resolver) envNumber(vm *exec.VirtualMachine) int64 {
	num := int64(r.Service.GetBlockNumber().Uint64())
	log.Info("envNumber", "num", num)
	return num
}

// define: int64_t gasLimit();
func (r *Resolver) envGasLimit(vm *exec.VirtualMachine) int64 {
	gasLimit := int64(r.Service.GetGasLimit())
	log.Info("envGasLimit", "gasLimit", gasLimit)
	return gasLimit
}

// define: int64_t timestamp();
func (r *Resolver) envTimestamp(vm *exec.VirtualMachine) int64 {
	time := r.Service.GetTime().Int64()
	log.Info("envTimestamp", "time", time)
	return time
}

// define: void coinbase(char addr[22]);
func (r *Resolver) envCoinbase(vm *exec.VirtualMachine) int64 {
	offset := int(int32(vm.GetCurrentFrame().Locals[0]))
	coinBase := r.Service.GetCoinBase()
	copy(vm.Memory.Memory[offset:], coinBase.Bytes())
	return 0
}

// define: u256 balance();
func (r *Resolver) envBalance(vm *exec.VirtualMachine) int64 {
	addr := int(int32(vm.GetCurrentFrame().Locals[0]))
	log.Info("the currentFrame is:", "frame", vm.GetCurrentFrame(), "local", vm.GetCurrentFrame().Locals)
	addrLen := int(int32(vm.GetCurrentFrame().Locals[1]))
	ptr := int(int32(vm.GetCurrentFrame().Locals[2]))

	address := vm.Memory.Memory[addr : addr+addrLen]
	balance := r.Service.GetBalance(common.BytesToAddress(address))
	// 256 bits
	if len(balance.Bytes()) > 32 {
		panic(fmt.Sprintf("balance overflow(%d>32)", len(balance.Bytes())))
	}
	// bigendiangasLimit
	offset := 32 - len(balance.Bytes())
	if offset > 0 {
		empty := make([]byte, offset)
		copy(vm.Memory.Memory[ptr:ptr+offset], empty)
	}
	copy(vm.Memory.Memory[ptr+offset:], balance.Bytes())
	return 0
}

// define: void origin(char addr[22]);
func (r *Resolver) envOrigin(vm *exec.VirtualMachine) int64 {
	offset := int(int32(vm.GetCurrentFrame().Locals[0]))
	address := r.Service.GetOrigin()
	copy(vm.Memory.Memory[offset:], address.Bytes())
	return 0
}

// define: void caller(char addr[22]);
func (r *Resolver) envCaller(vm *exec.VirtualMachine) int64 {
	offset := int(int32(vm.GetCurrentFrame().Locals[0]))
	caller := r.Service.Caller().Address()
	log.Info("envCaller", "caller", caller)
	copy(vm.Memory.Memory[offset:], caller.Bytes())
	return 0
}

// define: int64_t callValue();
func (r *Resolver) envCallValue(vm *exec.VirtualMachine) int64 {
	value := r.Service.CallValue()
	ptr := int(int32(vm.GetCurrentFrame().Locals[0]))
	if len(value.Bytes()) > 32 {
		panic(fmt.Sprintf("balance overflow(%d > 32)", len(value.Bytes())))
	}
	// bigendian
	offset := 32 - len(value.Bytes())
	if offset > 0 {
		empty := make([]byte, offset)
		copy(vm.Memory.Memory[ptr:ptr+offset], empty)
	}
	copy(vm.Memory.Memory[ptr+offset:], value.Bytes())
	return 0
}

// define: void address(char hash[22]);
func (r *Resolver) envAddress(vm *exec.VirtualMachine) int64 {
	offset := int(int32(vm.GetCurrentFrame().Locals[0]))
	address := r.Service.Address()
	copy(vm.Memory.Memory[offset:], address.Bytes())
	return 0
}

// define: void sha3(char *src, size_t srcLen, char *dest, size_t destLen);
func envSha3(vm *exec.VirtualMachine) int64 {
	offset := int(int32(vm.GetCurrentFrame().Locals[0]))
	size := int(int32(vm.GetCurrentFrame().Locals[1]))
	destOffset := int(int32(vm.GetCurrentFrame().Locals[2]))
	destSize := int(int32(vm.GetCurrentFrame().Locals[3]))
	data := vm.Memory.Memory[offset : offset+size]
	hash := crypto.Keccak256(data)
	log.Info("envSha3 called", "hash", hash, "hasHex", common.Bytes2Hex(hash))
	if destSize < len(hash) {
		// todo
		return 0
	}
	//fmt.Printf("Sha3:%v, 0:%v, 1:%v, (-2):%v, (-1):%v. \n", common.Bytes2Hex(hash), hash[0], fmt.Sprintf("%b", hash[1]), hash[len(hash)-2], hash[len(hash)-1])
	copy(vm.Memory.Memory[destOffset:], hash)
	return 0
}

func envHexStringSameWithVM(vm *exec.VirtualMachine) int64 {
	log.Debug("envHexStringSameWithVM execute")
	offset := int(int32(vm.GetCurrentFrame().Locals[0]))
	size := int(int32(vm.GetCurrentFrame().Locals[1]))
	destOffset := int(int32(vm.GetCurrentFrame().Locals[2]))
	//destSize := int(int32(vm.GetCurrentFrame().Locals[3]))
	data := vm.Memory.Memory[offset : offset+size]
	str := common.HexStringSameWithVM(string(data))
	log.Info("envHexStringSameWithVM  ", "data", data, "str", str)
	copy(vm.Memory.Memory[destOffset:], str)
	return 0
}

func envSha3GasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}

func envHexStringSameWithVMGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}

func constGasFunc(gas uint64) exec.GasCost {
	return func(vm *exec.VirtualMachine) (uint64, error) {
		return gas, nil
	}
}

func envEmitEventGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}

func envSetStateGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}

func envGetStateGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}

func envGetStateSizeGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}

// define: int64_t getNonce();
func (r *Resolver) envGetCallerNonce(vm *exec.VirtualMachine) int64 {
	addr := r.Service.Caller().Address()
	nonce, _ := r.Service.StateDBService.GetNonce(addr)
	log.Info("envGetCallerNonce", "nonce", nonce)
	return int64(nonce)
}

/*func (r *Resolver) envCurrentTime(vm *exec.VirtualMachine) int64 {
	curTime := time.Now().UnixNano()
	log.Info("envCurrentTime", "time", curTime)
	return curTime
}*/

func (r *Resolver) envCallTransfer(vm *exec.VirtualMachine) int64 {
	key := int(int32(vm.GetCurrentFrame().Locals[0]))
	keyLen := int(int32(vm.GetCurrentFrame().Locals[1]))
	value := int(vm.GetCurrentFrame().Locals[2])
	bValue := new(big.Int)
	// 256 bits
	bValue.SetBytes(vm.Memory.Memory[value : value+32])
	value256 := math.U256(bValue)
	addr := common.BytesToAddress(vm.Memory.Memory[key : key+keyLen])
	_, returnGas, err := r.Service.Transfer(addr, value256)

	//先使用在life　vm中添加的字段，待后续看是否可以使用life自带gas机制
	log.Info("envCallTransfer", "GasUsed", vm.GasUsed, "returnGas", returnGas, "err", err)
	vm.GasUsed -= returnGas
	if err != nil {
		return 1
	} else {
		return 0
	}
}


func (r *Resolver) envCallTransferUDIP(vm *exec.VirtualMachine) int64 {
	key := int(int32(vm.GetCurrentFrame().Locals[0]))
	keyLen := int(int32(vm.GetCurrentFrame().Locals[1]))
	value := int64(vm.GetCurrentFrame().Locals[2])
	value256 := math.U256(new(big.Int).Mul(new(big.Int).SetInt64(value), math.BigPow(10,15)))
	addr := common.BytesToAddress(vm.Memory.Memory[key : key+keyLen])
	_, returnGas, err := r.Service.Transfer(addr, value256)

	//先使用在life　vm中添加的字段，待后续看是否可以使用life自带gas机制
	log.Info("envCallTransferUDIP", "GasUsed", vm.GasUsed, "returnGas", returnGas, "err", err, "transfer value", value256)
	vm.GasUsed -= returnGas
	if err != nil {
		return 1
	} else {
		return 0
	}
}

func (r *Resolver) envDipperCall(vm *exec.VirtualMachine) int64 {
	addr := int(int32(vm.GetCurrentFrame().Locals[0]))
	params := int(int32(vm.GetCurrentFrame().Locals[1]))
	paramsLen := int(int32(vm.GetCurrentFrame().Locals[2]))

	contractAddr := vm.Memory.Memory[addr : addr+common.AddressLength]
	inputs := vm.Memory.Memory[params : params+paramsLen]
	log.Info("envDipperCall", "contractAddr", contractAddr, "inputs", inputs)
	_, err := r.Service.ResolverCall(contractAddr, inputs)
	if err != nil {
		fmt.Printf("call error,%s", err.Error())
		return 0
	}
	return 0
}

func (r *Resolver) envDipperDelegateCall(vm *exec.VirtualMachine) int64 {
	addr := int(int32(vm.GetCurrentFrame().Locals[0]))
	params := int(int32(vm.GetCurrentFrame().Locals[1]))
	paramsLen := int(int32(vm.GetCurrentFrame().Locals[2]))

	contractAddr := vm.Memory.Memory[addr : addr+common.AddressLength]
	inputs := vm.Memory.Memory[params : params+paramsLen]
	log.Info("envDipperDelegateCall", "contractAddr", contractAddr, "inputs", inputs)
	_, err := r.Service.ResolverDelegateCall(contractAddr, inputs)
	if err != nil {
		fmt.Printf("call error,%s", err.Error())
		return 0
	}
	return 0
}

func (r *Resolver) envDipperCallInt64(vm *exec.VirtualMachine) int64 {
	addr := int(int32(vm.GetCurrentFrame().Locals[0]))
	params := int(int32(vm.GetCurrentFrame().Locals[1]))
	paramsLen := int(int32(vm.GetCurrentFrame().Locals[2]))

	contractAddr := vm.Memory.Memory[addr : addr+common.AddressLength]
	inputs := vm.Memory.Memory[params : params+paramsLen]
	log.Info("envDipperCallInt64", "contractAddr", contractAddr, "inputs", inputs)
	ret, err := r.Service.ResolverCall(contractAddr, inputs)
	if err != nil {
		fmt.Printf("call error,%s", err.Error())
		return 0
	}
	res := utils.Align32BytesConverter(ret, "int64")
	log.Info("envDipperCallInt64", "ret", res)
	return res.(int64)
}

func (r *Resolver) envDipperDelegateCallInt64(vm *exec.VirtualMachine) int64 {
	addr := int(int32(vm.GetCurrentFrame().Locals[0]))
	params := int(int32(vm.GetCurrentFrame().Locals[1]))
	paramsLen := int(int32(vm.GetCurrentFrame().Locals[2]))

	contractAddr := vm.Memory.Memory[addr : addr+common.AddressLength]
	inputs := vm.Memory.Memory[params : params+paramsLen]
	log.Info("envDipperDelegateCallInt64", "contractAddr", contractAddr, "inputs", inputs)
	ret, err := r.Service.ResolverDelegateCall(contractAddr, inputs)
	if err != nil {
		fmt.Printf("call error,%s", err.Error())
		return 0
	}
	res := utils.Align32BytesConverter(ret, "int64")
	log.Info("envDipperDelegateCallInt64", "ret", res)
	return res.(int64)
}

func (r *Resolver) envDipperCallString(vm *exec.VirtualMachine) int64 {
	addr := int(int32(vm.GetCurrentFrame().Locals[0]))
	params := int(int32(vm.GetCurrentFrame().Locals[1]))
	paramsLen := int(int32(vm.GetCurrentFrame().Locals[2]))

	contractAddr := vm.Memory.Memory[addr : addr+common.AddressLength]
	inputs := vm.Memory.Memory[params : params+paramsLen]
	log.Info("envDipperCallString", "contractAddr", contractAddr, "inputs", inputs)
	ret, err := r.Service.ResolverCall(contractAddr, inputs)
	if err != nil {
		fmt.Printf("call error,%s", err.Error())
		return 0
	}
	res := utils.Align32BytesConverter(ret, "string")
	log.Info("envDipperCallString", "ret", res)
	return MallocString(vm, string(ret))
}

func (r *Resolver) envDipperDelegateCallString(vm *exec.VirtualMachine) int64 {
	addr := int(int32(vm.GetCurrentFrame().Locals[0]))
	params := int(int32(vm.GetCurrentFrame().Locals[1]))
	paramsLen := int(int32(vm.GetCurrentFrame().Locals[2]))

	contractAddr := vm.Memory.Memory[addr : addr+common.AddressLength]
	inputs := vm.Memory.Memory[params : params+paramsLen]
	log.Info("envDipperDelegateCallString", "contractAddr", contractAddr, "inputs", inputs)
	ret, err := r.Service.ResolverDelegateCall(contractAddr, inputs)
	if err != nil {
		fmt.Printf("call error,%s", err.Error())
		return 0
	}
	res := utils.Align32BytesConverter(ret, "string")
	log.Info("envDipperDelegateCallString", "ret", res)
	return MallocString(vm, string(ret))
}

func envDipperCallGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}

func envDipperCallInt64GasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}

func envDipperCallStringGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}
