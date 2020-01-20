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
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMakeVerifiersReader(t *testing.T) {
	vReader := MakeVerifiersReader(fakeChain{})

	type result struct {
		data interface{}
	}
	testCases := []struct {
		name   string
		given  func() interface{}
		expect result
	}{
		{
			name:"CurrentVerifiers",
			given: func() interface{} {
				curVerifiers := vReader.CurrentVerifiers()
				return curVerifiers[0]
			},
			expect:result{aliceAddr},
		},
		{
			name:"NextVerifiers",
			given: func() interface{} {
				nextVerifiers := vReader.NextVerifiers()
				return nextVerifiers[0]
			},
			expect:result{bobAddr},
		},
		{
			name:"PrimaryNode",
			given: func() interface{} {
				pNode := vReader.PrimaryNode()
				return pNode
			},
			expect:result{aliceAddr},
		},
		{
			name:"VerifiersTotalCount",
			given: func() interface{} {
				count := vReader.VerifiersTotalCount()
				return count
			},
			expect:result{1},
		},
		{
			name:"GetPBFTPrimaryNode",
			given: func() interface{} {
				pNode := vReader.GetPBFTPrimaryNode()
				return pNode
			},
			expect:result{aliceAddr},
		},
		{
			name:"GetPBFTPrimaryNode",
			given: func() interface{} {
				vReader = MakeVerifiersReader(fakeChain{createBlock(20)})
				pNode := vReader.GetPBFTPrimaryNode()
				return pNode
			},
			expect:result{bobAddr},
		},
	}

	for _,tc:=range testCases{
		data:=tc.given()
		assert.Equal(t,tc.expect.data,data)
	}
}

