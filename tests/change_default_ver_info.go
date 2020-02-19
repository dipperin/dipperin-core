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

package tests

import (
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/chain"
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/dipperin/dipperin-core/third-party/crypto"
	"github.com/dipperin/dipperin-core/third-party/p2p/enode"
)

func ChangeVerBootNodeAddress() ([]Account, error) {
	config := chain_config.GetChainConfig()
	if config.VerifierBootNodeNumber != 4 {
		panic("config.VerifierBootNodeNumber != 4")
	}
	accounts := []Account{
		{Pk: crypto.HexToECDSAErrPanic("1e00aa89565549d616d43c4e71b61d46a963fdb69489093a57cacf06836ecd91")},
		{Pk: crypto.HexToECDSAErrPanic("2e00ab89565549d616d43c4e71b61d46a963fdb69489093a57cacf06836ecd92")},
		{Pk: crypto.HexToECDSAErrPanic("3e00ac89565549d616d43c4e71b61d46a963fdb69489093a57cacf06836ecd93")},
		{Pk: crypto.HexToECDSAErrPanic("4e00ad89565549d616d43c4e71b61d46a963fdb69489093a57cacf06836ecd94")},
	}

	for i := 0; i < 4; i++ {
		chain_config.VerBootNodeAddress[i] = accounts[i].Address()
	}

	n, _ := enode.ParseV4(fmt.Sprintf("enode://b832f4f2fe19dbc5604766bbb268a6d0f7ce9ce381b034b262a92f0ad8283a1b5fa058dea5269b66fbb2014a24fa7198c6dc2d8c9cbac7a348258fc20702561f@%v:%v", "127.0.0.1", 10003))
	chain_config.VerifierBootNodes = append(chain_config.VerifierBootNodes, n)
	n, _ = enode.ParseV4(fmt.Sprintf("enode://199cc6526cb63866dfa5dc81aed9952f2002b677560b6f3dc2a6a34a5576216f0ca25711c5b4268444fdef5fee4a01a669af90fd5b6049b2a5272b39c466b2ac@%v:%v", "127.0.0.1", 10006))
	chain_config.VerifierBootNodes = append(chain_config.VerifierBootNodes, n)
	n, _ = enode.ParseV4(fmt.Sprintf("enode://71112a581231af08a63d5a9079ea8dd690efd992f2cfbf98ad43697345de564441406133247d19c754c98051c64909c40db15094770a881a373ca1ff2f20bea2@%v:%v", "127.0.0.1", 10009))
	chain_config.VerifierBootNodes = append(chain_config.VerifierBootNodes, n)
	n, _ = enode.ParseV4(fmt.Sprintf("enode://07f3fdca9a07b048ea7d0cb642f69004e4fa5dd390888a9bb3e9fc382697c3634280cc8d327703b872d3711462da4aca96ee805069510375e7be2aded3dc5ad6@%v:%v", "127.0.0.1", 10012))
	chain_config.VerifierBootNodes = append(chain_config.VerifierBootNodes, n)

	return accounts, nil
}

func ChangeVerifierAddress(accounts []Account) ([]Account, error) {
	changeVerifiers := []Account{
		{Pk: crypto.HexToECDSAErrPanic("fe00ee89565549d616d43c4e71b61d46a963fdb69489093a57cacf06836ecd91")},
		{Pk: crypto.HexToECDSAErrPanic("fe00ee89565549d616d43c4e71b61d46a963fdb69489093a57cacf06836ecd92")},
		{Pk: crypto.HexToECDSAErrPanic("fe00ee89565549d616d43c4e71b61d46a963fdb69489093a57cacf06836ecd93")},
		{Pk: crypto.HexToECDSAErrPanic("fe00ee89565549d616d43c4e71b61d46a963fdb69489093a57cacf06836ecd94")},
		{Pk: crypto.HexToECDSAErrPanic("fe00ee89565549d616d43c4e71b61d46a963fdb69489093a57cacf06836ecd95")},
		{Pk: crypto.HexToECDSAErrPanic("fe00ee89565549d616d43c4e71b61d46a963fdb69489093a57cacf06836ecd96")},
		{Pk: crypto.HexToECDSAErrPanic("fe00ee89565549d616d43c4e71b61d46a963fdb69489093a57cacf06836ecd97")},
		{Pk: crypto.HexToECDSAErrPanic("fe00ee89565549d616d43c4e71b61d46a963fdb69489093a57cacf06836ecd98")},
		{Pk: crypto.HexToECDSAErrPanic("fe00ee89565549d616d43c4e71b61d46a963fdb69489093a57cacf06836ecd99")},
		{Pk: crypto.HexToECDSAErrPanic("fe00ee89565549d616d43c4e71b61d46a963fdb69489093a57cacf06836ecd9a")},
		{Pk: crypto.HexToECDSAErrPanic("fe00ee89565549d616d43c4e71b61d46a963fdb69489093a57cacf06836ecd9b")},
		{Pk: crypto.HexToECDSAErrPanic("fe00ee89565549d616d43c4e71b61d46a963fdb69489093a57cacf06836ecd9c")},
		{Pk: crypto.HexToECDSAErrPanic("fe00ee89565549d616d43c4e71b61d46a963fdb69489093a57cacf06836ecd9d")},
		{Pk: crypto.HexToECDSAErrPanic("fe00ee89565549d616d43c4e71b61d46a963fdb69489093a57cacf06836ecd9e")},
		{Pk: crypto.HexToECDSAErrPanic("fe00ee89565549d616d43c4e71b61d46a963fdb69489093a57cacf06836ecd9f")},
		{Pk: crypto.HexToECDSAErrPanic("fe00ee89565549d616d43c4e71b61d46a963fdb69489093a57cacf06836ecd1a")},
		{Pk: crypto.HexToECDSAErrPanic("fe00ee89565549d616d43c4e71b61d46a963fdb69489093a57cacf06836ecd2a")},
		{Pk: crypto.HexToECDSAErrPanic("fe00ee89565549d616d43c4e71b61d46a963fdb69489093a57cacf06836ecd3a")},
		{Pk: crypto.HexToECDSAErrPanic("fe00ee89565549d616d43c4e71b61d46a963fdb69489093a57cacf06836ecd4a")},
		{Pk: crypto.HexToECDSAErrPanic("fe00ee89565549d616d43c4e71b61d46a963fdb69489093a57cacf06836ecd5a")},
		{Pk: crypto.HexToECDSAErrPanic("fe00ee89565549d616d43c4e71b61d46a963fdb69489093a57cacf06836ecd6a")},
		{Pk: crypto.HexToECDSAErrPanic("fe00ee89565549d616d43c4e71b61d46a963fdb69489093a57cacf06836ecd7a")},
	}

	if accounts != nil {
		for i := 0; i < len(accounts); i++ {
			changeVerifiers[i] = accounts[i]
		}
	}

	chain.VerifierAddress = []common.Address{}
	for _, v := range changeVerifiers {
		chain.VerifierAddress = append(chain.VerifierAddress, v.Address())
	}
	return changeVerifiers, nil
}
