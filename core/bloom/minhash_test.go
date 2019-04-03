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


package iblt

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var defaultHashPoolConfig = NewHashPoolConfig(16, 4000)

func TestHashPool_MinHash(t *testing.T) {
	const (
		intersect = 100000
		aSize     = 1000
		bSize     = 1000
	)

	e1 := NewHashPool(defaultHashPoolConfig)
	e2 := NewHashPool(defaultHashPoolConfig)

	whole := generateRandomKV(intersect + aSize + bSize)

	// construct distinctive KV pairs
	cnt := 0
	for _, v := range whole {
		if cnt < intersect {
			e1.Encode(v)
			e2.Encode(v)
		} else if cnt < intersect+aSize {
			e1.Encode(v)
		} else if cnt < intersect+aSize+bSize {
			e2.Encode(v)
		}
		cnt++
	}

	diff := e1.Decode(e2)
	assert.True(t, int(float32(diff)*1.45) > (aSize + bSize))
}
