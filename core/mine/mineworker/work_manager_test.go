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

package mineworker

import (
	"testing"
)

func TestWorkManager_SubmitWork(t *testing.T) {
	/*master := ""
	var masterNode *enode.Node
	if masterNode, err := enode.ParseV4(master); err != nil {
		panic("fail")
	}
	remote, connect := MakeRemoteWorker(common.Address{},1)
	// init p2p
	p2pConf := dipperin.DefaultMinerP2PConf()
	p2pConf.StaticNodes = []*enode.Node{ masterNode }
	minerKeyPath := filepath.Join(util.HomeDir(), "dipperin_miner", coinbase)
	p2pConf.PrivateKey = loadNodeKeyFromFile(minerKeyPath)
	p2pConf.ListenAddr = p2pListenAddr
	p2pServer := &p2p.Server{Config: p2pConf}
	minePm := chain_communication.NewMineProtocolManager(connect)
	p2pServer.Protocols = append(p2pServer.Protocols, minePm.GetProtocol())

	// init miner
	n = NewFullNode(NodeConfig{}, []NodeService{ p2pServer })*/
}
