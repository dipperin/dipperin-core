package spv

import (
	"bytes"
	"errors"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third-party/trie"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
)

var invalidProof = errors.New("invalid proof")

// newSPVHeader builds a Header from a block
func newSPVHeader(block model.AbstractBlock) SPVHeader {
	return SPVHeader{
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
		newSPVHeader(block),
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

// Validate checks validity of SPVProof
func (p SPVProof) Validate() error {
	// Decode proof
	var rlpProof []RlpProof
	err := rlp.DecodeBytes(p.Proof, &rlpProof)
	if err != nil {
		return err
	}

	data := ethdb.NewMemDatabase()
	for i := 0; i < len(rlpProof); i++ {
		err = data.Put(rlpProof[i].Key, rlpProof[i].Value)
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

	value, _, err := trie.VerifyProof(p.Header.TxRoot, tx.CalTxId().Bytes(), data)
	if err != nil {
		return err
	}

	if !bytes.Equal(value, p.Transaction) {
		return invalidProof
	}

	return err
}
