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

package components

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	model2 "github.com/dipperin/dipperin-core/core/csbft/model"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third-party/crypto"
	"github.com/dipperin/dipperin-core/third-party/crypto/cs-crypto"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/dipperin/dipperin-core/third-party/p2p"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"

	"time"
)

//*******separating line***********************************

func TestCsBftFetcher_OnStart(t *testing.T) {
	fetcher := NewFetcher(&fackConn{})
	err := fetcher.Start()

	assert.NoError(t, err)
	assert.Equal(t, true, fetcher.IsRunning())

	fetcher.Stop()
	assert.Equal(t, false, fetcher.IsRunning())
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
	log.Debug("CreateKey", "add", add)
	return []*ecdsa.PrivateKey{key1, key2, key3, key4}, add
}

type FakePeer struct {
	id   string
	addr common.Address
	rw   p2p.MsgReadWriter
}

func newFakePeer(alice, bob string) (p1, p2 *FakePeer) {
	_, addr := CreateKey()
	rw1, rw2 := p2p.MsgPipe()
	p1 = &FakePeer{
		id:   alice,
		addr: addr[0],
		rw:   rw2,
	}
	p2 = &FakePeer{
		id:   bob,
		addr: addr[1],
		rw:   rw1,
	}
	return
}

type FakeCsBft struct {
	peerList []*FakePeer

	//FakeNodeContext
}

func (fp FakePeer) SendMsg(msgCode uint64, msg interface{}) error {
	//log.Info("send msg", "msgCode", msgCode, "peerid", fp.id)
	return p2p.Send(fp.rw, uint64(msgCode), msg)
}

func (fp FakePeer) ReadMsg() (p2p.Msg, error) {
	return fp.rw.ReadMsg()
}

func (fc FakeCsBft) SendFetchBlockMsg(msgCode uint64, from common.Address, msg *model2.FetchBlockReqDecodeMsg) error {
	//return nil
	return fc.peerList[0].SendMsg(msgCode, msg)
}

func newFakeCsBft(p1 *FakePeer) *FakeCsBft {
	return &FakeCsBft{peerList: []*FakePeer{p1}}
}

func TestCsBftFetcher_FetchBlock(t *testing.T) {
	log.InitLogger(log.LvlDebug)
	log.PBft.Logger = log.SetInitLogger(log.DefaultLogConf, "fetcher_test")
	_, addr := CreateKey()
	alice, bob := newFakePeer("alice", "bob")

	//I'm Alice
	fetcherA := NewFetcher(newFakeCsBft(bob))
	fetcherA.BaseService.Start()
	go func() {
		for {
			err := ReadDataMsg(fetcherA, "Alice", bob)
			assert.NoError(t, err)
		}
	}()

	//I'm Bob
	fetcherB := NewFetcher(newFakeCsBft(alice))
	errOS := fetcherB.OnStart()
	assert.Equal(t, errOS, nil)
	go func() {
		for {
			ReadDataMsg(fetcherB, "Bob", alice)
		}
	}()

	b1 := createBlock()
	block := fetcherA.FetchBlock(addr[1], b1.Hash())

	assert.Equal(t, false, fetcherA.IsFetching(b1.Hash()))
	assert.Equal(t, block.Hash(), b1.Hash())
	assert.Equal(t, len(fetcherA.requests), 0)
}

func TestCsBftFetcher_onFetchBlock(t *testing.T) {
	log.PBft.Logger = log.SetInitLogger(log.DefaultLogConf, "fetcher_test")
	alice, bob := newFakePeer("alice", "bob")
	//I'm Alice
	fetcherA := NewFetcher(newFakeCsBft(bob))
	fetcherA.Start()
	_, addr := CreateKey()
	go func() {
		for {
			err := ReadDataMsg(fetcherA, "Alice", bob)
			assert.NoError(t, err)
		}
	}()

	//I'm Bob
	fetcherB := NewFetcher(newFakeCsBft(alice))
	errOS := fetcherB.OnStart()
	assert.Equal(t, errOS, nil)

	go func() {
		for {
			ReadDataMsg(fetcherB, "Bob", alice)
		}
	}()

	b1 := createBlock()
	b2 := createBlock()
	b3 := createBlock()
	b4 := createBlock()
	b5 := createBlock()
	b6 := createBlock()

	fetcherA.FetchBlock(addr[0], b1.Hash())
	fetcherA.FetchBlock(addr[0], b2.Hash())
	fetcherA.FetchBlock(addr[0], b3.Hash())
	fetcherA.FetchBlock(addr[0], b4.Hash())
	fetcherA.FetchBlock(addr[0], b5.Hash())
	fetcherA.FetchBlock(addr[0], b6.Hash())

	fmt.Println("FetcherA are fetching: ", len(fetcherA.requests))
	err := fetcherA.IsFetching(b1.Hash())
	assert.Equal(t, err, false)
	fbrm := FetchBlockReqMsg{
		BlockHash:  b1.Hash(),
		ResultChan: make(chan model.AbstractBlock),
	}
	fetcherA.onFetchBlock(&fbrm)
	assert.Equal(t, fbrm.BlockHash, b1.Hash())
}

