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

const (
	active = "▸ {{ .Name | cyan | bold }}"
	inactive = "{{ .Name | cyan }}"
	label = `{{ "✔" | green }} %s: {{ .Name | cyan }}`
	Selected = `{{ "✔" | green }} {{ "%s" | cyan }}: {{ .Name }}`
)

var TplNodeType = &promptui.SelectTemplates{
	Label:    label,
	Active:   active,
	Inactive: inactive,
	Selected: fmt.Sprintf(Selected, "Node Type"),
}

var PromptTemplate = &promptui.PromptTemplates{
	Prompt: "{{ . }}:",
	Valid:   `{{ "✔" | green }} {{ . }}: `,
	Invalid: `{{ "✗" | red }} {{ . }}: `,
	Success: `{{ "✔" | green }} {{ . | cyan}}: `,
}