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


package rpc_interface

import (
	"github.com/dipperin/dipperin-core/common/hexutil"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/common"
)

// swagger:response CurBalanceResp
type CurBalanceResp struct {
	Balance *hexutil.Big `json:"balance"`
}

//TODO please confirm with chuang, may use CurrBalanceResp
// swagger:response CurStakeResp
type CurStakeResp struct {
	Stake *hexutil.Big `json:"balance"`
}

// swagger:response TransactionResp
type TransactionResp struct {
	Transaction *model.Transaction `json:"transaction"`
	BlockHash common.Hash `json:"blockHash"`
	BlockNumber uint64 `json:"blockNumber"`
	TxIndex uint64 `json:"transactionIndex"`
}

// swagger:response BlockResp
type BlockResp struct {
	Header model.Header `json:"header"`
	Body   model.Body   `json:"body"`
}

// swagger:response NewTransactionResp
type NewTransactionResp struct {
	TxHash common.Hash `json:"tx_hash"`
}

type SendTxReq struct {
	Tx *model.Transaction `json:"tx"`
}

type ERC20Resp struct {
	TxId common.Hash `json:"txid"`
	CtId common.Address `json:"ctid"`
}

//current practical verifiers resp
type PeerInfoResp struct {
	NodeId string
	Address common.Address
}


//Election resp
type ElectionResp struct {
	TxId common.Hash
	SendBlockNumber uint64
}

type ElectionStatus int

const (
	WaitPackaged ElectionStatus = iota
	Packaged
	Invalid
)

//get election status resp
type GetElectionStatus struct {
	ElectionStatus
	VerifierRound uint64
}

type VerifierStatus struct {
	Status string
	Stake *hexutil.Big
	Balance *hexutil.Big
	Reputation uint64
	IsCurrentVerifier bool
}







