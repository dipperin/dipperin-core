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


package csbftnode

import (
	"testing"
	"crypto/ecdsa"
	"math/big"
	"github.com/dipperin/dipperin-core/core/bloom"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/model"
	"fmt"
	"github.com/dipperin/dipperin-core/third-party/p2p"
	"net"
	"github.com/dipperin/dipperin-core/third-party/crypto"
	"github.com/dipperin/dipperin-core/third-party/crypto/cs-crypto"
	"github.com/dipperin/dipperin-core/core/csbft/components"
	"github.com/dipperin/dipperin-core/core/csbft/state-machine"
	//"github.com/ethereum/go-ethereum/node"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/stretchr/testify/assert"
	"github.com/dipperin/dipperin-core/third-party/log/pbft_log"
)

// New fakeMsgSender
type FackMsgSender struct{

}

func (fms *FackMsgSender) BroadcastMsg(msgCode uint64, msg interface{}) {
	fmt.Println("broad cast msg")
}

func (fms *FackMsgSender) SendReqRoundMsg(msgCode uint64, from []common.Address, msg interface{}) {
	fmt.Println("send req round msg")
}

func (fms *FackMsgSender) BroadcastEiBlock(block model.AbstractBlock) {
	fmt.Println("broadcast block")
}


// New FakeValidtor
type FakeValidtor struct{}

func (FakeValidtor) FullValid(block model.AbstractBlock) (error) {
	return nil
}

func (FakeValidtor) Valid(block model.AbstractBlock) error {
	return nil
}

func TestCsBft_AddPeer2(t *testing.T) {
	fc := NewFakeFullChain()
	sks, _ := CreateKey()
	fs := newFackSigner(sks[1])
	fetcher := components.NewFetcher(nil)
	config := &state_machine.BftConfig{fc,fetcher,fs,&FackMsgSender{}, &FakeValidtor{}}

	node1 := NewCsBft(config)
	node1.SetFetcher(fetcher)
	log.Info("TestCsBft_AddPeer2", "fetcher", node1.fetcher)
	err := node1.Start()
	assert.NoError(t, err)
	adderr := node1.AddPeer(nil)
	assert.NoError(t,adderr)
	//node1.FastSync()
	//assert.Equal(t, node1.FastSync(), true)
	//fmt.Println(node1.blockPool)
	//node1.OnNewWaitVerifyBlock()
	node1.Stop()


}

var as = "289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232031"
var bs = "289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232032"
var cs = "289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232033"
var ds = "289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232034"

func CreateKey() (keys []*ecdsa.PrivateKey, adds []common.Address) {
	key1, err1 := crypto.HexToECDSA(as)
	key2, err2 := crypto.HexToECDSA(bs)
	key3, err3 := crypto.HexToECDSA(cs)
	key4, err4 := crypto.HexToECDSA(ds)
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		return nil, nil
	}
	add := []common.Address{}
	add = append(add, cs_crypto.GetNormalAddress(key1.PublicKey))
	add = append(add, cs_crypto.GetNormalAddress(key2.PublicKey))
	add = append(add, cs_crypto.GetNormalAddress(key3.PublicKey))
	add = append(add, cs_crypto.GetNormalAddress(key4.PublicKey))
	return []*ecdsa.PrivateKey{key1, key2, key3, key4}, add
}



//-----------------
// Fake fullchain

func NewFakeFullChain() *FakeFullChain {
	_, v := CreateKey()
	fb := FakeBlock{uint64(0), common.HexToHash("0x232"), nil}
	return &FakeFullChain{
		Height:     uint64(0),
		Validators: v,
		Blocks:     &fb,
	}
}

type FakeFullChain struct {
	Height     uint64
	Validators []common.Address
	Blocks     model.AbstractBlock
	commits    []model.AbstractVerification
	notify  func(height uint64)
}

func (fc *FakeFullChain) GetSeenCommit(height uint64) []model.AbstractVerification {
	return fc.commits
}

func (fc *FakeFullChain) SaveBlock(block model.AbstractBlock, seenCommits []model.AbstractVerification) error {
	pbft_log.Debug("save block","height",block.Number())
	fc.Blocks = block
	fc.Height = block.Number()
	fc.commits = seenCommits
	fc.notify(block.Number())
	return nil
}

