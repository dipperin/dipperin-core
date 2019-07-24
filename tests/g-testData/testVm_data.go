package g_testData

import (
	"fmt"
	"github.com/dipperin/dipperin-core/common/util"
	"path/filepath"
)

var testCorePath = "go/src/github.com/dipperin/dipperin-core/core/vm/test-data"
var testDIPCPath = "c++/src/dipc/testcontract/"

type ContractPathType uint8
const (
	CoreVmTestData ContractPathType = iota
	DIPCTestContract
)

var contractPath =map[ContractPathType]string{
	CoreVmTestData:   testCorePath,
	DIPCTestContract: testDIPCPath,
}

func GetWASMPath(fileName string,pathType ContractPathType ) string {
	homeDir := util.HomeDir()
	path := filepath.Join(homeDir, contractPath[pathType])
	return filepath.Join(path, fmt.Sprintf("%s/%s.wasm", fileName, fileName))
}

func GetAbiPath(fileName string,pathType ContractPathType) string {
	homeDir := util.HomeDir()
	path := filepath.Join(homeDir,contractPath[pathType])
	return filepath.Join(path, fmt.Sprintf("%s/%s.cpp.abi.json", fileName, fileName))
}


