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

package verifiershaltcheck

import (
	"github.com/stretchr/testify/assert"
	"testing"
	
	"github.com/dipperin/dipperin-core/core/model"
)

func TestGenProposalMsg(t *testing.T) {
	_, err := MakeTestProposalMsg(t, 0)
	assert.NoError(t, err)
}

func TestProposalGenerator_GenProposal(t *testing.T) {
	cfg := MakeTestProposalConfig(t, model.VerBootNodeVoteMessage, 0)
	pg := ProposalGenerator{ProposalGeneratorConfig: cfg}
	_, err := pg.GenProposal()
	assert.NoError(t, err)
}

func TestProposalGenerator_GenEmptyBlock(t *testing.T) {
	cfg := MakeTestProposalConfig(t, model.VerBootNodeVoteMessage, 0)
	pg := ProposalGenerator{ProposalGeneratorConfig: cfg}
	_, err := pg.GenEmptyBlock()
	assert.NoError(t, err)
}


func TestProposalGenerator_getAddress(t *testing.T) {
	cfg := MakeTestProposalConfig(t, model.VerBootNodeVoteMessage, 0)
	pg := ProposalGenerator{ProposalGeneratorConfig: cfg}
	addr := pg.getAddress()
	assert.NotNil(t, addr)
}
