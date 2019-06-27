package utils

import (
	"fmt"
	"github.com/dipperin/dipperin-core/common/hexutil"
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

func TestIntConvertBytes(t *testing.T) {
	a := "123"
	bytes, err := StringConverter(a, "int64")
	log.Info("the bytes is:","bytes", bytes)
	assert.NoError(t, err)
	v := BytesConverter(bytes, "int64")
	fmt.Println(v.(int64))


	uint64Num := uint64(128)
	fmt.Println("uint64Num", Uint64ToBytes(uint64Num))
	fmt.Println("int64Num", Int64ToBytes(int64(128)))
	fmt.Println("bytesToUint64", BytesToUint64([]byte{128}))
	fmt.Println("bytesToInt64", BytesToInt64([]byte{128}))
	fmt.Println("bytesToUint64", BytesToUint64([]byte{0,0,0,0,0,0,0,128}))
	fmt.Println("bytesToInt64", BytesToInt64([]byte{0,0,0,0,0,0,0,128}))
	vi := int64(uint64Num)
	b := make([]byte, 8)
	//fmt.Println(v >> 56)

	_ = b[7] // early bounds check to guarantee safety of writes below
	b[0] = byte(vi >> 56)
	b[1] = byte(vi >> 48)
	b[2] = byte(vi >> 40)
	b[3] = byte(vi >> 32)
	b[4] = byte(vi >> 24)
	b[5] = byte(vi >> 16)
	b[6] = byte(vi >> 8)
	b[7] = byte(vi)
	//binary.BigEndian.PutUint64(buf, int64Num)
	fmt.Println("int64Num", b)
	fmt.Println(vi)
	fmt.Println(Int64ToBytes(vi))
	fmt.Println(BytesToInt64(Int64ToBytes(vi)))

	vin := int64(128)
	b[0] = byte(vin >> 56)
	b[1] = byte(vin >> 48)
	b[2] = byte(vin >> 40)
	b[3] = byte(vin >> 32)
	b[4] = byte(vin >> 24)
	b[5] = byte(vin >> 16)
	b[6] = byte(vin >> 8)
	b[7] = byte(vin)
	fmt.Println("vin", b)


	v = BytesConverter([]byte{18}, "int64")
	fmt.Println("byte8", v.(int64))

	bytes, err = StringConverter(a, "uint64")
	assert.NoError(t, err)
	v = BytesConverter(bytes, "uint64")
	fmt.Println(v.(uint64))

	a = "12345"
	bytes, err = StringConverter(a, "uint16")
	assert.NoError(t, err)

	fmt.Println(BytesToUint16(bytes))

	bytes, err = StringConverter(a, "int16")
	assert.NoError(t, err)
	fmt.Println(BytesToInt16(bytes))

}

func TestInt64ConvertBytes(t *testing.T)  {
	fmt.Println("bytesToUint64", BytesToUint64([]byte{128}))
	fmt.Println("bytesToInt64", BytesToInt64([]byte{128}))
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