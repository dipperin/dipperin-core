package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAlign32BytesConverter(t *testing.T) {
	num1 := uint64(1000000)
	value := MakeUpBytes(Uint64ToBytes(num1), "uint64")
	result, err := Align32BytesConverter(value, "uint64")
	assert.NoError(t, err)
	assert.Equal(t, num1, result)

	num2 := int64(-1000000)
	value = MakeUpBytes(Int64ToBytes(num2), "int64")
	result, err = Align32BytesConverter(value, "int64")
	assert.NoError(t, err)
	assert.Equal(t, num2, result)

	num3 := uint16(10000)
	value = MakeUpBytes(Uint16ToBytes(num3), "uint16")
	result, err = Align32BytesConverter(value, "uint16")
	assert.NoError(t, err)
	assert.Equal(t, num3, result)

	num4 := int16(-10000)
	value = MakeUpBytes(Int16ToBytes(num4), "int16")
	result, err = Align32BytesConverter(value, "int16")
	assert.NoError(t, err)
	assert.Equal(t, num4, result)

	num5 := uint32(10000)
	value = MakeUpBytes(Uint32ToBytes(num5), "uint32")
	result, err = Align32BytesConverter(value, "uint32")
	assert.NoError(t, err)
	assert.Equal(t, num5, result)

	num6 := int32(-10000)
	value = MakeUpBytes(Int32ToBytes(num6), "int32")
	result, err = Align32BytesConverter(value, "int32")
	assert.NoError(t, err)
	assert.Equal(t, num6, result)

	f1 := float32(1.11)
	value = MakeUpBytes(Float32ToBytes(f1), "float32")
	result, err = Align32BytesConverter(value, "float32")
	assert.NoError(t, err)
	assert.Equal(t, f1, result)

	f2 := float64(1.23)
	value = MakeUpBytes(Float64ToBytes(f2), "float64")
	result, err = Align32BytesConverter(value, "float64")
	assert.NoError(t, err)
	assert.Equal(t, f2, result)

	str := "hello"
	value = MakeUpBytes([]byte(str), "string")
	result, err = Align32BytesConverter(value, "string")
	assert.NoError(t, err)
	assert.Equal(t, str, result)

	b := true
	value = MakeUpBytes(BoolToBytes(b), "bool")
	result, err = Align32BytesConverter(value, "bool")
	assert.NoError(t, err)
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