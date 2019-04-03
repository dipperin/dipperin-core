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


package trie

import (
	"testing"
	"math/big"
	"github.com/tidwall/gjson"
	"fmt"
	"github.com/dipperin/dipperin-core/common/util"
	"strconv"
	"errors"
	"github.com/stretchr/testify/assert"
	"reflect"
)

type erc20 struct {
	Owners  []string `json:"owners"`
	Balance map[string]*big.Int `json:"balance"`
	Name    string `json:"name"`
	Dis     uint64 `json:"dis"`
}

func TestJson2kv(t *testing.T) {
	veryBig, ok := big.NewInt(0).SetString("999999999999999999999999999999999999999999999999999999999999999999999", 10)
	assert.True(t, ok)
	c := erc20{
		Owners: []string{"0x321aa", "0xaf871"},
		Balance: map[string]*big.Int{
			"0x123":  veryBig,
			"0x1231": big.NewInt(123321),
		},
		Name: "jk",
		Dis:  10002,
	}
	kv := obj2KV2(c)
	cJsonStr, err := kv2JsonStr(kv)
	assert.NoError(t, err)

	var tmpC erc20
	err = util.ParseJson(cJsonStr, &tmpC)
	assert.NoError(t, err)
	assert.Equal(t, c.Balance["0x123"], tmpC.Balance["0x123"])
}

func obj2KV2(obj interface{}) map[string]string {
	j := gjson.ParseBytes(util.StringifyJsonToBytes(obj))
	result := map[string]string{}
	json2KV2("", j, result)
	return result
}

func json2KV2(key string, json gjson.Result, result map[string]string) {
	if json.IsObject() {
		json.ForEach(func(key1, value gjson.Result) bool {
			if key == "" {
				json2KV2(key1.Str, value, result)
			} else {
				json2KV2(key+objSegS+key1.Str, value, result)
			}
			return true
		})

	} else if json.IsArray() {
		index := 0
		json.ForEach(func(key1, value gjson.Result) bool {
			if key == "" {
				json2KV2(strconv.Itoa(index), value, result)
			} else {
				json2KV2(key+arrSegS+strconv.Itoa(index), value, result)
			}
			index ++
			return true
		})

	} else {
		if json.Type != gjson.String {
			result[key] = json.String()
		} else {
			result[key] = `"` + json.String() + `"`
		}
	}
}


const (
	dataTypeNormal = iota
	dataTypeObj
	dataTypeArr
)

const (
	objSegC = '.'
	arrSegC = '^'
	objSegS = string(objSegC)
	arrSegS = string(arrSegC)
)

type jNode struct {
	key string
	nType int
	value string
	children []absNode
}

func (n *jNode) GetChild(cKey string) absNode {
	for _, c := range n.children {
		if c.Key() == cKey {
			return c
		}
	}
	return nil
}

func (n *jNode) AddChild(c absNode) {
	n.children = append(n.children, c)
}

func (n *jNode) SetValue(v string) {
	n.value = v
}

func (n *jNode) Key() string {
	return n.key
}

func (n *jNode) Value() string {
	return n.value
}

func (n *jNode) NType() int {
	return n.nType
}

func (n *jNode) Children() []absNode {
	return n.children
}

//
func kv2JsonStr(kv map[string]string) (result string, err error) {
	// children in root is parsed in kV
	root := &jNode{ nType: dataTypeObj, children: []absNode{} }
	if err = parseKVToJNode(root, kv); err != nil {
		return
	}

	// Parsing jNode and convert it into JSON string
	result, err = extractObj(root)
	return
}

type nodeInfoCache struct {
	n absNode
	kv map[string]string
}

