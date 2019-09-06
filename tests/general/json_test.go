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

package general

import (
	"fmt"
	"github.com/tidwall/gjson"
	"strconv"
	"testing"
)

const jsonStr = `{"header":{"num":1,"msg":"hi","extra":"aGVsbG8="},"body":{"txs":[{"id":1,"amount":1321},{"id":2,"amount":13210}]}}`

func TestJson(t *testing.T) {
	value := gjson.Parse(jsonStr)

	result := map[string]interface{}{}

	json2KV("", value, result)

	for k, v := range result {
		fmt.Println(k, ": ", v)
	}

}

func json2KV(key string, json gjson.Result, result map[string]interface{}) {
	if json.IsObject() {

		json.ForEach(func(key1, value gjson.Result) bool {

			if key == "" {
				json2KV(key1.Str, value, result)
			} else {
				json2KV(key+"."+key1.Str, value, result)
			}

			return true
		})

	} else if json.IsArray() {
		index := 0
		json.ForEach(func(key1, value gjson.Result) bool {

			if key == "" {
				json2KV(strconv.Itoa(index), value, result)
			} else {
				json2KV(key+"."+strconv.Itoa(index), value, result)
			}

			index++
			return true
		})

	} else {
		result[key] = json.Value()
	}
}
