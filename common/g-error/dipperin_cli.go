package g_error

import "errors"

var (
	ErrParseBigIntFromString = errors.New("get big int from string error")
	ErrMissNumber            = errors.New("missing number")
	ErrCharacterIsNotDigit   = errors.New("the first and last character should be 0~9")
	ErrInvalidDecimalLength  = errors.New("decimal length is invalid")
	ErrInvalidAddressLen     = errors.New("address length is invalid")
	ErrInvalidAddressPrefix  = errors.New("address prefix should be 0x or 0X")
)