func parseKVToJNode(parent absNode, kv map[string]string) error {
	childrenInfoCaches := map[string]*nodeInfoCache{}

	for k, v := range kv {
		k1, extKey, dt := extractKey(k)

		// add child to parent
		oldC := parent.GetChild(k1)
		rv := reflect.ValueOf(oldC)
		if !rv.IsValid() || rv.IsNil() {
			oldC = &jNode{ key: k1, nType: dt, children: []absNode{} }
			parent.AddChild(oldC)
		}
		if oldC.NType() != dt {
			return errors.New("data type not match for key: " + k1)
		}

		// extract extKey and value
		switch dt {
		case dataTypeNormal:
			oldC.SetValue(v)
		default:
			cache := childrenInfoCaches[k1]
			if cache == nil {
				cache = &nodeInfoCache{ n: oldC, kv: map[string]string{} }
				childrenInfoCaches[k1] = cache
			}
			cache.kv[extKey] = v
		}

	}

	for _, cache := range childrenInfoCaches {
		if err := parseKVToJNode(cache.n, cache.kv); err != nil {
			return err
		}
	}

	return nil
}

func extractKey(k string) (x, extKey string, dataType int) {
	kLen := len(k)
	for i := 0; i < kLen; i++ {
		tmpC := k[i]
		if tmpC == objSegC {
			return k[:i], k[i+1:], dataTypeObj
		}
		if  tmpC == arrSegC {
			return k[:i], k[i+1:], dataTypeArr
		}
	}
	return k, "", dataTypeNormal
}

type absNode interface {
	Key() string
	Value() string
	NType() int
	Children() []absNode

	GetChild(cKey string) absNode
	AddChild(c absNode)
	SetValue(v string)
}

func extractObj(n absNode) (string, error) {
	if n.NType() != dataTypeObj {
		return "", errors.New(fmt.Sprintf("extractArr expect obj data, bug got: %v", n.NType()))
	}

	result := "{"

	for _, v := range n.Children() {
		result += `"` + v.Key() + `":`
		switch v.NType() {
		case dataTypeObj:
			if tmpObjStr, err := extractObj(v); err != nil {
				return "", err
			} else {
				result += tmpObjStr
			}
		case dataTypeArr:
			if tmpArrStr, err := extractArr(v); err != nil {
				return "", err
			} else {
				result += tmpArrStr
			}
		case dataTypeNormal:
			result += v.Value()
		}
		result += ","
	}

	// remove the last comma
	result = result[:len(result) - 1]
	result += "}"

	return result, nil
}

func extractArr(n absNode) (string, error) {
	if n.NType() != dataTypeArr {
		return "", errors.New(fmt.Sprintf("extractArr expect arr data, bug got: %v", n.NType()))
	}

	children := n.Children()
	cLen := len(children)
	dataArr := make([]string, cLen)

	for _, v := range children {
		// TODO dealing with Index Transboundary Problem
		dataIndex, err := strconv.ParseInt(v.Key(), 10, 64)
		if err != nil {
			return "", errors.New(fmt.Sprintf("invalid arr index: %v, err: %v", v.Key(), err))
		}
		index := int(dataIndex)
		if index >= cLen {
			return "", errors.New(fmt.Sprintf("data index can't bigger than arr len, len: %v index: %v", cLen, v.Key()))
		}
		if dataArr[index] != "" {
			return "", errors.New(fmt.Sprintf("duplicate value at index: %v", v.Key()))
		}
		switch v.NType() {
		case dataTypeObj:
			if tmpObjStr, err := extractObj(v); err != nil {
				return "", err
			} else {
				dataArr[index] = tmpObjStr
			}
		case dataTypeArr:
			if tmpArrStr, err := extractArr(v); err != nil {
				return "", err
			} else {
				dataArr[index] = tmpArrStr
			}
		case dataTypeNormal:
			dataArr[index] = v.Value()
		}
	}

	result := "["
	for i := 0; i < cLen; i++ {
		result += dataArr[i] + ","
	}

	// remove the last comma
	result = result[:len(result) - 1]
	result += "]"

	return result, nil
}
