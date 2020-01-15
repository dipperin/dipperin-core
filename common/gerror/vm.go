package gerror

import "errors"

var (
	ErrOutOfGas                  = errors.New("out of gas")
	ErrCodeStoreOutOfGas         = errors.New("contract creation code storage out of gas")
	ErrDepth                     = errors.New("max call depth exceeded")
	ErrTraceLimitReached         = errors.New("the number of logs reached the specified limit")
	ErrInsufficientBalance       = errors.New("insufficient balance for transfer")
	ErrContractAddressCollision  = errors.New("contract address collision")
	ErrContractAddressCreate     = errors.New("contract address create fail")
	ErrVMTransfer                = errors.New("transfer to contract address err")
	ErrInsufficientBalanceForGas = errors.New("insufficient balance to pay for gas")
	ErrExecutionReverted         = errors.New("vm: execution reverted")
	ErrMaxCodeSizeExceeded       = errors.New("vm: max code size exceeded")
	ErrBoolParam                 = errors.New("invalid boolean param")
	ErrParamType                 = errors.New("invalid param type")
	ErrEmptyInput                = errors.New("vm_utils: empty input")
	ErrInvalidRlpFormat          = errors.New("vm_utils: invalid rlp format")
	ErrInsufficientParams        = errors.New("vm_utils: invalid input params")
	ErrInvalidOutputLength       = errors.New("vm_utils: invalid init function outputs length")
	ErrLengthInputAbiNotMatch    = errors.New("vm_utils: length of input and abi not match")
	ErrFuncNameNotFound          = errors.New("vm_utils: function name not found")
)
