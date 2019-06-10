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
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/bloom"
	crypto2 "github.com/dipperin/dipperin-core/third-party/crypto/cs-crypto"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/dipperin/dipperin-core/third-party/log/bloom_log"
	"github.com/dipperin/dipperin-core/third-party/log/pbft_log"
	"github.com/dipperin/dipperin-core/third-party/log/witch_log"
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
	"sort"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

var (
	EmptyTxRoot    = DeriveSha(Transactions{})
	EmptyVerfRoot  = DeriveSha(Verifications{})
	EmptyInterRoot = DeriveSha(InterLink{})
)

var (
	lock                    sync.Mutex
	DefaultBlockBloomConfig = iblt.NewBloomConfig(8, 4)
	DefaultInvBloomConfig   = iblt.NewInvBloomConfig(1<<12, 4)
	DefaultTxs              = 100
	DefaultGasLimit         = uint64(6666666666)
)

type Header struct {
	// version of this block generated
	Version uint64 `json:"version"  gencodec:"required"`
	// the height of the block
	Number uint64 `json:"number"  gencodec:"required"`
	// used for VRF
	Seed common.Hash `json:"seed"  gencodec:"required"`
	// seed proof
	Proof []byte `json:"proof"  gencodec:"required"`
	// miner public key
	MinerPubKey []byte `json:"miner_pub_key"  gencodec:"required"`
	//Previous block hash
	PreHash common.Hash `json:"pre_hash"  gencodec:"required"`
	// difficulty for this block
	Diff common.Difficulty `json:"diff"  gencodec:"required"`
	// timestamp for this block
	TimeStamp *big.Int `json:"timestamp"  gencodec:"required"`
	// the address of the miner who mined this block
	CoinBase common.Address `json:"coinbase"  gencodec:"required"`
	GasLimit uint64        `json:"gasLimit"         gencodec:"required"`
	GasUsed  uint64        `json:"gasUsed"          gencodec:"required"`
	// nonce needed to be mined by the miner
	Nonce common.BlockNonce `json:"nonce"  gencodec:"required"`
	//todo add bloom filter for Logs or txs
	Bloom *iblt.Bloom `json:"Bloom"        gencodec:"required"`
	// MPT trie Root for transaction
	TransactionRoot common.Hash `json:"txs_root"   gencodec:"required"`
	// MPT trie Root for accounts state
	//todo if we want to put normal accounts and contract accounts into different tree we can generate two MPT trie for them
	StateRoot common.Hash `json:"state_root" gencodec:"required"`
	// MPT trie Root for committed message
	VerificationRoot common.Hash `json:"verification_root"  gencodec:"required"`
	// MPT trie Root for interlink message
	InterlinkRoot common.Hash `json:"interlink_root"  gencodec:"required"`
	// MPT trie Root for register
	RegisterRoot common.Hash `json:"register_root"  gencodec:"required"`
	//add receipt hash
	ReceiptHash common.Hash `json:"receiptsRoot"     gencodec:"required"`
}

func (h *Header) GetGasLimit() uint64 {
	return h.GasLimit
}

func (h *Header) GetGasUsed() uint64 {
	return h.GasUsed
}

func (h *Header) IsEqual(header AbstractHeader) bool {
	panic("implement me")
}

func (h *Header) CoinBaseAddress() common.Address { return h.CoinBase }
func (h *Header) GetTimeStamp() *big.Int          { return h.TimeStamp }
func (h *Header) GetStateRoot() common.Hash {
	return h.StateRoot
}

func NewHeader(version uint64, num uint64, prehash common.Hash, seed common.Hash, diff common.Difficulty, time *big.Int, coinbase common.Address, nonce common.BlockNonce) *Header {

	return &Header{
		Version:     version,
		Number:      num,
		PreHash:     prehash,
		Seed:        seed,
		Diff:        diff,
		TimeStamp:   time,
		CoinBase:    coinbase,
		Nonce:       nonce,
		Bloom:       iblt.NewBloom(DefaultBlockBloomConfig),
		GasLimit:    DefaultGasLimit,
		Proof:       []byte{},
		MinerPubKey: []byte{},
	}
}

// CopyHeader creates a deep copy of a block header to prevent side effects from
// modifying a header variable.
func CopyHeader(h *Header) *Header {
	cpy := *h
	if cpy.TimeStamp = new(big.Int); h.TimeStamp != nil {
		cpy.TimeStamp.Set(h.TimeStamp)
	}
	return &cpy
}

