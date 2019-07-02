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

package tests

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/bloom"
	"github.com/dipperin/dipperin-core/core/csbft/model"
	"github.com/dipperin/dipperin-core/core/csbft/state-machine"
	model2 "github.com/dipperin/dipperin-core/core/model"
	"math/big"
	"sync/atomic"
)

func NewBftCluster(verifiers []Account) *BftCluster {
	count := len(verifiers)
	addresses := (Accounts)(verifiers).GetAddresses()

	var states []*state_machine.BftState
	for i := 0; i < count; i++ {
		s := &state_machine.BftState{Height: 0, BlockPoolNotEmpty: false}
		s.OnNewHeight(1, 1, addresses)
		s.OnBlockPoolNotEmpty()
		states = append(states, s)
	}
	return &BftCluster{BftStates: states, Verifiers: verifiers}
}

type BftCluster struct {
	BftStates []*state_machine.BftState
	Verifiers []Account
}

func (bc *BftCluster) NewRoundMsg(x int, h, round uint64) (r []*model.NewRoundMsg) {
	for i := 0; i < x; i++ {
		v := bc.Verifiers[i]
		m := model.NewRoundMsgWithSign(h, round, v.SignHash, v.Address())
		if m == nil {
			panic("can't create round msg for: " + v.Address().Hex())
		}
		r = append(r, m)
	}
	return
}

func (bc *BftCluster) NewProposal(x int, round uint64, block model2.AbstractBlock) (p []*model.Proposal) {
	for i := 0; i < x; i++ {
		v := bc.Verifiers[i]
		m := model.NewProposalWithSign(block.Number(), round, block.Hash(), v.SignHash, v.Address())
		if m == nil {
			panic("can't create round msg for: " + v.Address().Hex())
		}
		p = append(p, m)
	}
	return
}

func (bc *BftCluster) NewVote(x int, round uint64, vt model2.VoteMsgType, block model2.AbstractBlock) (p []*model2.VoteMsg) {
	for i := 0; i < x; i++ {
		v := bc.Verifiers[i]
		m, err := model2.NewVoteMsgWithSign(block.Number(), round, block.Hash(), vt, v.SignHash, v.Address())
		if err != nil {
			panic(err)
		}
		p = append(p, m)
	}
	return
}

func (bc *BftCluster) StatesIter(cb func(*state_machine.BftState)) {
	for _, bs := range bc.BftStates {
		cb(bs)
	}
}

type FakeBlockForBft struct {
	Num   uint64
	PHash common.Hash

	hash atomic.Value `rlp:"-"`
}

func (fb *FakeBlockForBft) Version() uint64 {
	panic("implement me")
}

func (fb *FakeBlockForBft) Number() uint64 {
	return fb.Num
}

func (fb *FakeBlockForBft) IsSpecial() bool {
	panic("implement me")
}

func (fb *FakeBlockForBft) Difficulty() common.Difficulty {
	panic("implement me")
}

func (fb *FakeBlockForBft) PreHash() common.Hash {
	return fb.PHash
}

func (fb *FakeBlockForBft) Seed() common.Hash {
	panic("implement me")
}

func (fb *FakeBlockForBft) RefreshHashCache() common.Hash {
	panic("implement me")
}

func (fb *FakeBlockForBft) Hash() common.Hash {
	if x := fb.hash.Load(); x != nil {
		return x.(common.Hash)
	}

	h := common.RlpHashKeccak256(fb)
	fb.hash.Store(h)
	return h
}

func (fb *FakeBlockForBft) EncodeRlpToBytes() ([]byte, error) {
	panic("implement me")
}

func (fb *FakeBlockForBft) TxIterator(cb func(int, model2.AbstractTransaction) error) error {
	panic("implement me")
}

func (fb *FakeBlockForBft) TxRoot() common.Hash {
	panic("implement me")
}

func (fb *FakeBlockForBft) Timestamp() *big.Int {
	panic("implement me")
}

func (fb *FakeBlockForBft) Nonce() common.BlockNonce {
	panic("implement me")
}

func (fb *FakeBlockForBft) StateRoot() common.Hash {
	panic("implement me")
}

func (fb *FakeBlockForBft) SetStateRoot(root common.Hash) {
	panic("implement me")
}

func (fb *FakeBlockForBft) GetRegisterRoot() common.Hash {
	panic("implement me")
}

func (fb *FakeBlockForBft) SetRegisterRoot(root common.Hash) {
	panic("implement me")
}

func (fb *FakeBlockForBft) FormatForRpc() interface{} {
	panic("implement me")
}

func (fb *FakeBlockForBft) SetNonce(nonce common.BlockNonce) {
	panic("implement me")
}

func (fb *FakeBlockForBft) CoinBaseAddress() common.Address {
	panic("implement me")
}

func (fb *FakeBlockForBft) GetTransactionFees() *big.Int {
	panic("implement me")
}

func (fb *FakeBlockForBft) CoinBase() *big.Int {
	panic("implement me")
}

func (fb *FakeBlockForBft) GetTransactions() []*model2.Transaction {
	panic("implement me")
}

func (fb *FakeBlockForBft) GetInterlinks() model2.InterLink {
	panic("implement me")
}

func (fb *FakeBlockForBft) SetInterLinkRoot(root common.Hash) {
	panic("implement me")
}

func (fb *FakeBlockForBft) GetInterLinkRoot() (root common.Hash) {
	panic("implement me")
}

func (fb *FakeBlockForBft) SetInterLinks(inter model2.InterLink) {
	panic("implement me")
}

func (fb *FakeBlockForBft) GetAbsTransactions() []model2.AbstractTransaction {
	panic("implement me")
}

func (fb *FakeBlockForBft) GetBloom() iblt.Bloom {
	panic("implement me")
}

func (fb *FakeBlockForBft) Header() model2.AbstractHeader {
	panic("implement me")
}

func (fb *FakeBlockForBft) Body() model2.AbstractBody {
	panic("implement me")
}

func (fb *FakeBlockForBft) TxCount() int {
	panic("implement me")
}

func (fb *FakeBlockForBft) GetEiBloomBlockData(reqEstimator *iblt.HybridEstimator) *model2.BloomBlockData {
	panic("implement me")
}

func (fb *FakeBlockForBft) GetBlockTxsBloom() *iblt.Bloom {
	panic("implement me")
}

func (fb *FakeBlockForBft) VerificationRoot() common.Hash {
	panic("implement me")
}

func (fb *FakeBlockForBft) SetVerifications(vs []model2.AbstractVerification) {
	panic("implement me")
}

func (fb *FakeBlockForBft) VersIterator(func(int, model2.AbstractVerification, model2.AbstractBlock) error) error {
	panic("implement me")
}

func (fb *FakeBlockForBft) GetVerifications() []model2.AbstractVerification {
	panic("implement me")
}
