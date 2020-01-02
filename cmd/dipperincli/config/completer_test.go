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
	"github.com/dipperin/dipperin-core/common/log"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"strings"
	"testing"
)

func TestDipperinCliCompleter(t *testing.T) {
	d := prompt.Document{}

	assert.Equal(t, DipperinCliCompleterNew(d), nilSuggest)

	b := prompt.NewBuffer()
	b.InsertText("test", false, true)

	d = *b.Document()

	assert.Equal(t, DipperinCliCompleterNew(d), []prompt.Suggest{})

	b = prompt.NewBuffer()
	b.InsertText("miner -test", false, true)

	d = *b.Document()

	assert.Equal(t, DipperinCliCompleterNew(d), []prompt.Suggest{})

	b = prompt.NewBuffer()
	b.InsertText("miner.A -A -F", false, true)

	d = *b.Document()

	DipperinCliCompleterNew(d)
	//assert.Equal(t, DipperinCliCompleterNew(d), []prompt.Suggest{})

	b = prompt.NewBuffer()
	b.InsertText("miner ", false, true)

	d = *b.Document()

	assert.Equal(t, DipperinCliCompleterNew(d), []prompt.Suggest{prompt.Suggest{Text: "SetMineGasConfig", Description: ""}, prompt.Suggest{Text: "SetMineCoinBase", Description: ""}, prompt.Suggest{Text: "StartMine", Description: ""}, prompt.Suggest{Text: "StopMine", Description: ""}})
}

func TestDipperinCliCompleterNew(t *testing.T) {
	d := prompt.Document{}

	args := strings.Split("tx ", " ")
	fmt.Println(len(args))
	b := prompt.NewBuffer()
	b.InsertText("tx ", false, true)

	d = *b.Document()

	suggest := DipperinCliCompleterNew(d)
	log.DLogger.Debug("TestDipperinCliCompleterNew", zap.Any("suggest", suggest))

	args = strings.Split("tx SendTx ", " ")
	fmt.Println(len(args))
	b = prompt.NewBuffer()
	b.InsertText("tx SendTx ", false, true)

	d = *b.Document()

	suggest = DipperinCliCompleterNew(d)
	log.DLogger.Debug("TestDipperinCliCompleterNew", zap.Any("suggest", suggest))

	args = strings.Split("tx SendTx -", " ")
	fmt.Println(len(args))
	b = prompt.NewBuffer()
	b.InsertText("tx SendTx -", false, true)

	d = *b.Document()

	suggest = DipperinCliCompleterNew(d)
	log.DLogger.Debug("TestDipperinCliCompleterNew", zap.Any("suggest", suggest))

	args = strings.Split("tx SendTransactionContract --abi ", " ")
	fmt.Println(len(args), strings.TrimLeft(args[len(args)-1], "--"))
	for _, arg := range args {
		fmt.Println(strings.TrimLeft(arg, "--"))
	}
	b = prompt.NewBuffer()
	b.InsertText("tx SendTransactionContract --abi ", false, true)

	d = *b.Document()

	suggest = DipperinCliCompleterNew(d)
	log.DLogger.Debug("TestDipperinCliCompleterNew", zap.Any("suggest", suggest))

	b = prompt.NewBuffer()
	b.InsertText("rpc.Add -A -F", false, true)

	d = *b.Document()

	suggest = DipperinCliCompleterNew(d)
	log.DLogger.Debug("TestDipperinCliCompleterNew", zap.Any("suggest", suggest))

	b = prompt.NewBuffer()
	b.InsertText("rp", false, true)

	d = *b.Document()

	suggest = DipperinCliCompleterNew(d)
	log.DLogger.Debug("TestDipperinCliCompleterNew", zap.Any("suggest", suggest))
	//assert.Equal(t, DipperinCliCompleterNew(d), []prompt.Suggest{})
}

func Test_argumentsCompleter(t *testing.T) {
	assert.Equal(t, argumentsCompleterNew([]string{"test"}), []prompt.Suggest{})
	assert.Equal(t, argumentsCompleterNew([]string{"miner", "a"}), []prompt.Suggest{})
	assert.Equal(t, argumentsCompleterNew([]string{"miner", "a", "b"}), nilSuggest)
	fmt.Println(argumentsCompleterNew([]string{"miner", "-a", "b"}))
	assert.Equal(t, argumentsCompleterNew([]string{"tx", "-p"}), []prompt.Suggest{})
}

func Test_excludeOptions(t *testing.T) {
	args := []string{"-test1", "test2", "test3"}
	assert.Equal(t, excludeOptions(args), []string{"test2", "test3"})
}

func TestCheckModuleMethodIsRight(t *testing.T) {
	assert.Equal(t, false, CheckModuleMethodIsRight("", ""))
	assert.Equal(t, true, CheckModuleMethodIsRight("", "-h"))
	assert.Equal(t, true, CheckModuleMethodIsRight("tx", "SendTx"))
	assert.Equal(t, true, CheckModuleMethodIsRight("chain", "CurrentBlock"))
	assert.Equal(t, true, CheckModuleMethodIsRight("verifier", "VerifierStatus"))
	assert.Equal(t, true, CheckModuleMethodIsRight("personal", "CurrentBalance"))
	assert.Equal(t, true, CheckModuleMethodIsRight("miner", "StartMine"))
}
