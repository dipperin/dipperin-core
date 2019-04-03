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


package utils

import (
	"os"
	"testing"

	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/stretchr/testify/assert"
)

func TestSetupGenesis(t *testing.T) {
	defer os.RemoveAll("/tmp/test_genesis")

	assert.Panics(t, func() {
		SetupGenesis("/dev/test_genesis", chain_config.GetChainConfig())
	})

	SetupGenesis("/tmp/test_genesis", chain_config.GetChainConfig())
}

func TestFatalf(t *testing.T) {
	assert.Panics(t, func() {
		Fatalf("test Fatalf")
	})
}
