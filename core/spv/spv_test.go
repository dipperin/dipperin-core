package spv

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/tests/factory"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSPVProof_Validate(t *testing.T) {
	block := model.CreateBlock(2, common.HexToHash("123"), 20)
	proof, err := NewSPVProof(*block.GetTransactions()[10], block)
	assert.NoError(t, err)

	err = proof.Validate()
	assert.NoError(t, err)

	tx, _ := factory.CreateTestTx()
	proof, err = NewSPVProof(*tx, block)
	assert.NoError(t, err)
	assert.NotNil(t, proof)

	err = proof.Validate()
	assert.Equal(t, invalidProof, err)
}
