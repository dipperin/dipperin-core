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

package rpcinterface

//go:generate mockgen -destination=./peer_manager_mock_test.go -package=rpcinterface github.com/dipperin/dipperin-core/core/chaincommunication PeerManager

//go:generate mockgen -destination=./chain_interface_mock_test.go -package=rpcinterface github.com/dipperin/dipperin-core/core/cs-chain/chain-writer/middleware ChainInterface

//go:generate mockgen -destination=./block_mock_test.go -package=rpcinterface github.com/dipperin/dipperin-core/core/model AbstractBlock

//go:generate mockgen -destination=./node_conf_mock_test.go -package=rpcinterface github.com/dipperin/dipperin-core/core/dipperin/service NodeConf

//go:generate mockgen -destination=./peer_mock_test.go -package=rpcinterface github.com/dipperin/dipperin-core/core/chaincommunication PmAbstractPeer

//go:generate mockgen -destination=./protocol_manager_mock_test.go -package=rpcinterface github.com/dipperin/dipperin-core/core/chaincommunication AbstractPbftProtocolManager
