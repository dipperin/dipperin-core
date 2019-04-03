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


package rpc_interface

import "github.com/dipperin/dipperin-core/third-party/p2p"

//go:generate mockgen -destination=./p2p_api_mock_test.go -package=rpc_interface github.com/caiqingfeng/dipperin-core/core/rpc-interface P2PAPI
type P2PAPI interface {
	AddPeer(url string) error
	RemovePeer(url string) error
	AddTrustedPeer(url string)  error
	RemoveTrustedPeer(url string)  error
	Peers() ([]*p2p.PeerInfo, error)
	CsPmInfo() (*p2p.CsPmPeerInfo, error)
}

type DipperinP2PApi struct {
	service P2PAPI
}

func (api *DipperinP2PApi) AddPeer(url string) error {
	return api.service.AddPeer(url)
}

func (api *DipperinP2PApi) RemovePeer(url string) error {
	return api.service.RemovePeer(url)
}

func (api *DipperinP2PApi) AddTrustedPeer(url string) error {
	return api.service.AddTrustedPeer(url)
}

func (api *DipperinP2PApi) RemoveTrustedPeer(url string) error {
	return api.service.RemoveTrustedPeer(url)
}

func (api *DipperinP2PApi) Peers() ([]*p2p.PeerInfo, error) {
	return api.service.Peers()
}

func (api *DipperinP2PApi) CsPmInfo() (*p2p.CsPmPeerInfo, error) {
	return api.service.CsPmInfo()
}

