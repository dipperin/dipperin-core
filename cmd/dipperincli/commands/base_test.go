package commands

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFilterParams_UnmarshalJSON(t *testing.T) {
	params := ""
	var filterParams FilterParams
	err := json.Unmarshal([]byte(params), &filterParams)
	assert.Error(t, err)

	params = `{
		"block_hash":"0x000023e18421a0abfceea172867b9b4a3bcf593edd0b504554bb7d1cf5f5e7b7",
		"addresses":["0x0014049F835be46352eD0Ec6B819272A2c8cF4feA10f"],
		"topics":[["0x0b5d2220daf8f0dfd95983d2ce625affbb7183c991271f49d818b4a64a268dbb"]]
	}`

	err = json.Unmarshal([]byte(params), &filterParams)
	assert.NoError(t, err)
}
