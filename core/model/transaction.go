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
	"container/heap"
	"crypto/ecdsa"
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/g-error"
	"github.com/dipperin/dipperin-core/common/log"
	"github.com/dipperin/dipperin-core/common/math"
	"github.com/dipperin/dipperin-core/core/vm/model"
	"github.com/ethereum/go-ethereum/rlp"
	"go.uber.org/zap"
	"io"
	"math/big"
	"sort"
	"sync/atomic"
)

type Transaction struct {
	data txData
	wit  witness
	// caches
	hash atomic.Value
	size atomic.Value
	from atomic.Value

	//add receipt cache
	receipt atomic.Value
	//add actual tx usedFee
	actualTxFee atomic.Value
}

type RpcTransaction struct {
	To       common.Address
	Value    *big.Int
	GasPrice *big.Int
	GasLimit uint64
	Data     []byte
	Nonce    uint64
}

// IntrinsicGas computes the 'intrinsic gas' for a message with the given data.
func IntrinsicGas(data []byte, contractCreation, homestead bool) (uint64, error) {
	// Set the starting gas for the raw transaction
	var gas uint64
	if contractCreation && homestead {
		gas = model.TxGasContractCreation
	} else {
		gas = model.TxGas
	}
	// Bump the required gas by the amount of transactional data
	if len(data) > 0 {
		// Zero and non-zero bytes are priced differently
		var nz uint64
		for _, byt := range data {
			if byt != 0 {
				nz++
			}
		}
		// Make sure we don't exceed uint64 for all data combinations
		if (math.MaxUint64-gas)/model.TxDataNonZeroGas < nz {
			return 0, g_error.ErrOutOfGas
		}
		gas += nz * model.TxDataNonZeroGas

		z := uint64(len(data)) - nz
		if (math.MaxUint64-gas)/model.TxDataZeroGas < z {
			return 0, g_error.ErrOutOfGas
		}
		gas += z * model.TxDataZeroGas
	}
	return gas, nil
}

func NewTransaction(nonce uint64, to common.Address, amount, gasPrice *big.Int, gasLimit uint64, data []byte) *Transaction {
	return newTransaction(nonce, &to, amount, gasPrice, gasLimit, data)
}

func NewContractCreation(nonce uint64, amount *big.Int, gasPrice *big.Int, gasLimit uint64, data []byte) *Transaction {
	return newTransaction(nonce, nil, amount, gasPrice, gasLimit, data)
}

func newTransaction(nonce uint64, to *common.Address, amount, gasPrice *big.Int, gasLimit uint64, data []byte) *Transaction {
	if len(data) > 0 {
		data = common.CopyBytes(data)
	}
	txdata := txData{
		AccountNonce: nonce,
		Recipient:    to,
		HashLock:     nil,
		TimeLock:     new(big.Int),
		Amount:       new(big.Int),
		//Fee:          new(big.Int),
		Price:     gasPrice,
		GasLimit:  gasLimit,
		ExtraData: data,
	}
	wit := witness{
		R:       new(big.Int),
		S:       new(big.Int),
		V:       new(big.Int),
		HashKey: nil,
	}
	if amount != nil {
		txdata.Amount.Set(amount)
	}

	return &Transaction{data: txdata, wit: wit}
}

func (tx Transaction) IsEqual(tempTx Transaction) bool {

	return tx.CalTxId().IsEqual(tempTx.CalTxId())
}

// Whether the transaction is a cross-chain transaction
//func (tx Transaction) IsCrossChain() bool {
//	// TODO: Cross chain criteria
//	hasTimelock := tx.data.TimeLock != nil
//	if hasTimelock {
//		// when TimeLock is not nil, Timelock not equal to zero
//		hasTimelock = tx.data.TimeLock.Cmp(big.NewInt(0)) != 0
//		//fmt.Println("hasTimelock", tx.data.TimeLock)
//	}
//
//	hasHashlock := tx.data.HashLock != nil
//	if hasHashlock {
//		// when Hashlock is not nil, len of Hashlock not equal to zero
//		hasHashlock = len(tx.data.HashLock) != 0
//		//fmt.Println("hasHashlock", tx.data.HashLock)
//	}
//
//	hasHashKey := tx.wit.HashKey != nil
//	if hasHashKey {
//		// when HashKey is not nil, len of HashKey not equal to zero
//		hasHashKey = len(tx.wit.HashKey) != 0
//		//fmt.Println("hasHashKey", tx.wit.HashKey)
//	}
//	return hasTimelock || hasHashlock || hasHashKey
//}

