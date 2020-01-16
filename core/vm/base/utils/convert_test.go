// Copyright 2019, Keychain Foundation Ltd.
// This file is part of the Dipperin-core library.
//
// The Dipperin-core library is free software: you can redistribute
// it and/or modify it under the terms of the GNU Lesser General Public License
// as published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// The Dipperin-core library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package utils

import (
	"github.com/dipperin/dipperin-core/common/gerror"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAlign32BytesConverter(t *testing.T) {

	testCases := []struct{
		name string
		given func() error
		expect error
	}{
		{
			name:"errInput",
			given: func() error {
				result := Align32BytesConverter(nil, "test")
				assert.Equal(t, []byte(nil), result)
				return nil
			},
			expect:nil,
		},
		{
			name:"test_uint64",
			given: func() error {
				num1 := uint64(1000000)
				result := Align32BytesConverter(Uint64ToBytes(num1), "uint64")
				assert.Equal(t, num1, result)
				return nil
			},
			expect:nil,
		},
		{
			name: "test_int64",
			given: func() error {
				num2 := int64(-1000000)
				result := Align32BytesConverter(Int64ToBytes(num2), "int64")
				assert.Equal(t, num2, result)
				return nil
			},
			expect:nil,
		},
		{
			name:"test_uint32",
			given: func() error {
				num3 := uint32(10000)
				result := Align32BytesConverter(Uint32ToBytes(num3), "uint32")
				assert.Equal(t, num3, result)
				return nil
			},
			expect:nil,
		},
		{
			name:"test_int32",
			given: func() error {
				num4 := int32(-10000)
				result := Align32BytesConverter(Int32ToBytes(num4), "int32")
				assert.Equal(t, num4, result)
				return nil
			},
			expect:nil,
		},
		{
			name:"test_uint16",
			given: func() error {
				num5 := uint16(10000)
				result := Align32BytesConverter(Uint16ToBytes(num5), "uint16")
				assert.Equal(t, num5, result)
				return nil
			},
			expect:nil,
		},
		{
			name:"test_int16",
			given: func() error {
				num6 := int16(-10000)
				result := Align32BytesConverter(Int16ToBytes(num6), "int16")
				assert.Equal(t, num6, result)
				return nil
			},
			expect:nil,
		},
		{
			name:"test_float32",
			given: func() error {
				f1 := float32(1.11)
				result := Align32BytesConverter(Float32ToBytes(f1), "float32")
				assert.Equal(t, f1, result)
				return nil
			},
			expect:nil,
		},
		{
			name:"test_float64",
			given: func() error {
				f2 := float64(1.23)
				result := Align32BytesConverter(Float64ToBytes(f2), "float64")
				assert.Equal(t, f2, result)
				return nil
			},
			expect:nil,
		},
		{
			name:"test_string",
			given: func() error {
				str := "hello world, I'm the best programmer who code dipperin-core"
				result := Align32BytesConverter([]byte(str), "string")
				assert.Equal(t, str, result)
				result = Align32BytesConverter([]byte(result.(string)), "string")
				assert.Equal(t, str, result)
				return nil
			},
			expect:nil,
		},
		{
			name:"test_bool",
			given: func() error {
				b := true
				result := Align32BytesConverter(BoolToBytes(b), "bool")
				assert.Equal(t, b, result)
				return nil
			},
			expect:nil,
		},
	}

	for _,tc := range testCases{
		assert.Equal(t, tc.expect,tc.given())
	}

}

func TestStringConverter(t *testing.T) {

	testCases := []struct{
		name string
		given func() error
		expect error
	}{
		{
			name:"test_uint64",
			given: func() error {
				num, err := StringConverter("1000", "uint64")
				assert.Equal(t, Uint64ToBytes(1000), num)
				return err
			},
			expect:nil,
		},
		{
			name:"test_int64",
			given: func() error {
				num, err := StringConverter("-1000", "int64")
				assert.Equal(t, Int64ToBytes(-1000), num)
				return err
			},
			expect:nil,
		},
		{
			name:"test_uint32",
			given: func() error {
				num, err := StringConverter("1000", "uint32")
				assert.Equal(t, Uint32ToBytes(1000), num)
				return err
			},
			expect:nil,
		},
		{
			name:"test_int32",
			given: func() error {
				num, err := StringConverter("-1000", "int32")
				assert.Equal(t, Int32ToBytes(-1000), num)
				return err
			},
			expect:nil,
		},
		{
			name:"test_uint16",
			given: func() error {
				num, err := StringConverter("1000", "uint16")
				assert.Equal(t, Uint16ToBytes(1000), num)
				return err
			},
			expect:nil,
		},
		{
			name:"test_uint16",
			given: func() error {
				num, err := StringConverter("1000", "uint16")
				assert.Equal(t, Uint16ToBytes(1000), num)
				return err
			},
			expect:nil,
		},
		{
			name:"test_int16",
			given: func() error {
				num, err := StringConverter("-1000", "int16")
				assert.Equal(t, Int16ToBytes(-1000), num)
				return err
			},
			expect:nil,
		},
		{
			name:"test_float32",
			given: func() error {
				num, err := StringConverter("1.11", "float32")
				assert.Equal(t, Float32ToBytes(1.11), num)
				return err
			},
			expect:nil,
		},
		{
			name:"test_float64",
			given: func() error {
				num, err := StringConverter("1.23", "float64")
				assert.Equal(t, Float64ToBytes(1.23), num)
				return err
			},
			expect:nil,
		},
		{
			name:"test_bool",
			given: func() error {
				num, err := StringConverter("true", "bool")
				assert.Equal(t, BoolToBytes(true), num)
				return err
			},
			expect:nil,
		},
		{
			name:"test_bool_err",
			given: func() error {
				_, err := StringConverter("test", "bool")
				assert.Error(t, err)
				return err
			},
			expect:gerror.ErrBoolParam,
		},
		{
			name:"test_string",
			given: func() error {
				num, err := StringConverter("test", "string")
				assert.Equal(t, []byte("test"), num)
				return err
			},
			expect:nil,
		},
		{
			name:"test_string",
			given: func() error {
				_, err := StringConverter("test", "test")
				assert.Error(t, err)
				return err
			},
			expect:gerror.ErrParamType,
		},
	}

	for _, tc := range testCases{
		assert.Equal(t, tc.expect, tc.given())
	}













}
