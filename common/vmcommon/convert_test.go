package vmcommon

import (
	"fmt"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
	"testing"
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



func TestIntConvertBytes(t *testing.T)  {
	a := "123"
	byte, err := StringConverter(a, "int64")
	log.Info("the byte is:","byte",byte)
	assert.NoError(t, err)
	v := BytesConverter(byte, "int64")
	fmt.Println(v.(int64))

	byte, err = StringConverter(a, "uint64")
	assert.NoError(t, err)
	v = BytesConverter(byte, "uint64")
	fmt.Println(v.(int64))

	a = "12345"
	byte, err = StringConverter(a, "uint16")
	assert.NoError(t, err)

	fmt.Println(BytesToUint16(byte))

	byte, err = StringConverter(a, "int16")
	assert.NoError(t, err)
	fmt.Println(BytesToInt16(byte))

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