package resolver

// #cgo CFLAGS: -I../softfloat/source/include
// #define SOFTFLOAT_FAST_INT64
// #include "softfloat.h"
//
// #cgo CXXFLAGS: -std=c++14
// #include "printqf.h"
// #include "print128.h"
import "C"
import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/dipperin/dipperin-core/third-party/crypto"
	"github.com/dipperin/dipperin-core/third-party/life/exec"
	"math"
	"math/big"
)



//void * memcpy ( void * destination, const void * source, size_t num );
func envMemcpy(vm *exec.VirtualMachine) int64 {
	dest := int(uint32(vm.GetCurrentFrame().Locals[0]))
	src := int(uint32(vm.GetCurrentFrame().Locals[1]))
	len := int(uint32(vm.GetCurrentFrame().Locals[2]))

	copy(vm.Memory.Memory[dest:dest+len], vm.Memory.Memory[src:src+len])
	return int64(dest)
}

func envMemcpyGasCost(vm *exec.VirtualMachine) (uint64, error) {
	len := int(uint32(vm.GetCurrentFrame().Locals[2]))
	return uint64(len), nil
}

//void * memmove ( void * destination, const void * source, size_t num );
func envMemmove(vm *exec.VirtualMachine) int64 {
	dest := int(uint32(vm.GetCurrentFrame().Locals[0]))
	src := int(uint32(vm.GetCurrentFrame().Locals[1]))
	len := int(uint32(vm.GetCurrentFrame().Locals[2]))

	copy(vm.Memory.Memory[dest:dest+len], vm.Memory.Memory[src:src+len])
	return int64(dest)
}

func envMemmoveGasCost(vm *exec.VirtualMachine) (uint64, error) {
	len := int(uint32(vm.GetCurrentFrame().Locals[2]))
	return uint64(len), nil
}

//int memcmp ( const void * ptr1, const void * ptr2, size_t num );
func envMemcmp(vm *exec.VirtualMachine) int64 {
	ptr1 := int(uint32(vm.GetCurrentFrame().Locals[0]))
	ptr2 := int(uint32(vm.GetCurrentFrame().Locals[1]))
	num := int(uint32(vm.GetCurrentFrame().Locals[2]))

	return int64(bytes.Compare(vm.Memory.Memory[ptr1:ptr1+num], vm.Memory.Memory[ptr2:ptr2+num]))
}

func envMemcmpGasCost(vm *exec.VirtualMachine) (uint64, error) {
	len := int(uint32(vm.GetCurrentFrame().Locals[2]))
	return uint64(len), nil
}

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
func (r *Resolver)envPrints(vm *exec.VirtualMachine) int64 {
	start := int(uint32(vm.GetCurrentFrame().Locals[0]))
	end := 0
	for end = start; end < len(vm.Memory.Memory); end++ {
		if vm.Memory.Memory[end] == 0 {
			break
		}
	}
	r.context.Log.Debug(string(vm.Memory.Memory[start:end]))

	//fmt.Printf("%s", string(vm.Memory.Memory[start:end]))
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
func (r *Resolver)envPrintsl(vm *exec.VirtualMachine) int64 {
	ptr := int(uint32(vm.GetCurrentFrame().Locals[0]))
	msgLen := int(uint32(vm.GetCurrentFrame().Locals[1]))
	msg := vm.Memory.Memory[ptr : ptr+msgLen]
	r.context.Log.Debug(string(msg))
	return 0
}

func envPrintslGasCost(vm *exec.VirtualMachine) (uint64, error) {
	msgLen := int(uint32(vm.GetCurrentFrame().Locals[1]))
	return uint64(msgLen), nil
}

//libc printi()
func (r *Resolver)envPrinti(vm *exec.VirtualMachine) int64 {
	r.context.Log.Debug(fmt.Sprintf("%d", vm.GetCurrentFrame().Locals[0]))
	return 0
}

func envPrintiGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}

func (r *Resolver)envPrintui(vm *exec.VirtualMachine) int64 {
	r.context.Log.Debug(fmt.Sprintf("%d", vm.GetCurrentFrame().Locals[0]))
	return 0
}

func envPrintuiGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}

func (r *Resolver)envPrinti128(vm *exec.VirtualMachine) int64 {
	pos := vm.GetCurrentFrame().Locals[0]
	buf := vm.Memory.Memory[pos : pos+16]
	lo := uint64(binary.LittleEndian.Uint64(buf[:8]))
	ho := uint64(binary.LittleEndian.Uint64(buf[8:]))
	ret := C.printi128(C.uint64_t(lo), C.uint64_t(ho))
	r.context.Log.Debug(fmt.Sprintf("%s", C.GoString(ret)))
	return 0
}

func envPrinti128GasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}

func (r *Resolver)envPrintui128(vm *exec.VirtualMachine) int64 {
	pos := vm.GetCurrentFrame().Locals[0]
	buf := vm.Memory.Memory[pos : pos+16]
	lo := uint64(binary.LittleEndian.Uint64(buf[:8]))
	ho := uint64(binary.LittleEndian.Uint64(buf[8:]))
	ret := C.printui128(C.uint64_t(lo), C.uint64_t(ho))
	r.context.Log.Debug(fmt.Sprintf("%s", C.GoString(ret)))
	return 0
}

func envPrintui128GasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}

func (r *Resolver)envPrintsf(vm *exec.VirtualMachine) int64 {
	pos := vm.GetCurrentFrame().Locals[0]
	float := math.Float32frombits(uint32(pos))
	r.context.Log.Debug(fmt.Sprintf("%g", float))
	return 0
}

func envPrintsfGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}

func (r *Resolver)envPrintdf(vm *exec.VirtualMachine) int64 {
	pos := vm.GetCurrentFrame().Locals[0]
	double := math.Float64frombits(uint64(pos))
	r.context.Log.Debug(fmt.Sprintf("%g", double))
	return 0
}

func envPrintdfGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}

func (r *Resolver)envPrintqf(vm *exec.VirtualMachine) int64 {
	frame := vm.GetCurrentFrame()
	pos := frame.Locals[0]

	low := C.uint64_t(binary.LittleEndian.Uint64(vm.Memory.Memory[pos : pos+8]))
	high := C.uint64_t(binary.LittleEndian.Uint64(vm.Memory.Memory[pos+8 : pos+16]))

	buf := C.GoString(C.__printqf(low, high))
	r.context.Log.Debug(fmt.Sprintf("%s", buf))
	return 0
}

func envPrintqfGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}

func (r *Resolver)envPrintn(vm *exec.VirtualMachine) int64 {
	r.context.Log.Debug(fmt.Sprintf("%d", int(uint32(vm.GetCurrentFrame().Locals[0]))))
	return 0
}

func envPrintnGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}

func (r *Resolver)envPrinthex(vm *exec.VirtualMachine) int64 {
	data := int(uint32(vm.GetCurrentFrame().Locals[0]))
	dataLen := int(uint32(vm.GetCurrentFrame().Locals[1]))
	r.context.Log.Debug(fmt.Sprintf("%x", vm.Memory.Memory[data:data+dataLen]))
	return 0
}

func envPrinthexGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}

//libc malloc()
func envMalloc(vm *exec.VirtualMachine) int64 {
	//mem := vm.Memory
	size := int(uint32(vm.GetCurrentFrame().Locals[0]))

	pos := vm.Memory.Malloc(size)
	if pos == -1 {
		panic("melloc error...")
	}

	return int64(pos)
}

func envMallocGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 0, nil
}

