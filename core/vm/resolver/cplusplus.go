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
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/addressutil"
	"github.com/dipperin/dipperin-core/common/log"
	"github.com/dipperin/dipperin-core/common/math"
	"github.com/dipperin/dipperin-core/core/vm/base/utils"
	"github.com/dipperin/dipperin-core/third_party/crypto"
	"github.com/dipperin/dipperin-core/third_party/life/exec"
	"go.uber.org/zap"
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
// todo
func (r *Resolver) envPrints(vm *exec.VirtualMachine) int64 {
	start := int(uint32(vm.GetCurrentFrame().Locals[0]))
	end := 0
	for end = start; end < len(vm.Memory.Memory); end++ {
		if vm.Memory.Memory[end] == 0 {
			break
		}
	}
	str := vm.Memory.Memory[start:end]
	log.DLogger.Debug(string(str))
	log.DLogger.Info("envPrints called", zap.String("string", string(str)))
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
// todo
func envPrintsl(vm *exec.VirtualMachine) int64 {
	ptr := int(uint32(vm.GetCurrentFrame().Locals[0]))
	msgLen := int(uint32(vm.GetCurrentFrame().Locals[1]))
	msg := vm.Memory.Memory[ptr : ptr+msgLen]
	log.DLogger.Debug(string(msg))
	log.DLogger.Info("envPrintsl called", zap.String("string", string(msg)))
	return 0
}


func envPrintslGasCost(vm *exec.VirtualMachine) (uint64, error) {
	msgLen := int(uint32(vm.GetCurrentFrame().Locals[1]))
	return uint64(msgLen), nil
}

//libc printi()
// todo
func envPrinti(vm *exec.VirtualMachine) int64 {
	num := vm.GetCurrentFrame().Locals[0]
	log.DLogger.Debug(fmt.Sprintf("%d", num))
	log.DLogger.Info("envPrinti called", zap.Int64("int", num))
	return 0
}

func envPrintiGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}

// todo
func envPrintui(vm *exec.VirtualMachine) int64 {
	num := vm.GetCurrentFrame().Locals[0]
	log.DLogger.Debug(fmt.Sprintf("%d", num))
	log.DLogger.Info("envPrintui called", zap.Int64("uint", num))
	return 0
}

func envPrintuiGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}

/*func envPrinti128(vm *exec.VirtualMachine) int64 {
	pos := vm.GetCurrentFrame().Locals[0]
	buf := vm.Memory.Memory[pos : pos+16]
	lo := binary.LittleEndian.Uint64(buf[:8])
	ho := binary.LittleEndian.Uint64(buf[8:])
	ret := C.printi128(C.uint64_t(lo), C.uint64_t(ho))

	num := C.GoString(ret)
	log.DLogger.Debug(fmt.Sprintf("%s", num))
	log.DLogger.Info("envPrinti128 called", zap.String("int128", num))
	return 0
}

func envPrinti128GasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}

func envPrintui128(vm *exec.VirtualMachine) int64 {
	pos := vm.GetCurrentFrame().Locals[0]
	buf := vm.Memory.Memory[pos : pos+16]
	lo := binary.LittleEndian.Uint64(buf[:8])
	ho := binary.LittleEndian.Uint64(buf[8:])
	ret := C.printui128(C.uint64_t(lo), C.uint64_t(ho))

	num := C.GoString(ret)
	log.DLogger.Debug(fmt.Sprintf("%s", num))
	log.DLogger.Info("envPrintui128 called", zap.String("uint128", num))
	return 0
}

func envPrintui128GasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}*/

/*func envPrintsf(vm *exec.VirtualMachine) int64 {
	pos := vm.GetCurrentFrame().Locals[0]
	float := math.Float32frombits(uint32(pos))
	log.DLogger.Debug(fmt.Sprintf("%g", float))
	log.DLogger.Info("envPrintsf called", "float", fmt.Sprintf("%g", float))
	return 0
}

func envPrintsfGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}

func envPrintdf(vm *exec.VirtualMachine) int64 {
	pos := vm.GetCurrentFrame().Locals[0]
	double := math.Float64frombits(uint64(pos))
	log.DLogger.Debug(fmt.Sprintf("%g", double))
	log.DLogger.Info("envPrintdf called", "double", fmt.Sprintf("%g", double))
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
	log.DLogger.Debug(fmt.Sprintf("%s", buf))
	log.DLogger.Info("envPrintqf called", "longDouble", fmt.Sprintf("%s", buf))
	return 0
}

func envPrintqfGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}
*/

// todo
func envPrintn(vm *exec.VirtualMachine) int64 {
	data := fmt.Sprintf("%d", int(uint32(vm.GetCurrentFrame().Locals[0])))
	log.DLogger.Debug(data)
	log.DLogger.Info("envPrintn called", zap.String("envPrintn", data))
	return 0
}

func envPrintnGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}

// todo
func envPrinthex(vm *exec.VirtualMachine) int64 {
	data := int(uint32(vm.GetCurrentFrame().Locals[0]))
	dataLen := int(uint32(vm.GetCurrentFrame().Locals[1]))
	hex := vm.Memory.Memory[data : data+dataLen]
	log.DLogger.Debug(fmt.Sprintf("%x", hex))
	log.DLogger.Info("envPrinthex called", zap.Uint8s("hex", hex))
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
// todo  a little troublesome
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

// todo  a little troublesome
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
	log.DLogger.Info("envGasPrice", zap.Int64("price", gasPrice))
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
	log.DLogger.Info("envNumber", zap.Int64("num", num))
	return num
}

// define: int64_t gasLimit();
func (r *Resolver) envGasLimit(vm *exec.VirtualMachine) int64 {
	gasLimit := int64(r.Service.GetGasLimit())
	log.DLogger.Info("envGasLimit", zap.Int64("gasLimit", gasLimit))
	return gasLimit
}

// define: int64_t timestamp();
func (r *Resolver) envTimestamp(vm *exec.VirtualMachine) int64 {
	time := r.Service.GetTime().Int64()
	log.DLogger.Info("envTimestamp", zap.Int64("time", time))
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
// todo
func (r *Resolver) envBalance(vm *exec.VirtualMachine) int64 {
	addr := int(int32(vm.GetCurrentFrame().Locals[0]))
	log.DLogger.Info("the currentFrame is:", zap.Any("frame", vm.GetCurrentFrame()), zap.Int64s("local", vm.GetCurrentFrame().Locals))
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


//  todo  it seems duplicate with envCaller
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
	log.DLogger.Info("envCaller", zap.Any("caller", caller))
	copy(vm.Memory.Memory[offset:], caller.Bytes())
	return 0
}

// define: int64_t callValue();
// todo
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

func (r *Resolver) envCallValueUDIP(vm *exec.VirtualMachine) int64 {
	value := r.Service.CallValue()
	result := int64(new(big.Int).Div(value, new(big.Int).Set(math.BigPow(10, 15))).Uint64())
	log.DLogger.Debug("envCallValueUDIP", zap.Any("value", value), zap.Int64("result", result))
	return result
}

// define: void address(char hash[22]);
func (r *Resolver) envAddress(vm *exec.VirtualMachine) int64 {
	offset := int(int32(vm.GetCurrentFrame().Locals[0]))
	address := r.Service.Address()
	copy(vm.Memory.Memory[offset:], address.Bytes())
	return 0
}

// define: void sha3(char *src, size_t srcLen, char *dest, size_t destLen);
func (r *Resolver) envSha3(vm *exec.VirtualMachine) int64 {
	offset := int(int32(vm.GetCurrentFrame().Locals[0]))
	size := int(int32(vm.GetCurrentFrame().Locals[1]))
	destOffset := int(int32(vm.GetCurrentFrame().Locals[2]))
	destSize := int(int32(vm.GetCurrentFrame().Locals[3]))
	data := vm.Memory.Memory[offset : offset+size]
	hash := crypto.Keccak256(data)
	log.DLogger.Info("envSha3 called", zap.Uint8s("hash", hash), zap.String("hasHex", common.Bytes2Hex(hash)), zap.Int("destSize", destSize), zap.Int("hash len", len(hash)))
	if destSize < len(hash) {
		// todo
		return 1
	}
	//fmt.Printf("Sha3:%v, 0:%v, 1:%v, (-2):%v, (-1):%v. \n", common.Bytes2Hex(hash), hash[0], fmt.Sprintf("%b", hash[1]), hash[len(hash)-2], hash[len(hash)-1])
	copy(vm.Memory.Memory[destOffset:], hash)
	return 0
}

func (r *Resolver) envHexStringSameWithVM(vm *exec.VirtualMachine) int64 {
	log.DLogger.Debug("envHexStringSameWithVM execute")
	offset := int(int32(vm.GetCurrentFrame().Locals[0]))
	size := int(int32(vm.GetCurrentFrame().Locals[1]))
	destOffset := int(int32(vm.GetCurrentFrame().Locals[2]))
	//destSize := int(int32(vm.GetCurrentFrame().Locals[3]))
	data := vm.Memory.Memory[offset : offset+size]
	str := common.HexStringSameWithVM(string(data))
	log.DLogger.Info("envHexStringSameWithVM  ", zap.Uint8s("data", data), zap.String("str", str))
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

func env__ashlti3GasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}

func env__multi3GasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}

func env__divti3GasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}

