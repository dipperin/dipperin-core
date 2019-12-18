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
	"encoding/binary"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/vm/model"
	"github.com/dipperin/dipperin-core/third-party/crypto"
	"github.com/dipperin/dipperin-core/third-party/life/exec"
	"github.com/dipperin/dipperin-core/third-party/log"
	"math"
)

type uint128 struct {
	high uint64
	low  uint64
}

func (u *uint128) lsh(shift uint) {
	if shift >= 128 {
		u.low = 0
		u.high = 0
	} else {
		var halfSize uint = 128 / 2

		if shift >= halfSize {
			shift -= halfSize
			u.high = u.low
			u.low = 0
		}

		if shift != 0 {
			u.high <<= shift
		}

		var mask uint64 = ^(math.MaxUint64 >> shift)
		u.high |= (u.low & mask) >> (halfSize - shift)
		u.low <<= shift
	}
}

func (u *uint128) rsh(shift uint) {
	if shift >= 128 {
		u.high = 0
		u.low = 0
	} else {
		var halfSize uint = 128 / 2

		if shift >= halfSize {
			shift -= halfSize
			u.low = u.high
			u.high = 0
		}

		if shift != 0 {
			u.low >>= shift
		}

		var mask uint64 = ^(math.MaxUint64 << shift)
		u.low |= (u.high & mask) << (halfSize - shift)
		u.high >>= shift
	}
}

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

	log.Info("Get Params From Memory ", "address", r.Service.Address(), "copyKey", string(copyKey), "copyValue", copyValue)
	r.Service.StateDBService.SetState(r.Service.Address(), copyKey, copyValue)
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
	//log.Info("Get Params key From Memory ", "copyKey", string(copyKey))
	val := r.Service.GetState(r.Service.Address(), copyKey)
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
	copyKey := make([]byte, keyLen)
	copy(copyKey, vm.Memory.Memory[key:key+keyLen])
	log.Info("Get Params key From Memory ", "copyKey", string(copyKey))
	val := r.Service.GetState(r.Service.Address(), copyKey)
	log.Info("Get valueLen", "valueLen", len(val))
	return int64(len(val))
}

// arithmetic long double
func (r *Resolver) env__ashlti3(vm *exec.VirtualMachine) int64 {
	frame := vm.GetCurrentFrame()
	pos := int(frame.Locals[0])

	u := &uint128{
		low:  uint64(frame.Locals[1]),
		high: uint64(frame.Locals[2]),
	}
	shift := uint(frame.Locals[3])
	u.lsh(shift)

	buf := make([]byte, 16)
	binary.LittleEndian.PutUint64(buf, u.low)
	binary.LittleEndian.PutUint64(buf[8:], u.high)
	copy(vm.Memory.Memory[pos:pos+16], buf)
	return 0
}

//func(r *Resolver) env__multi3(vm *exec.VirtualMachine) int64 {
//	frame := vm.GetCurrentFrame()
//	pos := int(frame.Locals[0])
//
//	ret := C.___multi3(
//		C.uint64_t(frame.Locals[1]),
//		C.uint64_t(frame.Locals[2]),
//		C.uint64_t(frame.Locals[3]),
//		C.uint64_t(frame.Locals[4]),
//	)
//
//	buf := C.GoBytes(unsafe.Pointer(&ret), C.sizeof___int128)
//	copy(vm.Memory.Memory[pos:pos+16], buf)
//	return 0
//}
//
//func (r *Resolver) env__divti3(vm *exec.VirtualMachine) int64  {
//	frame := vm.GetCurrentFrame()
//	pos := int(frame.Locals[0])
//
//	ret := C.___divti3(
//		C.uint64_t(frame.Locals[1]),
//		C.uint64_t(frame.Locals[2]),
//		C.uint64_t(frame.Locals[3]),
//		C.uint64_t(frame.Locals[4]),
//	)
//
//	buf := C.GoBytes(unsafe.Pointer(&ret), C.sizeof___int128)
//	copy(vm.Memory.Memory[pos:pos+16], buf)
//	return 0
//}

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
	//log.Info("Memory Moved", "dest", dest, "src", src, "valueLen", len, "value", vm.Memory.Memory[dest:dest+len])

	return int64(dest)
}

func MallocString(vm *exec.VirtualMachine, str string) int64 {
	mem := vm.Memory
	size := len([]byte(str)) + 1

	//log.Info("MallocString str", "str", str)

	pos := mem.Malloc(size)
	copy(mem.Memory[pos:pos+size], []byte(str))
	vm.ExternalParams = append(vm.ExternalParams, int64(pos))
	return int64(pos)
}
