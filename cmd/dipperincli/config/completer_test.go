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
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestDipperinCliCompleter(t *testing.T) {
	log.InitLogger(log.LvlDebug)
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
	b.InsertText("rpc.A -A -F", false, true)

	d = *b.Document()

	DipperinCliCompleterNew(d)
	//assert.Equal(t, DipperinCliCompleterNew(d), []prompt.Suggest{})


	b = prompt.NewBuffer()
	b.InsertText("RPC", false, true)

	d = *b.Document()

	assert.Equal(t, DipperinCliCompleter(d), []prompt.Suggest{{Text:"-h", Description:""}, {Text:"--help", Description:""}})
}


func TestDipperinCliCompleterNew(t *testing.T) {
	log.InitLogger(log.LvlDebug)
	d := prompt.Document{}


	args := strings.Split("tx ", " ")
	fmt.Println(len(args))
	b := prompt.NewBuffer()
	b.InsertText("tx ", false, true)

	d = *b.Document()

	suggest := DipperinCliCompleterNew(d)
	log.Debug("TestDipperinCliCompleterNew", "suggest", suggest)


	args = strings.Split("tx SendTx ", " ")
	fmt.Println(len(args))
	b = prompt.NewBuffer()
	b.InsertText("tx SendTx ", false, true)

	d = *b.Document()

	suggest = DipperinCliCompleterNew(d)
	log.Debug("TestDipperinCliCompleterNew", "suggest", suggest)

	args = strings.Split("tx SendTx -", " ")
	fmt.Println(len(args))
	b = prompt.NewBuffer()
	b.InsertText("tx SendTx -", false, true)

	d = *b.Document()

	suggest = DipperinCliCompleterNew(d)
	log.Debug("TestDipperinCliCompleterNew", "suggest", suggest)


	args = strings.Split("tx SendTransactionContract --abi ", " ")
	fmt.Println(len(args), strings.TrimLeft(args[len(args)-1], "--"))
	for _,arg := range args{
		fmt.Println(strings.TrimLeft(arg, "--"))
	}
	b = prompt.NewBuffer()
	b.InsertText("tx SendTransactionContract --abi ", false, true)

	d = *b.Document()

	suggest = DipperinCliCompleterNew(d)
	log.Debug("TestDipperinCliCompleterNew", "suggest", suggest)

	b = prompt.NewBuffer()
	b.InsertText("rpc.Add -A -F", false, true)

	d = *b.Document()

	suggest = DipperinCliCompleterNew(d)
	log.Debug("TestDipperinCliCompleterNew", "suggest", suggest)


	b = prompt.NewBuffer()
	b.InsertText("rp", false, true)

	d = *b.Document()

	suggest = DipperinCliCompleterNew(d)
	log.Debug("TestDipperinCliCompleterNew", "suggest", suggest)
	//assert.Equal(t, DipperinCliCompleterNew(d), []prompt.Suggest{})
}


func Test_argumentsCompleter(t *testing.T) {
	assert.Equal(t, argumentsCompleter([]string{"test"}), []prompt.Suggest{})
	assert.Equal(t, argumentsCompleter([]string{"rpc", "a"}), []prompt.Suggest{})
	assert.Equal(t, argumentsCompleter([]string{"rpc", "a", "b"}), nilSuggest)
	fmt.Println(argumentsCompleter([]string{"rpc", "-a", "b"}))
	assert.Equal(t, argumentsCompleter([]string{"rpc", "-a"}), nilSuggest)
}

func Test_excludeOptions(t *testing.T) {
	args := []string{"-test1", "test2", "test3"}
	assert.Equal(t, excludeOptions(args), []string{"test2", "test3"})
}