// define: int64_t getNonce();
func (r *Resolver) envGetCallerNonce(vm *exec.VirtualMachine) int64 {
	addr := r.Service.Caller().Address()
	nonce, _ := r.Service.StateDBService.GetNonce(addr)
	log.DLogger.Info("envGetCallerNonce", zap.Uint64("nonce", nonce))
	return int64(nonce)
}

/*func (r *Resolver) envCurrentTime(vm *exec.VirtualMachine) int64 {
	curTime := time.Now().UnixNano()
	log.DLogger.Info("envCurrentTime", "time", curTime)
	return curTime
}*/

// todo  is it reasonable to call the VM.Call() method to implement Transfer?
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
	log.DLogger.Info("envCallTransfer", zap.Uint64("GasUsed", vm.GasUsed), zap.Uint64("returnGas", returnGas), zap.Error(err))
	vm.GasUsed -= returnGas
	if err != nil {
		return 1
	} else {
		return 0
	}
}

func (r *Resolver) envGetSignerAddress(vm *exec.VirtualMachine) int64 {
	sha3DataStart := int(int32(vm.GetCurrentFrame().Locals[0]))
	sha3DataLen := int(int32(vm.GetCurrentFrame().Locals[1]))
	signatureStart := int(int32(vm.GetCurrentFrame().Locals[2]))
	signatureLen := int(int32(vm.GetCurrentFrame().Locals[3]))
	returnStart := int(int32(vm.GetCurrentFrame().Locals[4]))

	sha3Data := vm.Memory.Memory[sha3DataStart : sha3DataStart+sha3DataLen]
	signature := vm.Memory.Memory[signatureStart : signatureStart+signatureLen]

	//crypto.VerifySignature(r.Service.Self().Address().)

	log.DLogger.Info("Resolver#envVerifySignature", zap.String("sha3Data hex", common.Bytes2Hex(sha3Data)), zap.String("signature",   string(signature)))

	//sha3Byte := common.Hex2Bytes(sha3Data)
	signByte := common.Hex2Bytes(string(signature))
	log.DLogger.Info("signByte len",zap.Int("sha3 byte len", len(sha3Data)), zap.Int("hash byte len", len(signByte)))
	pK, err := crypto.SigToPub(sha3Data, signByte)
	if err != nil {
		log.DLogger.Error("Sig To Pub err ", zap.Error( err))
		return 1
	}

	log.DLogger.Info("Resolver#envVerifySignature addr", zap.Any("address ", addressutil.PubKeyToAddress(*pK, common.AddressTypeNormal)), zap.String("self address", r.Service.Self().Address().Hex()))

	addr := addressutil.PubKeyToAddress(*pK, common.AddressTypeNormal)
	copy(vm.Memory.Memory[returnStart:], addr.Bytes())
	return 0
}

// todo  is it reasonable to call the VM.Call() method to implement Transfer?
func (r *Resolver) envCallTransferUDIP(vm *exec.VirtualMachine) int64 {
	key := int(int32(vm.GetCurrentFrame().Locals[0]))
	keyLen := int(int32(vm.GetCurrentFrame().Locals[1]))
	value := int64(vm.GetCurrentFrame().Locals[2])
	value256 := math.U256(new(big.Int).Mul(new(big.Int).SetInt64(value), math.BigPow(10, 15)))
	addr := common.BytesToAddress(vm.Memory.Memory[key : key+keyLen])
	_, returnGas, err := r.Service.Transfer(addr, value256)

	//先使用在life　vm中添加的字段，待后续看是否可以使用life自带gas机制
	log.DLogger.Info("envCallTransferUDIP", zap.Uint64("GasUsed", vm.GasUsed), zap.Uint64("returnGas", returnGas), zap.Error(err), zap.Any("transfer value", value256))
	vm.GasUsed -= returnGas
	if err != nil {
		return 1
	} else {
		return 0
	}
}

// todo
func (r *Resolver) envDipperCall(vm *exec.VirtualMachine) int64 {
	addr := int(int32(vm.GetCurrentFrame().Locals[0]))
	params := int(int32(vm.GetCurrentFrame().Locals[1]))
	paramsLen := int(int32(vm.GetCurrentFrame().Locals[2]))

	contractAddr := vm.Memory.Memory[addr : addr+common.AddressLength]
	inputs := vm.Memory.Memory[params : params+paramsLen]
	log.DLogger.Info("envDipperCall", zap.Uint8s("contractAddr", contractAddr), zap.Uint8s("inputs", inputs))
	_, err := r.Service.ResolverCall(contractAddr, inputs)
	if err != nil {
		fmt.Printf("call error,%s", err.Error())
		return 0
	}
	return 0
}

