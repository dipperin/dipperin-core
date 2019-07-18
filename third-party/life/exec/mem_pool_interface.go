package exec

type MemPoolInterface interface {
	Get(pages int) []byte
	Put(mem []byte)
}
