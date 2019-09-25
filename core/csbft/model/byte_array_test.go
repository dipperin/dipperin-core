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

package model

import (
	"fmt"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

//Size, GetIndex, SetIndex, String
func TestBitArray_Nil(t *testing.T) {
	ba := NewBitArray(0)
	ba.Size()
	ba.String()
	ba.GetIndex(1)
	ba.SetIndex(0, true)
	ba1 := ba.Copy()
	assert.Equal(t, ba, ba1)
}

//Size, GetIndex, SetIndex, String
func TestBitArray_NotNil(t *testing.T) {
	ba := NewBitArray(5)
	ba.Size()
	ba.String()
	ba.GetIndex(1)
	ba.GetIndex(6)
	ba.SetIndex(4, true)
	ba.SetIndex(4, false)
	ba.SetIndex(6, true)
	ba1 := ba.Copy()
	assert.Equal(t, ba, ba1)
}

func TestBitArray_stringIndented2(t *testing.T) {
	ba := NewBitArray(100)
	ba.Size()
	ba.String()

}

//Or, And, Not
func TestBitArray_Nil2(t *testing.T) {
	ba := NewBitArray(0)
	ba.Or(nil)
	ba1 := NewBitArray(2)
	ba.Or(ba1)
	ba.And(ba)
	ba.Not()
	ba.Sub(ba)
}

func TestBitArray_NotNil2(t *testing.T) {
	ba := NewBitArray(5)
	ba1 := NewBitArray(4)
	ba2 := NewBitArray(6)
	ba.Or(nil)
	ba.Or(ba)
	ba.And(ba)
	ba.Not()
	ba.Sub(ba1)
	ba.Sub(ba2)
	//time.Sleep(time.Second*1)
	//ba.Sub(ba1)
}

func TestNewBitArray_IsEmpty(t *testing.T) {
	ba := NewBitArray(0)
	assert.Equal(t, ba.IsEmpty(), true)
	ba1 := NewBitArray(5)
	assert.Equal(t, ba1.IsEmpty(), true)
	ba1.Elems[0] = 12
	assert.Equal(t, ba1.IsEmpty(), false)
}

func TestBitArray_IsFull2(t *testing.T) {
	ba := NewBitArray(0)
	ba.IsFull()
	ba1 := NewBitArray(2)
	ba1.Elems[0] = 1
	//ba1.Elems[1] = 2
	assert.Equal(t, ba1.IsFull(), false)
}

func TestBitArray_PickRandom2(t *testing.T) {
	ba := NewBitArray(0)
	ba.PickRandom()
	ba = NewBitArray(6)
	ba.PickRandom()
}

//Byte, Update
func TestBitArray_Bytes2(t *testing.T) {
	ba := NewBitArray(5)
	ba.Bytes()
	ba1 := NewBitArray(0)
	ba1.Update(nil)
	ba2 := NewBitArray(4)
	ba2.Update(ba)
}

func TestBitArray_MarshalJSON2(t *testing.T) {
	ba := NewBitArray(0)
	ba.MarshalJSON()
	ba1 := NewBitArray(5)
	ba1.MarshalJSON()
}

func TestBitArray_UnmarshalJSON(t *testing.T) {
	ba := NewBitArray(5)
	bz := []byte("null")
	ba.UnmarshalJSON(bz)
	bz1 := []byte("hello")
	ba.UnmarshalJSON(bz1)
	bz2 := []byte("\"_x\"")
	ba.UnmarshalJSON(bz2)
}

func TestBitArray_print(t *testing.T) {
	var bitArrayJSONRegexp = regexp.MustCompile(`\A"([_x]*)"\z`)
	var bz = []byte("\"_x\"")
	var b = string(bz)
	match := bitArrayJSONRegexp.FindStringSubmatch(b)
	log.Info("TestBitArray_print", "result", bitArrayJSONRegexp.Match([]byte("\"_x\"")))
	log.Info("TestBitArray_print", "result2", bitArrayJSONRegexp.Match([]byte("_x___yt")))
	fmt.Println(match)
}
