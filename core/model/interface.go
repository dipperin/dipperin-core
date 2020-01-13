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

package model

import (
	"crypto/ecdsa"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/bloom"
	"math/big"
)

//go:generate mockgen -destination=../../tests/mock/model-mock/abstract_header_mock.go -package=model_mock github.com/dipperin/dipperin-core/core/model AbstractHeader
type AbstractHeader interface {
	GetNumber() uint64
	Hash() common.Hash
	GetPreHash() common.Hash
	EncodeRlpToBytes() ([]byte, error)
	GetStateRoot() common.Hash
	CoinBaseAddress() common.Address
	DuplicateHeader() AbstractHeader
	IsEqual(header AbstractHeader) bool
	SetVerificationRoot(newRoot common.Hash)
	GetSeed() common.Hash
	GetProof() []byte
	GetGasLimit() uint64
	GetGasUsed() uint64
	GetMinerPubKey() *ecdsa.PublicKey
	GetTimeStamp() *big.Int
	GetInterLinkRoot() common.Hash
	GetDifficulty() common.Difficulty
	GetRegisterRoot() common.Hash
	SetRegisterRoot(root common.Hash)
	//GetBloomLog() model.Bloom
	//IsEqual(header *Header) bool
}

type AbstractBody interface {
	GetTxsSize() int
	GetTxByIndex(i int) AbstractTransaction
	EncodeRlpToBytes() ([]byte, error)
	GetInterLinks() InterLink
	//GetReceipts() ([]*model.Receipt, error)
}

type AbstractBlock interface {
	Version() uint64
	Number() uint64
	IsSpecial() bool
	Difficulty() common.Difficulty
	PreHash() common.Hash
	Seed() common.Hash
	RefreshHashCache() common.Hash
	Hash() common.Hash
	EncodeRlpToBytes() ([]byte, error)
	TxIterator(cb func(int, AbstractTransaction) error) error
	TxRoot() common.Hash
	Timestamp() *big.Int
	Nonce() common.BlockNonce
	StateRoot() common.Hash
	SetStateRoot(root common.Hash)
	GetRegisterRoot() common.Hash
	SetRegisterRoot(root common.Hash)
	FormatForRpc() interface{}
	SetNonce(nonce common.BlockNonce)
	CoinBaseAddress() common.Address
	GetTransactionFees() *big.Int
	CoinBase() *big.Int
	GetTransactions() []*Transaction
	GetInterlinks() InterLink
	SetInterLinkRoot(root common.Hash)
	GetInterLinkRoot() (root common.Hash)
	SetInterLinks(inter InterLink)
	GetAbsTransactions() []AbstractTransaction
	GetBloom() iblt.Bloom
	Header() AbstractHeader
	Body() AbstractBody
	TxCount() int
	GetEiBloomBlockData(reqEstimator *iblt.HybridEstimator) *BloomBlockData
	GetBlockTxsBloom() *iblt.Bloom
	VerificationRoot() common.Hash
	SetVerifications(vs []AbstractVerification)
	VersIterator(func(int, AbstractVerification, AbstractBlock) error) error
	GetVerifications() []AbstractVerification
	SetReceiptHash(receiptHash common.Hash)
	GetReceiptHash() common.Hash
	//GetBloomLog() model.Bloom
	//SetBloomLog(bloom model.Bloom)
	//GasLimit() uint64
}

type PriofityCalculator interface {
	GetElectPriority(common.Hash, uint64, *big.Int, uint64) (uint64, error)
	GetReputation(uint64, *big.Int, uint64) (uint64, error)
}

//go:generate mockgen -destination=./../chain-communication/transaction_mock_test.go -package=chain_communication github.com/dipperin/dipperin-core/core/model AbstractTransaction
type AbstractTransaction interface {
	Size() common.StorageSize
	Amount() *big.Int
	CalTxId() common.Hash
	Nonce() uint64
	To() *common.Address
	Sender(singer Signer) (common.Address, error)
	SenderPublicKey(signer Signer) (*ecdsa.PublicKey, error)
	EncodeRlpToBytes() ([]byte, error)
	GetSigner() Signer
	GetType() common.TxType
	ExtraData() []byte
	Cost() *big.Int
	EstimateFee() *big.Int
	GetGasPrice() *big.Int
	GetGasLimit() uint64
	AsMessage(checkNonce bool) (Message, error)
	PaddingReceipt(parameters ReceiptPara)
	PaddingActualTxFee(fee *big.Int)
	GetReceipt() *Receipt
	GetActualTxFee() (fee *big.Int)
}

//go:generate mockgen -destination=./../economy-model/verification_mock_test.go -package=economy_model github.com/dipperin/dipperin-core/core/model AbstractVerification
//go:generate mockgen -destination=./../../cmd/utils/ver-halt-check/verification_mock_test.go -package=ver_halt_check github.com/dipperin/dipperin-core/core/model AbstractVerification
//go:generate mockgen -destination=./../cs-chain/chain-writer/middleware/verification_mock_test.go -package=middleware github.com/dipperin/dipperin-core/core/model AbstractVerification
type AbstractVerification interface {
	GetHeight() uint64
	GetRound() uint64
	GetViewID() uint64
	GetType() VoteMsgType
	GetBlockId() common.Hash
	GetAddress() common.Address
	GetBlockHash() string
	Valid() error
	HaltedVoteValid(verifiers []common.Address) error
}

// TxDifference returns a new set which is the difference between a and b.
func TxDifference(a, b []AbstractTransaction) []AbstractTransaction {
	keep := make([]AbstractTransaction, 0, len(a))

	remove := make(map[common.Hash]struct{})
	for _, tx := range b {
		remove[tx.CalTxId()] = struct{}{}
	}

	for _, tx := range a {
		if _, ok := remove[tx.CalTxId()]; !ok {
			keep = append(keep, tx)
		}
	}

	return keep
}