func TestCsBftFetcher_FetchBlockResp(t *testing.T) {
	_, bob := newFakePeer("alice", "bob")

	//I'm Alice
	fetcherA := NewFetcher(newFakeCsBft(bob))
	err := fetcherA.OnStart()
	assert.Equal(t, err, nil)
	b1 := &FakeBlock{uint64(1), common.HexToHash("0x232"), nil}
	fetcherA.FetchBlockResp(&FetchBlockRespMsg{2, b1})

}

func TestCsBftFetcher_FetchBlockResp2(t *testing.T) {
	log.InitLogger(log.LvlDebug)
	_, bob := newFakePeer("alice", "bob")

	//I'm Alice
	fetcherA := NewFetcher(newFakeCsBft(bob))
	fetcherA.Start()
	b1 := &FakeBlock{uint64(1), common.HexToHash("0x232"), nil}
	b2 := &FakeBlock{uint64(1), common.HexToHash("0x233"), nil}
	fetcherA.FetchBlockResp(&FetchBlockRespMsg{1, b1})
	fetcherA.FetchBlockResp(&FetchBlockRespMsg{2, b2})
	fetcherA.BaseService.Stop()

}

func createBlock() (block *model.Block) {
	header1 := model.NewHeader(1, 100, common.HexToHash("1111"), common.HexToHash("1111"), common.HexToDiff("1fffffff"), big.NewInt(324234), common.HexToAddress("032f14"), common.BlockNonceFromInt(432423))
	var txs1 []*model.Transaction
	var msg1 []model.AbstractVerification
	block1 := model.NewBlock(header1, txs1, msg1)
	return block1
}

func ReadDataMsg(fetcher *CsBftFetcher, name string, peer *FakePeer) error {
	msg, readErr := peer.ReadMsg()
	if readErr != nil {
		return readErr
	}
	switch model2.CsBftMsgType(msg.Code) {
	case model2.TypeOfFetchBlockReqMsg:
		log.Info(fmt.Sprintf("%v receive FetchBlock Request Msg", name))
		var m model2.FetchBlockReqDecodeMsg
		if err := msg.Decode(&m); err != nil {
			log.Debug("Decode FetchBlock Request Msg fail")
			return err
		}
		if err := peer.SendMsg(uint64(model2.TypeOfFetchBlockRespMsg), &FetchBlockRespMsg{
			MsgId: m.MsgId,
			Block: createBlock(),
		}); err != nil {
			log.Warn("send fetch block to client failed", "err", err)
			return err
		}
		log.Info(fmt.Sprintf("%v send FetchBlock Response Msg", name))

	case model2.TypeOfFetchBlockRespMsg:
		log.Info(fmt.Sprintf("%v receive FetchBlock Response Msg", name))
		var m model2.FetchBlockRespDecodeMsg
		if err := msg.Decode(&m); err != nil {
			log.Debug("Decode FetchBlock Response Msg fail")
			return err
		}
		fetcher.fetchRespChan <- (&FetchBlockRespMsg{
			MsgId: m.MsgId,
			Block: m.Block,
		})
	}
	return nil
}

//*******separating line***********************************

type fackConn struct {
}

func (kc *fackConn) SendFetchBlockMsg(msgCode uint64, from common.Address, msg *model2.FetchBlockReqDecodeMsg) error {
	//fmt.Println("send fetch block")
	return nil
}

func TestCsBftFetcher_IsNotRunning(t *testing.T) {
	_, bob := newFakePeer("alice", "bob")
	fetcherA := NewFetcher(newFakeCsBft(bob))
	_, addr := CreateKey()
	b1 := createBlock()
	fetcherA.FetchBlock(addr[0], b1.Hash())
	fetcherA.IsFetching(b1.Hash())
	fetcherA.OnReset()
}

// TODO
func TestCsBftFetcher_IsRunning2(t *testing.T) {
	//pbft_log.InitPbftLogger()
	log.InitLogger(log.LvlInfo)
	_, bob := newFakePeer("alice", "bob")
	fetcherA := NewFetcher(newFakeCsBft(bob))
	_, addr := CreateKey()
	b1 := createBlock()
	fetcherA.Start()
	go func() {
		time.Sleep(time.Millisecond * 5)
		fetcherA.Stop()
	}()
	fetcherA.FetchBlock(addr[0], b1.Hash())
	fetcherA.IsFetching(b1.Hash())
	fbrm := FetchBlockReqMsg{BlockHash: b1.Hash(), ResultChan: make(chan model.AbstractBlock)}
	fetcherA.onFetchBlock(&fbrm)
}
