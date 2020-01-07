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
	"crypto/ecdsa"
	"errors"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/bloom"
	"github.com/dipperin/dipperin-core/core/model"
	model2 "github.com/dipperin/dipperin-core/core/vm/model"

	"github.com/dipperin/dipperin-core/tests/factory"
	"github.com/dipperin/dipperin-core/tests/factory/g-testData"
	"github.com/ethereum/go-ethereum/ethdb"
	"math/big"
	"time"
)

var (
	DataBaseErr = errors.New("fakeDataBase test error")
	DecoderErr  = errors.New("fakeDecoder test error")
	BatchErr    = errors.New("fakeBatch test error")
	BodyErr     = errors.New("fakeBody test error")
	HeaderErr   = errors.New("fakeHeader test error")
)

func createSignedTx(nonce uint64, amount *big.Int, to common.Address) *model.Transaction {
	key1, _ := model.CreateKey()
	fs1 := model.NewSigner(big.NewInt(1))
	testTx1 := model.NewTransaction(nonce, to, amount, g_testData.TestGasPrice, g_testData.TestGasLimit, []byte{})
	signedTx, _ := testTx1.SignTx(key1, fs1)
	return signedTx
}

func createBlock(num uint64) *model.Block {
	header := model.NewHeader(1, num, common.Hash{}, common.HexToHash("123456"), common.HexToDiff("1fffffff"), big.NewInt(time.Now().UnixNano()), factory.AliceAddrV, common.BlockNonce{})

	// tx list
	to := common.HexToAddress(common.AddressContractCreate)
	tx1 := createSignedTx(0, g_testData.TestValue, to)
	tx2 := createSignedTx(0, g_testData.TestValue, factory.BobAddrV)
	txList := []*model.Transaction{tx1, tx2}

	// vote
	var voteList []model.AbstractVerification
	block := model.NewBlock(header, txList, voteList)

	// calculate block nonce
	model.CalNonce(block)
	block.RefreshHashCache()
	return block
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

type fakeBody struct{}

func (body fakeBody) Version() uint64 {
	panic("implement me")
}

func (body fakeBody) Number() uint64 {
	panic("implement me")
}

func (body fakeBody) IsSpecial() bool {
	panic("implement me")
}

func (body fakeBody) Difficulty() common.Difficulty {
	panic("implement me")
}

func (body fakeBody) PreHash() common.Hash {
	panic("implement me")
}

func (body fakeBody) Seed() common.Hash {
	panic("implement me")
}

func (body fakeBody) RefreshHashCache() common.Hash {
	panic("implement me")
}

func (body fakeBody) Hash() common.Hash {
	panic("implement me")
}

func (body fakeBody) TxIterator(cb func(int, model.AbstractTransaction) error) error {
	panic("implement me")
}

func (body fakeBody) TxRoot() common.Hash {
	panic("implement me")
}

func (body fakeBody) Timestamp() *big.Int {
	panic("implement me")
}

func (body fakeBody) Nonce() common.BlockNonce {
	panic("implement me")
}

func (body fakeBody) StateRoot() common.Hash {
	panic("implement me")
}

func (body fakeBody) SetStateRoot(root common.Hash) {
	panic("implement me")
}

func (body fakeBody) GetRegisterRoot() common.Hash {
	panic("implement me")
}

func (body fakeBody) SetRegisterRoot(root common.Hash) {
	panic("implement me")
}

func (body fakeBody) FormatForRpc() interface{} {
	panic("implement me")
}

func (body fakeBody) SetNonce(nonce common.BlockNonce) {
	panic("implement me")
}

func (body fakeBody) CoinBaseAddress() common.Address {
	panic("implement me")
}

func (body fakeBody) GetTransactionFees() *big.Int {
	panic("implement me")
}

func (body fakeBody) CoinBase() *big.Int {
	panic("implement me")
}

func (body fakeBody) GetTransactions() []*model.Transaction {
	panic("implement me")
}

func (body fakeBody) GetInterlinks() model.InterLink {
	panic("implement me")
}

func (body fakeBody) SetInterLinkRoot(root common.Hash) {
	panic("implement me")
}

func (body fakeBody) GetInterLinkRoot() (root common.Hash) {
	panic("implement me")
}

func (body fakeBody) SetInterLinks(inter model.InterLink) {
	panic("implement me")
}

func (body fakeBody) GetAbsTransactions() []model.AbstractTransaction {
	panic("implement me")
}

func (body fakeBody) GetBloom() iblt.Bloom {
	panic("implement me")
}

func (body fakeBody) Header() model.AbstractHeader {
	panic("implement me")
}

func (body fakeBody) Body() model.AbstractBody {
	panic("implement me")
}

func (body fakeBody) TxCount() int {
	panic("implement me")
}

func (body fakeBody) GetEiBloomBlockData(reqEstimator *iblt.HybridEstimator) *model.BloomBlockData {
	panic("implement me")
}

func (body fakeBody) GetBlockTxsBloom() *iblt.Bloom {
	panic("implement me")
}

func (body fakeBody) VerificationRoot() common.Hash {
	panic("implement me")
}

func (body fakeBody) SetVerifications(vs []model.AbstractVerification) {
	panic("implement me")
}

func (body fakeBody) VersIterator(func(int, model.AbstractVerification, model.AbstractBlock) error) error {
	panic("implement me")
}

func (body fakeBody) GetVerifications() []model.AbstractVerification {
	panic("implement me")
}

func (body fakeBody) SetReceiptHash(receiptHash common.Hash) {
	panic("implement me")
}

func (body fakeBody) GetReceiptHash() common.Hash {
	panic("implement me")
}

func (body fakeBody) GetBloomLog() model2.Bloom {
	panic("implement me")
}

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

type fakeHeader struct{}

func (h fakeHeader) GetBloomLog() model2.Bloom {
	panic("implement me")
}

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

func (h fakeHeader) GetMinerPubKey() *ecdsa.PublicKey {
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