type writeCounter common.StorageSize

func (c *writeCounter) Write(b []byte) (int, error) {
	*c += writeCounter(len(b))
	return len(b), nil
}

func (h Header) RlpBlockWithoutNonce() (rb []byte) {
	tmpH := h
	//without nonce
	tmpH.Nonce = common.BlockNonce{}
	//header rlp
	rb, _ = rlp.EncodeToBytes(tmpH)
	return
}

// compute Hash of the header
func (h *Header) Hash() common.Hash {
	ret := h.RlpBlockWithoutNonce()
	//header's rlp + nonce
	splice := append(ret, h.Nonce[:]...)

	return crypto2.Keccak256Hash(splice)
}

// compute VRFHash of the header
//func (h *Header) Hash() common.Hash {
//	return common.RlpHashKeccak256(h)
//}

func (h *Header) GetNumber() uint64 {
	return h.Number
}

func (h *Header) GetSeed() common.Hash {
	return h.Seed
}

func (h *Header) GetProof() []byte {
	return h.Proof
}

func (h *Header) GetMinerPubKey() *ecdsa.PublicKey {
	return crypto2.ToECDSAPub(h.MinerPubKey)
}

func (h *Header) GetPreHash() common.Hash {
	return h.PreHash
}

func (h *Header) GetInterLinkRoot() common.Hash {
	return h.InterlinkRoot
}

func (bh *Header) GetDifficulty() common.Difficulty {
	return bh.Diff
}

func (h *Header) GetRegisterRoot() common.Hash {
	return h.RegisterRoot
}

func (h *Header) SetRegisterRoot(root common.Hash) {
	h.RegisterRoot = root
}

func (h *Header) DuplicateHeader() AbstractHeader {
	cpy := *h
	return &cpy
}

func (h *Header) SetVerificationRoot(newRoot common.Hash) {
	h.VerificationRoot = newRoot
}

//todo for download
func (h *Header) Size() common.StorageSize {

	return common.StorageSize(unsafe.Sizeof(*h)) + common.StorageSize((len(h.Diff)+h.TimeStamp.BitLen())/8)
}

// the nonce is not comprised to calculate the hash value
func (h *Header) HashWithoutNonce() common.Hash {
	tmpH := *h
	tmpH.Nonce = common.BlockNonce{}
	return common.RlpHashKeccak256(tmpH)
}

func rlpHash(x interface{}) (h common.Hash, err error) {
	if b, e := rlp.EncodeToBytes(x); e != nil {
		return h, e
	} else {
		return crypto2.Keccak256Hash(b), nil
	}
	return
}

func (h *Header) String() string {
	return fmt.Sprintf(`Header(%s):
[	Version:	        %d
	Number:	            %d
	Seed:				%s
	PreHash:	        %s
	Difficulty:	        %s
	TimeStamp:	        %v
	CoinBase:           %s
	GasLimit        	%d
	GasUsed             %d
	Nonce:		        %s
	Bloom：         		%v
	TransactionRoot:    %s
	StateRoot:	        %s
	VerificationRoot:   %s
	InterlinkRoot:      %s
	RegisterRoot     	%s
	ReceiptHash      	%s]`, h.Hash().Hex(), h.Version, h.Number, h.Seed.Hex(), h.PreHash.Hex(), h.Diff.Hex(), h.TimeStamp, h.CoinBase.Hex(), h.GasLimit, h.GasUsed, h.Nonce.Hex(), h.Bloom.Hex(), h.TransactionRoot.Hex(), h.StateRoot.Hex(), h.VerificationRoot.Hex(), h.InterlinkRoot.Hex(), h.RegisterRoot.Hex(), h.ReceiptHash.Hex())
}

// swagger:response Body
type Body struct {
	Txs    []*Transaction         `json:"transactions"`
	Vers   []AbstractVerification `json:"commit_msg"`
	Inters InterLink              `json:"interlinks"`
}

func (b *Body) GetTxsSize() int {
	return len(b.Txs)
}

func (b *Body) GetTxByIndex(i int) AbstractTransaction {
	return b.Txs[i]
}

func (b *Body) GetInterLinks() InterLink {
	return b.Inters
}

type Block struct {
	header *Header
	body   *Body

	// caches
	hash atomic.Value `json:"-"`
	size atomic.Value `json:"-"`
	receipts atomic.Value `json:"-"`
}

