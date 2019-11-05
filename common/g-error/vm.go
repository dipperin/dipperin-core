package g_error

import "errors"

var (
	ErrOutOfGas                  = errors.New("out of gas")
	ErrCodeStoreOutOfGas         = errors.New("contract creation code storage out of gas")
	ErrDepth                     = errors.New("max call depth exceeded")
	ErrTraceLimitReached         = errors.New("the number of logs reached the specified limit")
	ErrInsufficientBalance       = errors.New("insufficient balance for transfer")
	ErrContractAddressCollision  = errors.New("contract address collision")
	ErrInsufficientBalanceForGas = errors.New("insufficient balance to pay for gas")
	ErrExecutionReverted         = errors.New("vm: execution reverted")
	ErrMaxCodeSizeExceeded       = errors.New("vm: max code size exceeded")
)
