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
	"crypto/ecdsa"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/consts"
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third-party/crypto"
	"github.com/dipperin/dipperin-core/third-party/crypto/cs-crypto"
	"math/big"
)

type FBlock struct {
	coinbase common.Address
	number uint64
	perHash common.Hash
	difficult common.Difficulty
	txs       []*model.Transaction
	vers      []model.AbstractVerification
	interLink model.InterLink
}

func NewFBlock(number uint64, perHash common.Hash, difficulty common.Difficulty, coinbase common.Address) *FBlock {
	return &FBlock{
		number:    number,
		perHash:   perHash,
		difficult: difficulty,
		coinbase:  coinbase,
	}
}

func (receiver *FBlock) AddTx(tx *model.Transaction) *FBlock {
	receiver.txs = append(receiver.txs, tx)
	return receiver
}

func (receiver *FBlock) AddTxs(txs []*model.Transaction) *FBlock {
	receiver.txs = append(receiver.txs, txs...)
	return receiver
}

func (receiver *FBlock) AddVers(vers []model.AbstractVerification) *FBlock {
	receiver.vers = append(receiver.vers, vers...)
	return receiver
}

func (receiver *FBlock) AddInterLink(il model.InterLink) *FBlock {
	receiver.interLink = append(receiver.interLink, il...)
	return receiver
}

//func (receiver *FBlock) CreateBlock(reader state_processor.AccountDBChainReader, storage state_processor.StateStorage, preBlockStateRoot common.Hash, economyModel economy_model.EconomyModel) *model.Block {
//	header := model.NewHeader(1, receiver.number, receiver.perHash, common.HexToHash("seed"), receiver.difficult, big.NewInt(time.Now().UnixNano()), receiver.coinbase, common.BlockNonceFromInt(432423))
//
//	block := model.NewBlock(header, receiver.txs, receiver.vers)
//
//	// process txs for state root
//	aDB, err := chain.NewBlockProcessor(reader, preBlockStateRoot, storage)
//	if err != nil {
//		panic("create account db failed: " + err.Error())
//	}
//	if err = aDB.Process(block, economyModel); err != nil {
//		panic("process block failed: " + err.Error())
//	}
//	sRoot, err := aDB.Finalise()
//	if err != nil {
//		panic("finalise account db failed: " + err.Error())
//	}
//	block.SetStateRoot(sRoot)
//
//	if receiver.interLink.Len() > 0 {
//		block.SetInterLinks(receiver.interLink)
//	}
//
//	return block
//}

type FTx struct {
	prk   *ecdsa.PrivateKey
	nonce uint64

	txs []*model.Transaction
}

func NewFTx(prkStr *ecdsa.PrivateKey, nonce uint64) *FTx {
	return &FTx{
		prk:   prkStr,
		nonce: nonce,
	}
}


func (receiver *FTx) CreateTx(toAddress common.Address, amount *big.Int, fee *big.Int, data []byte) *FTx {
	tx := FactoryCreateTx(receiver.prk, receiver.nonce, toAddress, amount, fee, data)
	receiver.nonce++

	receiver.txs = append(receiver.txs, tx)
	return receiver
}

func (receiver *FTx) GetTxs() []*model.Transaction {
	return receiver.txs
}

// address string to common.Address
func str2Address(addr string) common.Address {
	return common.StringToAddress(addr)
}

// prk -> common.Address
func factoryGetAddress(prk *ecdsa.PrivateKey) common.Address {
	return cs_crypto.GetNormalAddress(prk.PublicKey)
}

func factoryCreatePrivateKey(prkStr string) *ecdsa.PrivateKey {
	prk, err := crypto.HexToECDSA(prkStr)

	if err != nil {
		panic(err)
	}

	return prk
}

func FactoryCreateTx(senderPrk *ecdsa.PrivateKey, nonce uint64, toAddress common.Address,
	amount *big.Int, fee *big.Int, data []byte) *model.Transaction {

	signer := model.NewMercurySigner(chain_config.GetChainConfig().ChainId)

	tx := model.NewTransaction(nonce, toAddress, amount, fee, data)

	tmpTx, err := tx.SignTx(senderPrk, signer)

	if err != nil {
		panic(err)
	}

	return tmpTx
}

func factoryGenPrk() *ecdsa.PrivateKey {
	prk, err := crypto.GenerateKey()

	if err != nil {
		panic(err)
	}

	return prk
}

type FAddress struct {
	address2Prk map[common.Address]*ecdsa.PrivateKey

	addressList []common.Address

	count int
}

func NewFAddress(count int) *FAddress {
	fa := &FAddress{
		address2Prk: make(map[common.Address]*ecdsa.PrivateKey),
	}

	for i := 0; i < count; i++ {
		prk := factoryGenPrk()
		address := cs_crypto.GetNormalAddress(prk.PublicKey)
		fa.address2Prk[address] = prk
		fa.addressList = append(fa.addressList, address)
	}

	fa.count = count

	return fa
}

func (receiver *FAddress) Alloc() map[common.Address]*big.Int {
	tmpMap := make(map[common.Address]*big.Int)

	for i := 0; i < receiver.count; i++ {
		tmpMap[receiver.addressList[0]] = big.NewInt(100 * consts.DIP)
	}

	return tmpMap
}

func (receiver *FAddress) GetPrk(index int) *ecdsa.PrivateKey {
	if index > receiver.count-1 {
		panic("error: Transboundary")
	}

	return receiver.address2Prk[receiver.addressList[index]]
}

func (receiver *FAddress) GetAddress(index int) common.Address {
	if index > receiver.count-1 {
		panic("error: Transboundary")
	}

	return receiver.addressList[index]
}
