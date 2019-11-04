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
	"github.com/dipperin/dipperin-core/third-party/log"
	"strings"
)

/*func optionCompleter(args []string, long bool) []prompt.Suggest {
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
		suggests = txPromptFlags
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
*/
func optionCompleterNew(args []string, long bool) []prompt.Suggest {
	l := len(args)
	if l <= 1 {
		if long {
			return prompt.FilterHasPrefix(optionHelp, "--", false)
		}
		return optionHelp
	}

	var suggests []prompt.Suggest
	commandArgs := excludeOptions(args)
	//fmt.Println("optionCompleterNew", "commandArgs", commandArgs)
	log.Debug("optionCompleterNew", "commandArgs", commandArgs)

	if len(args) == 2 {
		suggests = getSuggestFromModuleName(commandArgs[0])
	} else {
		switch commandArgs[0] {
		case "tx":
			suggests = txPromptFlags
		case "chain", "verifier", "personal", "miner":
			suggests = commonFlags
		}
	}

	log.Debug("optionCompleterNew", "suggests", suggests)
	defer log.Debug("optionCompleterNew", "suggests  defer", suggests)
	arg := args[l-1]
	for i := l - 1; arg == ""; i-- {
		arg = args[i]
	}
	if long {
		return prompt.FilterContains(
			suggests,
			strings.TrimLeft(args[l-1], "--"),
			true,
		)
	}
	log.Debug("optionCompleterNew", "suggests1", suggests)
	return prompt.FilterContains(suggests, strings.TrimLeft(args[l-1], "--"), true)
}

var optionHelp = []prompt.Suggest{
	{Text: "-h"},
	{Text: "--help"},
}

var commonFlags = []prompt.Suggest{
	{Text: "-p", Description: "parameters"},
}

var txPromptFlags = []prompt.Suggest{
	{Text: "-p", Description: "parameters"},
	{Text: "--abi", Description: "abi path"},
	{Text: "--wasm", Description: "wasm path"},
	{Text: "--input", Description: "contract parameters"},
	{Text: "--is-create", Description: "create contract or call, default is call"},
	{Text: "--func-name", Description: "the function to call"},
}

/*func callMethod(args []string, long bool) []prompt.Suggest {
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
*/
var personalMethods = []prompt.Suggest{
	// personal
	{Text: "CurrentBalance", Description: ""},
	{Text: "CloseWallet", Description: ""},
	{Text: "AddAccount", Description: ""},
	{Text: "CurrentStake", Description: ""},
	{Text: "CurrentReputation", Description: ""},
	{Text: "EstablishWallet", Description: ""},
	{Text: "GetAddressNonceFromWallet", Description: ""},
	{Text: "GetDefaultAccountBalance", Description: ""},
	{Text: "GetDefaultAccountStake", Description: ""},
	{Text: "GetTransactionNonce", Description: ""},
	{Text: "ListWallet", Description: "list wallet"},
	{Text: "ListWalletAccount", Description: ""},
	{Text: "OpenWallet", Description: ""},
	{Text: "RestoreWallet", Description: ""},
	{Text: "SetBftSigner", Description: ""},
}

var minerMethods = []prompt.Suggest{
	{Text: "SetMineGasConfig", Description: ""},
	{Text: "SetMineCoinBase", Description: ""},
	{Text: "StartMine", Description: ""},
	{Text: "StopMine", Description: ""},
}

var txMethods = []prompt.Suggest{
	{Text: "AnnounceERC20", Description: ""},
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
	{Text: "SendCancelTransaction", Description: ""},
	{Text: "SendCancelTx", Description: ""},
	{Text: "SendUnStakeTransaction", Description: ""},
	{Text: "SendUnStakeTx", Description: ""},
	{Text: "SendRegisterTransaction", Description: ""},
	{Text: "SendRegisterTx", Description: ""},
	{Text: "SendTransaction", Description: ""},
	{Text: "SendTransactionContract", Description: ""},
	{Text: "SendTx", Description: ""},
	{Text: "TransferEDIPToDIP", Description: ""},
	{Text: "GetContractAddressByTxHash", Description: ""},
	{Text: "CallContract", Description: ""},
	{Text: "EstimateGas", Description: ""},
	{Text: "Transaction", Description: ""},
}

