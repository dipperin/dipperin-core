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

func TestBftMiddleware(t *testing.T) {
	bc := NewBftBlockContext(nil, nil, nil)
	//assert.Equal(t, true, len(bc.middlewares) == 0)
	assert.Len(t, bc.middlewares, 0)
	bc.Use(CheckBlock(&bc.BlockContext))
	bc.Use(ValidateBlockNumber(&bc.BlockContext))
	bc.Use(UpdateStateRoot(&bc.BlockContext))
	bc.Use(UpdateBlockVerifier(&bc.BlockContext))
	bc.Use(InsertBlock(&bc.BlockContext))
	//assert.Equal(t, true, len(bc.middlewares) == 5)
	assert.Len(t, bc.middlewares, 5)
	err := bc.Process()
	assert.Error(t, err)
}

func TestBftMiddleware2(t *testing.T) {
	x := NewBftBlockValidator(nil)
	assert.NotNil(t, x)
	assert.Error(t, x.FullValid(nil))

	assert.NotNil(t, NewBftBlockContextWithoutVotes(nil, nil))
}
