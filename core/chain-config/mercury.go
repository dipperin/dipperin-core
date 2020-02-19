// Copyright 2019, Keychain Foundation Ltd.
// This file is part of the dipperin-core library.
//
// The dipperin-core library is free software: you can redistribute
// it and/or modify it under the terms of the GNU Lesser General Public License
// as published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// The dipperin-core library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package chain_config

import (
	"crypto/ecdsa"
	"encoding/hex"
	"github.com/dipperin/dipperin-core/common/hexutil"
	"github.com/dipperin/dipperin-core/core/chain-config/env-conf"
	"github.com/dipperin/dipperin-core/third-party/crypto"
	"github.com/dipperin/dipperin-core/third-party/crypto/cs-crypto"
	"github.com/dipperin/dipperin-core/third-party/p2p/enode"
	"net"
)

const (
	mercuryHostIP = "14.17.65.122"
	InitVBootPort = 10000
)

func mercuryKBoots() []*enode.Node {

	pkByte, _ := hexutil.Decode(env_conf.MercuryBootNodePublicKey)
	cs_crypto.ToECDSAPub(pkByte)

	n := enode.NewV4(cs_crypto.ToECDSAPub(pkByte), net.ParseIP(mercuryHostIP), 30301, 30301)
	return []*enode.Node{n}
}

func NewMercuryVBoots() []*enode.Node {
	config := GetChainConfig()
	vBoots := make([]*enode.Node, 0)
	for i := 0; i < config.VerifierBootNodeNumber; i++ {
		pkByte, _ := hexutil.Decode(env_conf.MercuryVerBootPublicKey[i])
		cs_crypto.ToECDSAPub(pkByte)

		n := enode.NewV4(cs_crypto.ToECDSAPub(pkByte), net.ParseIP(mercuryHostIP), InitVBootPort+(i+1)*3, InitVBootPort+(i+1)*3)
		vBoots = append(vBoots, n)
	}
	return vBoots
}

func mercuryVBoots() []*enode.Node {
	n, _ := enode.ParseV4("enode://7a035400458c476d52f49287d062445349fa3c3b5dd101392baf4f1953d47687b53d3191abfa144576e22bf979c3d0d6bae5ecac7a83aeb4c9230fc5253179fa@14.17.65.122:10000")
	return []*enode.Node{n}
}

// Unable to get the port number, so not configured temporarily
//func mercuryVBoots() []*enode.Node {
//}

func pkStrToPk(keyStr string) *ecdsa.PrivateKey {
	key, err := hex.DecodeString(keyStr)
	if err != nil {
		panic(err)
	}
	pk, err := crypto.ToECDSA(key)
	if err != nil {
		panic(err)
	}
	return pk
}
