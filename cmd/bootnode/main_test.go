// Copyright 2015 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

// bootnode runs a bootstrap node for the Ethereum Discovery Protocol.
package main

import (
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func Test_main(t *testing.T) {
	assert.Panics(t, func() {
		main()
	})

	*genKey = "bt_test_k"
	main()
	*genKey = ""

	*nodeKeyHex = "xxy"
	*nodeKeyFile = "xx"
	assert.Panics(t, func() {
		main()
	})

	*nodeKeyFile = ""
	assert.Panics(t, func() {
		main()
	})

	*nodeKeyFile = "xx"
	*nodeKeyHex = ""
	assert.Panics(t, func() {
		main()
	})

	*nodeKeyFile = ""
	*nodeKeyHex = "81910c1adab446e4ff8624913f4c051ea6654de1e807504bd8be4be658af5545"
	*writeAddr = true
	main()

	*writeAddr = false
	assert.NoError(t, os.Setenv(chain_config.BootEnvTagName, "mercury"))
	*netrestrict = "123"
	assert.Panics(t, func() {
		main()
	})

	*netrestrict = ""
	assert.NoError(t, os.Setenv(chain_config.BootEnvTagName, "test"))
	*listenAddr = "xxx"
	assert.Panics(t, func() {
		main()
	})

	*listenAddr = ":3731"
	*natdesc = "any"
	go main()
	time.Sleep(50 * time.Millisecond)

	*listenAddr = ":3732"
	*natdesc = ""
	go main()
	time.Sleep(50 * time.Millisecond)

	*listenAddr = ":3733"
	*runv5 = true
	go main()
	time.Sleep(50 * time.Millisecond)

	assert.Panics(t, func() {
		main()
	})

	time.Sleep(500 * time.Millisecond)
}
