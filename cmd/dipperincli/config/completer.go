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
	"strings"
	"unicode"
)

var nilSuggest []prompt.Suggest

var commands = []prompt.Suggest{
	{Text: "rpc", Description: "rpc method"},
	{Text: "exit", Description: "exit"},
}

func DipperinCliCompleter(d prompt.Document) []prompt.Suggest {
	if d.TextBeforeCursor() == "" {
		return nilSuggest
	}

	args := strings.Split(d.TextBeforeCursor(), " ")
	w := d.GetWordBeforeCursor()
	if strings.HasPrefix(w, "-") {
		return optionCompleter(args, strings.HasPrefix(w, "--"))
	}

	for i, r := range w {
		if i == 0 {
			if unicode.IsUpper(r) {
				return callMethod(args, strings.HasPrefix(w, "--"))
			}
		}
	}

	return argumentsCompleter(excludeOptions(args))
}

func argumentsCompleter(args []string) []prompt.Suggest {
	l := len(args)

	if l <= 1 {
		return prompt.FilterHasPrefix(commands, args[0], true)
	}

	first := args[0]

	switch first {
	case "rpc", "r":
		if l == 2 {
			second := args[1]
			var subCommands []prompt.Suggest
			return prompt.FilterHasPrefix(subCommands, second, true)
		}
	}

	return nilSuggest
}

func excludeOptions(args []string) []string {
	ret := make([]string, 0, len(args))
	for i := range args {
		if !strings.HasPrefix(args[i], "-") {
			ret = append(ret, args[i])
		}
	}
	return ret
}
