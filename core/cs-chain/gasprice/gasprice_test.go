package gasprice

import (
	config2 "github.com/dipperin/dipperin-core/common/config"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestNewOracle(t *testing.T) {
	config := GasPriceConfig{
		Blocks:     -20,
		Percentile: -60,
		Default:    big.NewInt(config2.DEFAULT_GAS_PRICE),
	}
	oracle := NewOracle(nil, config)
	assert.Equal(t, 1, oracle.checkBlocks)
	assert.Equal(t, 5, oracle.maxBlocks)
	assert.Equal(t, 0, oracle.maxEmpty)
	assert.Equal(t, 0, oracle.percentile)

	config.Percentile = 200
	oracle = NewOracle(nil, config)
	assert.Equal(t, 100, oracle.percentile)
}

func TestOracle_SuggestPrice(t *testing.T) {
	csChain := createCsChain(nil)
	config := GasPriceConfig{
		Blocks:     20,
		Percentile: 60,
		Default:    big.NewInt(config2.DEFAULT_GAS_PRICE),
	}

	insertBlockToChain(t, csChain, 5, nil)
	oracle := NewOracle(csChain, config)
	gasPrice, err := oracle.SuggestPrice()
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(1), gasPrice)

	// insert 10 blocks with txs
	var txs1 []*model.Transaction
	for i := 0; i < 50; i++ {
		tx := createSignedTx3(uint64(i), big.NewInt(0), big.NewInt(int64(i)))
		txs1 = append(txs1, tx)
	}
	for i := 0; i < 10; i++ {
		insertBlockToChain(t, csChain, 1, txs1[i*5:(i+1)*5])
	}
	gasPrice, err = oracle.SuggestPrice()
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(25), gasPrice)

	// insert 15 blocks without txs
	insertBlockToChain(t, csChain, 15, nil)
	gasPrice, err = oracle.SuggestPrice()
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(25), gasPrice)

	// insert 10 blocks with txs
	var txs2 []*model.Transaction
	for i := 0; i < 50; i++ {
		tx := createSignedTx3(uint64(i+50), big.NewInt(0), big.NewInt(int64(100-i)))
		txs2 = append(txs2, tx)
	}
	for i := 0; i < 10; i++ {
		insertBlockToChain(t, csChain, 1, txs2[i*5:(i+1)*5])
	}
	gasPrice, err = oracle.SuggestPrice()
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(76), gasPrice)

	// get price again
	gasPrice, err = oracle.SuggestPrice()
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(76), gasPrice)
}
