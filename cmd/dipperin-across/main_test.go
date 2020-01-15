package main

import (
	"github.com/dipperin/dipperin-core/cmd/dipperin-across/node_cluster"
	"github.com/dipperin/dipperin-core/cmd/dipperin-across/sidechain"
	"github.com/dipperin/dipperin-core/cmd/dipperincli/commands"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_main(t *testing.T) {
	t.Skip()
	clusterA, err := node_cluster.CreateIpcNodeCluster("local")
	assert.NoError(t, err)

	clusterB, err := node_cluster.CreateIpcNodeCluster("tps")
	assert.NoError(t, err)

	_, err = clusterA.GetAddressBalance(alice)
	assert.NoError(t, err)

	_, err = clusterA.GetAddressBalance(bob)
	assert.NoError(t, err)

	_, err = clusterB.GetAddressBalance(alice)
	assert.NoError(t, err)

	_, err = clusterB.GetAddressBalance(bob)
	assert.NoError(t, err)
}

func Test_GetSPVProof(t *testing.T) {
	t.Skip()
	clusterB, err := node_cluster.CreateIpcNodeCluster("tps")
	assert.NoError(t, err)

	amount, _ := commands.MoneyValueToCSCoin("1dip")
	_, _, err = sidechain.GetSPVProof(clusterB, bob, alice, amount)
	assert.NoError(t, err)
}
