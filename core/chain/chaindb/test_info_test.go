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

package chaindb

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/ethereum/go-ethereum/ethdb"
	"errors"
	"crypto/ecdsa"
	"math/big"
)

var (
	DataBaseErr = errors.New("fakeDataBase test error")
	DecoderErr  = errors.New("fakeDecoder test error")
	BatchErr    = errors.New("fakeBatch test error")
	BodyErr     = errors.New("fakeBody test error")
	HeaderErr     = errors.New("fakeHeader test error")
)

func createBlock(number uint64) *model.Block {
	return model.CreateBlock(number, common.HexToHash("123456"), 2)
}

func newDb() ethdb.Database {
	return ethdb.NewMemDatabase()
}

func newDecoder() model.BlockDecoder {
	return model.MakeDefaultBlockDecoder()
}

func newChainDB() *ChainDB {
	return NewChainDB(newDb(), newDecoder())
}

type fakeDataBase struct {
	err error
}

func (data fakeDataBase) Put(key []byte, value []byte) error {
	return DataBaseErr
}

func (data fakeDataBase) Delete(key []byte) error {
	return DataBaseErr
}

func (data fakeDataBase) Get(key []byte) ([]byte, error) {
	return nil, DataBaseErr
}

func (data fakeDataBase) Has(key []byte) (bool, error) {
	panic("implement me")
}

func (data fakeDataBase) Close() {
	panic("implement me")
}

func (data fakeDataBase) NewBatch() ethdb.Batch {
	return fakeBatch{err: data.err}
}

type fakeDecoder struct{}

func (decoder fakeDecoder) DecodeRlpHeaderFromBytes(data []byte) (model.AbstractHeader, error) {
	return nil, DecoderErr
}

func (decoder fakeDecoder) DecodeRlpBodyFromBytes(data []byte) (model.AbstractBody, error) {
	return nil, DecoderErr
}

func (decoder fakeDecoder) DecodeRlpBlockFromHeaderAndBodyBytes(headerB []byte, bodyB []byte) (model.AbstractBlock, error) {
	return nil, DecoderErr
}

func (decoder fakeDecoder) DecodeRlpBlockFromBytes(data []byte) (model.AbstractBlock, error) {
	panic("implement me")
}

func (decoder fakeDecoder) DecodeRlpTransactionFromBytes(data []byte) (model.AbstractTransaction, error) {
	panic("implement me")
}

type fakeBatch struct {
	err error
}

func (batch fakeBatch) Put(key []byte, value []byte) error {
	return batch.err
}

func (batch fakeBatch) Delete(key []byte) error {
	panic("implement me")
}

func (batch fakeBatch) ValueSize() int {
	panic("implement me")
}

func (batch fakeBatch) Write() error {
	return BatchErr
}

func (batch fakeBatch) Reset() {
	panic("implement me")
}

type fakeBody struct {}

func (body fakeBody) GetTxsSize() int {
	panic("implement me")
}

func (body fakeBody) GetTxByIndex(i int) model.AbstractTransaction {
	panic("implement me")
}

func (body fakeBody) EncodeRlpToBytes() ([]byte, error) {
	return nil, BodyErr
}

func (body fakeBody) GetInterLinks() model.InterLink {
	panic("implement me")
}

type fakeHeader struct {}

func (h fakeHeader) GetTimeStamp() *big.Int {
	panic("implement me")
}

func (h fakeHeader) GetGasLimit() uint64 {
	panic("implement me")
}

func (h fakeHeader) GetGasUsed() uint64 {
	panic("implement me")
}

func (h fakeHeader) GetNumber() uint64 {
	return uint64(0)
}

func (h fakeHeader) Hash() common.Hash {
	return common.Hash{}
}

func (h fakeHeader) GetPreHash() common.Hash {
	panic("implement me")
}

func (h fakeHeader) EncodeRlpToBytes() ([]byte, error) {
	return nil, HeaderErr
}

func (h fakeHeader) GetStateRoot() common.Hash {
	panic("implement me")
}

func (h fakeHeader) CoinBaseAddress() common.Address {
	panic("implement me")
}

func (h fakeHeader) DuplicateHeader() model.AbstractHeader {
	panic("implement me")
}

func (h fakeHeader) IsEqual(header model.AbstractHeader) bool {
	panic("implement me")
}

func (h fakeHeader) SetVerificationRoot(newRoot common.Hash) {
	panic("implement me")
}

func (h fakeHeader) GetSeed() common.Hash {
	panic("implement me")
}

func (h fakeHeader) GetProof() []byte {
	panic("implement me")
}

func (h fakeHeader) GetMinerPubKey() (*ecdsa.PublicKey) {
	panic("implement me")
}

func (h fakeHeader) GetInterLinkRoot() common.Hash {
	panic("implement me")
}

func (h fakeHeader) GetDifficulty() common.Difficulty {
	panic("implement me")
}

func (h fakeHeader) GetRegisterRoot() common.Hash {
	panic("implement me")
}

func (h fakeHeader) SetRegisterRoot(root common.Hash) {
	panic("implement me")
}