func (tx Transaction) GetType() common.TxType {
	return tx.data.Recipient.GetAddressType()
}

func (tx Transaction) String() string {
	//return "data:" + util.StringifyJson(tx.data) + " wit:" + util.StringifyJson(tx.wit)
	var from, to string

	if tx.wit.V != nil {
		//todo  this getsinger method need to be implement later.
		signer := tx.GetSigner()
		if f, err := tx.Sender(signer); err != nil {
			from = fmt.Sprintf("%s", err)
		} else {
			from = fmt.Sprintf("%x", f[:])
		}
	} else {
		from = "[invalid sender: nil V field]"
	}

	if tx.data.Recipient == nil {
		to = "[contract creation]"
	} else {
		to = fmt.Sprintf("%x", tx.data.Recipient[:])
	}
	return fmt.Sprintf(`
	TXID:	  %s
	Type:     %s
	From:     0x%s
	To:       0x%s
	Nonce:    %v
	GasPrice: %v
	GasLimit: %v
	Hashlock: %v
	Timelock: %#x
	Value:    %d CSC
	Data:     0x%x
	V:        %#x
	R:        %#x
	S:        %#x
	HashKey:  0x%x    
`,
		tx.CalTxId().Hex(),
		tx.data.Recipient.GetAddressTypeStr(),
		from,
		to,
		tx.data.AccountNonce,
		tx.data.Price,
		tx.data.GasLimit,
		tx.data.HashLock,
		tx.data.TimeLock,
		tx.data.Amount,
		tx.data.ExtraData,
		tx.wit.V,
		tx.wit.R,
		tx.wit.S,
		tx.wit.HashKey,
	)
}

type txData struct {
	AccountNonce uint64          `json:"nonce"    gencodec:"required"`
	Recipient    *common.Address `json:"to"       rlp:"nil"`
	HashLock     *common.Hash    `json:"hashLock" rlp:"nil"`
	TimeLock     *big.Int        `json:"timeLock22" gencodec:"required"`
	Amount       *big.Int        `json:"Value"    gencodec:"required"`
	Price        *big.Int        `json:"gasPrice" gencodec:"required"`
	GasLimit     uint64          `json:"gas"      gencodec:"required"`
	ExtraData    []byte          `json:"input"    gencodec:"required"`
}

type witness struct {
	// Signature values
	R *big.Int `json:"r" gencodec:"required"`
	S *big.Int `json:"s" gencodec:"required"`
	V *big.Int `json:"v" gencodec:"required"`
	// hash_key
	HashKey []byte `json:"hashKey"    gencodec:"required"`
}

type TransactionRLP struct {
	Txdata txData
	Wit    witness
}

//EncodeRLP implements rlp.Encoder
func (tx *Transaction) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, TransactionRLP{
		tx.data,
		tx.wit,
	})
}

func (tx *Transaction) EncodeRlpToBytes() ([]byte, error) {
	return rlp.EncodeToBytes(TransactionRLP{
		tx.data,
		tx.wit,
	})
}

func (tx *Transaction) GetGasPrice() *big.Int {
	return tx.data.Price
}

func (tx *Transaction) GetGasLimit() uint64 {
	return tx.data.GasLimit
}

//DecodeRLP implements rlp.Decoder
func (tx *Transaction) DecodeRLP(s *rlp.Stream) error {
	var dtx TransactionRLP
	_, size, _ := s.Kind()
	if err := s.Decode(&dtx); err != nil {
		return err
	}
	tx.data, tx.wit = dtx.Txdata, dtx.Wit
	tx.size.Store(common.StorageSize(rlp.ListSize(size)))
	return nil
}

func (tx *Transaction) Nonce() uint64 { return tx.data.AccountNonce }

