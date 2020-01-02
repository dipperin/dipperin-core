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

package state_machine

import (
	"errors"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/log"
	model2 "github.com/dipperin/dipperin-core/core/csbft/model"
	"github.com/dipperin/dipperin-core/core/model"
	"go.uber.org/zap"
	"sync"
)

/*

Record the latest round value of all nodes at the current height

1. Uniquely record the latest round value of all nodes
2. It can judge whether the number of votes for the round of the current node reaches at least 2/3, in which case the external starts a new round of pbft process.
3. In case where the current node is backward, it can judge the value of the most recent round value where the consensus can be reached, and inform the outside that it should move to that round value.

*/
type NewRoundSet struct {
	Height        uint64
	maj32         uint64
	halfUp        uint64
	RoundMessages map[uint64]map[common.Address]*model2.NewRoundMsg
	verifiers     []common.Address
	lock          sync.Mutex
}

func NewNRoundSet(height uint64, vers []common.Address) *NewRoundSet {
	return &NewRoundSet{
		RoundMessages: make(map[uint64]map[common.Address]*model2.NewRoundMsg),
		Height:        height,
		maj32:         0,
		halfUp:        0,
		verifiers:     vers}
}

func (rs *NewRoundSet) MissingAtRound(round uint64) *model2.BitArray {
	rs.lock.Lock()
	defer rs.lock.Unlock()
	byteArray := model2.NewBitArray(len(rs.verifiers))
	if rs.RoundMessages[round] == nil {
		return byteArray
	}
	for index, add := range rs.verifiers {
		if rs.RoundMessages[round][add] != nil {
			byteArray.SetIndex(index, true)
		}
	}
	return byteArray
}

// check have enough msg at round
func (rs *NewRoundSet) EnoughAtRound(round uint64) bool {
	rs.lock.Lock()
	defer rs.lock.Unlock()

	if rs.RoundMessages[round] == nil {
		return false
	}
	vLen := len(rs.verifiers)
	count := len(rs.RoundMessages[round])
	need := vLen * 2 / 3
	// should be more than need
	log.DLogger.Info("NewRoundSet#EnoughAtRound   check round msg enough", zap.Uint64("round", round), zap.Uint64("height", rs.Height), zap.Int("count", count), zap.Int("need", need+1), zap.Int("vLen", vLen))
	// more than 2/3
	if count > need {
		return true
	}
	return false
}

func (rs *NewRoundSet) enoughAtRound(round uint64) bool {
	rs.lock.Lock()
	defer rs.lock.Unlock()

	if rs.RoundMessages[round] == nil {
		return false
	}
	vLen := len(rs.verifiers)
	count := len(rs.RoundMessages[round])
	need := vLen * 2 / 3
	// should be more than need
	log.DLogger.Info("NewRoundSet#EnoughAtRound   check round msg enough", zap.Uint64("round", round), zap.Uint64("height", rs.Height), zap.Int("count", count), zap.Int("need", need+1), zap.Int("vLen", vLen))
	// more than 2/3
	if count > need {
		return true
	}
	return false
}

func (rs *NewRoundSet) shouldCatchUpTo() uint64 {
	return rs.maj32
}

// check height valid outside
func (rs *NewRoundSet) Add(nr *model2.NewRoundMsg) error {
	rs.lock.Lock()
	defer rs.lock.Unlock()

	if nr.Height != rs.Height {
		log.DLogger.Info("the Height is:", zap.Uint64("nr", nr.Height), zap.Uint64("rs", rs.Height))
		return errors.New("new round msg height not match")
	}

	if rs.RoundMessages[nr.Round] == nil {
		rs.RoundMessages[nr.Round] = make(map[common.Address]*model2.NewRoundMsg)
	}

	if err := nr.Valid(); err != nil {
		return err
	}

	if !rs.isCurrentVerifier(nr.Witness.Address) {
		log.DLogger.Info("NewRoundSet#Add Witness.Address is not current verifier address", zap.Any("address", nr.Witness.Address))
		return errors.New("new round msg not from verifier")
	}

	// only accept higher round
	if rs.RoundMessages[nr.Round][nr.Witness.Address] == nil {
		rs.RoundMessages[nr.Round][nr.Witness.Address] = nr
	}

	if rs.hasMaj32(nr.Round) {
		if rs.maj32 < nr.Round {
			rs.maj32 = nr.Round
		}
	}

	return nil
}

func (rs *NewRoundSet) isCurrentVerifier(vAddr common.Address) bool {
	for _, curV := range rs.verifiers {
		if curV.IsEqual(vAddr) {
			return true
		}
	}
	return false
}

// for the same round only a valid proposal is received, so the round key will support multiple rounds of information at the same height.
type ProposalSet struct {
	// key is round
	Proposals map[uint64]*model2.Proposal
}

func NewProposalSet() *ProposalSet {
	return &ProposalSet{Proposals: make(map[uint64]*model2.Proposal)}
}

// check have proposal for round
func (ps *ProposalSet) Have(round uint64) bool {
	return ps.Proposals[round] != nil
}

// check valid before Add outside
func (ps *ProposalSet) Add(p *model2.Proposal) {
	ps.Proposals[p.Round] = p
}

func (ps *ProposalSet) GetProposal(round uint64) *model2.Proposal {
	return ps.Proposals[round]
}

type BlockSet struct {
	// key is round
	Blocks map[uint64]model.AbstractBlock
}

func NewBlockSet() *BlockSet {
	return &BlockSet{Blocks: make(map[uint64]model.AbstractBlock)}
}

func (bs *BlockSet) GetBlock(round uint64) model.AbstractBlock {
	return bs.Blocks[round]
}

