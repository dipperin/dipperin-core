package g_testData

import (
	"fmt"
	"github.com/dipperin/dipperin-core/common/util"
	"io/ioutil"
	"os"
	"path/filepath"
)

var testCorePath = "src/github.com/dipperin/dipperin-core/core/vm/test-data"
var testDIPCPath = "c++/src/dipc/testcontract/"

type ContractPathType uint8

const (
	CoreVmTestData ContractPathType = iota
	DIPCTestContract
)

var contractPath = map[ContractPathType]string{
	CoreVmTestData:   testCorePath,
	DIPCTestContract: testDIPCPath,
}

func GetWASMPath(fileName string, pathType ContractPathType) string {
	var tmpPath string
	if pathType == CoreVmTestData {
		tmpPath = os.Getenv("GOPATH")
	} else {
		tmpPath = util.HomeDir()
	}
	path := filepath.Join(tmpPath, contractPath[pathType])
	return filepath.Join(path, fmt.Sprintf("%s/%s.wasm", fileName, fileName))
}

func GetAbiPath(fileName string, pathType ContractPathType) string {
	var tmpPath string
	if pathType == CoreVmTestData {
		tmpPath = os.Getenv("GOPATH")
	} else {
		tmpPath = util.HomeDir()
	}
	path := filepath.Join(tmpPath, contractPath[pathType])
	return filepath.Join(path, fmt.Sprintf("%s/%s.cpp.abi.json", fileName, fileName))
}

func GetCodeAbi(code, abi string) ([]byte, []byte) {
	fileCode, err := ioutil.ReadFile(code)
	if err != nil {
		panic(fmt.Sprintf("Read code failed, err=%s", err.Error()))
	}
	fileABI, err := ioutil.ReadFile(abi)
	if err != nil {
		panic(fmt.Sprintf("Read abi failed, err=%s", err.Error()))
	}
	return fileCode, fileABI
}
