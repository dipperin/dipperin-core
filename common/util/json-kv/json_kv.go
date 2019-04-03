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


package json_kv

import (
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/tidwall/gjson"
	"strconv"
	"reflect"
	"errors"
	"fmt"
	"strings"
)

const (
	objSegC = '.'
	arrSegC = '^'
	objSegS = string(objSegC)
	arrSegS = string(arrSegC)

	invalidKeyChars = objSegS + arrSegS
)

func Obj2KV(obj interface{}) (map[string]string, error) {
	j := gjson.ParseBytes(util.StringifyJsonToBytes(obj))
	return Json2KV(j)
}

func JsonBytes2KV(b []byte) (map[string]string, error) {
	j := gjson.ParseBytes(b)
	return Json2KV(j)
}

func JsonStr2KV(str string) (map[string]string, error) {
	j := gjson.Parse(str)
	return Json2KV(j)
}

func Json2KV(json gjson.Result) (map[string]string, error) {

	result := map[string]string{}
	err := json2KV("", json, result)
	return result, err
}

// Determine if the key is an illegal key
func invalidJsonKey(key string) bool {
	return strings.ContainsAny(key, invalidKeyChars)
}

// todo All keys in obj can't take two special characters, or you can escape special characters. Do more border checks
func json2KV(key string, json gjson.Result, result map[string]string) (err error) {
	if json.IsObject() {
		json.ForEach(func(key1, value gjson.Result) bool {
			if key == "" {
				err = json2KV(key1.Str, value, result)
			} else {
				err = json2KV(key + objSegS + key1.Str, value, result)
			}
			return true
		})

	} else if json.IsArray() {
		index := 0
		json.ForEach(func(key1, value gjson.Result) bool {
			if key == "" {
				err = json2KV(strconv.Itoa(index), value, result)
			} else {
				err = json2KV(key + arrSegS + strconv.Itoa(index), value, result)
			}
			index++
			return true
		})

	} else {
		if json.Type != gjson.String {
			result[key] = json.String()
		} else {
			//result[key] = `"` + json.String() + `"`
			// Handling strings that need to be escaped
			result[key] = util.StringifyJson(json.String())
			//fmt.Println("===", key, result[key])
		}
	}
	return nil
}

const (
	dataTypeNormal = iota
	dataTypeObj
	dataTypeArr
)


func KV2JsonObj(kv map[string]string, result interface{}) error {
	str, err := KV2JsonStr(kv)
	if err != nil {
		return err
	}

	return util.ParseJson(str, result)
}

func KV2JsonStr(kv map[string]string) (result string, err error) {
	// The children in root are parsed in kv
	root := &jNode{ nType: dataTypeObj, children: []absNode{} }
	if err = parseKVToJNode(root, kv); err != nil {
		return
	}

	// Parsing jNode and finally converting it into json string
	result, err = extractObj(root)
	return
}

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
	resultLen := len(result)
	if result[resultLen - 1] == ',' {
		result = result[:resultLen - 1]
	}
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
		index, err := strconv.Atoi(v.Key())

		if index < 0 || err != nil {
			return "", errors.New(fmt.Sprintf("invalid arr index: %v, err: %v", v.Key(), err))
		}
		if index >= cLen {
			return "", errors.New(fmt.Sprintf("data index can't bigger than arr len, len: %v index: %v", cLen, v.Key()))
		}

		// Map structure cannot have the same key
		//if dataArr[index] != "" {
		//	return "", errors.New(fmt.Sprintf("duplicate value at index: %v", v.Key()))
		//}
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

	// finishing results
	result := "["
	for i := 0; i < cLen; i++ {
		result += dataArr[i] + ","
	}
	// remove the last comma
	resultLen := len(result)
	if result[resultLen - 1] == ',' {
		result = result[:resultLen - 1]
	}
	result += "]"

	return result, nil
}