//func (tx *Transaction) Version() uint64        { return tx.data.Version }
func (tx *Transaction) HashLock() *common.Hash { return tx.data.HashLock }
func (tx *Transaction) TimeLock() *big.Int     { return new(big.Int).Set(tx.data.TimeLock) }
func (tx *Transaction) ExtraData() []byte      { return tx.data.ExtraData }
func (tx *Transaction) Amount() *big.Int       { return new(big.Int).Set(tx.data.Amount) }

//func (tx *Transaction) Fee() *big.Int          { return new(big.Int).Set(tx.data.Fee) }
func (tx *Transaction) RawSignatureValues() (*big.Int, *big.Int, *big.Int) {
	return new(big.Int).Set(tx.wit.V), new(big.Int).Set(tx.wit.R), new(big.Int).Set(tx.wit.S)
}

func (tx *Transaction) HashKey() []byte {
	return tx.wit.HashKey
}

func (tx *Transaction) To() *common.Address {
	if tx.data.Recipient == nil {
		return nil
	}
	to := *tx.data.Recipient
	return &to
}

func (tx *Transaction) Sender(signer Signer) (common.Address, error) {
	if signer == nil {
		signer = tx.GetSigner()
	}

	if sc := tx.from.Load(); sc != nil {
		sigCache := sc.(sigCache)
		// If the signer used to derive from in a previous
		// call is not the same as used current, invalidate
		// the cache.
		//fmt.Println("the tx sender cache")
		//fmt.Println(sigCache.signer, signer)
		if sigCache.signer.Equal(signer) {
			//fmt.Println("+++++++++++++++++++++++++++++==sigCache signer match==")
			return sigCache.from, nil
		}
	}

	addr, err := signer.GetSender(tx)
	if err != nil {
		return common.Address{}, err
	}

	tx.from.Store(sigCache{signer: signer, from: addr})
	return addr, nil
}

func (tx *Transaction) SenderPublicKey(signer Signer) (*ecdsa.PublicKey, error) {
	if signer == nil {
		signer = tx.GetSigner()
	}
	pubKey, err := signer.GetSenderPublicKey(tx)
	if err != nil {
		return pubKey, err
	}
	return pubKey, nil

}

func (tx *Transaction) Size() common.StorageSize {
	if size := tx.size.Load(); size != nil {
		return size.(common.StorageSize)
	}
	c := writeCounter(0)
	if err := rlp.Encode(&c, tx); err != nil {
		panic(err)
	}
	tx.size.Store(common.StorageSize(c))
	return common.StorageSize(c)
}

// get the transaction ID
func (tx *Transaction) CalTxId() common.Hash {
	if h := tx.hash.Load(); h != nil {
		//fmt.Println("Transaction#CalTxId", "get load hash", h)
		return h.(common.Hash)
	}

	//get tx address
	address, err := tx.Sender(nil)
	if err != nil {
		return common.Hash{}
	}

	//calculate TxId = hash(tx.data+address)
	txId, err := rlpHash([]interface{}{tx.data, address})
	if err != nil {
		return common.Hash{}
	}

	tx.hash.Store(txId)

	return txId
}

func (tx Transaction) ChainId() *big.Int {
	return deriveChainId(tx.wit.V)
}

//todo: currently use default signer ,later need get a way to determine the signer from the tx itself.
func (tx Transaction) GetSigner() Signer {
	id := deriveChainId(tx.wit.V)

	return DipperinSigner{id}
}

// Cost returns amount + fee
func (tx *Transaction) Cost() *big.Int {
	total := new(big.Int).Mul(tx.data.Price, new(big.Int).SetUint64(tx.data.GasLimit))
	total.Add(total, tx.data.Amount)
	return total
}

func (tx *Transaction) EstimateFee() *big.Int {
	encoded, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return big.NewInt(0)
	}
	size := len(encoded)
	rate := big.NewInt(int64(1))
	fee := big.NewInt(0).Mul(big.NewInt(int64(size)), rate)
	return fee
}

