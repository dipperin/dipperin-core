package state_processor

import (
	"testing"
	"io/ioutil"
	"github.com/stretchr/testify/assert"
	"github.com/dipperin/dipperin-core/core/model"
)

func TestProcessContractCreate (t *testing.T){
	code, err := ioutil.ReadFile("/tmp/sample.wasm")
	assert.NoError(t,err)

	tx := newContractCreateTx(code)
	processor := newProcessor()

	err = processor.ProcessTx(tx)
	assert.NoError(t, err)
}

func TestProcessContractCall (t *testing.T){

}

func newContractCreateTx (code []byte) model.AbstractTransaction{
	return nil
}

func newProcessor() *AccountStateDB{
	return nil
}