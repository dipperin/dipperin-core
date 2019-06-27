package utils

import (
	"fmt"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
	"testing"
	"github.com/dipperin/dipperin-core/common/hexutil"
)

func TestInt64ToBytes(t *testing.T) {
	/*	a := []byte("hello")
		fmt.Println(a)
		b := BytesToInt64(a)
		fmt.Println(b)
		c := Int64ToBytes(b)
		fmt.Println(c)*/

	a := int64(2323)
	b := Int64ToBytes(a)
	fmt.Println(b)
	c := BytesToInt64(b)
	fmt.Println(c)
}

func TestFloat32ToBytes(t *testing.T) {
	var a float32
	a = 1.1

	byte := Float32ToBytes(a)
	log.Info("the byte is:","byte",byte)
}

func TestInputRlpData(t *testing.T){
	input := "test,123,456"
	funcName := "testFunc"

	inputRlp, err := rlp.EncodeToBytes([]interface{}{
		funcName, input,
	})

	assert.NoError(t,err)

	log.Info("call inputRlp is:","inputRlp",inputRlp)

	//initInput := "test,123,456"
	//params := getRpcParamFromString(input)
}

func TestAlign32BytesConverter(t *testing.T) {
	value, innerErr := hexutil.Decode("0x00000000000000000000000000000000000000000000000000000000000f4240")
	assert.NoError(t, innerErr)

	result, err := Align32BytesConverter(value, "uint64")
	assert.NoError(t, err)
	fmt.Println(result)
}