// ProposalSet will verify whether there is already a proposal for a round, so here there is no need to check
func (bs *BlockSet) AddBlock(block model.AbstractBlock, round uint64) {
	bs.Blocks[round] = block
}

func (bs *BlockSet) GetBlockByHash(h common.Hash) model.AbstractBlock {
	for _, b := range bs.Blocks {
		if b.Hash().IsEqual(h) {
			return b
		}
	}
	return nil
}

type VoteSet struct {
	// key is round
	Height    uint64
	Votes     map[uint64]map[common.Address]*model.VoteMsg
	blockVote map[uint64]map[common.Hash]int
	verifiers []common.Address
	lock      sync.Mutex
}

func NewVoteSet(height uint64, vers []common.Address) *VoteSet {
	return &VoteSet{
		Height:    height,
		Votes:     make(map[uint64]map[common.Address]*model.VoteMsg),
		blockVote: make(map[uint64]map[common.Hash]int),
		verifiers: vers,
	}
}

// get verifications for broadcast
func (vs *VoteSet) FinalVerifications(round uint64) []model.AbstractVerification {
	votes := vs.roundVotes(round)
	hash := vs.VotesEnough(round)
	if hash.IsEqual(common.Hash{}) {
		return nil
	}

	var resultVotes []model.AbstractVerification
	for _, tmpV := range votes {
		if tmpV.GetBlockId().IsEqual(hash) {
			resultVotes = append(resultVotes, tmpV)
		}
	}
	return resultVotes
}

// check votes enough for a round
func (vs *VoteSet) VotesEnough(round uint64) common.Hash {
	vs.lock.Lock()
	defer vs.lock.Unlock()

	blockVotes := vs.roundBlockVotes(round)
	verifier := vs.verifiers
	for index, count := range blockVotes {
		if count > len(verifier)*2/3 {
			log.DLogger.Info("the votes is enough", zap.String("hash", index.Hex()))
			return index
		}
	}
	return common.Hash{}
}

// get votes for round
func (vs *VoteSet) roundVotes(round uint64) map[common.Address]*model.VoteMsg {
	vm := vs.Votes[round]
	if vm == nil {
		vs.Votes[round] = make(map[common.Address]*model.VoteMsg)
	}
	return vm
}

func (vs *VoteSet) roundBlockVotes(round uint64) map[common.Hash]int {
	vm := vs.blockVote[round]
	if vm == nil {
		vs.blockVote[round] = make(map[common.Hash]int)
	}
	return vm
}

// check a vote is valid
func (vs *VoteSet) validVote(v *model.VoteMsg) error {
	// check have wit
	if v.Witness == nil {
		return errors.New("invalid vote, witness is nil")
	}

	// check is cur verifier
	if !vs.isCurrentVerifier(v.Witness.Address) {
		return errors.New("vote addr is not current verifier, addr: " + v.Witness.Address.Hex())
	}
	// check sign valid
	if err := v.Witness.Valid(v.Hash().Bytes()); err != nil {
		return err
	}
	// check is already have
	votes := vs.roundVotes(v.Round)
	if votes[v.Witness.Address] != nil {
		return errors.New("already have vote from addr: " + v.Witness.Address.Hex())
	}
	return nil
}

// check is current verifier
func (vs *VoteSet) isCurrentVerifier(vAddr common.Address) bool {
	for _, curV := range vs.verifiers {
		if curV.IsEqual(vAddr) {
			return true
		}
	}
	return false
}

// Add a valid vote
func (vs *VoteSet) AddVote(v *model.VoteMsg) error {
	vs.lock.Lock()
	defer vs.lock.Unlock()

	if v.Height != vs.Height {
		log.DLogger.Debug("[AddVote] vote not valid", zap.String("height", "not valid"))
		return errors.New("vote height not match")
	}
	if err := vs.validVote(v); err != nil {
		log.DLogger.Debug("[AddVote] vote not valid", zap.Error(err))
		return err
	}
	//vs.Votes[v.Round] = make(map[common.Address]*VoteMsg)
	vs.roundVotes(v.Round)[v.Witness.Address] = v
	if vs.roundBlockVotes(v.Round)[v.BlockID] == 0 {
		vs.roundBlockVotes(v.Round)[v.BlockID] = 1
	} else {
		vs.roundBlockVotes(v.Round)[v.BlockID]++
	}
	log.DLogger.Debug("[AddVote] success add vote")
	return nil
}

func (rs NewRoundSet) hasMaj32(round uint64) bool {
	if rs.RoundMessages[round] == nil {
		return false
	}

	num := len(rs.RoundMessages[round])

	if num > int(len(rs.verifiers))*2/3 {
		return true
	}

	return false
}

func (rs NewRoundSet) hasHalfUp(round uint64) bool {
	if rs.RoundMessages[round] == nil {
		return false
	}
	num := len(rs.RoundMessages[round])
	bakRound := round
	result := false
	for {
		bakRound++
		if bak := rs.RoundMessages[bakRound]; bak != nil {
			if len(bak) > num && len(bak) > int(len(rs.verifiers))/2 {
				rs.halfUp = bakRound
				result = true
			}
		} else {
			break
		}
	}
	if num > int(len(rs.verifiers)/2) && rs.halfUp < round {
		rs.halfUp = round
		result = true
	}
	return result
}

/*

In the case where the round of the current node is backward, it is possible to determine the value of the round value that may be reached recently, and inform the outside that it should move to the round value.
For example: 4 nodes round are 1, 2, 3, 4 respectively, then because round can only rise and can not fall, for a single node it will finally get 3, 3, 3, 4, because by calling this function 1 will get 3, 2 will also get 3, 3 and 4 will not change after the call.

*/
