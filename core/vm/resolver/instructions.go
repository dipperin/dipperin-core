package resolver

import (
	"github.com/dipperin/dipperin-core/core/vm/model"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/third-party/crypto"
	"github.com/dipperin/dipperin-core/third-party/life/exec"
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

	log.Info("Get Params From Memory ", "copyKey", string(copyKey), "copyValue", copyValue)
	r.Service.ReSolverSetState(copyKey, copyValue)
	return 0
}

//Sload
func (r *Resolver) envGetState(vm *exec.VirtualMachine) int64 {
	log.Info("envGetState Called")
	key := int(int32(vm.GetCurrentFrame().Locals[0]))
	keyLen := int(int32(vm.GetCurrentFrame().Locals[1]))
	value := int(int32(vm.GetCurrentFrame().Locals[2]))
	valueLen := int(int32(vm.GetCurrentFrame().Locals[3]))

	copyKey := make([]byte, keyLen)
	copy(copyKey, vm.Memory.Memory[key:key+keyLen])
	log.Info("Get Params key From Memory ", "copyKey", string(copyKey))
	val := r.Service.ReSolverGetState(copyKey)
	if len(val) > valueLen {
		return 0
	}
	copy(vm.Memory.Memory[value:value+valueLen], val)
	log.Info("Save Value Into Memory", "valuePos", value, "valueLen", valueLen, "value", string(val))
	return 0
}

func (r *Resolver) envGetStateSize(vm *exec.VirtualMachine) int64 {
	log.Info("envGetStateSize Called")
	key := int(int32(vm.GetCurrentFrame().Locals[0]))
	keyLen := int(int32(vm.GetCurrentFrame().Locals[1]))
	val := r.Service.ReSolverGetState(vm.Memory.Memory[key: key+keyLen])
	log.Info("Get valueLen", "valueLen", len(val))
	return int64(len(val))
}

//void emitEvent(const char *topic, size_t topicLen, const uint8_t *data, size_t dataLen);
//topic = funcName
//data = param...
func (r *Resolver) envEmitEvent(vm *exec.VirtualMachine) int64 {
	log.Info("emitEvent Called")

	topic := int(int32(vm.GetCurrentFrame().Locals[0]))
	topicLen := int(int32(vm.GetCurrentFrame().Locals[1]))
	dataSrc := int(int32(vm.GetCurrentFrame().Locals[2]))
	dataLen := int(int32(vm.GetCurrentFrame().Locals[3]))

	t := make([]byte, topicLen)
	d := make([]byte, dataLen)
	copy(t, vm.Memory.Memory[topic:topic+topicLen])
	copy(d, vm.Memory.Memory[dataSrc:dataSrc+dataLen])

	log.Info("the blockNumber is:", "blockNumber", r.Service.GetBlockNumber())
	log.Info("envEmitEvent", "TopicName", string(t), "Len", topicLen)
	addedLog := &model.Log{
		Address:     r.Service.Address(),
		Topics:      []common.Hash{common.BytesToHash(crypto.Keccak256(t))},
		TopicName:   string(t),
		Data:        d,
		BlockNumber: r.Service.GetBlockNumber().Uint64(),
		TxHash:      r.Service.GetTxHash(),
		TxIndex:     uint(r.Service.GetTxIdx()),
		BlockHash:   r.Service.GetCurBlockHash(),
	}
	r.Service.AddLog(addedLog)
	return 0
}

func envMalloc(vm *exec.VirtualMachine) int64 {
	//log.Info("envMalloc Called")
	size := int(uint32(vm.GetCurrentFrame().Locals[0]))

	pos := vm.Memory.Malloc(size)
	if pos == -1 {
		panic("melloc error...")
	}

	//log.Info("Malloc Memory", "pos", pos, "size", size)
	return int64(pos)
}

func envFree(vm *exec.VirtualMachine) int64 {
	/*	if vmcommon.Config.DisableFree {
			return 0
		}*/
	//log.Info("envFree Called")
	mem := vm.Memory
	offset := int(uint32(vm.GetCurrentFrame().Locals[0]))

	err := mem.Free(offset)
	if err != nil {
		panic("free error...")
	}
	//log.Info("Malloc Free", "offset", offset)
	return 0
}

//void * memory copy ( void * destination, const void * source, size_t num );
func envMemcpy(vm *exec.VirtualMachine) int64 {
	//log.Info("envMemcpy Called")
	dest := int(uint32(vm.GetCurrentFrame().Locals[0]))
	src := int(uint32(vm.GetCurrentFrame().Locals[1]))
	len := int(uint32(vm.GetCurrentFrame().Locals[2]))

	copy(vm.Memory.Memory[dest:dest+len], vm.Memory.Memory[src:src+len])
	//log.Info("Memory Copyed", "dest", dest, "src", src, "valueLen", len, "value", vm.Memory.Memory[dest:dest+len])
	return int64(dest)
}

//void * memmove ( void * destination, const void * source, size_t num );
func envMemmove(vm *exec.VirtualMachine) int64 {
	//log.Info("envMemmove Called")
	dest := int(uint32(vm.GetCurrentFrame().Locals[0]))
	src := int(uint32(vm.GetCurrentFrame().Locals[1]))
	len := int(uint32(vm.GetCurrentFrame().Locals[2]))

	copy(vm.Memory.Memory[dest:dest+len], vm.Memory.Memory[src:src+len])
	log.Info("Memory Moved", "dest", dest, "src", src, "valueLen", len, "value", vm.Memory.Memory[dest:dest+len])
	return int64(dest)
}
