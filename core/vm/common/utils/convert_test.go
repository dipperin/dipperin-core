package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAlign32BytesConverter(t *testing.T) {
	num1 := uint64(1000000)
	result := Align32BytesConverter(Uint64ToBytes(num1), "uint64")
	assert.Equal(t, num1, result)

	num2 := int64(-1000000)
	result = Align32BytesConverter(Int64ToBytes(num2), "int64")
	assert.Equal(t, num2, result)

	num3 := uint16(10000)
	result = Align32BytesConverter(Uint16ToBytes(num3), "uint16")
	assert.Equal(t, num3, result)

	num4 := int16(-10000)
	result = Align32BytesConverter(Int16ToBytes(num4), "int16")
	assert.Equal(t, num4, result)

	num5 := uint32(10000)
	result = Align32BytesConverter(Uint32ToBytes(num5), "uint32")
	assert.Equal(t, num5, result)

	num6 := int32(-10000)
	result = Align32BytesConverter(Int32ToBytes(num6), "int32")
	assert.Equal(t, num6, result)

	f1 := float32(1.11)
	result = Align32BytesConverter(Float32ToBytes(f1), "float32")
	assert.Equal(t, f1, result)

	f2 := float64(1.23)
	result = Align32BytesConverter(Float64ToBytes(f2), "float64")
	assert.Equal(t, f2, result)

	str := "hello"
	result = Align32BytesConverter([]byte(str), "string")
	assert.Equal(t, str, result)

	b := true
	result = Align32BytesConverter(BoolToBytes(b), "bool")
	assert.Equal(t, b, result)
}

func TestStringConverter(t *testing.T) {
	num, err := StringConverter("1000", "uint64")
	assert.NoError(t, err)
	assert.Equal(t, Uint64ToBytes(1000), num)

	num, err = StringConverter("-1000", "int64")
	assert.NoError(t, err)
	assert.Equal(t, Int64ToBytes(-1000), num)

	num, err = StringConverter("1000", "uint32")
	assert.NoError(t, err)
	assert.Equal(t, Uint32ToBytes(1000), num)

	num, err = StringConverter("-1000", "int32")
	assert.NoError(t, err)
	assert.Equal(t, Int32ToBytes(-1000), num)

	num, err = StringConverter("1000", "uint16")
	assert.NoError(t, err)
	assert.Equal(t, Uint16ToBytes(1000), num)

	num, err = StringConverter("-1000", "int16")
	assert.NoError(t, err)
	assert.Equal(t, Int16ToBytes(-1000), num)

	num, err = StringConverter("1.11", "float32")
	assert.NoError(t, err)
	assert.Equal(t, Float32ToBytes(1.11), num)

	num, err = StringConverter("1.23", "float64")
	assert.NoError(t, err)
	assert.Equal(t, Float64ToBytes(1.23), num)

	num, err = StringConverter("true", "bool")
	assert.NoError(t, err)
	assert.Equal(t, BoolToBytes(true), num)
}
