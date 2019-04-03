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


package config

import (
	"github.com/c-bata/go-prompt"
	"strings"
)

func optionCompleter(args []string, long bool) []prompt.Suggest {
	l := len(args)
	if l <= 1 {
		if long {
			return prompt.FilterHasPrefix(optionHelp, "--", false)
		}
		return optionHelp
	}

	var suggests []prompt.Suggest
	commandArgs := excludeOptions(args)
	switch commandArgs[0] {
	case "rpc":
		suggests = rpcFlags
	}

	if long {
		return prompt.FilterContains(
			prompt.FilterHasPrefix(suggests, "--", false),
			strings.TrimLeft(args[l-1], "--"),
			true,
		)
	}
	return prompt.FilterContains(suggests, strings.TrimLeft(args[l-1], "-"), true)
}

var optionHelp = []prompt.Suggest{
	{Text: "-h"},
	{Text: "--help"},
}

var rpcFlags = []prompt.Suggest{
	{Text: "-m", Description: "operation（this method must match rpc server's method）"},
	{Text: "-p", Description: "parameters"},
}

func callMethod(args []string, long bool) []prompt.Suggest {
	l := len(args)
	if l <= 2 {
		if long {
			return prompt.FilterHasPrefix(optionHelp, "--", false)
		}
		return optionHelp
	}

	var suggests []prompt.Suggest
	suggests = methodFlags

	if long {
		return prompt.FilterContains(
			prompt.FilterHasPrefix(suggests, "--", false),
			strings.TrimLeft(args[l-1], "--"),
			true,
		)
	}
	return prompt.FilterContains(suggests, strings.TrimLeft(args[l-1], "-"), true)
}

var methodFlags = []prompt.Suggest{
	{Text: "AddAccount", Description: ""},
	{Text: "AddPeer", Description: ""},
	{Text: "AnnounceERC20", Description: ""},
	{Text: "CloseWallet", Description: ""},
	{Text: "CurrentBalance", Description: ""},
	{Text: "CurrentBlock", Description: ""},
	{Text: "CurrentStake", Description: ""},
	{Text: "CurrentReputation", Description: ""},
	{Text: "ERC20Allowance", Description: ""},
	{Text: "ERC20Approve", Description: ""},
	{Text: "ERC20Balance", Description: ""},
	{Text: "ERC20GetInfo", Description: ""},
	//{Text: "ERC20TokenDecimals", Description: ""},
	//{Text: "ERC20TokenName", Description: ""},
	//{Text: "ERC20TokenSymbol", Description: ""},
	//{Text: "ERC20TotalSupply", Description: ""},
	{Text: "ERC20Transfer", Description: ""},
	{Text: "ERC20TransferFrom", Description: ""},
	{Text: "EstablishWallet", Description: ""},
	{Text: "GetAddressNonceFromWallet", Description: ""},
	{Text: "GetBlockByHash", Description: ""},
	{Text: "GetBlockByNumber", Description: ""},
	{Text: "GetCurVerifiers", Description: ""},
	{Text: "GetDefaultAccountBalance", Description: ""},
	{Text: "GetDefaultAccountStake", Description: ""},
	{Text: "GetGenesis", Description: ""},
	{Text: "GetNextVerifiers", Description: ""},
	{Text: "GetTransactionNonce", Description: ""},
	{Text: "GetVerifiersBySlot", Description: ""},
	{Text: "ListWallet", Description: ""},
	{Text: "ListWalletAccount", Description: ""},
	{Text: "OpenWallet", Description: ""},
	{Text: "Peers", Description: ""},
	{Text: "RestoreWallet", Description: ""},
	{Text: "SendCancelTransaction", Description: ""},
	{Text: "SendCancelTx", Description: ""},
	{Text: "SendUnStakeTransaction", Description: ""},
	{Text: "SendUnStakeTx", Description: ""},
	{Text: "SendRegisterTransaction", Description: ""},
	{Text: "SendRegisterTx", Description: ""},
	{Text: "SendTransaction", Description: ""},
	{Text: "SendTx", Description: ""},
	{Text: "SetExchangeRate", Description: ""},
	{Text: "SetMineCoinBase", Description: ""},
	{Text: "SetBftSigner", Description: ""},
	{Text: "StartMine", Description: ""},
	{Text: "StopMine", Description: ""},
	{Text: "Transaction", Description: ""},
	{Text: "TransferEDIPToDIP", Description: ""},
	{Text: "VerifierStatus", Description: ""},
	{Text: "GetBlockDiffVerifierInfo", Description: ""},
	{Text: "CheckVerifierType", Description: ""},
}
