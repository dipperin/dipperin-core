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
		given func() ([]byte,string)
		expectResult interface{}
	}{
		{
			name:"errInput",
			given: func() ([]byte,string) {
				return nil,"test"
			},
			expectResult:[]byte(nil),
		},
		{
			name:"test_uint64",
			given: func() ([]byte,string) {
				num1 := uint64(1000000)
				return Uint64ToBytes(num1), "uint64"
			},
			expectResult:uint64(1000000),
		},
		{
			name: "test_int64",
			given: func() ([]byte,string) {
				num2 := int64(-1000000)
				return Int64ToBytes(num2), "int64"
			},
			expectResult:int64(-1000000),
		},
		{
			name:"test_uint32",
			given: func() ([]byte,string) {
				num3 := uint32(10000)
				return Uint32ToBytes(num3), "uint32"
			},
			expectResult:uint32(10000),
		},
		{
			name:"test_int32",
			given: func() ([]byte,string) {
				num4 := int32(-10000)
				return Int32ToBytes(num4), "int32"
			},
			expectResult:int32(-10000),
		},
		{
			name:"test_uint16",
			given: func() ([]byte,string) {
				num5 := uint16(10000)
				return Uint16ToBytes(num5), "uint16"
			},
			expectResult:uint16(10000),
		},
		{
			name:"test_int16",
			given: func() ([]byte,string) {
				num6 := int16(-10000)
				return Int16ToBytes(num6), "int16"
			},
			expectResult:int16(-10000),
		},
		{
			name:"test_float32",
			given: func() ([]byte,string) {
				f1 := float32(1.11)
				return Float32ToBytes(f1), "float32"
			},
			expectResult:float32(1.11),
		},
		{
			name:"test_float64",
			given: func() ([]byte,string) {
				f2 := float64(1.23)
				return Float64ToBytes(f2), "float64"
			},
			expectResult:float64(1.23),
		},
		{
			name:"test_string",
			given: func() ([]byte,string) {
				str := "hello world, I'm the best programmer who code dipperin-core"
				return []byte(str), "string"
			},
			expectResult:"hello world, I'm the best programmer who code dipperin-core",
		},
		{
			name:"test_bool",
			given: func() ([]byte,string) {
				b := true
				return BoolToBytes(b), "bool"
			},
			expectResult:true,
		},
	}

	for _,tc := range testCases{
		input,fieldType := tc.given()
		assert.Equal(t, tc.expectResult,Align32BytesConverter(input, fieldType) )
	}

}

func TestStringConverter(t *testing.T) {

	testCases := []struct{
		name string
		given func() (string, string)
		expect error
		expectResult []byte
	}{
		{
			name:"test_uint64",
			given: func() (string, string) {
				return "1000", "uint64"
			},
			expect:nil,
			expectResult:Uint64ToBytes(1000),
		},
		{
			name:"test_int64",
			given: func() (string, string) {
				return "-1000", "int64"
			},
			expect:nil,
			expectResult:Int64ToBytes(-1000),
		},
		{
			name:"test_uint32",
			given: func() (string, string) {
				return "1000", "uint32"
			},
			expect:nil,
			expectResult:Uint32ToBytes(1000),
		},
		{
			name:"test_int32",
			given: func() (string, string) {
				return "-1000", "int32"
			},
			expect:nil,
			expectResult:Int32ToBytes(-1000),
		},
		{
			name:"test_uint16",
			given: func() (string, string) {
				return "1000", "uint16"
			},
			expect:nil,
			expectResult:Uint16ToBytes(1000),
		},
		{
			name:"test_uint16",
			given: func() (string, string) {
				return "1000", "uint16"
			},
			expect:nil,
			expectResult:Uint16ToBytes(1000),
		},
		{
			name:"test_int16",
			given: func() (string, string) {
				return "-1000", "int16"
			},
			expect:nil,
			expectResult:Int16ToBytes(-1000),
		},
		{
			name:"test_float32",
			given: func() (string, string) {
				return "1.11", "float32"
			},
			expect:nil,
			expectResult:Float32ToBytes(1.11),
		},
		{
			name:"test_float64",
			given: func() (string, string) {
				return "1.23", "float64"
			},
			expect:nil,
			expectResult:Float64ToBytes(1.23),
		},
		{
			name:"test_bool",
			given: func() (string, string) {
				return "true", "bool"
			},
			expect:nil,
			expectResult:BoolToBytes(true),
		},
		{
			name:"test_bool_err",
			given: func() (string, string) {
				return "test", "bool"
			},
			expect:gerror.ErrBoolParam,
		},
		{
			name:"test_string",
			given: func() (string, string) {
				return "test", "string"
			},
			expect:nil,
			expectResult:[]byte("test"),
		},
		{
			name:"test_string",
			given: func() (string, string) {
				return "test", "test"
			},
			expect:gerror.ErrParamType,
		},
	}

	for _, tc := range testCases{
		t.Log(tc.name)
		input, fieldType := tc.given()
		result, err := StringConverter(input, fieldType)
		if err != nil {
			assert.Equal(t, tc.expect, err)
		} else {
			assert.Equal(t, tc.expectResult, result)
		}
	}













}