func (fc *FakeFullChain) CurrentBlock() model.AbstractBlock {
	return fc.Blocks
}

func (fc *FakeFullChain) IsChangePoint(block model.AbstractBlock, isProcessPackageBlock bool) bool {
	if block.Number()%10 == 9 {
		return true
	}
	return false
}

func (fc *FakeFullChain) GetNextVerifiers() []common.Address {
	_, v := CreateKey()
	if fc.Height==9{
		return nil
	}
	return v
}

func (fc *FakeFullChain) GetCurrVerifiers() []common.Address {
	_, v := CreateKey()
	return v
}

func (fc *FakeFullChain) SetNewHeightNotifier(nc func(height uint64)){
	fc.notify = nc
}

//--------------------------
// Fake block
type FakeBlock struct {
	Height     uint64
	HeaderHash common.Hash
	Headers    model.AbstractHeader
}

func (fb *FakeBlock) SetReceiptHash(receiptHash common.Hash) {
	panic("implement me")
}

func (fb *FakeBlock) GetReceiptHash() common.Hash {
	panic("implement me")
}

func (fb *FakeBlock) IsSpecial() bool {
	panic("implement me")
}

func (fb *FakeBlock) GetRegisterRoot() common.Hash {
	panic("implement me")
}

func (fb *FakeBlock) SetRegisterRoot(root common.Hash) {
	panic("implement me")
}

func (fb *FakeBlock) Body() model.AbstractBody {
	panic("implement me")
}

func (fb *FakeBlock) GetBlockTxsBloom() *iblt.Bloom {
	panic("implement me")
}

func (fb *FakeBlock) Version() uint64 {
	panic("implement me")
}

func (fb *FakeBlock) Number() uint64 {
	return fb.Height
}

func (fb *FakeBlock) Difficulty() common.Difficulty {
	panic("implement me")
}

func (fb *FakeBlock) PreHash() common.Hash {
	panic("implement me")
}

func (fb *FakeBlock) Seed() common.Hash {
	panic("implement me")
}

func (fb *FakeBlock) RefreshHashCache() common.Hash {
	panic("implement me")
}

func (fb *FakeBlock) Hash() common.Hash {
	return fb.HeaderHash
}

func (fb *FakeBlock) EncodeRlpToBytes() ([]byte, error) {
	panic("implement me")
}

func (fb *FakeBlock) TxIterator(cb func(int, model.AbstractTransaction) (error)) (error) {
	panic("implement me")
}

func (fb *FakeBlock) TxRoot() common.Hash {
	panic("implement me")
}

func (fb *FakeBlock) Timestamp() *big.Int {
	panic("implement me")
}

func (fb *FakeBlock) Nonce() common.BlockNonce {
	panic("implement me")
}

func (fb *FakeBlock) StateRoot() common.Hash {
	panic("implement me")
}

func (fb *FakeBlock) SetStateRoot(root common.Hash) {
	panic("implement me")
}

func (fb *FakeBlock) FormatForRpc() interface{} {
	panic("implement me")
}

func (fb *FakeBlock) SetNonce(nonce common.BlockNonce) {
	panic("implement me")
}

func (fb *FakeBlock) CoinBaseAddress() common.Address {
	panic("implement me")
}

func (fb *FakeBlock) GetInterlinks() model.InterLink {
	panic("implement me")
}

func (fb *FakeBlock) SetInterLinkRoot(root common.Hash) {
	panic("implement me")
}

func (fb *FakeBlock) GetInterLinkRoot() (root common.Hash) {
	panic("implement me")
}

func (fb *FakeBlock) SetInterLinks(inter model.InterLink) {
	panic("implement me")
}

func (fb *FakeBlock) Header() model.AbstractHeader {
	return fb.Headers
}

func (fb *FakeBlock) GetEiBloomBlockData(reqEstimator *iblt.HybridEstimator) *model.BloomBlockData {
	panic("implement me")
}

func (fb *FakeBlock) SetVerifications(vs []model.AbstractVerification) {
	panic("implement me")
}

func (fb *FakeBlock) VersIterator(func(int, model.AbstractVerification, model.AbstractBlock) error) (error) {
	panic("implement me")
}

func (fb *FakeBlock) GetVerifications() ([]model.AbstractVerification) {
	panic("implement me")
}

