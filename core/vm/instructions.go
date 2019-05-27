package vm

import (
	"fmt"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/third-party/crypto"
	"github.com/dipperin/dipperin-core/third-party/life/exec"
	"github.com/dipperin/dipperin-core/core/vm/model"
)

// Sstore
func (r *Resolver) envSetState(vm *exec.VirtualMachine) int64 {
	log.Info("envSetState Called")
	key := int(int32(vm.GetCurrentFrame().Locals[0]))
	keyLen := int(int32(vm.GetCurrentFrame().Locals[1]))
	value := int(int32(vm.GetCurrentFrame().Locals[2]))
	valueLen := int(int32(vm.GetCurrentFrame().Locals[3]))

	log.Info("Frame Locals", "keyPos", key, "keyLen", keyLen, "valuePos", value, "valueLen", valueLen)
	copyKey := make([]byte, keyLen)
	copyValue := make([]byte, valueLen)
	copy(copyKey, vm.Memory.Memory[key:key+keyLen])
	copy(copyValue, vm.Memory.Memory[value:value+valueLen])

	log.Info("Get Params From Memory ", "copyKey", copyKey, "copyValue", copyValue)
	r.state.SetState(r.contract.self.Address(), copyKey, copyValue)
	return 0
}

//Sload
func (r *Resolver) envGetState(vm *exec.VirtualMachine) int64 {
	log.Info("envGetState Called")
	key := int(int32(vm.GetCurrentFrame().Locals[0]))
	keyLen := int(int32(vm.GetCurrentFrame().Locals[1]))
	value := int(int32(vm.GetCurrentFrame().Locals[2]))
	valueLen := int(int32(vm.GetCurrentFrame().Locals[3]))

	val := r.state.GetState(r.contract.self.Address(), vm.Memory.Memory[key: key+keyLen])
	if len(val) > valueLen {
		return 0
	}
	copy(vm.Memory.Memory[value:value+valueLen], val)
	log.Info("Save Value Into Memory", "valuePos", value, "valueLen", valueLen, "value", val)
	return 0
}

func (r *Resolver) envGetStateSize(vm *exec.VirtualMachine) int64 {
	log.Info("envGetStateSize Called")
	key := int(int32(vm.GetCurrentFrame().Locals[0]))
	keyLen := int(int32(vm.GetCurrentFrame().Locals[1]))
	val := r.state.GetState(r.contract.self.Address(), vm.Memory.Memory[key: key+keyLen])
	log.Info("Get valueLen", "valueLen", len(val))
	return int64(len(val))
}

//void emitEvent(const char *topic, size_t topicLen, const uint8_t *data, size_t dataLen);
//topic = funcName
//data = param...
func (r *Resolver) emitEvent(vm *exec.VirtualMachine) int64 {
	log.Info("emitEvent Called")

	topic := int(int32(vm.GetCurrentFrame().Locals[0]))
	topicLen := int(int32(vm.GetCurrentFrame().Locals[1]))
	dataSrc := int(int32(vm.GetCurrentFrame().Locals[2]))
	dataLen := int(int32(vm.GetCurrentFrame().Locals[3]))

	t := make([]byte, topicLen)
	d := make([]byte, dataLen)
	copy(t, vm.Memory.Memory[topic:topic+topicLen])
	copy(d, vm.Memory.Memory[dataSrc:dataSrc+dataLen])

	addedLog := &model.Log{
		ContractAddr:r.contract.self.Address(),
		Topic:common.BytesToHash(crypto.Keccak256(t)),
		Data:d,
		BlockNumber:r.context.BlockNumber.Uint64(),
		TxHash:r.state.TxHash(),
		TxIndex:uint(r.state.TxIdx()),
		BlockHash:r.context.BlockHash,
	}

	r.state.AddLog(addedLog)
	return 0
}

func (r *Resolver) envMalloc(vm *exec.VirtualMachine) int64 {
	log.Info("envMalloc Called")
	size := int(uint32(vm.GetCurrentFrame().Locals[0]))

	pos := vm.Memory.Malloc(size)
	if pos == -1 {
		panic("melloc error...")
	}

	log.Info("Malloc Memory", "pos", pos, "size", size)
	return int64(pos)
}

