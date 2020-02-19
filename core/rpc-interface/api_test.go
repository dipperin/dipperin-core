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

package rpc_interface

import (
	"fmt"
	"reflect"
	"testing"
)

func TestGetMethodName(t *testing.T) {
	t.Skip("no run this test")
	var api interface{} = (*DipperinVenusApi)(nil)

	value := reflect.ValueOf(api)
	typ := value.Type()

	for i := 0; i < value.NumMethod(); i++ {
		fmt.Printf(`{Text: "%s", Description: ""}`, typ.Method(i).Name)
		fmt.Println(",")
	}
}
