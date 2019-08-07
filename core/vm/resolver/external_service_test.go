package resolver

import (
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestResolverNeedExternalService_Transfer(t *testing.T) {
	vmValue := &fakeVmContextService{}
	contract := &fakeContractService{}
	state := NewFakeStateDBService()
	service := &resolverNeedExternalService{
		contract,
		vmValue,
		state,
	}

	resp, gasLeft, err := service.Transfer(aliceAddr, big.NewInt(100))
	assert.NoError(t, err)
	assert.Equal(t, []byte(nil), resp)
	assert.Equal(t, uint64(0), gasLeft)
}
