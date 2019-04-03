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
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNodeName(t *testing.T) {
	s, err := NodeName()
	assert.Error(t, err)
	assert.Equal(t, "", s)
}

func TestDataDir(t *testing.T) {
	_, err := DataDir()
	assert.Error(t, err)
	//assert.Equal(t, "/home/qydev/tmp/dipperin_apps/node", s)
}

func TestP2PListener(t *testing.T) {
	s, err := P2PListener()
	assert.Error(t, err)
	assert.Equal(t, "", s)
}

func TestHTTPPort(t *testing.T) {
	s, err := HTTPPort()
	assert.Error(t, err)
	assert.Equal(t, "", s)
}

func TestWSPort(t *testing.T) {
	s, err := WSPort()
	assert.Error(t, err)
	assert.Equal(t, "", s)
}

func TestWalletPassword(t *testing.T) {
	s, err := WalletPassword()
	assert.Error(t, err)
	assert.Equal(t, "", s)
}

func TestWalletPassPhrase(t *testing.T) {
	s, err := WalletPassPhrase()
	assert.Error(t, err)
	assert.Equal(t, "", s)
}

func TestWalletPath(t *testing.T) {
	s, err := WalletPath("/tmp/test_wallet_path")
	assert.Error(t, err)
	assert.Equal(t, "", s)
}

func Test_portValidate(t *testing.T) {
	err := portValidate("xx123")
	assert.Error(t, err)
	err = portValidate("8123")
	assert.NoError(t, err)
}

func Test_filepathValidate(t *testing.T) {
	err := filepathValidate("sdfx.sdf/wef")
	assert.Error(t, err)
	err = filepathValidate("/tmp/ok")
	assert.NoError(t, err)
	err = filepathValidate("")
	assert.NoError(t, err)
}

func Test_emptyValidate(t *testing.T) {
	err := emptyValidate("")
	assert.Error(t, err)
	err = emptyValidate("123")
	assert.NoError(t, err)
}
