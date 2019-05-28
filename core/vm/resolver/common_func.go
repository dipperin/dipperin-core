package resolver

import "github.com/dipperin/dipperin-core/third-party/life/exec"

func MallocString(vm *exec.VirtualMachine, str string) int64 {
	mem := vm.Memory
	size := len([]byte(str)) + 1

	pos := mem.Malloc(size)
	copy(mem.Memory[pos:pos+size], []byte(str))
	return int64(pos)
}
