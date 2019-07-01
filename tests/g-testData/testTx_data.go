package g_testData

import (
	"github.com/dipperin/dipperin-core/core/vm/model"
	"math/big"
)

var (
	TestGasPrice = big.NewInt(1)
	TestGasLimit = 2 * model.TxGas
	TestValue    = big.NewInt(100)
)
