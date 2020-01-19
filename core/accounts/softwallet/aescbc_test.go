// Copyright 2019, Keychain Foundation Ltd.
// This file is part of the dipperin-core library.
//
// The dipperin-core library is free software: you can redistribute
// it and/or modify it under the terms of the GNU Lesser General Public License
// as published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// The dipperin-core library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package softwallet

import (
	"crypto/aes"
	"github.com/dipperin/dipperin-core/common/gerror"
	"github.com/stretchr/testify/assert"
	"testing"
)

var commonInput = []byte{
	0x6b, 0xc1, 0xbe, 0xe2, 0x2e, 0x40, 0x9f, 0x96, 0xe9, 0x3d, 0x7e, 0x11, 0x73, 0x93, 0x17, 0x2a,
	0xae, 0x2d, 0x8a, 0x57, 0x1e, 0x03, 0xac, 0x9c, 0x9e, 0xb7, 0x6f, 0xac, 0x45, 0xaf, 0x8e, 0x51,
	0x30, 0xc8, 0x1c, 0x46, 0xa3, 0x5c, 0xe4, 0x11, 0xe5, 0xfb, 0xc1, 0x19, 0x1a, 0x0a, 0x52, 0xef,
	0xf6, 0x9f, 0x24, 0x45, 0xdf, 0x4f, 0x9b, 0x17, 0xad, 0x2b, 0x41, 0x7b, 0xe6, 0x6c, 0x37, 0x10,
}

var commonKey256 = []byte{
	0x60, 0x3d, 0xeb, 0x10, 0x15, 0xca, 0x71, 0xbe, 0x2b, 0x73, 0xae, 0xf0, 0x85, 0x7d, 0x77, 0x81,
	0x1f, 0x35, 0x2c, 0x07, 0x3b, 0x61, 0x08, 0xd7, 0x2d, 0x98, 0x10, 0xa3, 0x09, 0x14, 0xdf, 0xf4,
}

var commonIV = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}

var cbcAESTests = struct {
	name string
	key  []byte
	iv   []byte
	in   []byte
	out  []byte
}{
	"CBC-AES256",
	commonKey256,
	commonIV,
	commonInput,
	[]byte{
		0xf5, 0x8c, 0x4c, 0x04, 0xd6, 0xe5, 0xf1, 0xba, 0x77, 0x9e, 0xab, 0xfb, 0x5f, 0x7b, 0xfb, 0xd6,
		0x9c, 0xfc, 0x4e, 0x96, 0x7e, 0xdb, 0x80, 0x8d, 0x67, 0x9f, 0x77, 0x7b, 0xc6, 0x70, 0x2c, 0x7d,
		0x39, 0xf2, 0x33, 0x69, 0xa9, 0xd9, 0xba, 0xcf, 0xa5, 0x30, 0xe2, 0x63, 0x04, 0x23, 0x14, 0x61,
		0xb2, 0xeb, 0x05, 0xe2, 0xc3, 0x9b, 0xe9, 0xfc, 0xda, 0x6c, 0x19, 0x07, 0x8c, 0x6a, 0x9d, 0x1b,
	},
}

var errorKey = []byte{
	0x12, 0x23, 0x34, 0x56,
}

func TestAesEncryptCBC(t *testing.T) {
	type result struct {
		data []byte
		err error
	}

	testCases := []struct{
		name string
		given func()([]byte, error)
		expect result
	}{
		{
			name:"errorKey",
			given: func() ([]byte, error) {
				return AesEncryptCBC(cbcAESTests.iv, errorKey, cbcAESTests.in)
			},
			expect:result{[]byte(nil), aes.KeySizeError(4)},
		},
		{
			name:"AesEncryptCBCRight",
			given: func() ([]byte, error) {
				return AesEncryptCBC(cbcAESTests.iv, cbcAESTests.key, cbcAESTests.in)
			},
			expect:result{cbcAESTests.out, nil},
		},
	}
	for _,tc := range testCases{
		t.Log("TestAesDecryptCBC", tc.name)
		data, err := tc.given()
		assert.Equal(t , tc.expect.err, err)
		assert.Equal(t, tc.expect.data, data)
	}

}

func TestAesDecryptCBC(t *testing.T) {
	type result struct {
		data []byte
		err error
	}

	testCases := []struct{
		name string
		given func()([]byte, error)
		expect result
	}{
		{
			name:"errCipherText",
			given: func() ([]byte, error) {
				return AesDecryptCBC(cbcAESTests.iv, cbcAESTests.key, cbcAESTests.out[:12])
			},
			expect:result{[]byte(nil), gerror.ErrAESInvalidParameter},
		},
		{
			name:"errKey",
			given: func() ([]byte, error) {
				return AesDecryptCBC(cbcAESTests.iv, errorKey, cbcAESTests.out)
			},
			expect:result{[]byte(nil), aes.KeySizeError(4)},
		},
		{
			name:"AesDecryptCBCRight",
			given: func() ([]byte, error) {
				return AesDecryptCBC(cbcAESTests.iv, cbcAESTests.key, cbcAESTests.out)
			},
			expect:result{cbcAESTests.in, nil},
		},
	}


	for _,tc := range testCases{
		t.Log("TestAesDecryptCBC", tc.name)
		data, err := tc.given()
		assert.Equal(t , tc.expect.err, err)
		assert.Equal(t, tc.expect.data, data)
	}

}
