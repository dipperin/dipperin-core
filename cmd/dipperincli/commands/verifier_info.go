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

package commands

import (
	"github.com/dipperin/dipperin-core/cmd/dipperincli/config"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/consts"
	"github.com/dipperin/dipperin-core/core/accounts"
	"github.com/dipperin/dipperin-core/core/rpc-interface"
	"sync"
	"time"
)

var (
	defaultAccountStake     = "0" + consts.CoinDIPName
	electionMap             sync.Map
	trackingAccounts        []accounts.Account
	logElectionTxTickerTime = 30 * time.Second
)

func loadDefaultAccountStake() {
	var resp rpc_interface.CurBalanceResp

	if err := client.Call(&resp, getDipperinRpcMethodByName("CurrentStake"), defaultAccount); err != nil {
		l.Error("call get current deposit error", "err", err)
		return
	}

	stake, err := CSCoinToMoneyValue(resp.Balance)
	if err == nil {
		defaultAccountStake = stake + consts.CoinDIPName
	}
}

func PrintCommandsModuleName() {
	var moduleName string
	for _, c := range config.Commands {
		moduleName = moduleName + c.Text + ","
	}
	if len(config.Commands) > 0 {
		moduleName = moduleName[:len(moduleName)-1]
	}
	l.Info("you can use the base command to interactive with the node :" + moduleName)
}

func PrintDefaultAccountStake() {
	l.Info("address current stake is:", "address", defaultAccount, "stake", defaultAccountStake)
}

// tag sending campaigns tx
//func markElectionTx(txHash common.Hash, address common.Address) {
//	electionMap.Store(txHash, address)
//}

//func AsyncCheckElectionTx() {
//	timer := time.NewTicker(5 * time.Second)
//
//	for {
//		select {
//		case <-timer.C:
//
//			if checkSync() {
//				return
//			}
//
//			electionMap.Range(func(key, value interface{}) bool {
//				handleElectionTx(key.(common.Hash), value.(common.Address))
//				return true
//			})
//		}
//	}
//}

func AsyncLogElectionTx() *time.Ticker {
	loadRegistedAccounts()
	timer := time.NewTicker(logElectionTxTickerTime)
	go func() {
		for {
			select {
			case <-timer.C:
				if !checkSync() {
					logElection()
				}
			}
		}
	}()
	return timer
}

func loadRegistedAccounts() {
	var resp []accounts.WalletIdentifier
	if err := client.Call(&resp, getDipperinRpcMethodByName("ListWallet")); err != nil {
		l.Error("Call ListWallet", "err", err)
		return
	}
	var respA []accounts.Account

	if err := client.Call(&respA, getDipperinRpcMethodByName("ListWalletAccount"), resp[0]); err != nil {
		l.Error("Call ListWallet", "err", err)
		return
	}

	for i := range respA {
		var resp rpc_interface.VerifierStatus

		if err := client.Call(&resp, getDipperinRpcMethodByName("VerifierStatus"), respA[i].Address.Hex()); err != nil {
			l.Error("call verifier status error", "err", err)
			return
		}

		if resp.Status == VerifierStatusRegistered {
			trackingAccounts = append(trackingAccounts, respA[i])
		}
	}

}

func logElection() {
	if len(trackingAccounts) == 0 {
		return
	}

	//var config chain_config.ChainConfig
	//if err := client.Call(&config, getDipperinRpcMethodByName("GetChainConfig")); err != nil {
	//	l.Error("get chain config", "err", err)
	//	return
	//}

	var respBlock rpc_interface.BlockResp
	if err := client.Call(&respBlock, getDipperinRpcMethodByName("CurrentBlock")); err != nil {
		l.Error("get current Block error", "err", err)
		return
	}

	var respSlot uint64
	if err := client.Call(&respSlot, getDipperinRpcMethodByName("GetSlotByNum"), respBlock.Header.Number); err != nil {
		l.Error("get current Block error", "err", err)
		return
	}

	var resp []common.Address
	if err := client.Call(&resp, getDipperinRpcMethodByName("GetVerifiersBySlot"), respSlot); err != nil {
		l.Error("get verifiers by slot error", "err", err)
	}

	for i := range trackingAccounts {
		isV := isVerifier(trackingAccounts[i].Address, resp)
		var resp rpc_interface.VerifierStatus
		if err := client.Call(&resp, getDipperinRpcMethodByName("VerifierStatus"), trackingAccounts[i].Address.Hex()); err != nil {
			l.Error("call verifier status error", "err", err)
			continue
		}

		balance, err := CSCoinToMoneyValue(resp.Balance)
		if err != nil {
			l.Error("The address has no balance, balance = 0 DIP")
			balance = "0 DIP"
		}

		l.Info("[Verifier Tracking]", "current height", respBlock.Header.Number, "slot", respSlot, trackingAccounts[i].Address.String()+" is verifier", isV)

		if resp.Status == VerifierStatusNoRegistered || resp.Status == VerifiedStatusUnstaked {

			l.Info("[Verifier Tracking]", "Verifier status", resp.Status, "balance", balance+" DIP")
			continue
		}

		stake, err := CSCoinToMoneyValue(resp.Stake)

		if err != nil {
			l.Error("The address has no stake, stake = 0 DIP")
		}

		l.Info("[Verifier Tracking]", "Verifier status", resp.Status, "balance", balance+" DIP", "stake", stake+" DIP", "reputation", resp.Reputation)

	}

}