//libc free()
func envFree(vm *exec.VirtualMachine) int64 {
	mem := vm.Memory
	offset := int(uint32(vm.GetCurrentFrame().Locals[0]))

	err := mem.Free(offset)
	if err != nil {
		panic("free error...")
	}

	return 0
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
	//ptr := int(int32(vm.GetCurrentFrame().Locals[0]))
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
func (r *Resolver)envGasPrice(vm *exec.VirtualMachine) int64 {
	gasPrice := r.stateDb.GasPrice()
	return gasPrice
}

// define: void blockHash(int num, char hash[20]);
func (r *Resolver)envBlockHash(vm *exec.VirtualMachine) int64 {
	num := int(int32(vm.GetCurrentFrame().Locals[0]))
	offset := int(int32(vm.GetCurrentFrame().Locals[1]))
	blockHash := r.stateDb.BlockHash(uint64(num))
	//fmt.Printf("Number:%v ,Num:%v ,0:%v, 1:%v, (-2):%v, (-1):%v. \n", num, blockHash.Hex(), " -> ", blockHash[0], blockHash[1], blockHash[len(blockHash)-2], blockHash[len(blockHash)-1])
	copy(vm.Memory.Memory[offset:], blockHash.Bytes())
	return 0
}

// define: int64_t number();
func (r *Resolver)envNumber(vm *exec.VirtualMachine) int64 {
	return int64(r.context.StateDB.BlockNumber().Uint64())
}

// define: int64_t gasLimit();
func (r *Resolver)envGasLimit(vm *exec.VirtualMachine) int64 {
	return int64(r.context.StateDB.GasLimimt())
}

// define: int64_t timestamp();
func (r *Resolver)envTimestamp(vm *exec.VirtualMachine) int64 {
	return r.context.StateDB.Time().Int64()
}

// define: void coinbase(char addr[20]);
func (r *Resolver)envCoinbase(vm *exec.VirtualMachine) int64 {
	offset := int(int32(vm.GetCurrentFrame().Locals[0]))
	coinBase := r.context.StateDB.Coinbase()
	//fmt.Println("CoinBase:", coinBase.Hex(), " -> ", coinBase[0], coinBase[1], coinBase[len(coinBase)-2], coinBase[len(coinBase)-1])
	copy(vm.Memory.Memory[offset:], coinBase.Bytes())
	return 0
}

// define: u256 balance();
func (r *Resolver)envBalance(vm *exec.VirtualMachine) int64 {
	balance := r.context.StateDB.GetBalance(r.context.StateDB.Address())
	ptr := int(int32(vm.GetCurrentFrame().Locals[0]))
	// 256 bits
	if len(balance.Bytes()) > 32 {
		panic(fmt.Sprintf("balance overflow(%d>32)", len(balance.Bytes())))
	}
	// bigendian
	offset := 32 - len(balance.Bytes())
	if offset > 0 {
		empty := make([]byte, offset)
		copy(vm.Memory.Memory[ptr:ptr+offset], empty)
	}
	copy(vm.Memory.Memory[ptr+offset:], balance.Bytes())
	return 0
}

// define: void origin(char addr[20]);
func (r *Resolver)envOrigin(vm *exec.VirtualMachine) int64 {
	offset := int(int32(vm.GetCurrentFrame().Locals[0]))
	address := r.context.StateDB.Origin()
	//fmt.Println("Origin:", address.Hex(), " -> ", address[0], address[1], address[len(address)-2], address[len(address)-1])
	copy(vm.Memory.Memory[offset:], address.Bytes())
	return 0
}

// define: void caller(char addr[20]);
func (r *Resolver)envCaller(vm *exec.VirtualMachine) int64 {
	offset := int(int32(vm.GetCurrentFrame().Locals[0]))
	caller := r.context.StateDB.Caller()
	//fmt.Println("Caller:", caller.Hex(), " -> ", caller[0], caller[1], caller[len(caller)-2], caller[len(caller)-1])
	copy(vm.Memory.Memory[offset:], caller.Bytes())
	return 0
}

// define: int64_t callValue();
func (r *Resolver)envCallValue(vm *exec.VirtualMachine) int64 {
	value := r.context.StateDB.CallValue()
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

// define: void address(char hash[20]);
func (r *Resolver)envAddress(vm *exec.VirtualMachine) int64 {
	offset := int(int32(vm.GetCurrentFrame().Locals[0]))
	address := r.context.StateDB.Address()
	//fmt.Println("Address:", address.Hex(), " -> ", address[0], address[1], address[len(address)-2], address[len(address)-1])
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
	//fmt.Println(common.Bytes2Hex(hash))
	if destSize < len(hash) {
		return 0
	}
	//fmt.Printf("Sha3:%v, 0:%v, 1:%v, (-2):%v, (-1):%v. \n", common.Bytes2Hex(hash), hash[0], fmt.Sprintf("%b", hash[1]), hash[len(hash)-2], hash[len(hash)-1])
	copy(vm.Memory.Memory[destOffset:], hash)
	return 0
}

func envSha3GasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}

func constGasFunc(gas uint64) exec.GasCost {
	return func(vm *exec.VirtualMachine) (uint64, error) {
		return gas, nil
	}
}

//void emitEvent(const char *topic, size_t topicLen, const uint8_t *data, size_t dataLen);
func (r *Resolver)envEmitEvent(vm *exec.VirtualMachine) int64 {
	topic := int(int32(vm.GetCurrentFrame().Locals[0]))
	topicLen := int(int32(vm.GetCurrentFrame().Locals[1]))
	dataSrc := int(int32(vm.GetCurrentFrame().Locals[2]))
	dataLen := int(int32(vm.GetCurrentFrame().Locals[3]))

	t := make([]byte, topicLen)
	d := make([]byte, dataLen)
	copy(t, vm.Memory.Memory[topic:topic+topicLen])
	copy(d, vm.Memory.Memory[dataSrc:dataSrc+dataLen])
	address := r.context.StateDB.Address()
	topics := []common.Hash{common.BytesToHash(crypto.Keccak256(t))}
	bn := r.context.StateDB.BlockNumber().Uint64()

	r.context.StateDB.AddLog(address, topics, d, bn)
	return 0
}

func envEmitEventGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}

