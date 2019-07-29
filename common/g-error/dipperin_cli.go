package g_error

import "errors"

var (
	ErrParseBigIntFromString = errors.New("get big int from string error")
	ErrMissNumber            = errors.New("missing number")
	ErrCharacterIsNotDigit   = errors.New("the first and last character should be 0~9")
	ErrInvalidDecimalLength  = errors.New("decimal length is invalid")
)
