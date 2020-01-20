package verifiershaltcheck

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/golang/mock/gomock"
	"testing"
)

var (
	testPubKey = []byte{4, 10, 248, 89, 242, 124, 79, 120, 34, 127, 39, 85, 139, 225, 252, 185, 104, 94, 58, 117, 134, 106, 184, 94, 254, 220, 118, 164, 33, 107, 190, 232, 137, 245, 225, 221, 184, 4, 147, 86, 29, 48, 186, 31, 244, 0, 246, 57, 222, 205, 135, 226, 106, 243, 108, 65, 17, 77, 253, 57, 141, 172, 195, 29, 35}
	testBootSign = []byte{168, 169, 140, 244, 123, 191, 221, 245, 38, 20, 198, 105, 66, 62, 102, 73, 24, 103, 205, 93, 224, 192, 206, 67, 46, 71, 109, 120, 16, 58, 10, 242, 14, 10, 41, 189, 204, 65, 128, 227, 98, 56, 241, 254, 105, 28, 76, 229, 104, 210, 130, 44, 91, 212, 12, 103, 153, 56, 50, 36, 184, 28, 51, 149, 0}
	bootNodeIndex = 0
)

// New A Proposal Config For Test
func MakeTestProposalConfig(t *testing.T, voteType model.VoteMsgType, verBootIndex int) ProposalGeneratorConfig {
	bootNodeIndex = verBootIndex
	
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	
	block := NewMockAbstractBlock(ctrl)
	block.EXPECT().Number().Return(uint64(2)).AnyTimes()
	block.EXPECT().Version().Return(uint64(0)).AnyTimes()
	block.EXPECT().Hash().Return(common.Hash{}).AnyTimes()
	block.EXPECT().StateRoot().Return(common.Hash{}).AnyTimes()
	block.EXPECT().GetRegisterRoot().Return(common.Hash{}).AnyTimes()
	block.EXPECT().GetInterLinkRoot().Return(common.Hash{}).AnyTimes()
	block.EXPECT().GetInterlinks().Return(model.InterLink{common.Hash{}}).AnyTimes()
	
	return ProposalGeneratorConfig{
		CurBlock: block,
		PubKey:   testPubKey,
		VoteType: voteType,
		ProcessStateFunc: processState,
		SignHashFunc: signHash,
	}
}

func processState(_ model.AbstractBlock, _, _ common.Hash) (stateRoot, registerRoot common.Hash, err error) {
	return common.Hash{}, common.Hash{}, nil
}

func signHash(hash []byte) (bytes []byte, e error) {
	return testBootSign, nil
}

// New A Proposal Msg For Test
func MakeTestProposalMsg(t *testing.T, verBootIndex int) (*ProposalMsg, error) {
	config := MakeTestProposalConfig(t, model.VerBootNodeVoteMessage, verBootIndex)
	return GenProposalMsg(config)
}

// New A HaltHandle For Test
func MakeTestHaltHandle(t *testing.T, verBootIndex int) *VBHaltHandler {
	config := MakeTestProposalConfig(t, model.VerBootNodeVoteMessage, verBootIndex)
	return NewHaltHandler(config)
}

func ProposeEmptyBlockForTest(t *testing.T, verBootIndex int) (*VBHaltHandler, error) {
	handle := MakeTestHaltHandle(t, verBootIndex)
	if _, err := handle.ProposeEmptyBlock(); err != nil {
		return nil, err
	}
	return handle, nil
}