package main

import (
	"github.com/dipperin/dipperin-core/cmd/dipperin-across/node_cluster"
	"github.com/dipperin/dipperin-core/cmd/dipperin-across/oracle"
	"github.com/dipperin/dipperin-core/cmd/dipperin-across/sidechain"
	"github.com/dipperin/dipperin-core/cmd/dipperincli/commands"
)

var (
	miner  = "default_m0"
	alice  = "default_v0"
	bob    = "default_v1"
	chainB = uint64(1602)
)

/*
1. miner在链A上部署Oracle智能合约
2. alice在链A上部署Sidechain智能合约
3. bob在链B上发送给alice的转账交易
4. bob生成转账交易的SPVProof
5. miner调用Oracle合约的SetHeader方法（B链区块头存入A链）
6. bob在调用Sidechain合约的Transfer方法（验证SPVProof，实现交易原子性）
*/

func main() {
	clusterA, err := node_cluster.CreateIpcNodeCluster("local")
	if err != nil {
		panic(err)
	}

	clusterB, err := node_cluster.CreateIpcNodeCluster("tps")
	if err != nil {
		panic(err)
	}

	oracleAddr, err := oracle.InitOracleContract(clusterA, miner)
	if err != nil {
		panic(err)
	}

	//go oracle.SubscribeChainBlock(clusterA, clusterB, miner, oracleAddr)

	amount, _ := commands.MoneyValueToCSCoin("10dip")
	sidechainAddr, err := sidechain.InitSidechainContract(clusterA, alice, bob, amount, chainB)
	if err != nil {
		panic(err)
	}

	spvProof, height, err := sidechain.GetSPVProof(clusterB, bob, alice, amount)
	if err != nil {
		panic(err)
	}

	err = oracle.CallOracleContract(clusterA, clusterB, miner, oracleAddr, height)
	if err != nil {
		panic(err)
	}

	err = sidechain.CallSidechainContract(clusterA, bob, sidechainAddr, oracleAddr, spvProof)
	if err != nil {
		panic(err)
	}
}
