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

// This is a duplicated and slightly modified version of "gopkg.in/karalabe/cookiejar.v2/collections/prque".

package prque

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	p := New(func(a interface{}, i int) {})
	x := 1
	p.Push(x, 2)
	pX, pri := p.Pop()
	assert.Equal(t, 1, pX)
	assert.Equal(t, int64(2), pri)
	p.Push(x, 2)
	pX = p.PopItem()
	assert.Equal(t, 1, pX)
	assert.Nil(t, p.Remove(-1))
	p.Push(x, 2)
	assert.Equal(t, 1, p.Size())
	assert.NotNil(t, p.Remove(0))
	assert.True(t, p.Empty())
	p.Reset()
	p.cont.Reset()
}
