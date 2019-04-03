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


package contract

import (
	"testing"
	"bufio"
	"bytes"
	"reflect"
	"github.com/stretchr/testify/assert"
	"github.com/dipperin/dipperin-core/common"
	"fmt"
	"github.com/dipperin/dipperin-core/common/util"
)

type RWC struct {
	*bufio.ReadWriter
}

func (rwc *RWC) Close() error {
	return nil
}

func TestParamsCodec_ParseRequestArguments(t *testing.T) {
	req := bytes.NewBufferString(`[11, 22]`)

	codec := NewParamsCodec(bufio.NewReader(req))
	assert.NotNil(t, codec)
	rValue, pErr := codec.ParseRequestArguments([]reflect.Type{reflect.TypeOf(1), reflect.TypeOf(1)})
	assert.NoError(t, pErr)
	assert.Equal(t, 11, rValue[0].Interface())
	assert.Equal(t, 22, rValue[1].Interface())

	rValue, pErr = codec.ParseRequestArguments([]reflect.Type{reflect.TypeOf(1), reflect.TypeOf('1')})
	assert.Error(t, pErr)
	assert.Nil(t, rValue)

	req = bytes.NewBufferString(`["0x0000918c773880B462929ACE4F975CcfED9Be2d8Efc9",1000000000]`)
	codec = NewParamsCodec(bufio.NewReader(req))
	rValue, pErr = codec.ParseRequestArguments([]reflect.Type{reflect.TypeOf(common.HexToAddress("0x0000918c773880B462929ACE4F975CcfED9Be2d8Efc9")), reflect.TypeOf(1)})
	fmt.Println(rValue[0].Interface(), rValue[1].Interface())
	fmt.Println(pErr)

	params := fmt.Sprintf("[\"%v\",%v]", common.HexToAddress("0x0000918c773880B462929ACE4F975CcfED9Be2d8Efc9"), 12345)
	fmt.Println(params)

	req = bytes.NewBufferString(params)
	codec = NewParamsCodec(bufio.NewReader(req))
	rValue, pErr = codec.ParseRequestArguments([]reflect.Type{reflect.TypeOf(common.HexToAddress("0x0000918c773880B462929ACE4F975CcfED9Be2d8Efc9")), reflect.TypeOf(1)})
	fmt.Println(rValue[0].Interface(), rValue[1].Interface())
	fmt.Println(pErr)

	destStr := fmt.Sprintf("%v", common.HexToAddress("0x0000918c773880B462929ACE4F975CcfED9Be2d8Efc9"))
	params = util.StringifyJson([]interface{}{ destStr, 5678 })
	fmt.Println(params)
	req = bytes.NewBufferString(params)
	codec = NewParamsCodec(bufio.NewReader(req))
	rValue, pErr = codec.ParseRequestArguments([]reflect.Type{reflect.TypeOf(common.HexToAddress("0x0000918c773880B462929ACE4F975CcfED9Be2d8Efc9")), reflect.TypeOf(1)})
	fmt.Println(rValue[0].Interface(), rValue[1].Interface())
	fmt.Println(pErr)

	reqe := bytes.NewBufferString(`k`)
	codece := NewParamsCodec(bufio.NewReader(reqe))
	rValue, pErr = codece.ParseRequestArguments([]reflect.Type{reflect.TypeOf(1)})
	assert.Error(t, pErr)
	assert.Nil(t, rValue)

	reqe = bytes.NewBufferString(`["k"]`)
	codece = NewParamsCodec(bufio.NewReader(reqe))
	rValue, pErr = codece.ParseRequestArguments([]reflect.Type{reflect.TypeOf(1)})
	assert.Error(t, pErr)
	assert.Nil(t, rValue)

	reqe = bytes.NewBufferString(`"k"`)
	codece = NewParamsCodec(bufio.NewReader(reqe))
	rValue, pErr = codece.ParseRequestArguments([]reflect.Type{reflect.TypeOf(1)})
	assert.Error(t, pErr)
	assert.Nil(t, rValue)

	reqe = bytes.NewBufferString(`["k"]`)
	codece = NewParamsCodec(bufio.NewReader(reqe))
	rValue, pErr = codece.ParseRequestArguments([]reflect.Type{reflect.TypeOf("l"), reflect.TypeOf(1)})
	assert.Error(t, pErr)
	assert.Nil(t, rValue)

	reqe = bytes.NewBufferString(`[11,"k"]`)
	codece = NewParamsCodec(bufio.NewReader(reqe))
	rValue, pErr = codece.ParseRequestArguments([]reflect.Type{reflect.TypeOf(1)})
	assert.Error(t, pErr)
	assert.Nil(t, rValue)


}
