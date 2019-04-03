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


package dipperin

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/core/chain-communication"
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/dipperin/dipperin-core/third-party/p2p"
	"github.com/dipperin/dipperin-core/third-party/p2p/enode"
	"path/filepath"
	"runtime"
	"sync"
	"os"
	"crypto/ecdsa"
	"github.com/dipperin/dipperin-core/third-party/crypto"
	"fmt"
	"github.com/dipperin/dipperin-core/core/cs-chain"
)

const (
	nodeInfoDirName = "app_nodes_info"

	staticNodes     = "static-nodes.json"
	trustedNodes    = "trusted-nodes.json"
)

// DefaultDataDir is the default data directory to use for the databases and other
// persistence requirements.
func DefaultDataDir() string {
	// Try to place the data folder in the user's home dir
	home := util.HomeDir()
	if home != "" {
		if runtime.GOOS == "darwin" {
			return filepath.Join(home, "Library", chain_config.AppName)
		} else if runtime.GOOS == "windows" {
			return filepath.Join(home, "AppData", "Roaming", chain_config.AppName)
		} else {
			return filepath.Join(home, "."+chain_config.AppName)
		}
	}
	// As we cannot guess a stable location, return empty and handle later
	return ""
}

// get the default node config
func DefaultNodeConf() NodeConfig {
	c := NodeConfig{Name: chain_config.AppName, DataDir: DefaultDataDir(), IPCPath: "/tmp/dipperin.ipc", HTTPHost: "127.0.0.1", HTTPPort: 7777, WSHost: "127.0.0.1", WSPort: 8888, IsUploadNodeData: 0, UploadURL: ""}
	return c
}

// get the default p2p configuration
func DefaultP2PConf() p2p.Config {
	conf := p2p.Config{NoDiscovery: false, MaxPeers: chain_communication.P2PMaxPeerCount, ListenAddr: ":60606"}
	return conf
}

func DefaultMinerP2PConf() p2p.Config {
	return p2p.Config{NoDiscovery: true, MaxPeers: chain_communication.P2PMaxPeerCount, ListenAddr: ":68080"}
}

/*// add nodes to static nodes exclude local ip
func addStaticNodeWithoutCurIp(conf *p2p.Config, n *enode.Node, localAddrs []net.Addr) {
	for _, a := range localAddrs {
		if strings.Contains(a.String(), n.IP().String()) {
			log.Debug("Is the local node configuration, not added to the static node", "local ip", a.String(), "n ip", n.IP().String())
			return
		}
	}
	log.Debug("Not a native node configuration, join a static node", "n ip", n.IP().String())
	conf.StaticNodes = append(conf.StaticNodes, n)
}

// add nodes to static nodes exclude local ip
func addStaticNodeWithoutCurNodeId(conf *p2p.Config, n *enode.Node, curNodeId string) {
	if n.ID().String() == curNodeId {
		log.Debug("Is the local node configuration, not added to the static node", "curNodeId", curNodeId, "n ip", n.IP().String())
		return
	}
	log.Debug("Not a native node configuration, join a static node", "node id", n.ID().String())
	conf.StaticNodes = append(conf.StaticNodes, n)
}

// add a static node to the conf
func AddStaticNodesToP2PConf(conf *p2p.Config, curNodeId string) {
	//for _, vn := range chain_config.VerifyNodes {
	//	addStaticNodeWithoutCurNodeId(conf, vn, curNodeId)
	//}
	//for _, mn := range chain_config.MineMasterNodes {
	//	addStaticNodeWithoutCurNodeId(conf, mn, curNodeId)
	//}
}*/

