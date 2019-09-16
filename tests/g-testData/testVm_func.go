package g_testData

import (
	"github.com/ethereum/go-ethereum/rlp"
	"io/ioutil"
)

func GetCallExtraData(funcName, param string) ([]byte, error) {
	input := []interface{}{
		funcName,
		param,
	}

	result, err := rlp.EncodeToBytes(input)
	return result, err
}

func GetCreateExtraData(wasmPath, abiPath string, init string) ([]byte, error) {
	// GetContractExtraData
	WASMBytes, err := ioutil.ReadFile(wasmPath)
	if err != nil {
		return WASMBytes, err
	}

	abiBytes, err := ioutil.ReadFile(abiPath)
	if err != nil {
		return abiBytes, err
	}

	var rlpParams []interface{}
	if init == "" {
		rlpParams = []interface{}{
			WASMBytes, abiBytes,
		}
	} else {
		rlpParams = []interface{}{
			WASMBytes, abiBytes, init,
		}
	}

	data, err := rlp.EncodeToBytes(rlpParams)
	return data, err
}
