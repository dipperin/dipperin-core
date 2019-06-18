package utils

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"fmt"
)

func TestConvertInputs(t *testing.T) {
	inputParam1 := InputParam{
		Type: "string",
	}
	inputParam2 := InputParam{
		Type: "string",
	}
	inputParam3 := InputParam{
		Type: "uint64",
	}
	data := []byte{242, 128, 172, 48, 48, 48, 48, 52, 49, 55, 57, 100, 53, 55, 101, 52, 53, 99, 98, 51, 98, 53, 52, 100, 54, 102, 97, 101, 102, 54, 57, 101, 55, 52, 54, 98, 102, 50, 52, 48, 101, 50, 56, 55, 57, 55, 56, 131, 15, 66, 64}
	convertData, err := ConvertInputs(data, []InputParam{inputParam1, inputParam2, inputParam3})
	assert.NoError(t, err)
	fmt.Println(string(convertData))
}