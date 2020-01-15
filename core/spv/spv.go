package spv

import (
	"bytes"
	"errors"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third-party/trie"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
)

var (
	invalidChainID = errors.New("invalid chain id")
	invalidHash    = errors.New("invalid header hash")
	invalidHeight  = errors.New("invalid height")
	invalidTxRoot  = errors.New("invalid tx root")
	invalidProof   = errors.New("invalid proof")
	invalidFrom    = errors.New("invalid from")
	invalidTo      = errors.New("invalid to")
	invalidAmount  = errors.New("invalid amount")
)

// NewSPVHeader builds a Header from a block
func NewSPVHeader(block model.AbstractBlock) SPVHeader {
	return SPVHeader{
		block.ChainID(),
		block.Hash(),
		block.Number(),
		block.TxRoot(),
	}
}

// NewSPVProof builds a proof from a tx and block
func NewSPVProof(tx model.Transaction, block model.AbstractBlock) (SPVProof, error) {
	rlpTx, err := tx.EncodeRlpToBytes()
	if err != nil {
		return SPVProof{}, err
	}

	// calculate MPT rlpProof
	rlpProof, err := getRlpProof(tx.CalTxId(), block.GetTransactions())
	if err != nil {
		return SPVProof{}, err
	}

	spvProof := SPVProof{
		rlpTx,
		NewSPVHeader(block),
		rlpProof,
	}

	return spvProof, nil
}

// getRlpProof get transaction proof
func getRlpProof(tx common.Hash, txs model.Transactions) ([]byte, error) {
	tree := new(trie.Trie)
	for i := 0; i < txs.Len(); i++ {
		tree.Update(txs.GetKey(i), txs.GetRlp(i))
	}

	proof := ethdb.NewMemDatabase()
	err := tree.Prove(tx[:], 0, proof)
	if err != nil {
		return nil, err
	}

	var rlpProof []RlpProof
	keys := proof.Keys()
	for _, key := range keys {
		value, innerErr := proof.Get(key)
		if innerErr != nil {
			return nil, innerErr
		}
		tmp := RlpProof{Key: key, Value: value}
		rlpProof = append(rlpProof, tmp)
	}

	rlpByte, err := rlp.EncodeToBytes(rlpProof)
	if err != nil {
		return nil, err
	}

	return rlpByte, nil
}

func (p SPVProof) validateHeader(header SPVHeader, id, height uint64) error {
	if p.Header.ChainID != id || p.Header.ChainID != header.ChainID {
		return invalidChainID
	}

	if p.Header.Hash != header.Hash {
		return invalidHash
	}

	if p.Header.Height > height || p.Header.Height != header.Height {
		return invalidHeight
	}

	if p.Header.TxRoot != header.TxRoot {
		return invalidTxRoot
	}

	return nil
}

func (p SPVProof) validateTx(tx model.Transaction, from, to common.Address, amount *big.Int) error {
	sender, err := tx.Sender(nil)
	if err != nil {
		return err
	}

	if tx.ChainId().Uint64() != p.Header.ChainID {
		return invalidChainID
	}

	if sender != from {
		return invalidFrom
	}

	if *tx.To() != to {
		return invalidTo
	}

	if tx.Amount().Cmp(amount) != 0 {
		return invalidAmount
	}

	return nil
}

func (p SPVProof) validateProof(txHash common.Hash, proof trie.DatabaseReader) error {
	value, _, err := trie.VerifyProof(p.Header.TxRoot, txHash.Bytes(), proof)
	if err != nil {
		return err
	}

	if !bytes.Equal(value, p.Transaction) {
		return invalidProof
	}

	return nil
}

// Validate checks validity of SPVProof
func (p SPVProof) Validate(header SPVHeader, id, height uint64, from, to common.Address, amount *big.Int) error {
	// Verify block
	err := p.validateHeader(header, id, height)
	if err != nil {
		return err
	}

	// Decode proof
	var rlpProof []RlpProof
	err = rlp.DecodeBytes(p.Proof, &rlpProof)
	if err != nil {
		return err
	}

	// Create tx trie
	tree := ethdb.NewMemDatabase()
	for i := 0; i < len(rlpProof); i++ {
		err = tree.Put(rlpProof[i].Key, rlpProof[i].Value)
		if err != nil {
			return err
		}
	}

	// Decode tx
	var tx model.Transaction
	err = rlp.DecodeBytes(p.Transaction, &tx)
	if err != nil {
		return err
	}

	// Verify proof
	err = p.validateProof(tx.CalTxId(), tree)
	if err != nil {
		return err
	}

	// Verify tx
	return p.validateTx(tx, from, to, amount)
}