func (tx *Transaction) AsMessage(checkNonce bool) (Message, error) {
	msg := Message{
		nonce:      tx.data.AccountNonce,
		gasLimit:   tx.data.GasLimit,
		gasPrice:   new(big.Int).Set(tx.data.Price),
		to:         tx.data.Recipient,
		amount:     tx.data.Amount,
		data:       tx.data.ExtraData,
		checkNonce: checkNonce,
	}

	var err error
	msg.from, err = tx.Sender(tx.GetSigner())
	return msg, err
}

type ReceiptPara struct {
	Root              []byte
	HandlerResult     bool
	CumulativeGasUsed uint64
	Logs              []*model.Log
}

func (tx *Transaction) PaddingActualTxFee(fee *big.Int) {
	tx.actualTxFee.Store(fee)
}

func (tx *Transaction) GetActualTxFee() (fee *big.Int) {
	if feeLoad := tx.actualTxFee.Load(); feeLoad != nil {
		return feeLoad.(*big.Int)
	}
	log.DLogger.Error("the transaction fee cache is nil")
	return nil
}

func (tx *Transaction) PaddingReceipt(parameters ReceiptPara) {
	log.DLogger.Info("Call PaddingReceipt", zap.Bool("handlerResult", parameters.HandlerResult))
	receipt := model.NewReceipt(parameters.Root, parameters.HandlerResult, parameters.CumulativeGasUsed, parameters.Logs)
	tx.receipt.Store(receipt)
}

func (tx *Transaction) GetReceipt() *model.Receipt {
	if receiptLoad := tx.receipt.Load(); receiptLoad != nil {
		return receiptLoad.(*model.Receipt)
	}
	log.DLogger.Error("the receipt cache is nil")
	return nil
}

// Transactions is a Transaction slice type for basic sorting.
type Transactions []*Transaction

func (ss Transactions) GetKey(i int) []byte {
	res := ss[i].CalTxId().Bytes()
	return res

}

// GetRlp implements Rlpable and returns the i'th element of s in rlp.
func (ss Transactions) GetRlp(i int) []byte {
	enc, _ := rlp.EncodeToBytes(ss[i])
	return enc
}

// Len returns the length of s.
func (ss Transactions) Len() int { return len(ss) }

// Swap swaps the i'th and the j'th element in s.
func (ss Transactions) Swap(i, j int) { ss[i], ss[j] = ss[j], ss[i] }

func (ss Transactions) String() string {
	var out string
	for _, s := range ss {
		out += fmt.Sprintf("%v/n", s)
	}
	return out
}

func (ss Transactions) Less(i, j int) bool { return ss[i].CalTxId().Cmp(ss[j].CalTxId()) == -1 }

type TransactionBy func(b1, b2 *Block) bool

func (self TransactionBy) Sort(blocks Blocks) {
	ts := txSorter{
		blocks: blocks,
		by:     self,
	}
	sort.Sort(ts)
}

type txSorter struct {
	blocks Blocks
	by     func(b1, b2 *Block) bool
}

func (self txSorter) Len() int { return len(self.blocks) }
func (self txSorter) Swap(i, j int) {
	self.blocks[i], self.blocks[j] = self.blocks[j], self.blocks[i]
}
func (self txSorter) Less(i, j int) bool { return self.by(self.blocks[i], self.blocks[j]) }

// TxByPrice implements both the sort and the heap interface, making it useful
// for all at once sorting as well as individually adding and removing elements.
type TxByFee []AbstractTransaction

func (s TxByFee) Len() int           { return len(s) }
func (s TxByFee) Less(i, j int) bool { return s[i].GetGasPrice().Cmp(s[j].GetGasPrice()) > 0 }
func (s TxByFee) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func (s *TxByFee) Push(x interface{}) {
	*s = append(*s, x.(*Transaction))
}

func (s *TxByFee) Pop() interface{} {
	old := *s
	n := len(old)
	x := old[n-1]
	*s = old[0 : n-1]
	return x
}

// TransactionsByPriceAndNonce represents a set of transactions that can return
// transactions in a profit-maximizing sorted order, while supporting removing
// entire batches of transactions for non-executable accounts.
type TransactionsByFeeAndNonce struct {
	txs    map[common.Address][]AbstractTransaction // Per account nonce-sorted list of transactions
	heads  TxByFee                                  // NextMiss transaction for each unique account (price heap)
	signer Signer                                   // Signer for the set of transactions
}

