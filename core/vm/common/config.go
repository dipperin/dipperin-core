package common

import (
	"github.com/dipperin/dipperin-core/third_party/life/exec"
	"github.com/dipperin/dipperin-core/third_party/life/mem-manage"
)

var DEFAULT_VM_CONFIG = exec.VMConfig{
	EnableJIT:          false,
	DefaultMemoryPages: mem_manage.DefaultPageSize,
}
