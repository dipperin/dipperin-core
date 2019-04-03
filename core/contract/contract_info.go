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


package contract

// contract infomation for user
type ContractInfo struct {
	// contract description
	Description string `json:"description"`
	// initilazing parameters
	InitArgs []*ContractArg `json:"init_args"`
	// contract methods
	Methods []*ContractMethod `json:"methods"`
}

// contract methods
type ContractMethod struct {
	// method name
	Name string `json:"name"`
	// method description
	Description string `json:"description"`
	// parameters
	Args []*ContractArg `json:"args"`
	// method return
	Return *ContractArg `json:"return"`
	// whether need triggering transaction
	TxMethod bool `json:"tx_method"`
}

// contract parameter
type ContractArg struct {
	// parameter's name
	Name string `json:"name"`
	// description
	Description string `json:"description"`
	// type
	ArgType string `json:"arg_type"`
}