func isVerifier(addr common.Address, verifiers []common.Address) bool {
	for i := range verifiers {
		if addr.IsEqual(verifiers[i]) {
			return true
		}
	}
	return false
}

func addTrackingAccount(adds common.Address) {
	acc := accounts.Account{Address: adds}
	for i := range trackingAccounts {
		if trackingAccounts[i].Address.IsEqual(acc.Address) {
			return
		}
	}
	trackingAccounts = append(trackingAccounts, acc)
}

func removeTrackingAccount(adds common.Address) {
	acc := accounts.Account{Address: adds}
	var newTrackingAccounts []accounts.Account
	for i := range trackingAccounts {
		if !trackingAccounts[i].Address.IsEqual(acc.Address) {
			newTrackingAccounts = append(newTrackingAccounts, trackingAccounts[i])
		}
	}
	trackingAccounts = newTrackingAccounts
}

//func handleElectionTx(txHash common.Hash, address common.Address) {
//	//l.Info("check election tx hash", "hash", txHash)
//
//	//var config chain_config.ChainConfig
//	//if err := client.Call(&config, getDipperinRpcMethodByName("GetChainConfig")); err != nil {
//	//	l.Error("get chain config", "err", err)
//	//	return
//	//}
//	//
//	//
//	//var electionStatus rpc_interface.GetElectionStatus
//	//if err := client.Call(&electionStatus, getDipperinRpcMethodByName("GetElectionStatus"), txHash); err != nil {
//	//	l.Error("GetElectionStatus error", "err", err)
//	//	return
//	//}
//	//
//	//startHeight := electionStatus.VerifierRound*config.SlotSize
//	//endHeight  := startHeight + config.SlotSize-1
//	//if electionStatus.ElectionStatus == rpc_interface.Invalid {
//	//	l.Error("the election has been invalid")
//	//	electionMap.Delete(txHash)
//	//} else if electionStatus.ElectionStatus == rpc_interface.Packaged {
//	//	l.Info("election success the round as a verifier is:", "round", electionStatus.VerifierRound,"startHeight",startHeight,"endHeight",endHeight)
//	//	electionMap.Delete(txHash)
//	//} else if electionStatus.ElectionStatus == rpc_interface.WaitPackaged {
//	//	l.Info("the Election Status is WaitPackaged")
//	//} else {
//	//	l.Error("the Election Status error")
//	//}
//	return
//}

// get tx
/*func handleElectionTx(txHash common.Hash, address common.Address) {

	 l.Info("check election tx hash", "hash", txHash)

	// go to the chain to query tx
	var resp rpc_interface.TransactionResp
	if err := client.Call(&resp, getDipperinRpcMethodByName("Transaction"), txHash); err != nil {
		l.Error("Call Transaction", "err", err)
		return
	}

	if resp.BlockNumber == 0 {
		l.Debug("no get tx in chain")
		return
	}

	// get the block number of the transaction on
	blockNumber := resp.BlockNumber

	// block number converted to block current round
	var config chain_config.ChainConfig
	if err := client.Call(&config, getDipperinRpcMethodByName("GetChainConfig")); err != nil {
		l.Error("get chain config", "err", err)
		return
	}

	curRound := config.GetRound(blockNumber)

	l.Debug("curRound", "round", curRound)

	// get the verifiers of the specified round
	targetRound := curRound + 2

	l.Debug("targetRound", "targetRound", targetRound)
	var verifiers []common.Address
	if err := client.Call(&verifiers, getDipperinRpcMethodByName("GetVerifiersBySlot"), targetRound); err != nil {
		l.Error("GetVerifiersBySlot", "err", err)
		return
	}

	l.Info("verifiers len", "length", len(verifiers))

	// determine if it is a verifier
	isVerifier := false
	for i := range verifiers {
		if verifiers[i].IsEqual(address) {
			isVerifier = true
			break
		}
	}

	if isVerifier {
		startHeight := targetRound* config.SlotSize
		endHeight := targetRound*config.SlotSize - 1

		l.Info("address is verifier", "address", address.Hex(), "round", targetRound, "start height",
			startHeight, "end height", endHeight)

	}

	electionMap.Delete(txHash)
}*/
