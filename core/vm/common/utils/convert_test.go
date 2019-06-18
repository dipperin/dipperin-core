package utils

import (
	"fmt"
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
	byte, err := StringConverter(a, "int64")
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