func (b *Block) GasLimit() uint64 {
	return b.header.GasLimit
}

func (b *Block) GasUsed() uint64{
	return b.header.GasUsed
}

func (b *Block) IsSpecial() bool {
	if b.Difficulty().Equal(common.Difficulty{0}) && b.Nonce().IsEqual(common.BlockNonce{0}) {
		return true
	}
	return false
}

func (b *Block) SetDifficulty(diff common.Difficulty) {
	b.header.Diff = diff
}

func (b *Block) SetTimeStamp(timeStamp *big.Int) {
	b.header.TimeStamp = timeStamp
}

func (b *Block) GetRegisterRoot() common.Hash {
	return b.header.RegisterRoot
}

func (b *Block) SetRegisterRoot(root common.Hash) {
	b.header.RegisterRoot = root
}

func (b *Block) SetReceiptHash(receiptHash common.Hash) {
	b.header.ReceiptHash = receiptHash
}

func (b *Block) GetReceiptHash() common.Hash {
	return b.header.ReceiptHash
}

/*func (b *Block) PaddingReceipts(receipts model.Receipts){
	b.receipts.Store(receipts)
}

func (b *Block) GetReceipts() (model.Receipts,error){
	if r:=b.receipts.Load();r!=nil{
		return r.(model.Receipts),nil
	}

	return nil,g_error.BlockReceiptsAreEmpty
}*/

// Get block txs bloom
func (b *Block) GetBlockTxsBloom() *iblt.Bloom {
	bloom := iblt.NewBloom(iblt.DeriveBloomConfig(len(b.GetTransactions())))

	txs := b.GetTransactions()
	for _, tx := range txs {
		bloom.Digest(tx.CalTxId().Bytes())
	}
	bloom_log.Info("GetBlockTxsBloom", "txs len", len(txs))

	return bloom
}

// Get block invBloom
func (b *Block) getBlockInvBloom(reqEstimator *iblt.HybridEstimator) *iblt.InvBloom {
	estimator := iblt.NewHybridEstimator(reqEstimator.Config())

	for _, tx := range b.GetTransactions() {
		d := reqEstimator.NewData()
		d.SetBytes(tx.CalTxId().Bytes())
		estimator.Encode(d)
	}

	// hack optimize

	// peer knows the set difference, then peer constructs corresponding IBLT
	estimatedDiff := estimator.Decode(reqEstimator)
	estimatedConfig := iblt.NewInvBloomConfig(estimatedDiff*40, 4)
	invBloom := iblt.NewInvBloom(estimatedConfig)

	for _, tx := range b.GetTransactions() {
		invBloom.InsertRLP(tx.CalTxId(), tx)
	}

	return invBloom
}

func (b *Block) GetEiBloomBlockData(reqEstimator *iblt.HybridEstimator) *BloomBlockData {
	//startAt := time.Now()
	txs := b.GetTransactions()
	estimator := iblt.NewHybridEstimator(reqEstimator.Config())
	//log.Info("==block.go ==GetEiBloomBlockData==iblt.NewHybridEstimator()==1", "t", time.Now().Sub(startAt))

	startAt := time.Now()
	if len(txs) > DefaultTxs {
		fmt.Println("1")
		// run multiple procedure
		witch_log.Info("start MapWork")
		mapWorkEstimator := newMapWorkHybridEstimator(estimator)
		if err := RunWorkMap(mapWorkEstimator, txs); err != nil {
			witch_log.Info("RunWorkMap by mapWorkEstimator failed", "err", err)
		}
	} else {
		for _, tx := range txs {
			estimator.EncodeByte(tx.CalTxId().Bytes())
		}
	}
	estimatedConfig := estimator.DeriveConfig(reqEstimator)
	bloomConfig := iblt.DeriveBloomConfig(len(b.GetTransactions()))
	invBloom := iblt.NewGraphene(estimatedConfig, bloomConfig)
	bloom_log.Info("GetEiBloomBlockData", "txs len", len(txs))
	if len(txs) > DefaultTxs {

		// run multiple procedure
		mapWorkBloom := newMapWorkInvBloom(invBloom)
		if err := RunWorkMap(mapWorkBloom, txs); err != nil {
			witch_log.Info("RunWorkMap by mapWorkBloom failed", "err", err)
		}
	} else {
		for _, tx := range txs {
			invBloom.InsertRLP(tx.CalTxId(), tx)
			invBloom.Bloom().Digest(tx.CalTxId().Bytes())
		}
	}

	diff := time.Now().Sub(startAt)
	if diff > time.Second {
		witch_log.Info("GetEiBloomBlockData", "cost", diff, "len", len(txs), "num", b.Number())
	}

	invBloomRLP, err := rlp.EncodeToBytes(invBloom)
	if err != nil {
		log.Error("con't rlp invBloom", "block hash", b.Hash().Hex())
		return nil
	}
	return &BloomBlockData{
		Header:          b.header,
		BloomRLP:        invBloomRLP,
		PreVerification: b.Verifications(),
		Interlinks:      b.GetInterlinks(),
	}

}

