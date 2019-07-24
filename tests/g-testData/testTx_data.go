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

type ContractTestParameter struct {
	NodeName      string
	WASMPath      string
	AbiPath       string
	InitInputPara string
	CallFuncName  string
	CallInputPara string
}

var ContractTestPar = ContractTestParameter{
	NodeName:      "default_v0",
	WASMPath:      GetWASMPath("dipclib_test",DIPCTestContract),
	AbiPath:       GetAbiPath("dipclib_test",DIPCTestContract),
	InitInputPara: "",
	CallFuncName:  "libTest",
	CallInputPara: "",
}
