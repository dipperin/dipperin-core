// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package accountsbase

import (
	"fmt"
	"github.com/dipperin/dipperin-core/common/gerror"
	"github.com/issue9/assert"
	"math"
	"testing"
)

// Tests that HD derivation paths can be correctly parsed into our internal binary
// representation.
func TestHDPathParsing(t *testing.T) {
	testCases := []struct {
		name   string
		given  string
		output DerivationPath
		err    error
	}{
		{"ErrEmptyDerivedPath","", nil, gerror.ErrDerivedPath},
		{"ErrDerivedPath", "m", nil, gerror.ErrEmptyDerivedPath},
		{"ErrInvalidcomponent", "m/xyzc", nil, fmt.Errorf("invalid component: %s", "xyzc")},
		{"ErrOverflowsHardenedRange", "m/2147483648'", nil, fmt.Errorf("component %v out of allowed hardened range [0, %d]", 2147483648, math.MaxUint32-0x80000000)},
		{"ErrOverflowsAllowedRange", "m/214748364811", nil, fmt.Errorf("component %v out of allowed range [0, %d]", 214748364811, math.MaxUint32)}, // Overflows 32 bit integer
		{"ParseDerivationPathRight","m/44'/60'/0'/0", DerivationPath{0x80000000 + 44, 0x80000000 + 60, 0x80000000 + 0, 0}, nil},
	}
	for _, tt := range testCases {
		path, err := ParseDerivationPath(tt.given)
		if err == nil {
			assert.Equal(t, tt.output, path)
		} else {
			assert.Equal(t, tt.err, err)
		}
	}
}

func TestDerivationPath_String(t *testing.T) {
	testCases := []struct {
		name   string
		given  DerivationPath
		output string
	}{
		{"StringRight", DerivationPath{0x80000000 + 44, 0x80000000 + 60, 0x80000000 + 0, 0}, "m/44'/60'/0'/0"},
	}

	for _, tt := range testCases{
		assert.Equal(t, tt.output, tt.given.String())
	}
}