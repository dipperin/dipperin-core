package utils

type CodeAbi struct {
	Code  []byte `json:"code"`
	Abi   []byte `json:"abi"`
	Input []byte `json:"Input"`
}
