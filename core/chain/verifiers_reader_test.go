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

package chain

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestMakeVerifiersReader(t *testing.T) {
	vReader := MakeVerifiersReader(fakeChain{})
	curVerifiers := vReader.CurrentVerifiers()
	assert.Equal(t, aliceAddr, curVerifiers[0])

	nextVerifiers := vReader.NextVerifiers()
	assert.Equal(t, bobAddr, nextVerifiers[0])

	pNode := vReader.PrimaryNode()
	assert.Equal(t, aliceAddr, pNode)

	count := vReader.VerifiersTotalCount()
	assert.Equal(t, 1, count)

	pNode = vReader.GetPBFTPrimaryNode()
	assert.Equal(t, aliceAddr, pNode)

	vReader = MakeVerifiersReader(fakeChain{createBlock(20)})
	pNode = vReader.GetPBFTPrimaryNode()
	assert.Equal(t, bobAddr, pNode)
}
