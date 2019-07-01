package g_testData

import (
	"github.com/dipperin/dipperin-core/common/util"
	"path/filepath"
	"fmt"
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