// NewTransactionsByPriceAndNonce creates a transaction set that can retrieve
// price sorted transactions in a nonce-honouring way.
//
// Note, the input map is reowned so the caller should not interact any more with
// if after providing it to the constructor.
// TODO  errors will occur if all transactions are deleted
func NewTransactionsByFeeAndNonce(signer Signer, txs map[common.Address][]AbstractTransaction) *TransactionsByFeeAndNonce {
	// Initialize a price based heap with the head transactions
	heads := make(TxByFee, 0, len(txs))
	for from, accTxs := range txs {
		log.DLogger.Info("NewTransactionsByFeeAndNonce ", zap.Any("from", from), zap.Int("len(heads)", len(heads)))
		// 此处 ethereum　的写法,假设from != acc这种异常情况出现，txs map会被新增acc字段交易或将原acc字段替换成from的相关交易
		//　导致txs 异常．此外会导致range 不确定性，修改的acc　txs有可能遍历到，也有可能遍历不到
		//　因此统一修改逻辑为:当出现此异常时，将此from的txs直接删除，heads里也不处理此类交易．
		/*		heads = append(heads, accTxs[0])
				// Ensure the sender address is from the signer
				acc, _ := accTxs[0].Sender(signer)
				txs[acc] = accTxs[1:]
				if from != acc {
					delete(txs, from)
				}*/
		if len(accTxs) == 0 {
			log.DLogger.Warn("theaccTxs is nil")
			delete(txs, from)
		} else {
			acc, _ := accTxs[0].Sender(signer)
			if from != acc {
				log.DLogger.Warn("the tx sender and from is different")
				delete(txs, from)
			} else {
				heads = append(heads, accTxs[0])
				txs[acc] = accTxs[1:]
			}
		}
	}
	heap.Init(&heads)

	// Assemble and return the transaction set
	return &TransactionsByFeeAndNonce{
		txs:    txs,
		heads:  heads,
		signer: signer,
	}
}

// Peek returns the next transaction by price.
func (t *TransactionsByFeeAndNonce) Peek() AbstractTransaction {
	if len(t.heads) == 0 {
		return nil
	}
	return t.heads[0]
}

// Shift replaces the current best head with the next one from the same account.
func (t *TransactionsByFeeAndNonce) Shift() {
	acc, _ := t.heads[0].Sender(t.signer)
	if txs, ok := t.txs[acc]; ok && len(txs) > 0 {
		t.heads[0], t.txs[acc] = txs[0], txs[1:]
		heap.Fix(&t.heads, 0)
	} else {
		heap.Pop(&t.heads)
	}
}

// Pop removes the best transaction, *not* replacing it with the next one from
// the same account. This should be used when a transaction cannot be executed
// and hence all subsequent ones should be discarded from the same account.
func (t *TransactionsByFeeAndNonce) Pop() {
	heap.Pop(&t.heads)
}

type Message struct {
	to         *common.Address
	from       common.Address
	nonce      uint64
	amount     *big.Int
	gasLimit   uint64
	gasPrice   *big.Int
	data       []byte
	checkNonce bool
}

/*func NewMessage(from common.Address, to *common.Address, nonce uint64, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte, checkNonce bool) Message {
	return Message{
		from:       from,
		to:         to,
		nonce:      nonce,
		amount:     amount,
		gasLimit:   gasLimit,
		gasPrice:   gasPrice,
		data:       data,
		checkNonce: checkNonce,
	}
}*/

func (m Message) From() common.Address { return m.from }
func (m Message) To() *common.Address  { return m.to }
func (m Message) GasPrice() *big.Int   { return m.gasPrice }
func (m Message) Value() *big.Int      { return m.amount }
func (m Message) Gas() uint64          { return m.gasLimit }
func (m Message) Nonce() uint64        { return m.nonce }
func (m Message) Data() []byte         { return m.data }
func (m Message) CheckNonce() bool     { return m.checkNonce }
func (m *Message) SetGas(gas uint64) {
	m.gasLimit = gas
}
