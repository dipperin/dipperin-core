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
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

type header struct {
	Num   uint64 `json:"num"`
	Msg   string `json:"msg"`
	Extra []byte `json:"extra"`
}

type body struct {
	Txs []transaction `json:"txs"`
}

type transaction struct {
	Id     int      `json:"id"`
	Amount *big.Int `json:"amount"`
}

type block struct {
	Header header `json:"header"`
	Body   body   `json:"body"`
}

// JSON flattened into K-V form
func TestJson2Kv(t *testing.T) {
	b := &block{
		Header: header{
			Num:   1,
			Msg:   "hi",
			Extra: []byte("hello"),
		},
		Body: body{
			Txs: []transaction{
				{Id: 1, Amount: big.NewInt(1321)}, {Id: 2, Amount: big.NewInt(13210)},
			},
		},
	}
	bB := util.StringifyJsonToBytes(b)
	bStr := string(bB)
	fmt.Println(bStr)

	var tmp block
	assert.NoError(t, util.ParseJsonFromBytes(bB, &tmp))
	assert.Equal(t, "hello", string(tmp.Header.Extra))
}