func (r *Resolver)envSetState(vm *exec.VirtualMachine) int64 {
	key := int(int32(vm.GetCurrentFrame().Locals[0]))
	keyLen := int(int32(vm.GetCurrentFrame().Locals[1]))
	value := int(int32(vm.GetCurrentFrame().Locals[2]))
	valueLen := int(int32(vm.GetCurrentFrame().Locals[3]))

	copyKey := make([]byte, keyLen)
	copyValue := make([]byte, valueLen)
	copy(copyKey, vm.Memory.Memory[key:key+keyLen])
	copy(copyValue, vm.Memory.Memory[value:value+valueLen])
	r.context.StateDB.SetState(copyKey, copyValue)
	return 0
}

func envSetStateGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}

func (r *Resolver)envGetState(vm *exec.VirtualMachine) int64 {
	key := int(int32(vm.GetCurrentFrame().Locals[0]))
	keyLen := int(int32(vm.GetCurrentFrame().Locals[1]))
	value := int(int32(vm.GetCurrentFrame().Locals[2]))
	valueLen := int(int32(vm.GetCurrentFrame().Locals[3]))

	val := r.context.StateDB.GetState(vm.Memory.Memory[key : key+keyLen])

	if len(val) > valueLen {
		return 0
	}

	copy(vm.Memory.Memory[value:value+valueLen], val)
	return 0
}

func envGetStateGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}

func (r *Resolver)envGetStateSize(vm *exec.VirtualMachine) int64 {
	key := int(int32(vm.GetCurrentFrame().Locals[0]))
	keyLen := int(int32(vm.GetCurrentFrame().Locals[1]))
	val := r.context.StateDB.GetState(vm.Memory.Memory[key : key+keyLen])

	return int64(len(val))
}

func envGetStateSizeGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}

// define: int64_t getNonce();
func (r *Resolver)envGetCallerNonce(vm *exec.VirtualMachine) int64 {
	return r.context.StateDB.GetCallerNonce()
}

func (r *Resolver)envCallTransfer(vm *exec.VirtualMachine) int64 {
	key := int(int32(vm.GetCurrentFrame().Locals[0]))
	keyLen := int(int32(vm.GetCurrentFrame().Locals[1]))
	value := int(vm.GetCurrentFrame().Locals[2])
	bValue := new(big.Int)
	// 256 bits
	bValue.SetBytes(vm.Memory.Memory[value : value+32])
	value256 := inner.U256(bValue)
	addr := common.BytesToAddress(vm.Memory.Memory[key : key+keyLen])

	_, returnGas, err := r.context.StateDB.Transfer(addr, value256)

	r.context.GasUsed -= returnGas
	if err != nil {
		return 1
	} else {
		return 0
	}
}

