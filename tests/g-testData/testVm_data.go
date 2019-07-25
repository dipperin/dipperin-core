package g_testData

import (
	"fmt"
	"github.com/dipperin/dipperin-core/common/util"
	"path/filepath"
	"io/ioutil"
)

func GetWasmPath(fileName string) string {
	homeDir := util.HomeDir()
	path := filepath.Join(homeDir, "go/src/github.com/dipperin/dipperin-core/core/vm/test-data")
	return filepath.Join(path, fmt.Sprintf("%s/%s.wasm", fileName, fileName))
}

func GetAbiPath(fileName string) string {
	homeDir := util.HomeDir()
	path := filepath.Join(homeDir, "go/src/github.com/dipperin/dipperin-core/core/vm/test-data")
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