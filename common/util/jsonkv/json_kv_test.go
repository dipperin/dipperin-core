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

package jsonkv

import (
	"fmt"
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/stretchr/testify/assert"
	"math/big"
	"reflect"
	"testing"
)

type erc20 struct {
	// todo special characters cause conversion errors
	Owners  []string            `json:"owne.rs"`
	Balance map[string]*big.Int `json:"balance"`
	Name    string              `json:"name"`
	Name2   string              `json:"name2"`
	Dis     uint64              `json:"dis"`
}

var kvResult = map[string]string{
	"owne.rs^0":      `"0x321aa"`,
	"owne.rs^1":      `"0xaf871"`,
	"balance.0x123":  "999999999999999999999999999999999999999999999999999999999999999999999",
	"balance.0x1231": "123321",
	"name":           `"\njk\""`,
	"name2":          `"asdfdsf"`,
	"dis":            "10002",
}

func buildTestStruct() erc20 {
	veryBig, _ := big.NewInt(0).SetString("999999999999999999999999999999999999999999999999999999999999999999999", 10)
	c := erc20{
		Owners: []string{"0x321aa", "0xaf871"},
		Balance: map[string]*big.Int{
			"0x123":  veryBig,
			"0x1231": big.NewInt(123321),
		},
		Name:  "\njk\"",
		Name2: "asdfdsf",
		Dis:   10002,
	}
	return c
}

func TestObj2KV(t *testing.T) {
	c := buildTestStruct()

	kv, err := Obj2KV(c)
	assert.NoError(t, err)

	assert.True(t, reflect.DeepEqual(kvResult, kv))

}

func TestJsonBytes2KV(t *testing.T) {
	c := buildTestStruct()

	kv, err := JsonBytes2KV(util.StringifyJsonToBytes(c))

	assert.NoError(t, err)
	assert.True(t, reflect.DeepEqual(kvResult, kv))
}

func TestJsonStr2KV(t *testing.T) {
	c := buildTestStruct()

	kv, err := JsonStr2KV(util.StringifyJson(c))

	assert.NoError(t, err)
	assert.True(t, reflect.DeepEqual(kvResult, kv))
}

func TestJson2kv(t *testing.T) {
	c := buildTestStruct()

	fmt.Println(util.StringifyJson(c))

	kv, err := Obj2KV(c)
	assert.NoError(t, err)
	// assemble the result into a json string
	//for k, v := range kv {
	//	fmt.Println(k, v)
	//}
	_, err = KV2JsonStr(kv)
	assert.NoError(t, err)
	//fmt.Println(cJsonStr)

	var tmpC erc20
	err = KV2JsonObj(kv, &tmpC)
	assert.NoError(t, err)
	assert.Equal(t, c.Balance["0x123"], tmpC.Balance["0x123"])
	assert.Equal(t, "\njk\"", tmpC.Name)
	assert.Equal(t, "asdfdsf", tmpC.Name2)
	fmt.Println(util.StringifyJson(tmpC))
}

func TestInvalidJsonKey(t *testing.T) {
	assert.True(t, invalidJsonKey("sdf"+objSegS))
	assert.True(t, invalidJsonKey("sdf"+arrSegS))
	assert.False(t, invalidJsonKey("sdf"))
}

func TestKV2JsonStr(t *testing.T) {
	kv := map[string]string{
		"owne.rs^0":      `"0x321aa"`,
		"owne.rs^1":      `"0xaf871"`,
		"owne.rs^5":      `"0x111"`,
		"balance.0x123":  "999999999999999999999999999999999999999999999999999999999999999999999",
		"balance.0x1231": "123321",
		"name":           `"\njk\""`,
		"name2":          `"asdfdsf"`,
		"dis":            "10002",
	}

	_, err := KV2JsonStr(kv)

	assert.Error(t, err)

	kv2 := map[string]string{
		"owne.rs^0":      `"0x321aa"`,
		"owne.rs^1":      `"0xaf871"`,
		"owne.rs^-1":     `"0x111"`,
		"balance.0x123":  "999999999999999999999999999999999999999999999999999999999999999999999",
		"balance.0x1231": "123321",
		"name":           `"\njk\""`,
		"name2":          `"asdfdsf"`,
		"dis":            "10002",
	}

	_, err = KV2JsonStr(kv2)

	assert.Error(t, err)

	kv3 := map[string]string{
		"owne.rs^0^0":    `"0x321aa"`,
		"owne.rs^0^1":    `"0x321aa"`,
		"owne.rs^1":      `"0xaf871"`,
		"balance.0x123":  "999999999999999999999999999999999999999999999999999999999999999999999",
		"balance.0x1231": "123321",
		"name":           `"\njk\""`,
		"name2":          `"asdfdsf"`,
		"dis":            "10002",
	}

	s1, err := KV2JsonStr(kv3)

	assert.NoError(t, err)

	rKv3, err := JsonStr2KV(s1)

	assert.NoError(t, err)

	assert.True(t, reflect.DeepEqual(kv3, rKv3))

	kv4 := map[string]string{
		"owne.rs^0.a":    `"0x321aa"`,
		"owne.rs^0.b":    `"0x321aa"`,
		"owne.rs^1":      `"0xaf871"`,
		"balance.0x123":  "999999999999999999999999999999999999999999999999999999999999999999999",
		"balance.0x1231": "123321",
		"name":           `"\njk\""`,
		"name2":          `"asdfdsf"`,
		"dis":            "10002",
	}

	s2, err := KV2JsonStr(kv4)

	assert.NoError(t, err)

	rKv4, err := JsonStr2KV(s2)

	assert.NoError(t, err)

	assert.True(t, reflect.DeepEqual(kv4, rKv4))

	kv5 := map[string]string{
		"owne.rs^0":      `"0x321aa"`,
		"owne.rs^1":      `"0xaf871"`,
		"owne.rs^2^0":    `"0x111"`,
		"owne.rs^2^5":    `"0x111"`,
		"balance.0x123":  "999999999999999999999999999999999999999999999999999999999999999999999",
		"balance.0x1231": "123321",
		"name":           `"\njk\""`,
		"name2":          `"asdfdsf"`,
		"dis":            "10002",
	}

	_, err = KV2JsonStr(kv5)

	assert.Error(t, err)

}