func (b *Block) SetVerifications(vs []AbstractVerification) {
	b.body.Vers = vs
}

func (b *Block) GetVerifications() (result []AbstractVerification) {
	return b.body.Vers
}

func (b *Block) GetTransactions() []*Transaction {
	return b.body.Txs
}

func (b *Block) GetInterlinks() InterLink {
	return b.body.Inters
}

func (b *Block) GetAbsTransactions() []AbstractTransaction {
	var res []AbstractTransaction
	for _, tx := range b.body.Txs {
		res = append(res, tx)
	}
	return res
}

func (b *Block) SetNonce(nonce common.BlockNonce) {
	//log.Debug("set block nonce", "nonce", nonce.Hex())
	b.header.Nonce = nonce
	//log.Debug("after change nonce", "header nonce", b.header.Nonce.Hex())
}

func (b *Block) FormatForRpc() interface{} {
	return nil
}

func (b *Block) SetStateRoot(root common.Hash) {
	b.header.StateRoot = root
}

func (b *Block) SetInterLinkRoot(root common.Hash) {
	b.header.InterlinkRoot = root
}

func (b *Block) GetInterLinkRoot() (root common.Hash) {
	root = b.header.InterlinkRoot
	return
}

func (b *Block) SetInterLinks(inter InterLink) {
	b.body.Inters = inter
}

func (b Block) TxIterator(cb func(int, AbstractTransaction) error) error {
	for i, tx := range b.body.Txs {
		if err := cb(i, tx); err != nil {
			return err
		}
	}
	return nil
}

func (b Block) VersIterator(cb func(int, AbstractVerification, AbstractBlock) error) error {
	for i, verification := range b.body.Vers {
		if err := cb(i, verification, &b); err != nil {
			return err
		}
	}
	return nil
}

func (b Block) GetCoinbase() *big.Int {
	return nil
}

func (b Block) GetTransactionFees() *big.Int {
	tempfee := big.NewInt(0)
	for _, tx := range b.body.Txs {
		var addFee *big.Int
		if tx.GetType() == common.AddressTypeContractCreate || tx.GetType()==common.AddressTypeContract{
			addFee = tx.contractTxFee.Load().(*big.Int)
		}else{
			addFee = tx.Fee()
		}
		tempfee.Add(tempfee, addFee)
	}
	return tempfee
}

func NewBlock(header *Header, txs []*Transaction, msgs []AbstractVerification) *Block {
	var vers []AbstractVerification
	for _, m := range msgs {
		vers = append(vers, m)
	}

	body := &Body{Txs: txs, Vers: vers, Inters: InterLink{}} //TODO proofs should add wuhao
	// maybe don't use copy header
	b := &Block{header: CopyHeader(header), body: body}

	if len(txs) == 0 {
		b.header.TransactionRoot = EmptyTxRoot
	} else {
		b.header.TransactionRoot = DeriveSha(Transactions(txs))
		b.body.Txs = make(Transactions, len(txs))
		copy(b.body.Txs, txs)
	}

	pbft_log.Info("the calculated tx root is:","root",b.header.TransactionRoot.Hex())
	pbft_log.Info("the block txs is:","len",len(txs))
	for _,tx := range txs{
		pbft_log.Info("the tx is:","tx",tx)
	}

	// calculate verification Root
	if len(msgs) == 0 {
		b.header.VerificationRoot = EmptyVerfRoot
	} else {
		b.header.VerificationRoot = DeriveSha(Verifications(msgs))
	}
	//TODO: set default interlinkroot
	//todo add other mpt trie check
	//TODO proofs mpt trie check

	return b
}

