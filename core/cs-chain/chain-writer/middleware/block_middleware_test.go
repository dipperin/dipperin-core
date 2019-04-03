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


package middleware

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestChainOrBlockCannotBeNull(t *testing.T) {
	bc := NewBlockContext(nil, nil)
	assert.Equal(t, true, len(bc.middlewares) == 0)
	bc.Use(ValidateBlockNumber(bc))
	bc.Use(UpdateStateRoot(bc))
	bc.Use(UpdateBlockVerifier(bc))
	bc.Use(InsertBlock(bc))
	assert.Equal(t, true, len(bc.middlewares) == 4)
	err := bc.Process()
	assert.Error(t, err)
}
