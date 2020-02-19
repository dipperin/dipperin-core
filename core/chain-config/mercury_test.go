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

package chain_config

import (
	"encoding/hex"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewMercuryVBoots(t *testing.T) {
	node := NewMercuryVBoots()
	assert.Equal(t, config.VerifierBootNodeNumber, len(node))
}

func TestMercuryVBoots(t *testing.T) {
	node := mercuryVBoots()
	fmt.Println(node)

	node = venusVBoots()
	fmt.Println(node)
}

func TestPkStrToPk(t *testing.T) {
	pk := pkStrToPk("fe10ee89565549d616d43c4e71b61d46a963fdb69489093a57cacf06836ecd91")
	fmt.Println(pk)

	assert.Panics(t, func() {
		pkStrToPk("123")
	})

	input := hex.EncodeToString([]byte{123})
	assert.Panics(t, func() {
		pkStrToPk(input)
	})
}

func TestMercuryKBoots(t *testing.T) {
	node := mercuryKBoots()
	fmt.Println(node)
}
