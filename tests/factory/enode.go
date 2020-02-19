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

package factory

import (
	"encoding/hex"
	"github.com/dipperin/dipperin-core/third-party/p2p/enode"
	"github.com/dipperin/dipperin-core/third-party/p2p/enr"
	"github.com/ethereum/go-ethereum/rlp"
)

func GetEnode(str string) *enode.Node {
	//t.Skip("aaa")
	var pyRecord, _ = hex.DecodeString(str)
	var r enr.Record
	if err := rlp.DecodeBytes(pyRecord, &r); err != nil {
		panic(err)
	}
	n, err := enode.New(enode.ValidSchemes, &r)
	if err != nil {
		panic(err)
	}

	return n
}
