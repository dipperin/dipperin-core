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

package chain

//import (
//	"errors"
//	"fmt"
//	"github.com/dipperin/dipperin-core/common"
//	"github.com/dipperin/dipperin-core/core/chain/state-processor"
//	"github.com/dipperin/dipperin-core/core/model"
//	"github.com/dipperin/dipperin-core/third-party/log"
//	"github.com/ethereum/go-ethereum/rlp"
//	"math/big"
//	"sort"
//)
//
//type PredicateFunc func(LightProofs) bool
//type Proofs []Proof
//
//type LightProof struct {
//	Header model.AbstractHeader
//	Link   model.InterLink
//}
//
//type LightProofs []*LightProof
//type Headers []model.AbstractHeader
//
//// TODO: modified when deployed
//var reader = NewFakeFullChain()
//var TestGenesis = reader.Genesis()
//var minDiff = common.HexToDiff("0x20ffffff")
//
//type Proof struct {
//	Prefix LightProofs
//	Suffix LightProofs
//	PreMap HeadersMap
//	SufMap HeadersMap
//}
//
//func NewLightProof(h model.AbstractHeader, l model.InterLink) *LightProof {
//	return &LightProof{
//		Header: h,
//		Link:   l,
//	}
//}
//
//// validlightproof tries to verify the given header, interlinks are indeed valid
//// by checking whether the derived root is equal to the header root
//func (l LightProof) ValidLightProof() bool {
//	if root := model.DeriveSha(l.Link); !root.IsEqual(l.Header.GetInterLinkRoot()) {
//		return false
//	}
//	return true
//}
//
//func (l LightProof) Hash() common.Hash {
//	return l.Header.Hash()
//}
//
//func (l LightProof) GetNumber() uint64 {
//	return l.Header.GetNumber()
//}
//
//func (p *Proof) InsertPrefix(h model.AbstractHeader, l model.InterLink) {
//	if !p.PreMap.Insert(h) {
//		p.Prefix = append(p.Prefix, NewLightProof(h, l))
//	}
//}
//
//func (p *Proof) InsertBatchPrefix(headers LightProofs) {
//	for _, h := range headers {
//		p.InsertPrefix(h.Header, h.Link)
//	}
//}
//
//func (p *Proof) InsertSuffix(h model.AbstractHeader, l model.InterLink) {
//	if !p.SufMap.Insert(h) {
//		p.Suffix = append(p.Suffix, NewLightProof(h, l))
//	}
//}
//
//func (p *Proof) InsertBatchSuffix(headers LightProofs) {
//	for _, h := range headers {
//		p.InsertSuffix(h.Header, h.Link)
//	}
//}
//
//func (p *Proof) Sort() {
//	sort.Sort(p.Suffix)
//	sort.Sort(p.Prefix)
//}
//
//func validProof(headers LightProofs) bool {
//	for i, header := range headers {
//		if i == 0 {
//			continue
//		}
//
//		//linkRoot := header.GetInterLinkRoot()
//		links := header.Link
//
//		found := false
//
//		preHeader := headers[i-1]
//
//		for _, h := range links {
//			if !found {
//				found = h.IsEqual(preHeader.Hash())
//			}
//		}
//
//		if !found {
//			return false
//		}
//	}
//
//	return true
//}
//
////todo How to ensure that this verification can be passed when proof is generated normally?
//func (p *Proof) Valid(fullchain Chain) (err error) {
//	if !validProof(p.Prefix) {
//		for _, h := range p.Prefix {
//			fmt.Println(h.Hash(), h.GetNumber())
//			fmt.Println(h.Link)
//		}
//		log.Error("interlink prefix invalid")
//		err = errors.New("interlink prefix invalid")
//		return
//	}
//
//	if !validProof(p.Suffix) {
//		log.Error("interlink suffix invalid")
//		err = errors.New("interlink suffix invalid")
//		return
//	}
//
//	if !validProof(append(p.Prefix, p.Suffix...)) {
//		log.Error("interlink proof invalid")
//		err = errors.New("interlink proof invalid")
//		return
//	}
//
//	return nil
//}
//
//func NewProofs() *Proof {
//	return &Proof{
//		Prefix: make(LightProofs, 0),
//		Suffix: make(LightProofs, 0),
//		PreMap: *NewHeadersMap(),
//		SufMap: *NewHeadersMap(),
//	}
//}
//
//type FakeFullChainInterlink struct {
//	blocks []model.AbstractBlock
//}
//
//func (f FakeFullChainInterlink) GetLastChangePoint(block model.AbstractBlock) *uint64 {
//	panic("implement me")
//}
//
//func (f FakeFullChainInterlink) GetSlot(block model.AbstractBlock) *uint64 {
//	panic("implement me")
//}
//
//func (f FakeFullChainInterlink) IsChangePoint(block model.AbstractBlock, isProcessPackageBlock bool) bool {
//	panic("implement me")
//}
//
//func (f FakeFullChainInterlink) GetSlotByNum(num uint64) *uint64 {
//	panic("implement me")
//}
//
//func (FakeFullChainInterlink) InsertBlocks(blocks []model.AbstractBlock) error {
//	panic("implement me")
//}
//
//func (f FakeFullChainInterlink) Genesis() model.AbstractBlock {
//	return f.blocks[0]
//}
//
//func (f FakeFullChainInterlink) CurrentBlock() model.AbstractBlock {
//	return f.blocks[len(f.blocks)-1]
//}
//
//func (FakeFullChainInterlink) CurrentHeader() model.AbstractHeader {
//	panic("implement me")
//}
//
//func (FakeFullChainInterlink) GetCurrVerifiers() []common.Address {
//	panic("implement me")
//}
//
//func (FakeFullChainInterlink) GetVerifiers(round uint64) []common.Address {
//	panic("implement me")
//}
//
//func (FakeFullChainInterlink) GetNextVerifiers() []common.Address {
//	panic("implement me")
//}
//
//func (FakeFullChainInterlink) CurrentSeed() (common.Hash, uint64) {
//	panic("implement me")
//}
//
//func (f FakeFullChainInterlink) GetBlock(hash common.Hash, number uint64) model.AbstractBlock {
//	n := f.GetBlockNumber(hash)
//	if n != nil {
//		return f.GetBlockByNumber(*n)
//	}
//	return nil
//}
//
//func (f FakeFullChainInterlink) GetBlockByHash(hash common.Hash) model.AbstractBlock {
//	for _, b := range f.blocks {
//		if hash.IsEqual(b.Hash()) {
//			return b
//		}
//	}
//	return nil
//}
//
//func (f FakeFullChainInterlink) GetLatestNormalBlock() model.AbstractBlock {
//	return nil
//}
//
//func (f FakeFullChainInterlink) GetBlockByNumber(number uint64) model.AbstractBlock {
//	return f.blocks[number]
//}
//
//func (FakeFullChainInterlink) HasBlock(hash common.Hash, number uint64) bool {
//	panic("implement me")
//}
//
//func (FakeFullChainInterlink) GetBody(hash common.Hash) model.AbstractBody {
//	panic("implement me")
//}
//
//func (FakeFullChainInterlink) GetBodyRLP(hash common.Hash) rlp.RawValue {
//	panic("implement me")
//}
//
//func (FakeFullChainInterlink) GetHeader(hash common.Hash, number uint64) model.AbstractHeader {
//	panic("implement me")
//}
//
//func (f FakeFullChainInterlink) GetHeaderByHash(hash common.Hash) model.AbstractHeader {
//	return f.GetBlockByHash(hash).Header()
//}
//
//func (f FakeFullChainInterlink) GetHeaderByNumber(number uint64) model.AbstractHeader {
//	if number > uint64(len(f.blocks)) {
//		return nil
//	}
//	return f.blocks[number].Header()
//}
//
//func (FakeFullChainInterlink) GetHeaderRLP(hash common.Hash) rlp.RawValue {
//	panic("implement me")
//}
//
//func (FakeFullChainInterlink) HasHeader(hash common.Hash, number uint64) bool {
//	panic("implement me")
//}
//
//func (f FakeFullChainInterlink) GetBlockNumber(hash common.Hash) *uint64 {
//	number := uint64(0)
//	for _, h := range f.blocks {
//		if hash.IsEqual(h.Hash()) {
//			number = h.Number()
//			return &number
//		}
//	}
//
//	return nil
//}
//
//func (FakeFullChainInterlink) GetTransaction(txHash common.Hash) (model.AbstractTransaction, common.Hash, uint64, uint64) {
//	panic("implement me")
//}
//
//func (FakeFullChainInterlink) GetStateStorage() state_processor.StateStorage {
//	panic("implement me")
//}
//
//func (FakeFullChainInterlink) CurrentState() (*state_processor.AccountStateDB, error) {
//	panic("implement me")
//}
//
//func (FakeFullChainInterlink) StateAtByBlockNumber(num uint64) (*state_processor.AccountStateDB, error) {
//	panic("implement me")
//}
//
//func (FakeFullChainInterlink) StateAtByStateRoot(root common.Hash) (*state_processor.AccountStateDB, error) {
//	panic("implement me")
//}
//
//func (FakeFullChainInterlink) TxMaxSize(block model.AbstractBlock) common.StorageSize {
//	panic("implement me")
//}
//
//func (FakeFullChainInterlink) TxMinSize(block model.AbstractBlock) common.StorageSize {
//	panic("implement me")
//}
//
//func (FakeFullChainInterlink) TxMaxAmount(block model.AbstractBlock) *big.Int {
//	panic("implement me")
//}
//
//func (FakeFullChainInterlink) TxMinAmount(block model.AbstractBlock) *big.Int {
//	panic("implement me")
//}
//
//func (FakeFullChainInterlink) ValidTx(tx model.AbstractTransaction) error {
//	panic("implement me")
//}
//
//func (FakeFullChainInterlink) ValidateBlock(block model.AbstractBlock) error {
//	panic("implement me")
//}
//
//func (FakeFullChainInterlink) GetSeenCommit(height uint64) []model.AbstractVerification {
//	panic("implement me")
//}
//
//func (FakeFullChainInterlink) GetVerifiersByBlock(height uint64) []common.Address {
//	panic("implement me")
//}
//
//func (f *FakeFullChainInterlink) SaveBlock(block model.AbstractBlock, seenCommits []model.AbstractVerification) error {
//	f.blocks = append(f.blocks, block)
//	return nil
//}
//
//func (FakeFullChainInterlink) LastNumberBySlot() uint64 {
//	panic("implement me")
//}
//
//func NewFakeFullChain() *FakeFullChainInterlink {
//	header1 := model.NewHeader(1, 0, common.HexToHash("1111"), common.HexToHash("1111"), minDiff, big.NewInt(324234), common.HexToAddress("032f14"), common.BlockNonceFromInt(432423))
//	block1 := model.NewBlock(header1, nil, nil)
//	inter := common.HexToHash("")
//	block1.SetInterLinks([]common.Hash{inter})
//	block1.SetInterLinkRoot(model.DeriveSha(block1.GetInterlinks()))
//
//	f := &FakeFullChainInterlink{
//		blocks: make([]model.AbstractBlock, 0),
//	}
//
//	f.blocks = append(f.blocks, block1)
//	return f
//}
//
//func (bs LightProofs) Len() int {
//	return len(bs)
//}
//
//func (bs LightProofs) Less(i, j int) bool {
//	return bs[i].GetNumber() < bs[j].GetNumber()
//}
//
//func (bs LightProofs) Swap(i, j int) {
//	bs[i], bs[j] = bs[j], bs[i]
//}
//
//func NewHeadersMap() *HeadersMap {
//	return &HeadersMap{
//		headers: make(map[common.Hash]bool),
//	}
//}
//
//type HeadersMap struct {
//	headers map[common.Hash]bool
//}
//
//func (h *HeadersMap) Insert(header model.AbstractHeader) bool {
//	if !h.In(header) {
//		h.headers[header.Hash()] = true
//		return false
//	}
//	return true
//}
//
//func (h HeadersMap) In(header model.AbstractHeader) bool {
//	_, ok := h.headers[header.Hash()]
//	return ok
//}
//
//// getHeaders implements the [A:Z], returns the sub-chain of the fullChain from
//// A (inclusive) to Z (exclusive)
//func getHeaders(fullchain Chain, start uint64, end uint64) (headers LightProofs) {
//	for start < end {
//		var block model.AbstractBlock
//		block = fullchain.GetBlockByNumber(start)
//		proof := NewLightProof(block.Header(), block.GetInterlinks())
//
//		if block == nil {
//			return nil
//		}
//		headers = append(headers, proof)
//		start++
//	}
//	return headers
//}
//
//// AfterTarget implements [b:] operator that returns all the headers
//// after and including b. Noted that, the parameter headers are assumed
//// sorted, and the b is not necessarily in headers, but is in the underlying
//// support block chain of headers.
//func AfterTarget(headers LightProofs, b model.AbstractHeader) (res LightProofs) {
//	left := 0
//	right := headers.Len() - 1
//	tar := b.GetNumber()
//
//	if len(headers) == 0 {
//		return
//	}
//
//	//sort.Sort(headers)
//	var mid int
//	for left <= right {
//		mid = (left + right) / 2
//		midNumber := headers[mid].GetNumber()
//
//		if midNumber < tar {
//			left = mid + 1
//		} else if midNumber > tar {
//			right = mid - 1
//		} else {
//			res = make(LightProofs, len(headers[mid:]))
//			copy(res, headers[mid:])
//			return
//		}
//	}
//
//	res = make(LightProofs, len(headers[left:]))
//	copy(res, headers[left:])
//	return
//}
//
//func GetSuffixProof(fullchain Chain, end uint64, m int, k uint64) (p *Proof, err error) {
//	if end > fullchain.CurrentBlock().Number() {
//		err = errors.New("blockchain is too short to get enough proof")
//		return
//	}
//
//	mu := len(fullchain.GetBlockByNumber(end - k).GetInterlinks())
//	mu--
//	//fmt.Println("end", end)
//	//fmt.Println("k", k)
//	//fmt.Println("end - k", end-k)
//	//fmt.Println(fullchain.GetBlockByNumber(end - k).GetInterlinks())
//	//fmt.Println("mu", mu)
//
//	genesis := fullchain.Genesis()
//	header := NewLightProof(genesis.Header(), genesis.GetInterlinks())
//	//todo: Why didn't the genesis block be added to the suffix proof here?
//
//	//todo: Whether it should be end+1-k, end+1 here (input end is the block number of the last block)
//	suf := getHeaders(fullchain, end-k, end)
//
//	// using map to ensure the items are unique
//	//p := NewHeadersMap()
//	p = NewProofs()
//	p.InsertBatchSuffix(suf)
//	for index := mu; index >= 0; index-- {
//		subChain := getHeaders(fullchain, header.GetNumber(), end-k)
//		alpha := UpChain(subChain, index)
//		p.InsertBatchPrefix(alpha)
//
//		if len(alpha) < m {
//			continue
//		}
//		header = alpha[len(alpha)-m]
//	}
//	p.Sort()
//	return p, nil
//}
//
//func followDown(fullchain Chain, lo model.AbstractHeader, hi model.AbstractHeader) (auxiliary LightProofs) {
//	header := hi
//	mu := model.HashLevel(header.Hash(), header.GetDifficulty().DiffToTarget())
//	//fmt.Println(lo.Hash())
//	for !header.Hash().IsEqual(lo.Hash()) {
//		link := fullchain.GetBlockByNumber(header.GetNumber()).GetInterlinks()
//		//fmt.Println("mu", mu)
//		//fmt.Println("len(link)", len(link))
//		//fmt.Println("number", header.GetNumber(), header.Hash())
//		//fmt.Println(link)
//		if mu >= len(link) {
//			mu = len(link) - 1
//		}
//		prime := fullchain.GetHeaderByHash(link[mu])
//
//		if prime.GetNumber() < lo.GetNumber() {
//			mu--
//		} else {
//			auxiliary = append(auxiliary, NewLightProof(header, link))
//			header = prime
//		}
//	}
//	sort.Sort(LightProofs(auxiliary))
//	return
//}
//
//func getLevel(b model.AbstractHeader) int {
//	return model.HashLevel(b.Hash(), b.GetDifficulty().DiffToTarget())
//}
//
//// UpChain select the blocks that have more links than
//// the given level, and forms a super chain
//func UpChain(headers LightProofs, level int) (superChain LightProofs) {
//	for _, block := range headers {
//		if getLevel(block.Header) >= level {
//			superChain = append(superChain, block)
//		}
//	}
//	return superChain
//}
//
//func BestArg(pre LightProofs, b model.AbstractHeader, m int) (arg int) {
//	length := m
//	subChain := AfterTarget(pre, b)
//	//fmt.Println(pre)
//	//fmt.Println(b)
//	//fmt.Println(len(subChain))
//	for mu := 0; length >= m; mu ++ {
//		//fmt.Println("m:", m, "|upchain|:", length, "mu:", mu, "2^mu:", 1<<uint(mu))
//		//fmt.Println("2^mu * length", length*(1<<uint(mu)))
//		if length*(1<<uint(mu)) > arg {
//			arg = length * (1 << uint(mu))
//		}
//
//		subChain = UpChain(subChain, mu)
//		length = len(subChain)
//	}
//	//fmt.Println("---------------------------------------------")
//	//fmt.Println("best arg", arg)
//	return
//}
//
//// BlockGreater implements the $ >=_m $ operator defined in
//// the NIPPOW paper Algorithm 4.
//// We assume all the input LightProofs are sorted.
//func BlockGreater(A LightProofs, B LightProofs, threshold int) bool {
//	La := len(A) - 1
//	Lb := len(B) - 1
//
//	var b model.AbstractHeader
//
//	for La >= 0 && Lb >= 0 {
//		if A[La].GetNumber() == B[Lb].GetNumber() {
//			if A[La].Hash().IsEqual(B[Lb].Hash()) {
//				// Found greatest common ancestor
//				b = A[La].Header
//				break
//			} else {
//				// Equal number but hash is different, diverged here
//				// continue to find common ancestor
//				La--
//				Lb--
//			}
//		} else if A[La].GetNumber() > B[Lb].GetNumber() {
//			La--
//		} else {
//			Lb--
//		}
//	}
//
//	if b == nil {
//		b = TestGenesis.Header()
//	}
//	//fmt.Println(commonBlocks[len(commonBlocks)-1])
//	//fmt.Println(A)
//	//fmt.Println(B)
//	//fmt.Println(A[La])
//	//fmt.Println(B[Lb])
//	bestA := BestArg(A, b, threshold)
//	bestB := BestArg(B, b, threshold)
//	return bestA > bestB
//}
//
//func GetInfixProof(fullchain Chain, endPoint uint64, threshold int, suffix uint64, target uint64) (p *Proof, err error) {
//	p, err = GetSuffixProof(fullchain, endPoint, threshold, suffix)
//	targetHeader := fullchain.GetHeaderByNumber(target)
//	if err != nil {
//		return
//	}
//	for _, header := range p.Prefix {
//		if header.GetNumber() >= target {
//			auxiliary := followDown(fullchain, targetHeader, header.Header)
//			p.InsertBatchPrefix(auxiliary)
//			break
//		}
//	}
//
//	p.Sort()
//	return
//}
//
//func GetGoodInfixProof(fullchain Chain, endPoint uint64, threshold int, suffix uint64, target uint64, delta float64) (p *Proof, err error) {
//	p, err = GetGoodSuffixProof(fullchain, endPoint, threshold, suffix, delta)
//	if err != nil {
//		return
//	}
//	targetHeader := fullchain.GetHeaderByNumber(target)
//	for _, header := range p.Prefix {
//		if header.GetNumber() >= target {
//			auxiliary := followDown(fullchain, targetHeader, header.Header)
//			p.InsertBatchPrefix(auxiliary)
//			break
//		}
//	}
//
//	p.Sort()
//	return
//}
//
//func GetGoodSuffixProof(fullchain Chain, end uint64, m int, k uint64, delta float64) (p *Proof, err error) {
//	if end > fullchain.CurrentBlock().Number() {
//		err = errors.New("blockchain is too short to get enough proof")
//		return
//	}
//
//	mu := len(fullchain.GetBlockByNumber(end - k).GetInterlinks())
//	mu--
//	//fmt.Println("end", end)
//	//fmt.Println("k", k)
//	//fmt.Println("end - k", end-k)
//	//fmt.Println(fullchain.GetBlockByNumber(end - k).GetInterlinks())
//	//fmt.Println("mu", mu)
//	genesis := fullchain.Genesis()
//	header := NewLightProof(genesis.Header(), genesis.GetInterlinks())
//	suf := getHeaders(fullchain, end-k, end)
//
//	// using map to ensure the items are unique
//	//p := NewHeadersMap()
//	p = NewProofs()
//	p.InsertBatchSuffix(suf)
//	for index := mu; index >= 0; index-- {
//		subChain := getHeaders(fullchain, header.GetNumber(), end-k)
//		alpha := UpChain(subChain, index)
//		p.InsertBatchPrefix(alpha)
//
//		if len(alpha) < m {
//			continue
//		}
//
//		under := GetUnderlyingHeaders(fullchain, alpha)
//		if Good(alpha, under, m, index, delta) {
//			header = alpha[len(alpha)-m]
//		} else {
//			log.Info("header hash is not good, using lower level")
//		}
//	}
//	p.Sort()
//	return p, nil
//}
//
//// super is assumed sorted
//func GetUnderlyingHeaders(reader Chain, super LightProofs) (under LightProofs) {
//	if len(super) < 1 {
//		return nil
//	}
//	start := super[0].GetNumber()
//	end := super[len(super)-1].GetNumber()
//
//	under = make(LightProofs, end-start+1)
//	for i := start; i <= end; i ++ {
//		b := reader.GetBlockByNumber(i)
//		under[i-start] = NewLightProof(b.Header(), b.GetInterlinks())
//	}
//
//	return
//}
//
//// we assume under contains all the elements in super
//func DownChain(under, super LightProofs) (down LightProofs) {
//	if len(under) < len(super) {
//		return nil
//	}
//
//	tail := super[len(super)-1].GetNumber()
//	after := AfterTarget(under, super[0].Header)
//	for _, h := range after {
//		down = append(down, h)
//		if h.GetNumber() > tail {
//			break
//		}
//	}
//
//	return
//}
//
//func LocalGood(super, under LightProofs, mu int, delta float64) bool {
//	// TODO: we only need to compare the length
//	// input argument be modified to use integer
//	return len(super)*(1<<uint(mu)) > int((1-delta)*float64(len(under)))
//}
//
//// chain is assumed to be the 0-chain
//func SuperchainQuality(chain LightProofs, m int, mu int, delta float64) bool {
//	mP := m
//	superchain := UpChain(chain, mu)
//	length := len(superchain)
//
//	for mP < length {
//		super := superchain[(length - mP):]
//		//Since the chain is ordered by blockNumber, use AfterTarget directly here.，
//		// You don't have to perform underlying on the chain as written by paper, and you can't get chainReader here.
//		under := AfterTarget(chain, super[0].Header)
//		if !LocalGood(super, under, mu, delta) {
//			return false
//		}
//		mP++
//	}
//
//	return true
//}
//
//func MultilevelQuality(super, under LightProofs, k1 int, mu int, delta float64) bool {
//	muP := 0
//	for muP < mu {
//		//Here is not verified by any original set of under the original paper, directly use the under
//		//The author believes that this problem can only be solved by special attacks described in paper, which can be restricted by chain consensus conditions.
//		star := DownChain(under, super)
//		star = UpChain(star, muP)
//		length := len(UpChain(star, muP))
//		if length >= k1 {
//			//fmt.Println("|C* mu-up|*2^mu", len(UpChain(star, mu))*(1<<uint(mu)), "(1-delta)", delta, "mu'", muP)
//			//fmt.Println("length*2^mu'", (1<<uint(muP))*length, "length", length)
//			if len(UpChain(star, mu))*(1<<uint(mu)) < int((1-delta)*float64((1<<uint(muP))*length)) {
//				return false
//			}
//		} else {
//			// higher level assumed to has less block
//			break
//		}
//		muP++
//		//fmt.Println("---------------------------")
//	}
//	return true
//}
//
//func Good(super, under LightProofs, m int, mu int, delta float64) bool {
//	if !SuperchainQuality(under, m, mu, delta) {
//		return false
//	}
//
//	if !MultilevelQuality(super, under, m, mu, delta) {
//		return false
//	}
//
//	return true
//}
//
//// VerifySuffix implements algorithm 2
//func VerifySuffix(proofs Proofs, k uint64, m int, Q PredicateFunc) bool {
//	prePrime := LightProofs{NewLightProof(TestGenesis.Header(), TestGenesis.GetInterlinks())}
//	var sufPrime LightProofs
//
//	for _, p := range proofs {
//		// if validChain && |suf| = k && pre >=_m prePrime
//		if len(p.Suffix) == int(k) && BlockGreater(p.Prefix, prePrime, m) {
//			prePrime = p.Prefix
//			sufPrime = p.Suffix
//		}
//	}
//
//	return Q(sufPrime)
//}
//
//type BlockByHash struct {
//	M    map[common.Hash]*LightProof
//	Flag map[common.Hash]bool
//}
//
//func NewBlockByHash() *BlockByHash {
//	return &BlockByHash{
//		M:    make(map[common.Hash]*LightProof),
//		Flag: make(map[common.Hash]bool),
//	}
//}
//
//func (b *BlockByHash) Insert(h common.Hash, proof *LightProof) {
//	b.M[h] = proof
//}
//
//func (b *BlockByHash) Traverse(h common.Hash) {
//	b.Flag[h] = true
//}
//
////determine if the block hash has been recorded
//func (b BlockByHash) Traversed(h common.Hash) bool {
//	return b.Flag[h]
//}
//
//func (b BlockByHash) Get(h common.Hash) (*LightProof, bool) {
//	p, ok := b.M[h]
//	return p, ok
//}
//
//func Ancestors(b *LightProof, header *BlockByHash) LightProofs {
//	if b.Hash().IsEqual(TestGenesis.Hash()) {
//		return LightProofs{b}
//	}
//
//	var headers LightProofs
//	headersMap := NewHeadersMap()
//	//fmt.Println("finding ancestor at")
//	//fmt.Println(b.Hash(), b.GetNumber())
//	//fmt.Println("link")
//	//fmt.Println(b.Link)
//	//fmt.Println("iterating interlink")
//	for i := len(b.Link) - 1; i > 0; i-- {
//		hash := b.Link[i]
//		//fmt.Println("hash", hash)
//		if h, ok := header.Get(hash); ok {
//			// To ensure every edge in the DAG is traversed only once,
//			// otherwise it would cost so much time if we follow the algorithm
//			// in the original paper
//			if !header.Traversed(hash) {
//				//fmt.Println(h.Hash(), h.GetNumber())
//				//fmt.Println("entering recursive")
//				header.Traverse(hash)
//				aux := Ancestors(h, header)
//				for _, h := range aux {
//					if !headersMap.Insert(h.Header) {
//						headers = append(headers, h)
//					}
//				}
//			}
//		}
//	}
//
//	if !headersMap.Insert(b.Header) {
//		headers = append(headers, b)
//	}
//	return headers
//}
//
//// VerifyInfix implements algorithm 7
//func VerifyInfix(proofs Proofs, k uint64, m int, D PredicateFunc) bool {
//	headers := NewBlockByHash()
//	headers.Insert(TestGenesis.Hash(), NewLightProof(TestGenesis.Header(), TestGenesis.GetInterlinks()))
//	//fmt.Println("Setup blockById")
//	for _, p := range proofs {
//		for _, B := range p.Prefix {
//			headers.Insert(B.Hash(), B)
//			//fmt.Println(B.Hash(), B.GetNumber())
//		}
//	}
//
//	prePrime := LightProofs{NewLightProof(TestGenesis.Header(), TestGenesis.GetInterlinks())}
//
//	for _, p := range proofs {
//		// if validChain && |suf| = k && pre >=_m prePrime
//		if len(p.Suffix) == int(k) && BlockGreater(p.Prefix, prePrime, m) {
//			prePrime = p.Prefix
//		}
//	}
//
//	an := Ancestors(prePrime[len(prePrime)-1], headers)
//
//	return D(an)
//}
