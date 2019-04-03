package state_machine

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func MakeVoteSet(height uint64) *VoteSet{
	_, v := CreateKey()
	return NewVoteSet(height, v)
}
func TestVoteSet_AddVote(t *testing.T) {
	vs := MakeVoteSet(1)

	vote := MakeNewVote(uint64(2),uint64(1),&FakeBlock{},0)
	err := vs.AddVote(vote)

	assert.Error(t,err,"vote height not match")
}