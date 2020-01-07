package spv

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestSPVProof_Validate(t *testing.T) {
	block := model.CreateBlock(2, common.HexToHash("123"), 20)
	tx := block.GetTransactions()[10]
	proof, err := NewSPVProof(*tx, block)
	assert.NoError(t, err)

	chainID := chain_config.GetChainConfig().ChainId.Uint64()
	height := uint64(2)
	err = proof.Validate(chainID, height, model.AliceAddr, model.BobAddr, tx.Amount())
	assert.NoError(t, err)

	err = proof.Validate(uint64(1), height, model.AliceAddr, model.BobAddr, tx.Amount())
	assert.Equal(t, invalidChainID, err)

	err = proof.Validate(chainID, uint64(1), model.AliceAddr, model.BobAddr, tx.Amount())
	assert.Equal(t, invalidHeight, err)

	testAddr := common.HexToAddress("123456")
	err = proof.Validate(chainID, height, testAddr, model.BobAddr, tx.Amount())
	assert.Equal(t, invalidFrom, err)

	err = proof.Validate(chainID, height, model.AliceAddr, testAddr, tx.Amount())
	assert.Equal(t, invalidTo, err)

	err = proof.Validate(chainID, height, model.AliceAddr, model.BobAddr, big.NewInt(100))
	assert.Equal(t, invalidAmount, err)

	tx = model.CreateSignedTx(100, big.NewInt(5000))
	proof, err = NewSPVProof(*tx, block)
	assert.NoError(t, err)
	assert.NotNil(t, proof)

	err = proof.Validate(chainID, height, model.AliceAddr, model.BobAddr, tx.Amount())
	assert.Equal(t, invalidProof, err)
}
