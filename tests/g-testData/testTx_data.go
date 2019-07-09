package g_testData

import (
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/core/vm/model"
	"math/big"
	"path/filepath"
)

var (
	TestGasPrice = big.NewInt(1)
	TestGasLimit = 2 * model.TxGas
	TestValue    = big.NewInt(100)
)

type ContractTestParameter struct {
	NodeName string
	WASMPath string
	AbiPath string
	InitInputPara string
	CallFuncName string
	CallInputPara string
}

var ContractTestPar = ContractTestParameter{
	NodeName:      "default_v0",
	WASMPath:      filepath.Join(util.HomeDir(), "c++/src/dipc/testcontract/dipclib_test/dipclib_test.wasm"),
	AbiPath:       filepath.Join(util.HomeDir(), "c++/src/dipc/testcontract/dipclib_test/dipclib_test.cpp.abi.json"),
	InitInputPara: "",
	CallFuncName:  "libTest",
	CallInputPara: "",
}