func (r *Resolver) envFree(vm *exec.VirtualMachine) int64 {
	/*	if vm.Config.DisableFree {
			return 0
		}*/
	log.Info("envFree Called")
	mem := vm.Memory
	offset := int(uint32(vm.GetCurrentFrame().Locals[0]))

	err := mem.Free(offset)
	if err != nil {
		panic("free error...")
	}
	log.Info("Malloc Free", "offset", offset)
	return 0
}

//void * memory copy ( void * destination, const void * source, size_t num );
func (r *Resolver) envMemcpy(vm *exec.VirtualMachine) int64 {
	log.Info("envMemcpy Called")
	dest := int(uint32(vm.GetCurrentFrame().Locals[0]))
	src := int(uint32(vm.GetCurrentFrame().Locals[1]))
	len := int(uint32(vm.GetCurrentFrame().Locals[2]))

	copy(vm.Memory.Memory[dest:dest+len], vm.Memory.Memory[src:src+len])
	log.Info("Memory Copyed", "dest", dest, "src", src, "valueLen", len, "value", vm.Memory.Memory[dest:dest+len])
	return int64(dest)
}

//void * memmove ( void * destination, const void * source, size_t num );
func (r *Resolver) envMemmove(vm *exec.VirtualMachine) int64 {
	log.Info("envMemmove Called")
	dest := int(uint32(vm.GetCurrentFrame().Locals[0]))
	src := int(uint32(vm.GetCurrentFrame().Locals[1]))
	len := int(uint32(vm.GetCurrentFrame().Locals[2]))

	copy(vm.Memory.Memory[dest:dest+len], vm.Memory.Memory[src:src+len])
	log.Info("Memory Moved", "dest", dest, "src", src, "valueLen", len, "value", vm.Memory.Memory[dest:dest+len])
	return int64(dest)
}

func (r *Resolver) envPrints(vm *exec.VirtualMachine) int64 {
	fmt.Println("Printing...")
	start := int(uint32(vm.GetCurrentFrame().Locals[0]))
	end := 0
	for end = start; end < len(vm.Memory.Memory); end++ {
		if vm.Memory.Memory[end] == 0 {
			break
		}
	}
	//fmt.Println(start)
	fmt.Println(string(vm.Memory.Memory[start:end]), start)
	return 0
}

func (r *Resolver) envPrinti(vm *exec.VirtualMachine) int64 {
	fmt.Println("Printing int64...")
	fmt.Println(fmt.Sprintf("%d", vm.GetCurrentFrame().Locals[0]))
	return 0
}

func (r *Resolver) envPrintui(vm *exec.VirtualMachine) int64 {
	fmt.Println("Printing uint64...")
	fmt.Println(fmt.Sprintf("%d", vm.GetCurrentFrame().Locals[0]))
	return 0
}

func (r *Resolver) envPrints_l(vm *exec.VirtualMachine) int64 {
	fmt.Println("Prints_l")
	ptr := int(uint32(vm.GetCurrentFrame().Locals[0]))
	msgLen := int(uint32(vm.GetCurrentFrame().Locals[1]))

	msg := vm.Memory.Memory[ptr: ptr+msgLen]
	fmt.Println(string(msg))
	return 0
}

func (r *Resolver) envCoinBase(vm *exec.VirtualMachine) int64 {
	offset := int(int32(vm.GetCurrentFrame().Locals[0]))
	coinBase := r.context.Coinbase
	copy(vm.Memory.Memory[offset:], coinBase.Bytes())
	return 0
}

func (r *Resolver) envBlockNum(vm *exec.VirtualMachine) int64 {
	offset := int(int32(vm.GetCurrentFrame().Locals[0]))
	num := r.context.BlockNumber
	copy(vm.Memory.Memory[offset:], num.Bytes())
	return 0
}

func (r *Resolver) envDifficulty(vm *exec.VirtualMachine) int64 {
	offset := int(int32(vm.GetCurrentFrame().Locals[0]))
	diff := r.context.Difficulty
	copy(vm.Memory.Memory[offset:], diff.Bytes())
	return 0
}

func (r *Resolver) envAbort(vm *exec.VirtualMachine) int64 {
	panic("abort")
}