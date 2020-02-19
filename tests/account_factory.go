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

import "github.com/dipperin/dipperin-core/third-party/crypto"

var AccFactory = &AccountFactory{}

var defaultAccounts = []Account{
	{Pk: crypto.HexToECDSAErrPanic("fe10ee89565549d616d43c4e71b61d46a963fdb69489093a57cacf06836ecd91")},
	{Pk: crypto.HexToECDSAErrPanic("fe20ee89565549d616d43c4e71b61d46a963fdb69489093a57cacf06836ecd91")},
	{Pk: crypto.HexToECDSAErrPanic("fe30ee89565549d616d43c4e71b61d46a963fdb69489093a57cacf06836ecd91")},
	{Pk: crypto.HexToECDSAErrPanic("fe40ee89565549d616d43c4e71b61d46a963fdb69489093a57cacf06836ecd91")},
	{Pk: crypto.HexToECDSAErrPanic("fe50ee89565549d616d43c4e71b61d46a963fdb69489093a57cacf06836ecd91")},
	{Pk: crypto.HexToECDSAErrPanic("fe60ee89565549d616d43c4e71b61d46a963fdb69489093a57cacf06836ecd91")},
	{Pk: crypto.HexToECDSAErrPanic("fe70ee89565549d616d43c4e71b61d46a963fdb69489093a57cacf06836ecd91")},
	{Pk: crypto.HexToECDSAErrPanic("fe80ee89565549d616d43c4e71b61d46a963fdb69489093a57cacf06836ecd91")},
	{Pk: crypto.HexToECDSAErrPanic("fe90ee89565549d616d43c4e71b61d46a963fdb69489093a57cacf06836ecd91")},
	{Pk: crypto.HexToECDSAErrPanic("fea0ee89565549d616d43c4e71b61d46a963fdb69489093a57cacf06836ecd91")},
	{Pk: crypto.HexToECDSAErrPanic("feb0ee89565549d616d43c4e71b61d46a963fdb69489093a57cacf06836ecd91")},
}

// get or gen account
type AccountFactory struct{}

// gen account
func (acc *AccountFactory) GenAccount() Account {
	pk, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}
	return Account{Pk: pk}
}

// gen x accounts
func (acc *AccountFactory) GenAccounts(x int) (r []Account) {
	for i := 0; i < x; i++ {
		pk, err := crypto.GenerateKey()
		if err != nil {
			panic(err)
		}
		r = append(r, Account{Pk: pk})
	}
	return
}

// get default account(have certain address)
func (acc *AccountFactory) GetAccount(index int) Account {
	return defaultAccounts[index]
}
