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

package chain_writer

import "github.com/dipperin/dipperin-core/core/cs-chain/chain-writer/middleware"

type ChainWriter interface {
	SaveBlock() error
}

type AbstractChainWriterFactory interface {
	NewWriter(context interface{}) ChainWriter
	SetChain(chain middleware.ChainInterface)
}

//go:generate mockgen -destination=./chain_interface_mock_test.go -package=chain_writer github.com/dipperin/dipperin-core/core/cs-chain/chain-writer/middleware ChainInterface

//go:generate mockgen -destination=./block_mock_test.go -package=chain_writer github.com/dipperin/dipperin-core/core/model AbstractBlock
