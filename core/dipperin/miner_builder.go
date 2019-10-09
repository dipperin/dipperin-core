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
	"errors"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/core/chain-communication"
	"github.com/dipperin/dipperin-core/core/mine/mineworker"
	"github.com/dipperin/dipperin-core/third-party/p2p"
	"github.com/dipperin/dipperin-core/third-party/p2p/enode"
	"path/filepath"
)

func NewMinerNode(master string, coinbase string, minerCount int, p2pListenAddr string) (n Node, err error) {
	if coinbase == "" || minerCount <= 0 {
		err = errors.New("coinbase or miner count not right")
		return
	}
	var masterNode *enode.Node
	if masterNode, err = enode.ParseV4(master); err != nil {
		return nil, errors.New("parse master node failed:" + err.Error())
	}
	// init p2p
	p2pConf := DefaultMinerP2PConf()
	p2pConf.StaticNodes = []*enode.Node{masterNode}
	minerKeyPath := filepath.Join(util.HomeDir(), "dipperin_miner", coinbase)
	p2pConf.PrivateKey = loadNodeKeyFromFile(minerKeyPath)
	p2pConf.ListenAddr = p2pListenAddr
	p2pServer := &p2p.Server{Config: p2pConf}

	// init mine protocols
	_, remoteConn := mineworker.MakeRemoteWorker(common.HexToAddress(coinbase), minerCount)

	//p2pServer.Protocols
	minePm := chain_communication.NewMineProtocolManager(remoteConn)
	p2pServer.Protocols = append(p2pServer.Protocols, minePm.GetProtocol())

	// init miner
	n = NewCsNode([]NodeService{p2pServer}, NodeConfig{})
	return
}
