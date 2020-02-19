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

package peer_spec

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/chain-communication"
	"github.com/dipperin/dipperin-core/third-party/p2p"
	"net"
)

func PeerBuilder() chain_communication.PmAbstractPeer {
	return &FakePeer{}
}

type FakePeer struct {
	TestMsg uint64
}

func (p *FakePeer) NodeName() string {
	panic("implement me")
}

func (p *FakePeer) NodeType() uint64 {
	panic("implement me")
}

func (p *FakePeer) SendMsg(msgCode uint64, msg interface{}) error {
	p.TestMsg = msgCode
	return nil
}

func (p *FakePeer) ID() string {
	return "123"
}

func (p *FakePeer) ReadMsg() (p2p.Msg, error) {
	return p2p.Msg{}, nil
}

func (p *FakePeer) GetHead() (common.Hash, uint64) {
	panic("implement me")
}

func (p *FakePeer) SetHead(head common.Hash, height uint64) {
	panic("implement me")
}

func (p *FakePeer) GetPeerRawUrl() string {
	panic("implement me")
}

func (p *FakePeer) DisconnectPeer() {
	panic("implement me")
}

func (p *FakePeer) RemoteVerifierAddress() (addr common.Address) {
	panic("implement me")
}

func (p *FakePeer) RemoteAddress() net.Addr {
	panic("implement me")
}

func (p *FakePeer) SetRemoteVerifierAddress(addr common.Address) {
	panic("implement me")
}

func (p *FakePeer) SetNodeName(name string) {
	panic("implement me")
}

func (p *FakePeer) SetNodeType(nt uint64) {
	panic("implement me")
}

func (p *FakePeer) SetPeerRawUrl(rawUrl string) {
	panic("implement me")
}

func (p *FakePeer) SetNotRunning() {
	panic("implement me")
}

func (p *FakePeer) IsRunning() bool {
	panic("implement me")
}

func (p *FakePeer) GetCsPeerInfo() *p2p.CsPeerInfo {
	panic("implement me")
}