func (fb *FakeBlock) GetTransactionFees() *big.Int {
	panic("implement me")
}

func (fb *FakeBlock) CoinBase() *big.Int {
	panic("implement me")
}

func (fb *FakeBlock) GetTransactions() []*model.Transaction {
	panic("implement me")
}

func (fb *FakeBlock) GetAbsTransactions() []model.AbstractTransaction {
	panic("implement me")
}

func (fb *FakeBlock) GetBloom() iblt.Bloom {
	panic("implement me")
}

func (fb *FakeBlock) VerificationRoot() common.Hash {
	panic("implement me")
}

func (fb *FakeBlock) TxCount() int {
	panic("implement me")
}


type fakeSigner struct {
	baseAddr   common.Address
	privateKey *ecdsa.PrivateKey
}

func newFackSigner(sk *ecdsa.PrivateKey) *fakeSigner {
	return &fakeSigner{privateKey: sk, baseAddr: cs_crypto.GetNormalAddress(sk.PublicKey)}
}

func (signer *fakeSigner) GetAddress() common.Address {
	return signer.baseAddr
}

func (signer *fakeSigner) SignHash(hash []byte) ([]byte, error) {
	//pb := crypto.CompressPubkey(&signer.privateKey.PublicKey)
	//log.Info("fake signer sign", "p k", hexutil.Encode(pb))
	return crypto.Sign(hash, signer.privateKey)
}

//-----------------------------
// Fake fetcher
type FakeFetcher struct {
	Block model.AbstractBlock
}

func (fc *FakeFetcher) Start() error {
	log.Info("FakeFetcher start")
	return nil
}

func (fc *FakeFetcher) OnStart() error {
	panic("implement me")
}

func (fc *FakeFetcher) Stop() {
	panic("implement me")
}

func (fc *FakeFetcher) OnStop() {
	panic("implement me")
}

func (fc *FakeFetcher) Reset() error {
	panic("implement me")
}

func (fc *FakeFetcher) OnReset() error {
	panic("implement me")
}

func (fc *FakeFetcher) IsRunning() bool {
	panic("implement me")
}

func (fc *FakeFetcher) Quit() <-chan struct{} {
	panic("implement me")
}

func (fc *FakeFetcher) String() string {
	panic("implement me")
}

func (fc *FakeFetcher) SetLogger(logger log.Logger) {
	panic("implement me")
}

func (fc *FakeFetcher) FetchBlock(from common.Address, blockHash common.Hash) model.AbstractBlock {
	return fc.Block
}

// mock peer
type tPeer struct {
	nodeType uint64
	name     string
	id       string
	address  common.Address
}

func (p *tPeer) NodeName() string {
	return p.name
}

func (p *tPeer) NodeType() uint64 {
	return p.nodeType
}

func (p *tPeer) SendMsg(msgCode uint64, msg interface{}) error {
	fmt.Println("send","code",msgCode)
	return nil
}

func (p *tPeer) ID() string {
	return p.id
}

func (p *tPeer) ReadMsg() (p2p.Msg, error) {
	panic("implement me")
}

func (p *tPeer) GetHead() (common.Hash, uint64) {
	panic("implement me")
}

func (p *tPeer) SetHead(head common.Hash, height uint64) {
	panic("implement me")
}

func (p *tPeer) GetPeerRawUrl() string {
	panic("implement me")
}

func (p *tPeer) DisconnectPeer() {
	panic("implement me")
}

func (p *tPeer) RemoteVerifierAddress() (addr common.Address) {
	return p.address
}

func (p *tPeer) RemoteAddress() net.Addr {
	panic("implement me")
}

func (p *tPeer) SetRemoteVerifierAddress(addr common.Address) {
	panic("implement me")
}

func (p *tPeer) SetNodeName(name string) {
	panic("implement me")
}

func (p *tPeer) SetNodeType(nt uint64) {
	panic("implement me")
}

func (p *tPeer) SetPeerRawUrl(rawUrl string) {
	panic("implement me")
}

func (p *tPeer) SetNotRunning() {
	panic("implement me")
}

func (p *tPeer) IsRunning() bool {
	panic("implement me")
}

func (p *tPeer) GetCsPeerInfo() *p2p.CsPeerInfo {
	panic("implement me")
}
