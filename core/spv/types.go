package spv

import (
	"fmt"
	"github.com/dipperin/dipperin-core/common"
)

// HexBytes is a type alias to make JSON hex ser/deser easier
type HexBytes []byte

// SPVHeader is a parsed header
type SPVHeader struct {
	Hash   common.Hash `json:"hash"`
	Height uint64      `json:"height"`
	TxRoot common.Hash `json:"tx_root"`
}

// SPVProof is the base struct for an SPV proof
type SPVProof struct {
	Transaction HexBytes  `json:"transaction"`
	Header      SPVHeader `json:"header"`
	Proof       HexBytes  `json:"proof"`
}

// RlpProof is a rlp proof
type RlpProof struct {
	Key   []byte
	Value []byte
}

func (h *SPVHeader) String() string {
	return fmt.Sprintf(`SPVHeader:
	[
		Hash:	    %s
		Height:		%d
		TxRoot		%s
	]`, h.Hash.Hex(), h.Height, h.TxRoot.Hex())
}

func (p *SPVProof) String() string {
	return fmt.Sprintf(`SPVProof:
	[
		Transaction:	    %d
		Header:	%s
		Proof:	        	%d
	]`, p.Transaction, p.Header.String(), p.Proof)
}
