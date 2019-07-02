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

package address_util

import (
	"crypto/ecdsa"
	"encoding/binary"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/third-party/crypto"
)

func GenERC20Address() (common.Address, error) {
	return GenContractAddress(common.AddressTypeERC20)
}

func GenContractAddress(addressType int) (common.Address, error) {
	pk, err := crypto.GenerateKey()
	if err != nil {
		return common.Address{}, err
	}
	return PubKeyToAddress(pk.PublicKey, addressType), nil
}

func PubKeyToAddress(p ecdsa.PublicKey, addressType int) common.Address {
	pubBytes := crypto.FromECDSAPub(&p)
	var tmpTypeB [2]byte
	binary.BigEndian.PutUint16(tmpTypeB[:], uint16(addressType))
	tmpAddr := crypto.Keccak256(pubBytes[1:])[12:]
	return common.BytesToAddress(append(tmpTypeB[:], tmpAddr...))
}