/*// read static node configuration into global
func ReadStaticNodesFromDataDir(dataDir string) {
	log.Info("read static node configuration")
	//chain_config.VerifyNodes = readNodesFromFile(filepath.Join(dataDir, chain_config.StaticVerifiersFileName))
	//chain_config.MineMasterNodes = readNodesFromFile(filepath.Join(dataDir, chain_config.StaticMinerMastersFileName))
}

// read the configured node from the file
func readNodesFromFile(confFile string) (result []*enode.Node) {
	vb, err := ioutil.ReadFile(confFile)
	if err != nil {
		log.Warn("unable to read node configuration", "err", err)
		return
	}
	var vs []string
	if err = util.ParseJsonFromBytes(vb, &vs); err != nil {
		log.Warn("unable to parse node configuration", "err", err)
		return
	}
	for _, s := range vs {
		if n, pErr := enode.ParseV4(s); pErr != nil {
			log.Warn("parsing node error", "err", err)
		} else {
			result = append(result, n)
		}
	}
	return
}*/

func loadNodeKeyFromFile(dataDir string) *ecdsa.PrivateKey {
	nodeKeyFilePath := filepath.Join(dataDir, "nodekey")

	if !util.FileExist(dataDir) {
		os.MkdirAll(dataDir, 0755)
	}
	if key, err := crypto.LoadECDSA(nodeKeyFilePath); err == nil {
		return key
	} else {
		log.Info("can't load nodekey from data dir, gen a new key", "key file path", nodeKeyFilePath, "load err", err)
	}
	key, err := crypto.GenerateKey()
	if err != nil {
		log.Error("gen node key failed", "err", err)
		return nil
	}
	if err = crypto.SaveECDSA(nodeKeyFilePath, key); err != nil {
		log.Error("save node key failed", "err", err)
	}
	return key
}

func getNodeList(path string) []*enode.Node {
	if _, err := os.Stat(path); err != nil {
		return nil
	}

	var nodelist []string
	if err := common.LoadJSON(path, &nodelist); err != nil {
		log.Error(fmt.Sprintf("Can't load node file %s: %v", path, err))
		return nil
	}

	var nodes []*enode.Node
	for _, url := range nodelist {
		if url == "" {
			continue
		}
		node, err := enode.ParseV4(url)
		if err != nil {
			log.Error(fmt.Sprintf("Node URL %s: %v\n", url, err))
			continue
		}
		nodes = append(nodes, node)
	}
	return nodes

}

//TODO: change back to the permission
func MakeVerifiersReader(fullChain cs_chain.Chain) *ChainVerifiersReader {
	return &ChainVerifiersReader{
		fullChain: fullChain,
	}
}

type ChainVerifiersReader struct {
	fullChain cs_chain.Chain
	lock sync.Mutex
}


func (verifier *ChainVerifiersReader) CurrentVerifiers() []common.Address {
	return verifier.fullChain.GetCurrVerifiers()
}

func (verifier *ChainVerifiersReader) NextVerifiers() []common.Address {
	return verifier.fullChain.GetNextVerifiers()
}

// fetch the primary node that have proposed the current block
func (verifier *ChainVerifiersReader) PrimaryNode() common.Address {
	verifiers := verifier.fullChain.GetCurrVerifiers()
	return verifiers[0]
}

// fetch the primary node that have proposed the current block. But if it's the last block of the current round, then it should return the primary node of the next round.
func (verifier *ChainVerifiersReader) GetPBFTPrimaryNode() common.Address {
	var verifiers []common.Address
	if verifier.ShouldChangeVerifier() {
		log.Info("GetPBFTPrimaryNode# The current height on the chain is the last block of the round, and verifiers of the next round should be taken.")
		verifiers = verifier.fullChain.GetNextVerifiers()
	} else {
		verifiers = verifier.fullChain.GetCurrVerifiers()
	}

	return verifiers[0]
}

func (verifier *ChainVerifiersReader) VerifiersTotalCount() int {
	verifiers := verifier.fullChain.GetCurrVerifiers()
	return len(verifiers)
}

func (verifier *ChainVerifiersReader) ShouldChangeVerifier() bool {
	currentBlock := verifier.fullChain.CurrentBlock()
	// If 10 blocks counts one round and the 9th is on the chain, then the 10th should be verified by verifiers of the next round.
	return verifier.fullChain.IsChangePoint(currentBlock, false)
}
