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
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewGraphene(t *testing.T) {
	c := DeriveBloomConfig(32000)
	i := NewInvBloomConfig(100, 4)
	g := NewGraphene(i, c)

	assert.Equal(t, g.Bloom().config.BloomBits, uint(4))
	assert.Equal(t, g.InvBloom().config.BucketNum, uint(100))
}

func TestGraphene_InsertRLP(t *testing.T) {
	c := DeriveBloomConfig(32000)
	i := NewInvBloomConfig(100, 4)
	g := NewGraphene(i, c)
	g.InsertRLP("test", "test")
	assert.Equal(t, g.InvBloom().buckets[8].Count, int32(1))
}

func TestGraphene_InvBloomConfig(t *testing.T) {
	c := DeriveBloomConfig(32000)
	i := NewInvBloomConfig(100, 4)
	g := NewGraphene(i, c)
	assert.Equal(t, g.InvBloomConfig(), i)
}

func TestGraphene_BloomConfig(t *testing.T) {
	c := DeriveBloomConfig(32000)
	i := NewInvBloomConfig(100, 4)
	g := NewGraphene(i, c)
	assert.Equal(t, g.BloomConfig(), c)
}

func TestGraphene_FilterRLP(t *testing.T) {

	c := DeriveBloomConfig(16)
	i := NewInvBloomConfig(10, 4)
	g := NewGraphene(i, c)

	g.InsertRLP("t", "test1")

	filterMap := make(map[interface{}]interface{})
	filterMap["t"] = "test1"
	filterMap["test2"] = "test2"

	assert.Equal(t, len(g.FilterRLP(filterMap)), 1)

	filterMap2 := make(map[interface{}]interface{})
	filterMap2["t"] = "test1"
	filterMap2[0] = "test2"

	assert.Equal(t, len(g.FilterRLP(filterMap2)), 0)

	filterMap3 := make(map[interface{}]interface{})
	filterMap3["t"] = "test1"
	filterMap3[nil] = "test2"

	assert.Equal(t, len(g.FilterRLP(filterMap3)), 0)
}

func TestGraphene_FilterListRLP(t *testing.T) {
	c := DeriveBloomConfig(16)
	i := NewInvBloomConfig(10, 4)
	g := NewGraphene(i, c)

	g.InsertRLP("t", "1")
	g.InsertRLP("a", "2")

	filterMap := make(map[interface{}]interface{})

	filterMap["t"] = "1"
	filterMap["a"] = "1"
	filterMap["b"] = "3"

	a, err := g.FilterListRLP(filterMap)

	assert.NoError(t, err)

	assert.Equal(t, a, [][]byte{[]byte("2")})
}

func TestGraphene_ListRLP(t *testing.T) {
	c := DeriveBloomConfig(16)
	i := NewInvBloomConfig(10, 4)
	g := NewGraphene(i, c)

	g.InsertRLP("1", "1")
	g.InsertRLP("2", "2")

	a, b, err := g.ListRLP()

	assert.NoError(t, err)

	assert.Equal(t, a, [][]byte{[]byte("2"), []byte("1")})

	assert.Equal(t, b, [][]byte(nil))

}

func TestGraphene_Recover(t *testing.T) {
	c := DeriveBloomConfig(16)
	i := NewInvBloomConfig(10, 4)
	g := NewGraphene(i, c)

	g.InsertRLP("t", "1")
	g.InsertRLP("a", "2")

	filterMap := make(map[interface{}]interface{})

	filterMap["t"] = "1"
	filterMap["a"] = "2"
	filterMap["b"] = "3"

	a, err := g.Recover(filterMap)

	assert.NoError(t, err)

	isCorrect := 0

	for _,b := range a {
		if string(b) == "1" || string(b) == "2" {
			isCorrect++
		}
	}

	assert.Equal(t, isCorrect, 2)
}
