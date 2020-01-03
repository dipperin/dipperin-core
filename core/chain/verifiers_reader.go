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
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/log"
	"sync"
)

func MakeVerifiersReader(fullChain Chain) *VerifiersReader {
	return &VerifiersReader{
		fullChain: fullChain,
	}
}

type VerifiersReader struct {
	fullChain Chain
	lock      sync.Mutex
}

func (verifier *VerifiersReader) CurrentVerifiers() []common.Address {
	return verifier.fullChain.GetCurrVerifiers()
}

func (verifier *VerifiersReader) NextVerifiers() []common.Address {
	return verifier.fullChain.GetNextVerifiers()
}

// Get the main node of the current block, not considering the situation of 9
func (verifier *VerifiersReader) PrimaryNode() common.Address {
	verifiers := verifier.fullChain.GetCurrVerifiers()
	return verifiers[0]
}

// To get the PBFT master node, use the next considering the situation of 9.
func (verifier *VerifiersReader) GetPBFTPrimaryNode() common.Address {
	var verifiers []common.Address
	if verifier.ShouldChangeVerifier() {
		//The current height on the chain is the last block of the round, and the next round of verifiers should be taken
		log.DLogger.Info("[GetPBFTPrimaryNode]:The current height on the chain is the last block of the round, and the next round of verifiers should be taken")
		verifiers = verifier.fullChain.GetNextVerifiers()
	} else {
		verifiers = verifier.fullChain.GetCurrVerifiers()
	}
	return verifiers[0]
}

func (verifier *VerifiersReader) VerifiersTotalCount() int {
	verifiers := verifier.fullChain.GetCurrVerifiers()
	return len(verifiers)
}

func (verifier *VerifiersReader) ShouldChangeVerifier() bool {
	currentBlock := verifier.fullChain.CurrentBlock()
	//If there are 10 blocks in one round then when the 9th is on the chain, the 10th should be verified by the next round of verifiers.
	return verifier.fullChain.IsChangePoint(currentBlock, false)
}
