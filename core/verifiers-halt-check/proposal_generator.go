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

package verifiers_halt_check

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/bloom"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third-party/crypto/cs-crypto"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/dipperin/dipperin-core/third-party/log/ver_halt_check_log"
	"math/big"
	"time"
)

// generate ProposalMsg
func GenProposalMsg(config ProposalGeneratorConfig) (*ProposalMsg, error) {
	g := &ProposalGenerator{ProposalGeneratorConfig: config}
	return g.GenProposal()
}

type SignHashFunc func(hash []byte) ([]byte, error)

type ProcessFunc func(block model.AbstractBlock, preStateRoot, preRegisterRoot common.Hash) (stateRoot, registerRoot common.Hash, err error)

type ProposalGeneratorConfig struct {
	// block need
	CurBlock          model.AbstractBlock
	NewBlockSeed      common.Hash
	NewBlockProof     []byte
	LastVerifications []model.AbstractVerification
	PubKey            []byte

	// vote msg need
	SignHashFunc SignHashFunc

	//process account and register need
	ProcessStateFunc ProcessFunc

	VoteType model.VoteMsgType
}

type ProposalGenerator struct {
	ProposalGeneratorConfig

	address common.Address
}

func (g *ProposalGenerator) GenProposal() (*ProposalMsg, error) {
	emptyBlock, err := g.GenEmptyBlock()
	if err != nil {
		return nil, err
	}
	vm, err := GenVoteMsg(emptyBlock, g.SignHashFunc, g.getAddress(), g.VoteType)
	if err != nil {
		return nil, err
	}
	return &ProposalMsg{
		EmptyBlock: *emptyBlock,
		VoteMsg:    *vm,
	}, nil
}

func (g *ProposalGenerator) GenEmptyBlock() (*model.Block, error) {
	curBlock := g.CurBlock
	log.Info("GenerateEmptyBlock", "num", curBlock.Number())

	currentHeight := curBlock.Number()
	header := &model.Header{
		Version:     curBlock.Version(),
		Number:      currentHeight + 1,
		Seed:        g.NewBlockSeed,
		Proof:       g.NewBlockProof,
		MinerPubKey: g.PubKey,
		PreHash:     curBlock.Hash(),
		Diff:        common.Difficulty{},
		TimeStamp:   big.NewInt(time.Now().Add(time.Second * 3).UnixNano()),
		CoinBase:    g.getAddress(),
		Bloom:       iblt.NewBloom(model.DefaultBlockBloomConfig),
	}

	block := model.NewBlock(header, []*model.Transaction{}, g.LastVerifications)

	sateRoot, registerRoot, err := g.ProcessStateFunc(block, curBlock.StateRoot(), curBlock.GetRegisterRoot())
	if err != nil {
		return nil, err
	}

	block.SetStateRoot(sateRoot)
	block.SetRegisterRoot(registerRoot)

	// set interlink root
	linkList := model.NewInterLink(curBlock.GetInterlinks(), block)
	block.SetInterLinks(linkList)
	linkRoot := model.DeriveSha(linkList)
	block.SetInterLinkRoot(linkRoot)
	//avoid hash error
	block.RefreshHashCache()

	ver_halt_check_log.Log.Info("the GenEmptyBlock block hash is:", "hash", block.Hash().Hex())
	return block, nil
}

func (g *ProposalGenerator) getAddress() common.Address {
	if !g.address.IsEmpty() {
		return g.address
	}

	pk := cs_crypto.ToECDSAPub(g.PubKey)
	g.address = cs_crypto.GetNormalAddress(*pk)

	return g.address
}
