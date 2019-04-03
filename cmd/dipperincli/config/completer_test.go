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
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDipperinCliCompleter(t *testing.T) {
	d := prompt.Document{}

	assert.Equal(t, DipperinCliCompleter(d), nilSuggest)

	b := prompt.NewBuffer()
	b.InsertText("test", false, true)

	d = *b.Document()

	assert.Equal(t, DipperinCliCompleter(d), []prompt.Suggest{})

	b = prompt.NewBuffer()
	b.InsertText("rpc -test", false, true)

	d = *b.Document()

	assert.Equal(t, DipperinCliCompleter(d), []prompt.Suggest{})

	b = prompt.NewBuffer()
	b.InsertText("RPC", false, true)

	d = *b.Document()

	assert.Equal(t, DipperinCliCompleter(d), []prompt.Suggest{{Text:"-h", Description:""}, {Text:"--help", Description:""}})
}

func Test_argumentsCompleter(t *testing.T) {
	assert.Equal(t, argumentsCompleter([]string{"test"}), []prompt.Suggest{})
	assert.Equal(t, argumentsCompleter([]string{"rpc", "a"}), []prompt.Suggest{})
	assert.Equal(t, argumentsCompleter([]string{"rpc", "a", "b"}), nilSuggest)
}

func Test_excludeOptions(t *testing.T) {
	args := []string{"-test1", "test2", "test3"}
	assert.Equal(t, excludeOptions(args), []string{"test2", "test3"})
}
