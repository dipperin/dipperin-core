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

import "C"
import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/log"
	model2 "github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third_party/crypto"
	"github.com/dipperin/dipperin-core/third_party/life/exec"
	"go.uber.org/zap"
)

// Sstore
func (r *Resolver) envSetState(vm *exec.VirtualMachine) int64 {
	log.DLogger.Info("envSetState Called")
	key := int(int32(vm.GetCurrentFrame().Locals[0]))
	keyLen := int(int32(vm.GetCurrentFrame().Locals[1]))
	value := int(int32(vm.GetCurrentFrame().Locals[2]))
	valueLen := int(int32(vm.GetCurrentFrame().Locals[3]))

	log.DLogger.Info("Frame Locals", zap.Int("keyPos", key), zap.Int("keyLen", keyLen), zap.Int("valuePos", value), zap.Int("valueLen", valueLen))
	copyKey := make([]byte, keyLen)
	copyValue := make([]byte, valueLen)
	copy(copyKey, vm.Memory.Memory[key:key+keyLen])
	copy(copyValue, vm.Memory.Memory[value:value+valueLen])

	log.DLogger.Info("Get Params From Memory ", zap.Any("address", r.Service.Address()), zap.String("copyKey", string(copyKey)), zap.Uint8s("copyValue", copyValue))
	r.Service.StateDBService.SetState(r.Service.Address(), copyKey, copyValue)
	return 0
}

//Sload
func (r *Resolver) envGetState(vm *exec.VirtualMachine) int64 {
	log.DLogger.Info("envGetState Called")
	key := int(int32(vm.GetCurrentFrame().Locals[0]))
	keyLen := int(int32(vm.GetCurrentFrame().Locals[1]))
	value := int(int32(vm.GetCurrentFrame().Locals[2]))
	valueLen := int(int32(vm.GetCurrentFrame().Locals[3]))

	copyKey := make([]byte, keyLen)
	copy(copyKey, vm.Memory.Memory[key:key+keyLen])
	//log.DLogger.Info("Get Params key From Memory ", "copyKey", string(copyKey))
	val := r.Service.GetState(r.Service.Address(), copyKey)
	if len(val) > valueLen {
		return 0
	}
	copy(vm.Memory.Memory[value:value+valueLen], val)
	log.DLogger.Info("Save Value Into Memory", zap.Int("valuePos", value), zap.Int("valueLen", valueLen), zap.Uint8s("value", val))
	return 0
}

func (r *Resolver) envGetStateSize(vm *exec.VirtualMachine) int64 {
	log.DLogger.Info("envGetStateSize Called")
	key := int(int32(vm.GetCurrentFrame().Locals[0]))
	keyLen := int(int32(vm.GetCurrentFrame().Locals[1]))
	copyKey := make([]byte, keyLen)
	copy(copyKey, vm.Memory.Memory[key:key+keyLen])
	log.DLogger.Info("Get Params key From Memory ", zap.String("copyKey", string(copyKey)))
	val := r.Service.GetState(r.Service.Address(), copyKey)
	log.DLogger.Info("Get valueLen", zap.Int("valueLen", len(val)))
	return int64(len(val))
}

//void emitEvent(const char *topic, size_t topicLen, const uint8_t *data, size_t dataLen);
//topic = funcName
//data = param...
func (r *Resolver) envEmitEvent(vm *exec.VirtualMachine) int64 {
	log.DLogger.Info("emitEvent Called")

	topic := int(int32(vm.GetCurrentFrame().Locals[0]))
	topicLen := int(int32(vm.GetCurrentFrame().Locals[1]))
	dataSrc := int(int32(vm.GetCurrentFrame().Locals[2]))
	dataLen := int(int32(vm.GetCurrentFrame().Locals[3]))

	t := make([]byte, topicLen)
	d := make([]byte, dataLen)
	copy(t, vm.Memory.Memory[topic:topic+topicLen])
	copy(d, vm.Memory.Memory[dataSrc:dataSrc+dataLen])

	log.DLogger.Info("the blockNumber is:", zap.Any("blockNumber", r.Service.GetBlockNumber()))
	log.DLogger.Info("envEmitEvent", zap.String("TopicName", string(t)), zap.Int("Len", topicLen))
	addedLog := &model2.Log{
		Address:   r.Service.Address(),
		Topics:    []common.Hash{common.BytesToHash(crypto.Keccak256(t))},
		TopicName: string(t),
		Data:      d,
		TxHash:    r.Service.GetTxHash(),
	}
	r.Service.AddLog(addedLog)
	return 0
}

func envMalloc(vm *exec.VirtualMachine) int64 {
	//log.DLogger.Info("envMalloc Called")
	size := int(uint32(vm.GetCurrentFrame().Locals[0]))

	pos := vm.Memory.Malloc(size)
	if pos == -1 {
		panic("melloc error...")
	}

	//log.DLogger.Info("Malloc Memory", "pos", pos, "size", size)
	return int64(pos)
}

func envFree(vm *exec.VirtualMachine) int64 {
	/*	if vmcommon.Config.DisableFree {
		return 0
	}*/

	//log.DLogger.Info("envFree Called")
	mem := vm.Memory
	offset := int(uint32(vm.GetCurrentFrame().Locals[0]))

	err := mem.Free(offset)
	if err != nil {
		panic("free error...")
	}
	//log.DLogger.Info("Malloc Free", "offset", offset)
	return 0
}

//void * memory copy ( void * destination, const void * source, size_t num );
func envMemcpy(vm *exec.VirtualMachine) int64 {
	//log.DLogger.Info("envMemcpy Called")

	dest := int(uint32(vm.GetCurrentFrame().Locals[0]))
	src := int(uint32(vm.GetCurrentFrame().Locals[1]))
	len := int(uint32(vm.GetCurrentFrame().Locals[2]))

	copy(vm.Memory.Memory[dest:dest+len], vm.Memory.Memory[src:src+len])

	//log.DLogger.Info("Memory Copyed", "dest", dest, "src", src, "valueLen", len, "value", vm.Memory.Memory[dest:dest+len])
	return int64(dest)
}

//void * memmove ( void * destination, const void * source, size_t num );
func envMemmove(vm *exec.VirtualMachine) int64 {
	//log.DLogger.Info("envMemmove Called")
	dest := int(uint32(vm.GetCurrentFrame().Locals[0]))
	src := int(uint32(vm.GetCurrentFrame().Locals[1]))
	len := int(uint32(vm.GetCurrentFrame().Locals[2]))

	copy(vm.Memory.Memory[dest:dest+len], vm.Memory.Memory[src:src+len])
	//log.DLogger.Info("Memory Moved", "dest", dest, "src", src, "valueLen", len, "value", vm.Memory.Memory[dest:dest+len])

	return int64(dest)
}

func MallocString(vm *exec.VirtualMachine, str string) int64 {
	mem := vm.Memory
	size := len([]byte(str)) + 1

	//log.DLogger.Info("MallocString str", "str", str)

	pos := mem.Malloc(size)
	copy(mem.Memory[pos:pos+size], []byte(str))
	vm.ExternalParams = append(vm.ExternalParams, int64(pos))
	return int64(pos)
}
