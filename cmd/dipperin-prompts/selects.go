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


package dipperin_prompts

import (
	"github.com/manifoldco/promptui"
	"fmt"
)

type nodeType struct {
	Name    string
	Value 	int
}

func NodeType() (int, error) {
	nodeTypes := []nodeType{
		{Name: "Normal", Value: 0},
		{Name: "MineMaster", Value: 1},
		{Name: "Verifier", Value: 2},
		//{Name: "Verifier Boot", Value: 3},
	}

	prompt := promptui.Select{
		Label:     "Node Type",
		Items:     nodeTypes,
		Templates: TplNodeType,
	}

	i, _, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return nodeTypes[0].Value, err
	}

	return nodeTypes[i].Value, err
}