func NewBlockWithLink(header *Header, txs []*Transaction, msgs []AbstractVerification, preLink InterLink) *Block {
	b := NewBlock(header, txs, msgs)

	linkList := NewInterLink(preLink, b)
	b.SetInterLinks(linkList)
	b.SetInterLinkRoot(DeriveSha(linkList))

	return b
}

func (b *Block) String() string {
	return fmt.Sprintf(`Header(%s):
[
	Header:	            %v
	Transactions:	    %v
    commit message:     %v
    interlins:          %v

]`, b.Hash().Hex(), b.header, b.body.Txs, b.body.Vers, b.body.Inters)
}

func (b *Block) RefreshHashCache() common.Hash {
	v := b.header.Hash()
	b.hash.Store(v)
	return v
}

func (b *Block) Hash() common.Hash {
	if h := b.hash.Load(); h != nil {
		return h.(common.Hash)
	}
	v := b.header.Hash()
	b.hash.Store(v)
	return v
}

func (b *Body) EncodeRlpToBytes() ([]byte, error) {
	return rlp.EncodeToBytes(b)
}
func (h *Header) EncodeRlpToBytes() ([]byte, error) {
	return rlp.EncodeToBytes(h)
}
func (b *Block) EncodeRlpToBytes() ([]byte, error) {
	return rlp.EncodeToBytes(blockForRlp{
		Header: b.header,
		Body:   b.body,
	})
}
func (b *Block) Version() uint64                 { return b.header.Version }
func (b *Block) Number() uint64                  { return b.header.Number }
func (b *Block) PreHash() common.Hash            { return b.header.PreHash }
func (b *Block) Seed() common.Hash               { return b.header.Seed }
func (b *Block) Timestamp() *big.Int             { return new(big.Int).Set(b.header.TimeStamp) }
func (b *Block) CoinBaseAddress() common.Address { return b.header.CoinBase }
func (b *Block) Nonce() common.BlockNonce        { return b.header.Nonce }
func (b *Block) Difficulty() common.Difficulty   { return b.header.Diff }
func (b *Block) TxRoot() common.Hash             { return b.header.TransactionRoot }
func (b *Block) StateRoot() common.Hash          { return b.header.StateRoot }
func (b *Block) VerificationRoot() common.Hash   { return b.header.VerificationRoot }
func (b *Block) GetBloom() iblt.Bloom            { return *b.header.Bloom }

func (b *Block) CoinBase() *big.Int {
	// TODO
	return big.NewInt(0)
}

func (b *Block) TxCount() int { return len(b.body.Txs) }

func (b *Block) Header() AbstractHeader {
	if b.header != nil {
		return CopyHeader(b.header)
	}
	return nil
}

func (b *Block) GetHeader() AbstractHeader {
	if b.header != nil {
		return CopyHeader(b.header)
	}
	return nil
}

func (b *Block) Body() AbstractBody {
	return b.body
	//return &Body{Txs: b.body.Txs, Vers: b.body.Vers}
}

func (b *Block) Verifications() []AbstractVerification { return b.GetVerifications() }

//func (b *Block) Proofs() []*Proof 	{ return b.body.Proofs }

func (b *Block) Transaction(hash common.Hash) *Transaction {
	for _, transaction := range b.body.Txs {
		if transaction.CalTxId() == hash {
			return transaction
		}
	}
	return nil
}

// EncodeToIBLT returns an Invertible Bloom LookUp Table
// from the block's all transactions
func (b *Block) EncodeToIBLT() *iblt.Graphene {
	//res := iblt.NewInvBloom(DefaultInvBloomConfig)
	res := iblt.NewGraphene(DefaultInvBloomConfig, iblt.DeriveBloomConfig(len(b.GetTransactions())))

	for _, tx := range b.body.Txs {
		res.InsertRLP(tx.CalTxId(), tx)
	}

	return res
}

type Blocks []*Block

type BlockBy func(b1, b2 *Block) bool

func (self BlockBy) Sort(blocks Blocks) {
	bs := blockSorter{
		blocks: blocks,
		by:     self,
	}
	sort.Sort(bs)
}

type blockSorter struct {
	blocks Blocks
	by     func(b1, b2 *Block) bool
}

func (self blockSorter) Len() int { return len(self.blocks) }
func (self blockSorter) Swap(i, j int) {
	self.blocks[i], self.blocks[j] = self.blocks[j], self.blocks[i]
}
func (self blockSorter) Less(i, j int) bool { return self.by(self.blocks[i], self.blocks[j]) }
