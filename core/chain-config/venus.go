package chain_config

import (
	"github.com/dipperin/dipperin-core/common/hexutil"
	env_conf "github.com/dipperin/dipperin-core/core/chain-config/env-conf"
	cs_crypto "github.com/dipperin/dipperin-core/third-party/crypto/cs-crypto"
	"github.com/dipperin/dipperin-core/third-party/p2p/enode"
	"net"
)

const (
	venusHostIP       = "14.17.65.122"
	InitVenusBootPort = 10000
)

func venusKBoots() []*enode.Node {

	pkByte, _ := hexutil.Decode(env_conf.VenusBootNodePublicKey)
	cs_crypto.ToECDSAPub(pkByte)

	n := enode.NewV4(cs_crypto.ToECDSAPub(pkByte), net.ParseIP(venusHostIP), 30301, 30301)
	return []*enode.Node{n}
}

func NewVenusVBoots() []*enode.Node {
	config := GetChainConfig()
	vBoots := make([]*enode.Node, 0)
	for i := 0; i < config.VerifierBootNodeNumber; i++ {
		pkByte, _ := hexutil.Decode(env_conf.MercuryVerBootPublicKey[i])
		cs_crypto.ToECDSAPub(pkByte)

		n := enode.NewV4(cs_crypto.ToECDSAPub(pkByte), net.ParseIP(venusHostIP), InitVenusBootPort+(i+1)*3, InitVenusBootPort+(i+1)*3)
		vBoots = append(vBoots, n)
	}
	return vBoots
}

func venusVBoots() []*enode.Node {
	n, _ := enode.ParseV4("enode://7a035400458c476d52f49287d062445349fa3c3b5dd101392baf4f1953d47687b53d3191abfa144576e22bf979c3d0d6bae5ecac7a83aeb4c9230fc5253179fa@14.17.65.122:10000")
	return []*enode.Node{n}
}
