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
	ErrInvalidOutputLength       = errors.New("vm_utils: invalid init function outputs length")
	ErrLengthInputAbiNotMatch    = errors.New("vm_utils: length of input and abi not match")
	ErrEmptyInput         = errors.New("interpreter_life: empty input")
	ErrEmptyABI           = errors.New("interpreter_life: empty abi")
	ErrInvalidRlpFormat   = errors.New("interpreter_life: invalid rlp format")
	ErrInsufficientParams = errors.New("interpreter_life: invalid input params")
	ErrInvalidAbi         = errors.New("interpreter_life: invalid abi, from json fail")
	ErrInputAbiNotMatch   = errors.New("interpreter_life: length of input and abi not match")
	ErrInvalidReturnType  = errors.New("interpreter_life: return type not void")
	ErrFuncNameNotFound   = errors.New("interpreter_life: function name not found")
	ErrCallContractAddrIsWrong   = errors.New("call contract addr is wrong")

)

