package mem_manage

//memory pool abstract interface
type MemPoolInterface interface {
	Get(pages int) []byte
	Put(mem []byte)
}

//memory management abstract interface
type MemManagementInterface interface {
	Malloc(size int) int
	Free(offset int) error
}