func (r *Resolver)envPlatonCall(vm *exec.VirtualMachine) int64 {
	addr := int(int32(vm.GetCurrentFrame().Locals[0]))
	params := int(int32(vm.GetCurrentFrame().Locals[1]))
	paramsLen := int(int32(vm.GetCurrentFrame().Locals[2]))
	_, err := r.context.StateDB.Call(vm.Memory.Memory[addr:addr+20], vm.Memory.Memory[params:params+paramsLen])
	if err != nil {
		fmt.Printf("call error,%s", err.Error())
		return 0
	}
	return 0
}
func (r *Resolver)envPlatonDelegateCall(vm *exec.VirtualMachine) int64 {
	addr := int(int32(vm.GetCurrentFrame().Locals[0]))
	params := int(int32(vm.GetCurrentFrame().Locals[1]))
	paramsLen := int(int32(vm.GetCurrentFrame().Locals[2]))

	_, err := r.context.StateDB.DelegateCall(vm.Memory.Memory[addr:addr+20], vm.Memory.Memory[params:params+paramsLen])
	if err != nil {
		fmt.Printf("call error,%s", err.Error())
		return 0
	}
	return 0
}

func (r *Resolver)envPlatonCallInt64(vm *exec.VirtualMachine) int64 {
	addr := int(int32(vm.GetCurrentFrame().Locals[0]))
	params := int(int32(vm.GetCurrentFrame().Locals[1]))
	paramsLen := int(int32(vm.GetCurrentFrame().Locals[2]))

	ret, err := r.context.StateDB.Call(vm.Memory.Memory[addr:addr+20], vm.Memory.Memory[params:params+paramsLen])
	if err != nil {
		fmt.Printf("call error,%s", err.Error())
		return 0
	}
	return common.BytesToInt64(ret)
}

func (r *Resolver)envPlatonDelegateCallInt64(vm *exec.VirtualMachine) int64 {
	addr := int(int32(vm.GetCurrentFrame().Locals[0]))
	params := int(int32(vm.GetCurrentFrame().Locals[1]))
	paramsLen := int(int32(vm.GetCurrentFrame().Locals[2]))

	ret, err := r.context.StateDB.DelegateCall(vm.Memory.Memory[addr:addr+20], vm.Memory.Memory[params:params+paramsLen])
	if err != nil {
		fmt.Printf("call error,%s", err.Error())
		return 0
	}
	return common.BytesToInt64(ret)
}

func (r *Resolver)envPlatonCallString(vm *exec.VirtualMachine) int64 {
	addr := int(int32(vm.GetCurrentFrame().Locals[0]))
	params := int(int32(vm.GetCurrentFrame().Locals[1]))
	paramsLen := int(int32(vm.GetCurrentFrame().Locals[2]))

	ret, err := r.context.StateDB.Call(vm.Memory.Memory[addr:addr+20], vm.Memory.Memory[params:params+paramsLen])
	if err != nil {
		fmt.Printf("call error,%s", err.Error())
		return 0
	}
	return MallocString(vm, string(ret))
}

func (r *Resolver)envPlatonDelegateCallString(vm *exec.VirtualMachine) int64 {
	addr := int(int32(vm.GetCurrentFrame().Locals[0]))
	params := int(int32(vm.GetCurrentFrame().Locals[1]))
	paramsLen := int(int32(vm.GetCurrentFrame().Locals[2]))

	ret, err := r.context.StateDB.DelegateCall(vm.Memory.Memory[addr:addr+20], vm.Memory.Memory[params:params+paramsLen])
	if err != nil {
		fmt.Printf("call error,%s", err.Error())
		return 0
	}
	return MallocString(vm, string(ret))
}

func envPlatonCallGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}

func envPlatonCallInt64GasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}

func envPlatonCallStringGasCost(vm *exec.VirtualMachine) (uint64, error) {
	return 1, nil
}