// todo
func (r *Resolver) envDipperDelegateCall(vm *exec.VirtualMachine) int64 {
	addr := int(int32(vm.GetCurrentFrame().Locals[0]))
	params := int(int32(vm.GetCurrentFrame().Locals[1]))
	paramsLen := int(int32(vm.GetCurrentFrame().Locals[2]))

	contractAddr := vm.Memory.Memory[addr : addr+common.AddressLength]
	inputs := vm.Memory.Memory[params : params+paramsLen]
	log.DLogger.Info("envDipperDelegateCall", zap.Uint8s("contractAddr", contractAddr), zap.Uint8s("inputs", inputs))
	_, err := r.Service.ResolverDelegateCall(contractAddr, inputs)
	if err != nil {
		fmt.Printf("call error,%s", err.Error())
		return 0
	}
	return 0
}

// todo
func (r *Resolver) envDipperCallInt64(vm *exec.VirtualMachine) int64 {
	addr := int(int32(vm.GetCurrentFrame().Locals[0]))
	params := int(int32(vm.GetCurrentFrame().Locals[1]))
	paramsLen := int(int32(vm.GetCurrentFrame().Locals[2]))

	contractAddr := vm.Memory.Memory[addr : addr+common.AddressLength]
	inputs := vm.Memory.Memory[params : params+paramsLen]
	log.DLogger.Info("envDipperCallInt64", zap.Uint8s("contractAddr", contractAddr), zap.Uint8s("inputs", inputs))
	ret, err := r.Service.ResolverCall(contractAddr, inputs)
	if err != nil {
		fmt.Printf("call error,%s", err.Error())
		return 0
	}
	res := utils.Align32BytesConverter(ret, "int64")
	log.DLogger.Info("envDipperCallInt64", zap.Any("ret", res))
	return res.(int64)
}

// todo
func (r *Resolver) envDipperDelegateCallInt64(vm *exec.VirtualMachine) int64 {
	addr := int(int32(vm.GetCurrentFrame().Locals[0]))
	params := int(int32(vm.GetCurrentFrame().Locals[1]))
	paramsLen := int(int32(vm.GetCurrentFrame().Locals[2]))

	contractAddr := vm.Memory.Memory[addr : addr+common.AddressLength]
	inputs := vm.Memory.Memory[params : params+paramsLen]
	log.DLogger.Info("envDipperDelegateCallInt64", zap.Uint8s("contractAddr", contractAddr), zap.Uint8s("inputs", inputs))
	ret, err := r.Service.ResolverDelegateCall(contractAddr, inputs)
	if err != nil {
		fmt.Printf("call error,%s", err.Error())
		return 0
	}
	res := utils.Align32BytesConverter(ret, "int64")
	log.DLogger.Info("envDipperDelegateCallInt64", zap.Any("ret", res))
	return res.(int64)
}


// todo
func (r *Resolver) envDipperCallString(vm *exec.VirtualMachine) int64 {
	addr := int(int32(vm.GetCurrentFrame().Locals[0]))
	params := int(int32(vm.GetCurrentFrame().Locals[1]))
	paramsLen := int(int32(vm.GetCurrentFrame().Locals[2]))

	contractAddr := vm.Memory.Memory[addr : addr+common.AddressLength]
	inputs := vm.Memory.Memory[params : params+paramsLen]
	log.DLogger.Info("envDipperCallString", zap.Uint8s("contractAddr", contractAddr), zap.Uint8s("inputs", inputs))
	ret, err := r.Service.ResolverCall(contractAddr, inputs)
	if err != nil {
		fmt.Printf("call error,%s", err.Error())
		return 0
	}
	res := utils.Align32BytesConverter(ret, "string")
	log.DLogger.Info("envDipperCallString", zap.Any("ret", res))
	return MallocString(vm, string(ret))
}

// todo
func (r *Resolver) envDipperDelegateCallString(vm *exec.VirtualMachine) int64 {
	addr := int(int32(vm.GetCurrentFrame().Locals[0]))
	params := int(int32(vm.GetCurrentFrame().Locals[1]))
	paramsLen := int(int32(vm.GetCurrentFrame().Locals[2]))

	contractAddr := vm.Memory.Memory[addr : addr+common.AddressLength]
	inputs := vm.Memory.Memory[params : params+paramsLen]
	log.DLogger.Info("envDipperDelegateCallString", zap.Uint8s("contractAddr", contractAddr), zap.Uint8s("inputs", inputs))
	ret, err := r.Service.ResolverDelegateCall(contractAddr, inputs)
	if err != nil {
		fmt.Printf("call error,%s", err.Error())
		return 0
	}
	res := utils.Align32BytesConverter(ret, "string")
	log.DLogger.Info("envDipperDelegateCallString", zap.Any("ret", res))
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
