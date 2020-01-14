package mem_manage

const (
	// DefaultCallStackSize is the default call stack size.
	DefaultCallStackSize = 512

	// DefaultPageSize is the linear memory page size.  65536
	DefaultPageSize = 65536

	// JITCodeSizeThreshold is the lower-bound code size threshold for the JIT compiler.
	JITCodeSizeThreshold = 30

	DefaultMemoryPages = 16
	DynamicMemoryPages = 16
	//DynamicMemoryPages = 160

	DefaultMemPoolCount = 5
	DefaultMemBlockSize = 5
)
