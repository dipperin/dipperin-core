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

package config

import (
	"github.com/c-bata/go-prompt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_optionCompleter(t *testing.T) {
	assert.Equal(t, optionCompleterNew([]string{}, false), optionHelp)
	assert.Equal(t, optionCompleterNew([]string{"--test1"}, true), []prompt.Suggest{{Text: "--help", Description: ""}})
	assert.Equal(t, optionCompleterNew([]string{"rpc", "-test"}, false), []prompt.Suggest{})
	assert.Equal(t, optionCompleterNew([]string{"rpc", "--test"}, true), []prompt.Suggest{})
}

//func Test_callMethod(t *testing.T) {
//	assert.Equal(t, callMethod([]string{}, false), optionHelp)
//	assert.Equal(t, callMethod([]string{}, true), []prompt.Suggest{{Text: "--help", Description: ""}})
//	assert.Equal(t, callMethod([]string{"rpc", "-test1", "-test2"}, false), []prompt.Suggest{})
//	assert.Equal(t, callMethod([]string{"rpc", "--test1", "--test2"}, true), []prompt.Suggest{})
//}
