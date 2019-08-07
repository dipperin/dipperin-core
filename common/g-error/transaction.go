package g_error

import "errors"

var (
	ErrEmptyReceipt       = errors.New("empty receipt")
	ErrTxRootNotMatch     = errors.New("transaction root not match")
	ErrTxInSpecialBlock   = errors.New("special block have transactions")
	ErrDelegatesNotEnough = errors.New("the register tx delegate is lower than MiniPledgeValue")
)
