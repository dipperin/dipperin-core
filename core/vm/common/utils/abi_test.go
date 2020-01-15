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
	"github.com/dipperin/dipperin-core/tests/factory/vminfo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWasmAbi_FromJson(t *testing.T) {
	abiByte := new(WasmAbi)
	err := abiByte.FromJson(nil)
	assert.Equal(t, gerror.ErrEmptyInput, err)

	_, abi := vminfo.GetTestData("event")

	err = abiByte.FromJson(abi)
	assert.NoError(t, err)
}
