package service

import (
	"testing"
	config2 "github.com/dipperin/dipperin-core/common/config"
	"math/big"
	"github.com/stretchr/testify/assert"
	"github.com/dipperin/dipperin-core/core/model"
	"fmt"
)

func TestNewOracle(t *testing.T) {
	config := Config{
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
	config := Config{
		Blocks:     20,
		Percentile: 60,
		Default:    big.NewInt(config2.DEFAULT_GAS_PRICE),
	}
	oracle := NewOracle(csChain, config)
	gasPrice, err := oracle.SuggestPrice()
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(1), gasPrice)

	// create txs
	var txs1 []*model.Transaction
	for i := 0; i < 50; i++ {
		tx := createSignedTx3(uint64(i), big.NewInt(0), big.NewInt(int64(i)))
		txs1 = append(txs1, tx)
	}

	// insert block to chain
	for i := 0; i < 10; i++ {
		insertBlockToChain(t, csChain, 1, txs1[i*5:(i+1)*5])
	}

	fmt.Println("------------------------------------------------")

	// get suggest gas price
	gasPrice, err = oracle.SuggestPrice()
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(25), gasPrice)

	// create txs
	var txs2 []*model.Transaction
	for i := 0; i < 50; i++ {
		tx := createSignedTx3(uint64(i+50), big.NewInt(0), big.NewInt(int64(50-i)))
		txs2 = append(txs2, tx)
	}

	// insert block to chain
	for i := 0; i < 10; i++ {
		insertBlockToChain(t, csChain, 1, txs2[i*5:(i+1)*5])
	}

	// get suggest gas price
	gasPrice, err = oracle.SuggestPrice()
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(26), gasPrice)

	gasPrice, err = oracle.SuggestPrice()
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(26), gasPrice)
}