var chainMethods = []prompt.Suggest{
	{Text: "AddPeer", Description: ""},
	{Text: "CurrentBlock", Description: ""},
	{Text: "GetBlockByHash", Description: ""},
	{Text: "GetBlockByNumber", Description: ""},
	{Text: "GetSlotByNumber", Description: ""},
	{Text: "GetGenesis", Description: ""},
	{Text: "Peers", Description: ""},
	{Text: "SetExchangeRate", Description: ""},
	{Text: "GetLogs", Description: ""},
	{Text: "GetReceiptByTxHash", Description: ""},
	{Text: "GetReceiptsByBlockNum", Description: ""},
	{Text: "GetTxActualFee", Description: ""},
	{Text: "SuggestGasPrice", Description: ""},
}

var verifierMethods = []prompt.Suggest{
	// verifier
	{Text: "GetCurVerifiers", Description: ""},
	{Text: "GetNextVerifiers", Description: ""},
	{Text: "GetVerifiersBySlot", Description: ""},
	{Text: "VerifierStatus", Description: ""},
	{Text: "GetBlockDiffVerifierInfo", Description: ""},
	{Text: "CheckVerifierType", Description: ""},
}

/*var methodFlags = []prompt.Suggest{
	// personal
	{Text: "CurrentBalance", Description: ""},
	{Text: "CloseWallet", Description: ""},
	{Text: "AddAccount", Description: ""},
	{Text: "CurrentStake", Description: ""},
	{Text: "CurrentReputation", Description: ""},
	{Text: "EstablishWallet", Description: ""},
	{Text: "GetAddressNonceFromWallet", Description: ""},
	{Text: "GetDefaultAccountBalance", Description: ""},
	{Text: "GetDefaultAccountStake", Description: ""},
	{Text: "GetTransactionNonce", Description: ""},
	{Text: "ListWallet", Description: "list wallet"},
	{Text: "ListWalletAccount", Description: ""},
	{Text: "OpenWallet", Description: ""},
	{Text: "RestoreWallet", Description: ""},

	// miner
	{Text: "SetMineGasConfig", Description: ""},
	{Text: "SetMineCoinBase", Description: ""},
	{Text: "StartMine", Description: ""},
	{Text: "StopMine", Description: ""},

	// chain
	{Text: "AddPeer", Description: ""},
	{Text: "CurrentBlock", Description: ""},
	{Text: "GetBlockByHash", Description: ""},
	{Text: "GetBlockByNumber", Description: ""},
	{Text: "GetGenesis", Description: ""},
	{Text: "Peers", Description: ""},

	// tx
	{Text: "AnnounceERC20", Description: ""},
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
	{Text: "SendCancelTransaction", Description: ""},
	{Text: "SendCancelTx", Description: ""},
	{Text: "SendUnStakeTransaction", Description: ""},
	{Text: "SendUnStakeTx", Description: ""},
	{Text: "SendRegisterTransaction", Description: ""},
	{Text: "SendRegisterTx", Description: ""},
	{Text: "SendTransaction", Description: ""},
	{Text: "SendTransactionContract", Description: ""},
	{Text: "SendTx", Description: ""},
	{Text: "Transaction", Description: ""},
	{Text: "GetContractAddressByTxHash", Description: ""},
	{Text: "GetConvertReceiptByTxHash", Description: ""},
	{Text: "GetReceiptByTxHash", Description: ""},
	{Text: "GetReceiptsByBlockNum", Description: ""},
	{Text: "TransferEDIPToDIP", Description: ""},
	{Text: "SetExchangeRate", Description: ""},
	{Text: "EstimateGas", Description: ""},

	// verifier
	{Text: "GetCurVerifiers", Description: ""},
	{Text: "GetNextVerifiers", Description: ""},
	{Text: "GetVerifiersBySlot", Description: ""},
	{Text: "SetBftSigner", Description: ""},
	{Text: "VerifierStatus", Description: ""},
	{Text: "GetBlockDiffVerifierInfo", Description: ""},
	{Text: "CheckVerifierType", Description: ""},
}*/
