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

package softwallet

import (
	"github.com/stretchr/testify/assert"
	"testing"
)


func TestClearSensitiveData(t *testing.T) {
	testKey, err := NewMaster(testSeed, &DipperinChainCfg)
	assert.NoError(t, err)

	ClearSensitiveData(testKey)
	assert.Equal(t, ExtendedKey{}, *testKey)

	testWalletInfo := NewHdWalletInfo()

	testWalletInfo.Seed = []byte{0x00, 0x01, 0x02, 0x03}
	ClearSensitiveData(&testWalletInfo)
	assert.Equal(t, (*WalletInfo)(nil), testWalletInfo)
}
