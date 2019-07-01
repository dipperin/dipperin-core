package g_testData

import (
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func GetCallExtraData(t *testing.T, funcName, param string) []byte {
	input := []interface{}{
		funcName,
		param,
	}

	result, err := rlp.EncodeToBytes(input)
	assert.NoError(t, err)
	return result
}

func GetCreateExtraData(t *testing.T, wasmPath, abiPath string, init string) []byte {
	// GetContractExtraData
	wasmBytes, err := ioutil.ReadFile(wasmPath)
	assert.NoError(t, err)

	abiBytes, err := ioutil.ReadFile(abiPath)
	assert.NoError(t, err)

	var rlpParams []interface{}
	if init == "" {
		rlpParams = []interface{}{
			wasmBytes, abiBytes,
		}
	} else {
		rlpParams = []interface{}{
			wasmBytes, abiBytes, init,
		}
	}

	data, err := rlp.EncodeToBytes(rlpParams)
	assert.NoError(t, err)
	return data
}