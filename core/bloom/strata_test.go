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

var defaultEstimatorConfig = NewEstimatorConfig(16)

func TestEstimator_TrailingZeros(t *testing.T) {
	e := NewEstimator(defaultEstimatorConfig)
	h := e.NewDataHash()

	h.SetBytes([]byte{0, 0, 0, 1})

	for i := 0; i < 26; i++ {
		zeros := e.TrailingZeros(h)
		assert.EqualValues(t, 8*len(h)-i-1, zeros)
		h.lsh()
	}
}

func TestEstimator_DecodeData(t *testing.T) {
	const (
		intersect = 100
		aSize     = 100
		bSize     = 100
	)

	e1 := NewEstimator(defaultEstimatorConfig)
	e2 := NewEstimator(defaultEstimatorConfig)

	whole := generateRandomKV(intersect + aSize + bSize)

	// construct distinctive KV pairs
	cnt := 0
	for _, v := range whole {
		if cnt < intersect {
			e1.EncodeData(v)
			e2.EncodeData(v)
		} else if cnt < intersect+aSize {
			e1.EncodeData(v)
		} else if cnt < intersect+aSize+bSize {
			e2.EncodeData(v)
		}
		cnt++
	}

	diff := e1.DecodeData(e2)
	assert.True(t, int(float32(diff)*1.45) > (aSize+bSize))
}

//added by caiqingfeng 2019.1.21
func TestString(t *testing.T) {
	e1 := NewEstimator(defaultEstimatorConfig)
	println(e1.String())
}
